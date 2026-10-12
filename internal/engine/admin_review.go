package engine

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"meuubuntu/internal/helper"
	"meuubuntu/internal/plan"
)

// Administrative batch previews bind the entire trusted snapshot chain, the
// invoking account, actual state, target and participating journal IDs.
func (e *Engine) adminUndoValues(ctx context.Context, item string, admin helper.Admin) ([]plan.UndoValue, map[string][]string, error) {
	changes, err := e.Store.ReadChanges()
	if err != nil {
		return nil, nil, err
	}
	var ids []string
	var options map[string][]string
	for _, c := range changes {
		if c.Item == item && c.Kind == "admin" && c.Reverted == "" {
			ids = append(ids, c.ID)
			if c.After != "" {
				options = map[string][]string{"value": strings.Split(c.After, ",")}
			}
		}
	}
	if len(ids) == 0 {
		return nil, nil, nil
	}
	if admin == nil {
		admin = e.NewAdmin()
		defer admin.Close()
	}
	reader, ok := admin.(helper.PreviewAdmin)
	if !ok {
		return nil, nil, errors.New("modo administrativo sem revisão verificável")
	}
	p, err := reader.Preview(ctx, helper.Request{Item: item, Options: options, Batch: true})
	if err != nil {
		return nil, nil, err
	}
	if p == nil || p.StateToken == "" {
		return nil, nil, errors.New("revisão administrativa incompleta")
	}
	if !p.CanRevert {
		return nil, nil, errors.New("configuração administrativa alterada; preserve-a antes de reverter")
	}
	return []plan.UndoValue{{Kind: "admin", Key: e.Cat.Item(item).AdminOp, Current: p.Current, Target: p.Target, ChangeID: strings.Join(ids, ","), Mode: "admin", StateToken: p.StateToken}}, options, nil
}

func (e *Engine) validateAdminUndo(ctx context.Context, s plan.Step, admin helper.Admin) error {
	values, _, err := e.adminUndoValues(ctx, s.Item, admin)
	if err != nil {
		return err
	}
	a, _ := json.Marshal(values)
	b, _ := json.Marshal(s.Undo)
	if string(a) != string(b) {
		return errors.New("o estado ou alvo administrativo mudou desde a revisão; reveja o plano")
	}
	return nil
}
