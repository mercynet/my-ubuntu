package selection

import (
	"testing"

	"meuubuntu/internal/catalog"
)

func TestMergeDedupConflicts(t *testing.T) {
	cat, _ := catalog.Load()
	cur := map[string]catalog.Selection{
		"htop":          {},
		"dock-position": {Options: map[string][]string{"pos": {"'LEFT'"}}},
		"power-saver":   {},
		"docker-engine": {},
	}
	orig := map[string][]string{"htop": {Manual}, "dock-position": {Manual}, "power-saver": {"perfil:Portátil"}, "docker-engine": {Manual}}
	in := map[string]catalog.Selection{
		"htop":              {},                                                  // duplicate -> deduplicated
		"dock-position":     {Options: map[string][]string{"pos": {"'BOTTOM'"}}}, // option conflict
		"power-performance": {},                                                  // group conflict
		"docker-ubuntu":     {},                                                  // hard conflict
		"ncdu":              {},                                                  // added
		"unknown-item":      {},
	}
	r := Merge(cat, cur, orig, in, "perfil:Infra")
	if len(r.Added) != 1 || r.Added[0] != "ncdu" {
		t.Fatalf("added %v", r.Added)
	}
	if got := r.Origins["htop"]; len(got) != 2 || got[0] != Manual || got[1] != "perfil:Infra" {
		t.Fatalf("htop origins %v", got)
	}
	kinds := map[string]string{}
	for _, c := range r.Conflicts {
		kinds[c.IncomingItem] = c.Kind
	}
	if kinds["dock-position"] != "options" || kinds["power-performance"] != "group" || kinds["docker-ubuntu"] != "exclusive" || len(r.Conflicts) != 3 {
		t.Fatalf("conflicts %+v", r.Conflicts)
	}
	// Manual choices are preserved until the user resolves the conflict.
	if r.Selections["dock-position"].Options["pos"][0] != "'LEFT'" || r.Selections["power-performance"].Options != nil {
		t.Fatalf("existing choices must not be overridden: %+v", r.Selections)
	}
	if _, ok := r.Selections["power-performance"]; ok {
		t.Fatal("conflicting incoming item must not be added")
	}
}

func TestMergeTwoProfilesNotDecidedByOrder(t *testing.T) {
	cat, _ := catalog.Load()
	a := map[string]catalog.Selection{"power-balanced": {}}
	b := map[string]catalog.Selection{"power-performance": {}}
	r := Merge(cat, nil, nil, a, "perfil:A")
	r = Merge(cat, r.Selections, r.Origins, b, "perfil:B")
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected a conflict between A and B, got %+v", r.Conflicts)
	}
}
