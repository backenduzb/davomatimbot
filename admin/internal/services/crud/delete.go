package crud

import (
	"github.com/gin-gonic/gin"
	"net/http"
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