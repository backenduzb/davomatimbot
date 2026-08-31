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

	for _, record := range input.Records {
		if !validStatuses[record.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status: " + record.Status})
			return
		}

		var student models.Student
		if err := db.First(&student, record.StudentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
			return
		}
		if student.ClassID != input.ClassID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student does not belong to class"})
			return
		}

		var existing models.Attendance
		err := db.Where("student_id = ? AND date = ?", record.StudentID, date).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newRecord := models.Attendance{
				StudentID: record.StudentID,
				ClassID:   input.ClassID,
				Date:      date,
				Status:    record.Status,
				Reason:    record.Reason,
			}
			if err := db.Create(&newRecord).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			continue
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		existing.Status = record.Status
		existing.Reason = record.Reason
		existing.ClassID = input.ClassID
		if err := db.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance saved"})
}
