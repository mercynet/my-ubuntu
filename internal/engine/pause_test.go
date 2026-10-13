package engine

import (
	"context"
	"sync"
	"testing"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/store"
)

type boundaryAdmin struct {
	helper.Admin
	started chan struct{}
	release chan struct{}
	once    sync.Once
}

func (b *boundaryAdmin) Do(ctx context.Context, req helper.Request, log func(string)) error {
	b.once.Do(func() { close(b.started); <-b.release })
	return b.Admin.Do(ctx, req, log)
}

func pauseFixture(t *testing.T) (*fixture, *boundaryAdmin, *plan.Plan) {
	t.Helper()
	f := setup(t)
	b := &boundaryAdmin{Admin: f.eng.NewAdmin(), started: make(chan struct{}), release: make(chan struct{})}
	f.eng.NewAdmin = func() helper.Admin { return b }
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"htop": {}, "color-scheme": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.eng.Cancel()
		select {
		case <-b.release:
		default:
			close(b.release)
		}
		awaitEngine(t, f.eng)
	})
	select {
	case <-b.started:
	case <-time.After(5 * time.Second):
		t.Fatal("operation never started")
	}
	return f, b, p
}

func awaitEngine(t *testing.T, e *Engine) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for e.Running() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if e.Running() {
		t.Error("engine did not stop")
	}
}

func awaitPaused(t *testing.T, f *fixture) store.Run {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		runs := f.st.Runs()
		if len(runs) > 0 && runs[0].Status == "paused" {
			return runs[0]
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("pause was not persisted")
	return store.Run{}
}

func TestPauseFinishesCurrentOperationAndResumeCompletes(t *testing.T) {
	f, b, _ := pauseFixture(t)
	if err := f.eng.Pause(); err != nil {
		t.Fatal(err)
	}
	// The current process must remain alive until explicitly released.
	if !f.eng.Running() {
		t.Fatal("pause killed the active operation")
	}
	close(b.release)
	run := awaitPaused(t, f)
	if run.Steps[0].Status != StDone || run.Steps[1].Status != StPending {
		t.Fatalf("unsafe boundary: %+v", run.Steps)
	}
	if !f.eng.PreventClose() {
		t.Fatal("a paused run must prevent window close")
	}
	if err := f.eng.Resume(); err != nil {
		t.Fatal(err)
	}
	awaitEngine(t, f.eng)
	if got := f.st.Runs()[0].Status; got != "done" {
		t.Fatalf("resume: %s %+v", got, f.st.Runs()[0].Steps)
	}
}

func TestCancelWhilePausedStartsNoMoreOperations(t *testing.T) {
	f, b, _ := pauseFixture(t)
	if err := f.eng.Pause(); err != nil {
		t.Fatal(err)
	}
	close(b.release)
	awaitPaused(t, f)
	f.eng.Cancel()
	awaitEngine(t, f.eng)
	run := f.st.Runs()[0]
	if run.Status != "cancelled" {
		t.Fatalf("status: %s", run.Status)
	}
	for _, step := range run.Steps[1:] {
		if step.Status != StCancelled {
			t.Fatalf("later step started: %+v", step)
		}
	}
}

func TestResumeRefusesChangesMadeDuringPause(t *testing.T) {
	f, b, _ := pauseFixture(t)
	if err := f.eng.Pause(); err != nil {
		t.Fatal(err)
	}
	close(b.release)
	awaitPaused(t, f)
	f.sim.Settings["org.gnome.desktop.interface color-scheme"] = "'prefer-light'"
	if err := f.eng.Resume(); err != nil {
		t.Fatal(err)
	}
	awaitEngine(t, f.eng)
	if run := f.st.Runs()[0]; run.Status != StIntervention {
		t.Fatalf("changed state was accepted: %+v", run)
	}
	if got := f.sim.Settings["org.gnome.desktop.interface color-scheme"]; got != "'prefer-light'" {
		t.Fatalf("external change overwritten: %s", got)
	}
}

func TestResumeBeforeBoundaryCancelsPauseRequest(t *testing.T) {
	f, b, _ := pauseFixture(t)
	if err := f.eng.Pause(); err != nil {
		t.Fatal(err)
	}
	if err := f.eng.Resume(); err != nil {
		t.Fatal(err)
	}
	f.eng.mu.Lock()
	state := f.eng.runState
	f.eng.mu.Unlock()
	if state != "running" {
		t.Fatalf("UI stuck while resuming: %s", state)
	}
	close(b.release)
	awaitEngine(t, f.eng)
	if got := f.st.Runs()[0].Status; got != "done" {
		t.Fatalf("status: %s", got)
	}
}
