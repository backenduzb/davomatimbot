package report

import (
	"bytes"
	"fmt"
	"time"

	"scheduler/internal/models"
	"scheduler/internal/repository/attendance"

	"github.com/xuri/excelize/v2"
)

const (
	// FileName — guruhga yuboriladigan fayl nomi.
	FileName = "davomat.xlsx"

	// SheetName — varaq nomi.
	SheetName = "Davomat"

	// DateLayout — hisobotdagi sana ko'rinishi (2026/05/07).
	DateLayout = "2006/01/02"

	// HashTagLayout — hashtagdagi sana ko'rinishi (#2026_05_07).
	HashTagLayout = "2006_01_02"
)

// thinBorders barcha tomonlari (chap/o'ng/yuqori/past) yupka chiziqli
// chegara uslubini qaytaradi.
func thinBorders(color string) []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: color, Style: 1},
		{Type: "right", Color: color, Style: 1},
		{Type: "top", Color: color, Style: 1},
		{Type: "bottom", Color: color, Style: 1},
	}
}

// GenerateXLSX kunlik davomat qatorlaridan davomat.xlsx faylini (bitta
// varaq, bosh qator, filtr va muzlatilgan sarlavha bilan) generatsiya
// qiladi va xotiradagi []byte ko'rinishida qaytaradi.
func GenerateXLSX(rows []attendance.ReportRow) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	if err := f.SetSheetName("Sheet1", SheetName); err != nil {
		return nil, fmt.Errorf("varaqa nomini o'zgartirish: %w", err)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4472C4"}},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    thinBorders("4472C4"),
	})
	if err != nil {
		return nil, fmt.Errorf("sarlavha uslubini yaratish: %w", err)
	}

	leftStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border:    thinBorders("D9D9D9"),
	})
	if err != nil {
		return nil, fmt.Errorf("chap uslubini yaratish: %w", err)
	}

	centerStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    thinBorders("D9D9D9"),
	})
	if err != nil {
		return nil, fmt.Errorf("markaz uslubini yaratish: %w", err)
	}

	dateStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: strPtr("yyyy/mm/dd"),
		Alignment:    &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:       thinBorders("D9D9D9"),
	})
	if err != nil {
		return nil, fmt.Errorf("sana uslubini yaratish: %w", err)
	}

	// Sarlavha qatori.
	headers := []string{"Class", "Student", "Date", "Attendance", "Reason"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(SheetName, cell, header); err != nil {
			return nil, fmt.Errorf("sarlavha yozish: %w", err)
		}
		if err := f.SetCellStyle(SheetName, cell, cell, headerStyle); err != nil {
			return nil, fmt.Errorf("sarlavha uslubini qo'llash: %w", err)
		}
	}

	// Ma'lumot qatorlari.
	for i, row := range rows {
		r := i + 2
		if err := f.SetCellValue(SheetName, fmt.Sprintf("A%d", r), row.ClassName); err != nil {
			return nil, fmt.Errorf("qator yozish: %w", err)
		}
		if err := f.SetCellValue(SheetName, fmt.Sprintf("B%d", r), NormalizeName(row.Student)); err != nil {
			return nil, fmt.Errorf("qator yozish: %w", err)
		}
		if err := f.SetCellValue(SheetName, fmt.Sprintf("C%d", r), row.Date); err != nil {
			return nil, fmt.Errorf("qator yozish: %w", err)
		}
		if err := f.SetCellValue(SheetName, fmt.Sprintf("D%d", r), DisplayStatus(row.Status)); err != nil {
			return nil, fmt.Errorf("qator yozish: %w", err)
		}
		if err := f.SetCellValue(SheetName, fmt.Sprintf("E%d", r), row.Reason); err != nil {
			return nil, fmt.Errorf("qator yozish: %w", err)
		}
	}

	lastRow := len(rows) + 1

	// Uslublar (faqat bor ma'lumot qatorlari bo'ylab).
	if lastRow >= 2 {
		if err := f.SetCellStyle(SheetName, "A2", fmt.Sprintf("A%d", lastRow), leftStyle); err != nil {
			return nil, fmt.Errorf("uslubni qo'llash: %w", err)
		}
		if err := f.SetCellStyle(SheetName, "B2", fmt.Sprintf("B%d", lastRow), leftStyle); err != nil {
			return nil, fmt.Errorf("uslubni qo'llash: %w", err)
		}
		if err := f.SetCellStyle(SheetName, "C2", fmt.Sprintf("C%d", lastRow), dateStyle); err != nil {
			return nil, fmt.Errorf("uslubni qo'llash: %w", err)
		}
		if err := f.SetCellStyle(SheetName, "D2", fmt.Sprintf("D%d", lastRow), centerStyle); err != nil {
			return nil, fmt.Errorf("uslubni qo'llash: %w", err)
		}
		if err := f.SetCellStyle(SheetName, "E2", fmt.Sprintf("E%d", lastRow), leftStyle); err != nil {
			return nil, fmt.Errorf("uslubni qo'llash: %w", err)
		}
	}

	// Ustunlar kengligi.
	widths := [5]float64{12, 32, 12, 14, 24}
	cols := [5]string{"A", "B", "C", "D", "E"}
	for i := range cols {
		if err := f.SetColWidth(SheetName, cols[i], cols[i], widths[i]); err != nil {
			return nil, fmt.Errorf("ustun kengligi: %w", err)
		}
	}
	if err := f.SetRowHeight(SheetName, 1, 20); err != nil {
		return nil, fmt.Errorf("qator balandligi: %w", err)
	}

	// Muzlatilgan sarlavha qatori.
	if err := f.SetPanes(SheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return nil, fmt.Errorf("panes: %w", err)
	}

	// Excel filtri (Class/Student/Date/Attendance/Reason bo'yicha).
	if err := f.AutoFilter(SheetName, fmt.Sprintf("A1:E%d", lastRow), nil); err != nil {
		return nil, fmt.Errorf("avtofiltr: %w", err)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("faylni yozish: %w", err)
	}

	return buf.Bytes(), nil
}

// DisplayStatus DB'dagi davomat holatini hisobotda o'qiladigan ko'rinishiga
// o'tkazadi (masalan, "present" -> "Present").
func DisplayStatus(status string) string {
	switch status {
	case models.AttendancePresent:
		return "Present"
	case models.AttendanceAbsent:
		return "Absent"
	case models.AttendanceExcused:
		return "Excused"
	case models.AttendanceLate:
		return "Late"
	case models.AttendanceNotMarked:
		return "Not Marked"
	default:
		return status
	}
}

// CaptionFor Telegram xabarini (caption) sanani dinamik hisoblab tuzadi:
//
//	davomat.xlsx
//
//	📅 Sana: 2026/09/02
//	#2026_09_02
func CaptionFor(date time.Time) string {
	return fmt.Sprintf("%s\n\n📅 Sana: %s\n#%s", FileName, date.Format(DateLayout), date.Format(HashTagLayout))
}

func strPtr(s string) *string { return &s }
