package students

import (
	"net/http"
	"strings"

	"admin/internal/database"
	"admin/internal/models"
	"admin/internal/schemas/crud"
	"admin/internal/schemas/student"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// List o'quvchilar ro'yxatini qaytaradi.
// Ixtiyoriy filtrlar: class_id (sinf bo'yicha) va search (F.I.Sh bo'yicha).
//
// Tartib: o'quvchilar avval SINF bo'yicha (sinf nomi alifbo tartibida),
// keyin sinf ichida F.I.Sh bo'yicha saralanadi. Sinf nomi bilan saralash
// uchun classes va class_names jadvallariga JOIN qilinadi — bu N+1
// so'rovlarsiz bitta so'rovda bajariladi.
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

	classID := c.Query("class_id")
	search := strings.TrimSpace(c.Query("search"))

	// Har bir so'rov uchun alohida query quramiz — bitta *gorm.DB obyektini
	// Count va Find uchun qayta ishlatish shartlarning takrorlanishiga olib
	// kelishi mumkin.
	baseQuery := func() *gorm.DB {
		q := database.DB.Model(&models.Student{})
		if classID != "" {
			q = q.Where("students.class_id = ?", classID)
		}
		if search != "" {
			q = q.Where("students.full_name ILIKE ?", "%"+search+"%")
		}
		return q
	}

	var count int64
	if err := baseQuery().Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var students []models.Student
	offset := (params.Page - 1) * params.PageSize
	if err := baseQuery().
		Joins("LEFT JOIN classes ON classes.id = students.class_id").
		Joins("LEFT JOIN class_names ON class_names.id = classes.class_name_id").
		Order("class_names.name ASC NULLS LAST").
		Order("students.class_id ASC").
		Order("students.full_name ASC").
		Select("students.*").
		Limit(params.PageSize).
		Offset(offset).
		Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

// BulkDeleteRequest — bir nechta o'quvchini bitta so'rovda o'chirish uchun.
type BulkDeleteRequest struct {
	IDs     []uint `json:"ids"`
	ClassID *uint  `json:"class_id"`
	All     bool   `json:"all"`
}

// BulkDelete tanlangan o'quvchilarni (yoki butun sinf o'quvchilarini)
// bitta SQL so'rovi bilan o'chiradi. Davomat yozuvlari ham o'chiriladi,
// chunki ular o'quvchiga bog'liq.
func BulkDelete(c *gin.Context) {
	var req BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 && req.ClassID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids yoki class_id ko'rsatilishi kerak"})
		return
	}

	var deleted int64
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		scope := tx.Model(&models.Student{})
		attScope := tx.Model(&models.Attendance{})

		if len(req.IDs) > 0 {
			ids := uniqueIDs(req.IDs)
			if len(ids) == 0 {
				return nil
			}
			scope = scope.Where("id IN ?", ids)
			attScope = attScope.Where("student_id IN ?", ids)
		} else {
			scope = scope.Where("class_id = ?", *req.ClassID)
			attScope = attScope.Where("class_id = ?", *req.ClassID)
		}

		if err := attScope.Delete(&models.Attendance{}).Error; err != nil {
			return err
		}
		res := scope.Delete(&models.Student{})
		if res.Error != nil {
			return res.Error
		}
		deleted = res.RowsAffected
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": deleted})
}

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
