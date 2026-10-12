package aptsrc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplacementMustCoverActiveRecipe(t *testing.T) {
	want := []byte("Types: deb\nURIs: https://repo.example/ubuntu\nSuites: stable\nComponents: main contrib\nArchitectures: amd64\n")
	cases := []struct {
		name, content, ext string
		match              bool
	}{
		{"disabled", "Enabled: no\nTypes: deb\nURIs: https://repo.example/ubuntu\nSuites: stable\nComponents: main contrib\n", ".sources", false},
		{"other suite", "deb https://repo.example/ubuntu testing main contrib", ".list", false},
		{"missing component", "deb https://repo.example/ubuntu stable main", ".list", false},
		{"other architecture", "deb [arch=arm64] https://repo.example/ubuntu stable main contrib", ".list", false},
		{"source packages only", "deb-src https://repo.example/ubuntu stable main contrib", ".list", false},
		{"restricted targets", "Types: deb\nURIs: https://repo.example/ubuntu\nSuites: stable\nComponents: main contrib\nTargets: Sources\n", ".sources", false},
		{"same recipe", "deb [arch=amd64 signed-by=/etc/apt/keyrings/vendor.gpg] https://repo.example/ubuntu/ stable main contrib", ".list", true},
		{"deb822 continuation", "Types: deb\nURIs: https://repo.example/ubuntu\nSuites: stable\nComponents: main\n contrib\n", ".sources", true},
		{"multiple stanzas", "Enabled: no\nTypes: deb\nURIs: https://unused.example\nSuites: stable\n\nTypes: deb\nURIs: https://repo.example/ubuntu\nSuites: stable\nComponents: main contrib\n", ".sources", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "etc/apt/sources.list.d")
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "vendor"+tc.ext), []byte(tc.content), 0644); err != nil {
				t.Fatal(err)
			}
			if got := len(ReplacementFiles(root, want, "")) > 0; got != tc.match {
				t.Fatalf("matched=%v, want %v", got, tc.match)
			}
		})
	}
}
