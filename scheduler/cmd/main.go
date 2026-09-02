package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"scheduler/config/settings"
	"scheduler/internal/database"
	"scheduler/internal/repository/attendance"
	"scheduler/internal/repository/schedulerstate"
	"scheduler/internal/services/scheduler"
	"scheduler/internal/services/telegram"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	settings.LoadEnv()

	// Serverning lokal soat sohasiga ishonmaymiz — soat sohasi aniq
	// yuklanadi (default: Asia/Tashkent).
	loc, err := time.LoadLocation(settings.Envs.TIMEZONE)
	if err != nil {
		log.Fatalf("timezone %q yuklab bo'lmadi: %v", settings.Envs.TIMEZONE, err)
	}

	log.Println("scheduler started")
	log.Printf("timezone: %s", settings.Envs.TIMEZONE)

	database.Connect(settings.Envs.DB_URL)

	chatID, err := strconv.ParseInt(settings.Envs.CHAT_ID, 10, 64)
	if err != nil {
		log.Fatalf("TELEGRAM_CHAT_ID yaroqsiz: %q", settings.Envs.CHAT_ID)
	}

	bot, err := gotgbot.NewBot(settings.Envs.BOT_TOKEN, nil)
	if err != nil {
		log.Fatalf("telegram bot yaratilmadi: %v", err)
	}

	s := scheduler.New(
		loc,
		chatID,
		telegram.NewSender(bot),
		attendance.ListForDate,
		schedulerstate.Default,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := s.Run(ctx); err != nil {
		log.Fatalf("scheduler xato: %v", err)
	}
}
