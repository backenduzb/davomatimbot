package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	AttendancePresent   = "present"
	AttendanceAbsent    = "absent"
	AttendanceExcused   = "excused"
	AttendanceLate      = "late"
	AttendanceNotMarked = "not_marked"
)

type Attendance struct {
	gorm.Model
	ID        uint
	StudentID uint      `gorm:"not null;uniqueIndex:idx_attendance_student_date"`
	Student   Student   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ClassID   uint      `gorm:"not null"`
	Class     Class     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Date      time.Time `gorm:"type:date;not null;uniqueIndex:idx_attendance_student_date"`
	Status    string    `gorm:"size:32;not null"`
	Reason    string    `gorm:"size:256"`
}
