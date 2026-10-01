package settings

import (
	"bot/internal/database"
	"bot/internal/models"

	"gorm.io/gorm/clause"
)

// KeyGuestCongrats — ustoz bo'lmagan foydalanuvchilar uchun tabrik
// funksiyasini yoqish/o'chirish kaliti ("1" = yoqilgan, "0" = o'chirilgan).
const KeyGuestCongrats = "guest_congrats_enabled"

// IsGuestCongratsEnabled tabrik funksiyasi yoqilganligini qaytaradi.
// Sozlama hali saqlanmagan bo'lsa default HOLAT — yoqilgan.
func IsGuestCongratsEnabled() bool {
	var s models.BotSetting
	if err := database.DB.Where("key = ?", KeyGuestCongrats).First(&s).Error; err != nil {
		return true
	}
	return s.Value == "1"
}

// SetGuestCongratsEnabled tabrik funksiyasini yoqadi/o'chiradi.
func SetGuestCongratsEnabled(enabled bool) error {
	value := "0"
	if enabled {
		value = "1"
	}
	s := models.BotSetting{Key: KeyGuestCongrats, Value: value}
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&s).Error
}
