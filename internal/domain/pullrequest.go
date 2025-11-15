package domain

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	AuthorID    string    `json:"authorId"`
	Status      PRStatus  `json:"status"`
	ReviewerIDs []string  `json:"reviewerIds"`
}

type PullRequestRepository interface {
	Create(pr *PullRequest) error
	GetByID(id string) (*PullRequest, error)
	GetByAuthorID(authorID string) ([]*PullRequest, error)
	GetByReviewerID(reviewerID string) ([]*PullRequest, error)
	Update(pr *PullRequest) error
	GetAll() ([]*PullRequest, error)
}

