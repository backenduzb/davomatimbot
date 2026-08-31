package statistics

import (
	"admin/internal/database"
	statisticsService "admin/internal/services/statistics"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetToday(c *gin.Context) {
	db := database.DB
	dateStr := c.Query("date")

	result, err := statisticsService.GetTodayStatistics(db, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
