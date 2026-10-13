package engine

import (
	"context"
	"encoding/json"
	"errors"

	"meuubuntu/internal/plan"
)

// Pause holds between completed operations, without interrupting APT/dpkg.
func (e *Engine) Pause() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.running {
		return errors.New("não existe uma execução em curso")
	}
	if e.cancel.Load() {
		return errors.New("já foi pedido cancelamento")
	}
	if e.pause.Load() {
		return nil
	}
	e.pause.Store(true)
	e.stateLocked("pausing")
	return nil
}

func (e *Engine) Resume() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.running || !e.pause.Load() {
		return errors.New("a execução não está em pausa")
	}
	if e.cancel.Load() {
		return errors.New("já foi pedido cancelamento")
	}
	e.pause.Store(false)
	if e.runState == "paused" {
		e.stateLocked("resuming")
	} else {
		e.stateLocked("running")
	}
	return nil
}

// Call with mu held. A sequence also protects clients from delayed events.
func (e *Engine) stateLocked(state string) {
	e.runState = state
	e.stateVersion++
	e.Emit("run:state", map[string]any{"runId": e.runID, "state": state, "version": e.stateVersion})
}

func (e *Engine) validateRemainingAfterPause(ctx context.Context, original *plan.Plan, from int) error {
	e.Det.Refresh(ctx)
	fresh, err := plan.Build(ctx, e.Cat, e.Det, original.Selection, original.Reverts)
	if err == nil {
		err = e.ReviewUndo(ctx, fresh)
	}
	if err != nil || fresh.Blocked {
		return errors.New("o estado mudou durante a pausa; reveja o que falta antes de continuar")
	}
	e.ReviewRecovery(fresh)
	for _, expected := range original.Steps[from:] {
		if expected.Item == "" {
			continue
		}
		matched := false
		for _, actual := range fresh.Steps {
			if actual.Item != expected.Item || actual.Op != expected.Op {
				continue
			}
			a, b := expected, actual
			a.ID, b.ID = "", ""
			if a.Op == "user-script" {
				b.Expect, b.Detail = a.Expect, a.Detail
			}
			x, _ := json.Marshal(a)
			y, _ := json.Marshal(b)
			matched = string(x) == string(y)
			if matched {
				break
			}
		}
		if !matched {
			return errors.New("o estado mudou durante a pausa; reveja o que falta antes de continuar")
		}
	}
	return nil
}
