package utils

import (
	"bot/internal/repository/sessions"
	"bot/internal/services/filters"
	replyKeyboards "bot/internal/services/keyboards/reply"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	sessions.DeleteSession(userID)

	if filters.CheckIsAdmin(userID) {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Jarayon bekor qilindi.\nMenyudan birini tanlang:", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
		})
		return handlers.EndConversation()
	}

	if filters.CheckIsTeacher(userID) {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Davomat jarayoni bekor qilindi.\nQayta boshlash uchun 👇 tugmani bosing.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.MainKeyboard(),
		})
		return handlers.EndConversation()
	}

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Jarayon bekor qilindi.\nQayta boshlash uchun /start yuboring.", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
	})
	return handlers.EndConversation() 
}
