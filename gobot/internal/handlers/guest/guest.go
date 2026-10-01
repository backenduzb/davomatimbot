package guest

import (
	"fmt"
	"strings"

	"bot/internal/repository/classes"
	"bot/internal/repository/sessions"
	"bot/internal/repository/states"
	replyKeyboards "bot/internal/services/keyboards/reply"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func StartGuestFlow(b *gotgbot.Bot, ctx *ext.Context) error {
	teachers := classes.GetAllTeachers()
	if len(teachers) == 0 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Hozircha ustozlar ro'yxati topilmadi.", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
		})
		sessions.DeleteSession(uint(ctx.EffectiveUser.Id))
		return handlers.EndConversation()
	}

	names := make([]string, 0, len(teachers))
	for _, t := range teachers {
		names = append(names, t.FullName)
	}
	_, _ = b.SendMessage(ctx.EffectiveChat.Id, "👋 Assalomu alaykum!\nTabriklamoqchi bo'lgan ustozingizni tanlang:", &gotgbot.SendMessageOpts{
		ReplyMarkup: replyKeyboards.TeachersKeyboard(names),
	})
	return handlers.NextConversationState(states.StateWaitingGuestTeacherChoice)
}

func HandleGuestTeacherChoice(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	choice := ""
	if ctx.EffectiveMessage != nil {
		choice = strings.TrimSpace(ctx.EffectiveMessage.Text)
	}

	session := sessions.GetSession(userID)
	if session == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. Qayta boshlash uchun /start yuboring.", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
		})
		return handlers.EndConversation()
	}

	teachers := classes.GetAllTeachers()
	teacher, ok := classes.FindTeacherByName(teachers, choice)
	if !ok {
		if len(teachers) == 0 {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Ustozlar ro'yxati topilmadi. Qayta boshlash uchun /start yuboring.", &gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
			})
			sessions.DeleteSession(userID)
			return handlers.EndConversation()
		}
		names := make([]string, 0, len(teachers))
		for _, t := range teachers {
			names = append(names, t.FullName)
		}
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Bunday ustoz topilmadi. Iltimos, tugmalardan birini tanlang.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.TeachersKeyboard(names),
		})
		return handlers.NextConversationState(states.StateWaitingGuestTeacherChoice)
	}

	session.GuestTeacherID = teacher.TelegramID
	session.GuestTeacherName = teacher.FullName

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("✍️ <b>%s</b> uchun tabrik matningizni yozing (matn, sticker yoki GIF yuborishingiz mumkin).\n\nBekor qilish uchun /cancel buyrug'ini yuboring.", teacher.FullName), &gotgbot.SendMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
	})
	return handlers.NextConversationState(states.StateWaitingGuestCongratsInput)
}

func senderLabel(ctx *ext.Context) string {
	u := ctx.EffectiveUser
	if u == nil {
		return "Mehmon"
	}
	name := strings.TrimSpace(strings.TrimSpace(u.FirstName) + " " + strings.TrimSpace(u.LastName))
	if u.Username != "" {
		if name != "" {
			return fmt.Sprintf("%s (@%s)", name, u.Username)
		}
		return "@" + u.Username
	}
	if name != "" {
		return name
	}
	return "Mehmon"
}

func HandleGuestCongratsInput(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	session := sessions.GetSession(userID)
	if session == nil || session.GuestTeacherID == 0 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. Qayta boshlash uchun /start yuboring.", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
		})
		return handlers.EndConversation()
	}

	msg := ctx.EffectiveMessage
	if msg == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Xabarni o'qib bo'lmadi. Tabrik matnini qayta yuboring.", nil)
		return handlers.NextConversationState(states.StateWaitingGuestCongratsInput)
	}

	teacherID := session.GuestTeacherID
	teacherName := session.GuestTeacherName
	from := senderLabel(ctx)

	if strings.TrimSpace(msg.Text) != "" {
		text := fmt.Sprintf("🎉 <b>Sizga tabrik keldi!</b>\n\n%s\n\n👤 <i>Yuboruvchi: %s</i>", msg.Text, from)
		if _, err := b.SendMessage(teacherID, text, &gotgbot.SendMessageOpts{ParseMode: "HTML"}); err != nil {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, "❌ Tabrikni yuborib bo'lmadi. Ustoz botni ishga tushirmagan bo'lishi mumkin.", &gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
			})
			sessions.DeleteSession(userID)
			return handlers.EndConversation()
		}
	} else {
		header := fmt.Sprintf("🎉 <b>Sizga tabrik keldi!</b>\n👤 <i>Yuboruvchi: %s</i>", from)
		if _, err := b.SendMessage(teacherID, header, &gotgbot.SendMessageOpts{ParseMode: "HTML"}); err != nil {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, "❌ Tabrikni yuborib bo'lmadi. Ustoz botni ishga tushirmagan bo'lishi mumkin.", &gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
			})
			sessions.DeleteSession(userID)
			return handlers.EndConversation()
		}
		if _, err := b.CopyMessage(teacherID, ctx.EffectiveChat.Id, msg.MessageId, nil); err != nil {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, "❌ Tabrik nusxasini yuborib bo'lmadi.", &gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
			})
			sessions.DeleteSession(userID)
			return handlers.EndConversation()
		}
	}

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("✅ Tabrigingiz <b>%s</b> ga yuborildi! Rahmat!", teacherName), &gotgbot.SendMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
	})
	sessions.DeleteSession(userID)
	return handlers.EndConversation()
}
