// Package plan turns a selection into a concrete, ordered, previewable plan.
package plan

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"meuubuntu/internal/aptsrc"
	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/sys"
)

type Step struct {
	ID       string                `json:"id"`
	Item     string                `json:"item"`
	Title    string                `json:"title"`
	Kind     string                `json:"kind"` // admin | user
	Op       string                `json:"op"`
	Options  map[string][]string   `json:"options,omitempty"`
	Detail   []string              `json:"detail"`
	Changes  []detect.SettingState `json:"changes,omitempty"`
	Category string                `json:"category"`
	Confirm  string                `json:"confirm,omitempty"` // must be confirmed explicitly before running
	Expect   string                `json:"expect,omitempty"`  // e.g. SHA-256 of the reviewed script
	Undo     []UndoValue           `json:"undo,omitempty"`
}

type UndoValue struct {
	Kind       string `json:"kind"`
	Schema     string `json:"schema,omitempty"`
	Key        string `json:"key"`
	Current    string `json:"current"`
	Target     string `json:"target"`
	ChangeID   string `json:"changeId,omitempty"`
	Mode       string `json:"mode"` // journal | reset | skip
	StateToken string `json:"stateToken,omitempty"`
}

type Conflict struct {
	A, B   string
	Hard   bool
	Reason string
}

type Skipped struct {
	Item   string `json:"item"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
	Status string `json:"status"`
}

type ManualStep struct {
	Item  string   `json:"item"`
	Name  string   `json:"name"`
	Steps []string `json:"steps"`
	URL   string   `json:"url"`
}

type ConflictView struct {
	A      string `json:"a"`
	B      string `json:"b"`
	Hard   bool   `json:"hard"`
	Reason string `json:"reason"`
}

type Plan struct {
	ReviewToken string                       `json:"reviewToken"`
	Steps       []Step                       `json:"steps"`
	Conflicts   []ConflictView               `json:"conflicts"`
	Blocked     bool                         `json:"blocked"`
	Added       []string                     `json:"added"`
	Skipped     []Skipped                    `json:"skipped"`
	Manual      []ManualStep                 `json:"manual"`
	Relogin     bool                         `json:"relogin"`
	Reboot      bool                         `json:"reboot"`
	NeedsRoot   bool                         `json:"needsRoot"`
	Warnings    []string                     `json:"warnings"`
	Selection   map[string]catalog.Selection `json:"selection"`
	Reverts     []string                     `json:"reverts"`
	Kept        []Skipped                    `json:"kept"` // revert requested but not possible automatically
}

// Fingerprint binds execution to the reviewed steps, options and current-value
// details. The GUI sends only this token, never commands or a serialized plan.
func (p *Plan) Fingerprint() string {
	copy := *p
	copy.ReviewToken = ""
	// Added only explains how dependencies entered the selection. Rebuilding
	// the reviewed, expanded selection must retain the same execution token.
	copy.Added = nil
	b, _ := json.Marshal(copy)
	return fmt.Sprintf("%x", sha256.Sum256(b))
}

var rank = map[string]int{
	catalog.MethodAPTRepo: 1, catalog.MethodAPT: 2, catalog.MethodAdmin: 3, catalog.MethodFlatpak: 4,
	catalog.MethodExtension: 5, catalog.MethodGSettings: 6, catalog.MethodGitConfig: 7,
	catalog.MethodTemplate: 8, catalog.MethodPowerProfile: 9, catalog.MethodManual: 10, catalog.MethodUserScript: 4,
	catalog.MethodNPMUser: 4,
}

// Build validates the selection and produces a plan using the real state.
// Build plans the selected items plus the explicit reverts the user asked for.
// Unselected items are never touched.
func Build(ctx context.Context, cat *catalog.Catalog, d *detect.Detector, sel map[string]catalog.Selection, reverts []string) (*Plan, error) {
	d.Refresh(ctx)
	p := &Plan{Steps: []Step{}, Selection: map[string]catalog.Selection{}}
	for id, s := range sel {
		v, err := cat.ValidateSelection(id, s)
		if err != nil {
			return nil, err
		}
		p.Selection[id] = v
	}
	// Dependencies.
	var addDeps func(id string)
	addDeps = func(id string) {
		it := cat.Item(id)
		if status := d.State(ctx, it, p.Selection[id]).Status; status == detect.Installed || status == detect.Provided {
			return // An already satisfied choice must not install missing optional prerequisites.
		}
		for _, r := range cat.Item(id).Requires {
			if _, ok := p.Selection[r]; !ok {
				v, _ := cat.ValidateSelection(r, catalog.Selection{})
				p.Selection[r] = v
				p.Added = append(p.Added, r)
				addDeps(r)
			}
		}
	}
	ids := sortedIDs(p.Selection)
	for _, id := range ids {
		addDeps(id)
	}
	ids = sortedIDs(p.Selection)

	states := map[string]detect.ItemState{}
	for _, id := range ids {
		states[id] = d.State(ctx, cat.Item(id), p.Selection[id])
	}

	// Conflicts: between selected items, against installed items, and groups.
	seen := map[string]bool{}
	addConf := func(a, b string, hard bool, reason string) {
		k := a + "|" + b
		if a > b {
			k = b + "|" + a
		}
		if seen[k] {
			return
		}
		seen[k] = true
		p.Conflicts = append(p.Conflicts, ConflictView{A: a, B: b, Hard: hard, Reason: reason})
		if hard {
			p.Blocked = true
		}
	}
	for _, id := range ids {
		it := cat.Item(id)
		for _, c := range it.Conflicts {
			if _, both := p.Selection[c.ID]; both {
				addConf(id, c.ID, c.Hard, c.Reason)
				continue
			}
			other := cat.Item(c.ID)
			if st := d.State(ctx, other, catalog.Selection{Options: defaults(cat, other)}); st.Status == detect.Installed && other.Method != catalog.MethodGSettings && other.Method != catalog.MethodPowerProfile {
				addConf(id, c.ID, c.Hard, c.Reason+" ("+other.Name+" já está instalado no sistema)")
			}
		}
		if it.Group != "" {
			for _, j := range ids {
				if j != id && cat.Item(j).Group == it.Group {
					addConf(id, j, true, "Alternativas exclusivas ("+it.Group+"): escolha apenas uma")
				}
			}
		}
	}

	sort.SliceStable(ids, func(i, j int) bool {
		a, b := cat.Item(ids[i]), cat.Item(ids[j])
		if rank[a.Method] != rank[b.Method] {
			return rank[a.Method] < rank[b.Method]
		}
		return a.ID < b.ID
	})
	// Keep stable method ordering where possible, but always visit selected
	// prerequisites before their dependants (e.g. PHP before Composer).
	var ordered []string
	visited := map[string]bool{}
	var visit func(string)
	visit = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true
		for _, dep := range cat.Item(id).Requires {
			if _, ok := p.Selection[dep]; ok {
				visit(dep)
			}
		}
		ordered = append(ordered, id)
	}
	for _, id := range ids {
		visit(id)
	}
	ids = ordered

	n := 0
	step := func(s Step) {
		n++
		s.ID = fmt.Sprintf("s%02d", n)
		if s.Kind == "admin" {
			p.NeedsRoot = true
		}
		p.Steps = append(p.Steps, s)
	}
	needUpdate := false
	flathubDone := false
	for _, id := range ids {
		it := cat.Item(id)
		s := p.Selection[id]
		st := states[id]
		switch st.Status {
		case detect.Installed, detect.Provided:
			p.Skipped = append(p.Skipped, Skipped{id, it.Name, st.Detail, st.Status})
			continue
		case detect.Incompatible, detect.Unknown:
			p.Skipped = append(p.Skipped, Skipped{id, it.Name, st.Detail, st.Status})
			continue
		}
		if it.Method == catalog.MethodManual {
			m := ManualStep{Item: id, Name: it.Name}
			if it.Manual != nil {
				m.Steps, m.URL = it.Manual.Steps, it.Manual.URL
			}
			p.Manual = append(p.Manual, m)
			continue
		}
		p.Relogin = p.Relogin || it.Relogin
		p.Reboot = p.Reboot || it.Reboot
		if it.Service {
			p.Warnings = append(p.Warnings, it.Name+": instala um serviço local que pode arrancar automaticamente.")
		}
		base := Step{Item: id, Category: it.Category, Options: s.Options, Confirm: it.Confirm}
		switch it.Method {
		case catalog.MethodAPTRepo:
			r := cat.Repo(it.Repo)
			x := base
			x.Kind, x.Op, x.Title = "admin", "repo-add", "Adicionar repositório oficial: "+r.Name
			x.Detail = []string{"Chave: " + r.KeyURL, "Impressões digitais exigidas: " + strings.Join(r.Fingerprints, ", "), "Fonte: " + r.URL + " (" + strings.ReplaceAll(r.Suite, "{codename}", d.Sys.Codename) + ")", "Ficheiro: /etc/apt/sources.list.d/meu-ubuntu-" + r.ID + ".sources"}
			if !r.Verified() {
				p.Warnings = append(p.Warnings, r.Name+": impressão digital não verificada — este passo vai falhar de propósito.")
			}
			step(x)
			fallthrough
		case catalog.MethodAPT:
			present := d.PackagesPresent(it, s)
			if !present && it.Method == catalog.MethodAPT && !needUpdate {
				needUpdate = true
				step(Step{Kind: "admin", Op: "apt-update", Title: "Atualizar índices APT", Detail: []string{sys.Argv("apt-get", "update")}, Category: "Sistema"})
			}
			x := base
			pk := cat.ExpandPackages(it, s)
			x.Kind, x.Op, x.Title = "admin", "apt-install", "Instalar "+it.Name
			x.Detail = []string{sys.Argv("apt-get", append([]string{"install", "-y", "--no-remove"}, pk...)...)}
			if !present {
				step(x)
			}
			if len(it.Post) > 0 {
				y := base
				y.Kind, y.Op, y.Title = "admin", "post", "Configurar "+it.Name
				y.Detail = append([]string{}, it.Post...)
				if contains(it.Post, "xdebug-off") {
					y.Detail = []string{"Xdebug instalado mas desativado (phpdismod xdebug); ative por projeto quando precisar"}
				}
				step(y)
			} else if present {
				step(Step{Item: id, Title: "Verificar " + it.Name, Kind: "user", Op: "verify", Category: it.Category, Options: s.Options, Detail: []string{"Repetir a verificação sem reinstalar pacotes"}})
			}
		case catalog.MethodAdmin:
			x := base
			x.Kind, x.Op, x.Title = "admin", "admin", it.Name
			x.Changes = st.Settings
			x.Detail = []string{it.Summary}
			if v := s.First("value"); v != "" {
				x.Detail = append(x.Detail, "Valor: "+v)
			}
			for _, setting := range st.Settings {
				x.Detail = append(x.Detail, setting.Key+":\nAnterior:\n"+setting.Current+"\nProposto:\n"+setting.Proposed)
			}
			step(x)
		case catalog.MethodFlatpak:
			if !flathubDone {
				flathubDone = true
				step(Step{Kind: "user", Op: "flatpak-remote", Title: "Garantir Flathub (utilizador)", Category: "Sistema",
					Detail: []string{sys.Argv("flatpak", "remote-add", "--user", "--if-not-exists", "flathub", "https://dl.flathub.org/repo/flathub.flatpakrepo")}})
			}
			x := base
			x.Kind, x.Op, x.Title = "user", "flatpak-install", "Instalar "+it.Name+" (Flatpak)"
			x.Detail = []string{sys.Argv("flatpak", "install", "--user", "-y", "--noninteractive", "flathub", it.Flatpak)}
			step(x)
		case catalog.MethodExtension:
			x := base
			x.Kind, x.Op, x.Title = "user", "extension", "Instalar extensão "+it.Name
			x.Detail = []string{fmt.Sprintf("extensions.gnome.org: %s para GNOME %d (verificado antes de transferir)", it.Extension, d.Sys.GNOMEMaj), "Ativação após terminar sessão (Wayland)"}
			if st.Status == detect.NeedsConfig {
				x.Op, x.Title, x.Detail = "extension-enable", "Ativar extensão "+it.Name, []string{"Extensão já instalada; ativar " + it.Extension}
			}
			step(x)
		case catalog.MethodGSettings, catalog.MethodGitConfig, catalog.MethodPowerProfile:
			x := base
			x.Kind, x.Op, x.Title = "user", it.Method, it.Name
			x.Changes = st.Settings
			for _, ss := range st.Settings {
				if ss.Current != ss.Proposed {
					cur := ss.Current
					if cur == "" {
						cur = "(sem valor)"
					}
					x.Detail = append(x.Detail, fmt.Sprintf("%s %s: %s → %s", ss.Schema, ss.Key, cur, ss.Proposed))
				}
			}
			step(x)
		case catalog.MethodUserScript:
			x := base
			x.Kind, x.Op, x.Title = "user", "user-script", "Instalar "+it.Name+" (sem root)"
			x.Detail = []string{"Transferir " + ScriptURL(it, s) + " para ficheiro (sem curl|bash)",
				"Recusar se contiver sudo ou hosts fora de: " + strings.Join(it.Script.AllowHosts, ", "),
				"Backup antes: " + strings.Join(it.Script.Backup, ", "),
				"Executar como utilizador: bash <ficheiro>; SHA-256 registado no log",
				"Atenção: o instalador do fornecedor não verifica checksums dos binários"}
			step(x)
		case catalog.MethodNPMUser:
			x := base
			x.Kind, x.Op, x.Title = "user", "npm-user", "Instalar "+it.Name+" para o seu utilizador"
			x.Detail = []string{"Pacote oficial: " + it.NPM.Package + "@" + it.NPM.Version,
				"Instalar em ~/.local, sem sudo e sem alterar o npm global do sistema",
				"Verificar o comando " + it.NPM.Binary + " --version depois da instalação"}
			step(x)
		case catalog.MethodTemplate:
			x := base
			x.Kind, x.Op, x.Title = "user", "template", "Modelo: "+it.Template.Name
			x.Detail = []string{st.Detail + " (nunca sobrescreve)"}
			step(x)
		}
	}
	p.undoSteps(ctx, cat, d, reverts, step)
	// Before any APT work, repair duplicated sources left by vendor packages.
	for i, x := range p.Steps {
		if x.Op == "apt-update" || x.Op == "apt-install" || x.Op == "repo-add" {
			chk := Step{Kind: "admin", Op: "sources-check", Title: "Verificar fontes APT (duplicados)", Category: "Sistema"}
			if dups := aptsrc.Duplicates(d.Root, cat); len(dups) > 0 {
				for _, x := range dups {
					chk.Title = "Limpar fontes APT duplicadas"
					files := x.Ours
					if x.Keyring != "" {
						files += " e " + x.Keyring
					}
					chk.Detail = append(chk.Detail, fmt.Sprintf("Remover %s: o repositório %s já está configurado em %s (backup antes em /var/backups/meu-ubuntu; chaves usadas por outras fontes são preservadas)", files, x.Name, strings.Join(x.Others, ", ")))
				}
			} else {
				chk.Detail = []string{"Nenhum duplicado agora. Volta a verificar durante a execução, porque alguns pacotes (ex.: warp-terminal) criam a sua própria fonte ao instalar; só ficheiros do Meu Ubuntu são removidos, com backup."}
			}
			p.Steps = append(p.Steps[:i], append([]Step{chk}, p.Steps[i:]...)...)
			for j := range p.Steps {
				p.Steps[j].ID = fmt.Sprintf("s%02d", j+1)
			}
			break
		}
	}
	return p, nil
}

// undoSteps reverts applied items the user explicitly asked to revert.
func (p *Plan) undoSteps(ctx context.Context, cat *catalog.Catalog, d *detect.Detector, reverts []string, step func(Step)) {
	for _, id := range reverts {
		it := cat.Item(id)
		if it == nil || it.Group != "" {
			continue
		}
		if _, sel := p.Selection[id]; sel {
			continue
		}
		v, _ := cat.ValidateSelection(id, catalog.Selection{})
		st := d.State(ctx, it, v)
		reversible := it.Method == catalog.MethodGSettings || it.Method == catalog.MethodGitConfig || it.Method == catalog.MethodTemplate || (it.Method == catalog.MethodAdmin && (it.AdminOp == "inotify" || it.AdminOp == "journald"))
		if st.Status != detect.Installed && (!reversible || st.Status == detect.Incompatible || st.Status == detect.Unknown || st.Status == detect.NeedsConfig) {
			continue
		}
		p.Reverts = append(p.Reverts, id)
		base := Step{Item: id, Category: it.Category, Title: "Reverter: " + it.Name}
		switch {
		case it.Method == catalog.MethodGSettings:
			base.Kind, base.Op = "user", "undo"
			for _, s := range it.Settings {
				base.Detail = append(base.Detail, s.Schema+" "+s.Key+": valor anterior registado, ou padrão do sistema")
			}
			step(base)
		case it.Method == catalog.MethodGitConfig || it.Method == catalog.MethodTemplate:
			base.Kind, base.Op = "user", "undo"
			base.Detail = []string{"Só desfaz o que o Meu Ubuntu registou; o resto mantém-se"}
			step(base)
		case it.Method == catalog.MethodAdmin && (it.AdminOp == "docker-group" || it.AdminOp == "inotify" || it.AdminOp == "journald"):
			base.Kind, base.Op = "admin", "revert"
			base.Changes = st.Settings
			base.Detail = []string{"Reverte a operação (só ficheiros escritos pelo Meu Ubuntu)"}
			if it.Relogin {
				p.Relogin = true
			}
			step(base)
		default:
			p.Kept = append(p.Kept, Skipped{id, it.Name, "Sem reversão automática (pacotes/aplicações não são removidos pelo Meu Ubuntu)", st.Status})
		}
	}
}

func defaults(cat *catalog.Catalog, it *catalog.Item) map[string][]string {
	s, _ := cat.ValidateSelection(it.ID, catalog.Selection{})
	return s.Options
}

func sortedIDs(m map[string]catalog.Selection) []string {
	var k []string
	for x := range m {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}

func contains(a []string, s string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

// ScriptURL resolves {opt:key} with a validated option value.
func ScriptURL(it *catalog.Item, s catalog.Selection) string {
	u := it.Script.URL
	if i := strings.Index(u, "{opt:"); i >= 0 {
		j := strings.Index(u[i:], "}")
		u = u[:i] + s.First(u[i+5:i+j]) + u[i+j+1:]
	}
	return u
}
