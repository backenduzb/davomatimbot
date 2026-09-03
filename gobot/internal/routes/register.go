package routes

import (
	"strings"

	"bot/internal/handlers"
	adminHandlers "bot/internal/handlers/admin"
	"bot/internal/handlers/users"
	"bot/internal/repository/states"
	"bot/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	tg "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// isPlainText buyruq bo'lmagan matnli xabarlarni filtrlaydi — shunda
// /cancel kabi buyruqlar fallback handlerga o'tadi.
func isPlainText(msg *gotgbot.Message) bool {
	return msg.Text != "" && !strings.HasPrefix(msg.Text, "/")
}

func RegisterSimpleHandler(dp *ext.Dispatcher) {
	attendanceConversation := tg.NewConversation(
		[]ext.Handler{
			tg.NewCommand("start", handlers.HandleStart),
		},

		map[string][]ext.Handler{
			// Admin sinfni reply klaviatura (matnli tugma) orqali tanlaydi.
			states.StateWaitingAdminClassChoice: {
				tg.NewMessage(isPlainText, adminHandlers.HandleAdminClassChoice),
			},
			states.StateWaitingAdminTeacherChoice: {
				tg.NewMessage(isPlainText, adminHandlers.HandleAdminTeacherChoice),
			},
			states.StateWaitingAbsentTypeChoice: {
				tg.NewCallback(func(cb *gotgbot.CallbackQuery) bool { return true }, users.HandleAbsentTypeChoice),
			},
			states.StateWaitingAbsentStudent: {
				tg.NewMessage(func(msg *gotgbot.Message) bool { return msg.Text != "" }, users.HandleAbsentStudentSelected),
			},
			states.StateWaitingAbsentConfirm: {
				tg.NewCallback(func(cb *gotgbot.CallbackQuery) bool { return true }, users.HandleAbsentConfirm),
			},
			states.StateWaitingReasonStudent: {
				tg.NewMessage(func(msg *gotgbot.Message) bool { return msg.Text != "" }, users.HandleReasonStudentSelected),
			},
			states.StateWaitingReasonInput: {
				tg.NewMessage(func(msg *gotgbot.Message) bool { return msg.Text != "" }, users.HandleReasonInput),
			},
			states.StateWaitingReasonConfirm: {
				tg.NewCallback(func(cb *gotgbot.CallbackQuery) bool { return true }, users.HandleReasonConfirm),
			},
			// Kech kelgan o'quvchilarni kiritish bosqichi.
			states.StateWaitingLateStudent: {
				tg.NewMessage(isPlainText, users.HandleLateStudentSelected),
			},
			states.StateWaitingLateConfirm: {
				tg.NewCallback(func(cb *gotgbot.CallbackQuery) bool { return true }, users.HandleLateConfirm),
			},
		},

		&tg.ConversationOpts{
			AllowReEntry: true,
			Fallbacks: []ext.Handler{
				tg.NewCommand("cancel", utils.HandleCancel),
			},
		},
	)

	dp.AddHandler(attendanceConversation)
}
