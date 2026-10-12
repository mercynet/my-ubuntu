package inventory

import (
	"testing"

	"meuubuntu/internal/catalog"
)

func TestClassify(t *testing.T) {
	cat, _ := catalog.Load()
	in := `git
linux-image-6.8.0-45-generic
grub-efi-amd64-signed
libicu74
libssl-dev
some-weird-tool
php8.1-cli
Name     Version  Rev   Tracking       Publisher   Notes
firefox  131.0    4793  latest/stable  mozilla✓    -
core22   20240904 1621  latest/stable  canonical✓  base
Spotify	com.spotify.Client	1.2	stable	flathub	system
org.freefilesync.FreeFileSync`
	got := map[string]Entry{}
	for _, e := range Parse(in, cat) {
		got[e.Name] = e
	}
	want := map[string]string{"git": "app", "linux-image-6.8.0-45-generic": "system", "grub-efi-amd64-signed": "system",
		"libicu74": "legacy", "libssl-dev": "dependency", "some-weird-tool": "unknown", "php8.1-cli": "legacy",
		"firefox": "app", "core22": "system", "com.spotify.Client": "unknown", "org.freefilesync.FreeFileSync": "app"}
	for n, c := range want {
		if got[n].Class != c {
			t.Errorf("%s: got %q want %q (%+v)", n, got[n].Class, c, got[n])
		}
	}
	if got["firefox"].Source != "snap" || got["org.freefilesync.FreeFileSync"].Source != "flatpak" {
		t.Errorf("sources: %+v %+v", got["firefox"], got["org.freefilesync.FreeFileSync"])
	}
}
