package integrations

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func setupFixture(t *testing.T) (*Preparer, *SetupReview, string) {
	t.Helper()
	pack := Package{Format: 2, ID: "officegest", Name: "Fixture", Version: "0.2.0", Setup: &ComposeSetup{Compose: "docker-compose.yml", Application: "app", Worker: "worker", Initializer: "redis-cluster-init", Template: "docker/bin/tenant-create.sh", Tools: []string{"git"}, Databases: []DatabaseSetup{{Service: "mariadb", Schemas: []string{"tenant1"}}}, Redis: []RedisNode{{Service: "redis-node-1", Port: 7001}, {Service: "redis-node-2", Port: 7002}, {Service: "redis-node-3", Port: 7003}}, Tenants: []TenantSetup{{Domain: "og-tenant1.test", Service: "mariadb", Schema: "tenant1"}}}}
	var err error
	root := t.TempDir()
	m := Manager{Data: filepath.Join(root, "data"), Config: filepath.Join(root, "config"), Simulated: true}
	if err = m.Install(pack); err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(root, "project")
	if err = os.Mkdir(project, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(project, "sentinel"), []byte("untouched"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = m.SetProject(pack.ID, project); err != nil {
		t.Fatal(err)
	}
	p := &Preparer{Manager: m}
	r, err := p.Review(context.Background(), pack.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if p.Running() {
			_ = p.Control("cancel")
		}
		deadline := time.Now().Add(3 * time.Second)
		for p.Running() && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
		}
		if p.Running() {
			t.Error("preparation did not stop")
		}
	})
	return p, r, project
}

func TestSetupOrphanProcessLockBlocksReviewAndTracksRestart(t *testing.T) {
	p, _, _ := setupFixture(t)
	path := filepath.Join(p.Manager.Data, "setup.lock")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		t.Fatal(err)
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if !p.Running() {
		t.Fatal("orphan lock ignored")
	}
	if _, err = p.Review(context.Background(), "officegest"); err == nil {
		t.Fatal("review raced orphan")
	}
	b, _ := json.Marshal(SetupJob{ID: "orphan", Package: "officegest", State: "running"})
	if err = os.WriteFile(filepath.Join(p.Manager.Data, "setup-last.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	j := p.Status()
	if j == nil || j.Finished {
		t.Fatal("live orphan reported interrupted")
	}
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		t.Fatal(err)
	}
	if p.Running() {
		t.Fatal("released lock still blocks")
	}
	j = p.Status()
	if !j.Finished || j.State != "interrupted" {
		t.Fatal("stopped orphan not marked interrupted")
	}
}

func waitSetup(t *testing.T, p *Preparer, want string) *SetupJob {
	t.Helper()
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		j := p.Status()
		if j != nil && j.State == want {
			return j
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("expected %s, got %+v", want, p.Status())
	return nil
}

func TestSetupReviewBindingAndCopies(t *testing.T) {
	p, r, project := setupFixture(t)
	r.Ports["http"] = 1
	r.Steps[0] = "changed"
	r.Databases[0].Schemas[0] = "changed"
	if p.review.Ports["http"] == 1 || p.review.Steps[0] == "changed" || p.review.Databases[0].Schemas[0] == "changed" {
		t.Fatal("review exposes mutable internal state")
	}
	if _, err := p.Start("officegest", "wrong"); err == nil {
		t.Fatal("wrong token accepted")
	}
	other := t.TempDir()
	if err := p.Manager.SetProject("officegest", other); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Start("officegest", r.Token); err == nil {
		t.Fatal("changed project accepted")
	}
	if b, _ := os.ReadFile(filepath.Join(project, "sentinel")); string(b) != "untouched" {
		t.Fatal("review changed project")
	}
}

func TestSetupSimulationPauseResumeAndPersistence(t *testing.T) {
	p, r, project := setupFixture(t)
	j, err := p.Start("officegest", r.Token)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = p.Start("officegest", r.Token); err == nil {
		t.Fatal("concurrent preparation accepted")
	}
	if err = p.Control("pause"); err != nil {
		t.Fatal(err)
	}
	paused := waitSetup(t, p, "paused")
	snapshot, _ := json.Marshal(paused.Steps)
	time.Sleep(150 * time.Millisecond)
	got, _ := json.Marshal(p.Status().Steps)
	if string(got) != string(snapshot) {
		t.Fatal("steps advanced while paused")
	}
	if err = p.Control("resume"); err != nil {
		t.Fatal(err)
	}
	finished := waitSetup(t, p, "done")
	if !finished.Finished || !finished.Simulated || finished.ID != j.ID {
		t.Fatalf("bad result: %+v", finished)
	}
	for _, s := range finished.Steps {
		if s.Status != "done" {
			t.Fatal("step incomplete")
		}
	}
	b, _ := os.ReadFile(filepath.Join(project, "sentinel"))
	entries, _ := os.ReadDir(project)
	if string(b) != "untouched" || len(entries) != 1 {
		t.Fatal("simulation mutated project")
	}
	for p.Running() {
		time.Sleep(10 * time.Millisecond)
	}
	loaded := (&Preparer{Manager: p.Manager}).Status()
	if loaded == nil || loaded.ID != j.ID || loaded.State != "done" {
		t.Fatalf("persistence: %+v", loaded)
	}
}

func TestSetupCancelWhilePausedAndRestart(t *testing.T) {
	p, r, _ := setupFixture(t)
	if _, err := p.Start("officegest", r.Token); err != nil {
		t.Fatal(err)
	}
	if err := p.Control("pause"); err != nil {
		t.Fatal(err)
	}
	waitSetup(t, p, "paused")
	restarted := (&Preparer{Manager: p.Manager}).Status()
	if restarted == nil || !restarted.Finished || restarted.State != "interrupted" {
		t.Fatalf("restart: %+v", restarted)
	}
	if err := p.Control("cancel"); err != nil {
		t.Fatal(err)
	}
	if err := p.Control("resume"); err == nil {
		t.Fatal("cancelled run resumed")
	}
	waitSetup(t, p, "cancelled")
}

func TestSetupValidationRejectsAmbiguousOrUnsafeData(t *testing.T) {
	_, r, _ := setupFixture(t)
	original, _ := json.Marshal(r.view.Setup)
	changes := []func(*ComposeSetup){
		func(s *ComposeSetup) { s.Compose = "../docker-compose.yml" },
		func(s *ComposeSetup) { s.Template = "/etc/passwd" },
		func(s *ComposeSetup) { s.Redis[0].Port = 60000 },
		func(s *ComposeSetup) { s.Redis[0].Service = s.Application },
		func(s *ComposeSetup) { s.Databases[0].Service = s.Application },
		func(s *ComposeSetup) { s.Databases[0].Schemas = []string{"mysql"} },
		func(s *ComposeSetup) { s.Databases[0].Schemas = []string{"tenant1", "tenant1"} },
		func(s *ComposeSetup) { s.Tenants[0].Domain = "name;command.test" },
	}
	for i, change := range changes {
		var s ComposeSetup
		_ = json.Unmarshal(original, &s)
		change(&s)
		if s.Validate() == nil {
			t.Errorf("unsafe case %d accepted", i)
		}
	}
}
