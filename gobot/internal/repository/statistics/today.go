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
	Name   string
	Reason string
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

	var records []struct {
		FullName string
		Status   string
		Reason   string
	}
	database.DB.Table("attendances").
		Select("students.full_name, attendances.status, attendances.reason").
		Joins("JOIN students ON students.id = attendances.student_id").
		Where("attendances.date = ? AND attendances.status IN ?", date, []string{models.AttendanceAbsent, models.AttendanceExcused}).
		Order("students.full_name asc").
		Scan(&records)

	for _, record := range records {
		switch record.Status {
		case models.AttendanceAbsent:
			overview.UnexcusedList = append(overview.UnexcusedList, AbsentStudent{Name: record.FullName})
		case models.AttendanceExcused:
			overview.ExcusedList = append(overview.ExcusedList, AbsentStudent{Name: record.FullName, Reason: record.Reason})
		}
	}

	return overview
}

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
	fmt.Fprintf(&b, "📈 Davomat foizi: <b>%.1f%%</b>\n\n", stats.AttendancePercent)

	b.WriteString("🚫 <b>Sababsiz kelmaganlar:</b>\n")
	if len(stats.UnexcusedList) == 0 {
		b.WriteString(" — Yo'q\n")
	} else {
		for _, student := range stats.UnexcusedList {
			fmt.Fprintf(&b, " • %s\n", student.Name)
		}
	}

	b.WriteString("\n📝 <b>Sababli kelmaganlar:</b>\n")
	if len(stats.ExcusedList) == 0 {
		b.WriteString(" — Yo'q\n")
	} else {
		for _, student := range stats.ExcusedList {
			if student.Reason != "" {
				fmt.Fprintf(&b, " • %s <i>(%s)</i>\n", student.Name, student.Reason)
			} else {
				fmt.Fprintf(&b, " • %s\n", student.Name)
			}
		}
	}

	b.WriteString("\n🏫 Davomat topshirish uchun sinfni tanlang:")
	return b.String()
}
