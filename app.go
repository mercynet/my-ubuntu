package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"runtime/debug"
	"slices"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/diag"
	"meuubuntu/internal/engine"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/inventory"
	"meuubuntu/internal/migrate"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/profiles"
	"meuubuntu/internal/selection"
	"meuubuntu/internal/store"
	"meuubuntu/internal/sys"
)

// App is the API exposed to the React UI through Wails bindings.
type App struct {
	ctx                context.Context
	pendingPath        string // profile JSON given on the command line (review only)
	cat                *catalog.Catalog
	det                *detect.Detector
	r                  sys.Runner
	st                 *store.Store
	eng                *engine.Engine
	home               string
	simulated          bool
	diagCache          diagnosticsCache
	integrationPending integrationPending
}

func NewApp(simulated bool) (*App, error) {
	cat, err := catalog.Load()
	if err != nil {
		return nil, err
	}
	name := "meu-ubuntu"
	var r sys.Runner = sys.Real{}
	if simulated {
		name = "meu-ubuntu-sim"
		r = sys.NewSim()
	}
	st, err := store.New(name)
	if err != nil {
		return nil, err
	}
	home, _ := os.UserHomeDir()
	fetch := helper.HTTPFetch
	if simulated {
		home = filepath.Join(st.Data, "home-simulada") // never touch the real home in simulation
		os.MkdirAll(home, 0o700)
		fetch = func(u string) ([]byte, error) { return simFetch(cat, u) }
	}
	a := &App{cat: cat, r: r, st: st, home: home, simulated: simulated}
	a.det = &detect.Detector{R: r, Cat: cat, Sys: detect.ReadSystem(r, simulated), Home: home}
	if simulated {
		a.det.Root = filepath.Join(st.Data, "root-simulado")
		watchers := filepath.Join(a.det.Root, "proc/sys/fs/inotify/max_user_watches")
		if err := os.MkdirAll(filepath.Dir(watchers), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(watchers, []byte("65536\n"), 0o600); err != nil {
			return nil, err
		}
	}
	a.eng = &engine.Engine{Cat: cat, Det: a.det, R: r, Store: st, Home: home, Fetch: fetch, Simulated: simulated,
		Emit: func(ev string, data any) {
			if a.ctx != nil {
				runtime.EventsEmit(a.ctx, ev, data)
			}
		}}
	a.eng.NewAdmin = func() helper.Admin {
		if simulated {
			root := filepath.Join(st.Data, "root-simulado")
			os.MkdirAll(root, 0o755)
			return helper.InProcess{Env: &helper.Env{R: r, Cat: cat, Root: root, Codename: a.det.Sys.Codename, UID: "1000", Fetch: fetch}}
		}
		return &helper.Pkexec{}
	}
	return a, nil
}

func simFetch(cat *catalog.Catalog, url string) ([]byte, error) {
	for _, r := range cat.Repos {
		if r.KeyURL == url {
			return []byte("SIM-FPR:" + strings.Join(r.Fingerprints, "\nSIM-FPR:")), nil
		}
	}
	switch {
	case strings.Contains(url, "extension-info"):
		uuid := url[strings.Index(url, "uuid=")+5 : strings.Index(url, "&")]
		uuid = strings.ReplaceAll(uuid, "%40", "@")
		return []byte(`{"uuid":"` + uuid + `","version":1,"download_url":"/download-extension/x.zip?version_tag=1"}`), nil
	case strings.Contains(url, "php.new"):
		return []byte("#!/bin/bash\necho 'instalador simulado'\n"), nil
	}
	return []byte("simulado"), nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// A run left "running" means the app died mid-run: mark it so the UI can offer resuming.
	for _, r := range a.st.Runs() {
		if r.Status == "running" {
			r.Status = "interrupted"
			a.st.SaveRun(r)
		}
	}
}

type ItemView struct {
	catalog.Item
	State detect.ItemState `json:"state"`
}

type InitData struct {
	Version     string            `json:"version"`
	System      detect.System     `json:"system"`
	Prefs       store.Prefs       `json:"prefs"`
	Profiles    []store.Profile   `json:"profiles"`
	Interrupted *store.Run        `json:"interrupted"`
	Repos       []catalog.Repo    `json:"repos"`
	Groups      []catalog.Group   `json:"groups"`
	Items       []catalog.Item    `json:"items"`
	Paths       map[string]string `json:"paths"`
	// Profile passed on the command line: shown for review, never executed.
	PendingImport *ProfileReview `json:"pendingImport"`
	PendingError  string         `json:"pendingError"`
}

func (a *App) Init() InitData {
	d := InitData{Version: version, System: a.det.Sys, Prefs: a.st.LoadPrefs(), Profiles: a.Profiles(), Repos: a.cat.Repos, Groups: a.cat.Groups, Items: a.cat.Items,
		Paths: map[string]string{"config": a.st.Config, "state": a.st.State, "data": a.st.Data}}
	d.Prefs.ManualProgress = a.cleanManualProgress(d.Prefs.ManualProgress)
	if runs := a.st.Runs(); len(runs) > 0 && runs[0].Status == "interrupted" {
		d.Interrupted = &runs[0]
	}
	if a.pendingPath != "" {
		r, err := a.ReviewProfileFile(a.pendingPath)
		if err != nil {
			d.PendingError = err.Error()
		} else {
			d.PendingImport = r
		}
		a.pendingPath = ""
	}
	return d
}

type Attribution struct {
	Name    string `json:"name"`
	License string `json:"license"`
	URL     string `json:"url"`
	Text    string `json:"text"`
}

type AboutInfo struct {
	Version      string        `json:"version"`
	Author       string        `json:"author"`
	GoVersion    string        `json:"goVersion"`
	WailsVersion string        `json:"wailsVersion"`
	Attributions []Attribution `json:"attributions"`
}

// About is read-only: the version comes from the executable being used.
func (a *App) About() (AboutInfo, error) {
	info := AboutInfo{Version: version, Author: "Paulo Silva", GoVersion: goruntime.Version()}
	info.Attributions = append(info.Attributions, Attribution{Name: "Meu Ubuntu", License: "GNU GPLv3", URL: "https://www.gnu.org/licenses/gpl-3.0.html", Text: projectLicense})
	if build, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range build.Deps {
			if dep.Path == "github.com/wailsapp/wails/v2" {
				info.WailsVersion = dep.Version
			}
		}
	}
	for _, source := range []struct{ name, license, url, file string }{
		{"Devicon", "MIT", "https://github.com/devicons/devicon", "DEVICON-LICENSE.txt"},
		{"Simple Icons", "CC0 1.0", "https://github.com/simple-icons/simple-icons", "SIMPLE-ICONS-LICENSE.txt"},
	} {
		text, err := logoLicenses.ReadFile("frontend/src/assets/apps/" + source.file)
		if err != nil {
			return AboutInfo{}, err
		}
		info.Attributions = append(info.Attributions, Attribution{Name: source.name, License: source.license, URL: source.url, Text: string(text)})
	}
	for _, l := range diag.ExportLicenses() {
		info.Attributions = append(info.Attributions, Attribution{Name: l.Name, License: l.License, URL: l.URL, Text: l.Text})
	}
	return info, nil
}

// SecondInstance is called when the app is opened again (e.g. with a JSON):
// the existing window reviews the profile instead of starting a second app.
func (a *App) SecondInstance(args []string) {
	if a.ctx == nil {
		return
	}
	runtime.WindowUnminimise(a.ctx)
	runtime.Show(a.ctx)
	if path := profileArg(args); path != "" {
		r, err := a.ReviewProfileFile(path)
		if err != nil {
			runtime.EventsEmit(a.ctx, "profile:review-error", err.Error())
			return
		}
		runtime.EventsEmit(a.ctx, "profile:review", r)
	}
}

// profileArg returns the first *.json argument (or --perfil=<file>).
func profileArg(args []string) string {
	for _, a := range args {
		if v, ok := strings.CutPrefix(a, "--perfil="); ok {
			return v
		}
		if strings.HasSuffix(strings.ToLower(a), ".json") && !strings.HasPrefix(a, "-") {
			return a
		}
	}
	return ""
}

// States re-detects the real system state for every item.
func (a *App) States(sel map[string]catalog.Selection) map[string]detect.ItemState {
	ctx := context.Background()
	a.det.Refresh(ctx)
	out := map[string]detect.ItemState{}
	for i := range a.cat.Items {
		it := &a.cat.Items[i]
		s, err := a.cat.ValidateSelection(it.ID, sel[it.ID])
		if err != nil {
			s, _ = a.cat.ValidateSelection(it.ID, catalog.Selection{})
		}
		out[it.ID] = a.det.State(ctx, it, s)
	}
	return out
}

// ItemState refreshes a changed option without rescanning the entire catalog.
func (a *App) ItemState(id string, sel catalog.Selection) (detect.ItemState, error) {
	v, err := a.cat.ValidateSelection(id, sel)
	if err != nil {
		return detect.ItemState{}, err
	}
	return a.det.State(context.Background(), a.cat.Item(id), v), nil
}

func (a *App) SavePrefs(p store.Prefs) error {
	var rev []string
	for _, id := range p.Reverts {
		if a.cat.Item(id) != nil {
			rev = append(rev, id)
		}
	}
	p.Reverts = rev
	for id := range p.Origins {
		if _, ok := p.Selections[id]; !ok {
			delete(p.Origins, id)
		}
	}
	clean := map[string]catalog.Selection{}
	for id, s := range p.Selections {
		if v, err := a.cat.ValidateSelection(id, s); err == nil {
			clean[id] = v
		}
	}
	p.Selections = clean
	var excluded []string
	seen := map[string]bool{}
	for _, id := range p.ExcludedInstalled {
		if a.cat.Item(id) != nil && !seen[id] {
			if _, chosen := clean[id]; !chosen {
				excluded = append(excluded, id)
				seen[id] = true
			}
		}
	}
	p.ExcludedInstalled = excluded
	p.ManualProgress = a.cleanManualProgress(p.ManualProgress)
	return a.st.SavePrefs(p)
}

// Checklist marks are personal notes, never installation evidence. Bind them to
// the exact curated instructions so an updated guide cannot inherit stale marks.
func (a *App) cleanManualProgress(progress map[string][]string) map[string][]string {
	clean := map[string][]string{}
	for id, steps := range progress {
		it := a.cat.Item(id)
		if it == nil || it.Manual == nil {
			continue
		}
		marked := map[string]bool{}
		for _, step := range steps {
			marked[step] = true
		}
		for _, step := range it.Manual.Steps {
			if marked[step] {
				clean[id] = append(clean[id], step)
				delete(marked, step)
			}
		}
	}
	if len(clean) == 0 {
		return nil
	}
	return clean
}

// BuildPlan plans the selection plus explicit reverts. Vendor scripts are
// downloaded now so the user reviews the exact content (SHA-256) they confirm.
func (a *App) BuildPlan(sel map[string]catalog.Selection, reverts []string) (*plan.Plan, error) {
	p, err := plan.Build(context.Background(), a.cat, a.det, sel, reverts)
	if err != nil {
		return nil, err
	}
	if err := a.eng.ReviewUndo(context.Background(), p); err != nil {
		return nil, err
	}
	a.eng.ReviewRecovery(p)
	for i := range p.Steps {
		st := &p.Steps[i]
		if st.Op != "user-script" {
			continue
		}
		body, err := a.eng.Fetch(plan.ScriptURL(a.cat.Item(st.Item), catalog.Selection{Options: st.Options}))
		if err != nil {
			st.Detail = append(st.Detail, "Não foi possível transferir o script para revisão: "+err.Error())
			continue
		}
		st.Expect = fmt.Sprintf("%x", sha256.Sum256(body))
		st.Detail = append(st.Detail, fmt.Sprintf("Script revisto: %d bytes, SHA-256 %s", len(body), st.Expect))
	}
	p.ReviewToken = p.Fingerprint()
	return p, nil
}

// ScriptPreview returns the vendor script content for review in the plan.
func (a *App) ScriptPreview(item string, options map[string][]string, expected string) (string, error) {
	it := a.cat.Item(item)
	if it == nil || it.Method != catalog.MethodUserScript {
		return "", errors.New("item sem script")
	}
	sel, err := a.cat.ValidateSelection(item, catalog.Selection{Options: options})
	if err != nil {
		return "", err
	}
	b, err := a.eng.Fetch(plan.ScriptURL(it, sel))
	if err != nil {
		return "", err
	}
	if expected == "" || fmt.Sprintf("%x", sha256.Sum256(b)) != expected {
		return "", errors.New("o script mudou; regenere o plano antes de o rever")
	}
	return string(b), nil
}

// Execute rebuilds the plan from the selection (the UI never sends commands)
// and requires the explicit confirmations shown in the reviewed plan.
func (a *App) Execute(sel map[string]catalog.Selection, reverts []string, confirmed map[string]string, reviewToken string) (string, error) {
	p, err := a.BuildPlan(sel, reverts)
	if err != nil {
		return "", err
	}
	if reviewToken == "" || reviewToken != p.ReviewToken {
		return "", errors.New("o estado do sistema ou o plano mudou; regenere e reveja o plano antes de executar")
	}
	if len(p.Steps) == 0 {
		return "", errors.New("nada para executar: tudo já está no estado pretendido ou é manual")
	}
	return a.eng.Start(p, confirmed)
}

// EndSession opens GNOME's own logout/reboot confirmation dialog (as the user).
func (a *App) EndSession(action string) error {
	if action != "logout" && action != "reboot" {
		return errors.New("ação inválida")
	}
	if a.simulated {
		return errors.New("modo simulação: a sessão não é terminada")
	}
	return exec.Command("gnome-session-quit", "--"+action).Start()
}

func (a *App) Cancel()       { a.eng.Cancel() }
func (a *App) Running() bool { return a.eng.Running() }

func (a *App) Runs() []store.Run { return a.st.Runs() }

func (a *App) RunLog(id string) string {
	if strings.ContainsAny(id, "/.") {
		return ""
	}
	b, _ := os.ReadFile(a.st.LogPath(id))
	if len(b) > 512*1024 {
		b = b[len(b)-512*1024:]
	}
	return string(b)
}

func (a *App) Changes() ([]store.Change, error) { return a.st.ReadChanges() }

func (a *App) Revert(id string, force bool) error {
	return errors.New("reversão direta desativada; reveja a alteração e use RevertReviewed")
}

func (a *App) PreviewRevert(id string) (*engine.RevertPreview, error) {
	return a.eng.PreviewRevert(context.Background(), id)
}

func (a *App) RevertReviewed(id string, force bool, token string) error {
	return a.eng.RevertReviewed(context.Background(), id, force, token)
}

// ---- Profiles ----

func (a *App) Profiles() []store.Profile {
	return append(profiles.Builtin(), a.st.ListProfiles(a.cat)...)
}

// SaveProfile stores a personal profile. With overwrite=false an existing
// profile with the same ID is reported (store.ErrExists) and kept.
func (a *App) SaveProfile(p store.Profile, overwrite bool) error {
	p.Schema, p.Version = store.ProfileSchema, store.ProfileVersion
	if p.ID == "" {
		p.ID = store.Slug(p.Name)
	}
	if _, err := store.ValidateProfile(a.cat, &p); err != nil {
		return err
	}
	return a.st.SaveProfile(p, store.SourceMine, overwrite)
}

// SaveImported keeps a reviewed imported profile in "Perfis importados".
func (a *App) SaveImported(p store.Profile, overwrite bool) error {
	p.Schema, p.Version = store.ProfileSchema, store.ProfileVersion
	if _, err := store.ValidateProfile(a.cat, &p); err != nil {
		return err
	}
	return a.st.SaveProfile(p, store.SourceImported, overwrite)
}

func (a *App) DeleteProfile(source, id string) error { return a.st.DeleteProfile(source, id) }

func (a *App) ExportProfile(p store.Profile) (string, error) {
	p.Schema, p.Version = store.ProfileSchema, store.ProfileVersion
	if p.ID == "" {
		p.ID = store.Slug(p.Name)
	}
	if _, err := store.ValidateProfile(a.cat, &p); err != nil {
		return "", err
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{Title: "Exportar perfil", DefaultFilename: p.ID + ".json",
		Filters: []runtime.FileFilter{{DisplayName: "Perfil JSON", Pattern: "*.json"}}})
	if err != nil || path == "" {
		return "", err
	}
	b, err := store.MarshalExport(p)
	if err != nil {
		return "", err
	}
	return path, os.WriteFile(path, b, 0o644)
}

type ReviewItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Area      string `json:"area"`
	Status    string `json:"status"`
	Detail    string `json:"detail"`
	Sensitive string `json:"sensitive,omitempty"`
}

// ProfileReview is what the user sees before a profile touches the selection.
type ProfileReview struct {
	Profile     store.Profile `json:"profile"`
	Warnings    []string      `json:"warnings"`
	Items       []ReviewItem  `json:"items"`
	Requirement []string      `json:"requirement"` // requirement mismatches with this machine
	ExistsMine  bool          `json:"existsMine"`
	ExistsImp   bool          `json:"existsImported"`
	Path        string        `json:"path,omitempty"`
}

func (a *App) review(p *store.Profile, warnings []string) *ProfileReview {
	r := &ProfileReview{Profile: *p, Warnings: warnings}
	ctx := context.Background()
	a.det.Refresh(ctx)
	for id, sel := range p.Selections {
		it := a.cat.Item(id)
		st := a.det.State(ctx, it, sel)
		r.Items = append(r.Items, ReviewItem{ID: id, Name: it.Name, Area: it.Area, Status: st.Status, Detail: st.Detail, Sensitive: a.cat.SensitiveReason(it)})
	}
	sort.Slice(r.Items, func(i, j int) bool { return r.Items[i].Name < r.Items[j].Name })
	sys := a.det.Sys
	if v := p.Requires.MinRelease; v != "" && sys.Release != "" && sys.Release < v {
		r.Requirement = append(r.Requirement, "O perfil pede Ubuntu "+v+" ou superior; este sistema é "+sys.Release)
	}
	if len(p.Requires.Arch) > 0 && !slices.Contains(p.Requires.Arch, sys.Arch) {
		r.Requirement = append(r.Requirement, "O perfil pede arquitetura "+strings.Join(p.Requires.Arch, "/")+"; este sistema é "+sys.Arch)
	}
	for _, x := range a.st.ListProfiles(a.cat) {
		if x.ID == p.ID {
			r.ExistsMine = r.ExistsMine || x.Source == store.SourceMine
			r.ExistsImp = r.ExistsImp || x.Source == store.SourceImported
		}
	}
	return r
}

// ReviewProfile compares any profile (built-in, mine, imported) with this PC.
func (a *App) ReviewProfile(p store.Profile) (*ProfileReview, error) {
	p.Schema, p.Version = store.ProfileSchema, store.ProfileVersion
	w, err := store.ValidateProfile(a.cat, &p)
	if err != nil {
		return nil, err
	}
	return a.review(&p, w), nil
}

// ImportProfile opens a file dialog and returns a review. Nothing is saved,
// selected or executed until the user decides.
func (a *App) ImportProfile() (*ProfileReview, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "Importar perfil",
		Filters: []runtime.FileFilter{{DisplayName: "Perfil JSON", Pattern: "*.json"}}})
	if err != nil || path == "" {
		return nil, err
	}
	return a.ReviewProfileFile(path)
}

func (a *App) ReviewProfileFile(path string) (*ProfileReview, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() || fi.Size() > 1<<20 {
		return nil, errors.New("não é um ficheiro de perfil válido (máx. 1 MB)")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, w, err := store.ParseProfile(a.cat, b)
	if err != nil {
		return nil, err
	}
	p.Source = store.SourceImported
	r := a.review(p, w)
	r.Path = path
	return r, nil
}

// MergeSelections combines a profile into the current selection; conflicts
// are returned for the user to resolve and existing choices are kept meanwhile.
func (a *App) MergeSelections(cur map[string]catalog.Selection, origins map[string][]string, incoming map[string]catalog.Selection, label string) selection.Result {
	return selection.Merge(a.cat, cur, origins, incoming, label)
}

// ---- Inventory ----

func (a *App) ParseInventory(text string) []inventory.Entry { return inventory.Parse(text, a.cat) }

func (a *App) OpenInventoryFile() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "Abrir lista de pacotes (apt-mark showmanual, snap list, flatpak list)"})
	if err != nil || path == "" {
		return "", err
	}
	st, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if st.Size() > 8<<20 {
		return "", errors.New("ficheiro demasiado grande")
	}
	b, err := os.ReadFile(path)
	return string(b), err
}

// ---- Diagnostics / migration ----

type FirefoxInfo struct {
	Roots   []migrate.FFRoot `json:"roots"`
	Running bool             `json:"running"`
}

func (a *App) Firefox() FirefoxInfo {
	return FirefoxInfo{Roots: migrate.Detect(a.home), Running: migrate.Running()}
}
func (a *App) FirefoxBackup(profile string) (string, error) { return migrate.Backup(a.home, profile) }
func (a *App) FirefoxMigrate(profile string) (string, error) {
	return migrate.MigrateToFlatpak(a.home, profile)
}

type SSHKey struct {
	File        string `json:"file"`
	Fingerprint string `json:"fingerprint"`
}

// SSHKeys lists public keys only (fingerprints); private keys are never read.
func (a *App) SSHKeys() []SSHKey {
	var out []SSHKey
	files, _ := filepath.Glob(filepath.Join(a.home, ".ssh", "*.pub"))
	for _, f := range files {
		fp, err := a.r.Output(context.Background(), "ssh-keygen", "-lf", f)
		if err != nil {
			fp = "(não foi possível ler)"
		}
		out = append(out, SSHKey{File: f, Fingerprint: strings.TrimSpace(fp)})
	}
	return out
}
