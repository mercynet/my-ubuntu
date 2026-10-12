package helper

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDisabledSourcePreservesActiveSource(t *testing.T) {
	e, _ := env(t)
	d := filepath.Join(e.Root, "etc/apt/sources.list.d")
	os.MkdirAll(d, 0755)
	os.MkdirAll(filepath.Join(e.Root, "etc/apt/keyrings"), 0755)
	ours := filepath.Join(d, "meu-ubuntu-warp.sources")
	os.WriteFile(ours, []byte("# Gerido por Meu Ubuntu — Warp\nTypes: deb\nURIs: https://releases.warp.dev/linux/deb\nSuites: stable\nComponents: main\n"), 0644)
	key := filepath.Join(e.Root, "etc/apt/keyrings/meu-ubuntu-warp.gpg")
	os.WriteFile(key, []byte("dummy key"), 0644)
	os.WriteFile(filepath.Join(d, "disabled.sources"), []byte("Enabled: no\nTypes: deb\nURIs: https://releases.warp.dev/linux/deb\nSuites: stable\nComponents: main\n"), 0644)
	if err := e.Execute(context.Background(), Request{Op: "sources-check"}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ours); err != nil {
		t.Fatal("active source must be preserved", err)
	}
}

func TestDedupePreservesSharedKeyring(t *testing.T) {
	e, _ := env(t)
	d := filepath.Join(e.Root, "etc/apt/sources.list.d")
	os.MkdirAll(d, 0755)
	os.MkdirAll(filepath.Join(e.Root, "etc/apt/keyrings"), 0755)
	os.WriteFile(filepath.Join(d, "meu-ubuntu-warp.sources"), []byte("# Gerido por Meu Ubuntu\nTypes: deb\nURIs: https://releases.warp.dev/linux/deb\nSuites: stable\nComponents: main\n"), 0644)
	key := filepath.Join(e.Root, "etc/apt/keyrings/meu-ubuntu-warp.gpg")
	os.WriteFile(key, []byte("dummy key"), 0644)
	vendor := filepath.Join(d, "vendor.sources")
	os.WriteFile(vendor, []byte("Types: deb\nURIs: https://releases.warp.dev/linux/deb\nSuites: stable\nComponents: main\nSigned-By: /etc/apt/keyrings/meu-ubuntu-warp.gpg\n"), 0644)
	if err := e.Execute(context.Background(), Request{Op: "sources-check"}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(key); err != nil {
		t.Fatal("shared keyring must be preserved", err)
	}
	if _, err := os.Stat(vendor); err != nil {
		t.Fatal(err)
	}
}

func TestDropinRevertRestoresOriginalFile(t *testing.T) {
	e, _ := env(t)
	p := filepath.Join(e.Root, "etc/systemd/journald.conf.d/60-meu-ubuntu.conf")
	os.MkdirAll(filepath.Dir(p), 0755)
	original := []byte("[Journal]\nSystemMaxUse=2G\n")
	os.WriteFile(p, original, 0644)
	opts := map[string][]string{"value": {"1G"}}
	for _, op := range []string{"admin", "revert"} {
		if err := e.Execute(context.Background(), Request{Op: op, Item: "journald-limit", Options: opts}, func(string) {}); err != nil {
			t.Fatal(err)
		}
	}
	if b, err := os.ReadFile(p); err != nil || string(b) != string(original) {
		t.Fatal("original file must be restored", err)
	}
}

type reviewRedirectTransport struct{ schemes []string }

func (r *reviewRedirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.schemes = append(r.schemes, req.URL.Scheme)
	if req.URL.Scheme == "https" {
		return &http.Response{StatusCode: 302, Header: http.Header{"Location": []string{"http://fixture.invalid/script"}}, Body: io.NopCloser(strings.NewReader("")), Request: req}, nil
	}
	return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("test fixture")), Request: req}, nil
}
func TestHTTPSRedirectDowngradeRejected(t *testing.T) {
	old := http.DefaultTransport
	defer func() { http.DefaultTransport = old }()
	tr := &reviewRedirectTransport{}
	http.DefaultTransport = tr
	// No network request is made: transport intercepts both URLs.
	b, err := HTTPFetch("https://fixture.invalid/script")
	if err == nil || len(b) != 0 {
		t.Fatal("downgrade must be rejected")
	}
	if strings.Join(tr.schemes, ",") != "https" {
		t.Fatal("plaintext request must never be sent")
	}
}

func TestInotifyRevertRestoresLiveLimit(t *testing.T) {
	e, sim := env(t)
	p := e.path("/proc/sys/fs/inotify/max_user_watches")
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("65536\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, op := range []string{"admin", "revert"} {
		if err := e.Execute(context.Background(), Request{Op: op, Item: "inotify-watches", Options: map[string][]string{"value": {"524288"}}}, func(string) {}); err != nil {
			t.Fatal(err)
		}
	}
	if !strings.Contains(strings.Join(sim.Calls, "\n"), "sysctl -w fs.inotify.max_user_watches=65536") {
		t.Fatal("live inotify value must be restored")
	}
	for _, c := range sim.Calls {
		if strings.Contains(c, "sysctl --system") {
			t.Fatal("unrelated sysctl settings must not be reloaded")
		}
	}
}

func TestDropinSnapshotFailurePreventsWriting(t *testing.T) {
	e, _ := env(t)
	stateDir := e.path("/var/lib/meu-ubuntu")
	if err := os.MkdirAll(stateDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(stateDir, "journald.json"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := e.Execute(context.Background(), Request{Op: "admin", Item: "journald-limit"}, func(string) {}); err == nil {
		t.Fatal("unwritable snapshot must fail")
	}
	if _, err := os.Stat(e.path("/etc/systemd/journald.conf.d/60-meu-ubuntu.conf")); !os.IsNotExist(err) {
		t.Fatal("config must not be changed without a snapshot")
	}
}

func TestUbuntuHelperRecipesRefusedOnFedora(t *testing.T) {
	e, sim := env(t)
	e.Distribution = "fedora"
	for _, op := range []string{"apt-install", "repo-add", "sources-check"} {
		if err := e.Execute(context.Background(), Request{Op: op, Item: "docker-engine"}, func(string) {}); err == nil {
			t.Fatalf("%s must not run on Fedora", op)
		}
	}
	if len(sim.Calls) > 0 {
		t.Fatal("no administrative command may run")
	}
}

func TestConsecutiveDropinRevertsRestoreEachPreviousValue(t *testing.T) {
	e, _ := env(t)
	for _, value := range []string{"1G", "2G"} {
		if err := e.Execute(context.Background(), Request{Op: "admin", Item: "journald-limit", Options: map[string][]string{"value": {value}}}, func(string) {}); err != nil {
			t.Fatal(err)
		}
	}
	if err := e.Execute(context.Background(), Request{Op: "revert", Item: "journald-limit", Options: map[string][]string{"value": {"1G"}}}, func(string) {}); err == nil {
		t.Fatal("older reversal must not overwrite a newer change")
	}
	if err := e.Execute(context.Background(), Request{Op: "revert", Item: "journald-limit", Options: map[string][]string{"value": {"2G"}}}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(e.path("/etc/systemd/journald.conf.d/60-meu-ubuntu.conf"))
	if err != nil || !strings.Contains(string(b), "SystemMaxUse=1G") {
		t.Fatal("previous 1G setting must be restored", err)
	}
	if err := e.Execute(context.Background(), Request{Op: "revert", Item: "journald-limit", Options: map[string][]string{"value": {"1G"}}}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(e.path("/etc/systemd/journald.conf.d/60-meu-ubuntu.conf")); !os.IsNotExist(err) {
		t.Fatal("original absence must be restored")
	}
}
