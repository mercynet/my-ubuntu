//go:build realsystem

// Real (non-simulated) user-level adapters. Run ONLY in a disposable
// container/VM: `go test -tags realsystem ./internal/engine -run Real -v`.
package engine

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/store"
	"meuubuntu/internal/sys"
)

func TestRealUserAdapters(t *testing.T) {
	if os.Getenv("MEU_UBUNTU_DISPOSABLE") != "1" {
		t.Skip("set MEU_UBUNTU_DISPOSABLE=1 (disposable environment only)")
	}
	tmp := t.TempDir()
	t.Setenv("XDG_STATE_HOME", filepath.Join(tmp, "s"))
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "c"))
	cat, _ := catalog.Load()
	r := sys.Real{}
	st, _ := store.New("real")
	home, _ := os.UserHomeDir()
	det := &detect.Detector{R: r, Cat: cat, Sys: detect.ReadSystem(r, false), Home: home}
	eng := &Engine{Cat: cat, Det: det, R: r, Store: st, Home: home, Fetch: helper.HTTPFetch, Emit: func(ev string, d any) {
		if l, ok := d.(LogEvent); ok {
			t.Log(l.StepID, l.Line)
		}
	}, NewAdmin: func() helper.Admin { t.Fatal("no admin steps expected"); return nil }}
	det.Refresh(context.Background())
	sel := map[string]catalog.Selection{"color-scheme": {}, "dock-click": {}, "git-defaults": {}, "tpl-sh": {}, "freefilesync": {}}
	delete(sel, "flathub") // flatpak preinstalled in the container
	p, err := plan.Build(context.Background(), cat, det, sel, nil)
	if err != nil {
		t.Fatal(err)
	}
	var steps []plan.Step
	for _, s := range p.Steps {
		if s.Kind == "user" {
			steps = append(steps, s)
		}
	}
	p.Steps = steps
	if _, err := eng.Start(p, nil); err != nil {
		t.Fatal(err)
	}
	for eng.Running() {
		time.Sleep(200 * time.Millisecond)
	}
	run := st.Runs()[0]
	for _, s := range run.Steps {
		t.Logf("%-40s %s %s", s.Title, s.Status, s.Error)
		if s.Status != StDone {
			t.Errorf("step %s: %s %s", s.Title, s.Status, s.Error)
		}
	}
	// Idempotent rerun: nothing left to do for the user-level items.
	det.Refresh(context.Background())
	p2, _ := plan.Build(context.Background(), cat, det, sel, nil)
	for _, s := range p2.Steps {
		if s.Kind == "user" {
			t.Errorf("rerun not idempotent: %s", s.Title)
		}
	}
	// Revert every journaled change and check real values came back.
	for _, c := range st.Changes() {
		if err := eng.Revert(context.Background(), c, false); err != nil {
			t.Errorf("revert %s %s: %v", c.Title, c.Key, err)
		}
	}
	out, _ := r.Output(context.Background(), "gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	if strings.TrimSpace(out) != "'default'" {
		t.Errorf("color-scheme not reverted: %s", out)
	}
	if _, err := os.Stat(filepath.Join(detect.TemplatesDir(home), "Script.sh")); err == nil {
		t.Error("template not removed by revert")
	}
}
