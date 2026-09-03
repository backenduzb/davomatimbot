package admin

import (
	"fmt"
	"strings"

	botHandlers "bot/internal/handlers"
	"bot/internal/models"
	"bot/internal/repository/classes"
	"bot/internal/repository/sessions"
	"bot/internal/repository/states"
	replyKeyboards "bot/internal/services/keyboards/reply"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// HandleAdminClassChoice admin reply klaviaturadan SINF NOMINI tanlaganda
// ishlaydi. Tugmalarda faqat sinf nomi bo'lgani uchun bu yerda nom bo'yicha
// mos sinf(lar) qidiriladi.
func HandleAdminClassChoice(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	choice := ""
	if ctx.EffectiveMessage != nil {
		choice = strings.TrimSpace(ctx.EffectiveMessage.Text)
	}

	session := sessions.GetSession(userID)
	if session == nil || !session.IsAdmin {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
		})
		return handlers.EndConversation()
	}

	classList := classes.GetAllClasses()
	matched := matchClassesByName(classList, choice)

	if len(matched) == 0 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Bunday sinf topilmadi. Iltimos, quyidagi tugmalardan birini tanlang.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.ClassNamesKeyboard(classList),
		})
		return handlers.NextConversationState(states.StateWaitingAdminClassChoice)
	}

	// Bir xil nomli bir nechta sinf bo'lsa — o'qituvchi bo'yicha aniqlashtiramiz.
	if len(matched) > 1 {
		session.PendingClassName = choice
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("ℹ️ <b>%s</b> nomli bir nechta sinf bor.\nO'qituvchini tanlang:", choice), &gotgbot.SendMessageOpts{
			ParseMode:   "HTML",
			ReplyMarkup: replyKeyboards.TeacherChoiceKeyboard(matched),
		})
		return handlers.NextConversationState(states.StateWaitingAdminTeacherChoice)
	}

	return applyClassChoice(b, ctx, session, matched[0])
}

// HandleAdminTeacherChoice bir xil nomli sinflar orasidan o'qituvchi bo'yicha
// aniq sinfni tanlash uchun ishlatiladi.
func HandleAdminTeacherChoice(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	choice := ""
	if ctx.EffectiveMessage != nil {
		choice = strings.TrimSpace(ctx.EffectiveMessage.Text)
	}

	session := sessions.GetSession(userID)
	if session == nil || !session.IsAdmin {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start buyrug'ini qayta yuboring.", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
		})
		return handlers.EndConversation()
	}

	classList := classes.GetAllClasses()
	matched := matchClassesByName(classList, session.PendingClassName)

	for _, class := range matched {
		if strings.EqualFold(strings.TrimSpace(class.TeacherFullName), choice) {
			session.PendingClassName = ""
			return applyClassChoice(b, ctx, session, class)
		}
	}

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Bunday o'qituvchi topilmadi. Iltimos, tugmalardan birini tanlang.", &gotgbot.SendMessageOpts{
		ReplyMarkup: replyKeyboards.TeacherChoiceKeyboard(matched),
	})
	return handlers.NextConversationState(states.StateWaitingAdminTeacherChoice)
}

// applyClassChoice tanlangan sinfni sessiyaga yozadi va davomat jarayonini
// boshlaydi.
func applyClassChoice(b *gotgbot.Bot, ctx *ext.Context, session *sessions.AttendanceSession, class models.Class) error {
	userID := uint(ctx.EffectiveUser.Id)
	session.SetClassContext(userID, class.ID)

	detail, ok := session.ClassInfo[userID]
	if !ok || detail.ClassName == "" || detail.ClassName == "-" {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Tanlangan sinf topilmadi. Iltimos, boshqa sinfni tanlang.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.ClassNamesKeyboard(classes.GetAllClasses()),
		})
		return handlers.NextConversationState(states.StateWaitingAdminClassChoice)
	}

	// Sinf tanlangach reply klaviatura olib tashlanadi — keyingi qadamda
	// inline tugmalar ishlatiladi.
	_, _ = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("✅ <b>%s</b> sinfi tanlandi.", detail.ClassName), &gotgbot.SendMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
	})

	return botHandlers.BeginAttendanceFlow(b, ctx)
}

// matchClassesByName sinf nomi bo'yicha mos sinflarni qaytaradi
// (registrga sezgir emas).
func matchClassesByName(classList []models.Class, name string) []models.Class {
	target := strings.TrimSpace(name)
	if target == "" {
		return nil
	}
	matched := make([]models.Class, 0, 2)
	for _, class := range classList {
		if strings.EqualFold(strings.TrimSpace(class.ClassName.Name), target) {
			matched = append(matched, class)
		}
	}
	return matched
}
