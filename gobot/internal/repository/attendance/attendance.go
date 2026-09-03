package attendance

import (
	"bot/internal/database"
	"bot/internal/models"
	"bot/internal/repository/classes"
	"bot/internal/repository/sessions"
	"time"
)

func SaveClassAttendance(classID uint, session *sessions.AttendanceSession) error {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var students []models.Student
	if err := database.DB.Where("class_id = ?", classID).Find(&students).Error; err != nil {
		return err
	}

	for _, student := range students {
		status := models.AttendancePresent
		reason := ""

		if _, ok := session.UnexcusedStudents[student.ID]; ok {
			status = models.AttendanceAbsent
		} else if detail, ok := session.ExcusedStudents[student.ID]; ok {
			status = models.AttendanceExcused
			reason = detail.Reason
		} else if _, ok := session.LateStudents[student.ID]; ok {
			// Kech kelgan o'quvchilar alohida "late" statusi bilan saqlanadi.
			status = models.AttendanceLate
		}

		var existing models.Attendance
		err := database.DB.Where("student_id = ? AND date = ?", student.ID, date).First(&existing).Error
		if err != nil {
			record := models.Attendance{
				StudentID: student.ID,
				ClassID:   classID,
				Date:      date,
				Status:    status,
				Reason:    reason,
			}
			if err := database.DB.Create(&record).Error; err != nil {
				return err
			}
			continue
		}

		existing.Status = status
		existing.Reason = reason
		existing.ClassID = classID
		if err := database.DB.Save(&existing).Error; err != nil {
			return err
		}
	}

	if err := classes.MarkClassUpdated(classID); err != nil {
		return err
	}

	return nil
}
