// Package engine executes a plan step by step, streaming real logs, verifying
// the real result of each step, and journaling every preference it changes.
package engine

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/store"
	"meuubuntu/internal/sys"
)

const (
	StPending      = "pending"
	StRunning      = "running"
	StDone         = "done"
	StFailed       = "failed"
	StSkipped      = "skipped"
	StIntervention = "intervention"
	StCancelled    = "cancelled"
)

type Engine struct {
	Cat       *catalog.Catalog
	Det       *detect.Detector
	R         sys.Runner // user-level runner
	NewAdmin  func() helper.Admin
	Store     *store.Store
	Home      string
	Emit      func(event string, data any)
	Fetch     func(string) ([]byte, error)
	Simulated bool

	mu      sync.Mutex
	running bool
	closing bool
	cancel  atomic.Bool
}

type StepEvent struct {
	RunID   string         `json:"runId"`
	StepID  string         `json:"stepId"`
	Status  string         `json:"status"`
	Error   string         `json:"error,omitempty"`
	Note    string         `json:"note,omitempty"`
	Time    string         `json:"time"`
	Problem *store.Problem `json:"problem,omitempty"`
}

type LogEvent struct {
	RunID  string `json:"runId"`
	StepID string `json:"stepId"`
	Time   string `json:"time"`
	Line   string `json:"line"`
}

func (e *Engine) Running() bool { e.mu.Lock(); defer e.mu.Unlock(); return e.running }

// PreventClose holds the same lock as Start so closing cannot race a new run.
// A normal window close waits for the running operation's safe boundary.
func (e *Engine) PreventClose() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.running {
		return true
	}
	e.closing = true
	return false
}

// Cancel requests a stop at the next safe boundary. A running APT/dpkg
// operation is never killed; it finishes first.
func (e *Engine) Cancel() { e.cancel.Store(true) }

// Start runs the plan asynchronously and returns the run ID.
// confirmed maps item -> the value the user confirmed in the reviewed plan
// ("yes", or the SHA-256 shown for a vendor script).
func (e *Engine) Start(p *plan.Plan, confirmed map[string]string) (string, error) {
	// Go callers may build plans directly; the public App API already binds
	// these values before checking the client's review token.
	for _, s := range p.Steps {
		if (s.Op == "undo" || s.Op == "revert") && s.Undo == nil {
			if err := e.ReviewUndo(context.Background(), p); err != nil {
				return "", err
			}
			break
		}
	}
	for _, s := range p.Steps {
		if s.Confirm == "" {
			continue
		}
		want := "yes"
		if s.Expect != "" {
			want = s.Expect
		}
		if confirmed[s.Item] != want {
			return "", fmt.Errorf("falta confirmar explicitamente: %s", s.Title)
		}
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closing {
		return "", errors.New("a aplicação está a encerrar; execução não iniciada")
	}
	if e.running {
		return "", errors.New("já existe uma execução em curso")
	}
	if p.Blocked {
		return "", errors.New("o plano tem conflitos que bloqueiam a execução")
	}
	run := store.Run{ID: store.NewID(), Started: time.Now().Format(time.RFC3339), Status: "running", Selections: p.Selection, Reverts: p.Reverts, Simulated: e.Simulated}
	run.LogFile = e.Store.LogPath(run.ID)
	for _, s := range p.Steps {
		run.Steps = append(run.Steps, store.StepRecord{ID: s.ID, Title: s.Title, Item: s.Item, Op: s.Op, Status: StPending})
	}
	logf, err := os.OpenFile(run.LogFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", fmt.Errorf("não foi possível criar o log; execução não iniciada: %w", err)
	}
	if err := e.Store.SaveRun(run); err != nil {
		logf.Close()
		os.Remove(run.LogFile)
		return "", fmt.Errorf("não foi possível guardar a execução; nada alterado: %w", err)
	}
	e.running = true
	e.cancel.Store(false)
	go e.exec(p, run, logf)
	return run.ID, nil
}

func (e *Engine) exec(p *plan.Plan, run store.Run, logf *os.File) {
	defer func() { e.mu.Lock(); e.running = false; e.mu.Unlock() }()
	defer logf.Close()
	var admin helper.Admin
	defer func() {
		if admin != nil {
			admin.Close()
		}
	}()
	ctx := context.Background() // never cancel a running process (dpkg safety)
	failed := false
	intervention := false
	failedItems := map[string]bool{}
	aptBroken := ""
	var persistenceErr error
	persist := func() bool {
		if persistenceErr == nil {
			persistenceErr = logf.Sync()
		}
		if persistenceErr == nil {
			persistenceErr = e.Store.SaveRun(run)
		}
		if persistenceErr != nil {
			run.Status = "failed"
			e.Emit("run:log", LogEvent{RunID: run.ID, Time: time.Now().Format(time.RFC3339), Line: "Falha ao guardar histórico/log: " + persistenceErr.Error() + ". Execução interrompida no limite seguro."})
			return false
		}
		return true
	}
	for i, s := range p.Steps {
		if e.cancel.Load() {
			for j := i; j < len(p.Steps); j++ {
				run.Steps[j].Status = StCancelled
				e.emitStep(run.ID, p.Steps[j].ID, StCancelled, "", "Cancelado antes de iniciar")
			}
			run.Status = "cancelled"
			break
		}
		if aptBroken != "" && isAPTStep(s) {
			run.Steps[i].Status = StSkipped
			run.Steps[i].Error = "Ignorado: as fontes APT estão inválidas (" + aptBroken + "). Corrija-as e tente novamente."
			if !persist() {
				break
			}
			e.emitStep(run.ID, s.ID, StSkipped, "", run.Steps[i].Error)
			continue
		}
		if blocker := e.blockedBy(s.Item, failedItems); blocker != "" {
			run.Steps[i].Status = StSkipped
			run.Steps[i].Error = "Ignorado: falhou um passo anterior de " + blocker
			failedItems[s.Item] = true
			if !persist() {
				break
			}
			e.emitStep(run.ID, s.ID, StSkipped, "", run.Steps[i].Error)
			continue
		}
		run.Steps[i].Status = StRunning
		if !persist() {
			break
		}
		e.emitStep(run.ID, s.ID, StRunning, "", "")
		var diagnosticLines []string
		log := func(line string) {
			diagnostic := redact(line)
			if len(diagnostic) > 4096 {
				diagnostic = diagnostic[:4096]
			}
			diagnosticLines = append(diagnosticLines, diagnostic)
			if len(diagnosticLines) > 40 {
				diagnosticLines = diagnosticLines[len(diagnosticLines)-40:]
			}
			if strings.Contains(line, "Conflicting values set for option") || strings.Contains(line, "The list of sources could not be read") || strings.Contains(line, "Malformed entry") {
				aptBroken = strings.TrimSpace(strings.TrimPrefix(line, "E: "))
			}
			t := time.Now()
			if _, err := fmt.Fprintf(logf, "%s [%s] %s\n", t.Format("15:04:05"), s.ID, redact(line)); err != nil {
				persistenceErr = err
			}
			e.Emit("run:log", LogEvent{RunID: run.ID, StepID: s.ID, Time: t.Format(time.RFC3339), Line: redact(line)})
		}
		log("▶ " + s.Title)
		var err error
		var changeID, changeAfter string
		if s.Op == "undo" {
			err = e.validateUndo(ctx, s)
		}
		if err == nil && s.Kind == "admin" {
			if admin == nil {
				log("A pedir autorização administrativa ao sistema (polkit)…")
				admin = e.NewAdmin()
			}
			opts := s.Options
			if s.Op == "admin" {
				if it := e.Cat.Item(s.Item); it != nil && (it.AdminOp == "docker-group" || it.AdminOp == "inotify" || it.AdminOp == "journald") {
					changeAfter = strings.Join(opts["value"], ",")
					changes, readErr := e.Store.ReadChanges()
					err = readErr
					for j := len(changes) - 1; j >= 0 && err == nil; j-- {
						c := changes[j]
						if c.Item == it.ID && c.Kind == "admin" && c.Reverted == "" {
							if c.After == changeAfter {
								changeID = c.ID
							}
							break
						}
					}
					if changeID == "" && err == nil {
						changeID, err = e.Store.BeginChange(store.Change{RunID: run.ID, Item: it.ID, Title: it.Name, Kind: "admin", Key: it.AdminOp, After: changeAfter})
					}
				}
			}
			if err == nil && persistenceErr == nil {
				if s.Op == "revert" {
					err = e.validateAdminUndo(ctx, s, admin)
					if err == nil {
						err = e.revertAdminItem(ctx, s.Item, admin, log, s.Undo[0].StateToken)
					}
				} else {
					err = admin.Do(ctx, helper.Request{Op: s.Op, Item: s.Item, Options: opts}, log)
				}
			}
			if err == nil && changeID != "" {
				err = e.Store.FinishChange(changeID, changeAfter)
			}
		} else if err == nil {
			if persistenceErr == nil {
				err = e.userOp(ctx, run.ID, s, log)
			}
		}
		if persistenceErr != nil {
			err = persistenceErr
		}
		status, note := StDone, ""
		if err != nil {
			status = StFailed
			log("✖ " + err.Error())
		} else {
			status, note = e.verify(ctx, s, log)
		}
		if status == StFailed {
			failed = true
			if s.Item != "" {
				failedItems[s.Item] = true
			}
			if err == nil {
				err = errors.New(note)
			}
		}
		run.Steps[i].Status = status
		run.Steps[i].Note = note
		if note != "" {
			log(note)
		}
		if status == StIntervention {
			intervention = true
			failedItems[s.Item] = true
		}
		if err != nil {
			run.Steps[i].Error = err.Error()
		}
		errS := ""
		if err != nil {
			errS = err.Error()
		}
		if status == StFailed {
			run.Steps[i].Problem = executionProblem(errS, strings.Join(diagnosticLines, "\n"))
		}
		e.Emit("run:step", StepEvent{RunID: run.ID, StepID: s.ID, Status: status, Error: errS, Note: note, Problem: run.Steps[i].Problem, Time: time.Now().Format(time.RFC3339)})
		if !persist() {
			break
		}
		// A failed repo or APT index step makes later APT steps meaningless.
		if status == StFailed && (s.Op == "apt-update") {
			log("Índices APT falharam; os restantes passos APT podem falhar.")
		}
	}
	if persistenceErr != nil {
		for i := range run.Steps {
			if run.Steps[i].Status == StPending || run.Steps[i].Status == StRunning {
				run.Steps[i].Status = StCancelled
				e.emitStep(run.ID, run.Steps[i].ID, StCancelled, "", "Interrompido por falha de persistência")
			}
		}
	}
	if run.Status == "running" {
		run.Status = "done"
		if failed {
			run.Status = "failed"
		} else if intervention {
			run.Status = StIntervention
		}
	}
	run.Finished = time.Now().Format(time.RFC3339)
	if err := e.Store.SaveRun(run); err != nil {
		run.Status = "failed"
		persistenceErr = err
	}
	done := map[string]any{"runId": run.ID, "status": run.Status}
	if persistenceErr != nil {
		done["error"] = "Não foi possível guardar o histórico: " + persistenceErr.Error()
	}
	e.Emit("run:done", done)
}

// blockedBy returns the name of a failed item this step depends on (itself or a requirement).
func (e *Engine) revertAdminItem(ctx context.Context, item string, admin helper.Admin, log func(string), reviewedToken ...string) error {
	changes, err := e.Store.ReadChanges()
	if err != nil {
		return err
	}
	found := false
	singleOperation := e.Cat.Item(item).AdminOp == "docker-group"
	applied := false
	for i := len(changes) - 1; i >= 0; i-- {
		c := changes[i]
		if c.Item != item || c.Kind != "admin" || c.Reverted != "" {
			continue
		}
		found = true
		var opts map[string][]string
		if c.After != "" {
			opts = map[string][]string{"value": strings.Split(c.After, ",")}
		}
		if !singleOperation || !applied {
			req := helper.Request{Op: "revert", Item: item, Options: opts}
			if !applied && len(reviewedToken) > 0 {
				req.Batch = true
				req.ExpectedStateToken = reviewedToken[0]
			}
			if err := admin.Do(ctx, req, log); err != nil {
				return err
			}
			applied = true
		}
		if err := e.Store.MarkReverted(c.ID); err != nil {
			return err
		}
	}
	if !found {
		return admin.Do(ctx, helper.Request{Op: "revert", Item: item}, log)
	}
	return nil
}

func isAPTStep(s plan.Step) bool {
	switch s.Op {
	case "apt-update", "apt-install", "repo-add", "post":
		return true
	}
	return false
}

func (e *Engine) blockedBy(id string, failed map[string]bool) string {
	if id == "" {
		return ""
	}
	if failed[id] {
		return e.Cat.Item(id).Name
	}
	for _, r := range e.Cat.Item(id).Requires {
		if failed[r] {
			return e.Cat.Item(r).Name
		}
		if blocker := e.blockedBy(r, failed); blocker != "" {
			return blocker
		}
	}
	return ""
}

func (e *Engine) emitStep(runID, stepID, status, errS, note string) {
	e.Emit("run:step", StepEvent{RunID: runID, StepID: stepID, Status: status, Error: errS, Note: note, Time: time.Now().Format(time.RFC3339)})
}

// verify re-detects the real state after a step.
func (e *Engine) verify(ctx context.Context, s plan.Step, log func(string)) (string, string) {
	if s.Item == "" || s.Op == "repo-add" || s.Op == "flatpak-remote" || s.Op == "undo" || s.Op == "revert" || s.Op == "sources-check" {
		return StDone, ""
	}
	it := e.Cat.Item(s.Item)
	if s.Op == "apt-install" || s.Op == "flatpak-install" || s.Op == "extension" || s.Op == "extension-enable" || s.Op == "admin" {
		e.Det.Refresh(ctx)
	}
	if s.Op == "apt-install" && len(it.Post) > 0 && e.Det.PackagesPresent(it, catalog.Selection{Options: s.Options}) {
		return StDone, "Pacotes instalados; configuração no passo seguinte"
	}
	st := e.Det.State(ctx, it, catalog.Selection{Options: s.Options})
	log("Verificação: " + st.Status + " — " + st.Detail)
	switch st.Status {
	case detect.Installed, detect.Provided:
		if len(it.Verify.Cmd) > 0 {
			out, err := e.R.Output(ctx, it.Verify.Cmd[0], it.Verify.Cmd[1:]...)
			if err != nil {
				return StIntervention, "Instalado, mas a verificação falhou: " + err.Error()
			}
			log("✔ " + strings.TrimSpace(firstLine(out)))
		}
		note := ""
		if it.Relogin {
			note = "Termine sessão para aplicar"
		}
		return StDone, note
	case detect.NeedsConfig, detect.Unknown:
		return StIntervention, st.Detail
	}
	return StFailed, "Resultado não confirmado: " + st.Detail
}

func firstLine(s string) string { l, _, _ := strings.Cut(s, "\n"); return l }

func (e *Engine) userOp(ctx context.Context, runID string, s plan.Step, log func(string)) error {
	it := e.Cat.Item(s.Item)
	sel := catalog.Selection{Options: s.Options}
	switch s.Op {
	case "verify":
		return nil // verify() performs the read-only check after this step.
	case "flatpak-remote":
		return e.R.Stream(ctx, log, "flatpak", "remote-add", "--user", "--if-not-exists", "flathub", "https://dl.flathub.org/repo/flathub.flatpakrepo")
	case "flatpak-install":
		return e.R.Stream(ctx, log, "flatpak", "install", "--user", "-y", "--noninteractive", "flathub", it.Flatpak)
	case catalog.MethodGSettings:
		for _, st := range it.Settings {
			schema, okp := detect.ResolveSchema(ctx, e.R, st.Schema)
			if !okp {
				return errors.New("perfil predefinido inexistente: abra a aplicação uma vez e tente novamente")
			}
			st.Schema = schema
			cur, ok := detect.GSettingsGet(ctx, e.R, st.Schema, st.Key)
			if !ok {
				return fmt.Errorf("chave inexistente %s %s", st.Schema, st.Key)
			}
			prop := detect.ProposedValue(st, sel)
			if cur == prop {
				log("Sem alteração: " + st.Key)
				continue
			}
			log(fmt.Sprintf("gsettings set %s %s %s (antes: %s)", st.Schema, st.Key, prop, cur))
			c := store.Change{RunID: runID, Item: it.ID, Title: it.Name, Kind: "gsettings", Schema: st.Schema, Key: st.Key, Before: cur, After: prop, HadValue: true}
			if err := e.change(c, func() error { return e.R.Stream(ctx, log, "gsettings", "set", st.Schema, st.Key, prop) }); err != nil {
				return err
			}
		}
		return nil
	case catalog.MethodGitConfig:
		for _, g := range it.Git {
			out, err := e.R.Output(ctx, "git", "config", "--global", "--get", g.Key)
			cur := strings.TrimSpace(out)
			prop := detect.GitProposed(g, sel)
			if cur == prop {
				continue
			}
			log(sys.Argv("git", "config", "--global", g.Key, prop))
			c := store.Change{RunID: runID, Item: it.ID, Title: it.Name, Kind: "git", Key: g.Key, Before: cur, After: prop, HadValue: err == nil}
			if err := e.change(c, func() error { return e.R.Stream(ctx, log, "git", "config", "--global", g.Key, prop) }); err != nil {
				return err
			}
		}
		return nil
	case catalog.MethodPowerProfile:
		out, err := e.R.Output(ctx, "powerprofilesctl", "get")
		if err != nil {
			return err
		}
		cur := strings.TrimSpace(out)
		if cur == it.Profile {
			return nil
		}
		return e.change(store.Change{RunID: runID, Item: it.ID, Title: it.Name, Kind: "power", Key: "perfil", Before: cur, After: it.Profile, HadValue: true}, func() error { return e.R.Stream(ctx, log, "powerprofilesctl", "set", it.Profile) })
	case "template":
		dir := detect.TemplatesDir(e.Home)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		p := filepath.Join(dir, filepath.Base(it.Template.Name))
		if _, err := os.Lstat(p); err == nil {
			log("Já existe, preservado: " + p)
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}
		id, err := e.Store.BeginChange(store.Change{RunID: runID, Item: it.ID, Title: it.Name, Kind: "template", Key: p, After: it.Template.Content})
		if err != nil {
			return err
		}
		f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			if os.IsExist(err) {
				log("Já existe, preservado: " + p)
				return errors.New("o modelo foi criado por outro processo; nada sobrescrito")
			}
			return err
		}
		_, err = f.WriteString(it.Template.Content)
		if err == nil {
			err = f.Sync()
		}
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			return err
		}
		log("Criado " + p)
		return e.Store.FinishChange(id, it.Template.Content)
	case "extension":
		return e.installExtension(ctx, runID, it, log)
	case "extension-enable":
		return e.enableExtension(ctx, runID, it, log)
	case "user-script":
		return e.runUserScript(ctx, it, sel, s.Expect, log)
	case catalog.MethodNPMUser:
		return e.installNPMUser(ctx, it, log)
	case "undo":
		return e.undo(ctx, runID, it, log, s.Undo...)
	}
	return fmt.Errorf("operação de utilizador desconhecida %q", s.Op)
}

func (e *Engine) change(c store.Change, apply func() error) error {
	id, err := e.Store.BeginChange(c)
	if err != nil {
		return fmt.Errorf("não foi possível guardar a alteração; nada aplicado: %w", err)
	}
	if err := apply(); err != nil {
		return err
	}
	return e.Store.FinishChange(id, c.After)
}

func (e *Engine) installExtension(ctx context.Context, runID string, it *catalog.Item, log func(string)) error {
	shell := e.Det.Sys.GNOMEMaj
	if shell == 0 {
		return errors.New("versão do GNOME desconhecida")
	}
	api := "https://extensions.gnome.org/extension-info/?uuid=" + url.QueryEscape(it.Extension) + "&shell_version=" + strconv.Itoa(shell)
	log("A consultar " + api)
	b, err := e.Fetch(api)
	if err != nil {
		return fmt.Errorf("sem versão compatível com GNOME %d em extensions.gnome.org (%v)", shell, err)
	}
	var info struct {
		DownloadURL string `json:"download_url"`
		Version     any    `json:"version"`
		UUID        string `json:"uuid"`
	}
	if err := json.Unmarshal(b, &info); err != nil || info.UUID != it.Extension || !strings.HasPrefix(info.DownloadURL, "/download-extension/") {
		return errors.New("resposta inesperada da API de extensões")
	}
	zip, err := e.Fetch("https://extensions.gnome.org" + info.DownloadURL)
	if err != nil {
		return err
	}
	tmp, err := os.MkdirTemp("", "meu-ubuntu-ext")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	zf := filepath.Join(tmp, "ext.zip")
	if err := os.WriteFile(zf, zip, 0o600); err != nil {
		return err
	}
	log(fmt.Sprintf("Versão %v para GNOME %d", info.Version, shell))
	if err := e.R.Stream(ctx, log, "gnome-extensions", "install", zf); err != nil {
		return err
	}
	return e.enableExtension(ctx, runID, it, log)
}

func (e *Engine) enableExtension(ctx context.Context, runID string, it *catalog.Item, log func(string)) error {
	cur, ok := detect.GSettingsGet(ctx, e.R, "org.gnome.shell", "enabled-extensions")
	if !ok {
		return errors.New("não foi possível ler as extensões ativas; nada alterado")
	}
	next := appendGVariantList(cur, it.Extension)
	return e.change(store.Change{RunID: runID, Item: it.ID, Title: it.Name, Kind: "gsettings", Schema: "org.gnome.shell", Key: "enabled-extensions", Before: cur, After: next, HadValue: true}, func() error {
		if err := e.R.Stream(ctx, log, "gnome-extensions", "enable", it.Extension); err == nil {
			return nil
		}
		// Wayland only loads new extensions at the next login: mark for enabling then.
		log("Será ativada no próximo início de sessão")
		return e.R.Stream(ctx, log, "gsettings", "set", "org.gnome.shell", "enabled-extensions", next)
	})
}

func appendGVariantList(cur, v string) string {
	cur = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(cur), "@as"))
	if cur == "" || cur == "[]" {
		return "['" + v + "']"
	}
	if strings.Contains(cur, "'"+v+"'") {
		return cur
	}
	return strings.TrimSuffix(cur, "]") + ", '" + v + "']"
}

// Revert undoes one journaled change, only if the value is still what we set.
// RevertChange serializes history actions with Start and other reversals.
func (e *Engine) RevertChange(ctx context.Context, id string, force bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.running {
		return errors.New("aguarde o fim da execução antes de reverter alterações")
	}
	changes, err := e.Store.ReadChanges()
	if err != nil {
		return err
	}
	for _, c := range changes {
		if c.ID != id {
			continue
		}
		if c.Reverted != "" {
			return errors.New("esta alteração já foi revertida")
		}
		return e.Revert(ctx, c, force)
	}
	return errors.New("alteração não encontrada")
}

func (e *Engine) Revert(ctx context.Context, c store.Change, force bool, stateToken ...string) error {
	switch c.Kind {
	case "gsettings":
		cur, ok := detect.GSettingsGet(ctx, e.R, c.Schema, c.Key)
		if !ok {
			return errors.New("não foi possível ler o valor atual; reversão recusada")
		}
		if cur != c.After && !force {
			return fmt.Errorf("valor atual (%s) já não é o que foi aplicado (%s); reversão não feita para não apagar alterações posteriores", cur, c.After)
		}
		if _, err := e.R.Output(ctx, "gsettings", "set", c.Schema, c.Key, c.Before); err != nil {
			return err
		}
	case "git":
		out, readErr := e.gitValue(ctx, c.Key)
		if readErr != nil {
			return readErr
		}
		if out != c.After && !force {
			return errors.New("valor atual do git mudou desde a aplicação; reversão recusada")
		}
		var err error
		if c.HadValue {
			_, err = e.R.Output(ctx, "git", "config", "--global", c.Key, c.Before)
		} else {
			_, err = e.R.Output(ctx, "git", "config", "--global", "--unset", c.Key)
		}
		if err != nil {
			return err
		}
	case "power":
		out, err := e.R.Output(ctx, "powerprofilesctl", "get")
		if err != nil {
			return err
		}
		if strings.TrimSpace(out) != c.After && !force {
			return errors.New("perfil de energia alterado depois da aplicação; reversão recusada")
		}
		if _, err := e.R.Output(ctx, "powerprofilesctl", "set", c.Before); err != nil {
			return err
		}
	case "template":
		if fi, err := os.Lstat(c.Key); err == nil && !fi.Mode().IsRegular() {
			return errors.New("modelo não é um ficheiro regular; reversão recusada")
		} else if err != nil && !os.IsNotExist(err) {
			return err
		}
		b, err := os.ReadFile(c.Key)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && string(b) != c.After && !force {
			return errors.New("o modelo foi editado depois de criado; não foi apagado")
		}
		if err == nil {
			if err := os.Remove(c.Key); err != nil {
				return err
			}
		}
	case "admin":
		a := e.NewAdmin()
		defer a.Close()
		var opts map[string][]string
		if c.After != "" {
			opts = map[string][]string{"value": strings.Split(c.After, ",")}
		}
		req := helper.Request{Op: "revert", Item: c.Item, Options: opts}
		if len(stateToken) > 0 {
			req.ExpectedStateToken = stateToken[0]
		}
		if err := a.Do(ctx, req, func(string) {}); err != nil {
			return err
		}
	default:
		return errors.New("tipo de alteração sem reversão automática")
	}
	return e.Store.MarkReverted(c.ID)
}

// redact hides anything that looks like a credential in logs.
func redact(s string) string {
	low := strings.ToLower(s)
	for _, k := range []string{"password=", "passwd=", "token=", "secret="} {
		if i := strings.Index(low, k); i >= 0 {
			return s[:i+len(k)] + "***"
		}
	}
	return s
}

var hostRe = regexp.MustCompile(`https?://([A-Za-z0-9.\-]+)`)

// runUserScript runs a vendor installer as the user. Safety comes from the
// user's explicit review of this exact script (SHA-256), not from the
// heuristics below, which only catch obvious surprises.
func (e *Engine) runUserScript(ctx context.Context, it *catalog.Item, sel catalog.Selection, expect string, log func(string)) error {
	u := plan.ScriptURL(it, sel)
	log("A transferir " + u)
	body, err := e.Fetch(u)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(body)
	log(fmt.Sprintf("SHA-256 %x (%d bytes)", sum, len(body)))
	if expect == "" || fmt.Sprintf("%x", sum) != expect {
		return fmt.Errorf("o script mudou desde a revisão do plano (esperado %s); nada foi executado", expect)
	}
	if regexp.MustCompile(`(^|[\s;&|(])sudo(\s|$)`).Match(body) {
		return errors.New("o script pede sudo; recusado (este método só corre sem root)")
	}
	allowed := map[string]bool{}
	for _, h := range it.Script.AllowHosts {
		allowed[h] = true
	}
	for _, m := range hostRe.FindAllSubmatch(body, -1) {
		if !allowed[string(m[1])] {
			return fmt.Errorf("o script referencia um host fora da lista permitida: %s; recusado", m[1])
		}
	}
	bdir := filepath.Join(e.Store.State, "backups", store.NewID())
	for _, b := range it.Script.Backup {
		src := filepath.Join(e.Home, strings.TrimPrefix(b, "~/"))
		data, err := os.ReadFile(src)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("não foi possível ler %s para backup: %w", src, err)
		}
		if err := os.MkdirAll(bdir, 0o700); err != nil {
			return err
		}
		dst := filepath.Join(bdir, filepath.Base(src))
		if err := os.WriteFile(dst, data, 0o600); err != nil {
			return err
		}
		log("Backup: " + src + " → " + dst)
	}
	tmp, err := os.MkdirTemp("", "meu-ubuntu-script")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	f := filepath.Join(tmp, "install.sh")
	if err := os.WriteFile(f, body, 0o600); err != nil {
		return err
	}
	keep := filepath.Join(e.Store.State, "scripts")
	if err := os.MkdirAll(keep, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(keep, it.ID+"-"+fmt.Sprintf("%x", sum[:6])+".sh"), body, 0o600); err != nil {
		return fmt.Errorf("não foi possível guardar o script revisto; nada executado: %w", err)
	}
	return e.R.Stream(ctx, log, "bash", f)
}

// undo disables an applied item the user unchecked. Preferences go back to the
// value recorded before Meu Ubuntu changed them, or to the system default when
// there is no record. Git settings and templates are only undone when recorded.
func (e *Engine) undo(ctx context.Context, runID string, it *catalog.Item, log func(string), reviewed ...plan.UndoValue) error {
	changes, err := e.Store.ReadChanges()
	if err != nil {
		return fmt.Errorf("histórico ilegível; reversão recusada: %w", err)
	}
	latest := func(kind, schema, key string) *store.Change {
		var found *store.Change
		for _, c := range changes {
			if c.Kind == kind && c.Item == it.ID && c.Schema == schema && c.Key == key {
				cc := c
				found = &cc
			}
		}
		return found
	}
	switch it.Method {
	case catalog.MethodGSettings:
		for _, st := range it.Settings {
			schema, ok := detect.ResolveSchema(ctx, e.R, st.Schema)
			if !ok {
				continue
			}
			cur, _ := detect.GSettingsGet(ctx, e.R, schema, st.Key)
			for _, v := range reviewed {
				if v.Schema == schema && v.Key == st.Key && v.Current != cur {
					return errors.New("o valor mudou depois da revisão; nada reposto nesta chave")
				}
			}
			if c := latest("gsettings", schema, st.Key); c != nil {
				if c.Reverted != "" {
					log("Já restaurado: " + schema + " " + st.Key)
					continue
				}
				log(fmt.Sprintf("Repor valor anterior %s %s: %s → %s", schema, st.Key, cur, c.Before))
				if err := e.Revert(ctx, *c, false); err != nil {
					return err
				}
				continue
			}
			id, err := e.Store.BeginChange(store.Change{RunID: runID, Item: it.ID, Title: it.Name + " (reposto)", Kind: "gsettings", Schema: schema, Key: st.Key, Before: cur, HadValue: true})
			if err != nil {
				return err
			}
			if _, err := e.R.Output(ctx, "gsettings", "reset", schema, st.Key); err != nil {
				return err
			}
			def, _ := detect.GSettingsGet(ctx, e.R, schema, st.Key)
			log(fmt.Sprintf("Repor padrão do sistema %s %s: %s → %s", schema, st.Key, cur, def))
			if err := e.Store.FinishReset(id, def); err != nil {
				return err
			}
		}
		return nil
	case catalog.MethodGitConfig:
		for _, g := range it.Git {
			if c := latest("git", "", g.Key); c != nil {
				if c.Reverted != "" {
					continue
				}
				log("Repor git " + g.Key)
				if err := e.Revert(ctx, *c, false); err != nil {
					return err
				}
			} else {
				log("Sem registo para git " + g.Key + ": mantém-se o valor atual")
			}
		}
		return nil
	case catalog.MethodTemplate:
		key := filepath.Join(detect.TemplatesDir(e.Home), filepath.Base(it.Template.Name))
		for i := len(changes) - 1; i >= 0; i-- {
			c := changes[i]
			if c.Kind == "template" && c.Item == it.ID && c.Key == key {
				if c.Reverted != "" {
					log("Modelo já restaurado; mantido")
					return nil
				}
				log("Remover modelo criado pelo Meu Ubuntu: " + c.Key)
				return e.Revert(ctx, c, false)
			}
		}
		log("Modelo não foi criado pelo Meu Ubuntu: mantido")
		return nil
	}
	return fmt.Errorf("sem desativação automática para %s", it.Name)
}
