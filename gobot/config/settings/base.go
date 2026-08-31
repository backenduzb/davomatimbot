package settings

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL string
	PORT string
	DEBUG string
	BOT_TOKEN string
	WEBHOOK_SECRET string
	WEBHOOK_URL string
}

var Envs Config

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	Envs = Config{
		DB_URL: os.Getenv("DATABASE_URL"),
		PORT: os.Getenv("PORT"),
		WEBHOOK_SECRET: os.Getenv("JWT_SECRET"),
		BOT_TOKEN: os.Getenv("BOT_TOKEN"),
		DEBUG: os.Getenv("DEBUG"),
		WEBHOOK_URL: os.Getenv("WEBHOOK_URL"),
	}
}
