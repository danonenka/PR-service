package postgres

import (
	"database/sql"
	"pr-service-task/internal/domain"
)

type ReviewerAssignmentRepository struct {
	db *sql.DB
}

func NewReviewerAssignmentRepository(db *sql.DB) *ReviewerAssignmentRepository {
	return &ReviewerAssignmentRepository{db: db}
}

func (r *ReviewerAssignmentRepository) Create(assignment *domain.ReviewerAssignment) error {
	query := `INSERT INTO reviewer_assignments (pr_id, reviewer_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, assignment.PRID, assignment.ReviewerID)
	return err
}

func (r *ReviewerAssignmentRepository) Delete(prID string, reviewerID string) error {
	query := `DELETE FROM reviewer_assignments WHERE pr_id = $1 AND reviewer_id = $2`
	_, err := r.db.Exec(query, prID, reviewerID)
	return err
}

func (r *ReviewerAssignmentRepository) GetByPRID(prID string) ([]*domain.ReviewerAssignment, error) {
	query := `SELECT pr_id, reviewer_id FROM reviewer_assignments WHERE pr_id = $1`
	rows, err := r.db.Query(query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := make([]*domain.ReviewerAssignment, 0)
	for rows.Next() {
		assignment := &domain.ReviewerAssignment{}
		if err := rows.Scan(&assignment.PRID, &assignment.ReviewerID); err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}
	return assignments, rows.Err()
}

func (r *ReviewerAssignmentRepository) GetByReviewerID(reviewerID string) ([]*domain.ReviewerAssignment, error) {
	query := `SELECT pr_id, reviewer_id FROM reviewer_assignments WHERE reviewer_id = $1`
	rows, err := r.db.Query(query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := make([]*domain.ReviewerAssignment, 0)
	for rows.Next() {
		assignment := &domain.ReviewerAssignment{}
		if err := rows.Scan(&assignment.PRID, &assignment.ReviewerID); err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}
	return assignments, rows.Err()
}

func (r *ReviewerAssignmentRepository) DeleteByPRID(prID string) error {
	query := `DELETE FROM reviewer_assignments WHERE pr_id = $1`
	_, err := r.db.Exec(query, prID)
	return err
}

