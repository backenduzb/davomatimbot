package classes

import (
	"bot/internal/database"
	"bot/internal/models"
	"strconv"
)

func GetClassInfo(telegramID uint) (TeacherFullName string, ClassName string) {
	var class models.Class

	err := database.DB.Preload("ClassName").Where("teacher_telegram_id = ?", formatTelegramID(telegramID)).First(&class).Error

	if err != nil {
		return "-", "-"
	}
	return class.TeacherFullName, class.ClassName.Name
}

func GetClassInfoByID(classID uint) (TeacherFullName string, ClassName string) {
	var class models.Class
	err := database.DB.Preload("ClassName").First(&class, classID).Error
	if err != nil {
		return "-", "-"
	}
	return class.TeacherFullName, class.ClassName.Name
}

func GetClassID(telegramID uint) uint {
	var class models.Class
	err := database.DB.Where("teacher_telegram_id = ?", formatTelegramID(telegramID)).First(&class).Error
	if err != nil {
		return 0
	}
	return class.ID
}

func GetAllClasses() []models.Class {
	var list []models.Class
	database.DB.Preload("ClassName").Order("id asc").Find(&list)
	return list
}

func MarkClassUpdated(classID uint) error {
	return database.DB.Model(&models.Class{}).Where("id = ?", classID).Update("updated", true).Error
}

func formatTelegramID(telegramID uint) string {
	return strconv.FormatUint(uint64(telegramID), 10)
}
