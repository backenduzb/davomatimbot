package classes

import (
	"bot/internal/database"
	"bot/internal/models"
	"sort"
	"strconv"
	"strings"
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
	database.DB.Model(&models.Class{}).
		Preload("ClassName").
		Select("classes.*").
		Joins("LEFT JOIN class_names ON class_names.id = classes.class_name_id").
		Order("class_names.name ASC NULLS LAST").
		Order("classes.id ASC").
		Find(&list)
	return list
}

func MarkClassUpdated(classID uint) error {
	return database.DB.Model(&models.Class{}).Where("id = ?", classID).Update("updated", true).Error
}

func GetAllTeacherChatIDs() []int64 {
	var rawIDs []string
	database.DB.Model(&models.Class{}).
		Where("teacher_telegram_id <> '' AND teacher_telegram_id IS NOT NULL").
		Distinct().
		Pluck("teacher_telegram_id", &rawIDs)

	seen := make(map[int64]struct{}, len(rawIDs))
	out := make([]int64, 0, len(rawIDs))
	for _, s := range rawIDs {
		s = strings.TrimSpace(s)
		if s == "" || s == "-" || s == "0" {
			continue
		}
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil || id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// TeacherInfo — tabrik/menyu uchun bitta ustoz (telegram ID bo'yicha yagona).
type TeacherInfo struct {
	TelegramID int64
	FullName   string
}

// GetAllTeachers barcha ustozlarni ismi bo'yicha alifbo tartibida qaytaradi
// (bir ustoz bitta marta — telegram ID bo'yicha). Ismi yoki telegram ID si
// yaroqsiz yozuvlar tashlab yuboriladi.
func GetAllTeachers() []TeacherInfo {
	classList := GetAllClasses()
	seen := make(map[int64]struct{}, len(classList))
	out := make([]TeacherInfo, 0, len(classList))
	for _, c := range classList {
		name := strings.TrimSpace(c.TeacherFullName)
		tidStr := strings.TrimSpace(c.TeacherTelegramId)
		if name == "" || name == "-" || tidStr == "" || tidStr == "-" || tidStr == "0" {
			continue
		}
		tid, err := strconv.ParseInt(tidStr, 10, 64)
		if err != nil || tid == 0 {
			continue
		}
		if _, ok := seen[tid]; ok {
			continue
		}
		seen[tid] = struct{}{}
		out = append(out, TeacherInfo{TelegramID: tid, FullName: name})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].FullName < out[j].FullName })
	return out
}

// FindTeacherByName ustozni F.I.Sh bo'yicha qidiradi (registrga sezgir emas).
func FindTeacherByName(teachers []TeacherInfo, name string) (TeacherInfo, bool) {
	target := strings.TrimSpace(name)
	for _, t := range teachers {
		if strings.EqualFold(t.FullName, target) {
			return t, true
		}
	}
	return TeacherInfo{}, false
}

func formatTelegramID(telegramID uint) string {
	return strconv.FormatUint(uint64(telegramID), 10)
}
