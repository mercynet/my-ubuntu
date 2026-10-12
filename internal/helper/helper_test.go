package helper

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

func env(t *testing.T) (*Env, *sys.Sim) {
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	sim := sys.NewSim()
	sim.Delay = 0
	return &Env{R: sim, Cat: cat, Root: t.TempDir(), Codename: "resolute", UID: "1000", Fetch: func(string) ([]byte, error) { return nil, nil }}, sim
}

func TestRejectsUnknownOpsAndItems(t *testing.T) {
	e, sim := env(t)
	for _, r := range []Request{
		{Op: "shell", Item: "htop"},
		{Op: "apt-install", Item: "does-not-exist"},
		{Op: "apt-install", Item: "color-scheme"}, // not an apt item
		{Op: "admin", Item: "htop"},
		{Op: "apt-install", Item: "php", Options: map[string][]string{"versions": {"8.5 && reboot"}}},
	} {
		if err := e.Execute(context.Background(), r, func(string) {}); err == nil {
			t.Errorf("expected rejection for %+v", r)
		}
	}
	for _, c := range sim.Calls {
		if strings.HasPrefix(c, "apt-get") {
			t.Errorf("no apt call expected, got %s", c)
		}
	}
}

func TestAptInstallUsesCatalogPackagesAndLockTimeout(t *testing.T) {
	e, sim := env(t)
	if err := e.Execute(context.Background(), Request{Op: "apt-install", Item: "smb"}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	last := sim.Calls[len(sim.Calls)-1]
	if !strings.Contains(last, "cifs-utils") || !strings.Contains(last, "smbclient") || !strings.Contains(last, "DPkg::Lock::Timeout") || !strings.Contains(last, "--no-remove") {
		t.Fatalf("unexpected call %s", last)
	}
}

func TestRepoRefusesWrongFingerprint(t *testing.T) {
	e, _ := env(t)
	e.R = fakeGPG{fpr: "0000000000000000000000000000000000000000"}
	e.Fetch = func(string) ([]byte, error) { return []byte("-----BEGIN PGP PUBLIC KEY BLOCK-----"), nil }
	err := e.Execute(context.Background(), Request{Op: "repo-add", Item: "docker-engine"}, func(string) {})
	if err == nil || !strings.Contains(err.Error(), "rejeitada") {
		t.Fatalf("expected fingerprint rejection, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(e.Root, "etc/apt/sources.list.d/meu-ubuntu-docker.sources")); err == nil {
		t.Fatal("sources file must not be written")
	}
}

func TestDropinRevertRefusesForeignFile(t *testing.T) {
	e, _ := env(t)
	p := filepath.Join(e.Root, "etc/systemd/journald.conf.d/60-meu-ubuntu.conf")
	os.MkdirAll(filepath.Dir(p), 0o755)
	os.WriteFile(p, []byte("[Journal]\nSystemMaxUse=9G\n"), 0o644)
	if err := e.Execute(context.Background(), Request{Op: "revert", Item: "journald-limit"}, func(string) {}); err == nil {
		t.Fatal("must refuse to remove a file it did not write")
	}
	os.Remove(p)
	if err := e.Execute(context.Background(), Request{Op: "admin", Item: "journald-limit", Options: map[string][]string{"value": {"1G"}}}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if !strings.Contains(string(b), "SystemMaxUse=1G") {
		t.Fatalf("content %q", b)
	}
	// Same header but edited content: refused (compared byte-for-byte, not by header).
	os.WriteFile(p, []byte(strings.Replace(string(b), "1G", "5G", 1)), 0o644)
	if err := e.Execute(context.Background(), Request{Op: "revert", Item: "journald-limit", Options: map[string][]string{"value": {"1G"}}}, func(string) {}); err == nil {
		t.Fatal("edited file must not be removed")
	}
	os.WriteFile(p, b, 0o644)
	if err := e.Execute(context.Background(), Request{Op: "revert", Item: "journald-limit", Options: map[string][]string{"value": {"1G"}}}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); err == nil {
		t.Fatal("not removed")
	}
}

type fakeGPG struct{ fpr string }

func (f fakeGPG) Output(ctx context.Context, name string, args ...string) (string, error) {
	return "pub:-:4096:1:X:1:::-:::scESC::::::23::0:\nfpr:::::::::" + f.fpr + ":\n", nil
}
func (f fakeGPG) Stream(ctx context.Context, on func(string), name string, args ...string) error {
	return nil
}
func (fakeGPG) LookPath(string) bool { return true }

func TestDedupeRemovesOnlyOurDuplicate(t *testing.T) {
	e, _ := env(t)
	d := filepath.Join(e.Root, "etc/apt/sources.list.d")
	os.MkdirAll(d, 0o755)
	os.MkdirAll(filepath.Join(e.Root, "etc/apt/keyrings"), 0o755)
	ours := filepath.Join(d, "meu-ubuntu-warp.sources")
	os.WriteFile(ours, []byte("# Gerido por Meu Ubuntu — Warp\nTypes: deb\nURIs: https://releases.warp.dev/linux/deb\nSuites: stable\nComponents: main\nArchitectures: amd64\n"), 0o644)
	os.WriteFile(filepath.Join(e.Root, "etc/apt/keyrings/meu-ubuntu-warp.gpg"), []byte("k"), 0o644)
	vendor := filepath.Join(d, "warpdotdev.list")
	os.WriteFile(vendor, []byte("deb [arch=amd64 signed-by=/etc/apt/trusted.gpg.d/warpdotdev.gpg] https://releases.warp.dev/linux/deb/ stable main\n"), 0o644)
	keep := filepath.Join(d, "meu-ubuntu-docker.sources")
	os.WriteFile(keep, []byte("# Gerido por Meu Ubuntu — Docker\nURIs: https://download.docker.com/linux/ubuntu\n"), 0o644)
	if err := e.Execute(context.Background(), Request{Op: "sources-check"}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ours); err == nil {
		t.Fatal("our duplicate must be removed")
	}
	if _, err := os.Stat(vendor); err != nil {
		t.Fatal("vendor source must be kept")
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatal("non-duplicated source must be kept")
	}
	// repo-add must not recreate it while the vendor file exists.
	if err := e.Execute(context.Background(), Request{Op: "repo-add", Item: "warp"}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ours); err == nil {
		t.Fatal("repo-add must not duplicate an existing source")
	}
}
