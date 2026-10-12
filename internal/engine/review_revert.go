package engine

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"meuubuntu/internal/detect"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/store"
	"meuubuntu/internal/sys"
)

type RevertPreview struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Current      string `json:"current"`
	Target       string `json:"target"`
	Warning      string `json:"warning,omitempty"`
	CanRevert    bool   `json:"canRevert"`
	ForceAllowed bool   `json:"forceAllowed"`
	NeedsRoot    bool   `json:"needsRoot"`
	ReviewToken  string `json:"reviewToken"`
	StateToken   string `json:"stateToken,omitempty"`
}

func (e *Engine) changeByID(id string) (store.Change, error) {
	changes, err := e.Store.ReadChanges()
	if err != nil {
		return store.Change{}, err
	}
	for _, c := range changes {
		if c.ID == id {
			if c.Reverted != "" {
				return store.Change{}, errors.New("esta alteração já foi revertida")
			}
			return c, nil
		}
	}
	return store.Change{}, errors.New("alteração não encontrada")
}

func (e *Engine) PreviewRevert(ctx context.Context, id string) (*RevertPreview, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.running {
		return nil, errors.New("aguarde o fim da execução antes de rever uma reversão")
	}
	c, err := e.changeByID(id)
	if err != nil {
		return nil, err
	}
	return e.revertPreview(ctx, c)
}

func (e *Engine) revertPreview(ctx context.Context, c store.Change) (*RevertPreview, error) {
	p := &RevertPreview{ID: c.ID, Title: c.Title, Target: c.Before, ForceAllowed: true}
	switch c.Kind {
	case "gsettings":
		cur, ok := detect.GSettingsGet(ctx, e.R, c.Schema, c.Key)
		if !ok {
			return nil, errors.New("não foi possível ler o valor atual; nada alterado")
		}
		p.Current, p.CanRevert = cur, cur == c.After
	case "git":
		cur, err := e.gitValue(ctx, c.Key)
		if err != nil {
			return nil, err
		}
		p.Current, p.CanRevert = cur, cur == c.After
		if !c.HadValue {
			p.Target = "(remover a preferência criada pela aplicação)"
		}
	case "power":
		cur, err := e.R.Output(ctx, "powerprofilesctl", "get")
		if err != nil {
			return nil, err
		}
		p.Current = strings.TrimSpace(cur)
		p.CanRevert = p.Current == c.After
	case "template":
		fi, err := os.Lstat(c.Key)
		if os.IsNotExist(err) {
			p.Current, p.CanRevert = "(ficheiro ausente)", true
		} else if err != nil {
			return nil, err
		} else if !fi.Mode().IsRegular() {
			return nil, errors.New("modelo não é um ficheiro regular; reversão recusada")
		} else if fi.Size() > 256*1024 {
			return nil, errors.New("modelo demasiado grande para revisão; reversão recusada")
		} else {
			b, err := os.ReadFile(c.Key)
			if err != nil {
				return nil, err
			}
			p.Current, p.CanRevert = string(b), string(b) == c.After
		}
		p.Target = "(remover o ficheiro criado pela aplicação)"
	case "admin":
		a := e.NewAdmin()
		defer a.Close()
		reader, ok := a.(helper.PreviewAdmin)
		if !ok {
			return nil, errors.New("modo administrativo não suporta revisão")
		}
		var options map[string][]string
		if c.After != "" {
			options = map[string][]string{"value": strings.Split(c.After, ",")}
		}
		r, err := reader.Preview(ctx, helper.Request{Op: "revert-preview", Item: c.Item, Options: options})
		if err != nil {
			return nil, err
		}
		if r == nil {
			return nil, errors.New("modo administrativo não devolveu a revisão")
		}
		p.Current, p.Target, p.CanRevert, p.Warning = r.Current, r.Target, r.CanRevert, r.Warning
		p.StateToken = r.StateToken
		p.NeedsRoot, p.ForceAllowed = true, false
	default:
		return nil, errors.New("tipo de alteração sem reversão automática")
	}
	if !p.CanRevert && p.Warning == "" {
		p.Warning = "O valor atual mudou depois da aplicação. Repor o anterior substituiria essa alteração posterior."
	}
	if c.Status == "pending" {
		p.Warning += " A aplicação original não tem confirmação final no histórico."
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	p.ReviewToken = fmt.Sprintf("%x", sha256.Sum256(b))
	return p, nil
}

// Re-read the actual value immediately before mutation. A stale review never authorises a new value.
func (e *Engine) RevertReviewed(ctx context.Context, id string, force bool, token string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.running {
		return errors.New("aguarde o fim da execução antes de reverter alterações")
	}
	c, err := e.changeByID(id)
	if err != nil {
		return err
	}
	p, err := e.revertPreview(ctx, c)
	if err != nil {
		return err
	}
	if token == "" || token != p.ReviewToken {
		return errors.New("o valor mudou desde a revisão; reveja a reversão novamente")
	}
	if !p.CanRevert && !(force && p.ForceAllowed) {
		return errors.New("a reversão precisa de uma decisão explícita sobre a alteração posterior")
	}
	return e.Revert(ctx, c, force, p.StateToken)
}

func (e *Engine) gitValue(ctx context.Context, key string) (string, error) {
	out, err := e.R.Output(ctx, "git", "config", "--global", "--get", key)
	if err != nil && sys.ExitCode(err) != 1 {
		return "", fmt.Errorf("não foi possível ler o valor atual do Git: %w", err)
	}
	return strings.TrimSpace(out), nil
}
