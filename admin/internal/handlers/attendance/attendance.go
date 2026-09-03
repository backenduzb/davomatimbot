package attendance

import (
	"admin/internal/database"
	"admin/internal/models"
	attendanceSchema "admin/internal/schemas/attendance"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}
	return time.Parse("2006-01-02", dateStr)
}

func GetToday(c *gin.Context) {
	db := database.DB

	dateStr := c.Query("date")
	classID := c.Query("class_id")

	date, err := parseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	studentQuery := db.Model(&models.Student{})
	if classID != "" {
		studentQuery = studentQuery.Where("class_id = ?", classID)
	}

	var students []models.Student
	if err := studentQuery.Order("full_name asc").Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	attendanceQuery := db.Where("date = ?", date)
	if classID != "" {
		attendanceQuery = attendanceQuery.Where("class_id = ?", classID)
	}

	var attendances []models.Attendance
	if err := attendanceQuery.Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	attendanceByStudent := make(map[uint]models.Attendance, len(attendances))
	for _, record := range attendances {
		attendanceByStudent[record.StudentID] = record
	}

	records := make([]attendanceSchema.Record, 0, len(students))
	for _, student := range students {
		status := models.AttendanceNotMarked
		reason := ""
		var recordID uint

		if record, ok := attendanceByStudent[student.ID]; ok {
			status = record.Status
			reason = record.Reason
			recordID = record.ID
		}

		records = append(records, attendanceSchema.Record{
			ID:          recordID,
			StudentID:   student.ID,
			StudentName: student.FullName,
			ClassID:     student.ClassID,
			Date:        date.Format("2006-01-02"),
			Status:      status,
			Reason:      reason,
		})
	}

	var classIDNum uint
	if classID != "" {
		parsed, err := strconv.ParseUint(classID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class_id"})
			return
		}
		classIDNum = uint(parsed)
	}

	c.JSON(http.StatusOK, attendanceSchema.TodayResponse{
		Date:    date.Format("2006-01-02"),
		ClassID: classIDNum,
		Records: records,
	})
}

func BatchUpsert(c *gin.Context) {
	db := database.DB

	var input attendanceSchema.BatchRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	validStatuses := map[string]bool{
		models.AttendancePresent: true,
		models.AttendanceAbsent:  true,
		models.AttendanceExcused: true,
		models.AttendanceLate:    true,
	}

	// Barcha o'quvchilarni va mavjud yozuvlarni BIR MARTA yuklaymiz —
	// ilgari har bir record uchun alohida SELECT bajarilardi (N+1).
	studentIDs := make([]uint, 0, len(input.Records))
	for _, record := range input.Records {
		if !validStatuses[record.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status: " + record.Status})
			return
		}
		studentIDs = append(studentIDs, record.StudentID)
	}

	if len(studentIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "attendance saved", "saved": 0})
		return
	}

	var students []models.Student
	if err := db.Where("id IN ?", studentIDs).Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	studentByID := make(map[uint]models.Student, len(students))
	for _, student := range students {
		studentByID[student.ID] = student
	}

	for _, record := range input.Records {
		student, ok := studentByID[record.StudentID]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
			return
		}
		if student.ClassID != input.ClassID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student does not belong to class"})
			return
		}
	}

	var existingRecords []models.Attendance
	if err := db.Where("date = ? AND student_id IN ?", date, studentIDs).Find(&existingRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	existingByStudent := make(map[uint]models.Attendance, len(existingRecords))
	for _, rec := range existingRecords {
		existingByStudent[rec.StudentID] = rec
	}

	toCreate := make([]models.Attendance, 0)
	toUpdate := make([]models.Attendance, 0)

	for _, record := range input.Records {
		if existing, ok := existingByStudent[record.StudentID]; ok {
			existing.Status = record.Status
			existing.Reason = record.Reason
			existing.ClassID = input.ClassID
			toUpdate = append(toUpdate, existing)
			continue
		}
		toCreate = append(toCreate, models.Attendance{
			StudentID: record.StudentID,
			ClassID:   input.ClassID,
			Date:      date,
			Status:    record.Status,
			Reason:    record.Reason,
		})
	}

	// Yozish tranzaksiyada va paketlar (batch) ko'rinishida bajariladi.
	if err := db.Transaction(func(tx *gorm.DB) error {
		if len(toCreate) > 0 {
			if err := tx.CreateInBatches(&toCreate, 200).Error; err != nil {
				return err
			}
		}
		for _, rec := range toUpdate {
			if err := tx.Model(&models.Attendance{}).
				Where("id = ?", rec.ID).
				Updates(map[string]any{
					"status":   rec.Status,
					"reason":   rec.Reason,
					"class_id": rec.ClassID,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance saved", "saved": len(input.Records)})
}
