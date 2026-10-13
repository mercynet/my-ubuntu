package diag

import (
	"archive/zip"
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-pdf/fpdf"
)

//go:embed assets/DejaVuSans.ttf assets/DejaVu-LICENSE.txt assets/FPDF-LICENSE.txt
var exportAssets embed.FS

type Metadata struct{ Version, System, Architecture, Session, Mode string }
type ExportLicense struct{ Name, License, URL, Text string }

func ExportLicenses() []ExportLicense {
	font, _ := exportAssets.ReadFile("assets/DejaVu-LICENSE.txt")
	pdf, _ := exportAssets.ReadFile("assets/FPDF-LICENSE.txt")
	return []ExportLicense{{"DejaVu Sans", "Bitstream Vera / DejaVu", "https://dejavu-fonts.github.io/License.html", string(font)}, {"go-pdf/fpdf", "MIT", "https://github.com/go-pdf/fpdf/tree/v0.9.0", string(pdf)}}
}

func exportRows(rep Report, meta Metadata, sections []string) ([][]string, error) {
	chosen := map[string]bool{}
	for _, s := range sections {
		if s != "facts" && s != "disks" && s != "startup" && s != "suggestions" {
			return nil, errors.New("secção de diagnóstico inválida")
		}
		chosen[s] = true
	}
	if len(chosen) == 0 {
		return nil, errors.New("escolha pelo menos uma secção")
	}
	rows := [][]string{{"Secção", "Registo", "Campo", "Valor"}}
	add := func(section, record, key, value string) { rows = append(rows, []string{section, record, key, value}) }
	for _, v := range [][2]string{{"Aplicação", "Meu Ubuntu " + meta.Version}, {"Recolhido em", rep.CollectedAt}, {"Sistema", meta.System}, {"Arquitetura", meta.Architecture}, {"Sessão", meta.Session}, {"Modo", meta.Mode}, {"Bateria detetada", strconv.FormatBool(rep.Battery)}} {
		add("Relatório", "Informações", v[0], v[1])
	}
	if chosen["facts"] {
		for i, f := range rep.Facts {
			add("Sistema", strconv.Itoa(i+1), f.Label, f.Value)
		}
		if len(rep.Facts) == 0 {
			add("Sistema", "", "Estado", "Sem dados disponíveis")
		}
	}
	if chosen["disks"] {
		for i, d := range rep.Disks {
			record := fmt.Sprintf("%d · %s", i+1, d.Name)
			for _, v := range [][2]string{{"Nome", d.Name}, {"Modelo", d.Model}, {"Capacidade", d.Size}, {"SSD", strconv.FormatBool(d.SSD)}, {"Suporta TRIM", strconv.FormatBool(d.Discard)}, {"Montagem", d.Mount}, {"Livre (%)", strconv.Itoa(d.FreePct)}, {"Espaço livre", d.FreeHint}} {
				add("Discos", record, v[0], v[1])
			}
		}
		if len(rep.Disks) == 0 {
			add("Discos", "", "Estado", "Nenhum disco detetado")
		}
	}
	if chosen["startup"] {
		for i, s := range rep.Startup {
			record := fmt.Sprintf("%d · %s", i+1, s.Name)
			for _, v := range [][2]string{{"Nome", s.Name}, {"Âmbito", s.Scope}, {"Estado", s.Status}, {"Detalhe", s.Detail}, {"Origem", s.Path}, {"Executável", s.Command}} {
				add("Arranque", record, v[0], v[1])
			}
		}
		if len(rep.Startup) == 0 {
			add("Arranque", "", "Estado", "Nenhuma entrada detetada")
		}
	}
	if chosen["suggestions"] {
		for i, s := range rep.Suggestions {
			for _, v := range [][2]string{{"Título", s.Title}, {"Detalhe", s.Detail}, {"Evidência", s.Evidence}, {"Item do catálogo", s.Item}, {"Nível", s.Level}} {
				add("Sugestões", strconv.Itoa(i+1), v[0], v[1])
			}
		}
		if len(rep.Suggestions) == 0 {
			add("Sugestões", "", "Estado", "Sem sugestões")
		}
	}
	if len(rows) > 10000 {
		return nil, errors.New("relatório demasiado extenso; exporte menos secções de cada vez")
	}
	for _, row := range rows {
		for i, s := range row {
			row[i] = cleanText(s)
			if len([]rune(row[i])) > 32767 {
				return nil, errors.New("campo demasiado extenso para o relatório")
			}
		}
	}
	return rows, nil
}

func cleanText(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' || r == 0xfffe || r == 0xffff {
			return '\ufffd'
		}
		return r
	}, s)
}

// Encode generates all formats from the same immutable collected snapshot.
func Encode(rep Report, meta Metadata, format string, sections []string) ([]byte, error) {
	rows, err := exportRows(rep, meta, sections)
	if err != nil {
		return nil, err
	}
	switch format {
	case "csv":
		return csvReport(rows)
	case "xlsx":
		return xlsxReport(rows)
	case "pdf":
		return pdfReport(rows, meta)
	default:
		return nil, errors.New("formato inválido; escolha PDF, CSV ou XLSX")
	}
}

func csvReport(rows [][]string) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(&b)
	w.Comma = ';'
	w.UseCRLF = true
	for _, row := range rows {
		cells := append([]string(nil), row...)
		for i, s := range cells {
			first := strings.TrimLeftFunc(s, unicode.IsSpace)
			if first != "" && strings.ContainsRune("=+-@", []rune(first)[0]) {
				cells[i] = "'" + s
			}
		}
		if err := w.Write(cells); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return b.Bytes(), w.Error()
}

func xmlText(s string) string { var b bytes.Buffer; xml.EscapeText(&b, []byte(s)); return b.String() }
func xlsxReport(rows [][]string) ([]byte, error) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	parts := map[string]string{
		"[Content_Types].xml":        `<?xml version="1.0" encoding="UTF-8"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/><Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/></Types>`,
		"_rels/.rels":                `<?xml version="1.0" encoding="UTF-8"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`,
		"xl/workbook.xml":            `<?xml version="1.0" encoding="UTF-8"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Diagnóstico" sheetId="1" r:id="rId1"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0" encoding="UTF-8"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/></Relationships>`,
		"xl/styles.xml":              `<?xml version="1.0" encoding="UTF-8"?><styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fonts count="1"><font><sz val="11"/><name val="Calibri"/></font></fonts><fills count="2"><fill><patternFill patternType="none"/></fill><fill><patternFill patternType="gray125"/></fill></fills><borders count="1"><border><left/><right/><top/><bottom/><diagonal/></border></borders><cellStyleXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"/></cellStyleXfs><cellXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0" applyAlignment="1"><alignment vertical="top" wrapText="1"/></xf></cellXfs><cellStyles count="1"><cellStyle name="Normal" xfId="0" builtinId="0"/></cellStyles></styleSheet>`,
	}
	var sheet strings.Builder
	sheet.WriteString(`<?xml version="1.0" encoding="UTF-8"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetViews><sheetView workbookViewId="0"><pane ySplit="1" topLeftCell="A2" activePane="bottomLeft" state="frozen"/></sheetView></sheetViews><cols><col min="1" max="2" width="24" customWidth="1"/><col min="3" max="3" width="30" customWidth="1"/><col min="4" max="4" width="90" customWidth="1"/></cols><sheetData>`)
	for i, row := range rows {
		fmt.Fprintf(&sheet, `<row r="%d">`, i+1)
		for j, s := range row {
			fmt.Fprintf(&sheet, `<c r="%c%d" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, 'A'+j, i+1, xmlText(s))
		}
		sheet.WriteString(`</row>`)
	}
	fmt.Fprintf(&sheet, `</sheetData><autoFilter ref="A1:D%d"/></worksheet>`, len(rows))
	parts["xl/worksheets/sheet1.xml"] = sheet.String()
	for _, name := range []string{"[Content_Types].xml", "_rels/.rels", "xl/workbook.xml", "xl/_rels/workbook.xml.rels", "xl/styles.xml", "xl/worksheets/sheet1.xml"} {
		w, err := z.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err = w.Write([]byte(parts[name])); err != nil {
			return nil, err
		}
	}
	if err := z.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func pdfReport(rows [][]string, meta Metadata) ([]byte, error) {
	f := fpdf.New("P", "mm", "A4", "")
	font, err := exportAssets.ReadFile("assets/DejaVuSans.ttf")
	if err != nil {
		return nil, err
	}
	f.AddUTF8FontFromBytes("DejaVu", "", font)
	f.SetMargins(15, 15, 15)
	f.SetAutoPageBreak(true, 18)
	f.SetTitle("Diagnóstico — Meu Ubuntu", true)
	f.SetCreator("Meu Ubuntu "+meta.Version, true)
	f.SetHeaderFunc(func() {
		// Automatic breaks can retain the previous cell's horizontal position.
		f.SetXY(15, 15)
		f.SetFont("DejaVu", "", 15)
		f.SetTextColor(200, 65, 15)
		f.CellFormat(180, 10, "Diagnóstico — Meu Ubuntu", "", 1, "L", false, 0, "")
		f.SetTextColor(35, 35, 35)
		f.Ln(2)
	})
	f.SetFooterFunc(func() {
		f.SetY(-12)
		f.SetFont("DejaVu", "", 8)
		f.CellFormat(180, 5, fmt.Sprintf("Meu Ubuntu %s · Página %d", meta.Version, f.PageNo()), "", 0, "R", false, 0, "")
	})
	f.AddPage()
	// This font backend only indexes the BMP. Keep supplementary codepoints
	// readable instead of allowing an out-of-range panic; CSV/XLSX stay intact.
	pdfText := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			if r > 0xffff {
				fmt.Fprintf(&b, "[U+%X]", r)
			} else {
				b.WriteRune(r)
			}
		}
		return b.String()
	}
	section, record := "", ""
	for _, row := range rows[1:] {
		if row[0] != section {
			// Keep each section title with the start of its first record.
			if f.GetY()+26 > 279 {
				f.AddPage()
			}
			section, record = row[0], ""
			f.Ln(3)
			f.SetFont("DejaVu", "", 12)
			f.SetFillColor(241, 241, 244)
			f.CellFormat(180, 8, pdfText(section), "", 1, "L", true, 0, "")
		}
		if row[1] != "" && row[1] != record {
			record = row[1]
			f.SetFont("DejaVu", "", 10)
			f.MultiCell(180, 6, pdfText(record), "", "L", false)
		}
		f.SetFont("DejaVu", "", 9)
		f.MultiCell(180, 5, pdfText(row[2]+": "+row[3]), "", "L", false)
	}
	var b bytes.Buffer
	if err := f.Output(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
