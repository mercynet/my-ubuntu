package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/sys"
)

func TestPartialWallpaperUndoDoesNotResetAlreadyRestoredKey(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	light, dark := "org.gnome.desktop.background picture-uri", "org.gnome.desktop.background picture-uri-dark"
	f.sim.Settings[light], f.sim.Settings[dark] = "'file:///personal-light.png'", "'file:///personal-dark.png'"
	f.run(t, map[string]catalog.Selection{"wallpaper": {}})
	appliedDark := f.sim.Settings[dark]
	f.sim.Settings[dark] = "'file:///later-choice.png'"
	if err := f.eng.undo(ctx, "first", f.cat.Item("wallpaper"), func(string) {}); err == nil {
		t.Fatal("must refuse the later dark choice")
	}
	if f.sim.Settings[light] != "'file:///personal-light.png'" {
		t.Fatal("first key was not restored")
	}
	f.sim.Settings[dark] = appliedDark
	if err := f.eng.undo(ctx, "retry", f.cat.Item("wallpaper"), func(string) {}); err != nil {
		t.Fatal(err)
	}
	if f.sim.Settings[light] != "'file:///personal-light.png'" || f.sim.Settings[dark] != "'file:///personal-dark.png'" {
		t.Fatalf("retry replaced personal settings: %+v", f.sim.Settings)
	}
}

func TestNoHistoryUndoBindsCurrentValueAndRefusesStaleExecution(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	key := "org.gnome.desktop.interface color-scheme"
	f.sim.Settings[key] = "'prefer-dark'"
	p, _ := plan.Build(ctx, f.cat, f.det, nil, []string{"color-scheme"})
	if err := f.eng.ReviewUndo(ctx, p); err != nil {
		t.Fatal(err)
	}
	token := p.Fingerprint()
	f.sim.Settings[key] = "'prefer-light'"
	newPlan, _ := plan.Build(ctx, f.cat, f.det, nil, []string{"color-scheme"})
	f.eng.ReviewUndo(ctx, newPlan)
	if token == newPlan.Fingerprint() {
		t.Fatal("token must change with the reviewed value")
	}
	id, err := f.eng.Start(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	waitRC(t, f)
	if f.sim.Settings[key] != "'prefer-light'" {
		t.Fatal("stale undo overwrote the later choice")
	}
	for _, r := range f.st.Runs() {
		if r.ID == id && r.Status != "failed" {
			t.Fatalf("stale undo reported success: %+v", r)
		}
	}
}

type rcRunner struct {
	*sys.Sim
	failPost, failVerify, failReload bool
}

func (r *rcRunner) Stream(ctx context.Context, log func(string), name string, args ...string) error {
	if r.failPost && name == "phpdismod" {
		return errors.New("injected post failure")
	}
	if r.failReload && name == "systemctl" && len(args) > 0 && args[0] == "restart" {
		return errors.New("injected reload failure")
	}
	return r.Sim.Stream(ctx, log, name, args...)
}
func (r *rcRunner) Output(ctx context.Context, name string, args ...string) (string, error) {
	if r.failVerify && name == "php" {
		return "", errors.New("injected PHP verification failure")
	}
	return r.Sim.Output(ctx, name, args...)
}
func rcUseRunner(t *testing.T, f *fixture, r *rcRunner) {
	f.eng.R, f.det.R = r, r
	root := t.TempDir()
	f.det.Root = root
	f.eng.NewAdmin = func() helper.Admin {
		return helper.InProcess{Env: &helper.Env{R: r, Cat: f.cat, Root: root, Codename: "resolute", Distribution: "ubuntu", UID: "1000"}}
	}
}
func waitRC(t *testing.T, f *fixture) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for f.eng.Running() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if f.eng.Running() {
		t.Fatal("run did not finish")
	}
}

func TestXdebugPostFailureRetriesConfigurationWithoutReinstalling(t *testing.T) {
	f := setup(t)
	r := &rcRunner{Sim: f.sim, failPost: true}
	rcUseRunner(t, f, r)
	sel := map[string]catalog.Selection{"xdebug": {}}
	_, first := f.run(t, sel)
	if first.Status != "failed" {
		t.Fatalf("expected post failure: %+v", first)
	}
	r.failPost = false
	p, next := f.run(t, sel)
	for _, s := range p.Steps {
		if s.Op == "apt-install" {
			t.Fatalf("retry reinstalled a present package: %+v", s)
		}
	}
	if next.Status != "done" || !f.sim.XdebugConfigured("8.5") {
		t.Fatalf("configuration was not recovered: %+v", next)
	}
	p, _ = plan.Build(context.Background(), f.cat, f.det, sel, nil)
	if len(p.Steps) != 0 {
		t.Fatalf("complete recipe must become a no-op: %+v", p.Steps)
	}
}

func TestFailedJournaldReloadIsNotDetectedAsComplete(t *testing.T) {
	f := setup(t)
	r := &rcRunner{Sim: f.sim, failReload: true}
	rcUseRunner(t, f, r)
	sel := map[string]catalog.Selection{"journald-limit": {}}
	_, first := f.run(t, sel)
	if first.Status != "failed" {
		t.Fatalf("expected reload failure: %+v", first)
	}
	chosen, _ := f.cat.ValidateSelection("journald-limit", catalog.Selection{})
	if st := f.det.State(context.Background(), f.cat.Item("journald-limit"), chosen); st.Status != detect.NeedsConfig {
		t.Fatalf("file must not prove runtime success: %+v", st)
	}
	r.failReload = false
	_, next := f.run(t, sel)
	if next.Status != "done" {
		t.Fatalf("reload retry failed: %+v", next)
	}
	if st := f.det.State(context.Background(), f.cat.Item("journald-limit"), chosen); st.Status != detect.Installed {
		t.Fatalf("successful reload not confirmed: %+v", st)
	}
	if len(f.st.Changes()) != 1 {
		t.Fatal("retry must preserve the original journal, not create another rollback baseline")
	}
	p, _ := plan.Build(context.Background(), f.cat, f.det, nil, []string{"journald-limit"})
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	waitRC(t, f)
	for _, c := range f.st.Changes() {
		if c.Reverted == "" {
			t.Fatal("successful retry could not be fully reverted")
		}
	}
}

func TestVerificationInterventionPersistsAndCanBeRetriedReadOnly(t *testing.T) {
	f := setup(t)
	r := &rcRunner{Sim: f.sim, failVerify: true}
	rcUseRunner(t, f, r)
	sel := map[string]catalog.Selection{"php": {}}
	_, first := f.run(t, sel)
	if first.Status != StIntervention {
		t.Fatalf("verification must require intervention: %+v", first)
	}
	found := false
	for _, s := range first.Steps {
		found = found || strings.Contains(s.Note, "injected PHP verification failure")
	}
	if !found {
		t.Fatal("verification explanation missing from saved history")
	}
	b, _ := os.ReadFile(first.LogFile)
	if !strings.Contains(string(b), "injected PHP verification failure") {
		t.Fatal("verification error missing from log")
	}
	r.failVerify = false
	p, _ := plan.Build(context.Background(), f.cat, f.det, sel, nil)
	f.eng.ReviewRecovery(p)
	if len(p.Steps) != 1 || p.Steps[0].Op != "verify" || p.NeedsRoot {
		t.Fatalf("retry must be read-only: %+v", p)
	}
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	waitRC(t, f)
	p, _ = plan.Build(context.Background(), f.cat, f.det, sel, nil)
	f.eng.ReviewRecovery(p)
	if len(p.Steps) != 0 {
		t.Fatal("successful verification was not cleared")
	}
}

type rcGroups struct{ *sys.Sim }

func (r rcGroups) Output(ctx context.Context, name string, args ...string) (string, error) {
	if name == "id" {
		if len(args) > 1 {
			return "user docker\n", nil
		}
		return "user\n", nil
	}
	return r.Sim.Output(ctx, name, args...)
}
func TestDockerGroupUsesPersistentMembershipBeforeNewLogin(t *testing.T) {
	f := setup(t)
	f.det.R = rcGroups{f.sim}
	st := f.det.State(context.Background(), f.cat.Item("docker-group"), catalog.Selection{})
	if st.Status != detect.Installed || !strings.Contains(st.Detail, "termine sessão") {
		t.Fatalf("persistent/session memberships confused: %+v", st)
	}
}

type rcGitAbsent struct{ *sys.Sim }

func (r rcGitAbsent) Output(ctx context.Context, name string, args ...string) (string, error) {
	if name == "git" {
		return "", fmt.Errorf("git: %w: absent", sys.ExitStatus(1))
	}
	return r.Sim.Output(ctx, name, args...)
}
func TestGitAbsenceHandlesWrappedExitCode(t *testing.T) {
	f := setup(t)
	f.eng.R = rcGitAbsent{f.sim}
	if value, err := f.eng.gitValue(context.Background(), "test.absent"); err != nil || value != "" {
		t.Fatalf("wrapped exit 1 is absence: %q %v", value, err)
	}
}

func TestNormalCloseWaitsForRunningOperationAndPreventsNewRuns(t *testing.T) {
	f := setup(t)
	p, _ := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"htop": {}}, nil)
	if _, err := f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	if !f.eng.PreventClose() {
		t.Fatal("close must be prevented while running")
	}
	f.eng.Cancel()
	waitRC(t, f)
	if f.eng.PreventClose() {
		t.Fatal("idle app should close")
	}
	if _, err := f.eng.Start(p, nil); err == nil {
		t.Fatal("must not start after closing has begun")
	}
}

func TestTemplateUndoReviewUsesActualTemplatesDirectory(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"tpl-md": {}})
	p, _ := plan.Build(context.Background(), f.cat, f.det, nil, []string{"tpl-md"})
	if err := f.eng.ReviewUndo(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if len(p.Steps) != 1 || !strings.HasPrefix(p.Steps[0].Undo[0].Key, filepath.Clean(detect.TemplatesDir(f.eng.Home))) {
		t.Fatalf("wrong review path: %+v", p)
	}
}
