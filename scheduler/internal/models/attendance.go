package models

import (
	"time"

	"gorm.io/gorm"
)

// Davomat holatlari (admin/ va gobot/ xizmatlaridagi qadriyatlar bilan bir xil).
const (
	AttendancePresent   = "present"
	AttendanceAbsent    = "absent"
	AttendanceExcused   = "excused"
	AttendanceLate      = "late"
	AttendanceNotMarked = "not_marked"
)

// Attendance — har bir o'quvchining har bir sana uchun bitta yozuvi
// (uniqueIndex: idx_attendance_student_date). Ma'lumotlar sana asosida
// ajratiladi: yangi kunning davomati yangi yozuvlar orqali keladi, eski
// kunning yozuvlari o'chirilmasdan qoladi.
type Attendance struct {
	gorm.Model
	ID        uint
	StudentID uint `gorm:"not null;uniqueIndex:idx_attendance_student_date"`
	ClassID   uint `gorm:"not null"`
	Date      time.Time `gorm:"type:date;not null;uniqueIndex:idx_attendance_student_date"`
	Status    string    `gorm:"size:32;not null"`
	Reason    string    `gorm:"size:256"`
}
