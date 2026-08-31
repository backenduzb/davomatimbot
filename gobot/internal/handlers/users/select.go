package users

import (
	"fmt"

	"bot/internal/repository/sessions"
	"bot/internal/repository/states"
	"bot/internal/repository/students"
	"bot/internal/services/keyboards/inline"
	replyKeyboards "bot/internal/services/keyboards/reply"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleAbsentStudentSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	studentName := ctx.Message.Text
	session := sessions.GetSession(userID)
	if session == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", nil)
		return handlers.EndConversation()
	}

	studentID := students.GetStudentIDByName(studentName, userID)

	if studentID == 0 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Bunday o'quvchi topilmadi. Iltimos, tugmalardan foydalaning!", nil)
		return handlers.NextConversationState(states.StateWaitingAbsentStudent)
	}

	if session.IsAlreadyMarked(studentID) {
		allStudents := students.GetAllStudents(userID)
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Bu o'quvchi allaqachon kiritilgan! Boshqasini tanlang:", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.FilteredStudentsKeyboard(allStudents, session),
		})
		return handlers.NextConversationState(states.StateWaitingAbsentStudent)
	}

	session.UnexcusedStudents[studentID] = studentName

	report := session.GenerateInfoText()
	msgText := fmt.Sprintf("✅ **%s** sababsizlar ro'yxatiga qo'shildi.\n%s\nKeyingi qadamni tanlang:", studentName, report)

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, msgText, &gotgbot.SendMessageOpts{
		ParseMode:   "Markdown",
		ReplyMarkup: inline.AbsentConfirmKeyboard(),
	})

	return handlers.NextConversationState(states.StateWaitingAbsentConfirm)
}

func HandleReasonStudentSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	studentName := ctx.Message.Text
	session := sessions.GetSession(userID)
	if session == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", nil)
		return handlers.EndConversation()
	}

	studentID := students.GetStudentIDByName(studentName, userID)
	if studentID == 0 {
		allStudents := students.GetAllStudents(userID)
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Bunday o'quvchi topilmadi. Iltimos, tugmalardan foydalaning!", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.FilteredStudentsKeyboard(allStudents, session),
		})
		return handlers.NextConversationState(states.StateWaitingReasonStudent)
	}

	if session.IsAlreadyMarked(studentID) {
		allStudents := students.GetAllStudents(userID)
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Bu o'quvchi allaqachon kiritilgan! Boshqasini tanlang:", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.FilteredStudentsKeyboard(allStudents, session),
		})
		return handlers.NextConversationState(states.StateWaitingReasonStudent)
	}

	session.LastSelectedStudentID = studentID
	session.LastSelectedStudentName = studentName

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("📝 **%s** nima sababdan kelmagan? Sababini matn ko'rinishida yozib yuboring:", studentName), &gotgbot.SendMessageOpts{
		ParseMode: "Markdown",
	})
	return handlers.NextConversationState(states.StateWaitingReasonInput)
}