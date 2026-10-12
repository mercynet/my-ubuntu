// Package diag performs read-only hardware/system diagnostics and derives
// contextual suggestions, each with the evidence that triggered it.
package diag

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"meuubuntu/internal/sys"
)

type Disk struct {
	Name     string `json:"name"`
	Model    string `json:"model"`
	Size     string `json:"size"`
	SSD      bool   `json:"ssd"`
	Discard  bool   `json:"discard"`
	Mount    string `json:"mount"`
	FreePct  int    `json:"freePct"`
	FreeHint string `json:"free"`
}

type Fact struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Suggestion struct {
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Evidence string `json:"evidence"`
	Item     string `json:"item,omitempty"`
	Level    string `json:"level"` // info | suggest | warn
}

type Report struct {
	SnapshotID  string         `json:"snapshotId"`
	CollectedAt string         `json:"collectedAt"`
	Facts       []Fact         `json:"facts"`
	Disks       []Disk         `json:"disks"`
	Startup     []StartupEntry `json:"startup"`
	Suggestions []Suggestion   `json:"suggestions"`
	Battery     bool           `json:"battery"`
}

func Run(ctx context.Context, r sys.Runner, simulated bool) Report {
	if simulated {
		return simReport()
	}
	var rep Report
	add := func(l, v string) {
		if strings.TrimSpace(v) != "" {
			rep.Facts = append(rep.Facts, Fact{l, strings.TrimSpace(v)})
		}
	}
	cpu, cores := cpuInfo()
	add("CPU", fmt.Sprintf("%s (%d threads)", cpu, cores))
	memKB, swapKB := memInfo()
	add("RAM", gb(memKB))
	add("Swap", gb(swapKB))
	if out, err := r.Output(ctx, "zramctl", "--noheadings"); err == nil && strings.TrimSpace(out) != "" {
		add("zram", strings.Fields(out)[0]+" ativo")
	} else {
		add("zram", "não ativo")
	}
	if out, err := r.Output(ctx, "lspci"); err == nil {
		for _, l := range strings.Split(out, "\n") {
			if strings.Contains(l, "VGA") || strings.Contains(l, "3D controller") {
				_, v, _ := strings.Cut(l, ": ")
				add("GPU", v)
			}
		}
	}
	add("Sessão", os.Getenv("XDG_SESSION_TYPE"))
	add("Ambiente", os.Getenv("XDG_CURRENT_DESKTOP"))
	add("Ecrãs", resolutions())
	pm := "desconhecido"
	if r.LookPath("powerprofilesctl") {
		out, _ := r.Output(ctx, "powerprofilesctl", "get")
		pm = "power-profiles-daemon (perfil: " + strings.TrimSpace(out) + ")"
	}
	tlp := r.LookPath("tlp")
	if tlp {
		pm += " + TLP instalado"
	}
	add("Gestor de energia", pm)
	bat := battery()
	rep.Battery = bat != ""
	if bat == "" {
		bat = "sem bateria"
	}
	add("Bateria", bat)
	trim, _ := r.Output(ctx, "systemctl", "is-enabled", "fstrim.timer")
	trim = strings.TrimSpace(trim)
	add("fstrim.timer", trim)
	jd, _ := r.Output(ctx, "journalctl", "--disk-usage")
	add("Journal (visível ao utilizador)", jd)
	inot, _ := os.ReadFile("/proc/sys/fs/inotify/max_user_watches")
	add("inotify max_user_watches", string(inot))
	if out, err := r.Output(ctx, "ubuntu-drivers", "list"); err == nil {
		add("Drivers adicionais", strings.Join(strings.Fields(out), ", "))
	}
	rep.Disks = disks(ctx, r)
	rep.Startup = startup(r)

	anySSD := false
	for _, d := range rep.Disks {
		anySSD = anySSD || (d.SSD && d.Discard)
		if d.Mount != "" && d.FreePct < 15 {
			rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "Pouco espaço em " + d.Mount, Detail: "Use ncdu para ver o que ocupa espaço antes de apagar algo.", Evidence: fmt.Sprintf("%d%% livre", d.FreePct), Item: "ncdu", Level: "warn"})
		}
	}
	if anySSD && trim != "enabled" {
		rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "TRIM periódico desativado", Detail: "Há SSD com suporte a discard; o padrão do Ubuntu é fstrim.timer semanal.", Evidence: "fstrim.timer: " + trim, Item: "fstrim-timer", Level: "suggest"})
	}
	if tlp && r.LookPath("powerprofilesctl") {
		rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "Dois gestores de energia", Detail: "TLP e power-profiles-daemon competem pelos mesmos ajustes. Mantenha apenas um.", Evidence: "tlp e powerprofilesctl presentes", Level: "warn"})
	}
	if rep.Battery {
		rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "Portátil detetado", Detail: "Considere o perfil «Portátil equilibrado».", Evidence: bat, Item: "power-balanced", Level: "suggest"})
	} else {
		rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "Sem bateria (ligado à corrente)", Detail: "O perfil «Desempenho» pode fazer sentido; o ganho depende do CPU e não foi medido.", Evidence: "sem bateria", Item: "power-performance", Level: "info"})
	}
	if f := parseSizeG(jd); f >= 1.0 {
		rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "Journal grande", Detail: "Limitar a retenção liberta espaço sem desligar os logs.", Evidence: strings.TrimSpace(jd), Item: "journald-limit", Level: "suggest"})
	}
	if n, _ := strconv.Atoi(strings.TrimSpace(string(inot))); n > 0 && n < 100000 {
		rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "Limite de inotify baixo", Detail: "Só aumente se o VS Code/Vite avisar que não consegue observar ficheiros.", Evidence: fmt.Sprintf("max_user_watches=%d", n), Item: "inotify-watches", Level: "info"})
	}
	if memKB > 0 && memKB < 8*1024*1024 && swapKB < 2*1024*1024 {
		rep.Suggestions = append(rep.Suggestions, Suggestion{Title: "Pouca RAM e swap", Detail: "Docker + IDE podem esgotar a memória. Avalie swap/zram manualmente; não é alterado automaticamente.", Evidence: gb(memKB) + " RAM, " + gb(swapKB) + " swap", Level: "info"})
	}
	return rep
}

func parseSizeG(s string) float64 {
	for _, f := range strings.Fields(s) {
		if strings.HasSuffix(f, "G") {
			v, _ := strconv.ParseFloat(strings.TrimSuffix(f, "G"), 64)
			return v
		}
	}
	return 0
}

func gb(kb int64) string { return fmt.Sprintf("%.1f GiB", float64(kb)/1024/1024) }

func cpuInfo() (string, int) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", 0
	}
	defer f.Close()
	name, n := "", 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		k, v, _ := strings.Cut(sc.Text(), ":")
		k = strings.TrimSpace(k)
		if k == "model name" && name == "" {
			name = strings.TrimSpace(v)
		}
		if k == "processor" {
			n++
		}
	}
	return name, n
}

func memInfo() (mem, swap int64) {
	b, _ := os.ReadFile("/proc/meminfo")
	for _, l := range strings.Split(string(b), "\n") {
		f := strings.Fields(l)
		if len(f) >= 2 {
			v, _ := strconv.ParseInt(f[1], 10, 64)
			switch f[0] {
			case "MemTotal:":
				mem = v
			case "SwapTotal:":
				swap = v
			}
		}
	}
	return
}

func resolutions() string {
	var out []string
	files, _ := filepath.Glob("/sys/class/drm/card*-*/status")
	for _, f := range files {
		st, _ := os.ReadFile(f)
		if strings.TrimSpace(string(st)) != "connected" {
			continue
		}
		modes, _ := os.ReadFile(filepath.Join(filepath.Dir(f), "modes"))
		m, _, _ := strings.Cut(string(modes), "\n")
		name := filepath.Base(filepath.Dir(f))
		if i := strings.Index(name, "-"); i >= 0 {
			name = name[i+1:]
		}
		out = append(out, name+" "+m)
	}
	return strings.Join(out, "; ")
}

func battery() string {
	files, _ := filepath.Glob("/sys/class/power_supply/BAT*")
	for _, d := range files {
		c, _ := os.ReadFile(filepath.Join(d, "capacity"))
		st, _ := os.ReadFile(filepath.Join(d, "status"))
		return strings.TrimSpace(string(c)) + "% (" + strings.TrimSpace(string(st)) + ")"
	}
	return ""
}

func disks(ctx context.Context, r sys.Runner) []Disk {
	out, err := r.Output(ctx, "lsblk", "-J", "-b", "-o", "NAME,MODEL,SIZE,ROTA,DISC-GRAN,TYPE,MOUNTPOINT,FSAVAIL,FSSIZE")
	if err != nil {
		return nil
	}
	type node struct {
		Name, Model, Type, Mountpoint string
		Size, Rota, Fsavail, Fssize   any
		DiscGran                      any `json:"disc-gran"`
		Children                      []node
	}
	var j struct{ Blockdevices []node }
	if json.Unmarshal([]byte(out), &j) != nil {
		return nil
	}
	var ds []Disk
	var walk func(n node, d *Disk)
	walk = func(n node, d *Disk) {
		if (n.Mountpoint == "/" || n.Mountpoint == "/home") && d.Mount != "/" {
			if fs := num(n.Fssize); fs > 0 {
				d.Mount, d.FreePct, d.FreeHint = n.Mountpoint, int(num(n.Fsavail)*100/fs), human(num(n.Fsavail))
			}
		}
		for _, c := range n.Children {
			walk(c, d)
		}
	}
	for _, b := range j.Blockdevices {
		if b.Type != "disk" || strings.HasPrefix(b.Name, "loop") || strings.HasPrefix(b.Name, "zram") {
			continue
		}
		d := Disk{Name: b.Name, Model: strings.TrimSpace(b.Model), Size: human(num(b.Size)), SSD: !truthy(b.Rota), Discard: num(b.DiscGran) > 0}
		walk(b, &d)
		ds = append(ds, d)
	}
	return ds
}

func num(v any) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	}
	return 0
}

func truthy(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return x == "1" || x == "true"
	case float64:
		return x != 0
	}
	return false
}

func human(b int64) string {
	f := float64(b)
	for _, u := range []string{"B", "KB", "MB", "GB", "TB"} {
		if f < 1024 {
			return fmt.Sprintf("%.1f %s", f, u)
		}
		f /= 1024
	}
	return fmt.Sprintf("%.1f PB", f)
}

func simReport() Report {
	return Report{
		Facts: []Fact{{"CPU", "AMD Ryzen 7 (simulado, 16 threads)"}, {"RAM", "31.3 GiB"}, {"Swap", "8.0 GiB"}, {"zram", "não ativo"},
			{"GPU", "AMD Radeon (simulado)"}, {"Sessão", "wayland"}, {"Ecrãs", "eDP-1 2560x1600"},
			{"Gestor de energia", "power-profiles-daemon (perfil: balanced)"}, {"Bateria", "84% (Discharging)"},
			{"fstrim.timer", "enabled"}, {"Journal (visível ao utilizador)", "Archived and active journals take up 1.4G in the file system."}},
		Disks: []Disk{{Name: "nvme0n1", Model: "Samsung SSD (simulado)", Size: "953.9 GB", SSD: true, Discard: true, Mount: "/", FreePct: 62, FreeHint: "580.2 GB"}},
		Startup: []StartupEntry{
			{Name: "Solaar", Scope: "utilizador", Path: "/sim/config/autostart/solaar.desktop", Status: "available", Detail: "Comando localizado. O arranque efetivo depende da sessão.", Command: "solaar"},
			{Name: "Rastreador de ficheiros", Scope: "sistema", Path: "/etc/xdg/autostart/tracker.desktop", Status: "available", Detail: "Comando localizado. O arranque efetivo depende da sessão.", Command: "tracker"},
		},
		Battery: true,
		Suggestions: []Suggestion{
			{Title: "Portátil detetado", Detail: "Considere o perfil «Portátil equilibrado».", Evidence: "84% (Discharging)", Item: "power-balanced", Level: "suggest"},
			{Title: "Journal grande", Detail: "Limitar a retenção liberta espaço sem desligar os logs.", Evidence: "1.4G", Item: "journald-limit", Level: "suggest"},
		},
	}
}
