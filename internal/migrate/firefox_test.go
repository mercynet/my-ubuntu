package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrationUsesUniqueProfileIndex(t *testing.T) {
	if Running() {
		t.Skip("Firefox is running")
	}
	home := t.TempDir()
	src := filepath.Join(home, "snap/firefox/common/.mozilla/firefox/default")
	os.MkdirAll(src, 0700)
	os.WriteFile(filepath.Join(src, "bookmarks"), []byte("test data"), 0600)
	root := filepath.Join(home, ".var/app/org.mozilla.firefox/config/mozilla/firefox")
	os.MkdirAll(root, 0700)
	ini := filepath.Join(root, "profiles.ini")
	os.WriteFile(ini, []byte("[Profile0]\nName=First\nPath=first\n\n[Profile2]\nName=Second\nPath=second\n"), 0600)
	if _, err := MigrateToFlatpak(home, src); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(ini)
	if strings.Count(string(b), "[Profile2]") != 1 || !strings.Contains(string(b), "[Profile3]") {
		t.Fatal("profile identifiers must remain unique")
	}
}

func TestBackupRejectsProfileSymlinkEscape(t *testing.T) {
	if Running() {
		t.Skip("Firefox is running")
	}
	home := t.TempDir()
	root := filepath.Join(home, "snap/firefox/common/.mozilla/firefox")
	os.MkdirAll(root, 0700)
	external := filepath.Join(home, "outside-firefox")
	src := filepath.Join(external, "dummy-private-directory")
	os.MkdirAll(src, 0700)
	os.WriteFile(filepath.Join(src, "dummy-private-data"), []byte("test fixture only"), 0600)
	os.Symlink(external, filepath.Join(root, "alias"))
	path := filepath.Join(root, "alias", "dummy-private-directory")
	if allowedProfile(home, path) {
		t.Fatal("symlink escape must be rejected")
	}
	if _, err := Backup(home, path); err == nil {
		t.Fatal("backup outside Firefox roots must fail")
	}
}
