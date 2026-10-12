package engine

import (
	"context"
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
)

func TestFailureKeepsActionableDiagnosticInHistoryAndEvent(t *testing.T) {
	f := setup(t)
	f.sim.FailPkgs["ncdu"] = true
	_, run := f.run(t, map[string]catalog.Selection{"ncdu": {}})
	failedSteps := 0
	for _, step := range run.Steps {
		if step.Status != StFailed {
			continue
		}
		failedSteps++
		foundEvent := false
		if step.Error != "exit status 100" || step.Problem == nil || step.Problem.Code != "package" {
			t.Fatalf("missing original error or actionable diagnostic: %+v", step)
		}
		f.mu.Lock()
		for _, event := range f.evs {
			if event.StepID == step.ID && event.Status == StFailed && event.Problem != nil && event.Problem.Code == "package" {
				foundEvent = true
			}
		}
		f.mu.Unlock()
		if !foundEvent {
			t.Fatal("diagnostic must also reach the live interface")
		}
	}
	if run.Status != "failed" || failedSteps == 0 {
		t.Fatal("run should fail")
	}
}

func TestFailureCauseNeedsConcreteEvidence(t *testing.T) {
	for _, tc := range []struct{ message, code string }{
		{"exit status 100", "unknown"},
		{"Fingerprint: ABC123; exit status 100", "unknown"},
		{"E: Could not get lock /var/lib/dpkg/lock-frontend", "busy"},
		{"No space left on device", "space"},
		{"Temporary failure resolving archive.ubuntu.com", "network"},
		{"E: The repository is not signed", "signature"},
		{"dpkg was interrupted", "dpkg"},
		{"E: Packages need to be removed but remove is disabled.", "removal-blocked"},
		{"autenticação recusada ou modo administrativo terminou", "authorization"},
	} {
		if p := executionProblem("exit status 100", tc.message); p.Code != tc.code {
			t.Errorf("%q: got %s, want %s", tc.message, p.Code, tc.code)
		}
	}
}

func TestReviewedRevertReadsActualValueAndRequiresUnchangedReview(t *testing.T) {
	f := setup(t)
	f.run(t, map[string]catalog.Selection{"dock-position": {}})
	c := f.st.Changes()[0]
	p, err := f.eng.PreviewRevert(context.Background(), c.ID)
	if err != nil {
		t.Fatal(err)
	}
	key := c.Schema + " " + c.Key
	if p.Current != c.After || p.Target != c.Before || !p.CanRevert || f.sim.Settings[key] != c.After {
		t.Fatalf("preview must read without changing anything: %+v", p)
	}
	if err := f.eng.RevertReviewed(context.Background(), c.ID, false, ""); err == nil {
		t.Fatal("a review is required")
	}
	f.sim.Settings[key] = "'RIGHT'"
	if err := f.eng.RevertReviewed(context.Background(), c.ID, true, p.ReviewToken); err == nil || !strings.Contains(err.Error(), "revisão") {
		t.Fatal("even forcing must reject a stale review")
	}
	if f.sim.Settings[key] != "'RIGHT'" {
		t.Fatal("later value was overwritten")
	}
	changed, err := f.eng.PreviewRevert(context.Background(), c.ID)
	if err != nil {
		t.Fatal(err)
	}
	if changed.CanRevert || !changed.ForceAllowed || changed.Current != "'RIGHT'" {
		t.Fatalf("later value must be shown with an explicit decision: %+v", changed)
	}
	if err := f.eng.RevertReviewed(context.Background(), c.ID, false, changed.ReviewToken); err == nil {
		t.Fatal("overwriting a later value must be explicit")
	}
	if err := f.eng.RevertReviewed(context.Background(), c.ID, true, changed.ReviewToken); err != nil {
		t.Fatal(err)
	}
	if f.sim.Settings[key] != c.Before {
		t.Fatal("reviewed target was not restored")
	}
}
