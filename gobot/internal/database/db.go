package database

import (
	"log"
	"os"
	"time"

	"bot/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("DB connections error", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.ClassName{},
		&models.Student{},
		&models.Attendance{},
	)
	DB = db
}
