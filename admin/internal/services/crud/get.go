package crud

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"admin/internal/schemas/crud"
)

func (ctrl *CRUDController[T, S, R]) Retrieve(c *gin.Context) {
	id := c.Param("id")
	var model T
	if err := ctrl.DB.First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}
	resp, err := MapModelToResponse[R](model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *CRUDController[T, S, R]) List(c *gin.Context) {
	var _ []T
	db := ctrl.DB.Model(new(T))
	
	paginated, err := Paginate[T](c, db) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	respResults, err := MapModelsToResponse[R](paginated.Results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := &crud.PaginatedResponse[R]{
		Count:       paginated.Count,
		TotalPages:  paginated.TotalPages,
		CurrentPage: paginated.CurrentPage,
		PageSize:    paginated.PageSize,
		Results:     respResults,
	}
	c.JSON(http.StatusOK, resp)
}