// Package aptsrc inspects APT sources without changing them. A source is only
// redundant when an active source covers its types, suites and architectures.
package aptsrc

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"meuubuntu/internal/catalog"
)

const Marker = "# Gerido por Meu Ubuntu"

type Duplicate struct {
	Repo    string   `json:"repo"`
	Name    string   `json:"name"`
	Ours    string   `json:"ours"`
	Keyring string   `json:"keyring,omitempty"` // empty if another source needs it
	Others  []string `json:"others"`
}

func OursPath(id string) string    { return "/etc/apt/sources.list.d/meu-ubuntu-" + id + ".sources" }
func KeyringPath(id string) string { return "/etc/apt/keyrings/meu-ubuntu-" + id + ".gpg" }
func norm(u string) string         { return strings.TrimSuffix(strings.TrimSpace(u), "/") }

type source struct {
	types, uris, suites, components, architectures []string
	active, supported                              bool
}

func files(root string) []string {
	out := []string{filepath.Join(root, "/etc/apt/sources.list")}
	for _, pat := range []string{"*.list", "*.sources"} {
		m, _ := filepath.Glob(filepath.Join(root, "/etc/apt/sources.list.d", pat))
		out = append(out, m...)
	}
	sort.Strings(out)
	return out
}

func parse(data []byte, deb822 bool) []source {
	var out []source
	if !deb822 {
		for _, line := range strings.Split(string(data), "\n") {
			line, _, _ = strings.Cut(line, "#")
			f := strings.Fields(line)
			if len(f) < 3 || (f[0] != "deb" && f[0] != "deb-src") {
				continue
			}
			s := source{types: f[:1], active: true, supported: true}
			f = f[1:]
			if strings.HasPrefix(f[0], "[") {
				var opts []string
				for len(f) > 0 {
					o := f[0]
					f = f[1:]
					opts = append(opts, strings.Trim(o, "[]"))
					if strings.HasSuffix(o, "]") {
						break
					}
				}
				for _, opt := range opts {
					k, v, _ := strings.Cut(opt, "=")
					switch k {
					case "arch":
						s.architectures = strings.Split(v, ",")
					case "signed-by":
					default:
						s.supported = false
					}
				}
			}
			if len(f) < 2 {
				continue
			}
			s.uris, s.suites, s.components = f[:1], f[1:2], f[2:]
			out = append(out, s)
		}
		return out
	}
	fields := map[string]string{}
	last, valid := "", true
	flush := func() {
		if len(fields) == 0 {
			return
		}
		s := source{
			types: strings.Fields(fields["types"]), uris: strings.Fields(fields["uris"]),
			suites: strings.Fields(fields["suites"]), components: strings.Fields(fields["components"]),
			architectures: strings.Fields(fields["architectures"]),
			active:        fields["enabled"] == "" || strings.EqualFold(fields["enabled"], "yes"), supported: valid,
		}
		for k := range fields {
			switch k {
			case "types", "uris", "suites", "components", "architectures", "signed-by", "enabled":
			default:
				if !strings.HasPrefix(k, "x-") {
					s.supported = false
				}
			}
		}
		if len(s.types) > 0 && len(s.uris) > 0 && len(s.suites) > 0 {
			out = append(out, s)
		}
		fields, last, valid = map[string]string{}, "", true
	}
	for _, line := range strings.Split(string(data), "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" {
			flush()
			continue
		}
		if strings.HasPrefix(trim, "#") {
			continue
		}
		if (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) && last != "" {
			fields[last] += " " + trim
			continue
		}
		k, v, ok := strings.Cut(trim, ":")
		if !ok {
			valid = false
			continue
		}
		last = strings.ToLower(k)
		if _, duplicate := fields[last]; duplicate {
			valid = false
		}
		fields[last] = strings.TrimSpace(v)
	}
	flush()
	return out
}

func covers(have, want []string, normalize bool) bool {
	for _, w := range want {
		found := false
		for _, h := range have {
			if h == w || (normalize && norm(h) == norm(w)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// ReplacementFiles only returns replacements when every desired entry is
// covered. Unknown restrictions prevent automatic cleanup.
func ReplacementFiles(root string, desired []byte, skip string) []string {
	wanted := parse(desired, true)
	if len(wanted) == 0 {
		return nil
	}
	covered := make([]bool, len(wanted))
	var out []string
	for _, f := range files(root) {
		if f == skip {
			continue
		}
		if fi, err := os.Lstat(f); err != nil || !fi.Mode().IsRegular() {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		matched := false
		for _, s := range parse(b, filepath.Ext(f) == ".sources") {
			if !s.active || !s.supported {
				continue
			}
			for i, w := range wanted {
				if !w.active || !w.supported {
					return nil
				}
				arch := len(s.architectures) == 0 || (len(w.architectures) > 0 && covers(s.architectures, w.architectures, false))
				if arch && covers(s.types, w.types, false) && covers(s.uris, w.uris, true) &&
					covers(s.suites, w.suites, false) && covers(s.components, w.components, false) {
					covered[i], matched = true, true
				}
			}
		}
		if matched {
			out = append(out, strings.TrimPrefix(f, root))
		}
	}
	for _, c := range covered {
		if !c {
			return nil
		}
	}
	return out
}

// KeyringInUse includes disabled sources, which may be enabled later. Any
// unreadable source makes deletion unsafe, so the keyring is retained.
func KeyringInUse(root, keyring, skip string) bool {
	for _, f := range files(root) {
		if f == skip {
			continue
		}
		b, err := os.ReadFile(f)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil || bytes.Contains(b, []byte(keyring)) {
			return true
		}
	}
	return false
}

func Duplicates(root string, cat *catalog.Catalog) []Duplicate {
	var out []Duplicate
	for _, r := range cat.Repos {
		ours := filepath.Join(root, OursPath(r.ID))
		if fi, err := os.Lstat(ours); err != nil || !fi.Mode().IsRegular() {
			continue
		}
		b, err := os.ReadFile(ours)
		if err != nil || !bytes.HasPrefix(b, []byte(Marker)) {
			continue
		}
		if others := ReplacementFiles(root, b, ours); len(others) > 0 {
			key := KeyringPath(r.ID)
			if KeyringInUse(root, key, ours) {
				key = ""
			}
			out = append(out, Duplicate{Repo: r.ID, Name: r.Name, Ours: OursPath(r.ID), Keyring: key, Others: others})
		}
	}
	return out
}
