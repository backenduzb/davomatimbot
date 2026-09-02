package attendance

import (
	"os"
	"testing"
	"time"

	"scheduler/internal/database"
	"scheduler/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestListForDateIntegration(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL o'rnatilmagan — integratsiya sinovi o'tkazib yuborildi")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("DB ulanish: %v", err)
	}

	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	if err := db.AutoMigrate(
		&models.ClassName{}, &models.Class{}, &models.Student{}, &models.Attendance{},
	); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	reportDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	suffix := time.Now().Format("20060102150405")
	className := "CI-TEST-" + suffix

	var (
		cn    models.ClassName
		cls   models.Class
		stA   models.Student
		stB   models.Student
	)
	cn = models.ClassName{Name: className}
	if err := db.Create(&cn).Error; err != nil {
		t.Fatalf("ClassName yaratish: %v", err)
	}
	cls = models.Class{
		ClassNameID:     cn.ID,
		TeacherFullName: "CI Teacher",
		TeacherTelegramId: "111",
	}
	if err := db.Create(&cls).Error; err != nil {
		t.Fatalf("Class yaratish: %v", err)
	}
	stA = models.Student{FullName: "Test B Oquvchi", ClassID: cls.ID}
	stB = models.Student{FullName: "Test A Oquvchi", ClassID: cls.ID}
	if err := db.Create(&stA).Error; err != nil {
		t.Fatalf("Student yaratish: %v", err)
	}
	if err := db.Create(&stB).Error; err != nil {
		t.Fatalf("Student yaratish: %v", err)
	}

	records := []models.Attendance{
		{StudentID: stA.ID, ClassID: cls.ID, Date: reportDate, Status: models.AttendanceAbsent, Reason: "Illness"},
		{StudentID: stB.ID, ClassID: cls.ID, Date: reportDate, Status: models.AttendancePresent},
	}
	if err := db.Create(&records).Error; err != nil {
		t.Fatalf("Attendance yaratish: %v", err)
	}
	defer func() {
		db.Unscoped().Delete(&records, "student_id IN ?", []uint{stA.ID, stB.ID})
		db.Unscoped().Delete(&[]models.Student{stA, stB})
		db.Unscoped().Delete(&cls)
		db.Unscoped().Delete(&cn)
	}()

	rows, err := ListForDate(reportDate)
	if err != nil {
		t.Fatalf("ListForDate: %v", err)
	}

	var mine []ReportRow
	for _, r := range rows {
		if r.ClassName == className {
			mine = append(mine, r)
		}
	}
	if len(mine) != 2 {
		t.Fatalf("qatorlar soni = %d, want 2: %+v", len(mine), mine)
	}
	// Ism bo'yicha tartib: "Test A Oquvchi" oldinda.
	if mine[0].Student != "Test A Oquvchi" || mine[0].Status != "present" {
		t.Errorf("1-qator = %+v", mine[0])
	}
	if mine[1].Student != "Test B Oquvchi" || mine[1].Status != "absent" || mine[1].Reason != "Illness" {
		t.Errorf("2-qator = %+v", mine[1])
	}
	if mine[0].ClassName != className {
		t.Errorf("sinf nomi = %q, want %q", mine[0].ClassName, className)
	}

	otherRows, err := ListForDate(time.Date(2020, 1, 16, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("ListForDate (boshqa sana): %v", err)
	}
	for _, r := range otherRows {
		if r.ClassName == className {
			t.Errorf("boshqa sanada sinf yozuvi topildi: %+v", r)
		}
	}
}

func TestListForDateSoftDeleted(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL o'rnatilmagan — integratsiya sinovi o'tkazib yuborildi")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("DB ulanish: %v", err)
	}

	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	if err := db.AutoMigrate(
		&models.ClassName{}, &models.Class{}, &models.Student{}, &models.Attendance{},
	); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	reportDate := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	suffix := time.Now().Format("20060102150405") + "s"
	className := "CI-TEST-" + suffix

	var (
		cn  models.ClassName
		cls models.Class
		st  models.Student
	)
	cn = models.ClassName{Name: className}
	if err := db.Create(&cn).Error; err != nil {
		t.Fatal(err)
	}
	cls = models.Class{ClassNameID: cn.ID, TeacherFullName: "CI Teacher", TeacherTelegramId: "222"}
	if err := db.Create(&cls).Error; err != nil {
		t.Fatal(err)
	}
	st = models.Student{FullName: "Deleted Oquvchi", ClassID: cls.ID}
	if err := db.Create(&st).Error; err != nil {
		t.Fatal(err)
	}
	rec := models.Attendance{StudentID: st.ID, ClassID: cls.ID, Date: reportDate, Status: models.AttendancePresent}
	if err := db.Create(&rec).Error; err != nil {
		t.Fatal(err)
	}

	defer func() {
		db.Unscoped().Delete(&rec)
		db.Unscoped().Delete(&st)
		db.Unscoped().Delete(&cls)
		db.Unscoped().Delete(&cn)
	}()

	if err := db.Delete(&st).Error; err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	rows, err := ListForDate(reportDate)
	if err != nil {
		t.Fatalf("ListForDate: %v", err)
	}
	for _, r := range rows {
		if r.Student == "Deleted Oquvchi" {
			t.Error("soft-delete qilingan o'quvchi hisobotda topildi")
		}
	}
}
