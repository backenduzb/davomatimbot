package admin

import (
	"fmt"
	"strings"

	botHandlers "bot/internal/handlers"
	"bot/internal/repository/classes"
	"bot/internal/repository/sessions"
	"bot/internal/repository/settings"
	"bot/internal/repository/states"
	"bot/internal/services/filters"
	replyKeyboards "bot/internal/services/keyboards/reply"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// HandleAdminMenuChoice admin birinchi menyudan tugma tanlaganda ishlaydi:
// "Davomat topshirish" -> sinf tanlash oqimi,
// "Hamma ustozlarga xabar berish" -> broadcast matnini kutish.
func HandleAdminMenuChoice(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	choice := ""
	if ctx.EffectiveMessage != nil {
		choice = strings.TrimSpace(ctx.EffectiveMessage.Text)
	}

	session := sessions.GetSession(userID)
	if session == nil || !session.IsAdmin {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start yuboring.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
		})
		return handlers.EndConversation()
	}

	if replyKeyboards.IsBroadcast(choice) {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "📢 Barcha ustozlarga yuboriladigan xabarni yozing yoki sticker/GIF yuboring.\n\nBekor qilish uchun /cancel buyrug'ini yuboring.", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
		})
		return handlers.NextConversationState(states.StateWaitingAdminBroadcast)
	}

	if replyKeyboards.IsDavomatTopshirish(choice) {
		return botHandlers.StartAdminAttendance(b, ctx)
	}

	if replyKeyboards.IsCongratsToggle(choice) {
		return handleCongratsToggle(b, ctx, session)
	}

	_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Iltimos, quyidagi tugmalardan birini tanlang.", &gotgbot.SendMessageOpts{
		ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
	})
	return handlers.NextConversationState(states.StateWaitingAdminMenu)
}

// handleCongratsToggle tabrik funksiyasini yoqadi/o'chiradi.
// Bu tugma faqat superadmin (bot egasi) uchun ishlaydi.
func handleCongratsToggle(b *gotgbot.Bot, ctx *ext.Context, session *sessions.AttendanceSession) error {
	userID := uint(ctx.EffectiveUser.Id)
	if !filters.IsSuperAdmin(userID) {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⛔ Bu tugma faqat bot egasi uchun.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
		})
		return handlers.NextConversationState(states.StateWaitingAdminMenu)
	}

	enabled := !settings.IsGuestCongratsEnabled()
	if err := settings.SetGuestCongratsEnabled(enabled); err != nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "❌ Sozlamani saqlashda xatolik yuz berdi.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
		})
		return handlers.NextConversationState(states.StateWaitingAdminMenu)
	}

	status := "✅ Yoqildi"
	if !enabled {
		status = "❌ O'chirildi"
	}
	_, _ = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("🎉 Tabrik funksiyasi: <b>%s</b>", status), &gotgbot.SendMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
	})
	return handlers.NextConversationState(states.StateWaitingAdminMenu)
}

// HandleAdminBroadcastInput admin xabarini (matn, sticker, GIF/animatsiya,
// foto, video va boshqalar) qabul qilib barcha ustozlarga yuboradi.
// Xabar nusxasi CopyMessage orqali yuboriladi — format saqlanib qoladi.
func HandleAdminBroadcastInput(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	session := sessions.GetSession(userID)
	if session == nil || !session.IsAdmin {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sessiya topilmadi. /start yuboring.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
		})
		return handlers.EndConversation()
	}

	msg := ctx.EffectiveMessage
	if msg == nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Xabarni o'qib bo'lmadi. Matn, sticker yoki GIF yuboring.", nil)
		return handlers.NextConversationState(states.StateWaitingAdminBroadcast)
	}

	teacherIDs := classes.GetAllTeacherChatIDs()
	if len(teacherIDs) == 0 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Tizimda ustozlar topilmadi.", &gotgbot.SendMessageOpts{
			ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
		})
		return handlers.NextConversationState(states.StateWaitingAdminMenu)
	}

	fromChatID := ctx.EffectiveChat.Id
	messageID := msg.MessageId

	sent := 0
	failed := 0
	for _, teacherID := range teacherIDs {
		_, err := b.CopyMessage(teacherID, fromChatID, messageID, nil)
		if err != nil {
			failed++
			continue
		}
		sent++
	}

	report := fmt.Sprintf("✅ Xabar yuborildi.\n\n👥 Jami ustozlar: <b>%d</b>\n📤 Yuborildi: <b>%d</b>", len(teacherIDs), sent)
	if failed > 0 {
		report += fmt.Sprintf("\n❌ Yuborilmadi: <b>%d</b>", failed)
	}
	_, _ = b.SendMessage(ctx.EffectiveChat.Id, report, &gotgbot.SendMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: replyKeyboards.AdminMenuKeyboard(),
	})
	return handlers.NextConversationState(states.StateWaitingAdminMenu)
}
