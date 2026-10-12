package main

import (
	"path/filepath"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/store"
)

func TestManualChecklistIsPersistentButNotInstallationEvidence(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "c"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(tmp, "s"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmp, "d"))
	a, err := NewApp(true)
	if err != nil {
		t.Fatal(err)
	}
	guide := a.cat.Item("oh-my-zsh").Manual.Steps
	prefs := store.Prefs{Selections: map[string]catalog.Selection{"oh-my-zsh": {}}, ManualProgress: map[string][]string{
		"oh-my-zsh": append(append([]string{}, guide...), guide[0], "uma instrução antiga"),
		"unknown":   {"comando inventado"},
		"git":       {"uma instrução inventada"},
	}}
	if err := a.SavePrefs(prefs); err != nil {
		t.Fatal(err)
	}
	after, err := NewApp(true)
	if err != nil {
		t.Fatal(err)
	}
	saved := after.Init().Prefs
	if len(saved.ManualProgress) != 1 || len(saved.ManualProgress["oh-my-zsh"]) != len(guide) {
		t.Fatalf("only current, distinct guide steps should survive: %+v", saved.ManualProgress)
	}
	p, err := after.BuildPlan(saved.Selections, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Manual) != 1 || p.Manual[0].Item != "oh-my-zsh" || after.eng.Running() || len(after.st.Runs()) != 0 {
		t.Fatal("marking a checklist must not confirm, execute or skip a manual installation")
	}
	// A changed guide is new work, even when its position in the checklist matches.
	after.cat.Item("oh-my-zsh").Manual.Steps[0] = "Uma nova instrução revista."
	if got := after.Init().Prefs.ManualProgress["oh-my-zsh"]; len(got) != len(guide)-1 {
		t.Fatalf("a changed instruction inherited an old completion mark: %+v", got)
	}
}
