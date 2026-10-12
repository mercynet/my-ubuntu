package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/plan"
)

func TestLocalServersRequireConfirmationAndReapplyWithoutChanges(t *testing.T) {
	f := setup(t)
	ids := []string{"apache2", "nginx", "redis-server", "rabbitmq-server", "postgresql-server", "memcached"}
	selection := map[string]catalog.Selection{}
	confirmed := map[string]string{}
	for _, id := range ids {
		selection[id] = catalog.Selection{}
		confirmed[id] = "yes"
	}
	p, err := plan.Build(context.Background(), f.cat, f.det, selection, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !p.NeedsRoot || p.Blocked || len(p.Warnings) < len(ids) {
		t.Fatalf("server plan must explain services and allow deliberate port configuration: %+v", p)
	}
	if _, err := f.eng.Start(p, nil); err == nil {
		t.Fatal("server installation must require explicit service confirmation")
	}
	runID, err := f.eng.Start(p, confirmed)
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for f.eng.Running() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if f.eng.Running() {
		t.Fatal("service installation did not finish")
	}
	status := ""
	for _, run := range f.st.Runs() {
		if run.ID == runID {
			status = run.Status
		}
	}
	if status != "done" {
		t.Fatalf("service packages must also pass installation verification: %+v", f.st.Runs())
	}
	for _, id := range ids {
		for _, pkg := range f.cat.Item(id).Packages {
			if !f.sim.Packages[pkg] {
				t.Fatalf("package %s was not installed", pkg)
			}
		}
	}
	for _, call := range f.sim.Calls {
		if strings.Contains(call, "systemctl enable") || strings.Contains(call, "systemctl start") {
			t.Fatalf("application must not add service management commands to package installation: %s", call)
		}
	}
	repeat, err := plan.Build(context.Background(), f.cat, f.det, selection, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(repeat.Steps) != 0 {
		t.Fatalf("installed local servers must not be reinstalled: %+v", repeat.Steps)
	}
}

func TestKafkaGuideDoesNotInstallASecondDockerEngine(t *testing.T) {
	f := setup(t)
	for _, pkg := range f.cat.Item("docker-ubuntu").Packages {
		f.sim.Packages[pkg] = true
	}
	p, err := plan.Build(context.Background(), f.cat, f.det, map[string]catalog.Selection{"kafka-local-guide": {}, "docker-ubuntu": {}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Blocked || len(p.Steps) != 0 || len(p.Manual) != 1 {
		t.Fatalf("Kafka stays a manual lab and preserves the existing Docker choice: %+v", p)
	}
	if _, added := p.Selection["docker-engine"]; added {
		t.Fatal("Kafka must not force the other Docker distribution")
	}
}
