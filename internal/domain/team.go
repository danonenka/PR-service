package domain

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TeamRepository interface {
	Create(team *Team) error
	GetByID(id string) (*Team, error)
	GetByName(name string) (*Team, error)
	GetAll() ([]*Team, error)
	Update(team *Team) error
	Delete(id string) error
}

