package report

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
	"time"

	"scheduler/internal/repository/attendance"

	"github.com/xuri/excelize/v2"
)

func sampleRows() []attendance.ReportRow {
	date := time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC)
	return []attendance.ReportRow{
		{ClassName: "10-A", Student: "Aliyev Vali", Date: date, Status: "present", Reason: ""},
		{ClassName: "10-A", Student: "Karimova Dilorom", Date: date, Status: "absent", Reason: "Illness"},
		{ClassName: "10-B", Student: "Rahimov Jasur", Date: date, Status: "absent", Reason: "Unexcused"},
	}
}

func TestGenerateXLSXContents(t *testing.T) {
	data, err := GenerateXLSX(sampleRows())
	if err != nil {
		t.Fatalf("GenerateXLSX: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("bo'sh fayl qaytarildi")
	}

	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("faylni qayta ochib bo'lmadi: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) != 1 || sheets[0] != SheetName {
		t.Fatalf("varaq ro'yxati = %v, want [Davomat]", sheets)
	}

	rows, err := f.GetRows(SheetName)
	if err != nil {
		t.Fatalf("GetRows: %v", err)
	}
	if len(rows) != 4 {
		t.Fatalf("qatorlar soni = %d, want 4 (sarlavha + 3)", len(rows))
	}

	wantHeaders := []string{"Class", "Student", "Date", "Attendance", "Reason"}
	for i, want := range wantHeaders {
		if rows[0][i] != want {
			t.Errorf("sarlavha[%d] = %q, want %q", i, rows[0][i], want)
		}
	}

	if rows[1][0] != "10-A" || rows[1][1] != "Aliyev Vali" || rows[1][3] != "Present" {
		t.Errorf("1-qator = %v", rows[1])
	}
	if rows[2][0] != "10-A" || rows[2][1] != "Karimova Dilorom" || rows[2][3] != "Absent" || rows[2][4] != "Illness" {
		t.Errorf("2-qator = %v", rows[2])
	}
	if rows[3][0] != "10-B" || rows[3][3] != "Absent" || rows[3][4] != "Unexcused" {
		t.Errorf("3-qator = %v", rows[3])
	}

	// Sana hujayrasida to'g'ri formatda ko'rinishi kerak.
	dateCell, err := f.GetCellValue(SheetName, "C2")
	if err != nil {
		t.Fatalf("GetCellValue: %v", err)
	}
	if dateCell != "2026/05/07" {
		t.Errorf("sana hujayra = %q, want 2026/05/07", dateCell)
	}

	// Muzlatilgan sarlavha.
	panes, err := f.GetPanes(SheetName)
	if err != nil {
		t.Fatalf("GetPanes: %v", err)
	}
	if panes == nil || !panes.Freeze || panes.YSplit != 1 || panes.TopLeftCell != "A2" {
		t.Errorf("panes = %+v, want frozen top row", panes)
	}
}

func TestGenerateXLSXAutoFilter(t *testing.T) {
	data, err := GenerateXLSX(sampleRows())
	if err != nil {
		t.Fatalf("GenerateXLSX: %v", err)
	}

	xml, err := readSheetXML(data)
	if err != nil {
		t.Fatalf("readSheetXML: %v", err)
	}
	if !strings.Contains(xml, `<autoFilter ref="A1:E4"`) {
		t.Error("avtofiltr (A1:E4) topilmadi — Excel filtri ishlamaydi")
	}
}

func TestGenerateXLSXEmpty(t *testing.T) {
	data, err := GenerateXLSX(nil)
	if err != nil {
		t.Fatalf("GenerateXLSX (bo'sh): %v", err)
	}

	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("faylni ochib bo'lmadi: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows(SheetName)
	if err != nil {
		t.Fatalf("GetRows: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("bo'sh hisobotda faqat sarlavha bo'lishi kerak, qatorlar = %d", len(rows))
	}

	xml, err := readSheetXML(data)
	if err != nil {
		t.Fatalf("readSheetXML: %v", err)
	}
	if !strings.Contains(xml, `<autoFilter ref="A1:E1"`) {
		t.Error("bo'sh hisobotda avtofiltr (A1:E1) topilmadi")
	}
}

func TestDisplayStatus(t *testing.T) {
	tests := map[string]string{
		"present":    "Present",
		"absent":     "Absent",
		"excused":    "Excused",
		"late":       "Late",
		"not_marked": "Not Marked",
		"unknown":    "unknown",
	}
	for in, want := range tests {
		if got := DisplayStatus(in); got != want {
			t.Errorf("DisplayStatus(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCaptionFor(t *testing.T) {
	date := time.Date(2026, 9, 2, 16, 0, 0, 0, time.UTC)
	want := "davomat.xlsx\n\n📅 Sana: 2026/09/02\n#2026_09_02"
	if got := CaptionFor(date); got != want {
		t.Errorf("CaptionFor = %q, want %q", got, want)
	}

	date2 := time.Date(2026, 12, 31, 16, 0, 0, 0, time.UTC)
	want2 := "davomat.xlsx\n\n📅 Sana: 2026/12/31\n#2026_12_31"
	if got := CaptionFor(date2); got != want2 {
		t.Errorf("CaptionFor = %q, want %q", got, want2)
	}
}

// readSheetXML xlsx (zip) ichidagi asosiy varaq XML'ini qaytaradi.
func readSheetXML(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	for _, entry := range zr.File {
		if strings.HasPrefix(entry.Name, "xl/worksheets/sheet") && strings.HasSuffix(entry.Name, ".xml") {
			rc, err := entry.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(rc); err != nil {
				return "", err
			}
			return buf.String(), nil
		}
	}
	return "", nil
}
