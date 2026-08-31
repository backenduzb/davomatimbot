package main

import (
	"bot/config/settings"
	"bot/internal/database"
	"bot/internal/routes"
	"bot/utils"
)

func main() {
	settings.LoadEnv()
	database.Connect(settings.Envs.DB_URL)
	bot, dispatcher := utils.BotInit()
	
	routes.RegisterSimpleHandler(dispatcher)
	
	if settings.Envs.DEBUG == "true" {
		utils.StartPolling(bot, dispatcher)
	} else {
		utils.StartWebhook(bot, dispatcher)
	}
}

