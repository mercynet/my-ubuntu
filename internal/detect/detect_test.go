package detect

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"meuubuntu/internal/adminconfig"
	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

func TestAdministrativeDetectionComparesChosenValue(t *testing.T) {
	c, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	r := sys.NewSim()
	r.Delay = 0
	d := &Detector{R: r, Cat: c, Sys: ReadSystem(r, true), Root: t.TempDir()}
	d.Refresh(context.Background())
	it := c.Item("journald-limit")
	sel, err := c.ValidateSelection(it.ID, catalog.Selection{Options: map[string][]string{"value": {"1G"}}})
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(d.Root, adminconfig.Dropins["journald"].Path)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("[Journal]\nSystemMaxUse=2G\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if st := d.State(context.Background(), it, sel); st.Status != Available {
		t.Fatalf("existing wrong value is not applied: %+v", st)
	}
	if err := os.WriteFile(p, []byte(fmt.Sprintf(adminconfig.Dropins["journald"].Template, "1G")), 0644); err != nil {
		t.Fatal(err)
	}
	if st := d.State(context.Background(), it, sel); st.Status != NeedsConfig {
		t.Fatalf("file alone must not prove a successful reload: %+v", st)
	}
	b, _ := os.ReadFile(p)
	if err := os.WriteFile(p+".applied", []byte(fmt.Sprintf("%x", sha256.Sum256(b))), 0644); err != nil {
		t.Fatal(err)
	}
	if st := d.State(context.Background(), it, sel); st.Status != Installed {
		t.Fatalf("expected verified setting: %+v", st)
	}
}

func TestUbuntuRecipesAreNotOfferedOnFedora(t *testing.T) {
	c, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	r := sys.NewSim()
	d := &Detector{R: r, Cat: c, Sys: System{ID: "fedora", Release: "44", Arch: "amd64"}}
	if st := d.State(context.Background(), c.Item("docker-engine"), catalog.Selection{}); st.Status != Incompatible {
		t.Fatalf("Ubuntu repository must be unavailable: %+v", st)
	}
}
