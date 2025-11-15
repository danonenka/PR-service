package handlers

import (
	"net/http"
	"pr-service-task/internal/domain"
	"pr-service-task/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

type PRHandler struct {
	prUsecase *usecase.PRUsecase
}

func NewPRHandler(prUsecase *usecase.PRUsecase) *PRHandler {
	return &PRHandler{prUsecase: prUsecase}
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type PRResponse struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

func (h *PRHandler) CreatePR(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	_, err := h.prUsecase.GetPRByID(req.PullRequestID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "PR_EXISTS",
				"message": "PR id already exists",
			},
		})
		return
	}

	pr := &domain.PullRequest{
		ID:          req.PullRequestID,
		Title:       req.PullRequestName,
		AuthorID:    req.AuthorID,
		Status:      domain.PRStatusOpen,
		ReviewerIDs: []string{},
	}

	if err := h.prUsecase.CreatePR(pr); err != nil {
		if err.Error() == "author not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "resource not found",
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

	pr, err = h.prUsecase.GetPRByID(pr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	now := time.Now().Format(time.RFC3339)
	c.JSON(http.StatusCreated, gin.H{
		"pr": PRResponse{
			PullRequestID:     pr.ID,
			PullRequestName:   pr.Title,
			AuthorID:          pr.AuthorID,
			Status:            string(pr.Status),
			AssignedReviewers: pr.ReviewerIDs,
			CreatedAt:         &now,
		},
	})
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

func (h *PRHandler) MergePR(c *gin.Context) {
	var req MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.prUsecase.MergePR(req.PullRequestID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	pr, err := h.prUsecase.GetPRByID(req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	mergedAt := time.Now().Format(time.RFC3339)
	c.JSON(http.StatusOK, gin.H{
		"pr": PRResponse{
			PullRequestID:     pr.ID,
			PullRequestName:   pr.Title,
			AuthorID:          pr.AuthorID,
			Status:            string(pr.Status),
			AssignedReviewers: pr.ReviewerIDs,
			MergedAt:          &mergedAt,
		},
	})
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_user_id" binding:"required"`
}

func (h *PRHandler) ReassignReviewer(c *gin.Context) {
	var req ReassignReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	pr, err := h.prUsecase.GetPRByID(req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	if pr.Status == domain.PRStatusMerged {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "PR_MERGED",
				"message": "cannot reassign on merged PR",
			},
		})
		return
	}

	found := false
	for _, reviewerID := range pr.ReviewerIDs {
		if reviewerID == req.OldUserID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "NOT_ASSIGNED",
				"message": "reviewer is not assigned to this PR",
			},
		})
		return
	}

	err = h.prUsecase.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		if err.Error() == "no available reviewers" {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "NO_CANDIDATE",
					"message": "no active replacement candidate in team",
				},
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	updatedPR, err := h.prUsecase.GetPRByID(req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var newReviewerID string
	for _, reviewerID := range updatedPR.ReviewerIDs {
		if reviewerID != req.OldUserID {
			newReviewerID = reviewerID
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": PRResponse{
			PullRequestID:     updatedPR.ID,
			PullRequestName:   updatedPR.Title,
			AuthorID:          updatedPR.AuthorID,
			Status:            string(updatedPR.Status),
			AssignedReviewers: updatedPR.ReviewerIDs,
		},
		"replaced_by": newReviewerID,
	})
}
