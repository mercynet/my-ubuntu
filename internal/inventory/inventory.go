// Package inventory classifies package lists exported from an old machine
// (apt-mark showmanual, snap list, flatpak list). The list is a source of
// candidates only: nothing is selected automatically.
package inventory

import (
	"regexp"
	"sort"
	"strings"

	"meuubuntu/internal/catalog"
)

type Entry struct {
	Name   string `json:"name"`
	Source string `json:"source"` // apt | snap | flatpak
	Class  string `json:"class"`  // app | dependency | system | legacy | unknown
	Item   string `json:"item,omitempty"`
	Reason string `json:"reason"`
}

var (
	systemRe  = regexp.MustCompile(`^(linux-(image|headers|modules|generic|hwe|signed|tools|firmware)|grub|shim|libc6|libc-bin|base-files|systemd|ubuntu-(desktop|minimal|standard|server)|initramfs|efibootmgr|mokutil|intel-microcode|amd64-microcode|apt$|dpkg$|snapd$|bash$|sudo$|core\d*$|bare$|gnome-\d+-|gtk-common-themes|snapd-desktop-integration|firmware-updater|snap-store|org\.freedesktop\.Platform|org\.gnome\.Platform|org\.kde\.Platform)`)
	legacyRe  = regexp.MustCompile(`^(libicu\d+|php[57]\.|php8\.[0-2]|python2|openjdk-8|wkhtmltopdf|mainline|ukuu|docker-compose$|.*-dbg$|libssl1\.|erlang)`)
	depRe     = regexp.MustCompile(`^(lib|fonts-|.*-dev$|.*-common$|.*-data$|.*-doc$|gir1\.2-|python3-[a-z0-9]|ruby-|perl-|node-|golang-.*-dev)`)
	nameRe    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9+._\-]{0,127}$`)
	flatpakRe = regexp.MustCompile(`^[a-z][a-z0-9]*\.[A-Za-z0-9_\-]+\.[A-Za-z0-9_.\-]+$`)
	snapLine  = regexp.MustCompile(`\s(latest/)?(stable|edge|beta|candidate)(/|\s|$)|\s\d+/stable`)
)

// Parse detects the format line by line and classifies each entry.
func Parse(text string, cat *catalog.Catalog) []Entry {
	byPkg := map[string]string{}
	byFlat := map[string]string{}
	bySnap := map[string]string{}
	for _, it := range cat.Items {
		s, _ := cat.ValidateSelection(it.ID, catalog.Selection{})
		for _, p := range cat.ExpandPackages(&it, s) {
			if _, ok := byPkg[p]; !ok && !depRe.MatchString(p) { // libraries are never "apps"
				byPkg[p] = it.ID
			}
		}
		if it.Flatpak != "" {
			byFlat[it.Flatpak] = it.ID
		}
		for _, t := range it.Tags {
			if strings.HasPrefix(t, "snap:") {
				bySnap[strings.TrimPrefix(t, "snap:")] = it.ID
			}
		}
	}
	seen := map[string]bool{}
	var out []Entry
	for _, raw := range strings.Split(text, "\n") {
		l := strings.TrimSpace(raw)
		if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, "Name ") || strings.HasPrefix(l, "Application ID") || strings.HasPrefix(l, "Listing") {
			continue
		}
		f := strings.Fields(l)
		name, src := f[0], "apt"
		if strings.Contains(name, "/") { // "apt list --installed" format
			name = strings.SplitN(name, "/", 2)[0]
		}
		if len(f) >= 4 && snapLine.MatchString(l) {
			src = "snap"
		}
		for _, x := range f {
			if flatpakRe.MatchString(x) {
				name, src = x, "flatpak"
				break
			}
		}
		name = strings.SplitN(name, ":", 2)[0]
		if !nameRe.MatchString(name) || seen[src+name] {
			continue
		}
		seen[src+name] = true
		e := Entry{Name: name, Source: src}
		switch {
		case src == "flatpak" && byFlat[name] != "":
			e.Class, e.Item, e.Reason = "app", byFlat[name], "Existe no catálogo"
		case src == "snap" && bySnap[name] != "":
			e.Class, e.Item, e.Reason = "app", bySnap[name], "Existe no catálogo (método recomendado pode diferir)"
		case src == "apt" && byPkg[name] != "":
			e.Class, e.Item, e.Reason = "app", byPkg[name], "Existe no catálogo"
		case systemRe.MatchString(name):
			e.Class, e.Reason = "system", "Componente do sistema/arranque — não transportar"
		case legacyRe.MatchString(name):
			e.Class, e.Reason = "legacy", "Versão antiga ou substituída — ver opções avançadas"
		case depRe.MatchString(name):
			e.Class, e.Reason = "dependency", "Biblioteca/dependência — instalada automaticamente quando necessária"
		default:
			e.Class, e.Reason = "unknown", "Não reconhecido — fica desmarcado"
		}
		out = append(out, e)
	}
	order := map[string]int{"app": 0, "unknown": 1, "legacy": 2, "dependency": 3, "system": 4}
	sort.SliceStable(out, func(i, j int) bool {
		if order[out[i].Class] != order[out[j].Class] {
			return order[out[i].Class] < order[out[j].Class]
		}
		return out[i].Name < out[j].Name
	})
	return out
}
