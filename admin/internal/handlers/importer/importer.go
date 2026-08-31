package importer

import (
	"fmt"
	"net/http"

	"admin/internal/database"
	"admin/internal/models"
	importerService "admin/internal/services/importer"

	"github.com/gin-gonic/gin"
)

// maxUploadSize — ruxsat etilgan maksimal fayl hajmi (10 MB).
const maxUploadSize = 10 << 20

// UploadXLSX xlsx faylni qabul qilib, o'quvchilarni/sinflarni bazaga import qiladi.
// Faqat admin buni bajarishi mumkin (login allaqachon admin bo'lishni talab qiladi,
// bu yerda qo'shimcha himoya sifatida yana tekshiramiz).
func UploadXLSX(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Avtorizatsiya talab qilinadi"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userIDVal).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Foydalanuvchi topilmadi"})
		return
	}
	if !user.IsAdmin || user.IsBanned {
		c.JSON(http.StatusForbidden, gin.H{"error": "Faqat admin fayl yuklashi mumkin"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fayl yuborilmadi."})
		return
	}
	if fileHeader.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fayl hajmi juda katta (maksimum 10 MB)."})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faylni ochib bo'lmadi."})
		return
	}
	defer file.Close()

	result, err := importerService.Import(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imported := result.StudentsCreated + result.StudentsLinked
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d ta o'quvchi muvaffaqiyatli import qilindi ✅", imported),
		"result":  result,
	})
}
