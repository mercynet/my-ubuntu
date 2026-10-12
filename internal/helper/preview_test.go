package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/sys"
)

func TestAdministrativePreviewUsesRootSnapshotAndDoesNotMutate(t *testing.T) {
	c, _ := catalog.Load()
	r := sys.NewSim()
	r.Delay = 0
	e := &Env{Cat: c, R: r, Root: t.TempDir(), UID: "1000"}
	path := e.path(dropins["journald"].Path)
	os.MkdirAll(filepath.Dir(path), 0755)
	original := []byte("[Journal]\nSystemMaxUse=256M\n")
	os.WriteFile(path, original, 0644)
	options := map[string][]string{"value": {"1G"}}
	if err := e.Execute(context.Background(), Request{Op: "admin", Item: "journald-limit", Options: options}, func(string) {}); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	request := Request{Op: "revert-preview", Item: "journald-limit", Options: options}
	p, err := e.PreviewRevert(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if p.Current != string(before) || p.Target != string(original) || !p.CanRevert {
		t.Fatalf("incorrect trusted preview: %+v", p)
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("preview mutated the config")
	}
	var input, output bytes.Buffer
	json.NewEncoder(&input).Encode(request)
	if err := Serve(e, &input, &output); err != nil {
		t.Fatal(err)
	}
	var result Message
	if err := json.Unmarshal(output.Bytes(), &result); err != nil || !result.OK || result.Preview == nil || result.Preview.Target != string(original) {
		t.Fatalf("preview did not cross the helper protocol: %s, %v", output.String(), err)
	}
	os.WriteFile(path, []byte("changed independently\n"), 0644)
	p, err = e.PreviewRevert(context.Background(), request)
	if err != nil || p.CanRevert || p.Warning == "" {
		t.Fatal("later administrative edit must prevent automatic restoration")
	}
	if _, err := e.PreviewRevert(context.Background(), Request{Item: "git"}); err == nil {
		t.Fatal("preview must not accept arbitrary catalog paths")
	}
}
