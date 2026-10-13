package sys

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Sim is a controlled executor for tests and demo mode. It never runs real
// commands; it models a small Ubuntu system in memory. It does NOT prove
// that real installations work.
type Sim struct {
	mu             sync.Mutex
	Packages       map[string]bool
	Flatpaks       map[string]bool
	FlatpakRemotes map[string]string
	Snaps          map[string]bool
	Settings       map[string]string // "schema key" -> value
	Defaults       map[string]string
	Exts           map[string]bool // installed extension uuids
	Git            map[string]string
	Groups         map[string]bool
	Profile        string
	FailPkgs       map[string]bool // apt-get install of these fails
	Delay          time.Duration
	Calls          []string
	Available      map[string]bool   // apt candidates; nil = everything available
	NPM            map[string]string // simulated CLI binary -> version
	XdebugDisabled map[string]bool
}

func (s *Sim) XdebugConfigured(version string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.XdebugDisabled[version]
}

func NewSim() *Sim {
	s := &Sim{
		XdebugDisabled: map[string]bool{},
		Packages: map[string]bool{"git": true, "curl": true, "snapd": true, "flatpak": false, "gnome-shell": true,
			"gnome-shell-extension-ubuntu-dock": true, "gnome-shell-extension-appindicator": true,
			"gnome-shell-extension-ubuntu-tiling-assistant": true, "power-profiles-daemon": true},
		Flatpaks: map[string]bool{}, Snaps: map[string]bool{"firefox": true, "snapd": true},
		FlatpakRemotes: map[string]string{},
		Settings: map[string]string{
			"org.gnome.desktop.interface color-scheme":                                                      "'default'",
			"org.gnome.desktop.interface text-scaling-factor":                                               "1.0",
			"org.gnome.shell.extensions.dash-to-dock dock-position":                                         "'LEFT'",
			"org.gnome.shell.extensions.dash-to-dock dash-max-icon-size":                                    "48",
			"org.gnome.shell.extensions.dash-to-dock dock-fixed":                                            "true",
			"org.gnome.desktop.interface accent-color":                                                      "'orange'",
			"org.gnome.Ptyxis default-profile-uuid":                                                         "'4f6a1c2e9b8d4e7fa1b2c3d4e5f60718'",
			"org.gnome.Ptyxis.Profile:/org/gnome/Ptyxis/Profiles/4f6a1c2e9b8d4e7fa1b2c3d4e5f60718/ opacity": "1.0",
		},
		Exts:     map[string]bool{"ubuntu-dock@ubuntu.com": true, "ubuntu-appindicators@ubuntu.com": true, "tiling-assistant@ubuntu.com": true, "ding@rastersoft.com": true},
		Git:      map[string]string{},
		Groups:   map[string]bool{},
		Profile:  "balanced",
		FailPkgs: map[string]bool{},
		NPM:      map[string]string{},
		Delay:    150 * time.Millisecond,
	}
	s.Defaults = map[string]string{}
	for k, v := range s.Settings {
		s.Defaults[k] = v
	}
	for _, p := range strings.Split(os.Getenv("MEU_UBUNTU_SIM_FAIL"), ",") {
		if p != "" {
			s.FailPkgs[p] = true
		}
	}
	return s
}

func (s *Sim) LookPath(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.NPM[filepath.Base(name)]; ok {
		return true
	}
	switch name {
	case "node", "npm":
		return s.Packages["nodejs"] && s.Packages["npm"]
	case "flatpak":
		return s.Packages["flatpak"]
	case "docker":
		return s.Packages["docker-ce"] || s.Packages["docker.io"]
	case "gsettings", "git", "dpkg-query", "apt-cache", "gnome-extensions", "powerprofilesctl", "snap", "gpg":
		return true
	}
	return s.Packages[name]
}

func (s *Sim) Output(ctx context.Context, name string, args ...string) (string, error) {
	var b strings.Builder
	err := s.Stream(ctx, func(l string) { b.WriteString(l + "\n") }, name, args...)
	return b.String(), err
}

func (s *Sim) Stream(ctx context.Context, onLine func(string), name string, args ...string) error {
	s.mu.Lock()
	s.Calls = append(s.Calls, Argv(name, args...))
	s.mu.Unlock()
	s.mu.Lock()
	version, isNPM := s.NPM[filepath.Base(name)]
	s.mu.Unlock()
	if isNPM {
		onLine(version)
		return nil
	}
	a := strings.Join(args, " ")
	switch name {
	case "node":
		onLine("v22.22.1")
		return nil
	case "npm":
		if len(args) == 0 {
			return fmt.Errorf("npm: argumentos em falta")
		}
		spec := args[len(args)-1]
		at := strings.LastIndex(spec, "@")
		if at < 1 {
			return fmt.Errorf("npm: versão fixa em falta")
		}
		binaries := map[string]string{"@anthropic-ai/claude-code": "claude", "@openai/codex": "codex", "@google/gemini-cli": "gemini", "opencode-ai": "opencode", "@earendil-works/pi-coding-agent": "pi", "@github/copilot": "copilot"}
		bin := binaries[spec[:at]]
		if bin == "" {
			return fmt.Errorf("npm: pacote desconhecido na simulação")
		}
		if s.FailPkgs[spec[:at]] {
			return fmt.Errorf("npm: falha simulada")
		}
		s.mu.Lock()
		if s.NPM == nil {
			s.NPM = map[string]string{}
		}
		s.NPM[bin] = spec[at+1:]
		s.mu.Unlock()
		onLine("Instalado " + spec + " (simulado)")
		return nil
	case "dpkg-query":
		s.mu.Lock()
		keys := sortedKeys(s.Packages)
		for _, k := range keys {
			st := "not-installed"
			if s.Packages[k] {
				st = "installed"
			}
			onLine(k + "\t" + st)
		}
		s.mu.Unlock()
		return nil
	case "apt-cache":
		for _, p := range args[1:] {
			if s.Available == nil || s.Available[p] {
				onLine(p + ":\n  Installed: (none)\n  Candidate: 1.0-sim")
			} else {
				onLine(p + ":\n  Installed: (none)\n  Candidate: (none)")
			}
		}
		return nil
	case "apt-get":
		if len(args) > 0 && args[len(args)-1] == "update" || strings.Contains(a, " update") {
			onLine("Hit:1 http://archive.ubuntu.com/ubuntu resolute InRelease (simulado)")
			return nil
		}
		var pkgs []string
		for _, x := range args {
			if !strings.HasPrefix(x, "-") && x != "install" && x != "remove" && !strings.Contains(x, "::") {
				pkgs = append(pkgs, x)
			}
		}
		for _, p := range pkgs {
			time.Sleep(s.Delay)
			if s.FailPkgs[p] {
				onLine("E: Unable to locate package " + p + " (simulado)")
				return fmt.Errorf("exit status 100")
			}
			onLine("Setting up " + p + " (1.0-sim) ... (simulado)")
			s.mu.Lock()
			s.Packages[p] = true
			if strings.HasPrefix(p, "php") && strings.HasSuffix(p, "-xdebug") {
				s.XdebugDisabled[strings.TrimSuffix(strings.TrimPrefix(p, "php"), "-xdebug")] = false
			}
			s.mu.Unlock()
		}
		return nil
	case "flatpak":
		s.mu.Lock()
		defer s.mu.Unlock()
		switch {
		case len(args) > 0 && args[0] == "--version":
			onLine("Flatpak 1.16.5 (simulado)")
		case len(args) > 0 && args[0] == "remote-add":
			if s.FlatpakRemotes == nil {
				s.FlatpakRemotes = map[string]string{}
			}
			s.FlatpakRemotes["flathub"] = "https://dl.flathub.org/repo/"
			onLine("Flathub configurado (simulado)")
		case len(args) > 0 && args[0] == "remote-info":
			id := args[len(args)-1]
			if s.FailPkgs[id] {
				return ExitStatus(1)
			}
			onLine("app/" + id + "/x86_64/stable")
		case len(args) > 0 && args[0] == "list":
			for k, v := range s.Flatpaks {
				if v {
					onLine(k)
				}
			}
		case len(args) > 0 && args[0] == "install":
			id := args[len(args)-1]
			if strings.Contains(a, "--from") {
				ref, err := os.ReadFile(id)
				if err != nil {
					return err
				}
				for _, line := range strings.Split(string(ref), "\n") {
					fields := strings.SplitN(line, "=", 2)
					if len(fields) == 2 && strings.TrimSpace(fields[0]) == "Name" {
						id = strings.TrimSpace(fields[1])
					}
				}
			}
			time.Sleep(s.Delay)
			if s.FailPkgs[id] {
				onLine("error: simulado")
				return ExitStatus(1)
			}
			s.Flatpaks[id] = true
			onLine("Installing " + id + " (simulado)")
		case len(args) > 0 && args[0] == "remotes":
			for name, remote := range s.FlatpakRemotes {
				onLine(name + "\t" + remote + "\t")
			}
		default:
			onLine("ok (simulado)")
		}
		return nil
	case "snap":
		s.mu.Lock()
		defer s.mu.Unlock()
		if len(args) > 0 && args[0] == "list" {
			onLine("Name  Version  Rev  Tracking  Publisher  Notes")
			for _, k := range sortedKeys(s.Snaps) {
				if s.Snaps[k] {
					onLine(k + "  1.0  1  latest/stable  canonical  -")
				}
			}
		}
		if len(args) > 1 && args[0] == "remove" {
			delete(s.Snaps, args[len(args)-1])
			onLine(args[len(args)-1] + " removed (simulado)")
		}
		return nil
	case "phpdismod":
		s.mu.Lock()
		defer s.mu.Unlock()
		if len(args) >= 3 && args[0] == "-v" {
			s.XdebugDisabled[args[1]] = true
		}
		return nil
	case "gsettings":
		s.mu.Lock()
		defer s.mu.Unlock()
		switch args[0] {
		case "get":
			v, ok := s.Settings[args[1]+" "+args[2]]
			if !ok {
				if strings.HasPrefix(args[1], "org.gnome") {
					onLine("false")
					return nil
				}
				return fmt.Errorf("No such schema “%s”", args[1])
			}
			onLine(v)
		case "set":
			s.Settings[args[1]+" "+args[2]] = args[3]
		case "reset":
			if d, ok := s.Defaults[args[1]+" "+args[2]]; ok {
				s.Settings[args[1]+" "+args[2]] = d
			} else {
				delete(s.Settings, args[1]+" "+args[2])
			}
		case "list-schemas":
			onLine("org.gnome.desktop.interface")
		}
		return nil
	case "gnome-extensions":
		s.mu.Lock()
		defer s.mu.Unlock()
		if args[0] == "list" {
			for _, k := range sortedKeys(s.Exts) {
				onLine(k)
			}
		}
		if args[0] == "install" {
			s.Exts["installed-sim@x"] = true
		}
		return nil
	case "git":
		s.mu.Lock()
		defer s.mu.Unlock()
		// git config --global [--get|--unset] key [value]
		if len(args) >= 4 && args[2] == "--get" {
			v, ok := s.Git[args[3]]
			if !ok {
				return ExitStatus(1)
			}
			onLine(v)
		} else if len(args) >= 4 && args[2] == "--unset" {
			delete(s.Git, args[3])
		} else if len(args) >= 4 {
			s.Git[args[2]] = args[3]
		}
		return nil
	case "gpg":
		// Simulated key files contain "SIM-FPR:<fpr>" lines (see simFetch).
		if len(args) > 0 {
			b, _ := os.ReadFile(args[len(args)-1])
			for _, l := range strings.Split(string(b), "\n") {
				if f, ok := strings.CutPrefix(l, "SIM-FPR:"); ok {
					onLine("pub:-:4096:1:SIM:1:::-:::scESC::::::23::0:")
					onLine("fpr:::::::::" + f + ":")
				}
			}
		}
		return nil
	case "powerprofilesctl":
		s.mu.Lock()
		defer s.mu.Unlock()
		if len(args) == 0 || args[0] == "get" {
			onLine(s.Profile)
		} else if args[0] == "set" {
			s.Profile = args[1]
		}
		return nil
	case "id":
		s.mu.Lock()
		defer s.mu.Unlock()
		g := "user adm sudo"
		if s.Groups["docker"] {
			g += " docker"
		}
		onLine(g)
		return nil
	case "usermod", "gpasswd":
		s.mu.Lock()
		s.Groups["docker"] = name == "usermod"
		s.mu.Unlock()
		onLine("ok (simulado)")
		return nil
	}
	time.Sleep(s.Delay / 3)
	onLine("(simulado) " + Argv(name, args...))
	return nil
}

func sortedKeys[V any](m map[string]V) []string {
	var k []string
	for x := range m {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}
