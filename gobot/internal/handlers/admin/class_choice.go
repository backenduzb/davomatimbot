package admin

import (
	"fmt"
	"strconv"
	"strings"

	"bot/internal/handlers"
	"bot/internal/repository/sessions"
	"bot/internal/repository/states"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleAdminClassChoice(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	userID := uint(ctx.EffectiveUser.Id)
	_, _ = query.Answer(b, nil)

	if !strings.HasPrefix(query.Data, "class_") {
		return nil
	}

	classID, err := strconv.ParseUint(strings.TrimPrefix(query.Data, "class_"), 10, 64)
	if err != nil || classID == 0 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Noto'g'ri sinf tanlandi.", nil)
		return handlers.NextConversationState(states.StateWaitingAdminClassChoice)
	}

	session := sessions.GetSession(userID)
	if session == nil || !session.IsAdmin {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", nil)
		return handlers.EndConversation()
	}

	session.SetClassContext(userID, uint(classID))

	className := session.ClassInfo[userID].ClassName
	_, _ = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("✅ <b>%s</b> sinfi tanlandi.", className), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return handlers.BeginAttendanceFlow(b, ctx)
}
