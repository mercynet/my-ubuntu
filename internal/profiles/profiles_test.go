package profiles

import (
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/store"
)

func TestBuiltinProfilesAreValid(t *testing.T) {
	cat, _ := catalog.Load()
	ps := Builtin()
	if len(ps) < 8 {
		t.Fatalf("expected built-in profiles, got %d", len(ps))
	}
	usage := 0
	for _, p := range ps {
		b, _ := store.MarshalExport(p)
		q, w, err := store.ParseProfile(cat, b)
		if err != nil || len(w) != 0 {
			t.Errorf("%s: %v %v", p.ID, err, w)
			continue
		}
		if q.Kind == "usage" {
			usage++
		}
	}
	if usage != 4 {
		t.Errorf("expected 4 usage profiles, got %d", usage)
	}
}
