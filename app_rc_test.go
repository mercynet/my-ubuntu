package main

import (
	"path/filepath"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

func rcApp(t *testing.T) *App {
	t.Helper()
	root := t.TempDir()
	for key, dir := range map[string]string{"XDG_CONFIG_HOME": "config", "XDG_STATE_HOME": "state", "XDG_DATA_HOME": "data"} {
		t.Setenv(key, filepath.Join(root, dir))
	}
	a, err := NewApp(true)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestLegacyBridgeRevertAlwaysRefusesMutation(t *testing.T) {
	a := rcApp(t)
	if err := a.Revert("any", true); err == nil {
		t.Fatal("unreviewed bridge mutation must be refused")
	}
	if len(a.st.Changes()) != 0 || a.eng.Running() {
		t.Fatal("legacy call mutated state")
	}
}

func TestExecuteRefusesChangedBatchUndoReview(t *testing.T) {
	a := rcApp(t)
	r := a.r.(*sys.Sim)
	key := "org.gnome.desktop.interface color-scheme"
	r.Settings[key] = "'prefer-dark'"
	p, err := a.BuildPlan(map[string]catalog.Selection{}, []string{"color-scheme"})
	if err != nil {
		t.Fatal(err)
	}
	r.Settings[key] = "'prefer-light'"
	if _, err := a.Execute(p.Selection, p.Reverts, nil, p.ReviewToken); err == nil {
		t.Fatal("review token must reject changed reset value")
	}
	if a.eng.Running() || r.Settings[key] != "'prefer-light'" {
		t.Fatal("stale plan changed the system")
	}
}

func TestDiagnosticExportsUseReviewedSnapshotAndDoNotMutate(t *testing.T) {
	a := rcApp(t)
	rep := a.Diagnostics()
	if rep.SnapshotID == "" || rep.CollectedAt == "" {
		t.Fatal("snapshot metadata missing")
	}
	for _, format := range []string{"pdf", "csv", "xlsx"} {
		b, err := a.diagnosticExport(rep.SnapshotID, format, []string{"facts"})
		if err != nil || len(b) == 0 {
			t.Fatalf("%s: %v", format, err)
		}
	}
	a.Diagnostics()
	if _, err := a.diagnosticExport(rep.SnapshotID, "csv", []string{"facts"}); err == nil {
		t.Fatal("stale snapshot must not export different unseen data")
	}
	if a.eng.Running() || len(a.st.Runs()) != 0 || len(a.st.Changes()) != 0 {
		t.Fatal("export executed or journaled a mutation")
	}
}
