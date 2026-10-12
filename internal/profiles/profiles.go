// Package profiles embeds the built-in, editable starting profiles.
package profiles

import (
	"embed"
	"encoding/json"
	"sort"

	"meuubuntu/internal/store"
)

//go:embed *.json
var files embed.FS

func Builtin() []store.Profile {
	ents, _ := files.ReadDir(".")
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	var out []store.Profile
	for _, n := range names {
		b, _ := files.ReadFile(n)
		var p store.Profile
		if json.Unmarshal(b, &p) == nil {
			p.Builtin, p.Source = true, store.SourceBuiltin
			out = append(out, p)
		}
	}
	return out
}
