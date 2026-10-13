package integrations

import (
	"bufio"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"meuubuntu/internal/safefile"
	"meuubuntu/internal/store"
)

// The executable belongs to the application, never to an imported archive.
//
//go:embed prepare-compose.py
var preparationScript []byte

type SetupReview struct {
	Token       string          `json:"token"`
	Project     string          `json:"project"`
	Namespace   string          `json:"namespace"`
	Services    []string        `json:"services"`
	Databases   []DatabaseSetup `json:"databases"`
	Steps       []string        `json:"steps"`
	Manual      []string        `json:"manual"`
	Ports       map[string]int  `json:"ports"`
	Address     string          `json:"address"`
	Simulated   bool            `json:"simulated"`
	helperToken string
	view        View
	dir         string
	workspace   string
	created     time.Time
}

type SetupStep struct {
	Title   string `json:"title"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
type SetupJob struct {
	ID        string      `json:"id"`
	Package   string      `json:"package"`
	Started   string      `json:"started"`
	Updated   string      `json:"updated"`
	State     string      `json:"state"`
	Finished  bool        `json:"finished"`
	Simulated bool        `json:"simulated"`
	Steps     []SetupStep `json:"steps"`
	Manual    []string    `json:"manual"`
	Address   string      `json:"address"`
	Terminal  string      `json:"terminal"`
	Message   string      `json:"message,omitempty"`
}

type Preparer struct {
	Manager Manager
	mu      sync.Mutex
	review  *SetupReview
	job     *SetupJob
	control string
	active  bool
}

func (p *Preparer) Running() bool {
	if p == nil {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active || p.externalActiveLocked()
}

// A child completing a safe operation can outlive the GUI after a crash.
// Share its kernel lock rather than trusting only the in-memory job flag.
func (p *Preparer) externalActiveLocked() bool {
	path := filepath.Join(p.Manager.Data, "setup.lock")
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil || !info.Mode().IsRegular() {
		return true
	}
	f, err := os.OpenFile(path, os.O_RDWR|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return true
	}
	defer f.Close()
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return true
	}
	_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return false
}
func binding(v View) string { b, _ := json.Marshal(v); return fmt.Sprintf("%x", sha256.Sum256(b)) }

func (p *Preparer) Review(ctx context.Context, id string) (*SetupReview, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active || p.externalActiveLocked() {
		return nil, errors.New("já existe uma preparação em curso")
	}
	v, err := p.Manager.Load(id)
	if err != nil {
		return nil, err
	}
	if v.Project == "" || v.Setup == nil {
		return nil, errors.New("escolha a pasta de um pacote com preparação automática")
	}
	if err := v.Setup.Validate(); err != nil {
		return nil, err
	}
	if p.review != nil {
		os.RemoveAll(p.review.dir)
		p.review = nil
	}
	if err := os.MkdirAll(p.Manager.Data, 0700); err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp(p.Manager.Data, "setup-review-")
	if err != nil {
		return nil, err
	}
	ok := false
	defer func() {
		if !ok {
			os.RemoveAll(dir)
		}
	}()
	workspace := filepath.Join(p.Manager.Data, "dev-environments", fmt.Sprintf("%x", sha256.Sum256([]byte(v.Project)))[:16])
	spec, _ := json.Marshal(v.Setup)
	if err = os.WriteFile(filepath.Join(dir, "spec.json"), spec, 0600); err != nil {
		return nil, err
	}
	if err = os.WriteFile(filepath.Join(dir, "prepare.py"), preparationScript, 0600); err != nil {
		return nil, err
	}
	r := &SetupReview{}
	if p.Manager.Simulated {
		r.Project, r.Namespace, r.Simulated = v.Project, "exemplo-controlado", true
		r.Steps = []string{"Preparar configuração pessoal", "Construir imagens", "Iniciar containers", "Criar/conferir schemas", "Conferir Redis Cluster", "Instalar Composer", "Conferir serviços"}
		r.Services, r.Databases = []string{v.Setup.Application, v.Setup.Initializer}, v.Setup.Databases
		r.Manual = []string{"Exemplo de simulação; nenhum projeto, container ou base de dados será alterado."}
		r.Ports = map[string]int{"http": 8089, "https": 8449}
		r.helperToken = "simulated"
	} else {
		cmd := exec.CommandContext(ctx, "python3", filepath.Join(dir, "prepare.py"), "--project", v.Project, "--spec", filepath.Join(dir, "spec.json"), "--workspace", workspace, "--lock", filepath.Join(p.Manager.Data, "setup.lock"), "--review")
		var output limitedBuffer
		cmd.Stdout, cmd.Stderr = &output, io.Discard
		err = cmd.Run()
		if decodeErr := json.Unmarshal(output.data, r); decodeErr != nil || err != nil {
			var failure struct {
				Message string `json:"message"`
			}
			json.Unmarshal(output.data, &failure)
			if failure.Message != "" {
				return nil, errors.New(failure.Message)
			}
			return nil, errors.New("não foi possível rever; confirme Python 3, Docker/Compose e a configuração do projeto")
		}
		r.helperToken = r.Token
	}
	r.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(store.NewID()+binding(v)+r.helperToken)))
	r.view, r.dir, r.workspace, r.created = v, dir, workspace, time.Now()
	p.review = r
	ok = true
	return cloneReview(r), nil
}

func cloneReview(r *SetupReview) *SetupReview {
	copy := *r
	copy.Steps = append([]string(nil), r.Steps...)
	copy.Manual = append([]string(nil), r.Manual...)
	copy.Services = append([]string(nil), r.Services...)
	copy.Databases = make([]DatabaseSetup, len(r.Databases))
	for i, d := range r.Databases {
		copy.Databases[i] = d
		copy.Databases[i].Schemas = append([]string(nil), d.Schemas...)
	}
	copy.Ports = make(map[string]int, len(r.Ports))
	for k, v := range r.Ports {
		copy.Ports[k] = v
	}
	return &copy
}

type limitedBuffer struct{ data []byte }

func (b *limitedBuffer) Write(v []byte) (int, error) {
	n := len(v)
	remaining := 256*1024 - len(b.data)
	if remaining > 0 {
		if len(v) > remaining {
			v = v[:remaining]
		}
		b.data = append(b.data, v...)
	}
	return n, nil
}

func (p *Preparer) Start(id, token string) (*SetupJob, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	r := p.review
	if p.active || p.externalActiveLocked() || r == nil || r.view.ID != id || token == "" || token != r.Token || time.Since(r.created) > 10*time.Minute {
		return nil, errors.New("reveja a preparação antes de executar")
	}
	current, err := p.Manager.Load(id)
	if err != nil || binding(current) != binding(r.view) {
		return nil, errors.New("o pacote ou a pasta mudou; reveja novamente")
	}
	job := &SetupJob{ID: store.NewID(), Package: id, Started: time.Now().Format(time.RFC3339), State: "running", Simulated: p.Manager.Simulated, Manual: r.Manual, Address: r.Address}
	if !p.Manager.Simulated {
		job.Terminal = filepath.Join(r.workspace, "dock")
	}
	for _, title := range r.Steps {
		job.Steps = append(job.Steps, SetupStep{Title: title, Status: "pending"})
	}
	p.job, p.control = job, filepath.Join(r.dir, "control")
	if err = p.saveLocked(); err != nil {
		p.job = nil
		return nil, err
	}
	p.active, p.review = true, nil
	go p.execute(r)
	return p.copyLocked(), nil
}

func (p *Preparer) copyLocked() *SetupJob {
	if p.job == nil {
		return nil
	}
	copy := *p.job
	copy.Steps = append([]SetupStep(nil), p.job.Steps...)
	copy.Manual = append([]string(nil), p.job.Manual...)
	return &copy
}
func (p *Preparer) Status() *SetupJob {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.job == nil {
		path := filepath.Join(p.Manager.Data, "setup-last.json")
		info, err := os.Lstat(path)
		var data []byte
		if err == nil && info.Mode().IsRegular() && info.Size() < 256*1024 {
			data, err = os.ReadFile(path)
		}
		if err == nil && len(data) > 0 {
			var j SetupJob
			if json.Unmarshal(data, &j) == nil {
				p.job = &j
			}
		}
	}
	if p.job != nil && !p.active && !p.job.Finished {
		if p.externalActiveLocked() {
			p.job.Message = "Processo anterior ainda termina a operação atual. Aguarde antes de rever novamente."
		} else {
			p.job.State, p.job.Finished, p.job.Message = "interrupted", true, "Preparação interrompida. Reveja novamente; o estado real será conferido."
		}
	}
	return p.copyLocked()
}
func (p *Preparer) saveLocked() error {
	p.job.Updated = time.Now().Format(time.RFC3339)
	b, err := json.Marshal(p.job)
	if err != nil {
		return err
	}
	return safefile.Write(filepath.Join(p.Manager.Data, "setup-last.json"), b, 0600)
}
func (p *Preparer) Control(action string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.active {
		return errors.New("não existe preparação em curso")
	}
	switch action {
	case "pause":
		if p.job.State == "cancelling" {
			return errors.New("já foi pedido cancelamento")
		}
		if err := os.WriteFile(p.control+".pause", nil, 0600); err != nil {
			return err
		}
		if p.job.State != "paused" {
			p.job.State = "pausing"
		}
	case "resume":
		if p.job.State == "cancelling" {
			return errors.New("já foi pedido cancelamento")
		}
		if err := os.Remove(p.control + ".pause"); err != nil {
			return errors.New("a preparação não está em pausa")
		}
		p.job.State = "running"
	case "cancel":
		if err := os.WriteFile(p.control+".cancel", nil, 0600); err != nil {
			return err
		}
		p.job.State = "cancelling"
	default:
		return errors.New("ação inválida")
	}
	return p.saveLocked()
}

type setupEvent struct {
	Index    *int   `json:"index"`
	Status   string `json:"status"`
	State    string `json:"state"`
	Message  string `json:"message"`
	Terminal string `json:"terminal"`
}

func (p *Preparer) receive(event setupEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if event.State == "paused" || event.State == "running" {
		if p.job.State != "cancelling" {
			_, err := os.Stat(p.control + ".pause")
			if event.State == "paused" && err == nil {
				p.job.State = "paused"
			}
			if event.State == "running" && err != nil {
				p.job.State = "running"
			}
		}
	}
	if event.Index != nil && *event.Index >= 0 && *event.Index < len(p.job.Steps) {
		step := &p.job.Steps[*event.Index]
		step.Status, step.Message = event.Status, event.Message
	} else if event.Status != "" {
		p.job.State, p.job.Message = event.Status, event.Message
		if event.Terminal != "" {
			p.job.Terminal = event.Terminal
		}
		p.job.Finished = true
	}
	if err := p.saveLocked(); err != nil {
		p.job.Message = "Falha ao guardar progresso; paragem pedida no próximo limite seguro."
		os.WriteFile(p.control+".cancel", nil, 0600)
	}
}

func (p *Preparer) execute(r *SetupReview) {
	defer os.RemoveAll(r.dir)
	defer func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.active = false
		if !p.job.Finished {
			p.job.State, p.job.Finished, p.job.Message = "failed", true, "Preparação não concluída; reveja novamente. Dados preservados."
		}
		for i := range p.job.Steps {
			if p.job.Steps[i].Status == "pending" || p.job.Steps[i].Status == "running" {
				p.job.Steps[i].Status = "cancelled"
			}
		}
		if err := p.saveLocked(); err != nil {
			p.job.State, p.job.Message = "failed", "Não foi possível guardar o histórico da preparação."
		}
	}()
	if p.Manager.Simulated {
		for i := range r.Steps {
			for {
				if _, err := os.Stat(p.control + ".cancel"); err == nil {
					p.receive(setupEvent{Status: "cancelled"})
					return
				}
				if _, err := os.Stat(p.control + ".pause"); err != nil {
					break
				}
				p.receive(setupEvent{State: "paused"})
				time.Sleep(100 * time.Millisecond)
			}
			p.receive(setupEvent{State: "running"})
			index := i
			p.receive(setupEvent{Index: &index, Status: "running"})
			time.Sleep(700 * time.Millisecond)
			p.receive(setupEvent{Index: &index, Status: "done"})
		}
		p.receive(setupEvent{Status: "done", Message: "Exemplo concluído; nada foi instalado."})
		return
	}
	ports, _ := json.Marshal(r.Ports)
	cmd := exec.Command("python3", filepath.Join(r.dir, "prepare.py"), "--project", r.Project, "--spec", filepath.Join(r.dir, "spec.json"), "--workspace", r.workspace, "--lock", filepath.Join(p.Manager.Data, "setup.lock"), "--token", r.helperToken, "--ports", string(ports), "--control", p.control)
	cmd.Stderr = io.Discard
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 4096), 64*1024)
	for scanner.Scan() {
		var event setupEvent
		if json.Unmarshal(scanner.Bytes(), &event) == nil {
			p.receive(event)
		}
	}
	if scanner.Err() != nil {
		os.WriteFile(p.control+".cancel", nil, 0600)
		io.Copy(io.Discard, stdout)
	}
	if err = cmd.Wait(); err != nil {
		p.mu.Lock()
		if p.job.State == "done" || !p.job.Finished {
			p.job.State, p.job.Finished, p.job.Message = "failed", true, "O processo de preparação falhou; reveja novamente antes de continuar."
		}
		p.mu.Unlock()
	}
}
