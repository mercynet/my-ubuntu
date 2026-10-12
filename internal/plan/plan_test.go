package plan

import (
	"context"
	"encoding/json"
	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/sys"
	"os"
	"path/filepath"
	"testing"
)

func TestInstalledManualToolDoesNotPullMissingPluginDependencies(t *testing.T) {
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	sim := sys.NewSim()
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".oh-my-zsh"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".oh-my-zsh", "oh-my-zsh.sh"), []byte("# existing installation"), 0600); err != nil {
		t.Fatal(err)
	}
	d := &detect.Detector{R: sim, Cat: cat, Sys: detect.ReadSystem(sim, true), Home: home, Root: t.TempDir()}
	p, err := Build(context.Background(), cat, d, map[string]catalog.Selection{"oh-my-zsh": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Steps) != 0 || len(p.Added) != 0 || len(p.Manual) != 0 || len(p.Skipped) != 1 {
		t.Fatalf("automatically including an installed tool must not install its optional plugins: %+v", p)
	}
}

func TestFingerprintBindsReviewedValuesAndReverts(t *testing.T) {
	p := &Plan{Steps: []Step{{ID: "s01", Item: "dock", Detail: []string{"LEFT → RIGHT"}}}, Selection: map[string]catalog.Selection{"dock": {}}}
	a := p.Fingerprint()
	p.ReviewToken = a
	if p.Fingerprint() != a {
		t.Fatal("token must not fingerprint itself")
	}
	p.Steps[0].Detail = []string{"BOTTOM → RIGHT"}
	if p.Fingerprint() == a {
		t.Fatal("current state changes must invalidate review")
	}
	a = p.Fingerprint()
	p.Reverts = []string{"journald-limit"}
	if p.Fingerprint() == a {
		t.Fatal("different reversals must invalidate review")
	}
}

func TestAlreadyInstalledChoicesReturnAnEmptyStepArrayForTheInterface(t *testing.T) {
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	sim := sys.NewSim()
	d := &detect.Detector{R: sim, Cat: cat, Sys: detect.ReadSystem(sim, true), Home: t.TempDir(), Root: t.TempDir()}
	p, err := Build(context.Background(), cat, d, map[string]catalog.Selection{"git": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var response struct {
		Steps []Step `json:"steps"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}
	if response.Steps == nil || len(response.Steps) != 0 {
		t.Fatalf("a completed profile must return steps: [], so reviewing it remains usable: %s", body)
	}
}

func TestExpandedDependencyPlanCanBeExecutedWithoutChangingReview(t *testing.T) {
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	sim := sys.NewSim()
	d := &detect.Detector{R: sim, Cat: cat, Sys: detect.ReadSystem(sim, true), Home: t.TempDir(), Root: t.TempDir()}
	ctx := context.Background()
	reviewed, err := Build(ctx, cat, d, map[string]catalog.Selection{"codex-cli": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviewed.Added) == 0 {
		t.Fatal("test needs an automatically added Node dependency")
	}
	rebuilt, err := Build(ctx, cat, d, reviewed.Selection, reviewed.Reverts)
	if err != nil {
		t.Fatal(err)
	}
	if reviewed.Fingerprint() != rebuilt.Fingerprint() {
		t.Fatal("the expanded selection must preserve the reviewed execution")
	}
	sim.Packages["nodejs"], sim.Packages["npm"] = true, true
	changed, err := Build(ctx, cat, d, reviewed.Selection, reviewed.Reverts)
	if err != nil {
		t.Fatal(err)
	}
	if changed.Fingerprint() == reviewed.Fingerprint() {
		t.Fatal("an actual change to dependency installation must invalidate review")
	}
}
