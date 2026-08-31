package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	IsOnline bool `gorm:"default:false"`
	IsBanned bool `gorm:"default:false"`
	IsAdmin bool `gorm:"default:false"`
	TelegramId string `gorm:"column:telegram_id;type:varchar(128)"`
	
	LastSeen time.Time
}