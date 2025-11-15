package postgres

import (
	"database/sql"
	"pr-service-task/internal/domain"
)

type PullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) Create(pr *domain.PullRequest) error {
	query := `INSERT INTO pull_requests (id, title, author_id, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, pr.ID, pr.Title, pr.AuthorID, pr.Status)
	return err
}

func (r *PullRequestRepository) GetByID(id string) (*domain.PullRequest, error) {
	query := `SELECT id, title, author_id, status FROM pull_requests WHERE id = $1`
	pr := &domain.PullRequest{}
	err := r.db.QueryRow(query, id).Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return pr, err
}

func (r *PullRequestRepository) GetByAuthorID(authorID string) ([]*domain.PullRequest, error) {
	query := `SELECT id, title, author_id, status FROM pull_requests WHERE author_id = $1`
	rows, err := r.db.Query(query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := make([]*domain.PullRequest, 0)
	for rows.Next() {
		pr := &domain.PullRequest{}
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, rows.Err()
}

func (r *PullRequestRepository) GetByReviewerID(reviewerID string) ([]*domain.PullRequest, error) {
	query := `
		SELECT pr.id, pr.title, pr.author_id, pr.status 
		FROM pull_requests pr
		INNER JOIN reviewer_assignments ra ON pr.id = ra.pr_id
		WHERE ra.reviewer_id = $1
	`
	rows, err := r.db.Query(query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := make([]*domain.PullRequest, 0)
	for rows.Next() {
		pr := &domain.PullRequest{}
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, rows.Err()
}

func (r *PullRequestRepository) Update(pr *domain.PullRequest) error {
	query := `UPDATE pull_requests SET title = $2, author_id = $3, status = $4 WHERE id = $1`
	_, err := r.db.Exec(query, pr.ID, pr.Title, pr.AuthorID, pr.Status)
	return err
}

func (r *PullRequestRepository) GetAll() ([]*domain.PullRequest, error) {
	query := `SELECT id, title, author_id, status FROM pull_requests`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := make([]*domain.PullRequest, 0)
	for rows.Next() {
		pr := &domain.PullRequest{}
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, rows.Err()
}

