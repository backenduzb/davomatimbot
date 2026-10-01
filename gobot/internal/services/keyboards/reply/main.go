package replyKeyboards

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// BtnDavomatTopshirish — davomat topshirishni qayta boshlash uchun asosiy
// reply tugma. /start yozib o'tirishni oldini oladi.
const BtnDavomatTopshirish = "📋 Davomat topshirish"

// BtnBroadcast — admin barcha ustozlarga xabar yuborish tugmasi.
const BtnBroadcast = "📢 Hamma ustozlarga xabar berish"

// BtnCongratsToggle — tabrik funksiyasini yoqish/o'chirish tugmasi.
// Faqat superadmin (bot egasi) bosa oladi.
const BtnCongratsToggle = "🎉 Tabrik funksiyasi"

// MainKeyboard yakuniy hisobotdan keyin (va bekor qilinganda) ko'rsatiladigan
// asosiy reply klaviatura — bitta "Davomat topshirish" tugmasi.
func MainKeyboard() gotgbot.ReplyKeyboardMarkup {
	return gotgbot.ReplyKeyboardMarkup{
		Keyboard: [][]gotgbot.KeyboardButton{
			{{Text: BtnDavomatTopshirish}},
		},
		ResizeKeyboard:        true,
		OneTimeKeyboard:       false,
		InputFieldPlaceholder: "Davomat topshirish",
	}
}

// IsDavomatTopshirish matn asosiy tugmaga mos kelishini tekshiradi.
func IsDavomatTopshirish(text string) bool {
	return strings.TrimSpace(text) == BtnDavomatTopshirish
}

// AdminMenuKeyboard admin /start bosganda chiqadigan birinchi menyu:
// davomat topshirish + broadcast + tabrik funksiyasi togglesi.
func AdminMenuKeyboard() gotgbot.ReplyKeyboardMarkup {
	return gotgbot.ReplyKeyboardMarkup{
		Keyboard: [][]gotgbot.KeyboardButton{
			{{Text: BtnDavomatTopshirish}},
			{{Text: BtnBroadcast}},
			{{Text: BtnCongratsToggle}},
		},
		ResizeKeyboard:        true,
		OneTimeKeyboard:       false,
		InputFieldPlaceholder: "Menyuni tanlang",
	}
}

// IsBroadcast matn broadcast tugmasiga mos kelishini tekshiradi.
func IsBroadcast(text string) bool {
	return strings.TrimSpace(text) == BtnBroadcast
}

// IsCongratsToggle matn tabrik funksiyasi tugmasiga mos kelishini tekshiradi.
func IsCongratsToggle(text string) bool {
	return strings.TrimSpace(text) == BtnCongratsToggle
}
