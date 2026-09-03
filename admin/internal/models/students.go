package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	ID uint
	FullName string `gorm:"size:128"`
	ClassID uint `gorm:"index"`
	Class Class `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}