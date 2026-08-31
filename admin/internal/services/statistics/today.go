package statistics

import (
	"admin/internal/models"
	attendanceSchema "admin/internal/schemas/attendance"
	"math"
	"time"

	"gorm.io/gorm"
)

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}
	return time.Parse("2006-01-02", dateStr)
}

func calcPercent(present, late, total int64) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(float64(present+late)/float64(total)*1000) / 10
}

func GetTodayStatistics(db *gorm.DB, dateStr string) (*attendanceSchema.TodayStatistics, error) {
	date, err := parseDate(dateStr)
	if err != nil {
		return nil, err
	}

	var classes []models.Class
	if err := db.Preload("ClassName").Find(&classes).Error; err != nil {
		return nil, err
	}

	type attendanceRow struct {
		ClassID uint
		Status  string
		Count   int64
	}

	var attendanceCounts []attendanceRow
	if err := db.Model(&models.Attendance{}).
		Select("class_id, status, count(*) as count").
		Where("date = ?", date).
		Group("class_id, status").
		Scan(&attendanceCounts).Error; err != nil {
		return nil, err
	}

	countsByClass := make(map[uint]map[string]int64)
	for _, row := range attendanceCounts {
		if countsByClass[row.ClassID] == nil {
			countsByClass[row.ClassID] = make(map[string]int64)
		}
		countsByClass[row.ClassID][row.Status] = row.Count
	}

	result := &attendanceSchema.TodayStatistics{
		Date:         date.Format("2006-01-02"),
		TotalClasses: int64(len(classes)),
		Classes:      make([]attendanceSchema.ClassStats, 0, len(classes)),
	}

	for _, class := range classes {
		var studentCount int64
		if err := db.Model(&models.Student{}).Where("class_id = ?", class.ID).Count(&studentCount).Error; err != nil {
			return nil, err
		}

		statusCounts := countsByClass[class.ID]
		present := statusCounts[models.AttendancePresent]
		absent := statusCounts[models.AttendanceAbsent]
		excused := statusCounts[models.AttendanceExcused]
		late := statusCounts[models.AttendanceLate]
		marked := present + absent + excused + late
		notMarked := studentCount - marked
		if notMarked < 0 {
			notMarked = 0
		}

		className := ""
		if class.ClassName.ID != 0 {
			className = class.ClassName.Name
		}

		classStats := attendanceSchema.ClassStats{
			ClassID:           class.ID,
			ClassName:         className,
			TotalStudents:     studentCount,
			Present:           present,
			Absent:            absent,
			Excused:           excused,
			Late:              late,
			NotMarked:         notMarked,
			AttendancePercent: calcPercent(present, late, studentCount),
		}

		result.TotalStudents += studentCount
		result.Present += present
		result.Absent += absent
		result.Excused += excused
		result.Late += late
		result.NotMarked += notMarked
		result.Classes = append(result.Classes, classStats)
	}

	result.AttendancePercent = calcPercent(result.Present, result.Late, result.TotalStudents)

	return result, nil
}
