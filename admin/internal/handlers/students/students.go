package students

import (
	"net/http"

	"admin/internal/database"
	"admin/internal/models"
	"admin/internal/schemas/crud"
	"admin/internal/schemas/student"

	"github.com/gin-gonic/gin"
)

// List o'quvchilar ro'yxatini qaytaradi.
// Ixtiyoriy filtrlar: class_id (sinf bo'yicha) va search (F.I.Sh bo'yicha).
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

	query := database.DB.Model(&models.Student{})

	if classID := c.Query("class_id"); classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("full_name ILIKE ?", "%"+search+"%")
	}

	var count int64
	query.Count(&count)

	var students []models.Student
	offset := (params.Page - 1) * params.PageSize
	query.Order("full_name asc").Limit(params.PageSize).Offset(offset).Find(&students)

	results := make([]student.Student, 0, len(students))
	for _, s := range students {
		results = append(results, student.Student{
			ID:       s.ID,
			FullName: s.FullName,
			ClassID:  s.ClassID,
		})
	}

	totalPages := int64(0)
	if count > 0 {
		totalPages = (count + int64(params.PageSize) - 1) / int64(params.PageSize)
	}

	c.JSON(http.StatusOK, crud.PaginatedResponse[student.Student]{
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: int64(params.Page),
		PageSize:    int64(params.PageSize),
		Results:     results,
	})
}
