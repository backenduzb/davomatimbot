package replyKeyboards

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"bot/internal/repository/students"
	"bot/internal/models"
	"bot/internal/repository/sessions"
)

func StudentsKeyboard(telegramID uint) gotgbot.ReplyKeyboardMarkup {
	var rows [][]gotgbot.KeyboardButton
	var currentRow []gotgbot.KeyboardButton
	allStudents := students.GetAllStudents(telegramID)
	for _, student := range allStudents {
		button := gotgbot.KeyboardButton{
            Text: student.FullName,
        }
        currentRow = append(currentRow, button)
        if len(currentRow) == 2 {
        	rows = append(rows, currentRow)
        	currentRow = []gotgbot.KeyboardButton{}
        }
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}
	return gotgbot.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
}

func FilteredStudentsKeyboard(allStudents []models.Student, session *sessions.AttendanceSession) gotgbot.ReplyKeyboardMarkup {
	var rows [][]gotgbot.KeyboardButton
	var currentRow []gotgbot.KeyboardButton

	for _, student := range allStudents {
		if session.IsAlreadyMarked(student.ID) {
			continue
		}

		button := gotgbot.KeyboardButton{
			Text: student.FullName,
		}

		currentRow = append(currentRow, button)
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
		OneTimeKeyboard:       false,
		InputFieldPlaceholder: "O'quvchini tanlang",
	}
}

// HasButtons klaviaturada kamida bitta tugma bor-yo'qligini tekshiradi.
// Bo'sh klaviatura yuborilsa Telegram uni e'tiborsiz qoldiradi va
// foydalanuvchi "tugmalar kelmadi" holatiga tushib qoladi.
func HasButtons(markup gotgbot.ReplyKeyboardMarkup) bool {
	for _, row := range markup.Keyboard {
		if len(row) > 0 {
			return true
		}
	}
	return false
}