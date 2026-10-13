package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

func TestReviewedNoChangePlanSerializesStepsAsArray(t *testing.T) {
	for _, tc := range []struct {
		name string
		sel  map[string]catalog.Selection
	}{
		{"empty selection", nil},
		{"already installed", map[string]catalog.Selection{"git": {}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := simulationApp(t)
			p, err := a.BuildPlan(tc.sel, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(p.Steps) != 0 {
				t.Fatalf("expected no changes, got %+v", p.Steps)
			}
			body, err := json.Marshal(p)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(body, &fields); err != nil {
				t.Fatal(err)
			}
			if string(fields["steps"]) != "[]" {
				t.Fatalf("frontend requires an iterable array, got steps=%s", fields["steps"])
			}
			if p.ReviewToken != p.Fingerprint() {
				t.Fatal("review token does not bind the returned plan")
			}
		})
	}
}

func simulationApp(t *testing.T) *App {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	a, err := NewApp(true)
	if err != nil {
		t.Fatal(err)
	}
	a.r.(*sys.Sim).Delay = 0
	return a
}

func TestExecuteRejectsChangedReview(t *testing.T) {
	a := simulationApp(t)
	sel := map[string]catalog.Selection{"dock-position": {}}
	p, err := a.BuildPlan(sel, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.r.(*sys.Sim).Settings["org.gnome.shell.extensions.dash-to-dock dock-position"] = "'RIGHT'"
	if _, err := a.Execute(sel, nil, nil, p.ReviewToken); err == nil || !strings.Contains(err.Error(), "mudou") {
		t.Fatalf("changed state must require new review: %v", err)
	}
	if a.eng.Running() || len(a.st.Runs()) != 0 {
		t.Fatal("stale review must not start execution")
	}
}

func TestOptionDetectionOnlyQueriesAffectedItem(t *testing.T) {
	a := simulationApp(t)
	a.det.Refresh(context.Background())
	sim := a.r.(*sys.Sim)
	sim.Calls = nil
	st, err := a.ItemState("dock-position", catalog.Selection{Options: map[string][]string{"pos": {"'RIGHT'"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Settings) != 1 || st.Settings[0].Proposed != "'RIGHT'" {
		t.Fatalf("incorrect proposed value: %+v", st)
	}
	if len(sim.Calls) != 1 || !strings.HasPrefix(sim.Calls[0], "gsettings get") {
		t.Fatalf("unnecessary catalog rescan: %v", sim.Calls)
	}
}

func TestScriptPreviewRejectsContentDifferentFromReview(t *testing.T) {
	a := simulationApp(t)
	p, err := a.BuildPlan(map[string]catalog.Selection{"php-herd-lite": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.eng.Fetch = func(string) ([]byte, error) { return []byte("echo changed"), nil }
	s := p.Steps[0]
	if _, err := a.ScriptPreview(s.Item, s.Options, s.Expect); err == nil {
		t.Fatal("preview must not show a different script as reviewed")
	}
}
