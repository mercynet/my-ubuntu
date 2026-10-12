package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/plan"
)

func TestUndoPreservesLaterUserChange(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"dock-position": {}})
	key := "org.gnome.shell.extensions.dash-to-dock dock-position"
	f.sim.Settings[key] = "'RIGHT'"
	if err := f.eng.undo(context.Background(), "probe", f.cat.Item("dock-position"), func(string) {}); err == nil {
		t.Fatal("later user change must prevent undo")
	}
	if got := f.sim.Settings[key]; got != "'RIGHT'" {
		t.Fatalf("later user change overwritten: %s", got)
	}
}

func TestPowerRevertPreservesLaterUserChange(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"power-performance": {}})
	changes := f.st.Changes()
	f.sim.Profile = "power-saver"
	if err := f.eng.Revert(context.Background(), changes[0], false); err == nil {
		t.Fatal("later power change must prevent undo")
	}
	if f.sim.Profile != "power-saver" {
		t.Fatal("later power change overwritten")
	}
}

func TestCustomOptionHasRevertPlan(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"dock-position": {Options: map[string][]string{"pos": {"'RIGHT'"}}}})
	p, err := plan.Build(context.Background(), f.cat, f.det, nil, []string{"dock-position"})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Steps) != 1 || p.Steps[0].Op != "undo" {
		t.Fatal("custom option must have an undo step")
	}
	if len(f.st.Changes()) != 1 {
		t.Fatal("journal must have a change that can be reverted")
	}
}

func TestRapidRunsHaveSeparateHistory(t *testing.T) {
	f := setup(t)
	// Rapid runs must retain independent history and logs.

	_, first := f.run(t, map[string]catalog.Selection{"dock-position": {}})
	_, second := f.run(t, map[string]catalog.Selection{"color-scheme": {}})
	if first.ID == second.ID {
		t.Fatal("run IDs must be unique")
	}
	if len(f.st.Runs()) != 2 {
		t.Fatal("both runs must be retained")
	}
}

func TestJournalFailurePreventsMutation(t *testing.T) {
	f := setup(t)
	os.Mkdir(filepath.Join(f.st.State, "changes.json"), 0700)
	f.run(t, map[string]catalog.Selection{"dock-position": {}})
	if f.sim.Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] != "'LEFT'" {
		t.Fatal("journal failure must prevent mutation")
	}
	if len(f.st.Changes()) != 0 {
		t.Fatal("failed mutation must not create a journal entry")
	}
}

func TestRunRecordFailurePreventsStarting(t *testing.T) {
	f := setup(t)
	if err := os.Remove(filepath.Join(f.st.State, "runs")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.st.State, "runs"), []byte("unavailable"), 0600); err != nil {
		t.Fatal(err)
	}
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"dock-position": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.eng.Start(p, nil); err == nil {
		t.Fatal("missing history storage must prevent execution")
	}
	if f.eng.Running() {
		t.Fatal("engine must remain available after startup failure")
	}
	if f.sim.Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] != "'LEFT'" {
		t.Fatal("setting must remain unchanged")
	}
}

func TestItemRevertRestoresConsecutiveAdministrativeChanges(t *testing.T) {
	f := setup(t)
	f.det.Root = filepath.Join(filepath.Dir(f.eng.Home), "root")
	for _, value := range []string{"1G", "2G"} {
		_, run := f.run(t, map[string]catalog.Selection{"journald-limit": {Options: map[string][]string{"value": {value}}}})
		if run.Status != "done" {
			t.Fatalf("apply failed: %+v", run.Steps)
		}
	}
	p, err := plan.Build(context.Background(), f.cat, f.det, nil, []string{"journald-limit"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	for f.eng.Running() {
		time.Sleep(time.Millisecond)
	}
	for _, c := range f.st.Changes() {
		if c.Reverted == "" {
			t.Fatal("every administrative change must be reverted")
		}
	}
	if _, err := os.Stat(filepath.Join(f.det.Root, "etc/systemd/journald.conf.d/60-meu-ubuntu.conf")); !os.IsNotExist(err) {
		t.Fatal("original absence must be restored")
	}
}

func TestPlanRefreshesExternalPackageChanges(t *testing.T) {
	f := setup(t)
	// External installation after the last Refresh (the full catalog recipe).
	installed := f.cat.Item("docker-ubuntu")
	choice, _ := f.cat.ValidateSelection(installed.ID, catalog.Selection{})
	for _, pkg := range f.cat.ExpandPackages(installed, choice) {
		f.sim.Packages[pkg] = true
	}
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"docker-engine": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Blocked {
		t.Fatal("planning must detect the externally installed conflict")
	}
	f.det.Refresh(context.Background())
	fresh, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"docker-engine": {}}, nil)
	if !fresh.Blocked {
		t.Fatal("refresh must reveal the conflict")
	}
}

type extensionRunner struct{ base *fixture }

func (r extensionRunner) LookPath(name string) bool { return r.base.sim.LookPath(name) }
func (r extensionRunner) Stream(ctx context.Context, log func(string), name string, args ...string) error {
	return r.base.sim.Stream(ctx, log, name, args...)
}
func (r extensionRunner) Output(ctx context.Context, name string, args ...string) (string, error) {
	if name == "gnome-extensions" {
		if len(args) == 1 {
			return "caffeine@patapon.info\n", nil
		}
		return "", nil
	}
	return r.base.sim.Output(ctx, name, args...)
}
func TestDisabledExtensionPlansActivation(t *testing.T) {
	f := setup(t)
	f.det.R = extensionRunner{f}
	f.det.Refresh(context.Background())
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"ext-caffeine": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Steps) != 1 || p.Steps[0].Op != "extension-enable" {
		t.Fatal("disabled extension must be activated without reinstalling")
	}
	if f.det.State(context.Background(), f.cat.Item("ext-caffeine"), catalog.Selection{}).Status != detect.NeedsConfig {
		t.Fatal("not detected as disabled")
	}
}
