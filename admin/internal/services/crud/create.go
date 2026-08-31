package crud

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ctrl *CRUDController[T, S, R]) Create(c *gin.Context) {
	var schema S
	if err := c.ShouldBindJSON(&schema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	model, err := MapSchemaToModel[T](schema)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.DB.Create(&model).Error; err != nil {
		if IsUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Bunday yozuv allaqachon mavjud"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp, err := MapModelToResponse[R](model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}