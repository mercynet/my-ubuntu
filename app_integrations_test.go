package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"meuubuntu/internal/integrations"
)

func integrationArchive(t *testing.T, path, version string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	e, err := w.Create("manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	p := integrations.Package{Format: 1, ID: "example-team", Name: "Example", Version: version}
	if err := json.NewEncoder(e).Encode(p); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestIntegrationPreparationBlocksOtherOperationsAndClose(t *testing.T) {
	a := simulationApp(t)
	pack := integrations.Package{Format: 2, ID: "example-team", Name: "Fixture", Version: "0.2.0", Setup: &integrations.ComposeSetup{Compose: "docker-compose.yml", Application: "app", Initializer: "init", Template: "docker/template.sh", Tools: []string{"git"}, Databases: []integrations.DatabaseSetup{{Service: "db", Schemas: []string{"fixture"}}}, Redis: []integrations.RedisNode{{Service: "redis-1", Port: 7001}, {Service: "redis-2", Port: 7002}, {Service: "redis-3", Port: 7003}}}}
	m := a.integrationManager()
	if err := m.Install(pack); err != nil {
		t.Fatal(err)
	}
	if err := m.SetProject(pack.ID, t.TempDir()); err != nil {
		t.Fatal(err)
	}
	r, err := a.IntegrationSetupReview(pack.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = a.IntegrationSetupStart(pack.ID, r.Token); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = a.Cancel()
		deadline := time.Now().Add(3 * time.Second)
		for a.Running() && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
		}
		if a.Running() {
			t.Error("preparation did not stop")
		}
	})
	if !a.Running() || !a.beforeClose(context.Background()) {
		t.Fatal("active preparation allowed window close")
	}
	if _, err = a.Execute(nil, nil, nil, "invalid"); err == nil {
		t.Fatal("concurrent plan accepted")
	}
	if _, err = a.IntegrationSetupReview(pack.ID); err == nil {
		t.Fatal("concurrent review accepted")
	}
	if err = a.InstallIntegration("invalid"); err == nil {
		t.Fatal("package modification accepted")
	}
	if err = a.Pause(); err != nil {
		t.Fatal(err)
	}
	if err = a.Cancel(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for a.Running() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if a.beforeClose(context.Background()) {
		t.Fatal("finished preparation blocked close")
	}
	if _, err = a.IntegrationSetupReview(pack.ID); err == nil {
		t.Fatal("review accepted during close")
	}
}

func TestIntegrationImportRequiresCurrentReview(t *testing.T) {
	a := simulationApp(t)
	path := filepath.Join(t.TempDir(), "package.zip")
	integrationArchive(t, path, "0.1.0")
	if err := a.InstallIntegration("unreviewed"); err == nil {
		t.Fatal("unreviewed import accepted")
	}
	r, err := a.reviewIntegrationImport(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.InstallIntegration("wrong"); err == nil {
		t.Fatal("wrong token accepted")
	}
	integrationArchive(t, path, "0.2.0")
	if err := a.InstallIntegration(r.Token); err == nil || len(a.Integrations()) != 0 {
		t.Fatal("changed archive installed")
	}
	r, err = a.reviewIntegrationImport(path)
	if err != nil {
		t.Fatal(err)
	}
	a.integrationPending.review.created = time.Now().Add(-11 * time.Minute)
	if err := a.InstallIntegration(r.Token); err == nil {
		t.Fatal("expired review accepted")
	}
	r, err = a.reviewIntegrationImport(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.InstallIntegration(r.Token); err != nil {
		t.Fatal(err)
	}
	if err := a.InstallIntegration(r.Token); err == nil {
		t.Fatal("consumed review reused")
	}
	if packs := a.Integrations(); len(packs) != 1 || packs[0].Version != "0.2.0" {
		t.Fatalf("incorrect installed package: %+v", packs)
	}
}
