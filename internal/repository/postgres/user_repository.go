package postgres

import (
	"database/sql"
	"github.com/danonenka/PR-service/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (id, name, is_active, team_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.IsActive, user.TeamID)
	return err
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, name, is_active, team_id FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamID)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return user, err
}

func (r *UserRepository) GetByTeamID(teamID string) ([]*domain.User, error) {
	query := `SELECT id, name, is_active, team_id FROM users WHERE team_id = $1`
	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) GetActiveByTeamID(teamID string) ([]*domain.User, error) {
	query := `SELECT id, name, is_active, team_id FROM users WHERE team_id = $1 AND is_active = true`
	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET name = $2, is_active = $3, team_id = $4 WHERE id = $1`
	_, err := r.db.Exec(query, user.ID, user.Name, user.IsActive, user.TeamID)
	return err
}

func (r *UserRepository) DeactivateUsers(teamID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	query := `UPDATE users SET is_active = false WHERE team_id = $1 AND id = ANY($2)`
	_, err := r.db.Exec(query, teamID, userIDs)
	return err
}

func (r *UserRepository) GetByIDs(ids []string) ([]*domain.User, error) {
	if len(ids) == 0 {
		return []*domain.User{}, nil
	}

	query := `SELECT id, name, is_active, team_id FROM users WHERE id = ANY($1)`
	rows, err := r.db.Query(query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}