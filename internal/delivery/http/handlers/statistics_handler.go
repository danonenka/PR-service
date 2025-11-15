package handlers

import (
	"net/http"
	"pr-service-task/internal/usecase"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statisticsUsecase *usecase.StatisticsUsecase
}

func NewStatisticsHandler(statisticsUsecase *usecase.StatisticsUsecase) *StatisticsHandler {
	return &StatisticsHandler{statisticsUsecase: statisticsUsecase}
}

func (h *StatisticsHandler) GetUserStats(c *gin.Context) {
	stats, err := h.statisticsUsecase.GetUserAssignmentStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *StatisticsHandler) GetPRStats(c *gin.Context) {
	stats, err := h.statisticsUsecase.GetPRAssignmentStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

