package users

import (
	"fmt"

	"bot/internal/repository/sessions"
	"bot/internal/repository/states"
	"bot/internal/services/keyboards/inline"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleReasonInput(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	reasonText := ctx.Message.Text
	session := sessions.GetSession(userID)
	if session == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", nil)
		return handlers.EndConversation()
	}

	session.ExcusedStudents[session.LastSelectedStudentID] = sessions.ExcusedDetail{
		Name:   session.LastSelectedStudentName,
		Reason: reasonText,
	}

	report := session.GenerateInfoText()
	msgText := fmt.Sprintf("✅ Sababli kelmagan o'quvchi qayd etildi.\n%s\nDavom etamizmi?", report)

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, msgText, &gotgbot.SendMessageOpts{
		ParseMode:   "Markdown",
		ReplyMarkup: inline.ReasonConfirmKeyboard(),
	})

	return handlers.NextConversationState(states.StateWaitingReasonConfirm)
}