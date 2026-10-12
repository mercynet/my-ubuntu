package helper

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Admin sends requests to a privileged executor.
type Admin interface {
	Do(ctx context.Context, req Request, log func(string)) error
	Close()
}

type PreviewAdmin interface {
	Preview(context.Context, Request) (*RevertPreview, error)
}

// InProcess runs the helper logic directly (tests / simulation only).
type InProcess struct{ Env *Env }

func (p InProcess) Do(ctx context.Context, req Request, log func(string)) error {
	return p.Env.Execute(ctx, req, log)
}
func (InProcess) Close() {}
func (p InProcess) Preview(ctx context.Context, req Request) (*RevertPreview, error) {
	return p.Env.PreviewRevert(ctx, req)
}

// Pkexec starts `pkexec <self> --helper` lazily, once, and keeps it for the run.
// The password prompt is the system polkit agent; nothing is stored.
type Pkexec struct {
	mu    sync.Mutex
	cmd   *exec.Cmd
	in    io.WriteCloser
	out   *bufio.Scanner
	start error
}

func (p *Pkexec) ensure() error {
	if p.cmd != nil || p.start != nil {
		return p.start
	}
	self, err := os.Executable()
	if err != nil {
		p.start = err
		return err
	}
	cmd := exec.Command("pkexec", self, "--helper")
	in, _ := cmd.StdinPipe()
	outp, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		p.start = err
		return err
	}
	p.cmd, p.in = cmd, in
	p.out = bufio.NewScanner(outp)
	p.out.Buffer(make([]byte, 64*1024), 1024*1024)
	return nil
}

// Do is serialized: only one administrative (APT) operation at a time.
func (p *Pkexec) Do(ctx context.Context, req Request, log func(string)) error {
	_, err := p.exchange(ctx, req, log)
	return err
}

func (p *Pkexec) Preview(ctx context.Context, req Request) (*RevertPreview, error) {
	req.Op = "revert-preview"
	return p.exchange(ctx, req, func(string) {})
}

func (p *Pkexec) exchange(ctx context.Context, req Request, log func(string)) (*RevertPreview, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ensure(); err != nil {
		return nil, err
	}
	b, _ := json.Marshal(req)
	if _, err := p.in.Write(append(b, '\n')); err != nil {
		return nil, errors.New("modo administrativo terminou (autenticação cancelada?)")
	}
	for p.out.Scan() {
		var m Message
		if json.Unmarshal(p.out.Bytes(), &m) != nil {
			continue
		}
		if m.Type == "log" {
			log(m.Line)
			continue
		}
		if m.OK {
			return m.Preview, nil
		}
		return nil, errors.New(m.Err)
	}
	p.start = errors.New("autenticação recusada ou modo administrativo terminou")
	return nil, p.start
}

func (p *Pkexec) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.in != nil {
		p.in.Close()
		p.cmd.Wait()
		p.cmd, p.in = nil, nil
	}
	p.start = nil
}
