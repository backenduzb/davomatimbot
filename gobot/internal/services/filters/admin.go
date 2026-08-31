package filters

import (
	"errors"
	"strconv"

	"bot/internal/database"
	"bot/internal/models"

	"gorm.io/gorm"
)

func formatTelegramID(telegramID uint) string {
	return strconv.FormatUint(uint64(telegramID), 10)
}

func CheckIsTeacher(telegramID uint) bool {
	var class models.Class
	err := database.DB.Where("teacher_telegram_id = ?", formatTelegramID(telegramID)).First(&class).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}
	return true
}

func CheckIsAdmin(telegramID uint) bool {
	var user models.User
	err := database.DB.Where("telegram_id = ? AND is_admin = ?", formatTelegramID(telegramID), true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}
	return true
}
