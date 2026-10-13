package detect

import (
	"context"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

func TestFlatpakMigrationPreservesNativeAndPreviousAppIDs(t *testing.T) {
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	r := sys.NewSim()
	r.Delay = 0
	r.Packages["vlc"] = true
	r.Packages["flatpak"] = true
	r.Flatpaks["org.mozilla.Thunderbird"] = true
	d := &Detector{R: r, Cat: cat, Sys: ReadSystem(r, true), Home: t.TempDir()}
	d.Refresh(context.Background())
	for _, id := range []string{"vlc", "thunderbird"} {
		if st := d.State(context.Background(), cat.Item(id), catalog.Selection{}); st.Status != Provided {
			t.Fatalf("existing %s lost: %+v", id, st)
		}
	}
	r.Flatpaks["org.mozilla.thunderbird"] = true
	d.Refresh(context.Background())
	if st := d.State(context.Background(), cat.Item("thunderbird"), catalog.Selection{}); st.Status != Installed {
		t.Fatalf("current app ID not detected: %+v", st)
	}
}

func TestFlatpakMissingPrerequisiteIsActionableAndEOLIsUnavailable(t *testing.T) {
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	r := sys.NewSim()
	r.Delay = 0
	d := &Detector{R: r, Cat: cat, Sys: ReadSystem(r, true), Home: t.TempDir()}
	d.Refresh(context.Background())
	if st := d.State(context.Background(), cat.Item("vlc"), catalog.Selection{}); st.Status != Available || st.Detail == "" {
		t.Fatalf("missing prerequisite hidden: %+v", st)
	}
	if st := d.State(context.Background(), cat.Item("metadata-cleaner"), catalog.Selection{}); st.Status != Incompatible {
		t.Fatalf("EOL install allowed: %+v", st)
	}
	r.Packages["flatpak"] = true
	r.Flatpaks["fr.romainvigier.MetadataCleaner"] = true
	d.Refresh(context.Background())
	if st := d.State(context.Background(), cat.Item("metadata-cleaner"), catalog.Selection{}); st.Status != Installed {
		t.Fatalf("existing EOL app lost: %+v", st)
	}
}
