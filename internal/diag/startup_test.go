package diag

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"meuubuntu/internal/sys"
)

func TestStartupDistinguishesMissingDisabledAndAvailableCommands(t *testing.T) {
	tests := []struct {
		name, desktop, extra, status, command string
	}{
		{"removed Skype", "GNOME", "Exec=/usr/bin/skypeforlinux\nX-GNOME-Autostart-enabled=true", "missing", "/usr/bin/skypeforlinux"},
		{"existing app", "ubuntu:GNOME", "Exec=example --minimized %U\nOnlyShowIn=GNOME;", "available", "example"},
		{"menu visibility does not disable startup", "GNOME", "Exec=example\nNoDisplay=true", "available", "example"},
		{"GNOME switch disabled", "GNOME", "Exec=example\nX-GNOME-Autostart-enabled=false", "disabled", ""},
		{"hidden entry", "GNOME", "Exec=example\nHidden=true", "disabled", ""},
		{"different session", "GNOME", "Exec=example\nOnlyShowIn=KDE;", "disabled", ""},
		{"excluded session", "ubuntu:GNOME", "Exec=example\nNotShowIn=GNOME;", "disabled", ""},
		{"TryExec guards an available launcher", "GNOME", "Exec=example\nTryExec=removed-target", "missing", "removed-target"},
		{"D-Bus is not proof of removal", "GNOME", "DBusActivatable=true", "unknown", ""},
		{"launcher is not proof of its target", "GNOME", "Exec=env SETTING=value removed-target", "unknown", "env"},
		{"condition is not evaluated", "GNOME", "Exec=example\nAutostartCondition=if-exists some-file", "unknown", "example"},
		{"action section cannot disable main entry", "GNOME", "Exec=example\n[Desktop Action stop]\nHidden=true\nName=Stop", "available", "example"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := StartupEntry{Name: "fallback"}
			keys := desktopKeys("[Desktop Entry]\nType=Application\nName=Example\n" + tt.extra + "\n")
			inspectStartup(&entry, keys, tt.desktop, func(name string) bool { return name == "example" || name == "env" })
			if entry.Status != tt.status || entry.Command != tt.command || entry.Name != "Example" || entry.Detail == "" {
				t.Fatalf("wrong diagnosis: %+v", entry)
			}
		})
	}
}

func TestDisabledUserEntryOverridesSystemStartup(t *testing.T) {
	user, system := t.TempDir(), t.TempDir()
	write := func(dir, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "example.desktop"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(system, "[Desktop Entry]\nType=Application\nName=System application\nExec=example\n")
	write(user, "[Desktop Entry]\nName=Disabled by user\nHidden=true\n")
	entries := readStartup([]startupDir{{user, "utilizador"}, {system, "sistema"}}, "GNOME", func(string) bool { return true })
	if len(entries) != 1 || entries[0].Status != "disabled" || entries[0].Scope != "utilizador" || entries[0].Name != "Disabled by user" {
		t.Fatalf("user override lost: %+v", entries)
	}
}

func TestStartupInspectsQuotedExecutableWithoutRunningIt(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "must-not-exist")
	command := filepath.Join(dir, "application with spaces")
	if err := os.WriteFile(command, []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "test.desktop")
	if err := os.WriteFile(path, []byte("[Desktop Entry]\nType=Application\nName=Quoted application\nExec=\""+command+"\" --option %U\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries := readStartup([]startupDir{{dir, "utilizador"}}, "GNOME", sys.Real{}.LookPath)
	if len(entries) != 1 || entries[0].Status != "available" || entries[0].Command != command {
		t.Fatalf("quoted executable was not found: %+v", entries)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("diagnostics must never execute an autostart command")
	}
	if err := os.Chmod(command, 0o644); err != nil {
		t.Fatal(err)
	}
	entries = readStartup([]startupDir{{dir, "utilizador"}}, "GNOME", sys.Real{}.LookPath)
	if entries[0].Status != "missing" {
		t.Fatal("a non-executable file is not an available command")
	}
}

func TestStartupUsesXDGConfigPaths(t *testing.T) {
	user, system := t.TempDir(), t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", user)
	t.Setenv("XDG_CONFIG_DIRS", system)
	t.Setenv("XDG_CURRENT_DESKTOP", "GNOME")
	dir := filepath.Join(user, "autostart")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "skype.desktop"), []byte("[Desktop Entry]\nType=Application\nName=Skype for Linux\nExec=/missing-skype-test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries := startup(sys.Real{})
	if len(entries) != 1 || entries[0].Name != "Skype for Linux" || entries[0].Status != "missing" || !strings.HasPrefix(entries[0].Path, user) {
		t.Fatalf("XDG startup diagnosis incorrect: %+v", entries)
	}
}
