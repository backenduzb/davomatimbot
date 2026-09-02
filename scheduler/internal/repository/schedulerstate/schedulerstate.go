package schedulerstate

import (
	"time"

	"scheduler/internal/database"
	"scheduler/internal/models"

	"gorm.io/gorm/clause"
)

const LastReportDateKey = "last_report_date"

type Store struct{}

var Default = Store{}

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
