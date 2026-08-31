package students

import (
	"bot/internal/database"
	"bot/internal/models"
	"bot/internal/repository/sessions"
)

func GetStudentsByClassID(classID uint) []models.Student {
	if classID == 0 {
		return nil
	}

	var students []models.Student
	err := database.DB.Model(&models.Student{}).
		Select("id", "full_name").
		Where("class_id = ?", classID).
		Order("full_name asc").
		Find(&students).Error
	if err != nil {
		return nil
	}
	return students
}

func GetAllStudents(telegramID uint) []models.Student {
	session := sessions.GetSession(telegramID)
	classID := session.ResolveClassID(telegramID)
	return GetStudentsByClassID(classID)
}

func GetStudentIDByName(name string, telegramID uint) uint {
	all := GetAllStudents(telegramID)
	for _, s := range all {
		if s.FullName == name {
			return s.ID
		}
	}
	return 0
}
