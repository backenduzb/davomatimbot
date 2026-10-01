package routes

import (
	"strings"

	"bot/internal/handlers"
	adminHandlers "bot/internal/handlers/admin"
	"bot/internal/handlers/guest"
	"bot/internal/handlers/users"
	"bot/internal/repository/states"
	replyKeyboards "bot/internal/services/keyboards/reply"
	"bot/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	tg "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func isPlainText(msg *gotgbot.Message) bool {
	return msg.Text != "" && !strings.HasPrefix(msg.Text, "/")
}

func isDavomatButton(msg *gotgbot.Message) bool {
	return msg.Text != "" && replyKeyboards.IsDavomatTopshirish(msg.Text)
}

func isBroadcastInput(msg *gotgbot.Message) bool {
	if msg == nil {
		return false
	}
	if msg.Text != "" {
		return !strings.HasPrefix(msg.Text, "/")
	}
	if msg.Caption != "" {
		return true
	}
	return msg.Sticker != nil ||
		msg.Animation != nil ||
		len(msg.Photo) > 0 ||
		msg.Video != nil ||
		msg.Document != nil ||
		msg.Voice != nil ||
		msg.Audio != nil
}

func RegisterSimpleHandler(dp *ext.Dispatcher) {
	attendanceConversation := tg.NewConversation(
		[]ext.Handler{
			tg.NewCommand("start", handlers.HandleStart),
			tg.NewMessage(isDavomatButton, handlers.HandleStart),
		},

		map[string][]ext.Handler{
			states.StateWaitingAdminMenu: {
				tg.NewMessage(isPlainText, adminHandlers.HandleAdminMenuChoice),
			},
			states.StateWaitingAdminBroadcast: {
				tg.NewMessage(isBroadcastInput, adminHandlers.HandleAdminBroadcastInput),
			},
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
			states.StateWaitingLateStudent: {
				tg.NewMessage(isPlainText, users.HandleLateStudentSelected),
			},
			states.StateWaitingLateConfirm: {
				tg.NewCallback(func(cb *gotgbot.CallbackQuery) bool { return true }, users.HandleLateConfirm),
			},
			states.StateWaitingGuestTeacherChoice: {
				tg.NewMessage(isPlainText, guest.HandleGuestTeacherChoice),
			},
			states.StateWaitingGuestCongratsInput: {
				tg.NewMessage(isBroadcastInput, guest.HandleGuestCongratsInput),
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
