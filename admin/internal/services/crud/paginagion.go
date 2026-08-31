package crud

import (
	"admin/internal/schemas/crud"
	"math"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginate[T any](c *gin.Context, db *gorm.DB) (*crud.PaginatedResponse[T], error) {
	var params crud.PaginationParams

	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, err
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	var count int64
	var results []T

	var model T
	if err := db.Model(&model).Count(&count).Error; err != nil {
		return nil, err
	}
	
	totalPages := int64(0)
	if count > 0 {
		totalPages = int64(math.Ceil(float64(count) / float64(params.PageSize)))
	}

	offset := (params.Page - 1) * params.PageSize
	if err := db.Limit(params.PageSize).Offset(offset).Find(&results).Error; err != nil {
		return nil,  err
	}

	return &crud.PaginatedResponse[T]{
		Count: count,
		TotalPages: totalPages,
		CurrentPage: int64(params.Page),
		PageSize: int64(params.PageSize),
		Results: results,
	}, nil	
}