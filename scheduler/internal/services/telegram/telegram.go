package telegram

import (
	"bytes"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Sender struct {
	bot *gotgbot.Bot
}

func NewSender(bot *gotgbot.Bot) *Sender {
	return &Sender{bot: bot}
}

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

func (s *Sender) SendMessage(chatID int64, text string) error {
	if s.bot == nil {
		return fmt.Errorf("telegram: bot yo'q")
	}
	_, err := s.bot.SendMessage(chatID, text, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("telegram: xabar yuborilmadi: %w", err)
	}
	return nil
}
