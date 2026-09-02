package settings

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	DefaultTimezone = "Asia/Tashkent"
	DefaultChatID   = "-1003013595617"
)

type Config struct {
	DB_URL    string
	BOT_TOKEN string
	CHAT_ID   string
	TIMEZONE  string
}

var Envs Config

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		timezone = DefaultTimezone
	}

	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID == "" {
		chatID = DefaultChatID
	}

	Envs = Config{
		DB_URL:    os.Getenv("DATABASE_URL"),
		BOT_TOKEN: os.Getenv("TELEGRAM_BOT_TOKEN"),
		CHAT_ID:   chatID,
		TIMEZONE:  timezone,
	}
}
