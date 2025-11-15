package handlers

import (
	"net/http"
	"github.com/danonenka/PR-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	prUsecase   *usecase.PRUsecase
	teamUsecase *usecase.TeamUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase, prUsecase *usecase.PRUsecase, teamUsecase *usecase.TeamUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		prUsecase:   prUsecase,
		teamUsecase: teamUsecase,
	}
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	var req SetIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	user, err := h.userUsecase.SetUserIsActive(req.UserID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	team, err := h.teamUsecase.GetTeamByID(user.TeamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": UserResponse{
			UserID:   user.ID,
			Username: user.Name,
			TeamName: team.Name,
			IsActive: user.IsActive,
		},
	})
}

type PullRequestShortResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type GetReviewResponse struct {
	UserID       string                      `json:"user_id"`
	PullRequests []PullRequestShortResponse `json:"pull_requests"`
}

func (h *UserHandler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "user_id parameter is required",
			},
		})
		return
	}

	prs, err := h.prUsecase.GetPRsByReviewerID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	prResponses := make([]PullRequestShortResponse, 0, len(prs))
	for _, pr := range prs {
		prResponses = append(prResponses, PullRequestShortResponse{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Title,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}

	c.JSON(http.StatusOK, GetReviewResponse{
		UserID:       userID,
		PullRequests: prResponses,
	})
}
