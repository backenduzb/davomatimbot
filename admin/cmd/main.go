package main

import (
	"github.com/gin-gonic/gin"
	"admin/internal/routes"
	"admin/internal/database"
	"admin/internal/middleware/cors"
	"admin/config/settings"
)

func main() {

	settings.LoadEnv()
	router := gin.Default()
	router.Use(cors.Middleware())
	
	database.Connect(settings.Envs.DB_URL)
	routes.SetupAuthRoutes(router, database.DB)
	router.Run(settings.Envs.PORT)
}