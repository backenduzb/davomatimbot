package crud

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *CRUDController[T, S, R]) Delete(c *gin.Context) {
	id := c.Param("id")
	var model T
	if err := ctrl.DB.First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}
	if err := ctrl.DB.Delete(&model, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// BulkDeleteRequest — bir nechta yozuvni bitta so'rovda o'chirish uchun.
type BulkDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1,max=5000"`
}

// BulkDeleteResponse — o'chirilgan yozuvlar soni.
type BulkDeleteResponse struct {
	Deleted int64 `json:"deleted"`
}

// BulkDelete berilgan ID'lar ro'yxatini BITTA SQL so'rovi bilan o'chiradi.
// Ilgari frontend har bir yozuv uchun alohida DELETE yuborardi — bu katta
// ro'yxatlarda juda sekin edi.
func (ctrl *CRUDController[T, S, R]) BulkDelete(c *gin.Context) {
	var req BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ids := uniqueIDs(req.IDs)
	if len(ids) == 0 {
		c.JSON(http.StatusOK, BulkDeleteResponse{Deleted: 0})
		return
	}

	var model T
	tx := ctrl.DB.Where("id IN ?", ids).Delete(&model)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, BulkDeleteResponse{Deleted: tx.RowsAffected})
}

// uniqueIDs takrorlanuvchi va nol ID'larni olib tashlaydi.
func uniqueIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
