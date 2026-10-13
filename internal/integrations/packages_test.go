package integrations

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixturePackage() Package {
	return Package{Format: 1, ID: "example-team", Name: "Example team", Version: "0.1.0", Publisher: "Example", Checks: []Check{{Type: "file", Title: "Project", Path: "compose.yml"}}, Guide: []Guide{{Title: "Prepare", Text: "Review first", Command: "echo manual-only"}}}
}

func fixtureManager(t *testing.T, simulated bool) Manager {
	t.Helper()
	root := t.TempDir()
	return Manager{Data: filepath.Join(root, "data"), Config: filepath.Join(root, "config"), Simulated: simulated}
}

func archive(t *testing.T, entries map[string][]byte, mode os.FileMode) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "package.zip")
	var data bytes.Buffer
	w := zip.NewWriter(&data)
	for name, content := range entries {
		h := &zip.FileHeader{Name: name, Method: zip.Deflate}
		h.SetMode(mode)
		entry, err := w.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write(content); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestArchiveRejectsUntrustedContents(t *testing.T) {
	b, _ := json.Marshal(fixturePackage())
	for _, tc := range []struct {
		name    string
		entries map[string][]byte
		mode    os.FileMode
	}{
		{"traversal", map[string][]byte{"../manifest.json": b}, 0600},
		{"absolute", map[string][]byte{"/manifest.json": b}, 0600},
		{"extra executable", map[string][]byte{"manifest.json": b, "setup.sh": []byte("exit 0")}, 0700},
		{"symlink", map[string][]byte{"manifest.json": []byte("/etc/passwd")}, os.ModeSymlink | 0600},
		{"decompression limit", map[string][]byte{"manifest.json": bytes.Repeat([]byte(" "), 256*1024+1)}, 0600},
		{"unknown command", map[string][]byte{"manifest.json": []byte(`{"format":1,"id":"example-team","name":"Example","version":"1","command":"sh"}`)}, 0600},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := ReadArchive(archive(t, tc.entries, tc.mode)); err == nil {
				t.Fatal("untrusted archive accepted")
			}
		})
	}
	path := archive(t, map[string][]byte{"manifest.json": b}, 0600)
	p, token, err := ReadArchive(path)
	if err != nil || p.ID != "example-team" || len(token) != 64 {
		t.Fatalf("valid archive: %+v %q %v", p, token, err)
	}
	symlink := filepath.Join(t.TempDir(), "linked.zip")
	if err := os.Symlink(path, symlink); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ReadArchive(symlink); err == nil {
		t.Fatal("archive symlink accepted")
	}
	if err := os.WriteFile(path, []byte("not a ZIP"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ReadArchive(path); err == nil {
		t.Fatal("invalid ZIP accepted")
	}
}

func TestDecodeRejectsUnsafeChecks(t *testing.T) {
	for _, path := range []string{"../secret", "/etc/passwd", "a/../../secret", "a\\secret", "a/./secret"} {
		p := fixturePackage()
		p.Checks[0].Path = path
		b, _ := json.Marshal(p)
		if _, err := Decode(b); err == nil {
			t.Errorf("unsafe path %q accepted", path)
		}
	}
	p := fixturePackage()
	p.Checks[0].Type = "shell"
	b, _ := json.Marshal(p)
	if _, err := Decode(b); err == nil {
		t.Fatal("arbitrary check type accepted")
	}
}

func TestPackageAndProjectPersistAcrossUpdates(t *testing.T) {
	m := fixtureManager(t, false)
	p := fixturePackage()
	if err := m.Install(p); err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	if err := m.SetProject(p.ID, project); err != nil {
		t.Fatal(err)
	}
	p.Version = "0.2.0"
	if err := m.Install(p); err != nil {
		t.Fatal(err)
	}
	reopened := Manager{Data: m.Data, Config: m.Config}
	v, err := reopened.Load(p.ID)
	if err != nil || v.Version != "0.2.0" || v.Project != project || len(reopened.List()) != 1 {
		t.Fatalf("lost update/project: %+v %v", v, err)
	}
	for _, path := range []string{filepath.Join(m.root(), p.ID+".json"), filepath.Join(m.Config, "integrations", p.ID+".json")} {
		i, err := os.Stat(path)
		if err != nil || i.Mode().Perm() != 0600 {
			t.Fatalf("private file permissions: %s %v", path, err)
		}
	}
	if err := reopened.SetProject(p.ID, "/"); err == nil {
		t.Fatal("root project accepted")
	}
	if err := reopened.SetProject("../escape", project); err == nil {
		t.Fatal("unknown/traversal ID accepted")
	}
}

func TestInstallRejectsSymlinkedStorage(t *testing.T) {
	m := fixtureManager(t, false)
	if err := os.MkdirAll(m.Data, 0700); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, m.root()); err != nil {
		t.Fatal(err)
	}
	if err := m.Install(fixturePackage()); err == nil {
		t.Fatal("symlinked package storage accepted")
	}
	entries, _ := os.ReadDir(outside)
	if len(entries) != 0 {
		t.Fatal("wrote through symlinked storage")
	}
}

func TestSimulationNeverCallsDockerOrGuide(t *testing.T) {
	m := fixtureManager(t, true)
	p := fixturePackage()
	p.Checks = []Check{{Type: "docker"}, {Type: "compose"}, {Type: "daemon"}, {Type: "compose-config"}, {Type: "file", Path: "compose.yml"}}
	marker := filepath.Join(t.TempDir(), "executed")
	p.Guide[0].Command = "touch " + marker
	if err := m.Install(p); err != nil {
		t.Fatal(err)
	}
	if err := m.SetProject(p.ID, t.TempDir()); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "docker"), []byte("#!/bin/sh\ntouch '"+marker+"'\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+":"+os.Getenv("PATH"))
	rep, err := m.Diagnose(context.Background(), p.ID)
	if err != nil || !rep.Simulated || len(rep.Results) != 5 {
		t.Fatalf("simulation report: %+v %v", rep, err)
	}
	for _, r := range rep.Results {
		if r.Status != "example" {
			t.Errorf("simulation claimed a real result: %+v", r)
		}
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("simulation executed a command")
	}
}

func TestFileDiagnosticsRejectEscapeAndDoNotReadContents(t *testing.T) {
	m := fixtureManager(t, false)
	project := t.TempDir()
	outside := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(outside, []byte("must-not-appear"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(project, "escape")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "present"), []byte("must-not-appear"), 0600); err != nil {
		t.Fatal(err)
	}
	p := fixturePackage()
	p.Checks = []Check{{Type: "file", Path: "escape"}, {Type: "file", Path: "missing"}, {Type: "file", Path: "present"}}
	if err := m.Install(p); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Diagnose(context.Background(), p.ID); err == nil {
		t.Fatal("diagnostics without a project accepted")
	}
	if err := m.SetProject(p.ID, project); err != nil {
		t.Fatal(err)
	}
	rep, err := m.Diagnose(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{"blocked", "missing", "ok"} {
		if rep.Results[i].Status != want || strings.Contains(rep.Results[i].Detail, "must-not-appear") {
			t.Fatalf("result %d: %+v", i, rep.Results[i])
		}
	}
}
