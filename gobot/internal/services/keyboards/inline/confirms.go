package inline

import (
	"fmt"

	"bot/internal/models"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func AbsentTypeKeyboard() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{Text: "✅ Sababsizlar yo'q", CallbackData: "no_unexcused"}},
			{{Text: "➕ Sababsizlarni qo'shish", CallbackData: "add_unexcused"}},
		},
	}
}

func AbsentConfirmKeyboard() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{Text: "➕ Yana qo'shish", CallbackData: "add_more_unexcused"}},
			{{Text: "✅ Sababsizlarni tugatish", CallbackData: "go_to_excused"}},
		},
	}
}

func ReasonConfirmKeyboard() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{Text: "➕ Sababli qo'shish", CallbackData: "add_more_excused"}},
			{{Text: "💾 Davomatni saqlash", CallbackData: "save_attendance"}},
		},
	}
}

func ClassChoiceKeyboard(classes []models.Class) gotgbot.InlineKeyboardMarkup {
	rows := make([][]gotgbot.InlineKeyboardButton, 0, len(classes))
	for _, class := range classes {
		label := class.TeacherFullName
		if class.ClassName.Name != "" {
			label = fmt.Sprintf("%s — %s", class.ClassName.Name, class.TeacherFullName)
		}
		rows = append(rows, []gotgbot.InlineKeyboardButton{
			{Text: label, CallbackData: fmt.Sprintf("class_%d", class.ID)},
		})
	}
	return gotgbot.InlineKeyboardMarkup{InlineKeyboard: rows}
}
