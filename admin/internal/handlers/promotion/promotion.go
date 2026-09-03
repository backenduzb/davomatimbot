package promotion

import (
	"net/http"

	"admin/internal/database"
	promotionService "admin/internal/services/promotion"

	"github.com/gin-gonic/gin"
)

// Preview sinflarni oshirishdan OLDIN nima o'zgarishini qaytaradi.
// Frontend shu ma'lumot asosida tasdiqlash oynasini ko'rsatadi.
func Preview(c *gin.Context) {
	plans, err := promotionService.BuildPlan(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	promote, toDelete, skip := 0, 0, 0
	var studentsToDelete int64
	for _, plan := range plans {
		switch plan.Action {
		case promotionService.ActionPromote:
			promote++
		case promotionService.ActionDelete:
			toDelete++
			studentsToDelete += plan.StudentCount
		case promotionService.ActionSkip:
			skip++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"plans":             plans,
		"promote":           promote,
		"delete":            toDelete,
		"students_to_delete": studentsToDelete,
		"skip":              skip,
		"total":             len(plans),
	})
}

// PromoteRequest — tasdiqlash uchun so'rov tanasi.
// Frontend tasodifiy bosishdan himoya sifatida confirm=true yuborishi shart.
type PromoteRequest struct {
	Confirm bool `json:"confirm"`
}

// Promote sinflarni keyingi o'quv yiliga oshiradi.
func Promote(c *gin.Context) {
	var req PromoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Amalni bajarish uchun tasdiqlash talab qilinadi (confirm: true)",
		})
		return
	}

	result, err := promotionService.Promote(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
