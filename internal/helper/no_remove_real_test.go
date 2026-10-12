//go:build realsystem

package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

// Map only the fixture's package name to its local archive. The helper still
// constructs the complete APT command, including its protection flags.
type localFixtureAPT struct {
	sys.Real
	name, archive string
}

func (r localFixtureAPT) Stream(ctx context.Context, log func(string), name string, args ...string) error {
	if name == "apt-get" {
		args = append([]string(nil), args...)
		for i, arg := range args {
			if arg == r.name {
				args[i] = r.archive
			}
		}
	}
	return r.Real.Stream(ctx, log, name, args...)
}

func TestRealAPTRefusesImplicitPackageRemoval(t *testing.T) {
	if os.Getenv("MEU_UBUNTU_DISPOSABLE") != "1" || os.Geteuid() != 0 {
		t.Skip("requires root in a disposable container with MEU_UBUNTU_DISPOSABLE=1")
	}
	ctx := context.Background()
	r := sys.Real{}
	kept, candidate := "meu-ubuntu-test-kept", "meu-ubuntu-test-conflict"
	for _, pkg := range []string{kept, candidate} {
		if out, _ := r.Output(ctx, "dpkg-query", "-W", "-f=${Status}", pkg); strings.TrimSpace(out) == "install ok installed" {
			t.Skip("fixture package already present; refusing to alter it")
		}
	}
	tmp := t.TempDir()
	build := func(name, conflict string) string {
		t.Helper()
		dir := filepath.Join(tmp, name)
		if err := os.MkdirAll(filepath.Join(dir, "DEBIAN"), 0o755); err != nil {
			t.Fatal(err)
		}
		control := fmt.Sprintf("Package: %s\nVersion: 1.0\nArchitecture: all\nMaintainer: Test <test@example.invalid>\nDescription: Disposable preservation test fixture\n", name)
		if conflict != "" {
			control += "Conflicts: " + conflict + "\n"
		}
		if err := os.WriteFile(filepath.Join(dir, "DEBIAN", "control"), []byte(control), 0o644); err != nil {
			t.Fatal(err)
		}
		archive := filepath.Join(tmp, name+".deb")
		if _, err := r.Output(ctx, "dpkg-deb", "--build", dir, archive); err != nil {
			t.Fatal(err)
		}
		return archive
	}
	keepArchive, conflictArchive := build(kept, ""), build(candidate, kept)
	if _, err := r.Output(ctx, "dpkg", "--install", keepArchive); err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = r.Output(ctx, "dpkg", "--remove", candidate, kept) }()
	// Prove that this transaction would remove the existing fixture without
	// the protection. This baseline is a dry run and makes no package changes.
	baseline, err := r.Output(ctx, "apt-get", "--simulate", "install", conflictArchive)
	if err != nil || !strings.Contains(baseline, "Remv "+kept) {
		t.Fatalf("fixture must require removal: %v\n%s", err, baseline)
	}
	data, _ := json.Marshal(catalog.Catalog{Items: []catalog.Item{{ID: "conflict-fixture", Name: "Conflict fixture", Method: catalog.MethodAPT, Packages: []string{candidate}}}})
	cat, err := catalog.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	e := &Env{R: localFixtureAPT{Real: r, name: candidate, archive: conflictArchive}, Cat: cat, Distribution: "ubuntu", Codename: "resolute"}
	var log []string
	err = e.Execute(ctx, Request{Op: "apt-install", Item: "conflict-fixture"}, func(line string) { log = append(log, line) })
	if err == nil || !strings.Contains(strings.ToLower(strings.Join(log, "\n")), "remove is disabled") {
		t.Fatalf("APT must refuse removal: %v\n%s", err, strings.Join(log, "\n"))
	}
	if out, _ := r.Output(ctx, "dpkg-query", "-W", "-f=${Status}", kept); strings.TrimSpace(out) != "install ok installed" {
		t.Fatalf("existing package was removed: %s", out)
	}
	if out, _ := r.Output(ctx, "dpkg-query", "-W", "-f=${Status}", candidate); strings.TrimSpace(out) == "install ok installed" {
		t.Fatal("conflicting package must not be installed")
	}
	t.Log("APT refused the conflicting installation and preserved the installed package")
}
