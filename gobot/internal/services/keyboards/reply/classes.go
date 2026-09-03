package replyKeyboards

import (
	"sort"

	"bot/internal/models"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// ClassNamesKeyboard admin uchun sinf tanlash reply klaviaturasini quradi.
// Tugmalarda FAQAT sinf nomi ko'rsatiladi (o'qituvchi F.I.Sh emas).
// Bir xil nomli sinflar bitta tugmaga birlashtiriladi — ular orasidan
// tanlash keyingi qadamda (o'qituvchi bo'yicha) amalga oshiriladi.
func ClassNamesKeyboard(classList []models.Class) gotgbot.ReplyKeyboardMarkup {
	seen := make(map[string]struct{}, len(classList))
	names := make([]string, 0, len(classList))

	for _, class := range classList {
		name := class.ClassName.Name
		if name == "" || name == "-" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}

	sort.Strings(names)

	return buildNameKeyboard(names)
}

// TeacherChoiceKeyboard bir xil nomli sinflar uchun o'qituvchini tanlash
// klaviaturasi (faqat bir nechta mos sinf topilganda ishlatiladi).
func TeacherChoiceKeyboard(classList []models.Class) gotgbot.ReplyKeyboardMarkup {
	names := make([]string, 0, len(classList))
	for _, class := range classList {
		if class.TeacherFullName == "" {
			continue
		}
		names = append(names, class.TeacherFullName)
	}
	sort.Strings(names)
	return buildNameKeyboard(names)
}

// buildNameKeyboard matnlar ro'yxatidan 2 ustunli reply klaviatura yasaydi.
func buildNameKeyboard(names []string) gotgbot.ReplyKeyboardMarkup {
	var rows [][]gotgbot.KeyboardButton
	var currentRow []gotgbot.KeyboardButton

	for _, name := range names {
		currentRow = append(currentRow, gotgbot.KeyboardButton{Text: name})
		if len(currentRow) == 2 {
			rows = append(rows, currentRow)
			currentRow = []gotgbot.KeyboardButton{}
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	return gotgbot.ReplyKeyboardMarkup{
		Keyboard:              rows,
		ResizeKeyboard:        true,
		OneTimeKeyboard:       true,
		InputFieldPlaceholder: "Sinfni tanlang",
	}
}
