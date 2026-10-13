// Package detect computes the real state of catalog items (read-only).
package detect

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"

	"meuubuntu/internal/adminconfig"
	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

const (
	Installed    = "installed"
	Available    = "available"
	NeedsConfig  = "config"
	Incompatible = "incompatible"
	Unknown      = "unknown"
	Manual       = "manual"
	Provided     = "provided"
)

type SettingState struct {
	Schema   string `json:"schema"`
	Key      string `json:"key"`
	Current  string `json:"current"`
	Proposed string `json:"proposed"`
	Missing  bool   `json:"missing"`
}

type ItemState struct {
	Status   string         `json:"status"`
	Detail   string         `json:"detail"`
	Settings []SettingState `json:"settings,omitempty"`
}

type System struct {
	ID        string `json:"id"`
	Release   string `json:"release"`
	Codename  string `json:"codename"`
	Pretty    string `json:"pretty"`
	Arch      string `json:"arch"`
	GNOME     string `json:"gnome"`
	GNOMEMaj  int    `json:"gnomeMajor"`
	Session   string `json:"session"`
	Simulated bool   `json:"simulated"`
}

type Detector struct {
	R    sys.Runner
	Cat  *catalog.Catalog
	Sys  System
	Home string
	Root string // filesystem prefix for read-only APT inspection ("" in production)

	mu       sync.Mutex
	pkgs     map[string]bool
	cand     map[string]bool // has apt candidate
	flatpaks map[string]bool
	snaps    map[string]bool
	exts     map[string]bool
	extsOn   map[string]bool
}

func ReadSystem(r sys.Runner, simulated bool) System {
	s := System{Arch: runtime.GOARCH, Session: os.Getenv("XDG_SESSION_TYPE"), Simulated: simulated}
	if simulated {
		s.ID = "ubuntu"
		s.Release, s.Codename, s.Pretty, s.GNOME, s.GNOMEMaj = "26.04", "resolute", "Ubuntu 26.04 LTS (simulado)", "50.1", 50
		if s.Session == "" {
			s.Session = "wayland"
		}
		return s
	}
	if f, err := os.Open("/etc/os-release"); err == nil {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			k, v, _ := strings.Cut(sc.Text(), "=")
			v = strings.Trim(v, `"`)
			switch k {
			case "ID":
				s.ID = v
			case "VERSION_ID":
				s.Release = v
			case "VERSION_CODENAME":
				s.Codename = v
			case "PRETTY_NAME":
				s.Pretty = v
			}
		}
		f.Close()
	}
	if out, err := r.Output(context.Background(), "gnome-shell", "--version"); err == nil {
		if m := regexp.MustCompile(`(\d+)(\.\d+)?`).FindStringSubmatch(out); m != nil {
			s.GNOME = m[0]
			s.GNOMEMaj, _ = strconv.Atoi(m[1])
		}
	}
	return s
}

// Refresh gathers batch data: dpkg, apt candidates, flatpak, snap, extensions.
func (d *Detector) Refresh(ctx context.Context) {
	pkgs := map[string]bool{}
	if out, err := d.R.Output(ctx, "dpkg-query", "-W", "-f=${Package}\t${db:Status-Status}\n"); err == nil {
		for _, l := range strings.Split(out, "\n") {
			n, st, ok := strings.Cut(l, "\t")
			if ok && st == "installed" {
				pkgs[strings.SplitN(n, ":", 2)[0]] = true
			}
		}
	} else {
		pkgs = nil
	}
	// Candidates only for archive packages (third-party repos may not be configured yet).
	var names []string
	for _, it := range d.Cat.Items {
		if it.Method == catalog.MethodAPT {
			names = append(names, d.Cat.ExpandPackages(&it, d.defaultSel(&it))...)
		}
	}
	cand := map[string]bool{}
	if len(names) > 0 {
		if out, err := d.R.Output(ctx, "apt-cache", append([]string{"policy"}, names...)...); err == nil {
			cur := ""
			for _, l := range strings.Split(out, "\n") {
				if !strings.HasPrefix(l, " ") && strings.HasSuffix(l, ":") {
					cur = strings.TrimSuffix(l, ":")
				} else if t := strings.TrimSpace(l); strings.HasPrefix(t, "Candidate:") && cur != "" {
					cand[cur] = !strings.Contains(t, "(none)")
				}
			}
		} else {
			cand = nil
		}
	}
	flat := map[string]bool{}
	if d.R.LookPath("flatpak") {
		if out, err := d.R.Output(ctx, "flatpak", "list", "--app", "--columns=application"); err == nil {
			for _, l := range strings.Fields(out) {
				flat[l] = true
			}
		}
	}
	snaps := map[string]bool{}
	if d.R.LookPath("snap") {
		if out, err := d.R.Output(ctx, "snap", "list"); err == nil {
			for i, l := range strings.Split(out, "\n") {
				if f := strings.Fields(l); i > 0 && len(f) > 0 {
					snaps[f[0]] = true
				}
			}
		}
	}
	var exts, on map[string]bool // nil = could not list (no GNOME session?)
	if out, err := d.R.Output(ctx, "gnome-extensions", "list"); err == nil {
		exts, on = map[string]bool{}, map[string]bool{}
		for _, l := range strings.Fields(out) {
			exts[l] = true
		}
	}
	if out, err := d.R.Output(ctx, "gnome-extensions", "list", "--enabled"); err == nil && on != nil {
		for _, l := range strings.Fields(out) {
			on[l] = true
		}
	} else {
		on = nil
	}
	d.mu.Lock()
	d.pkgs, d.cand, d.flatpaks, d.snaps, d.exts, d.extsOn = pkgs, cand, flat, snaps, exts, on
	d.mu.Unlock()
}

func (d *Detector) defaultSel(it *catalog.Item) catalog.Selection {
	s, _ := d.Cat.ValidateSelection(it.ID, catalog.Selection{})
	return s
}

func (d *Detector) HasPackage(p string) bool { d.mu.Lock(); defer d.mu.Unlock(); return d.pkgs[p] }
func (d *Detector) HasSnap(p string) bool    { d.mu.Lock(); defer d.mu.Unlock(); return d.snaps[p] }
func (d *Detector) HasFlatpak(p string) bool { d.mu.Lock(); defer d.mu.Unlock(); return d.flatpaks[p] }

func versionLess(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) && i < len(pb); i++ {
		x, _ := strconv.Atoi(pa[i])
		y, _ := strconv.Atoi(pb[i])
		if x != y {
			return x < y
		}
	}
	return len(pa) < len(pb)
}

// GSettingsGet reads one key. ok=false when schema/key does not exist.
func GSettingsGet(ctx context.Context, r sys.Runner, schema, key string) (string, bool) {
	out, err := r.Output(ctx, "gsettings", "get", schema, key)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(out), true
}

func ProposedValue(st catalog.Setting, sel catalog.Selection) string {
	if st.Option != "" {
		v := sel.First(st.Option)
		if st.Map != nil {
			return st.Map[v]
		}
		return v
	}
	return st.Value
}

// ResolveSchema fills runtime placeholders. ok=false when the profile does not
// exist yet (e.g. the terminal was never opened).
func ResolveSchema(ctx context.Context, r sys.Runner, schema string) (string, bool) {
	ph := map[string][2]string{
		"{ptyxis-profile}":         {"org.gnome.Ptyxis", "default-profile-uuid"},
		"{gnome-terminal-profile}": {"org.gnome.Terminal.ProfilesList", "default"},
	}
	for k, src := range ph {
		if strings.Contains(schema, k) {
			v, ok := GSettingsGet(ctx, r, src[0], src[1])
			v = strings.Trim(v, "'")
			if !ok || v == "" || !regexp.MustCompile(`^[A-Za-z0-9-]{8,64}$`).MatchString(v) {
				return "", false
			}
			schema = strings.ReplaceAll(schema, k, v)
		}
	}
	return schema, true
}

func GitProposed(g catalog.GitSetting, sel catalog.Selection) string {
	if g.Option != "" {
		return sel.First(g.Option)
	}
	return g.Value
}

// State computes the state of one item for a given selection.
func (d *Detector) State(ctx context.Context, it *catalog.Item, sel catalog.Selection) ItemState {
	if (it.Method == catalog.MethodAPT || it.Method == catalog.MethodAPTRepo || it.AdminOp == "ubuntu-drivers") && d.Sys.ID != "ubuntu" {
		return ItemState{Status: Incompatible, Detail: "Esta receita foi preparada para Ubuntu; distribuição detetada: " + d.Sys.ID}
	}
	if len(it.Compat.Arch) > 0 {
		ok := false
		for _, a := range it.Compat.Arch {
			ok = ok || a == d.Sys.Arch
		}
		if !ok {
			return ItemState{Status: Incompatible, Detail: "Arquitetura " + d.Sys.Arch + " não suportada"}
		}
	}
	if it.Compat.MinRelease != "" && d.Sys.Release != "" && versionLess(d.Sys.Release, it.Compat.MinRelease) {
		return ItemState{Status: Incompatible, Detail: "Requer Ubuntu " + it.Compat.MinRelease + "+"}
	}
	d.mu.Lock()
	pkgs, cand, flat, exts, on := d.pkgs, d.cand, d.flatpaks, d.exts, d.extsOn
	d.mu.Unlock()
	switch it.Method {
	case catalog.MethodAPT, catalog.MethodAPTRepo:
		if pkgs == nil {
			return ItemState{Status: Unknown, Detail: "Não foi possível consultar os pacotes instalados"}
		}
		want := it.Detect.Packages
		if len(want) == 0 {
			want = d.Cat.ExpandPackages(it, sel)
		}
		var missing, nocand []string
		for _, p := range want {
			if !pkgs[p] {
				missing = append(missing, p)
				if it.Method == catalog.MethodAPT && cand != nil && !cand[p] {
					if _, known := cand[p]; known {
						nocand = append(nocand, p)
					}
				}
			}
		}
		if len(missing) == 0 {
			if len(it.Post) > 0 && !d.PostReady(it, sel) {
				return ItemState{Status: NeedsConfig, Detail: "Pacotes instalados; configuração posterior por concluir"}
			}
			return ItemState{Status: Installed, Detail: "Pacotes instalados: " + strings.Join(want, ", ")}
		}
		if len(nocand) > 0 {
			return ItemState{Status: Incompatible, Detail: "Sem candidato APT nesta versão: " + strings.Join(nocand, ", ")}
		}
		if len(missing) < len(want) {
			return ItemState{Status: Available, Detail: "Parcial; falta: " + strings.Join(missing, ", ")}
		}
		return ItemState{Status: Available, Detail: "Por instalar: " + strings.Join(missing, ", ")}
	case catalog.MethodFlatpak:
		if flat[it.Flatpak] {
			return ItemState{Status: Installed, Detail: "Flatpak " + it.Flatpak + " instalado"}
		}
		for _, old := range it.LegacyFlatpaks {
			if flat[old] {
				return ItemState{Status: Provided, Detail: "ID Flatpak anterior já instalado: " + old + "; instalação preservada, migração não automática"}
			}
		}
		if len(it.LegacyPackages) > 0 {
			present := true
			for _, pkg := range it.LegacyPackages {
				if !pkgs[pkg] {
					present = false
					break
				}
			}
			if present {
				return ItemState{Status: Provided, Detail: "Aplicação já instalada por APT; preservada. Flatpak será usado nas novas instalações."}
			}
		}
		if it.FlatpakUnavailable != "" {
			return ItemState{Status: Incompatible, Detail: it.FlatpakUnavailable}
		}
		if !d.R.LookPath("flatpak") {
			return ItemState{Status: Available, Detail: "Requer Flatpak instalado; o plano inclui essa dependência antes da aplicação."}
		}
		return ItemState{Status: Available, Detail: it.Source + ": " + it.Flatpak + "; Flatpak e origem serão validados antes de instalar"}
	case catalog.MethodGSettings:
		st := ItemState{Status: Installed, Detail: "Já está com o valor proposto"}
		for _, s := range it.Settings {
			schema, okp := ResolveSchema(ctx, d.R, s.Schema)
			if !okp {
				return ItemState{Status: NeedsConfig, Detail: "Abra a aplicação uma vez para criar o perfil predefinido e volte a detetar"}
			}
			s.Schema = schema
			cur, ok := GSettingsGet(ctx, d.R, s.Schema, s.Key)
			prop := ProposedValue(s, sel)
			st.Settings = append(st.Settings, SettingState{Schema: s.Schema, Key: s.Key, Current: cur, Proposed: prop, Missing: !ok})
			if !ok {
				return ItemState{Status: Incompatible, Detail: "Chave inexistente: " + s.Schema + " " + s.Key, Settings: st.Settings}
			}
			if cur != prop {
				st.Status, st.Detail = Available, "Diferente do proposto"
			}
		}
		return st
	case catalog.MethodGitConfig:
		if !d.R.LookPath("git") {
			return ItemState{Status: NeedsConfig, Detail: "Requer Git instalado"}
		}
		st := ItemState{Status: Installed, Detail: "Já configurado"}
		for _, g := range it.Git {
			cur, _ := d.R.Output(ctx, "git", "config", "--global", "--get", g.Key)
			cur = strings.TrimSpace(cur)
			prop := GitProposed(g, sel)
			st.Settings = append(st.Settings, SettingState{Schema: "git", Key: g.Key, Current: cur, Proposed: prop})
			if cur != prop {
				st.Status, st.Detail = Available, "Diferente do proposto"
			}
		}
		return st
	case catalog.MethodTemplate:
		p := filepath.Join(TemplatesDir(d.Home), it.Template.Name)
		if _, err := os.Stat(p); err == nil {
			return ItemState{Status: Installed, Detail: "Já existe: " + p}
		}
		return ItemState{Status: Available, Detail: "Criar " + p}
	case catalog.MethodExtension:
		if exts == nil || on == nil {
			return ItemState{Status: Unknown, Detail: "Não foi possível listar as extensões instaladas (sessão GNOME ativa?)"}
		}
		for _, pid := range it.ProvidedBy {
			if exts[pid] {
				st := "ativa"
				if !on[pid] {
					st = "instalada, desativada"
				}
				return ItemState{Status: Provided, Detail: "Equivalente já presente neste sistema (" + pid + ", " + st + ")"}
			}
		}
		if exts[it.Extension] {
			if on[it.Extension] {
				return ItemState{Status: Installed, Detail: "Instalada e ativa"}
			}
			return ItemState{Status: NeedsConfig, Detail: "Instalada mas desativada (ative após reiniciar sessão)"}
		}
		if d.Sys.GNOMEMaj == 0 {
			return ItemState{Status: Unknown, Detail: "Versão do GNOME não detetada"}
		}
		return ItemState{Status: Available, Detail: "Compatibilidade com GNOME " + strconv.Itoa(d.Sys.GNOMEMaj) + " verificada na execução (extensions.gnome.org)"}
	case catalog.MethodPowerProfile:
		if !d.R.LookPath("powerprofilesctl") {
			return ItemState{Status: Incompatible, Detail: "power-profiles-daemon não encontrado"}
		}
		cur, _ := d.R.Output(ctx, "powerprofilesctl", "get")
		cur = strings.TrimSpace(cur)
		st := ItemState{Settings: []SettingState{{Schema: "power-profiles-daemon", Key: "perfil", Current: cur, Proposed: it.Profile}}}
		if cur == it.Profile {
			st.Status, st.Detail = Installed, "Perfil atual: "+cur
		} else {
			st.Status, st.Detail = Available, "Perfil atual: "+cur
		}
		return st
	case catalog.MethodAdmin:
		return d.adminState(ctx, it, sel)
	case catalog.MethodUserScript:
		bin := expand(d.Home, it.Script.Binary)
		if _, err := os.Stat(bin); err != nil {
			return ItemState{Status: Available, Detail: "Por instalar (sem root) em " + it.Script.Binary}
		}
		want := sel.First(it.Script.Option)
		out, err := d.R.Output(ctx, bin, it.Script.VersionArgs...)
		got := strings.TrimSpace(out)
		if err != nil {
			return ItemState{Status: NeedsConfig, Detail: "Binário presente mas não executa: " + err.Error()}
		}
		st := ItemState{Settings: []SettingState{{Schema: "versão", Key: it.Name, Current: got, Proposed: want}}}
		if want == "" || got == want {
			st.Status, st.Detail = Installed, "Instalado: "+got
		} else {
			st.Status, st.Detail = Available, "Instalado "+got+"; escolhido "+want+" (substitui o binário)"
		}
		return st
	case catalog.MethodNPMUser:
		if !it.NPM.Valid() {
			return ItemState{Status: Incompatible, Detail: "Receita npm inválida"}
		}
		bin := filepath.Join(d.Home, ".local", "bin", it.NPM.Binary)
		if !d.R.LookPath(bin) {
			if !d.R.LookPath(it.NPM.Binary) {
				return ItemState{Status: Available, Detail: "Por instalar para o seu utilizador"}
			}
			bin = it.NPM.Binary
		}
		out, err := d.R.Output(ctx, bin, "--version")
		if err != nil {
			return ItemState{Status: NeedsConfig, Detail: "Comando encontrado mas não executa: " + err.Error()}
		}
		return ItemState{Status: Installed, Detail: "Instalação existente preservada: " + strings.TrimSpace(out)}
	case catalog.MethodManual:
		for _, c := range it.Detect.Commands {
			if d.R.LookPath(c) {
				return ItemState{Status: Installed, Detail: "Encontrado no PATH: " + c}
			}
		}
		for _, f := range it.Detect.Files {
			if _, err := os.Stat(expand(d.Home, f)); err == nil {
				return ItemState{Status: Installed, Detail: "Encontrado: " + f}
			}
		}
		return ItemState{Status: Manual, Detail: "Instalação manual (receita automática não verificada)"}
	}
	return ItemState{Status: Unknown, Detail: "Método sem deteção"}
}

func (d *Detector) adminState(ctx context.Context, it *catalog.Item, sel catalog.Selection) ItemState {
	switch it.AdminOp {
	case "docker-group":
		out, err := d.R.Output(ctx, "id", "-nG", strconv.Itoa(os.Getuid()))
		if err != nil {
			return ItemState{Status: Unknown, Detail: "Não foi possível consultar os grupos configurados"}
		}
		if slices.Contains(strings.Fields(out), "docker") {
			active, _ := d.R.Output(ctx, "id", "-nG")
			detail := "Utilizador já pertence ao grupo docker"
			if !slices.Contains(strings.Fields(active), "docker") {
				detail += "; termine sessão para ativar a pertença"
			}
			return ItemState{Status: Installed, Detail: detail}
		}
		if !d.HasPackage("docker-ce") && !d.HasPackage("docker.io") {
			return ItemState{Status: NeedsConfig, Detail: "Requer Docker Engine instalado"}
		}
		return ItemState{Status: Available, Detail: "Utilizador fora do grupo docker"}
	case "fstrim-timer":
		out, _ := d.R.Output(ctx, "systemctl", "is-enabled", "fstrim.timer")
		if strings.TrimSpace(out) == "enabled" {
			return ItemState{Status: Installed, Detail: "fstrim.timer já ativo"}
		}
		return ItemState{Status: Available, Detail: "fstrim.timer: " + strings.TrimSpace(out)}
	case "snap-remove-firefox":
		if d.HasSnap("firefox") {
			return ItemState{Status: Available, Detail: "Firefox Snap instalado"}
		}
		return ItemState{Status: Installed, Detail: "Firefox Snap não está instalado"}
	case "inotify", "journald":
		cfg := adminconfig.Dropins[it.AdminOp]
		p := filepath.Join(d.Root, cfg.Path)
		fi, err := os.Lstat(p)
		if os.IsNotExist(err) {
			return ItemState{Status: Available, Detail: "Criará " + cfg.Path}
		}
		if err != nil {
			return ItemState{Status: Unknown, Detail: "Não foi possível ler " + cfg.Path + ": " + err.Error()}
		}
		if !fi.Mode().IsRegular() {
			return ItemState{Status: Incompatible, Detail: "Caminho não regular: " + cfg.Path}
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return ItemState{Status: Unknown, Detail: "Não foi possível ler " + cfg.Path}
		}
		want := fmt.Sprintf(cfg.Template, sel.First("value"))
		settings := []SettingState{{Schema: "Ficheiro", Key: cfg.Path, Current: string(b), Proposed: want}}
		if string(b) == want {
			if it.AdminOp == "inotify" {
				current, err := os.ReadFile(filepath.Join(d.Root, "/proc/sys/fs/inotify/max_user_watches"))
				if err != nil || strings.TrimSpace(string(current)) != sel.First("value") {
					return ItemState{Status: NeedsConfig, Detail: "Ficheiro gravado; limite inotify ainda não confirmado no sistema", Settings: settings}
				}
			} else {
				receipt, err := os.ReadFile(p + ".applied")
				if err != nil || strings.TrimSpace(string(receipt)) != fmt.Sprintf("%x", sha256.Sum256(b)) {
					return ItemState{Status: NeedsConfig, Detail: "Ficheiro gravado; recarga de journald ainda não confirmada", Settings: settings}
				}
			}
			return ItemState{Status: Installed, Detail: "Configuração verificada: " + cfg.Path, Settings: settings}
		}
		return ItemState{Status: Available, Detail: "Configuração diferente do valor proposto em " + cfg.Path + "; guarda o ficheiro anterior antes de alterar", Settings: settings}
	case "ubuntu-drivers":
		out, err := d.R.Output(ctx, "ubuntu-drivers", "devices")
		if err != nil {
			return ItemState{Status: Unknown, Detail: "Não foi possível consultar os drivers recomendados"}
		}
		matches := regexp.MustCompile(`(?m)^\s*driver\s*:\s*([a-zA-Z0-9.+:-]+)\s+.*\brecommended\b`).FindAllStringSubmatch(out, -1)
		if len(matches) == 0 {
			return ItemState{Status: Provided, Detail: "Nenhum driver adicional recomendado"}
		}
		var missing []string
		for _, m := range matches {
			if !d.HasPackage(m[1]) {
				missing = append(missing, m[1])
			}
		}
		if len(missing) == 0 {
			return ItemState{Status: Installed, Detail: "Drivers recomendados já instalados"}
		}
		return ItemState{Status: Available, Detail: "Drivers recomendados em falta: " + strings.Join(missing, ", ")}
	}
	return ItemState{Status: Unknown}
}

func TemplatesDir(home string) string {
	if b, err := os.ReadFile(filepath.Join(home, ".config/user-dirs.dirs")); err == nil {
		for _, l := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(l, "XDG_TEMPLATES_DIR=") {
				v := strings.Trim(strings.TrimPrefix(l, "XDG_TEMPLATES_DIR="), `"`)
				v = strings.ReplaceAll(v, "$HOME", home)
				if v != home && v != "" {
					return v
				}
			}
		}
	}
	return filepath.Join(home, "Modelos")
}

func expand(home, p string) string {
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(home, p[2:])
	}
	return p
}
