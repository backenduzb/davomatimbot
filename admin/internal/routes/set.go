package routes

import (
	attendanceHandlers "admin/internal/handlers/attendance"
	handlers "admin/internal/handlers/auth"
	statisticsHandlers "admin/internal/handlers/statistics"
	"admin/internal/middleware/auth"
	"admin/internal/models"
	"admin/internal/schemas/class"
	"admin/internal/schemas/student"
	"admin/internal/schemas/user/response"
	"admin/internal/services/crud"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)
	crud.RegisterCRUDRoutes[
		models.User,
		response.ProfileResponse,
		response.ProfileResponse,
	](api, "/users", db, auth.AuthMiddleware())
	crud.RegisterCRUDRoutes[
		models.Class,
		class.Class,
		class.Class,
	](api, "/classes", db, auth.AuthMiddleware())
	crud.RegisterCRUDRoutes[
		models.ClassName,
		class.ClassName,
		class.ClassName,
	](api, "/class/names", db, auth.AuthMiddleware())

	crud.RegisterCRUDRoutes[
		models.Student,
		student.Student,
		student.Student,
	](api, "/students", db, auth.AuthMiddleware())

	api.GET("/me",
		auth.AuthMiddleware(),
		handlers.Profile,
	)

	api.GET("/statistics/today",
		auth.AuthMiddleware(),
		statisticsHandlers.GetToday,
	)
	api.GET("/attendance/today",
		auth.AuthMiddleware(),
		attendanceHandlers.GetToday,
	)
	api.POST("/attendance/batch",
		auth.AuthMiddleware(),
		attendanceHandlers.BatchUpsert,
	)
}
