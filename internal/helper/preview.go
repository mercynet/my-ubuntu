package helper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"meuubuntu/internal/catalog"
)

type RevertPreview struct {
	Current    string `json:"current"`
	Target     string `json:"target"`
	CanRevert  bool   `json:"canRevert"`
	Warning    string `json:"warning,omitempty"`
	StateToken string `json:"stateToken,omitempty"`
}

// PreviewRevert reads only fixed catalog paths and trusted administrative snapshots.
func (e *Env) PreviewRevert(ctx context.Context, req Request) (*RevertPreview, error) {
	it := e.Cat.Item(req.Item)
	if it == nil || it.Method != catalog.MethodAdmin {
		return nil, errors.New("item administrativo inválido")
	}
	sel, err := e.Cat.ValidateSelection(it.ID, catalog.Selection{Options: req.Options})
	if err != nil {
		return nil, err
	}
	if it.AdminOp == "docker-group" {
		u, err := e.username()
		if err != nil {
			return nil, err
		}
		out, err := e.R.Output(ctx, "id", "-nG", u)
		if err != nil {
			return nil, err
		}
		member := false
		for _, group := range strings.Fields(out) {
			member = member || group == "docker"
		}
		current := "Fora do grupo docker"
		if member {
			current = "Membro do grupo docker"
		}
		p := &RevertPreview{Current: current, Target: "Fora do grupo docker", CanRevert: member || req.Batch}
		p.StateToken = previewFingerprint(p, nil, e.UID, req.Batch)
		return p, nil
	}
	op := it.AdminOp
	if op != "inotify" && op != "journald" {
		return nil, errors.New("operação sem reversão automática")
	}
	b, _, err := regularContents(e.snapshotPath(op))
	if err != nil {
		return nil, fmt.Errorf("sem configuração anterior verificável: %w", err)
	}
	var snapshot dropinSnapshot
	if err := json.Unmarshal(b, &snapshot); err != nil {
		return nil, err
	}
	if !bytes.Equal(snapshot.After, []byte(fmt.Sprintf(dropins[op].Template, sel.First("value")))) {
		return nil, errors.New("existe uma alteração administrativa posterior; reveja e reverta-a primeiro")
	}
	cur, _, err := regularContents(e.path(dropins[op].Path))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	p := &RevertPreview{Current: string(cur), Target: string(snapshot.Before)}
	target := snapshot
	if req.Batch {
		for target.Previous != nil {
			target = *target.Previous
		}
		p.Target = string(target.Before)
	}
	if os.IsNotExist(err) {
		p.Current = "(ficheiro ausente)"
	}
	if !target.HadValue {
		p.Target = "(remover ficheiro criado pela aplicação)"
	}
	p.CanRevert = (err == nil && bytes.Equal(cur, snapshot.After)) || (snapshot.HadValue && err == nil && bytes.Equal(cur, snapshot.Before)) || (!snapshot.HadValue && os.IsNotExist(err))
	if !p.CanRevert {
		p.Warning = "O ficheiro mudou depois da aplicação. Preserve essa alteração antes de reverter."
	}
	if op == "inotify" {
		b, err := os.ReadFile(e.path("/proc/sys/fs/inotify/max_user_watches"))
		if err != nil {
			return nil, err
		}
		p.Current += "\nLimite em vigor: " + strings.TrimSpace(string(b))
		p.Target += "\nLimite a restaurar: " + target.RuntimeBefore
	}
	p.StateToken = previewFingerprint(p, b, e.UID, req.Batch)
	return p, nil
}

func previewFingerprint(p *RevertPreview, snapshot []byte, uid string, batch bool) string {
	data, _ := json.Marshal(struct {
		Preview  *RevertPreview
		Snapshot []byte
		UID      string
		Batch    bool
	}{p, snapshot, uid, batch})
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
