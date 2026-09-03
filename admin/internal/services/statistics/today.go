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

	// Sinflar sinf nomi bo'yicha tartiblanadi — dashboard jadvali va
	// diagrammalarda ular doim bir xil, tushunarli tartibda chiqadi.
	var classes []models.Class
	if err := db.Model(&models.Class{}).
		Preload("ClassName").
		Select("classes.*").
		Joins("LEFT JOIN class_names ON class_names.id = classes.class_name_id").
		Order("class_names.name ASC NULLS LAST").
		Order("classes.id ASC").
		Find(&classes).Error; err != nil {
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

	// O'quvchilar sonini har bir sinf uchun alohida COUNT qilish o'rniga
	// (N+1 muammosi — 100 ta sinf = 100 ta so'rov) bitta GROUP BY so'rovi
	// bilan olamiz. Aynan shu joy dashboard statistikasini sekinlashtirardi.
	type studentCountRow struct {
		ClassID uint
		Count   int64
	}
	var studentCountRows []studentCountRow
	if err := db.Model(&models.Student{}).
		Select("class_id, count(*) as count").
		Where("deleted_at IS NULL").
		Group("class_id").
		Scan(&studentCountRows).Error; err != nil {
		return nil, err
	}
	studentCountByClass := make(map[uint]int64, len(studentCountRows))
	for _, row := range studentCountRows {
		studentCountByClass[row.ClassID] = row.Count
	}

	result := &attendanceSchema.TodayStatistics{
		Date:         date.Format("2006-01-02"),
		TotalClasses: int64(len(classes)),
		Classes:      make([]attendanceSchema.ClassStats, 0, len(classes)),
	}

	for _, class := range classes {
		studentCount := studentCountByClass[class.ID]

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
