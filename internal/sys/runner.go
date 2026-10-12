// Package sys wraps process execution. Every command is an argv slice
// (never a shell string), so nothing from profiles is interpreted by a shell.
package sys

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Runner interface {
	// Output runs argv and returns combined stdout (stderr is appended on error).
	Output(ctx context.Context, name string, args ...string) (string, error)
	// Stream runs argv, calling onLine for each stdout/stderr line.
	Stream(ctx context.Context, onLine func(string), name string, args ...string) error
	LookPath(name string) bool
}

type Real struct {
	Env []string // extra env vars
}

// ExitCode works through wrapped command errors without parsing error text.
func ExitCode(err error) int {
	var exit interface {
		error
		ExitCode() int
	}
	if errors.As(err, &exit) {
		return exit.ExitCode()
	}
	if err == nil {
		return 0
	}
	return -1
}

type ExitStatus int

func (e ExitStatus) Error() string { return fmt.Sprintf("exit status %d", e) }
func (e ExitStatus) ExitCode() int { return int(e) }

func (r Real) cmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	c := exec.CommandContext(ctx, name, args...)
	c.Env = append(os.Environ(), "LC_ALL=C.UTF-8", "DEBIAN_FRONTEND=noninteractive")
	c.Env = append(c.Env, r.Env...)
	return c
}

func (r Real) Output(ctx context.Context, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	c := r.cmd(ctx, name, args...)
	var out, errb bytes.Buffer
	c.Stdout, c.Stderr = &out, &errb
	err := c.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%s: %w: %s", name, err, strings.TrimSpace(errb.String()))
	}
	return out.String(), nil
}

func (r Real) Stream(ctx context.Context, onLine func(string), name string, args ...string) error {
	c := r.cmd(ctx, name, args...)
	pr, pw := io.Pipe()
	c.Stdout, c.Stderr = pw, pw
	if err := c.Start(); err != nil {
		pr.Close()
		pw.Close()
		return err
	}
	var wg sync.WaitGroup
	var scanErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		sc := bufio.NewScanner(pr)
		sc.Buffer(make([]byte, 64*1024), 1024*1024)
		for sc.Scan() {
			onLine(sc.Text())
		}
		scanErr = sc.Err()
		// A long line must not leave the writer blocked and deadlock c.Wait.
		if scanErr != nil {
			io.Copy(io.Discard, pr)
		}
		pr.Close()
	}()
	err := c.Wait()
	pw.Close()
	wg.Wait()
	if err == nil && scanErr != nil {
		return fmt.Errorf("log de %s: %w", name, scanErr)
	}
	return err
}

func (Real) LookPath(name string) bool { _, err := exec.LookPath(name); return err == nil }

// Argv renders a command for logs/previews (display only, never executed).
func Argv(name string, args ...string) string {
	parts := []string{name}
	for _, a := range args {
		if a == "" || strings.ContainsAny(a, " '\"$`\\") {
			a = "'" + strings.ReplaceAll(a, "'", `'\''`) + "'"
		}
		parts = append(parts, a)
	}
	return strings.Join(parts, " ")
}
