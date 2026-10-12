package diag

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixtureExport() (Report, Metadata) {
	r := simReport()
	r.CollectedAt = "2026-10-12T12:00:00Z"
	r.Facts = append(r.Facts, Fact{"Teste Unicode", "João · ação · configuração · Ω"}, Fact{"Formula", "=1+1"})
	r.Facts = append(r.Facts, Fact{"Emoji", "🚀"})
	return r, Metadata{Version: "1.0.0-rc.1", System: "Ubuntu 26.04 LTS", Architecture: "amd64", Session: "wayland", Mode: "simulação"}
}

func TestCSVSectionsUnicodeAndFormulaProtection(t *testing.T) {
	r, m := fixtureExport()
	b, err := Encode(r, m, "csv", []string{"facts"})
	if err != nil {
		t.Fatal(err)
	}
	reader := csv.NewReader(strings.NewReader(strings.TrimPrefix(string(b), "\ufeff")))
	reader.Comma = ';'
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	foundUnicode, foundFormula := false, false
	for _, row := range rows {
		if row[0] == "Arranque" || row[0] == "Discos" {
			t.Fatal("unchecked sections were exported")
		}
		foundUnicode = foundUnicode || row[3] == "João · ação · configuração · Ω"
		foundFormula = foundFormula || row[3] == "'=1+1"
	}
	if !foundUnicode || !foundFormula {
		t.Fatalf("Unicode/formula handling missing: %s", b)
	}
}

func TestXLSXContainsOnlyTextCellsAndWellFormedXML(t *testing.T) {
	r, m := fixtureExport()
	b, err := Encode(r, m, "xlsx", []string{"facts", "disks", "startup", "suggestions"})
	if err != nil {
		t.Fatal(err)
	}
	z, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range z.File {
		rd, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(rd)
		rd.Close()
		if err != nil {
			t.Fatal(err)
		}
		d := xml.NewDecoder(bytes.NewReader(data))
		for {
			_, err := d.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("invalid %s: %v", f.Name, err)
			}
		}
		if f.Name == "xl/worksheets/sheet1.xml" {
			found = true
			if bytes.Contains(data, []byte("<f>")) || !bytes.Contains(data, []byte("=1+1")) || !bytes.Contains(data, []byte("ação")) {
				t.Fatal("cells must preserve Unicode as text, never formulas")
			}
		}
	}
	if !found {
		t.Fatal("worksheet missing")
	}
}

func TestPDFAndIndependentReaderFixtures(t *testing.T) {
	r, m := fixtureExport()
	for _, format := range []string{"pdf", "csv", "xlsx"} {
		b, err := Encode(r, m, format, []string{"facts", "disks", "startup", "suggestions"})
		if err != nil {
			t.Fatal(err)
		}
		if format == "pdf" && (!bytes.HasPrefix(b, []byte("%PDF-")) || !bytes.Contains(b, []byte("%%EOF"))) {
			t.Fatal("invalid PDF framing")
		}
		if target := os.Getenv("MEU_UBUNTU_EXPORT_FIXTURES"); target != "" {
			if err := os.MkdirAll(target, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(target, "diagnostico."+format), b, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestExportRejectsUnknownFormatAndSections(t *testing.T) {
	r, m := fixtureExport()
	for _, c := range []struct {
		format   string
		sections []string
	}{{"html", []string{"facts"}}, {"csv", nil}, {"pdf", []string{"logs"}}} {
		if _, err := Encode(r, m, c.format, c.sections); err == nil {
			t.Fatalf("must reject %+v", c)
		}
	}
}
