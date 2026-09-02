package database

import (
	"log"
	"os"
	"time"

	"scheduler/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect admin/ va gobot/ xizmatlaridagi DB ulanish naqshini takrorlaydi.
func Connect(dsn string) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Warn,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("DB connections error", err)
	}

	// Ilova jadvalari allaqachon admin/gobot tomonidan yaratilgan bo'ladi;
	// AutoMigrate ularga o'zgarish kiritmaydi (faqat yetishmayotgan
	// scheduler_states jadvalini yaratadi).
	db.AutoMigrate(
		&models.ClassName{},
		&models.Class{},
		&models.Student{},
		&models.Attendance{},
		&models.SchedulerState{},
	)
	DB = db
}
