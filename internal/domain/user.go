package domain

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
	TeamID   string `json:"teamId"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByTeamID(teamID string) ([]*User, error)
	GetActiveByTeamID(teamID string) ([]*User, error)
	Update(user *User) error
	DeactivateUsers(teamID string, userIDs []string) error
	GetByIDs(ids []string) ([]*User, error)
}

