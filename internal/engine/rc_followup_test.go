package engine

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/sys"
)

func TestJournaldReapplyAfterRevertCannotReuseOldReceipt(t *testing.T) {
	f := setup(t)
	r := &rcRunner{Sim: f.sim}
	rcUseRunner(t, f, r)
	ctx := context.Background()
	sel := map[string]catalog.Selection{"journald-limit": {}}
	_, run := f.run(t, sel)
	if run.Status != "done" {
		t.Fatal(run.Status)
	}
	p, _ := plan.Build(ctx, f.cat, f.det, nil, []string{"journald-limit"})
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	waitRC(t, f)
	r.failReload = true
	_, run = f.run(t, sel)
	if run.Status != "failed" {
		t.Fatalf("expected reload failure: %+v", run)
	}
	choice, _ := f.cat.ValidateSelection("journald-limit", catalog.Selection{})
	if st := f.det.State(ctx, f.cat.Item("journald-limit"), choice); st.Status != detect.NeedsConfig {
		t.Fatalf("old receipt concealed failed reload: %+v", st)
	}
	r.failReload = false
	_, run = f.run(t, sel)
	if run.Status != "done" {
		t.Fatal("failed to repair reload")
	}
}

func TestRevertedLatestTemplateJournalCannotFallBackToOlderActiveChange(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	sel := map[string]catalog.Selection{"tpl-md": {}}
	f.run(t, sel)
	c1 := f.st.Changes()[0]
	os.Remove(c1.Key)
	f.run(t, sel)
	c2 := f.st.Changes()[1]
	if err := f.eng.Revert(ctx, c2, false); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(c1.Key, []byte(c1.After), 0600)
	p, _ := plan.Build(ctx, f.cat, f.det, nil, []string{"tpl-md"})
	if err := f.eng.ReviewUndo(ctx, p); err != nil {
		t.Fatal(err)
	}
	if p.Steps[0].Undo[0].Mode != "skip" {
		t.Fatal("review must keep the already restored path")
	}
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	waitRC(t, f)
	if _, err := os.Stat(c1.Key); err != nil {
		t.Fatal("older journal deleted a file that the review kept")
	}
}

type rcDriver struct {
	*sys.Sim
	fail bool
}

func (r rcDriver) Output(ctx context.Context, name string, args ...string) (string, error) {
	if name == "ubuntu-drivers" {
		if r.fail {
			return "", sys.ExitStatus(1)
		}
		return "driver   : nvidia-driver-590-open - distro non-free recommended\n", nil
	}
	return r.Sim.Output(ctx, name, args...)
}
func TestDriverStateRequiresRecommendedPackageAndReportsQueryFailure(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	f.det.R = rcDriver{Sim: f.sim}
	f.det.Refresh(ctx)
	it := f.cat.Item("ubuntu-drivers")
	if s := f.det.State(ctx, it, catalog.Selection{}); s.Status != detect.Available {
		t.Fatalf("missing driver: %+v", s)
	}
	f.sim.Packages["nvidia-driver-590-open"] = true
	f.det.Refresh(ctx)
	if s := f.det.State(ctx, it, catalog.Selection{}); s.Status != detect.Installed {
		t.Fatalf("installed driver: %+v", s)
	}
	f.det.R = rcDriver{Sim: f.sim, fail: true}
	if s := f.det.State(ctx, it, catalog.Selection{}); s.Status != detect.Unknown {
		t.Fatalf("failed query must not confirm success: %+v", s)
	}
}

func TestCancelledVerificationRemainsPendingAndPrecedesDependants(t *testing.T) {
	f := setup(t)
	r := &rcRunner{Sim: f.sim, failVerify: true}
	rcUseRunner(t, f, r)
	ctx := context.Background()
	sel := map[string]catalog.Selection{"php": {}}
	f.run(t, sel)
	p, _ := plan.Build(ctx, f.cat, f.det, sel, nil)
	f.eng.ReviewRecovery(p)
	f.eng.cancel.Store(true)
	// Save an actual cancelled retry without relying on goroutine scheduling.
	run := f.st.Runs()[0]
	run.ID = "zz-cancelled-retry"
	run.Status = "cancelled"
	run.Steps = run.Steps[len(run.Steps)-1:]
	run.Steps[0].Status = StCancelled
	run.Steps[0].Op = "verify"
	if err := f.st.SaveRun(run); err != nil {
		t.Fatal(err)
	}
	sel["composer"] = catalog.Selection{}
	p, _ = plan.Build(ctx, f.cat, f.det, sel, nil)
	f.eng.ReviewRecovery(p)
	verification, dependant := -1, -1
	for i, s := range p.Steps {
		if s.Item == "php" && s.Op == "verify" {
			verification = i
		}
		if s.Item == "composer" {
			dependant = i
		}
	}
	if verification < 0 || dependant < 0 || verification >= dependant {
		t.Fatalf("pending verification lost or after dependant: %+v", p.Steps)
	}
}

func TestTemplateDirectoryChangeCannotDeleteUnreviewedOldTemplate(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	f.run(t, map[string]catalog.Selection{"tpl-md": {}})
	old := filepath.Join(detect.TemplatesDir(f.eng.Home), filepath.Base(f.cat.Item("tpl-md").Template.Name))
	dir := filepath.Join(f.eng.Home, "new-templates")
	os.MkdirAll(dir, 0755)
	os.MkdirAll(filepath.Join(f.eng.Home, ".config"), 0755)
	os.WriteFile(filepath.Join(f.eng.Home, ".config/user-dirs.dirs"), []byte("XDG_TEMPLATES_DIR=\""+dir+"\"\n"), 0600)
	p, _ := plan.Build(ctx, f.cat, f.det, nil, []string{"tpl-md"})
	if err := f.eng.ReviewUndo(ctx, p); err != nil {
		t.Fatal(err)
	}
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	waitRC(t, f)
	if _, err := os.Stat(old); err != nil {
		t.Fatal("deleted the old template even though the review kept it")
	}
}

func TestAdministrativeBatchReviewBindsPrivateSnapshotTarget(t *testing.T) {
	f := setup(t)
	r := &rcRunner{Sim: f.sim}
	rcUseRunner(t, f, r)
	ctx := context.Background()
	for _, v := range []string{"1G", "2G"} {
		_, run := f.run(t, map[string]catalog.Selection{"journald-limit": {Options: map[string][]string{"value": {v}}}})
		if run.Status != "done" {
			t.Fatal(run.Status)
		}
	}
	p, _ := plan.Build(ctx, f.cat, f.det, nil, []string{"journald-limit"})
	if err := f.eng.ReviewUndo(ctx, p); err != nil {
		t.Fatal(err)
	}
	originalToken := p.Fingerprint()
	file := filepath.Join(f.det.Root, "var/lib/meu-ubuntu/journald.json")
	raw, _ := os.ReadFile(file)
	var snapshot map[string]any
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		t.Fatal(err)
	}
	oldest := snapshot
	for oldest["Previous"] != nil {
		oldest = oldest["Previous"].(map[string]any)
	}
	oldest["HadValue"] = true
	oldest["Before"] = base64.StdEncoding.EncodeToString([]byte("[Journal]\nSystemMaxUse=333M\n"))
	raw, _ = json.Marshal(snapshot)
	os.WriteFile(file, raw, 0600)
	updated, _ := plan.Build(ctx, f.cat, f.det, nil, []string{"journald-limit"})
	if err := f.eng.ReviewUndo(ctx, updated); err != nil {
		t.Fatal(err)
	}
	if originalToken == updated.Fingerprint() {
		t.Fatal("private target change did not invalidate review")
	}
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	waitRC(t, f)
	config, _ := os.ReadFile(filepath.Join(f.det.Root, "etc/systemd/journald.conf.d/60-meu-ubuntu.conf"))
	if !strings.Contains(string(config), "2G") {
		t.Fatal("stale administrative preview changed the file")
	}
	for _, c := range f.st.Changes() {
		if c.Reverted != "" {
			t.Fatal("stale administrative preview marked a change reverted")
		}
	}
}
