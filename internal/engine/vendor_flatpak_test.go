package engine

import (
	"context"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/plan"
)

func TestNavicatReferenceRequiresReviewAndInstallsInSimulation(t *testing.T) {
	f := setup(t)
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"navicat-premium-lite": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	confirmed := map[string]string{}
	found := false
	for _, s := range p.Steps {
		if s.Item == "navicat-premium-lite" && s.Op == "flatpak-install" {
			if s.Expect == "" || s.Confirm == "" {
				t.Fatal("vendor reference is not bound to confirmation")
			}
			confirmed[s.Item] = s.Expect
			found = true
		}
	}
	if !found {
		t.Fatal("no vendor install step")
	}
	if _, err = f.eng.Start(p, nil); err == nil {
		t.Fatal("terms accepted implicitly")
	}
	if _, err = f.eng.Start(p, confirmed); err != nil {
		t.Fatal(err)
	}
	awaitEngine(t, f.eng)
	if got := f.st.Runs()[0].Status; got != "done" {
		t.Fatalf("installation: %s %+v", got, f.st.Runs()[0].Steps)
	}
}
