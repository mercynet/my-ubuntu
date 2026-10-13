// Package integrations loads optional data packages without granting them execution privileges.
package integrations

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"meuubuntu/internal/safefile"
)

type Check struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Path  string `json:"path,omitempty"`
}
type Guide struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	Command string `json:"command,omitempty"`
}
type Package struct {
	Format      int           `json:"format"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Publisher   string        `json:"publisher"`
	Checks      []Check       `json:"checks"`
	Guide       []Guide       `json:"guide"`
	Setup       *ComposeSetup `json:"setup,omitempty"`
}
type View struct {
	Package
	Project string `json:"project"`
}
type Result struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}
type Report struct {
	CollectedAt string   `json:"collectedAt"`
	Simulated   bool     `json:"simulated"`
	Results     []Result `json:"results"`
}
type Manager struct {
	Data, Config string
	Simulated    bool
}

var validID = regexp.MustCompile(`^[a-z][a-z0-9-]{1,63}$`)

func Decode(data []byte) (Package, error) {
	var p Package
	if len(data) > 256*1024 {
		return p, errors.New("pacote demasiado grande")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		return p, errors.New("manifesto inválido")
	}
	if d.Decode(new(any)) != io.EOF || (p.Format != 1 && p.Format != 2) || !validID.MatchString(p.ID) || p.Name == "" || len(p.Name) > 80 || len(p.Description) > 1000 || len(p.Publisher) > 120 || p.Version == "" || len(p.Version) > 40 || len(p.Checks) > 30 || len(p.Guide) > 30 {
		return p, errors.New("formato de pacote inválido")
	}
	if (p.Format == 2) != (p.Setup != nil) {
		return p, errors.New("preparação automática exige formato 2")
	}
	if p.Setup != nil {
		if err := p.Setup.Validate(); err != nil {
			return p, err
		}
	}
	for _, c := range p.Checks {
		if len(c.Title) > 120 {
			return p, errors.New("título demasiado longo")
		}
		switch c.Type {
		case "docker", "compose", "daemon", "compose-config":
		case "file":
			if c.Path == "" || filepath.IsAbs(c.Path) || strings.Contains(c.Path, "\\") || strings.Contains(c.Path, "..") || filepath.Clean(c.Path) != c.Path {
				return p, errors.New("caminho do pacote inválido")
			}
		default:
			return p, errors.New("verificação não suportada")
		}
	}
	for _, g := range p.Guide {
		if len(g.Title) > 160 || len(g.Text) > 4000 || len(g.Command) > 2000 {
			return p, errors.New("guia demasiado grande")
		}
	}
	return p, nil
}

func ReadArchive(path string) (Package, string, error) {
	var p Package
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > 512*1024 {
		return p, "", errors.New("escolha um pacote ZIP até 512 KiB")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return p, "", err
	}
	z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return p, "", errors.New("ZIP inválido")
	}
	if len(z.File) != 1 || z.File[0].Name != "manifest.json" || !z.File[0].Mode().IsRegular() || z.File[0].UncompressedSize64 > 256*1024 {
		return p, "", errors.New("o ZIP deve conter apenas manifest.json, sem executáveis")
	}
	r, err := z.File[0].Open()
	if err != nil {
		return p, "", err
	}
	defer r.Close()
	b, err := io.ReadAll(io.LimitReader(r, 256*1024+1))
	if err != nil {
		return p, "", err
	}
	p, err = Decode(b)
	sum := sha256.Sum256(data)
	return p, hex.EncodeToString(sum[:]), err
}

func (m Manager) root() string { return filepath.Join(m.Data, "integrations") }
func (m Manager) Load(id string) (View, error) {
	var v View
	if !validID.MatchString(id) {
		return v, errors.New("ID inválido")
	}
	path := filepath.Join(m.root(), id+".json")
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > 256*1024 {
		return v, errors.New("pacote indisponível")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return v, err
	}
	p, err := Decode(b)
	if err != nil {
		return v, err
	}
	if p.ID != id {
		return v, errors.New("ID do pacote divergente")
	}
	v.Package = p
	b, _ = os.ReadFile(filepath.Join(m.Config, "integrations", id+".json"))
	_ = json.Unmarshal(b, &v.Project)
	return v, nil
}
func (m Manager) List() []View {
	out := []View{}
	entries, _ := os.ReadDir(m.root())
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			if v, err := m.Load(strings.TrimSuffix(e.Name(), ".json")); err == nil {
				out = append(out, v)
			}
		}
	}
	return out
}
func secureDir(path string) error {
	if i, err := os.Lstat(path); err == nil && (!i.IsDir() || i.Mode()&os.ModeSymlink != 0) {
		return errors.New("pasta de integração inválida")
	}
	return os.MkdirAll(path, 0700)
}
func (m Manager) Install(p Package) error {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	if _, err = Decode(b); err != nil {
		return err
	}
	if err = secureDir(m.root()); err != nil {
		return err
	}
	return safefile.Write(filepath.Join(m.root(), p.ID+".json"), b, 0600)
}
func (m Manager) SetProject(id, path string) error {
	if _, err := m.Load(id); err != nil {
		return err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	resolved, err = filepath.Abs(resolved)
	if err != nil {
		return err
	}
	i, err := os.Stat(resolved)
	if err != nil || !i.IsDir() || resolved == "/" {
		return errors.New("escolha a pasta do projeto")
	}
	dir := filepath.Join(m.Config, "integrations")
	if err = secureDir(dir); err != nil {
		return err
	}
	b, _ := json.Marshal(resolved)
	return safefile.Write(filepath.Join(dir, id+".json"), b, 0600)
}
func command(ctx context.Context, dir string, args ...string) (string, error) {
	child, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	cmd := exec.CommandContext(child, args[0], args[1:]...)
	cmd.Dir = dir
	var b bytes.Buffer
	cmd.Stdout = &limitedOutput{b: &b}
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return strings.TrimSpace(b.String()), err
}

type limitedOutput struct{ b *bytes.Buffer }

func (w *limitedOutput) Write(p []byte) (int, error) {
	n := len(p)
	remaining := 4096 - w.b.Len()
	if remaining > 0 {
		if len(p) > remaining {
			p = p[:remaining]
		}
		w.b.Write(p)
	}
	return n, nil
}

func (m Manager) Diagnose(ctx context.Context, id string) (Report, error) {
	v, err := m.Load(id)
	rep := Report{CollectedAt: time.Now().Format(time.RFC3339), Simulated: m.Simulated, Results: []Result{}}
	if err != nil {
		return rep, err
	}
	if v.Project == "" {
		return rep, errors.New("escolha primeiro a pasta do projeto")
	}
	for _, c := range v.Checks {
		r := Result{Title: c.Title, Status: "ok", Detail: "Verificado"}
		if m.Simulated {
			r.Status = "example"
			r.Detail = "Exemplo de simulação; o sistema não foi consultado"
			rep.Results = append(rep.Results, r)
			continue
		}
		switch c.Type {
		case "file":
			base, err := filepath.EvalSymlinks(v.Project)
			if err != nil {
				return rep, err
			}
			target, err := filepath.EvalSymlinks(filepath.Join(base, c.Path))
			if err != nil {
				r.Status = "missing"
				r.Detail = "Ficheiro em falta: " + c.Path
			} else {
				rel, er := filepath.Rel(base, target)
				if er != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
					r.Status = "blocked"
					r.Detail = "O ficheiro aponta para fora do projeto"
				} else if i, er := os.Stat(target); er != nil || !i.Mode().IsRegular() {
					r.Status = "missing"
					r.Detail = "Ficheiro regular em falta: " + c.Path
				} else {
					r.Detail = "Presente: " + c.Path
				}
			}
		case "docker":
			s, e := command(ctx, v.Project, "docker", "--version")
			if e != nil {
				r.Status = "missing"
				r.Detail = "Docker indisponível; instale-o no catálogo"
			} else {
				r.Detail = s
			}
		case "compose":
			s, e := command(ctx, v.Project, "docker", "compose", "version", "--short")
			if e != nil {
				r.Status = "missing"
				r.Detail = "Docker Compose v2 indisponível"
			} else {
				r.Detail = "Compose " + s
				var major, minor int
				if _, er := fmt.Sscanf(strings.TrimPrefix(s, "v"), "%d.%d", &major, &minor); er != nil || major < 2 || major == 2 && minor < 20 {
					r.Status = "blocked"
					r.Detail += "; este fluxo requer 2.20 ou superior"
				}
			}
		case "daemon":
			_, e := command(ctx, v.Project, "docker", "info", "--format", "{{.ServerVersion}}")
			if e != nil {
				r.Status = "blocked"
				r.Detail = "Sem acesso ao daemon; confira Docker ativo e permissões da sessão"
			} else {
				r.Detail = "Daemon acessível como utilizador"
			}
		case "compose-config":
			_, e := command(ctx, v.Project, "docker", "compose", "-f", "docker-compose.yml", "config", "--quiet")
			if e != nil {
				r.Status = "blocked"
				r.Detail = "Compose inválido ou dependência local ausente; confira versão, .env e ficheiros de configuração"
			} else {
				r.Detail = "Configuração Compose válida; não comprova que a aplicação está pronta"
			}
		}
		rep.Results = append(rep.Results, r)
	}
	return rep, nil
}
