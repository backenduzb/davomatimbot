package models

import (
	"gorm.io/gorm"
)

type ClassName struct {
	gorm.Model
	ID uint
	Name string `gorm:"size:64"`
}

type Class struct {
	gorm.Model
	ID uint
	Updated bool `gorm:"default:false"`
	ClassNameID uint
	ClassName ClassName `gorm:"foreignKey:ClassNameID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	TeacherFullName string `gorm:"size:128"`
	TeacherTelegramId string `gorm:"size:128"`
}