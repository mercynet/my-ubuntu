package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/plan"
	"meuubuntu/internal/store"
)

// ReviewUndo binds batch reversals to the journal and values actually shown.
func (e *Engine) ReviewUndo(ctx context.Context, p *plan.Plan) error {
	// The frontend iterates steps even when every action is already satisfied.
	keptSteps := make([]plan.Step, 0, len(p.Steps))
	for i := range p.Steps {
		s := &p.Steps[i]
		if s.Op != "undo" && s.Op != "revert" {
			keptSteps = append(keptSteps, *s)
			continue
		}
		var values []plan.UndoValue
		var err error
		if s.Op == "revert" {
			values, s.Options, err = e.adminUndoValues(ctx, s.Item, nil)
			if err == nil && len(values) == 0 {
				p.Kept = append(p.Kept, plan.Skipped{Item: s.Item, Name: e.Cat.Item(s.Item).Name, Reason: "Sem histórico administrativo verificável; configuração mantida", Status: detect.Unknown})
				continue
			}
		} else {
			values, err = e.undoValues(ctx, s.Item)
		}
		if err != nil {
			return err
		}
		s.Undo, s.Detail, s.Changes = values, nil, nil
		for _, v := range values {
			s.Detail = append(s.Detail, fmt.Sprintf("%s %s: %s → %s", v.Schema, v.Key, v.Current, v.Target))
			s.Changes = append(s.Changes, detect.SettingState{Schema: v.Schema, Key: v.Key, Current: v.Current, Proposed: v.Target})
		}
		keptSteps = append(keptSteps, *s)
	}
	p.Steps = keptSteps
	p.NeedsRoot = false
	for i := range p.Steps {
		p.Steps[i].ID = fmt.Sprintf("s%02d", i+1)
		p.NeedsRoot = p.NeedsRoot || p.Steps[i].Kind == "admin"
	}
	return nil
}

func (e *Engine) undoValues(ctx context.Context, item string) ([]plan.UndoValue, error) {
	it := e.Cat.Item(item)
	changes, err := e.Store.ReadChanges()
	if err != nil {
		return nil, fmt.Errorf("histórico ilegível; reversão recusada: %w", err)
	}
	latest := func(kind, schema, key string) *store.Change {
		for i := len(changes) - 1; i >= 0; i-- {
			c := changes[i]
			if c.Item == item && c.Kind == kind && c.Schema == schema && c.Key == key {
				return &c
			}
		}
		return nil
	}
	var values []plan.UndoValue
	add := func(kind, schema, key, current string) {
		v := plan.UndoValue{Kind: kind, Schema: schema, Key: key, Current: current, Mode: "reset", Target: "padrão do sistema"}
		if c := latest(kind, schema, key); c != nil {
			v.ChangeID, v.Target, v.Mode = c.ID, c.Before, "journal"
			if kind == "git" && !c.HadValue {
				v.Target = "(ausente)"
			}
			if kind == "template" {
				v.Target = "(remover o ficheiro criado)"
			}
			if c.Reverted != "" {
				v.Mode, v.Target = "skip", "já restaurado; manter o valor atual"
			}
		} else if kind != "gsettings" {
			v.Mode, v.Target = "skip", "sem histórico; manter o valor atual"
		}
		values = append(values, v)
	}
	switch it.Method {
	case catalog.MethodGSettings:
		for _, st := range it.Settings {
			schema, ok := detect.ResolveSchema(ctx, e.R, st.Schema)
			if !ok {
				return nil, errors.New("não foi possível ler o esquema; reversão recusada")
			}
			cur, ok := detect.GSettingsGet(ctx, e.R, schema, st.Key)
			if !ok {
				return nil, errors.New("não foi possível ler a preferência; reversão recusada")
			}
			add("gsettings", schema, st.Key, cur)
		}
	case catalog.MethodGitConfig:
		for _, st := range it.Git {
			cur, err := e.gitValue(ctx, st.Key)
			if err != nil {
				return nil, err
			}
			add("git", "", st.Key, cur)
		}
	case catalog.MethodTemplate:
		key := filepath.Join(detect.TemplatesDir(e.Home), filepath.Base(it.Template.Name))
		b, err := os.ReadFile(key)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		add("template", "", key, string(b))
	}
	return values, nil
}

func (e *Engine) validateUndo(ctx context.Context, s plan.Step) error {
	actual, err := e.undoValues(ctx, s.Item)
	if err != nil {
		return err
	}
	a, _ := json.Marshal(actual)
	b, _ := json.Marshal(s.Undo)
	if string(a) != string(b) {
		return errors.New("o valor ou histórico mudou desde a revisão; reveja a reversão novamente")
	}
	return nil
}
