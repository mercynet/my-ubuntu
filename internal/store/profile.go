package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"meuubuntu/internal/catalog"
)

// Profile format. v1 (0.1/0.2) is still accepted and migrated explicitly to v2.
const (
	ProfileSchema  = "meu-ubuntu/profile"
	ProfileVersion = 2
)

// Where a profile comes from. Only informative — never a proof of trust.
const (
	SourceBuiltin  = "builtin"
	SourceMine     = "mine"
	SourceImported = "imported"
)

var ErrExists = errors.New("já existe um perfil com este identificador")

type Requirements struct {
	MinRelease string   `json:"minRelease,omitempty"` // e.g. "26.04"
	Arch       []string `json:"arch,omitempty"`
}

// Profile holds only item IDs and declarative, catalog-validated options.
type Profile struct {
	Schema      string                       `json:"schema"`
	Version     int                          `json:"version"`
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description,omitempty"`
	Audience    string                       `json:"audience,omitempty"` // "para quem serve"
	Author      string                       `json:"author,omitempty"`   // optional, not verified
	Kind        string                       `json:"kind,omitempty"`     // usage | complement
	Requires    Requirements                 `json:"requires,omitempty"`
	Created     string                       `json:"created,omitempty"`
	Updated     string                       `json:"updated,omitempty"`
	Selections  map[string]catalog.Selection `json:"selections"`
	// Set by the loader, ignored on import.
	Source  string `json:"source,omitempty"`
	Builtin bool   `json:"builtin,omitempty"`
}

// v1 is the exact 0.1/0.2 format, decoded strictly.
type profileV1 struct {
	Schema      string                       `json:"schema"`
	Version     int                          `json:"version"`
	Name        string                       `json:"name"`
	Description string                       `json:"description,omitempty"`
	Builtin     bool                         `json:"builtin,omitempty"`
	Created     string                       `json:"created,omitempty"`
	Selections  map[string]catalog.Selection `json:"selections"`
}

var (
	idRe      = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,63}$`)
	releaseRe = regexp.MustCompile(`^\d{2}\.\d{2}$`)
	nameRe    = regexp.MustCompile(`[^a-z0-9-]+`)
	// Option keys/values that look like credentials are refused.
	secretRe = regexp.MustCompile(`(?i)(pass(word)?|senha|token|secret|private|apikey|api_key)`)
	// Text fields must not carry key material or well-known token formats.
	keyMaterialRe = regexp.MustCompile(`-----BEGIN [A-Z ]+-----|\bghp_[A-Za-z0-9]{20,}|\bsk_(live|test)_[A-Za-z0-9]{10,}|\bAKIA[0-9A-Z]{16}\b|\bxox[bp]-[A-Za-z0-9-]{10,}`)
)

func Slug(n string) string {
	s := strings.Trim(nameRe.ReplaceAllString(strings.ToLower(n), "-"), "-")
	if s == "" || !idRe.MatchString(s) {
		s = "perfil-" + s
	}
	if len(s) > 60 {
		s = strings.Trim(s[:60], "-")
	}
	return s
}

func strictDecode(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // e.g. "command", "script", "url" are rejected
	if err := dec.Decode(v); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("conteúdo extra depois do perfil")
	}
	return nil
}

// ParseProfile decodes a profile file strictly (unknown fields are an error),
// migrates v1 to v2 and validates it against the catalog.
func ParseProfile(cat *catalog.Catalog, data []byte) (*Profile, []string, error) {
	if len(data) > 1<<20 {
		return nil, nil, errors.New("ficheiro demasiado grande para um perfil")
	}
	var head struct {
		Schema  string `json:"schema"`
		Version int    `json:"version"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, nil, errors.New("JSON inválido: " + err.Error())
	}
	if head.Schema != ProfileSchema {
		return nil, nil, fmt.Errorf("ficheiro não é um perfil Meu Ubuntu (schema %q)", head.Schema)
	}
	var p Profile
	var warnings []string
	switch head.Version {
	case 1:
		var v1 profileV1
		if err := strictDecode(data, &v1); err != nil {
			return nil, nil, fmt.Errorf("perfil v1 com campos não permitidos: %w", err)
		}
		p = Profile{Schema: ProfileSchema, Version: ProfileVersion, ID: Slug(v1.Name), Name: v1.Name, Description: v1.Description,
			Created: v1.Created, Selections: v1.Selections}
		warnings = append(warnings, "Perfil no formato v1 migrado para v2 (identificador gerado a partir do nome)")
	case 2:
		if err := strictDecode(data, &p); err != nil {
			return nil, nil, fmt.Errorf("perfil com campos não permitidos: %w", err)
		}
	default:
		return nil, nil, fmt.Errorf("versão de perfil %d não suportada (suportadas: 1 e %d)", head.Version, ProfileVersion)
	}
	w, err := ValidateProfile(cat, &p)
	return &p, append(warnings, w...), err
}

// ValidateProfile checks metadata and every selection against the catalog.
// Unknown items are dropped and reported; invalid options are fatal.
func ValidateProfile(cat *catalog.Catalog, p *Profile) (warnings []string, err error) {
	if p.Schema != ProfileSchema || p.Version != ProfileVersion {
		return nil, fmt.Errorf("formato de perfil inválido (%s v%d)", p.Schema, p.Version)
	}
	if p.ID == "" {
		p.ID = Slug(p.Name)
	}
	switch {
	case !idRe.MatchString(p.ID):
		return nil, fmt.Errorf("identificador de perfil inválido %q", p.ID)
	case strings.TrimSpace(p.Name) == "" || len(p.Name) > 80:
		return nil, errors.New("nome do perfil inválido (1–80 caracteres)")
	case len(p.Description) > 600 || len(p.Audience) > 300 || len(p.Author) > 80:
		return nil, errors.New("descrição, público ou autoria demasiado longos")
	case p.Kind != "" && p.Kind != "usage" && p.Kind != "complement":
		return nil, fmt.Errorf("tipo de perfil inválido %q", p.Kind)
	case p.Requires.MinRelease != "" && !releaseRe.MatchString(p.Requires.MinRelease):
		return nil, fmt.Errorf("requisito de versão inválido %q", p.Requires.MinRelease)
	}
	for _, a := range p.Requires.Arch {
		if a != "amd64" && a != "arm64" {
			return nil, fmt.Errorf("arquitetura inválida %q", a)
		}
	}
	for _, t := range []string{p.Name, p.Description, p.Audience, p.Author} {
		if keyMaterialRe.MatchString(t) {
			return nil, errors.New("o perfil parece conter chaves, certificados ou tokens; recusado")
		}
	}
	clean := map[string]catalog.Selection{}
	ids := make([]string, 0, len(p.Selections))
	for id := range p.Selections {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		sel := p.Selections[id]
		if cat.Item(id) == nil {
			warnings = append(warnings, "Item desconhecido ignorado: "+id)
			continue
		}
		for k, vs := range sel.Options {
			for _, v := range vs {
				if secretRe.MatchString(k) || secretRe.MatchString(v) || keyMaterialRe.MatchString(v) {
					return nil, fmt.Errorf("%s: o perfil parece conter credenciais; recusado", id)
				}
			}
		}
		v, err := cat.ValidateSelection(id, sel)
		if err != nil {
			return nil, err
		}
		clean[id] = v
	}
	p.Selections = clean
	return warnings, nil
}

func (s *Store) dir(source string) string {
	if source == SourceImported {
		return filepath.Join(s.Config, "profiles", "imported")
	}
	return filepath.Join(s.Config, "profiles")
}

// SaveProfile writes a personal or imported profile. An existing profile with
// the same ID is never replaced unless overwrite is set.
func (s *Store) SaveProfile(p Profile, source string, overwrite bool) error {
	if source != SourceMine && source != SourceImported {
		return errors.New("destino de perfil inválido")
	}
	if !idRe.MatchString(p.ID) {
		return fmt.Errorf("identificador inválido %q", p.ID)
	}
	path := filepath.Join(s.dir(source), p.ID+".json")
	if _, err := os.Stat(path); err == nil && !overwrite {
		return ErrExists
	}
	p.Schema, p.Version, p.Source, p.Builtin = ProfileSchema, ProfileVersion, "", false
	now := time.Now().Format(time.RFC3339)
	if p.Created == "" {
		p.Created = now
	}
	p.Updated = now
	return writeJSON(path, p)
}

func (s *Store) DeleteProfile(source, id string) error {
	if !idRe.MatchString(id) || (source != SourceMine && source != SourceImported) {
		return errors.New("perfil inválido")
	}
	return os.Remove(filepath.Join(s.dir(source), id+".json"))
}

// ListProfiles returns personal and imported profiles (v1 files are migrated in memory).
func (s *Store) ListProfiles(cat *catalog.Catalog) []Profile {
	var out []Profile
	for _, src := range []string{SourceMine, SourceImported} {
		files, _ := filepath.Glob(filepath.Join(s.dir(src), "*.json"))
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				continue
			}
			p, _, err := ParseProfile(cat, b)
			if err != nil {
				continue
			}
			if src == SourceMine && strings.TrimSuffix(filepath.Base(f), ".json") != p.ID {
				p.ID = strings.TrimSuffix(filepath.Base(f), ".json") // v1 files keep their file name as ID
			}
			p.Source = src
			out = append(out, *p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// MarshalExport produces the shareable JSON (no loader-only fields).
func MarshalExport(p Profile) ([]byte, error) {
	p.Schema, p.Version, p.Source, p.Builtin = ProfileSchema, ProfileVersion, "", false
	return json.MarshalIndent(p, "", "  ")
}
