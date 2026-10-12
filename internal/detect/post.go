package detect

import (
	"os"
	"path/filepath"
	"strings"

	"meuubuntu/internal/catalog"
)

func (d *Detector) PackagesPresent(it *catalog.Item, sel catalog.Selection) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.pkgs == nil {
		return false
	}
	for _, p := range d.Cat.ExpandPackages(it, sel) {
		if !d.pkgs[p] {
			return false
		}
	}
	return true
}

// PostReady checks the state produced by the catalog's allowlisted post action.
func (d *Detector) PostReady(it *catalog.Item, sel catalog.Selection) bool {
	for _, action := range it.Post {
		if action != "xdebug-off" {
			return false
		}
		if len(sel.Options["versions"]) == 0 {
			return false
		}
		for _, version := range sel.Options["versions"] {
			if sim, ok := d.R.(interface{ XdebugConfigured(string) bool }); ok {
				if !sim.XdebugConfigured(version) {
					return false
				}
				continue
			}
			root := filepath.Join(d.Root, "/etc/php", version)
			sapis, err := os.ReadDir(root)
			if err != nil {
				return false
			}
			found := false
			for _, sapi := range sapis {
				if !sapi.IsDir() || sapi.Name() == "mods-available" {
					continue
				}
				entries, err := os.ReadDir(filepath.Join(root, sapi.Name(), "conf.d"))
				if err != nil {
					return false
				}
				found = true
				for _, entry := range entries {
					if strings.Contains(strings.ToLower(entry.Name()), "xdebug") && strings.HasSuffix(entry.Name(), ".ini") {
						return false
					}
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}
