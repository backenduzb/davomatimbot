package utils

import (
	"bot/internal/repository/sessions"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	sessions.DeleteSession(userID) 
	
	_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Davomat jarayoni bekor qilindi.", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
	})
	return handlers.EndConversation() 
}
