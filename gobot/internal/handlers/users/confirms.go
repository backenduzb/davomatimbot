package users

import (
	"bot/internal/repository/attendance"
	"bot/internal/repository/sessions"
	"bot/internal/repository/states"
	"bot/internal/repository/students"
	"bot/internal/services/keyboards/inline"
	replyKeyboards "bot/internal/services/keyboards/reply"
	"bot/utils"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleAbsentConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	userID := uint(ctx.EffectiveUser.Id)
	_, _ = query.Answer(b, nil)

	session := sessions.GetSession(userID)
	if session == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", nil)
		return handlers.EndConversation()
	}

	allStudents := students.GetAllStudents(userID)

	if query.Data == "add_more_unexcused" {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Keyingi sababsiz kelmagan o'quvchini tanlang:", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.FilteredStudentsKeyboard(allStudents, session),
		})
		return handlers.NextConversationState(states.StateWaitingAbsentStudent)
	}

	if query.Data == "go_to_excused" {
		report := session.GenerateInfoText()
		msgText := fmt.Sprintf("📝 Sababsiz kelmaganlar ro'yxati shakllandi.\n%s\n\nEndi sababli kelmagan o'quvchilar bormi?", report)

		_, _ = b.SendMessage(ctx.EffectiveChat.Id, msgText, &gotgbot.SendMessageOpts{
			ParseMode:   "Markdown",
			ReplyMarkup: inline.ReasonConfirmKeyboard(),
		})

		return handlers.NextConversationState(states.StateWaitingReasonConfirm)
	}
	return nil
}

func HandleReasonConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	userID := uint(ctx.EffectiveUser.Id)
	_, _ = query.Answer(b, nil)

	session := sessions.GetSession(userID)
	if session == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", nil)
		return handlers.EndConversation()
	}

	if query.Data == "add_more_excused" {
		allStudents := students.GetAllStudents(userID)
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Sababli kelmagan o'quvchini tanlang:", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.FilteredStudentsKeyboard(allStudents, session),
		})
		return handlers.NextConversationState(states.StateWaitingReasonStudent)
	}

	if query.Data == "save_attendance" {
		classID := session.ResolveClassID(userID)
		if classID == 0 {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sinf aniqlanmadi. /start buyrug'ini qayta yuboring.", nil)
			return handlers.EndConversation()
		}

		if err := attendance.SaveClassAttendance(classID, session); err != nil {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, "❌ Davomatni saqlashda xatolik yuz berdi. Qayta urinib ko'ring.", nil)
			return handlers.NextConversationState(states.StateWaitingReasonConfirm)
		}

		reportText := utils.GenerateFullReport(userID, true)

		_, err := b.SendMessage(ctx.EffectiveChat.Id, reportText, &gotgbot.SendMessageOpts{
			ParseMode:   "HTML",
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
		})
		if err != nil {
			return err
		}

		sessions.DeleteSession(userID)
		return handlers.EndConversation()
	}
	return nil
}
