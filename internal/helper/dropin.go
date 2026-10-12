package helper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"meuubuntu/internal/safefile"
)

// This snapshot is root-owned. Reversal never accepts file contents or paths
// from the GUI or its user-writable journal.
type dropinSnapshot struct {
	Before        []byte          `json:"before"`
	After         []byte          `json:"after"`
	HadValue      bool            `json:"hadValue"`
	Mode          os.FileMode     `json:"mode"`
	RuntimeBefore string          `json:"runtimeBefore,omitempty"`
	Previous      *dropinSnapshot `json:"previous,omitempty"`
}

func (e *Env) snapshotPath(op string) string { return e.path("/var/lib/meu-ubuntu/" + op + ".json") }

func regularContents(path string) ([]byte, os.FileMode, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, 0, err
	}
	if !fi.Mode().IsRegular() {
		return nil, 0, errors.New("caminho não regular; operação recusada")
	}
	b, err := os.ReadFile(path)
	return b, fi.Mode().Perm(), err
}

func (e *Env) applyDropin(ctx context.Context, op, path string, after []byte, log func(string)) error {
	p := e.path(path)
	old, mode, err := regularContents(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	s := dropinSnapshot{Before: old, After: after, HadValue: err == nil, Mode: mode}
	state := e.snapshotPath(op)
	if b, _, err := regularContents(state); err == nil {
		var previous dropinSnapshot
		if err := json.Unmarshal(b, &previous); err != nil {
			return fmt.Errorf("histórico administrativo ilegível: %w", err)
		}
		if !bytes.Equal(old, previous.After) && !(bytes.Equal(old, previous.Before) && bytes.Equal(after, previous.After)) {
			return errors.New("ficheiro alterado depois da aplicação; preserve ou reverta essa alteração primeiro")
		}
		if bytes.Equal(old, previous.After) && !bytes.Equal(old, after) {
			s.Previous = &previous
		} else {
			s = previous // retry the same change without losing its baseline
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if op == "inotify" && s.RuntimeBefore == "" {
		b, err := os.ReadFile(e.path("/proc/sys/fs/inotify/max_user_watches"))
		if err != nil {
			return fmt.Errorf("não foi possível guardar o limite inotify anterior: %w", err)
		}
		s.RuntimeBefore = strings.TrimSpace(string(b))
		if !regexp.MustCompile(`^[0-9]+$`).MatchString(s.RuntimeBefore) {
			return errors.New("limite inotify anterior inválido")
		}
	}
	if err := os.MkdirAll(filepath.Dir(state), 0o700); err != nil {
		return err
	}
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := safefile.Write(state, b, 0o600); err != nil {
		return fmt.Errorf("não foi possível guardar a configuração anterior; nada alterado: %w", err)
	}
	if err := os.Remove(p + ".applied"); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	if err := backupThenWrite(e, p, after, 0o644); err != nil {
		return err
	}
	log("Escrito " + path + " (configuração anterior guardada)")
	value := ""
	if op == "inotify" {
		m := regexp.MustCompile(`fs\.inotify\.max_user_watches = ([0-9]+)`).FindSubmatch(after)
		if len(m) != 2 {
			return errors.New("limite inotify inválido")
		}
		value = string(m[1])
	}
	if err := e.reloadDropin(ctx, op, value, log); err != nil {
		return err
	}
	// Root-owned proof is written only after a successful reload. A file alone
	// must not turn a failed reload into an apparently complete configuration.
	return safefile.Write(p+".applied", []byte(fmt.Sprintf("%x\n", sha256.Sum256(after))), 0644)
}

func (e *Env) revertDropin(ctx context.Context, op string, expected []byte, log func(string)) error {
	state := e.snapshotPath(op)
	b, _, err := regularContents(state)
	if err != nil {
		return fmt.Errorf("sem configuração anterior verificável; reversão recusada: %w", err)
	}
	var s dropinSnapshot
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !bytes.Equal(s.After, expected) {
		return errors.New("existe uma alteração administrativa posterior; reverta-a primeiro")
	}
	p := e.path(dropins[op].Path)
	cur, _, err := regularContents(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	alreadyRestored := (s.HadValue && err == nil && bytes.Equal(cur, s.Before)) || (!s.HadValue && os.IsNotExist(err))
	if err := os.Remove(p + ".applied"); err != nil && !os.IsNotExist(err) {
		return err
	}
	if !alreadyRestored {
		if err != nil || !bytes.Equal(cur, s.After) {
			return errors.New("ficheiro alterado depois da aplicação; reversão recusada")
		}
		if s.HadValue {
			if err := safefile.Write(p, s.Before, s.Mode); err != nil {
				return err
			}
		} else if err := os.Remove(p); err != nil {
			return err
		}
	}
	if err := e.reloadDropin(ctx, op, s.RuntimeBefore, log); err != nil {
		return err
	}
	if s.HadValue {
		if err := safefile.Write(p+".applied", []byte(fmt.Sprintf("%x\n", sha256.Sum256(s.Before))), 0644); err != nil {
			return err
		}
	}
	log("Restaurada a configuração anterior: " + dropins[op].Path)
	if s.Previous != nil {
		b, err := json.Marshal(s.Previous)
		if err != nil {
			return err
		}
		return safefile.Write(state, b, 0600)
	}
	return os.Remove(state)
}

func (e *Env) reloadDropin(ctx context.Context, op, runtimeBefore string, log func(string)) error {
	if op == "inotify" && runtimeBefore != "" {
		if !regexp.MustCompile(`^[0-9]+$`).MatchString(runtimeBefore) {
			return errors.New("limite anterior inválido")
		}
		if err := e.R.Stream(ctx, log, "sysctl", "-w", "fs.inotify.max_user_watches="+runtimeBefore); err != nil {
			return err
		}
		if e.Root != "" {
			return os.WriteFile(e.path("/proc/sys/fs/inotify/max_user_watches"), []byte(runtimeBefore+"\n"), 0600)
		}
		return nil
	}
	for _, c := range dropins[op].Apply {
		if err := e.R.Stream(ctx, log, c[0], c[1:]...); err != nil {
			return err
		}
	}
	return nil
}
