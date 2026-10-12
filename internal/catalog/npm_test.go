package catalog

import "testing"

func TestNPMRejectsCommandArgumentsAndUnpinnedSources(t *testing.T) {
	for _, npm := range []*NPM{
		nil,
		{Package: "--prefix=/tmp", Version: "1.2.3", Binary: "tool"},
		{Package: "https://example.com/tool.tgz", Version: "1.2.3", Binary: "tool"},
		{Package: "@openai/codex;touch /tmp/x", Version: "1.2.3", Binary: "codex"},
		{Package: "@openai/codex", Version: "latest", Binary: "codex"},
		{Package: "@openai/codex", Version: "1.2.3 --force", Binary: "codex"},
		{Package: "@openai/codex", Version: "1.2.3", Binary: "../codex"},
	} {
		if npm.Valid() {
			t.Fatalf("unsafe metadata accepted: %#v", npm)
		}
	}
	cat, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, item := range cat.Items {
		if item.Method != MethodNPMUser {
			continue
		}
		count++
		if !item.NPM.Valid() || item.Confirm == "" || len(item.Requires) == 0 || item.NPM.Binary == "" {
			t.Fatalf("incomplete AI agent recipe: %s", item.ID)
		}
	}
	if count != 6 {
		t.Fatalf("expected six curated agents, got %d", count)
	}
}
