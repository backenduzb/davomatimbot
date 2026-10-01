package classes

import (
	"strconv"
	"strings"
	"time"

	"scheduler/internal/database"
)

// UnsubmittedClass bugungi davomat hali topshirilmagan sinf.
type UnsubmittedClass struct {
	ClassID           uint
	ClassName         string
	TeacherFullName   string
	TeacherTelegramID int64
}

// ListUnsubmittedForDate berilgan sanada davomat yozuvi yo'q sinflarni
// qaytaradi. Faqat o'quvchisi bor va ustoz telegram ID si ko'rsatilgan
// sinflar olinadi — eslatmani yuborib bo'lmaydiganlar chiqarib tashlanadi.
func ListUnsubmittedForDate(date time.Time) ([]UnsubmittedClass, error) {
	dateStr := date.Format("2006-01-02")

	type row struct {
		ClassID           uint
		ClassName         string
		TeacherFullName   string
		TeacherTelegramID string
		StudentsCount     int64
		AttendanceCount   int64
	}

	var rows []row
	err := database.DB.
		Table("classes c").
		Select(`c.id AS class_id,
			COALESCE(cn.name, '') AS class_name,
			COALESCE(c.teacher_full_name, '') AS teacher_full_name,
			COALESCE(c.teacher_telegram_id, '') AS teacher_telegram_id,
			COUNT(DISTINCT s.id) AS students_count,
			COUNT(DISTINCT a.id) AS attendance_count`).
		Joins("LEFT JOIN class_names cn ON cn.id = c.class_name_id AND cn.deleted_at IS NULL").
		Joins("LEFT JOIN students s ON s.class_id = c.id AND s.deleted_at IS NULL").
		Joins("LEFT JOIN attendances a ON a.class_id = c.id AND a.date = ? AND a.deleted_at IS NULL", dateStr).
		Where("c.deleted_at IS NULL").
		Group("c.id, cn.name, c.teacher_full_name, c.teacher_telegram_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]UnsubmittedClass, 0, len(rows))
	for _, r := range rows {
		if r.StudentsCount == 0 {
			continue
		}
		if r.AttendanceCount > 0 {
			continue
		}
		tidStr := strings.TrimSpace(r.TeacherTelegramID)
		if tidStr == "" || tidStr == "-" || tidStr == "0" {
			continue
		}
		tid, err := strconv.ParseInt(tidStr, 10, 64)
		if err != nil || tid == 0 {
			continue
		}
		name := r.ClassName
		if name == "" {
			name = "Noma'lum sinf"
		}
		out = append(out, UnsubmittedClass{
			ClassID:           r.ClassID,
			ClassName:         name,
			TeacherFullName:   r.TeacherFullName,
			TeacherTelegramID: tid,
		})
	}
	return out, nil
}
