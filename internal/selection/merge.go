// Package selection combines profiles into the current selection. Items are
// deduplicated; anything that would silently override an existing choice is
// returned as a conflict for the user to resolve (never decided by order).
package selection

import (
	"sort"

	"meuubuntu/internal/catalog"
)

const Manual = "manual"

type Conflict struct {
	Kind          string            `json:"kind"` // options | group | exclusive
	Item          string            `json:"item"` // item currently selected
	Name          string            `json:"name"`
	Current       catalog.Selection `json:"current"`
	CurrentOrigin []string          `json:"currentOrigin"`
	IncomingItem  string            `json:"incomingItem"`
	IncomingName  string            `json:"incomingName"`
	Incoming      catalog.Selection `json:"incoming"`
	IncomingLabel string            `json:"incomingLabel"`
	Reason        string            `json:"reason"`
}

type Result struct {
	Selections map[string]catalog.Selection `json:"selections"`
	Origins    map[string][]string          `json:"origins"`
	Added      []string                     `json:"added"`
	Conflicts  []Conflict                   `json:"conflicts"`
}

func equal(a, b catalog.Selection) bool {
	if len(a.Options) != len(b.Options) {
		return false
	}
	for k, va := range a.Options {
		vb := b.Options[k]
		if len(va) != len(vb) {
			return false
		}
		sa, sb := append([]string(nil), va...), append([]string(nil), vb...)
		sort.Strings(sa)
		sort.Strings(sb)
		for i := range sa {
			if sa[i] != sb[i] {
				return false
			}
		}
	}
	return true
}

func addOrigin(o []string, label string) []string {
	for _, x := range o {
		if x == label {
			return o
		}
	}
	return append(o, label)
}

// Merge adds incoming (labelled, e.g. "perfil:Desenvolvimento PHP") to cur.
// Invalid incoming selections are skipped by the caller's validation first.
func Merge(cat *catalog.Catalog, cur map[string]catalog.Selection, origins map[string][]string, incoming map[string]catalog.Selection, label string) Result {
	res := Result{Selections: map[string]catalog.Selection{}, Origins: map[string][]string{}}
	for k, v := range cur {
		res.Selections[k] = v
	}
	for k, v := range origins {
		if _, ok := cur[k]; ok {
			res.Origins[k] = append([]string(nil), v...)
		}
	}
	ids := make([]string, 0, len(incoming))
	for id := range incoming {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		it := cat.Item(id)
		if it == nil {
			continue
		}
		in, err := cat.ValidateSelection(id, incoming[id])
		if err != nil {
			continue
		}
		if have, ok := res.Selections[id]; ok {
			have, _ = cat.ValidateSelection(id, have)
			if equal(have, in) {
				res.Origins[id] = addOrigin(res.Origins[id], label) // deduplicated
				continue
			}
			res.Conflicts = append(res.Conflicts, Conflict{Kind: "options", Item: id, Name: it.Name, Current: have, CurrentOrigin: res.Origins[id],
				IncomingItem: id, IncomingName: it.Name, Incoming: in, IncomingLabel: label, Reason: "Opções diferentes para o mesmo item"})
			continue
		}
		if c := blocking(cat, res, it, in, label); c != nil {
			res.Conflicts = append(res.Conflicts, *c)
			continue
		}
		res.Selections[id] = in
		res.Origins[id] = []string{label}
		res.Added = append(res.Added, id)
	}
	return res
}

// blocking finds a selected item that cannot coexist with it (same group or hard conflict).
func blocking(cat *catalog.Catalog, res Result, it *catalog.Item, in catalog.Selection, label string) *Conflict {
	ids := make([]string, 0, len(res.Selections))
	for id := range res.Selections {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		other := cat.Item(id)
		if other == nil {
			continue
		}
		kind, reason := "", ""
		if it.Group != "" && other.Group == it.Group {
			kind, reason = "group", "Alternativas exclusivas: só uma pode ficar"
		}
		for _, c := range it.Conflicts {
			if c.ID == id && c.Hard {
				kind, reason = "exclusive", c.Reason
			}
		}
		for _, c := range other.Conflicts {
			if c.ID == it.ID && c.Hard {
				kind, reason = "exclusive", c.Reason
			}
		}
		if kind != "" {
			return &Conflict{Kind: kind, Item: id, Name: other.Name, Current: res.Selections[id], CurrentOrigin: res.Origins[id],
				IncomingItem: it.ID, IncomingName: it.Name, Incoming: in, IncomingLabel: label, Reason: reason}
		}
	}
	return nil
}
