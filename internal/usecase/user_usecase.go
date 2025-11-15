package usecase

import (
	"errors"
	"github.com/danonenka/PR-service/internal/domain"
)

type UserUsecase struct {
	userRepo            domain.UserRepository
	teamRepo            domain.TeamRepository
	reassignmentUsecase *ReassignmentUsecase
}

func NewUserUsecase(userRepo domain.UserRepository, teamRepo domain.TeamRepository, reassignmentUsecase *ReassignmentUsecase) *UserUsecase {
	return &UserUsecase{
		userRepo:            userRepo,
		teamRepo:            teamRepo,
		reassignmentUsecase: reassignmentUsecase,
	}
}

func (u *UserUsecase) CreateUser(user *domain.User) error {
	_, err := u.teamRepo.GetByID(user.TeamID)
	if err != nil {
		return errors.New("team not found")
	}
	return u.userRepo.Create(user)
}

func (u *UserUsecase) GetUserByID(id string) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *UserUsecase) GetUsersByTeamID(teamID string) ([]*domain.User, error) {
	return u.userRepo.GetByTeamID(teamID)
}

func (u *UserUsecase) GetActiveUsersByTeamID(teamID string) ([]*domain.User, error) {
	return u.userRepo.GetActiveByTeamID(teamID)
}

func (u *UserUsecase) UpdateUser(user *domain.User) error {
	return u.userRepo.Update(user)
}

func (u *UserUsecase) DeactivateUsers(teamID string, userIDs []string) error {
	if u.reassignmentUsecase != nil {
		if err := u.reassignmentUsecase.ReassignDeactivatedReviewers(teamID, userIDs); err != nil {
		}
	}
	return u.userRepo.DeactivateUsers(teamID, userIDs)
}

func (u *UserUsecase) SetUserIsActive(userID string, isActive bool) (*domain.User, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.IsActive = isActive
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
