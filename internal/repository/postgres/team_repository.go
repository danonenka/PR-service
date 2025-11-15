package postgres

import (
	"database/sql"
	"github.com/danonenka/PR-service/internal/domain"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *domain.Team) error {
	query := `INSERT INTO teams (id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(query, team.ID, team.Name)
	return err
}

func (r *TeamRepository) GetByID(id string) (*domain.Team, error) {
	query := `SELECT id, name FROM teams WHERE id = $1`
	team := &domain.Team{}
	err := r.db.QueryRow(query, id).Scan(&team.ID, &team.Name)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return team, err
}

func (r *TeamRepository) GetByName(name string) (*domain.Team, error) {
	query := `SELECT id, name FROM teams WHERE name = $1`
	team := &domain.Team{}
	err := r.db.QueryRow(query, name).Scan(&team.ID, &team.Name)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return team, err
}

func (r *TeamRepository) GetAll() ([]*domain.Team, error) {
	query := `SELECT id, name FROM teams`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := make([]*domain.Team, 0)
	for rows.Next() {
		team := &domain.Team{}
		if err := rows.Scan(&team.ID, &team.Name); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, rows.Err()
}

func (r *TeamRepository) Update(team *domain.Team) error {
	query := `UPDATE teams SET name = $2 WHERE id = $1`
	_, err := r.db.Exec(query, team.ID, team.Name)
	return err
}

func (r *TeamRepository) Delete(id string) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

