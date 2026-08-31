package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	ID uint
	FullName string `gorm:"size:128"`
	ClassID uint
	Class Class `gorm:"consentraint:OnUpdate:CASCADE;OnDelete:RESTRICT;"`
}