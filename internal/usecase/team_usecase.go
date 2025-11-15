package usecase

import (
	"errors"
	"pr-service-task/internal/domain"
)

type TeamUsecase struct {
	teamRepo domain.TeamRepository
	userRepo domain.UserRepository
}

func NewTeamUsecase(teamRepo domain.TeamRepository, userRepo domain.UserRepository) *TeamUsecase {
	return &TeamUsecase{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (u *TeamUsecase) CreateTeam(team *domain.Team) error {
	return u.teamRepo.Create(team)
}

func (u *TeamUsecase) GetTeamByID(id string) (*domain.Team, error) {
	return u.teamRepo.GetByID(id)
}

func (u *TeamUsecase) GetTeamByName(name string) (*domain.Team, error) {
	return u.teamRepo.GetByName(name)
}

func (u *TeamUsecase) GetAllTeams() ([]*domain.Team, error) {
	return u.teamRepo.GetAll()
}

func (u *TeamUsecase) UpdateTeam(team *domain.Team) error {
	return u.teamRepo.Update(team)
}

func (u *TeamUsecase) DeleteTeam(id string) error {
	return u.teamRepo.Delete(id)
}

func (u *TeamUsecase) AddTeamWithMembers(teamName string, members []*domain.User) (*domain.Team, error) {
	_, err := u.teamRepo.GetByName(teamName)
	if err == nil {
		return nil, errors.New("TEAM_EXISTS")
	}

	team := &domain.Team{
		ID:   teamName,
		Name: teamName,
	}

	if err := u.teamRepo.Create(team); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"teams_name_key\"" {
			return nil, errors.New("TEAM_EXISTS")
		}
		return nil, err
	}

	for _, member := range members {
		member.TeamID = team.ID
		existingUser, err := u.userRepo.GetByID(member.ID)
		if err == nil && existingUser != nil {
			existingUser.Name = member.Name
			existingUser.IsActive = member.IsActive
			existingUser.TeamID = team.ID
			if err := u.userRepo.Update(existingUser); err != nil {
				return nil, err
			}
		} else {
			if err := u.userRepo.Create(member); err != nil {
				return nil, err
			}
		}
	}

	return team, nil
}

func (u *TeamUsecase) GetTeamWithMembers(teamName string) (*domain.Team, []*domain.User, error) {
	team, err := u.teamRepo.GetByName(teamName)
	if err != nil {
		return nil, nil, errors.New("team not found")
	}

	members, err := u.userRepo.GetByTeamID(team.ID)
	if err != nil {
		return nil, nil, err
	}

	return team, members, nil
}

