package routes

import (
	"bot/internal/handlers"
	adminHandlers "bot/internal/handlers/admin"
	"bot/internal/handlers/users"
	"bot/internal/repository/states"
	"bot/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	tg "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func RegisterSimpleHandler(dp *ext.Dispatcher) {
	attendanceConversation := tg.NewConversation(
		[]ext.Handler{
			tg.NewCommand("start", handlers.HandleStart),
		},

		map[string][]ext.Handler{
			states.StateWaitingAdminClassChoice: {
				tg.NewCallback(func(cb *gotgbot.CallbackQuery) bool { return true }, adminHandlers.HandleAdminClassChoice),
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
