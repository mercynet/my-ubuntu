// Package catalog holds the declarative, embedded catalog of applications,
// settings and recipes. Adding a common item only requires editing
// catalog.json; the UI renders whatever is declared here.
package catalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

//go:embed catalog.json
var embedded []byte

// Methods understood by the planner/executor.
const (
	MethodAPT          = "apt"           // packages from the Ubuntu archive
	MethodAPTRepo      = "apt-repo"      // packages from an official third-party repo (Repo)
	MethodFlatpak      = "flatpak"       // flatpak --user from Flathub
	MethodGSettings    = "gsettings"     // user preferences (dconf)
	MethodGitConfig    = "gitconfig"     // git config --global
	MethodTemplate     = "template"      // file in ~/Templates
	MethodExtension    = "extension"     // GNOME extension via extensions.gnome.org API
	MethodPowerProfile = "power-profile" // powerprofilesctl set
	MethodAdmin        = "admin"         // allowlisted root operation (AdminOp)
	MethodManual       = "manual"        // documented manual steps only
	MethodUserScript   = "user-script"   // vendor installer run as the user after inspection (Script)
	MethodNPMUser      = "npm-user"      // pinned official npm package, installed in the user's home
)

var idRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,63}$`)
var pkgRe = regexp.MustCompile(`^[a-z0-9][a-z0-9+.\-]{0,127}$`)
var optValRe = regexp.MustCompile(`^[A-Za-z0-9 ._:'@,+\-/\[\]()]{0,200}$`)

// Repo is an official third-party APT repository.
type Repo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	KeyURL       string   `json:"keyUrl"`
	Fingerprints []string `json:"fingerprints"` // allowlist of 40-hex primary key fingerprints; every primary key in the file must be listed
	URL          string   `json:"url"`
	Suite        string   `json:"suite"` // "{codename}" is replaced by the running release codename
	Components   []string `json:"components"`
	Arch         string   `json:"arch"`
	Docs         string   `json:"docs"`
}

type OptionValue struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Option struct {
	Key     string        `json:"key"`
	Label   string        `json:"label"`
	Multi   bool          `json:"multi,omitempty"`
	Values  []OptionValue `json:"values"`
	Default []string      `json:"default"`
	Help    string        `json:"help,omitempty"`
}

type Conflict struct {
	ID     string `json:"id"`
	Hard   bool   `json:"hard"`
	Reason string `json:"reason"`
}

// Setting is a single gsettings key. Value is a GVariant literal, or taken
// from the item option named in Option.
//
// Schema may contain a placeholder resolved at runtime from the user's
// configuration: {ptyxis-profile} or {gnome-terminal-profile} (relocatable
// profile schemas). With Map, the chosen option value is translated
// (e.g. one wallpaper choice -> light and dark URIs).
type Setting struct {
	Schema string            `json:"schema"`
	Key    string            `json:"key"`
	Value  string            `json:"value,omitempty"`
	Option string            `json:"option,omitempty"`
	Map    map[string]string `json:"map,omitempty"`
}

// Group describes mutually exclusive alternatives shown as one card with a select.
type Group struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Area    string `json:"area"`
	Summary string `json:"summary"`
	Why     string `json:"why"`
}

type GitSetting struct {
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	Option string `json:"option,omitempty"`
}

type Detect struct {
	Commands []string `json:"commands,omitempty"` // any of these on PATH => present (manual items)
	Files    []string `json:"files,omitempty"`    // any of these exist (~ expanded)
	Packages []string `json:"packages,omitempty"` // override packages to check
}

type Verify struct {
	Cmd []string `json:"cmd,omitempty"` // fixed argv from the embedded catalog only
}

type Compat struct {
	Arch       []string `json:"arch,omitempty"`
	MinRelease string   `json:"minRelease,omitempty"`
	Notes      string   `json:"notes,omitempty"`
}

type Manual struct {
	Steps []string `json:"steps"`
	URL   string   `json:"url,omitempty"`
}

type Item struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Area       string       `json:"area"` // apps, desktop, extensions, dev, energy
	Category   string       `json:"category"`
	Summary    string       `json:"summary"`
	Why        string       `json:"why"`
	Notes      string       `json:"notes,omitempty"`
	Source     string       `json:"source"`
	Method     string       `json:"method"`
	Packages   []string     `json:"packages,omitempty"` // may contain {opt:key}
	Repo       string       `json:"repo,omitempty"`
	Flatpak    string       `json:"flatpak,omitempty"`
	Extension  string       `json:"extension,omitempty"` // uuid
	ProvidedBy []string     `json:"providedBy,omitempty"`
	Template   *Template    `json:"template,omitempty"`
	Script     *Script      `json:"script,omitempty"`
	NPM        *NPM         `json:"npm,omitempty"`
	AdminOp    string       `json:"adminOp,omitempty"`
	Profile    string       `json:"powerProfile,omitempty"`
	Settings   []Setting    `json:"settings,omitempty"`
	Git        []GitSetting `json:"git,omitempty"`
	Requires   []string     `json:"requires,omitempty"`
	Conflicts  []Conflict   `json:"conflicts,omitempty"`
	Group      string       `json:"group,omitempty"`     // mutually exclusive alternatives
	AddonsFor  []string     `json:"addonsFor,omitempty"` // tools extended by this item; not exclusive
	Compat     Compat       `json:"compat"`
	Detect     Detect       `json:"detect"`
	Verify     Verify       `json:"verify"`
	Options    []Option     `json:"options,omitempty"`
	Advanced   bool         `json:"advanced,omitempty"`
	Service    bool         `json:"service,omitempty"`
	Relogin    bool         `json:"relogin,omitempty"`
	Reboot     bool         `json:"reboot,omitempty"`
	Manual     *Manual      `json:"manual,omitempty"`
	Post       []string     `json:"post,omitempty"` // allowlisted admin post-ops, e.g. "xdebug-off"
	Tags       []string     `json:"tags,omitempty"`
	Verified   string       `json:"verified,omitempty"` // where the recipe was verified
	Confirm    string       `json:"confirm,omitempty"`  // explicit confirmation required in the plan
	Sensitive  string       `json:"sensitive,omitempty"`
}

// SensitiveReason explains why an item deserves attention in imported
// profiles and plans (root-equivalent access, open ports, third-party code...).
func (c *Catalog) SensitiveReason(it *Item) string {
	switch {
	case it.Sensitive != "":
		return it.Sensitive
	case it.Method == MethodUserScript:
		return "Executa um instalador do fornecedor (confiança no fornecedor)"
	case it.Method == MethodNPMUser:
		return "Instala código do fornecedor como o seu utilizador"
	case it.Method == MethodAPTRepo:
		return "Adiciona um repositório de terceiros ao APT"
	case it.Service:
		return "Instala um serviço local que pode arrancar com o sistema"
	case it.Method == MethodAdmin:
		return "Operação administrativa no sistema"
	}
	return ""
}

// Script is a vendor installer that runs WITHOUT root. It is downloaded to a
// file, hashed, checked (no sudo, only allowlisted hosts) and then executed
// with bash as the user — never piped from curl.
type Script struct {
	URL         string   `json:"url"` // may contain {opt:key}
	AllowHosts  []string `json:"allowHosts"`
	Backup      []string `json:"backup"`      // files (~/...) copied before running
	Binary      string   `json:"binary"`      // ~/... path proving installation
	VersionArgs []string `json:"versionArgs"` // argv appended to Binary that prints "X.Y"
	Option      string   `json:"option"`      // option compared with the printed version
}

// NPM metadata belongs to the embedded catalog; profiles cannot supply it.
type NPM struct {
	Package       string `json:"package"`
	Version       string `json:"version"`
	Binary        string `json:"binary"`
	IgnoreScripts bool   `json:"ignoreScripts,omitempty"`
}

var npmPackageRe = regexp.MustCompile(`^(?:@[a-z0-9][a-z0-9._-]*/)?[a-z0-9][a-z0-9._-]*$`)
var npmVersionRe = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
var npmBinaryRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func (n *NPM) Valid() bool {
	return n != nil && npmPackageRe.MatchString(n.Package) && npmVersionRe.MatchString(n.Version) && npmBinaryRe.MatchString(n.Binary)
}

type Template struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Catalog struct {
	Version  int      `json:"version"`
	Repos    []Repo   `json:"repos"`
	Groups   []Group  `json:"groups"`
	Items    []Item   `json:"items"`
	Profiles []string `json:"-"`

	byID   map[string]*Item
	repoBy map[string]*Repo
}

func Load() (*Catalog, error) { return Parse(embedded) }

func Parse(data []byte) (*Catalog, error) {
	var c Catalog
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("catálogo inválido: %w", err)
	}
	return &c, c.index()
}

func (c *Catalog) index() error {
	c.byID = map[string]*Item{}
	c.repoBy = map[string]*Repo{}
	for i := range c.Repos {
		r := &c.Repos[i]
		if !idRe.MatchString(r.ID) {
			return fmt.Errorf("repo id inválido %q", r.ID)
		}
		c.repoBy[r.ID] = r
	}
	for i := range c.Items {
		it := &c.Items[i]
		if !idRe.MatchString(it.ID) {
			return fmt.Errorf("id inválido %q", it.ID)
		}
		if _, dup := c.byID[it.ID]; dup {
			return fmt.Errorf("id duplicado %q", it.ID)
		}
		c.byID[it.ID] = it
		if it.Method == MethodNPMUser && !it.NPM.Valid() {
			return fmt.Errorf("%s: pacote npm, versão fixa ou comando inválido", it.ID)
		}
		for _, p := range it.Packages {
			if !strings.Contains(p, "{opt:") && !pkgRe.MatchString(p) {
				return fmt.Errorf("%s: pacote inválido %q", it.ID, p)
			}
		}
		if it.Method == MethodAPTRepo && c.repoBy[it.Repo] == nil {
			return fmt.Errorf("%s: repo desconhecido %q", it.ID, it.Repo)
		}
		for _, o := range it.Options {
			for _, v := range o.Values {
				if !optValRe.MatchString(v.Value) {
					return fmt.Errorf("%s: valor de opção inválido %q", it.ID, v.Value)
				}
			}
		}
	}
	for _, it := range c.Items {
		for _, target := range it.AddonsFor {
			if target == it.ID || c.byID[target] == nil {
				return fmt.Errorf("%s: ferramenta de complemento inválida %q", it.ID, target)
			}
		}
		for _, r := range it.Requires {
			if c.byID[r] == nil {
				return fmt.Errorf("%s requer item desconhecido %q", it.ID, r)
			}
		}
		for _, cf := range it.Conflicts {
			if c.byID[cf.ID] == nil {
				return fmt.Errorf("%s conflita com item desconhecido %q", it.ID, cf.ID)
			}
		}
	}
	return nil
}

func (c *Catalog) Item(id string) *Item { return c.byID[id] }
func (c *Catalog) Repo(id string) *Repo { return c.repoBy[id] }

// Selection is the user's choice for one item: option key -> chosen values.
type Selection struct {
	Options map[string][]string `json:"options,omitempty"`
}

// ValidateSelection checks option values against the catalog and fills defaults.
func (c *Catalog) ValidateSelection(id string, s Selection) (Selection, error) {
	it := c.Item(id)
	if it == nil {
		return s, fmt.Errorf("item desconhecido %q", id)
	}
	out := Selection{Options: map[string][]string{}}
	for _, o := range it.Options {
		vals, ok := s.Options[o.Key]
		if !ok || len(vals) == 0 {
			vals = o.Default
		}
		if !o.Multi && len(vals) > 1 {
			return s, fmt.Errorf("%s: opção %s aceita um único valor", id, o.Key)
		}
		allowed := map[string]bool{}
		for _, v := range o.Values {
			allowed[v.Value] = true
		}
		for _, v := range vals {
			if !allowed[v] {
				return s, fmt.Errorf("%s: valor %q não permitido em %s", id, v, o.Key)
			}
		}
		out.Options[o.Key] = append([]string(nil), vals...)
	}
	for k := range s.Options {
		found := false
		for _, o := range it.Options {
			found = found || o.Key == k
		}
		if !found {
			return s, fmt.Errorf("%s: opção desconhecida %q", id, k)
		}
	}
	return out, nil
}

// ExpandPackages resolves {opt:key} placeholders using validated options.
func (c *Catalog) ExpandPackages(it *Item, s Selection) []string {
	var out []string
	seen := map[string]bool{}
	add := func(p string) {
		if pkgRe.MatchString(p) && !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	for _, p := range it.Packages {
		i := strings.Index(p, "{opt:")
		if i < 0 {
			add(p)
			continue
		}
		j := strings.Index(p[i:], "}")
		key := p[i+5 : i+j]
		for _, v := range s.Options[key] {
			add(p[:i] + v + p[i+j+1:])
		}
	}
	return out
}

// OptionValue returns the first chosen value for key.
func (s Selection) First(key string) string {
	if v := s.Options[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func (c *Catalog) Areas() []string {
	m := map[string]bool{}
	for _, it := range c.Items {
		m[it.Area] = true
	}
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

var Fingerprint = regexp.MustCompile(`^[0-9A-F]{40}$`)
var PackageName = pkgRe

// Verified reports whether the repo has a usable fingerprint allowlist.
func (r *Repo) Verified() bool {
	if len(r.Fingerprints) == 0 {
		return false
	}
	for _, f := range r.Fingerprints {
		if !Fingerprint.MatchString(f) {
			return false
		}
	}
	return true
}
