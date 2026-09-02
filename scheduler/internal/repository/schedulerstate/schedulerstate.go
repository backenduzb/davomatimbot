package schedulerstate

import (
	"time"

	"scheduler/internal/database"
	"scheduler/internal/models"

	"gorm.io/gorm/clause"
)

// LastReportDateKey — oxirgi muvaffaqiyatli yuborilgan hisobot sanasi
// uchun holat kaliti.
const LastReportDateKey = "last_report_date"

// Store scheduler holatini PostgreSQL'da saqlaydi. Bu kichik jadval
// scheduler restart/redeploy holatida takroriy hisobot yuborilishining
// oldini oladi (konteyner ichidagi fayl yo'qolib ketishi mumkin,
// baza esa barqaror).
type Store struct{}

// Default — database.DB global ulanishidan foydalanadigan instance.
var Default = Store{}

// GetLastReportDate oxirgi yuborilgan hisobot sanasini qaytaradi.
// Yozum mavjud bo'lmasa zero qaytaradi.
func (Store) GetLastReportDate() (date time.Time) {
	var state models.SchedulerState
	if err := database.DB.Where("state_key = ?", LastReportDateKey).First(&state).Error; err != nil {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", state.Value)
	if err != nil {
		return time.Time{}
	}
	return t
}

// SetLastReportDate oxirgi yuborilgan hisobot sanasini saqlaydi (upsert).
func (Store) SetLastReportDate(date time.Time) error {
	state := models.SchedulerState{
		StateKey: LastReportDateKey,
		Value:    date.Format("2006-01-02"),
	}
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "state_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&state).Error
}
