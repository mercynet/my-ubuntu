// Package helper implements the restricted administrative mode. The GUI runs
// as the normal user and starts `pkexec <self> --helper` once per execution.
// The helper reads JSON requests on stdin and only performs allowlisted
// operations. Package names, repositories and file contents are resolved from
// the embedded catalog — never taken verbatim from the request — so no
// arbitrary command can be smuggled in.
package helper

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"meuubuntu/internal/adminconfig"
	"meuubuntu/internal/aptsrc"
	"meuubuntu/internal/catalog"
	"meuubuntu/internal/safefile"
	"meuubuntu/internal/sys"
)

type Request struct {
	Op                 string              `json:"op"`
	Item               string              `json:"item,omitempty"`
	Options            map[string][]string `json:"options,omitempty"`
	Batch              bool                `json:"batch,omitempty"`
	ExpectedStateToken string              `json:"expectedStateToken,omitempty"`
}

type Message struct {
	Type    string         `json:"type"` // log | done
	Line    string         `json:"line,omitempty"`
	OK      bool           `json:"ok,omitempty"`
	Err     string         `json:"err,omitempty"`
	Preview *RevertPreview `json:"preview,omitempty"`
}

// Env abstracts the root environment so tests can run it unprivileged.
type Env struct {
	R            sys.Runner
	Cat          *catalog.Catalog
	Root         string // filesystem prefix ("" in production)
	Codename     string
	Distribution string // read from os-release by the helper, never from the request
	UID          string // the invoking user (from PKEXEC_UID), never from the request
	Fetch        func(url string) ([]byte, error)
}

const aptLockTimeout = "DPkg::Lock::Timeout=600" // wait for other APT users; never delete locks

var aptOpts = []string{"-y", "-o", aptLockTimeout, "-o", "Dpkg::Options::=--force-confdef", "-o", "Dpkg::Options::=--force-confold"}

func (e *Env) path(p string) string { return filepath.Join(e.Root, p) }

// Execute validates and runs one request.
func (e *Env) Execute(ctx context.Context, req Request, log func(string)) error {
	if req.Op == "revert" && req.ExpectedStateToken != "" {
		p, err := e.PreviewRevert(ctx, req)
		if err != nil {
			return err
		}
		if p.StateToken != req.ExpectedStateToken {
			return errors.New("a configuração administrativa mudou desde a revisão; nada revertido")
		}
	}
	item := e.Cat.Item(req.Item)
	apt := req.Op == "apt-update" || req.Op == "apt-install" || req.Op == "repo-add" || req.Op == "sources-check" || req.Op == "post" || (item != nil && item.AdminOp == "ubuntu-drivers")
	if apt && (e.Distribution != "" && e.Distribution != "ubuntu" || e.Root == "" && e.Distribution == "") {
		return errors.New("receitas administrativas Ubuntu não são suportadas nesta distribuição")
	}
	switch req.Op {
	case "apt-update":
		return e.R.Stream(ctx, log, "apt-get", append(aptOpts[:3:3], "update")...)
	case "apt-install":
		if item == nil || (item.Method != catalog.MethodAPT && item.Method != catalog.MethodAPTRepo) {
			return fmt.Errorf("item inválido para apt-install: %q", req.Item)
		}
		sel, err := e.Cat.ValidateSelection(item.ID, catalog.Selection{Options: req.Options})
		if err != nil {
			return err
		}
		pkgs := e.Cat.ExpandPackages(item, sel)
		if len(pkgs) == 0 {
			return errors.New("nenhum pacote")
		}
		args := append(append([]string{}, aptOpts...), "install", "--no-remove")
		args = append(args, pkgs...)
		if err := e.R.Stream(ctx, log, "apt-get", args...); err != nil {
			return err
		}
		if item.Method == catalog.MethodAPTRepo {
			// Some vendor packages (e.g. warp-terminal) add their own source for the
			// same URI on install; keeping ours too breaks APT ("Conflicting values").
			return e.dedupeSources(ctx, log)
		}
		return nil
	case "sources-check":
		return e.dedupeSources(ctx, log)
	case "repo-add":
		if item == nil || item.Method != catalog.MethodAPTRepo {
			return fmt.Errorf("item sem repositório: %q", req.Item)
		}
		return e.addRepo(ctx, e.Cat.Repo(item.Repo), log)
	case "post":
		return e.post(ctx, item, req, log)
	case "admin":
		if item == nil || item.Method != catalog.MethodAdmin {
			return fmt.Errorf("item inválido: %q", req.Item)
		}
		sel, err := e.Cat.ValidateSelection(item.ID, catalog.Selection{Options: req.Options})
		if err != nil {
			return err
		}
		return e.admin(ctx, item, sel, log)
	case "revert":
		if item == nil || item.Method != catalog.MethodAdmin {
			return fmt.Errorf("item inválido: %q", req.Item)
		}
		sel, err := e.Cat.ValidateSelection(item.ID, catalog.Selection{Options: req.Options})
		if err != nil {
			return err
		}
		return e.revert(ctx, item, sel, log)
	}
	return fmt.Errorf("operação não permitida: %q", req.Op)
}

func (e *Env) username() (string, error) {
	if !regexp.MustCompile(`^[0-9]+$`).MatchString(e.UID) || e.UID == "0" {
		return "", errors.New("utilizador de origem desconhecido")
	}
	u, err := user.LookupId(e.UID)
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func (e *Env) addRepo(ctx context.Context, r *catalog.Repo, log func(string)) error {
	if r == nil {
		return errors.New("repositório desconhecido")
	}
	if !r.Verified() {
		return fmt.Errorf("repositório %s sem impressão digital verificada; receita recusada", r.Name)
	}
	suite := strings.ReplaceAll(r.Suite, "{codename}", e.Codename)
	keyring := e.path("/etc/apt/keyrings/meu-ubuntu-" + r.ID + ".gpg")
	sources := e.path("/etc/apt/sources.list.d/meu-ubuntu-" + r.ID + ".sources")
	content := fmt.Sprintf("# Gerido por Meu Ubuntu — %s\n# Documentação: %s\nTypes: deb\nURIs: %s\nSuites: %s\nComponents: %s\nArchitectures: %s\nSigned-By: %s\n",
		r.Name, r.Docs, r.URL, suite, strings.Join(r.Components, " "), r.Arch, "/etc/apt/keyrings/meu-ubuntu-"+r.ID+".gpg")
	if others := aptsrc.ReplacementFiles(e.Root, []byte(content), sources); len(others) > 0 {
		log("Repositório " + r.Name + " já configurado por " + strings.Join(others, ", ") + "; não foi duplicado")
		return nil
	}
	if old, err := os.ReadFile(sources); err == nil && string(old) == content {
		if _, err := os.Stat(keyring); err == nil {
			log("Repositório " + r.Name + " já configurado (sem alterações)")
			return nil
		}
	}
	// The repo must publish this release before we touch APT config.
	rel := strings.TrimSuffix(r.URL, "/") + "/dists/" + suite + "/Release"
	log("A verificar " + rel)
	if _, err := e.Fetch(rel); err != nil {
		return fmt.Errorf("repositório não publica para %q (%v); nada foi alterado", suite, err)
	}
	log("A transferir chave " + r.KeyURL)
	key, err := e.Fetch(r.KeyURL)
	if err != nil {
		return err
	}
	tmp, err := os.MkdirTemp("", "meu-ubuntu-key")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	kf := filepath.Join(tmp, "key")
	if err := os.WriteFile(kf, key, 0o600); err != nil {
		return err
	}
	out, err := e.R.Output(ctx, "gpg", "--homedir", tmp, "--show-keys", "--with-colons", kf)
	if err != nil {
		return err
	}
	allowed := map[string]bool{}
	for _, f := range r.Fingerprints {
		allowed[f] = true
	}
	var primaries []string
	prev := ""
	for _, l := range strings.Split(out, "\n") {
		f := strings.Split(l, ":")
		if f[0] == "fpr" && prev == "pub" && len(f) > 9 {
			primaries = append(primaries, f[9])
		}
		if f[0] != "fpr" {
			prev = f[0]
		}
	}
	if len(primaries) == 0 {
		return errors.New("chave sem impressões digitais legíveis; rejeitada")
	}
	for _, f := range primaries {
		if !allowed[f] {
			return fmt.Errorf("chave inesperada %s no ficheiro transferido; rejeitada (esperado: %s)", f, strings.Join(r.Fingerprints, ", "))
		}
	}
	log("Impressões digitais verificadas: " + strings.Join(primaries, ", "))
	bin := key
	if bytes.Contains(key, []byte("BEGIN PGP PUBLIC KEY BLOCK")) {
		bin, err = runIn(ctx, e.R, "gpg", "--homedir", tmp, "--dearmor", "-o", filepath.Join(tmp, "key.gpg"), kf)
		if err != nil {
			return err
		}
		if bin, err = os.ReadFile(filepath.Join(tmp, "key.gpg")); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(keyring), 0o755); err != nil {
		return err
	}
	if err := backupThenWrite(e, keyring, bin, 0o644); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(sources), 0o755); err != nil {
		return err
	}
	if err := backupThenWrite(e, sources, []byte(content), 0o644); err != nil {
		return err
	}
	log("Fonte APT escrita: " + sources)
	return e.R.Stream(ctx, log, "apt-get", append(aptOpts[:3:3], "update")...)
}

// dedupeSources removes OUR source+keyring when another file configures the
// same repository. Only regular files written by Meu Ubuntu are removed, and
// both are backed up first.
func (e *Env) dedupeSources(ctx context.Context, log func(string)) error {
	dups := aptsrc.Duplicates(e.Root, e.Cat)
	if len(dups) == 0 {
		log("Fontes APT sem duplicados do Meu Ubuntu")
		return nil
	}
	bdir := e.path("/var/backups/meu-ubuntu")
	if err := os.MkdirAll(bdir, 0o700); err != nil {
		return err
	}
	stamp := time.Now().Format("20060102-150405.000000000")
	for _, d := range dups {
		log("Fonte duplicada de " + d.Name + " (também em " + strings.Join(d.Others, ", ") + "): a remover " + d.Ours)
		for _, f := range []string{d.Ours, d.Keyring} {
			if f == "" {
				continue
			}
			fi, err := os.Lstat(e.path(f))
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return err
			}
			if !fi.Mode().IsRegular() {
				return fmt.Errorf("caminho não regular; limpeza recusada: %s", f)
			}
			b, err := os.ReadFile(e.path(f))
			if err != nil {
				return err
			}
			dst := filepath.Join(bdir, filepath.Base(f)+"."+stamp)
			if err := os.WriteFile(dst, b, 0o600); err != nil {
				return fmt.Errorf("backup falhou, nada removido: %w", err)
			}
			log("Backup: " + f + " → " + strings.TrimPrefix(dst, e.Root))
		}
		if err := os.Remove(e.path(d.Ours)); err != nil {
			return err
		}
		if d.Keyring != "" && !aptsrc.KeyringInUse(e.Root, d.Keyring, "") {
			if err := os.Remove(e.path(d.Keyring)); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return e.R.Stream(ctx, log, "apt-get", append(aptOpts[:3:3], "update")...)
}

func runIn(ctx context.Context, r sys.Runner, name string, args ...string) ([]byte, error) {
	out, err := r.Output(ctx, name, args...)
	return []byte(out), err
}

func backupThenWrite(e *Env, p string, data []byte, mode os.FileMode) error {
	if fi, err := os.Lstat(p); err == nil && !fi.Mode().IsRegular() {
		return fmt.Errorf("caminho não regular; escrita recusada: %s", p)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if old, err := os.ReadFile(p); err == nil {
		if bytes.Equal(old, data) {
			return nil
		}
		bdir := e.path("/var/backups/meu-ubuntu")
		if err := os.MkdirAll(bdir, 0o700); err != nil {
			return err
		}
		sum := sha256.Sum256(old)
		name := strings.ReplaceAll(strings.TrimPrefix(p, e.Root), "/", "_") + "." + time.Now().Format("20060102-150405.000000000") + "." + hex.EncodeToString(sum[:4])
		if err := os.WriteFile(filepath.Join(bdir, name), old, 0o600); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return safefile.Write(p, data, mode)
}

func (e *Env) post(ctx context.Context, item *catalog.Item, req Request, log func(string)) error {
	if item == nil {
		return errors.New("item inválido")
	}
	sel, err := e.Cat.ValidateSelection(item.ID, catalog.Selection{Options: req.Options})
	if err != nil {
		return err
	}
	for _, p := range item.Post {
		switch p {
		case "xdebug-off":
			// Keep Xdebug installed but disabled for every selected PHP version.
			for _, v := range sel.Options["versions"] {
				log("Desativar Xdebug para PHP " + v + " (ativar depois com: sudo phpenmod -v " + v + " xdebug)")
				if err := e.R.Stream(ctx, log, "phpdismod", "-v", v, "xdebug"); err != nil {
					return err
				}
			}
			if len(sel.Options["versions"]) == 0 {
				if err := e.R.Stream(ctx, log, "phpdismod", "xdebug"); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("pós-operação desconhecida %q", p)
		}
	}
	return nil
}

var dropins = adminconfig.Dropins

func (e *Env) admin(ctx context.Context, item *catalog.Item, sel catalog.Selection, log func(string)) error {
	switch item.AdminOp {
	case "docker-group":
		u, err := e.username()
		if err != nil {
			return err
		}
		log("Adicionar " + u + " ao grupo docker (efetivo após terminar sessão)")
		return e.R.Stream(ctx, log, "usermod", "-aG", "docker", u)
	case "fstrim-timer":
		return e.R.Stream(ctx, log, "systemctl", "enable", "--now", "fstrim.timer")
	case "snap-remove-firefox":
		// Only the Firefox snap; snapd itself is never removed.
		return e.R.Stream(ctx, log, "snap", "remove", "firefox")
	case "ubuntu-drivers":
		return e.R.Stream(ctx, log, "ubuntu-drivers", "install")
	case "inotify", "journald":
		d := dropins[item.AdminOp]
		val := sel.First("value")
		if val == "" {
			return errors.New("valor obrigatório")
		}
		return e.applyDropin(ctx, item.AdminOp, d.Path, []byte(fmt.Sprintf(d.Template, val)), log)
	}
	return fmt.Errorf("operação administrativa desconhecida %q", item.AdminOp)
}

func (e *Env) revert(ctx context.Context, item *catalog.Item, sel catalog.Selection, log func(string)) error {
	switch item.AdminOp {
	case "docker-group":
		u, err := e.username()
		if err != nil {
			return err
		}
		out, err := e.R.Output(ctx, "id", "-nG", u)
		if err != nil {
			return err
		}
		member := false
		for _, g := range strings.Fields(out) {
			member = member || g == "docker"
		}
		if !member {
			return nil
		}
		return e.R.Stream(ctx, log, "gpasswd", "-d", u, "docker")
	case "inotify", "journald":
		return e.revertDropin(ctx, item.AdminOp, []byte(fmt.Sprintf(dropins[item.AdminOp].Template, sel.First("value"))), log)
	}
	return fmt.Errorf("sem reversão automática para %q", item.AdminOp)
}

// HTTPFetch downloads over HTTPS only, with a size cap.
func HTTPFetch(url string) ([]byte, error) {
	if !httpsURL(url) {
		return nil, errors.New("apenas HTTPS")
	}
	c := &http.Client{Timeout: 60 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if !httpsURL(req.URL.String()) {
			return errors.New("redirecionamento fora de HTTPS recusado")
		}
		if len(via) >= 10 {
			return errors.New("demasiados redirecionamentos")
		}
		return nil
	}}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	const limit = 20 << 20
	b, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if len(b) > limit {
		return nil, errors.New("transferência excede o limite de 20 MiB")
	}
	return b, nil
}

func httpsURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host != "" && u.User == nil
}

// Serve is the root side: one JSON request per line on in, messages on out.
// A system-wide lock prevents two administrative sessions from running at once.
func Serve(env *Env, in io.Reader, out io.Writer) error {
	enc := json.NewEncoder(out)
	if env.Root == "" {
		lf, err := os.OpenFile("/run/meu-ubuntu-helper.lock", os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW, 0o600)
		if err != nil {
			return err
		}
		defer lf.Close()
		if err := syscall.Flock(int(lf.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
			enc.Encode(Message{Type: "done", Err: "outra sessão administrativa do Meu Ubuntu está em curso"})
			return err
		}
	}
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		var req Request
		if err := json.Unmarshal(sc.Bytes(), &req); err != nil {
			enc.Encode(Message{Type: "done", Err: "pedido inválido"})
			continue
		}
		fmt.Fprintf(os.Stderr, "meu-ubuntu helper: %s %s\n", req.Op, req.Item) // audit trail (journal via pkexec)
		var preview *RevertPreview
		var err error
		if req.Op == "revert-preview" {
			preview, err = env.PreviewRevert(context.Background(), req)
		} else {
			err = env.Execute(context.Background(), req, func(l string) { enc.Encode(Message{Type: "log", Line: l}) })
		}
		m := Message{Type: "done", OK: err == nil}
		m.Preview = preview
		if err != nil {
			m.Err = err.Error()
		}
		enc.Encode(m)
	}
	return sc.Err()
}
