// Package migrate handles the Firefox profile backup/migration. It never
// touches passwords individually, never overwrites an existing profile and
// refuses to run while Firefox is open.
package migrate

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"meuubuntu/internal/safefile"
)

type FFProfile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size string `json:"size"`
}

type FFRoot struct {
	Kind     string      `json:"kind"` // snap | deb | flatpak
	Label    string      `json:"label"`
	Root     string      `json:"root"`
	Profiles []FFProfile `json:"profiles"`
}

func roots(home string) []FFRoot {
	return []FFRoot{
		{Kind: "snap", Label: "Firefox Snap", Root: filepath.Join(home, "snap/firefox/common/.mozilla/firefox")},
		{Kind: "snap", Label: "Firefox Snap (XDG)", Root: filepath.Join(home, "snap/firefox/common/.config/mozilla/firefox")},
		{Kind: "deb", Label: "Firefox .deb/Mozilla", Root: filepath.Join(home, ".mozilla/firefox")},
		{Kind: "deb", Label: "Firefox .deb (XDG)", Root: filepath.Join(home, ".config/mozilla/firefox")},
		{Kind: "flatpak", Label: "Firefox Flatpak", Root: filepath.Join(home, ".var/app/org.mozilla.firefox/.mozilla/firefox")},
		{Kind: "flatpak", Label: "Firefox Flatpak (XDG)", Root: filepath.Join(home, ".var/app/org.mozilla.firefox/config/mozilla/firefox")},
	}
}

// Detect lists every profile root that exists, reading profiles.ini.
func Detect(home string) []FFRoot {
	var out []FFRoot
	for _, r := range roots(home) {
		b, err := os.ReadFile(filepath.Join(r.Root, "profiles.ini"))
		if err != nil {
			continue
		}
		var cur FFProfile
		rel := true
		flush := func() {
			if cur.Path != "" {
				if rel {
					cur.Path = filepath.Join(r.Root, cur.Path)
				}
				if st, err := os.Stat(cur.Path); err == nil && st.IsDir() {
					cur.Size = dirSize(cur.Path)
					r.Profiles = append(r.Profiles, cur)
				}
			}
			cur, rel = FFProfile{}, true
		}
		for _, l := range strings.Split(string(b), "\n") {
			l = strings.TrimSpace(l)
			switch {
			case strings.HasPrefix(l, "["):
				flush()
			case strings.HasPrefix(l, "Name="):
				cur.Name = strings.TrimPrefix(l, "Name=")
			case strings.HasPrefix(l, "Path="):
				cur.Path = strings.TrimPrefix(l, "Path=")
			case l == "IsRelative=0":
				rel = false
			}
		}
		flush()
		out = append(out, r)
	}
	return out
}

func dirSize(p string) string {
	var n int64
	filepath.WalkDir(p, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if i, e := d.Info(); e == nil {
				n += i.Size()
			}
		}
		return nil
	})
	return fmt.Sprintf("%.0f MB", float64(n)/1e6)
}

// Running reports whether any Firefox process is alive (any packaging).
func Running() bool {
	ents, _ := os.ReadDir("/proc")
	for _, e := range ents {
		b, err := os.ReadFile(filepath.Join("/proc", e.Name(), "comm"))
		if err != nil {
			continue
		}
		c := strings.TrimSpace(string(b))
		if c == "firefox" || c == "firefox-bin" || c == "GeckoMain" || c == "firefox-esr" {
			return true
		}
	}
	return false
}

// known profile roots only; prevents using the API to read arbitrary dirs.
func allowedProfile(home, p string) bool {
	p, err := filepath.Abs(p)
	if err != nil {
		return false
	}
	actual, err := filepath.EvalSymlinks(p)
	if err != nil {
		return false
	}
	realHome, err := filepath.EvalSymlinks(home)
	if err != nil {
		return false
	}
	for _, r := range roots(home) {
		rel, err := filepath.Rel(home, r.Root)
		if err != nil {
			continue
		}
		expected := filepath.Join(realHome, rel)
		if strings.HasPrefix(p, r.Root+string(os.PathSeparator)) && strings.HasPrefix(actual, expected+string(os.PathSeparator)) {
			fi, err := os.Stat(actual)
			return err == nil && fi.IsDir()
		}
	}
	return false
}

// Backup writes a tar.gz of the profile into ~/MeuUbuntu-backups (0700).
func Backup(home, profile string) (string, error) {
	if !allowedProfile(home, profile) {
		return "", errors.New("caminho de perfil não reconhecido")
	}
	if Running() {
		return "", errors.New("feche o Firefox antes do backup (ficheiros em uso ficam inconsistentes)")
	}
	profile, _ = filepath.EvalSymlinks(profile)
	dir := filepath.Join(home, "MeuUbuntu-backups")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, "firefox-"+filepath.Base(profile)+"-"+time.Now().Format("20060102-150405.000000000")+".tar.gz")
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	base := filepath.Dir(profile)
	err = filepath.WalkDir(profile, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&fs.ModeSymlink != 0 || d.Name() == "lock" || d.Name() == ".parentlock" {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return errors.New("perfil contém um ficheiro especial; backup recusado")
		}
		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		h.Name, _ = filepath.Rel(base, p)
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		if d.Type().IsRegular() {
			src, err := os.Open(p)
			if err != nil {
				return err
			}
			_, err = io.Copy(tw, src)
			src.Close()
			return err
		}
		return nil
	})
	if err == nil {
		err = tw.Close()
	}
	if err == nil {
		err = gz.Close()
	}
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(dst)
		return "", err
	}
	return dst, nil
}

// MigrateToFlatpak copies a profile into the Flatpak root as a NEW profile.
// Existing destination profiles are never overwritten; profiles.ini is backed up.
func MigrateToFlatpak(home, profile string) (string, error) {
	if !allowedProfile(home, profile) {
		return "", errors.New("caminho de perfil não reconhecido")
	}
	if strings.Contains(profile, "/.var/app/org.mozilla.firefox/") {
		return "", errors.New("o perfil de origem já é do Flatpak")
	}
	if Running() {
		return "", errors.New("feche todas as janelas do Firefox (Snap e Flatpak) antes de migrar")
	}
	profile, _ = filepath.EvalSymlinks(profile)
	// Prefer an existing Flatpak root; otherwise the XDG location used by new installs.
	var root string
	for _, r := range roots(home) {
		if r.Kind == "flatpak" {
			if _, err := os.Stat(filepath.Join(r.Root, "profiles.ini")); err == nil {
				root = r.Root
				break
			}
		}
	}
	if root == "" {
		return "", errors.New("abra o Firefox Flatpak uma vez (para criar o seu perfil) e feche-o; depois repita a migração")
	}
	actualRoot, err := filepath.EvalSymlinks(root)
	realHome, homeErr := filepath.EvalSymlinks(home)
	relRoot, relErr := filepath.Rel(home, root)
	if err != nil || homeErr != nil || relErr != nil || actualRoot != filepath.Join(realHome, relRoot) {
		return "", errors.New("destino do Firefox contém um redirecionamento; migração recusada")
	}
	root = actualRoot
	name := filepath.Base(profile) + "-migrado"
	dst := filepath.Join(root, name)
	if err := os.Mkdir(dst, 0o700); err != nil {
		return "", fmt.Errorf("destino %s já existe; nada foi sobrescrito", dst)
	}
	completed := false
	defer func() {
		if !completed {
			os.RemoveAll(dst)
		}
	}()
	if err := copyTree(profile, dst); err != nil {
		os.RemoveAll(dst)
		return "", err
	}
	ini := filepath.Join(root, "profiles.ini")
	fi, err := os.Lstat(ini)
	if err != nil || !fi.Mode().IsRegular() {
		return "", errors.New("profiles.ini não é um ficheiro regular")
	}
	b, err := os.ReadFile(ini)
	if err != nil {
		return "", err
	}
	if err := safefile.Write(ini+".meu-ubuntu-"+time.Now().Format("20060102-150405.000000000")+".bak", b, 0o600); err != nil {
		return "", err
	}
	n := 0
	for _, m := range regexp.MustCompile(`(?m)^\s*\[Profile([0-9]+)\]\s*$`).FindAllStringSubmatch(string(b), -1) {
		id, _ := strconv.Atoi(m[1])
		if id >= n {
			n = id + 1
		}
	}
	sec := fmt.Sprintf("\n[Profile%d]\nName=Migrado (%s)\nIsRelative=1\nPath=%s\n", n, filepath.Base(profile), name)
	if err := safefile.Write(ini, append(b, []byte(sec)...), fi.Mode().Perm()); err != nil {
		return "", err
	}
	completed = true
	return dst, nil
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		t := filepath.Join(dst, rel)
		if d.Name() == "lock" || d.Name() == ".parentlock" || d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(t, 0o700)
		}
		if !d.Type().IsRegular() {
			return errors.New("perfil contém um ficheiro especial; migração recusada")
		}
		in, err := os.Open(p)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile(t, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, in)
		if cerr := out.Close(); err == nil {
			err = cerr
		}
		return err
	})
}
