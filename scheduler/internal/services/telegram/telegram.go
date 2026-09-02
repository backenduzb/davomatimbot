package telegram

import (
	"bytes"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Sender xizmatlarining yuboruvchi (DocumentSender) interfeysini
// Telegram boti orqali amalga oshiradi. Bitta hujjat yuborish uchun
// kerak bo'lgan minimal funksionallik — komandalar, polling yoki
// handlerlar yo'q.
type Sender struct {
	bot *gotgbot.Bot
}

// NewSender bot uchun yuboruvchi yaratadi.
func NewSender(bot *gotgbot.Bot) *Sender {
	return &Sender{bot: bot}
}

// SendDocument fayl ma'lumotini xotiradan (vaqtincha fayl yaratmasdan)
// hujjat sifatida so'ralgan chatga yuboradi.
func (s *Sender) SendDocument(chatID int64, fileName string, data []byte, caption string) error {
	if s.bot == nil {
		return fmt.Errorf("telegram: bot yo'q")
	}
	_, err := s.bot.SendDocument(chatID,
		gotgbot.InputFileByReader(fileName, bytes.NewReader(data)),
		&gotgbot.SendDocumentOpts{
			Caption: caption,
		},
	)
	if err != nil {
		return fmt.Errorf("telegram: hujjat yuborilmadi: %w", err)
	}
	return nil
}
