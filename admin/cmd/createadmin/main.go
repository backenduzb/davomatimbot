package main

import (
	"flag"
	"fmt"
	"os"

	"admin/config/settings"
	"admin/internal/database"
	"admin/internal/services/auth"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ishlatish tartibi: %s create-admin -username=... -password=... [-telegram-id=...]\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create-admin":
		adminCmd := flag.NewFlagSet("create-admin", flag.ContinueOnError)
		adminUsername := adminCmd.String("username", "", "Admin foydalanuvchi nomi")
		adminPassword := adminCmd.String("password", "", "Admin paroli")
		telegramID := adminCmd.String("telegram-id", "", "Adminning Telegram ID'si (bot uni shu orqali taniydi)")

		if err := adminCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Argumentlarni o'qishda xatolik:", err)
			os.Exit(1)
		}
		if *adminUsername == "" || *adminPassword == "" {
			fmt.Println("Xatolik: username va password majburiy!")
			adminCmd.Usage()
			os.Exit(1)
		}

		// Bazaga ulanish ma'lumotlarini .env faylidan olamiz
		settings.LoadEnv()
		if settings.Envs.DB_URL == "" {
			fmt.Println("Xatolik: .env faylida DATABASE_URL topilmadi!")
			os.Exit(1)
		}

		database.Connect(settings.Envs.DB_URL)

		if err := auth.CreateAdminUser(*adminUsername, *adminPassword, *telegramID); err != nil {
			fmt.Println("Admin yaratishda xatolik:", err)
			os.Exit(1)
		}

		fmt.Printf("Muvaffaqiyatli yaratildi! username: %s 🎉\n", *adminUsername)
		if *telegramID != "" {
			fmt.Printf("Telegram ID biriktirildi: %s\n", *telegramID)
		} else {
			fmt.Println("Ogohlantirish: telegram-id berilmadi — bot adminni taniy olmaydi. Uni keyinroq /users orqali yangilang.")
		}

	default:
		flag.Usage()
		os.Exit(1)
	}
}
