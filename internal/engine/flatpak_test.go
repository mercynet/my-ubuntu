package engine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/sys"
)

type flatpakProbe struct {
	sys.Runner
	missing      bool
	versionError bool
	remote       string
	ref          string
	refError     bool
	installs     int
}

func (r *flatpakProbe) LookPath(name string) bool {
	if name == "flatpak" && r.missing {
		return false
	}
	return r.Runner.LookPath(name)
}
func (r *flatpakProbe) Output(ctx context.Context, name string, args ...string) (string, error) {
	if name == "flatpak" && len(args) > 0 {
		switch args[0] {
		case "--version":
			if r.versionError {
				return "", errors.New("broken")
			}
		case "remotes":
			return r.remote, nil
		case "remote-info":
			if r.refError {
				return "", errors.New("offline")
			}
			return r.ref, nil
		}
	}
	return r.Runner.Output(ctx, name, args...)
}
func (r *flatpakProbe) Stream(ctx context.Context, log func(string), name string, args ...string) error {
	if name == "flatpak" && len(args) > 0 && args[0] == "install" {
		r.installs++
	}
	return r.Runner.Stream(ctx, log, name, args...)
}

func TestFlatpakValidationRefusesInstallationBeforeMutation(t *testing.T) {
	cases := []struct {
		name      string
		configure func(*flatpakProbe)
	}{
		{"missing", func(r *flatpakProbe) { r.missing = true }},
		{"broken", func(r *flatpakProbe) { r.versionError = true }},
		{"no-user-remote", func(r *flatpakProbe) { r.remote = "" }},
		{"foreign-remote", func(r *flatpakProbe) { r.remote = "flathub\thttps://example.org/repo\t\n" }},
		{"unsigned", func(r *flatpakProbe) { r.remote = "flathub\thttps://dl.flathub.org/repo/\tno-gpg-verify\n" }},
		{"disabled", func(r *flatpakProbe) { r.remote = "flathub\thttps://dl.flathub.org/repo/\tdisabled\n" }},
		{"offline-or-unpublished", func(r *flatpakProbe) { r.refError = true }},
		{"wrong-application", func(r *flatpakProbe) { r.ref = "app/another.Application/x86_64/stable" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := setup(t)
			f.sim.Packages["flatpak"] = true
			r := &flatpakProbe{Runner: f.sim, remote: "flathub\thttps://dl.flathub.org/repo/\t\n", ref: "app/org.videolan.VLC/x86_64/stable"}
			tc.configure(r)
			f.eng.R = r
			err := f.eng.userOp(context.Background(), "test", plan.Step{Item: "vlc", Op: "flatpak-install"}, func(string) {})
			if err == nil || r.installs != 0 {
				t.Fatalf("validation bypassed: error=%v installs=%d", err, r.installs)
			}
		})
	}
}

func TestFlatpakPreflightRunsBeforeSuccessfulInstallation(t *testing.T) {
	f := setup(t)
	f.sim.Packages["flatpak"] = true
	f.sim.FlatpakRemotes["flathub"] = flathubRepo
	if err := f.eng.userOp(context.Background(), "test", plan.Step{Item: "vlc", Op: "flatpak-install"}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	check, install := -1, -1
	for i, call := range f.sim.Calls {
		if strings.HasPrefix(call, "flatpak remote-info") {
			check = i
		}
		if strings.HasPrefix(call, "flatpak install") {
			install = i
		}
	}
	if check < 0 || install <= check || !f.sim.Flatpaks["org.videolan.VLC"] {
		t.Fatalf("installation preceded validation: %v", f.sim.Calls)
	}
}

func TestFailedFlathubSetupSkipsDependentApplications(t *testing.T) {
	f := setup(t)
	f.sim.Packages["flatpak"] = true
	f.sim.FlatpakRemotes["flathub"] = "https://example.org/repo"
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"vlc": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	awaitEngine(t, f.eng)
	found := false
	for _, step := range f.st.Runs()[0].Steps {
		if step.Item == "vlc" {
			found = true
			if step.Status != StSkipped {
				t.Fatalf("dependent app ran: %+v", step)
			}
		}
	}
	if !found || f.sim.Flatpaks["org.videolan.VLC"] {
		t.Fatal("unsafe source installed app")
	}
}
