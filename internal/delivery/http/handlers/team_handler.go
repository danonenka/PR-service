package handlers

import (
	"net/http"
	"github.com/danonenka/PR-service/internal/domain"
	"github.com/danonenka/PR-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamUsecase *usecase.TeamUsecase
}

func NewTeamHandler(teamUsecase *usecase.TeamUsecase) *TeamHandler {
	return &TeamHandler{teamUsecase: teamUsecase}
}

type TeamMemberRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type TeamRequest struct {
	TeamName string              `json:"team_name" binding:"required"`
	Members  []TeamMemberRequest `json:"members" binding:"required"`
}

type TeamMemberResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamResponse struct {
	TeamName string               `json:"team_name"`
	Members  []TeamMemberResponse `json:"members"`
}

func (h *TeamHandler) AddTeam(c *gin.Context) {
	var req TeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	members := make([]*domain.User, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, &domain.User{
			ID:       m.UserID,
			Name:     m.Username,
			IsActive: m.IsActive,
		})
	}

	team, err := h.teamUsecase.AddTeamWithMembers(req.TeamName, members)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	_, updatedMembers, err := h.teamUsecase.GetTeamWithMembers(team.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	memberResponses := make([]TeamMemberResponse, 0, len(updatedMembers))
	for _, m := range updatedMembers {
		memberResponses = append(memberResponses, TeamMemberResponse{
			UserID:   m.ID,
			Username: m.Name,
			IsActive: m.IsActive,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"team": TeamResponse{
			TeamName: team.Name,
			Members:  memberResponses,
		},
	})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "team_name parameter is required",
			},
		})
		return
	}

	team, members, err := h.teamUsecase.GetTeamWithMembers(teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	memberResponses := make([]TeamMemberResponse, 0, len(members))
	for _, m := range members {
		memberResponses = append(memberResponses, TeamMemberResponse{
			UserID:   m.ID,
			Username: m.Name,
			IsActive: m.IsActive,
		})
	}

	c.JSON(http.StatusOK, TeamResponse{
		TeamName: team.Name,
		Members:  memberResponses,
	})
}
