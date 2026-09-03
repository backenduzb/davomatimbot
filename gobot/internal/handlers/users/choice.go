package users

import (
	"bot/internal/repository/sessions"
	"bot/internal/repository/states"
	"bot/internal/repository/students"
	"bot/internal/services/keyboards/inline"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleAbsentTypeChoice(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	userID := uint(ctx.EffectiveUser.Id)
	_, _ = query.Answer(b, nil)

	session := sessions.GetSession(userID)
	if session == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", nil)
		return handlers.EndConversation()
	}

	allStudents := students.GetAllStudents(userID)

	if query.Data == "no_unexcused" {
		msgText := "✅ Sababsiz kelmaganlar yo'q.\n\nSababli kelmagan o'quvchilar bormi?"
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, msgText, &gotgbot.SendMessageOpts{
			ReplyMarkup: inline.ReasonConfirmKeyboard(),
		})
		return handlers.NextConversationState(states.StateWaitingReasonConfirm)
	}

	if query.Data == "add_unexcused" {
		return promptStudentChoice(b, ctx, session, allStudents,
			"Sababsiz kelmagan o'quvchini tanlang:",
			states.StateWaitingAbsentStudent)
	}
	return nil
}
