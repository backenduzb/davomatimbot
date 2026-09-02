package models

import (
	"gorm.io/gorm"
)

// ClassName va Class modellari admin/ hamda gobot/ xizmatlaridagi modellarga
// mos holda saqlanadi (bir xil PostgreSQL jadvallar bilan ishlaydi).

type ClassName struct {
	gorm.Model
	ID   uint
	Name string `gorm:"size:64"`
}

type Class struct {
	gorm.Model
	ID              uint
	Updated         bool   `gorm:"default:false"`
	ClassNameID     uint
	ClassName       ClassName
	TeacherFullName string `gorm:"size:128"`
	TeacherTelegramId string `gorm:"type:varchar(128)"`
}
