package attendance

import (
	"time"

	"scheduler/internal/database"
)

type ReportRow struct {
	ClassName string
	Student   string
	Date      time.Time
	Status    string
	Reason    string
}

func ListForDate(date time.Time) ([]ReportRow, error) {
	var rows []ReportRow

	err := database.DB.
		Table("attendances a").
		Select("cn.name AS class_name, s.full_name AS student, a.date AS date, a.status AS status, a.reason AS reason").
		Joins("JOIN students s ON s.id = a.student_id").
		Joins("JOIN classes c ON c.id = a.class_id").
		Joins("JOIN class_names cn ON cn.id = c.class_name_id").
		Where("a.date = ?", date.Format("2006-01-02")).
		Where("a.deleted_at IS NULL").
		Where("s.deleted_at IS NULL").
		Where("c.deleted_at IS NULL").
		Where("cn.deleted_at IS NULL").
		Order("cn.name ASC, s.full_name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
