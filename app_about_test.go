package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/store"
)

func TestAboutReadsExecutableAndAttributionsWithoutChangingUserState(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "c"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(tmp, "s"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmp, "d"))
	a, err := NewApp(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.SavePrefs(store.Prefs{Selections: map[string]catalog.Selection{"git": {}}, CatalogProfile: "builtin:terminal-produtivo"}); err != nil {
		t.Fatal(err)
	}
	before := a.Init()
	info, err := a.About()
	if err != nil {
		t.Fatal(err)
	}
	if info.Version != before.Version || info.Version != version || info.GoVersion == "" || info.Author == "" {
		t.Fatalf("incorrect executable metadata: %+v", info)
	}
	if len(info.Attributions) != 5 {
		t.Fatalf("missing attributions: %+v", info.Attributions)
	}
	for _, attribution := range info.Attributions {
		if attribution.Text == "" || !strings.HasPrefix(attribution.URL, "https://") || attribution.License == "" {
			t.Fatalf("incomplete attribution: %+v", attribution)
		}
	}
	after := a.Init()
	changes, err := a.st.ReadChanges()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before.Prefs, after.Prefs) || a.eng.Running() || len(a.st.Runs()) != 0 || len(changes) != 0 {
		t.Fatal("About must be read-only")
	}
}
