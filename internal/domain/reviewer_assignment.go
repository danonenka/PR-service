package domain

type ReviewerAssignment struct {
	PRID       string `json:"prId"`
	ReviewerID string `json:"reviewerId"`
}

type ReviewerAssignmentRepository interface {
	Create(assignment *ReviewerAssignment) error
	Delete(prID string, reviewerID string) error
	GetByPRID(prID string) ([]*ReviewerAssignment, error)
	GetByReviewerID(reviewerID string) ([]*ReviewerAssignment, error)
	DeleteByPRID(prID string) error
}

