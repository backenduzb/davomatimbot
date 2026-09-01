package routes

import (
	attendanceHandlers "admin/internal/handlers/attendance"
	handlers "admin/internal/handlers/auth"
	importerHandlers "admin/internal/handlers/importer"
	studentsHandlers "admin/internal/handlers/students"
	usersHandlers "admin/internal/handlers/users"
	statisticsHandlers "admin/internal/handlers/statistics"
	"admin/internal/middleware/auth"
	"admin/internal/models"
	"admin/internal/schemas/class"
	"admin/internal/schemas/student"
	"admin/internal/services/crud"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)
	// Foydalanuvchilar uchun maxsus handlerlar: parolni xeshlash,
	// duplicate username xatolarini to'g'ri qaytarish va bool maydonlarni
	// (is_admin, is_banned, ...) to'g'ri yangilash uchun.
	api.GET("/users", auth.AuthMiddleware(), usersHandlers.List)
	api.POST("/users", auth.AuthMiddleware(), usersHandlers.Create)
	api.GET("/users/:id", auth.AuthMiddleware(), usersHandlers.Retrieve)
	api.PUT("/users/:id", auth.AuthMiddleware(), usersHandlers.Update)
	api.PATCH("/users/:id", auth.AuthMiddleware(), usersHandlers.Update)
	api.DELETE("/users/:id", auth.AuthMiddleware(), usersHandlers.Delete)

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

	// O'quvchilar: ro'yxatni class_id/search bo'yicha filtrlash uchun maxsus
	// List handler, qolgan amallar (create/update/delete) generic CRUD'da.
	studentsCtrl := crud.NewCRUDController[models.Student, student.Student, student.Student](db)
	api.GET("/students", auth.AuthMiddleware(), studentsHandlers.List)
	api.POST("/students", auth.AuthMiddleware(), studentsCtrl.Create)
	api.GET("/students/:id", auth.AuthMiddleware(), studentsCtrl.Retrieve)
	api.PUT("/students/:id", auth.AuthMiddleware(), studentsCtrl.Update)
	api.PATCH("/students/:id", auth.AuthMiddleware(), studentsCtrl.PartialUpdate)
	api.DELETE("/students/:id", auth.AuthMiddleware(), studentsCtrl.Delete)

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
	api.POST("/import/xlsx",
		auth.AuthMiddleware(),
		importerHandlers.UploadXLSX,
	)
}
