package users

import (
	"net/http"
	"strconv"

	"admin/internal/database"
	"admin/internal/models"
	"admin/internal/schemas/crud"
	"admin/internal/schemas/user/response"
	services "admin/internal/services/response"
	"admin/internal/services/auth"

	"github.com/gin-gonic/gin"
)

// CreateRequest — yangi foydalanuvchi yaratish uchun so'rov.
type CreateRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	IsAdmin    bool   `json:"is_admin"`
	IsBanned   bool   `json:"is_banned"`
	TelegramId string `json:"telegram_id"`
}

// UpdateRequest — foydalanuvchini yangilash uchun so'rov.
// Pointer maydonlar faqat yuborilgan qiymatlarni yangilashga imkon beradi
// (false/bo'sh qiymatlar ham to'g'ri qo'llanadi).
type UpdateRequest struct {
	Username   *string `json:"username"`
	Password   *string `json:"password"`
	IsAdmin    *bool   `json:"is_admin"`
	IsBanned   *bool   `json:"is_banned"`
	IsOnline   *bool   `json:"is_online"`
	TelegramId *string `json:"telegram_id"`
}

// usernameTaken username allaqachon band ekanini tekshiradi (update paytida
// o'zini hisobga olmaydi).
func usernameTaken(username string, excludeID uint) bool {
	var count int64
	q := database.DB.Model(&models.User{}).Where("username = ?", username)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	q.Count(&count)
	return count > 0
}

func Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parol kamida 4 ta belgidan iborat bo'lishi kerak"})
		return
	}
	if usernameTaken(req.Username, 0) {
		c.JSON(http.StatusConflict, gin.H{"error": "Bu username allaqachon mavjud"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parolni xeshlashda xatolik"})
		return
	}

	user := models.User{
		Username:   req.Username,
		Password:   hash,
		IsAdmin:    req.IsAdmin,
		IsBanned:   req.IsBanned,
		TelegramId: req.TelegramId,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, services.NewProfileResponse(user))
}

// Update PUT va PATCH uchun ishlatiladi — faqat yuborilgan maydonlarni yangilaydi.
func Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Noto'g'ri id"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Foydalanuvchi topilmadi"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username != nil {
		if *req.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username bo'sh bo'lishi mumkin emas"})
			return
		}
		if usernameTaken(*req.Username, user.ID) {
			c.JSON(http.StatusConflict, gin.H{"error": "Bu username allaqachon mavjud"})
			return
		}
		user.Username = *req.Username
	}

	if req.Password != nil && *req.Password != "" {
		if len(*req.Password) < 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parol kamida 4 ta belgidan iborat bo'lishi kerak"})
			return
		}
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Parolni xeshlashda xatolik"})
			return
		}
		user.Password = hash
	}

	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}
	if req.IsBanned != nil {
		user.IsBanned = *req.IsBanned
	}
	if req.IsOnline != nil {
		user.IsOnline = *req.IsOnline
	}
	if req.TelegramId != nil {
		user.TelegramId = *req.TelegramId
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services.NewProfileResponse(user))
}

func List(c *gin.Context) {
	var params crud.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}

	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	var users []models.User
	offset := (params.Page - 1) * params.PageSize
	database.DB.Order("id asc").Limit(params.PageSize).Offset(offset).Find(&users)

	results := make([]response.ProfileResponse, 0, len(users))
	for _, u := range users {
		results = append(results, services.NewProfileResponse(u))
	}

	totalPages := int64(0)
	if count > 0 {
		totalPages = (count + int64(params.PageSize) - 1) / int64(params.PageSize)
	}

	c.JSON(http.StatusOK, crud.PaginatedResponse[response.ProfileResponse]{
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: int64(params.Page),
		PageSize:    int64(params.PageSize),
		Results:     results,
	})
}

func Retrieve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Noto'g'ri id"})
		return
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Foydalanuvchi topilmadi"})
		return
	}
	c.JSON(http.StatusOK, services.NewProfileResponse(user))
}

func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Noto'g'ri id"})
		return
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Foydalanuvchi topilmadi"})
		return
	}
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
