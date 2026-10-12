package engine

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/plan"
)

func TestNPMUserInstallsAfterDependenciesAndVerifiesTheAgent(t *testing.T) {
	f := setup(t)
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"codex-cli": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.eng.Start(p, nil); err == nil {
		t.Fatal("unconfirmed agent installation must be refused")
	}
	id, err := f.eng.Start(p, map[string]string{"codex-cli": "yes"})
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for f.eng.Running() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	status := ""
	for _, run := range f.st.Runs() {
		if run.ID == id {
			status = run.Status
		}
	}
	node, agent := -1, -1
	for i, step := range p.Steps {
		if step.Item == "nodejs" && step.Op == "apt-install" {
			node = i
		}
		if step.Op == catalog.MethodNPMUser {
			agent = i
			if step.Kind != "user" || step.Confirm == "" {
				t.Fatal("agent must require confirmation and run as user")
			}
		}
	}
	if node < 0 || agent <= node {
		t.Fatalf("dependency order: node=%d agent=%d", node, agent)
	}
	if status != "done" {
		t.Fatalf("agent install did not verify: %#v", f.st.Runs())
	}
	joined := strings.Join(f.sim.Calls, "\n")
	if !strings.Contains(joined, "--prefix "+filepath.Join(f.eng.Home, ".local")) || !strings.Contains(joined, "@openai/codex@0.162.1") {
		t.Fatalf("wrong package or install destination: %s", joined)
	}
	if strings.Contains(joined, "sudo") {
		t.Fatal("npm install must never use sudo")
	}
}

func TestNPMUserPreservesExistingAgentWithoutInstallingNode(t *testing.T) {
	f := setup(t)
	f.sim.NPM["codex"] = "0.160.0"
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"codex-cli": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Steps) != 0 || p.Selection["nodejs"].Options != nil {
		t.Fatalf("existing agent should be preserved: %#v", p)
	}
	if _, found := p.Selection["nodejs"]; found {
		t.Fatal("existing agent must not add an unnecessary runtime")
	}
}

func TestNPMUserFailsBeforeInstallationWithoutRuntime(t *testing.T) {
	f := setup(t)
	if err := f.eng.installNPMUser(context.Background(), f.cat.Item("codex-cli"), func(string) {}); err == nil {
		t.Fatal("missing Node/npm should stop installation")
	}
	for _, call := range f.sim.Calls {
		if strings.HasPrefix(call, "npm install") {
			t.Fatal("installation ran without runtime")
		}
	}
}

func TestPiDisablesDependencyLifecycleScripts(t *testing.T) {
	f := setup(t)
	f.sim.Packages["nodejs"], f.sim.Packages["npm"] = true, true
	if err := f.eng.installNPMUser(context.Background(), f.cat.Item("pi-coding-agent"), func(string) {}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(f.sim.Calls, "\n"), "--ignore-scripts") {
		t.Fatal("Pi must follow its documented install flags")
	}
}
