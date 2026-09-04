package handlers

import (
	"net/http"

	"admin/internal/database"
	"admin/internal/models"
	"admin/internal/services/auth"

	"github.com/gin-gonic/gin"
)

// MinPasswordLength — users.Update dagi qoida bilan bir xil.
const MinPasswordLength = 4

// ChangePasswordRequest — foydalanuvchi o'z parolini o'zgartirish so'rovi.
// Joriy parol majburiy: token o'g'irlangan bo'lsa ham parolni almashtirib
// bo'lmasligi uchun.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// ChangePassword — POST /api/me/password.
//
// Faqat so'rov yuborgan foydalanuvchining O'Z parolini almashtiradi
// (user_id token'dan olinadi, tanadan emas — boshqa hisobni nishonga olish
// mumkin emas).
func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Avtorizatsiya talab qilinadi"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Joriy va yangi parol majburiy"})
		return
	}

	user, err := auth.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Foydalanuvchi topilmadi"})
		return
	}

	if !auth.CheckPassword(user.Password, req.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Joriy parol noto'g'ri"})
		return
	}

	if len([]rune(req.NewPassword)) < MinPasswordLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Yangi parol kamida 4 ta belgidan iborat bo'lishi kerak",
		})
		return
	}

	if req.NewPassword == req.CurrentPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Yangi parol joriy paroldan farq qilishi kerak",
		})
		return
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parolni xeshlashda xatolik"})
		return
	}

	// Faqat parol ustunini yangilaymiz — boshqa maydonlarga (is_online,
	// last_seen va h.k.) tegmaslik uchun Save emas, Update ishlatiladi.
	if err := database.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("password", hash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parolni saqlashda xatolik"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parol muvaffaqiyatli o'zgartirildi"})
}
