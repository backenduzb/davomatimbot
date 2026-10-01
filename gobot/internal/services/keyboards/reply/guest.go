package replyKeyboards

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// TeachersKeyboard mehmon tabrik yo'llashi uchun ustoz tanlash klaviaturasi.
// Ismlar uzun bo'lishi mumkinligi uchun har bir tugma alohida qatorda.
func TeachersKeyboard(names []string) gotgbot.ReplyKeyboardMarkup {
	rows := make([][]gotgbot.KeyboardButton, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		rows = append(rows, []gotgbot.KeyboardButton{{Text: name}})
	}
	return gotgbot.ReplyKeyboardMarkup{
		Keyboard:              rows,
		ResizeKeyboard:        true,
		OneTimeKeyboard:       true,
		InputFieldPlaceholder: "Ustozni tanlang",
	}
}
