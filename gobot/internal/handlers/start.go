package handlers

import (
	"bot/internal/repository/classes"
	"bot/internal/repository/sessions"
	"bot/internal/repository/statistics"
	"bot/internal/repository/states"
	"bot/internal/services/filters"
	"bot/internal/services/keyboards/inline"
	replyKeyboards "bot/internal/services/keyboards/reply"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HandleStart(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := uint(ctx.EffectiveUser.Id)
	sessions.NewSession(userID)
	session := sessions.GetSession(userID)

	if filters.CheckIsAdmin(userID) {
		session.IsAdmin = true
		stats := statistics.GetTodayOverview()
		message := statistics.FormatAdminOverviewMessage(stats)
		classList := classes.GetAllClasses()

		if len(classList) == 0 {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, message+"\n\n⚠️ Sistemada sinflar topilmadi.", &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			})
			sessions.DeleteSession(userID)
			return handlers.EndConversation()
		}

		// Admin sinfni reply klaviaturadan tanlaydi: tugmalarda faqat
		// sinf nomlari ko'rsatiladi.
		classKeyboard := replyKeyboards.ClassNamesKeyboard(classList)
		if len(classKeyboard.Keyboard) == 0 {
			_, _ = b.SendMessage(ctx.EffectiveChat.Id, message+"\n\n⚠️ Sinflarga nom biriktirilmagan.", &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			})
			sessions.DeleteSession(userID)
			return handlers.EndConversation()
		}

		_, err := b.SendMessage(ctx.EffectiveChat.Id, message+"\n\n👇 Sinfni tanlang:", &gotgbot.SendMessageOpts{
			ParseMode:   "HTML",
			ReplyMarkup: classKeyboard,
		})
		if err != nil {
			return err
		}
		return handlers.NextConversationState(states.StateWaitingAdminClassChoice)
	}

	if !filters.CheckIsTeacher(userID) {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sizda davomat topshirish huquqi yo'q.", nil)
		sessions.DeleteSession(userID)
		return handlers.EndConversation()
	}

	classID := classes.GetClassID(userID)
	if classID == 0 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "⚠️ Sizga biriktirilgan sinf topilmadi.", nil)
		sessions.DeleteSession(userID)
		return handlers.EndConversation()
	}

	session.SetClassContext(userID, classID)
	return BeginAttendanceFlow(b, ctx)
}

func BeginAttendanceFlow(b *gotgbot.Bot, ctx *ext.Context) error {
	session := sessions.GetSession(uint(ctx.EffectiveUser.Id))
	className := "sinf"
	if detail, ok := session.ClassInfo[uint(ctx.EffectiveUser.Id)]; ok && detail.ClassName != "" {
		className = detail.ClassName
	}

	text := fmt.Sprintf("📋 <b>%s</b> uchun davomat jarayoni boshlandi.\nIltimos, tanlang:", className)
	_, err := b.SendMessage(ctx.EffectiveChat.Id, text, &gotgbot.SendMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: inline.AbsentTypeKeyboard(),
	})
	if err != nil {
		return err
	}
	return handlers.NextConversationState(states.StateWaitingAbsentTypeChoice)
}
