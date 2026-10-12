package engine

import (
	"fmt"
	"meuubuntu/internal/plan"
)

// ReviewRecovery adds read-only retries of checks that did not complete.
// The resulting steps participate in the reviewed plan fingerprint.
func (e *Engine) ReviewRecovery(p *plan.Plan) {
	seen := map[string]bool{}
	for _, s := range p.Steps {
		if s.Item != "" {
			seen[s.Item] = true
		}
	}
	for _, run := range e.Store.RecoveryRuns() {
		for i := len(run.Steps) - 1; i >= 0; i-- {
			s := run.Steps[i]
			if s.Status == StPending || s.Status == StRunning || s.Status == StCancelled || s.Status == StSkipped {
				continue
			}
			if seen[s.Item] {
				continue
			}
			seen[s.Item] = true
			sel, selected := p.Selection[s.Item]
			it := e.Cat.Item(s.Item)
			if !selected || it == nil || len(it.Verify.Cmd) == 0 || (s.Status != StIntervention && s.Status != StFailed) {
				continue
			}
			check := plan.Step{Item: it.ID, Title: "Repetir verificação: " + it.Name, Kind: "user", Op: "verify", Options: sel.Options, Category: it.Category, Detail: []string{"A verificação anterior ficou por concluir; não reinstala a aplicação"}}
			position := len(p.Steps)
			for j, step := range p.Steps {
				if e.dependsOn(step.Item, it.ID) {
					position = j
					break
				}
			}
			p.Steps = append(p.Steps, plan.Step{})
			copy(p.Steps[position+1:], p.Steps[position:])
			p.Steps[position] = check
			p.Warnings = append(p.Warnings, "Será repetida a verificação pendente de "+it.Name)
			var kept []plan.Skipped
			for _, skip := range p.Skipped {
				if skip.Item != it.ID {
					kept = append(kept, skip)
				}
			}
			p.Skipped = kept
		}
	}
	for i := range p.Steps {
		p.Steps[i].ID = fmt.Sprintf("s%02d", i+1)
	}
}

func (e *Engine) dependsOn(item, requirement string) bool {
	it := e.Cat.Item(item)
	if it == nil {
		return false
	}
	for _, r := range it.Requires {
		if r == requirement || e.dependsOn(r, requirement) {
			return true
		}
	}
	return false
}
