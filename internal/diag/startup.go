package diag

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"meuubuntu/internal/sys"
)

// A desktop entry is configuration, not proof that an application is installed
// or running. Inspect executable availability without ever executing its Exec.
type StartupEntry struct {
	Name    string `json:"name"`
	Scope   string `json:"scope"`
	Path    string `json:"path"`
	Status  string `json:"status"` // available | missing | disabled | unknown
	Detail  string `json:"detail"`
	Command string `json:"command,omitempty"`
}

type startupDir struct{ path, scope string }

func startup(r sys.Runner) []StartupEntry {
	home, _ := os.UserHomeDir()
	config := os.Getenv("XDG_CONFIG_HOME")
	if !filepath.IsAbs(config) {
		config = filepath.Join(home, ".config")
	}
	dirs := []startupDir{{filepath.Join(config, "autostart"), "utilizador"}}
	system := os.Getenv("XDG_CONFIG_DIRS")
	if system == "" {
		system = "/etc/xdg"
	}
	for _, dir := range filepath.SplitList(system) {
		if filepath.IsAbs(dir) {
			dirs = append(dirs, startupDir{filepath.Join(dir, "autostart"), "sistema"})
		}
	}
	return readStartup(dirs, os.Getenv("XDG_CURRENT_DESKTOP"), r.LookPath)
}

func readStartup(dirs []startupDir, desktop string, exists func(string) bool) []StartupEntry {
	var out []StartupEntry
	seen := map[string]bool{}
	for _, dir := range dirs {
		files, _ := filepath.Glob(filepath.Join(dir.path, "*.desktop"))
		for _, path := range files {
			id := filepath.Base(path)
			// A user entry overrides the system entry, even when disabled.
			if seen[id] {
				continue
			}
			seen[id] = true
			entry := StartupEntry{Name: id, Scope: dir.scope, Path: path, Status: "unknown"}
			body, err := os.ReadFile(path)
			if err != nil {
				entry.Detail = "Não foi possível ler esta entrada de arranque."
			} else {
				inspectStartup(&entry, desktopKeys(string(body)), desktop, exists)
			}
			out = append(out, entry)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out
}

func desktopKeys(body string) map[string]string {
	keys := map[string]string{}
	group := ""
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = line
			continue
		}
		if group != "[Desktop Entry]" || strings.HasPrefix(line, "#") {
			continue
		}
		if key, value, ok := strings.Cut(line, "="); ok {
			keys[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	return keys
}

func inspectStartup(entry *StartupEntry, keys map[string]string, desktop string, exists func(string) bool) {
	if keys["Name"] != "" {
		entry.Name = keys["Name"]
	}
	set := func(status, detail string) { entry.Status, entry.Detail = status, detail }
	if keys["Hidden"] == "true" || keys["X-GNOME-Autostart-enabled"] == "false" {
		set("disabled", "Esta entrada está desativada; não pede o arranque da aplicação.")
		return
	}
	if keys["Type"] != "Application" {
		set("unknown", "A entrada não identifica uma aplicação válida.")
		return
	}
	if keys["OnlyShowIn"] != "" || keys["NotShowIn"] != "" {
		if desktop == "" {
			set("unknown", "Não foi possível conferir o ambiente de sessão exigido por esta entrada.")
			return
		}
		match := func(list string) bool {
			for _, current := range strings.Split(desktop, ":") {
				for _, allowed := range strings.Split(list, ";") {
					if current != "" && current == allowed {
						return true
					}
				}
			}
			return false
		}
		if (keys["OnlyShowIn"] != "" && !match(keys["OnlyShowIn"])) || match(keys["NotShowIn"]) {
			set("disabled", "Esta entrada não se aplica ao ambiente de sessão atual.")
			return
		}
	}
	if target := keys["TryExec"]; target != "" && !exists(desktopUnescape(target)) {
		entry.Command = desktopUnescape(target)
		set("missing", "O executável exigido pela entrada não foi encontrado. Pode ser uma referência antiga.")
		return
	}
	if keys["DBusActivatable"] == "true" {
		set("unknown", "Esta entrada usa ativação D-Bus; a disponibilidade da aplicação não foi verificada.")
		return
	}
	command := desktopCommand(keys["Exec"])
	if command == "" {
		set("unknown", "Não foi possível identificar o executável desta entrada.")
		return
	}
	entry.Command = command
	if !exists(command) {
		set("missing", "O executável não foi encontrado. A entrada pode ter ficado de uma aplicação removida ou de uma configuração importada.")
		return
	}
	if keys["AutostartCondition"] != "" {
		set("unknown", "Comando localizado; esta entrada tem uma condição de arranque que não foi avaliada.")
		return
	}
	switch filepath.Base(command) {
	case "env", "sh", "bash", "dash", "zsh", "flatpak", "snap", "gio", "gtk-launch":
		set("unknown", "O lançador está disponível; a aplicação ou comando que ele inicia não foi verificado.")
	default:
		set("available", "Comando localizado. Isto não confirma que a aplicação esteja em execução ou que tenha arrancado com a sessão.")
	}
}

func desktopUnescape(value string) string {
	return strings.NewReplacer(`\s`, " ", `\n`, "\n", `\t`, "\t", `\r`, "\r", `\\`, `\`).Replace(value)
}

// Parse only the executable token. Arguments and field codes are never run.
func desktopCommand(value string) string {
	value = strings.TrimSpace(desktopUnescape(value))
	if value == "" {
		return ""
	}
	command := ""
	if value[0] == '"' {
		var token strings.Builder
		closed := false
		for i := 1; i < len(value); i++ {
			if value[i] == '\\' && i+1 < len(value) {
				i++
				token.WriteByte(value[i])
			} else if value[i] == '"' {
				if i+1 < len(value) && !strings.ContainsRune(" \t", rune(value[i+1])) {
					return ""
				}
				closed = true
				break
			} else {
				token.WriteByte(value[i])
			}
		}
		if !closed {
			return ""
		}
		command = token.String()
	} else {
		command = strings.Fields(value)[0]
	}
	if strings.ContainsAny(command, "=%\n\r") {
		return ""
	}
	return command
}
