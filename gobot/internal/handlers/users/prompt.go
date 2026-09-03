package users

import (
	"log"

	"bot/internal/models"
	"bot/internal/repository/sessions"
	"bot/internal/repository/states"
	"bot/internal/services/keyboards/inline"
	replyKeyboards "bot/internal/services/keyboards/reply"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// promptStudentChoice o'quvchi tanlash uchun reply klaviatura yuboradi.
//
// Muhim: agar tanlanadigan o'quvchi qolmagan bo'lsa (sinf bo'sh yoki hamma
// allaqachon belgilangan), Telegram BO'SH klaviaturani e'tiborsiz qoldiradi va
// foydalanuvchiga faqat matn ko'rinadi — "tugmalar kelmadi" holati aynan shu.
// Shuning uchun bunday holatda tushunarli xabar va keyingi qadam tugmalari
// yuboriladi.
func promptStudentChoice(
	b *gotgbot.Bot,
	ctx *ext.Context,
	session *sessions.AttendanceSession,
	allStudents []models.Student,
	promptText string,
	nextState string,
) error {
	if len(allStudents) == 0 {
		_, err := b.SendMessage(ctx.EffectiveChat.Id,
			"⚠️ Bu sinfda o'quvchilar topilmadi. Avval o'quvchilarni tizimga qo'shing.",
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
			})
		if err != nil {
			log.Printf("promptStudentChoice: bo'sh sinf xabarini yuborib bo'lmadi: %v", err)
		}
		return handlers.EndConversation()
	}

	keyboard := replyKeyboards.FilteredStudentsKeyboard(allStudents, session)

	if !replyKeyboards.HasButtons(keyboard) {
		// Barcha o'quvchilar allaqachon belgilangan — keyingi bosqichga o'tamiz.
		_, err := b.SendMessage(ctx.EffectiveChat.Id,
			"ℹ️ Barcha o'quvchilar allaqachon belgilangan.\n"+session.GenerateInfoText()+
				"\nKeyingi qadamni tanlang:",
			&gotgbot.SendMessageOpts{
				ParseMode:   "Markdown",
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
			})
		if err != nil {
			log.Printf("promptStudentChoice: xabarni yuborib bo'lmadi: %v", err)
		}

		_, err = b.SendMessage(ctx.EffectiveChat.Id, "Davomatni saqlaymizmi?", &gotgbot.SendMessageOpts{
			ReplyMarkup: inline.LateConfirmKeyboard(),
		})
		if err != nil {
			log.Printf("promptStudentChoice: tasdiq tugmalarini yuborib bo'lmadi: %v", err)
		}
		return handlers.NextConversationState(states.StateWaitingLateConfirm)
	}

	// Xato endi yutilmaydi — klaviatura yuborilmasa loglanadi.
	if _, err := b.SendMessage(ctx.EffectiveChat.Id, promptText, &gotgbot.SendMessageOpts{
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("promptStudentChoice: klaviaturani yuborib bo'lmadi: %v", err)
		return err
	}

	return handlers.NextConversationState(nextState)
}
