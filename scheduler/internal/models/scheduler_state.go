package models

import "time"

// SchedulerState — scheduler xizmatining o'z yordamchi holati
// (masalan, oxirgi yuborilgan kunlik hisobot sanasi). Ushbu jadvalni faqat
// scheduler xizmati yuritadi; ilova (admin/gobot) jadvallariga tegmaydi.
type SchedulerState struct {
	ID        uint   `gorm:"primaryKey"`
	StateKey  string `gorm:"column:state_key;size:64;not null;uniqueIndex"`
	Value     string `gorm:"column:value;size:256;not null"`
	UpdatedAt time.Time
}
