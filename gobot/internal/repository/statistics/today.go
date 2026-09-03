package statistics

import (
	"fmt"
	"math"
	"strings"
	"time"

	"bot/internal/database"
	"bot/internal/models"
)

type AbsentStudent struct {
	Name      string
	Reason    string
	ClassName string
}

// ClassBreakdown — bitta sinf kesimidagi davomat holati.
type ClassBreakdown struct {
	ClassName     string
	TeacherName   string
	TotalStudents int64
	Present       int64
	Absent        int64
	Excused       int64
	Late          int64
	NotMarked     int64
	Unexcused     []AbsentStudent
	ExcusedList   []AbsentStudent
	LateList      []AbsentStudent
}

type TodayOverview struct {
	Date              string
	TotalStudents     int64
	Present           int64
	Absent            int64
	Excused           int64
	Late              int64
	NotMarked         int64
	AttendancePercent float64
	UnexcusedList     []AbsentStudent
	ExcusedList       []AbsentStudent
	LateList          []AbsentStudent
	// Classes — sinflar kesimidagi taqsimot.
	Classes []ClassBreakdown
}

func todayDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func calcPercent(present, late, total int64) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(float64(present+late)/float64(total)*1000) / 10
}

func GetTodayOverview() TodayOverview {
	date := todayDate()
	overview := TodayOverview{Date: date.Format("02.01.2006")}

	database.DB.Model(&models.Student{}).Count(&overview.TotalStudents)

	type statusCount struct {
		Status string
		Count  int64
	}
	var counts []statusCount
	database.DB.Model(&models.Attendance{}).
		Select("status, count(*) as count").
		Where("date = ?", date).
		Group("status").
		Scan(&counts)

	for _, row := range counts {
		switch row.Status {
		case models.AttendancePresent:
			overview.Present = row.Count
		case models.AttendanceAbsent:
			overview.Absent = row.Count
		case models.AttendanceExcused:
			overview.Excused = row.Count
		case models.AttendanceLate:
			overview.Late = row.Count
		}
	}

	marked := overview.Present + overview.Absent + overview.Excused + overview.Late
	overview.NotMarked = overview.TotalStudents - marked
	if overview.NotMarked < 0 {
		overview.NotMarked = 0
	}
	overview.AttendancePercent = calcPercent(overview.Present, overview.Late, overview.TotalStudents)

	// Sinflar kesimi: har bir sinf uchun o'quvchilar soni.
	type classRow struct {
		ClassID     uint
		ClassName   string
		TeacherName string
		Total       int64
	}
	var classRows []classRow
	database.DB.Table("classes").
		Select("classes.id as class_id, COALESCE(class_names.name, '') as class_name, classes.teacher_full_name as teacher_name, COUNT(students.id) as total").
		Joins("LEFT JOIN class_names ON class_names.id = classes.class_name_id").
		Joins("LEFT JOIN students ON students.class_id = classes.id AND students.deleted_at IS NULL").
		Where("classes.deleted_at IS NULL").
		Group("classes.id, class_names.name, classes.teacher_full_name").
		Order("class_names.name ASC NULLS LAST").
		Scan(&classRows)

	breakdownIndex := make(map[uint]int, len(classRows))
	for _, row := range classRows {
		name := row.ClassName
		if name == "" {
			name = fmt.Sprintf("#%d", row.ClassID)
		}
		breakdownIndex[row.ClassID] = len(overview.Classes)
		overview.Classes = append(overview.Classes, ClassBreakdown{
			ClassName:     name,
			TeacherName:   row.TeacherName,
			TotalStudents: row.Total,
		})
	}

	var records []struct {
		ClassID   uint
		ClassName string
		FullName  string
		Status    string
		Reason    string
	}
	database.DB.Table("attendances").
		Select("attendances.class_id, COALESCE(class_names.name, '') as class_name, students.full_name, attendances.status, attendances.reason").
		Joins("JOIN students ON students.id = attendances.student_id").
		Joins("LEFT JOIN classes ON classes.id = attendances.class_id").
		Joins("LEFT JOIN class_names ON class_names.id = classes.class_name_id").
		Where("attendances.date = ? AND attendances.status IN ?", date, []string{
			models.AttendanceAbsent, models.AttendanceExcused, models.AttendanceLate,
		}).
		Order("class_names.name ASC NULLS LAST").
		Order("students.full_name asc").
		Scan(&records)

	for _, record := range records {
		entry := AbsentStudent{
			Name:      record.FullName,
			Reason:    record.Reason,
			ClassName: record.ClassName,
		}

		idx, hasClass := breakdownIndex[record.ClassID]

		switch record.Status {
		case models.AttendanceAbsent:
			overview.UnexcusedList = append(overview.UnexcusedList, entry)
			if hasClass {
				overview.Classes[idx].Absent++
				overview.Classes[idx].Unexcused = append(overview.Classes[idx].Unexcused, entry)
			}
		case models.AttendanceExcused:
			overview.ExcusedList = append(overview.ExcusedList, entry)
			if hasClass {
				overview.Classes[idx].Excused++
				overview.Classes[idx].ExcusedList = append(overview.Classes[idx].ExcusedList, entry)
			}
		case models.AttendanceLate:
			overview.LateList = append(overview.LateList, entry)
			if hasClass {
				overview.Classes[idx].Late++
				overview.Classes[idx].LateList = append(overview.Classes[idx].LateList, entry)
			}
		}
	}

	// Har bir sinf uchun "kelgan" va "belgilanmagan" sonini hisoblaymiz.
	type presentRow struct {
		ClassID uint
		Count   int64
	}
	var presentRows []presentRow
	database.DB.Table("attendances").
		Select("class_id, count(*) as count").
		Where("date = ? AND status = ?", date, models.AttendancePresent).
		Group("class_id").
		Scan(&presentRows)

	for _, row := range presentRows {
		if idx, ok := breakdownIndex[row.ClassID]; ok {
			overview.Classes[idx].Present = row.Count
		}
	}

	for i := range overview.Classes {
		cls := &overview.Classes[i]
		marked := cls.Present + cls.Absent + cls.Excused + cls.Late
		cls.NotMarked = cls.TotalStudents - marked
		if cls.NotMarked < 0 {
			cls.NotMarked = 0
		}
	}

	return overview
}

// FormatAdminOverviewMessage admin uchun bugungi davomat statistikasini
// SINFLAR KESIMIDA formatlaydi: avval umumiy yig'indi, keyin har bir sinf
// bo'yicha alohida blok (kelmaganlar va kech kelganlar ro'yxati bilan).
func FormatAdminOverviewMessage(stats TodayOverview) string {
	var b strings.Builder

	b.WriteString("📊 <b>Bugungi davomat statistikasi</b>\n")
	fmt.Fprintf(&b, "📅 Sana: <b>%s</b>\n\n", stats.Date)

	fmt.Fprintf(&b, "👥 Jami o'quvchilar: <b>%d</b>\n", stats.TotalStudents)
	fmt.Fprintf(&b, "✅ Kelgan: <b>%d</b>\n", stats.Present)
	fmt.Fprintf(&b, "🚫 Sababsiz: <b>%d</b>\n", stats.Absent)
	fmt.Fprintf(&b, "📝 Sababli: <b>%d</b>\n", stats.Excused)
	fmt.Fprintf(&b, "⏰ Kech kelgan: <b>%d</b>\n", stats.Late)
	fmt.Fprintf(&b, "❔ Belgilanmagan: <b>%d</b>\n", stats.NotMarked)
	fmt.Fprintf(&b, "📈 Davomat foizi: <b>%.1f%%</b>\n", stats.AttendancePercent)

	b.WriteString("\n———————————————\n")
	b.WriteString("🏫 <b>SINFLAR KESIMIDA</b>\n")

	if len(stats.Classes) == 0 {
		b.WriteString("\n — Sinflar topilmadi\n")
	}

	for _, cls := range stats.Classes {
		fmt.Fprintf(&b, "\n<b>%s</b>", cls.ClassName)
		if cls.TeacherName != "" && cls.TeacherName != "-" {
			fmt.Fprintf(&b, " <i>(%s)</i>", cls.TeacherName)
		}
		b.WriteString("\n")

		// Hali davomat topshirilmagan sinflarni alohida belgilaymiz.
		if cls.Present+cls.Absent+cls.Excused+cls.Late == 0 {
			fmt.Fprintf(&b, "   👥 %d ta o'quvchi — ⚠️ davomat topshirilmagan\n", cls.TotalStudents)
			continue
		}

		fmt.Fprintf(&b, "   👥 %d | ✅ %d | 🚫 %d | 📝 %d | ⏰ %d | ❔ %d\n",
			cls.TotalStudents, cls.Present, cls.Absent, cls.Excused, cls.Late, cls.NotMarked)

		if len(cls.Unexcused) > 0 {
			b.WriteString("   🚫 <b>Sababsiz:</b>\n")
			for _, student := range cls.Unexcused {
				fmt.Fprintf(&b, "      • %s\n", student.Name)
			}
		}

		if len(cls.ExcusedList) > 0 {
			b.WriteString("   📝 <b>Sababli:</b>\n")
			for _, student := range cls.ExcusedList {
				if student.Reason != "" {
					fmt.Fprintf(&b, "      • %s <i>(%s)</i>\n", student.Name, student.Reason)
				} else {
					fmt.Fprintf(&b, "      • %s\n", student.Name)
				}
			}
		}

		if len(cls.LateList) > 0 {
			b.WriteString("   ⏰ <b>Kech kelganlar:</b>\n")
			for _, student := range cls.LateList {
				fmt.Fprintf(&b, "      • %s\n", student.Name)
			}
		}
	}

	b.WriteString("\n———————————————\n")
	b.WriteString("🏫 Davomat topshirish uchun sinfni tanlang:")
	return b.String()
}

// telegramMessageLimit — Telegram bitta xabar uchun belgi cheklovi.
const telegramMessageLimit = 3800

// SplitMessage uzun hisobotni Telegram chekloviga mos bo'laklarga bo'ladi.
// Bo'lish qator chegarasida amalga oshiriladi, shunda HTML teglar buzilmaydi.
func SplitMessage(text string) []string {
	if len(text) <= telegramMessageLimit {
		return []string{text}
	}

	var parts []string
	var current strings.Builder

	for _, line := range strings.Split(text, "\n") {
		if current.Len()+len(line)+1 > telegramMessageLimit && current.Len() > 0 {
			parts = append(parts, current.String())
			current.Reset()
		}
		if current.Len() > 0 {
			current.WriteString("\n")
		}
		current.WriteString(line)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}
