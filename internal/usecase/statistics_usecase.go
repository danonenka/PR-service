package usecase

import (
	"github.com/danonenka/PR-service/internal/domain"
)

type StatisticsUsecase struct {
	prRepo         domain.PullRequestRepository
	assignmentRepo domain.ReviewerAssignmentRepository
	userRepo       domain.UserRepository
}

func NewStatisticsUsecase(
	prRepo domain.PullRequestRepository,
	assignmentRepo domain.ReviewerAssignmentRepository,
	userRepo domain.UserRepository,
) *StatisticsUsecase {
	return &StatisticsUsecase{
		prRepo:         prRepo,
		assignmentRepo: assignmentRepo,
		userRepo:       userRepo,
	}
}

type UserAssignmentStats struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	Assignments int    `json:"assignments"`
}

type PRAssignmentStats struct {
	PRID        string `json:"prId"`
	PRTitle     string `json:"prTitle"`
	Assignments int    `json:"assignments"`
}

func (u *StatisticsUsecase) GetUserAssignmentStats() ([]*UserAssignmentStats, error) {
	allPRs, err := u.prRepo.GetAll()
	if err != nil {
		return nil, err
	}

	userStatsMap := make(map[string]*UserAssignmentStats)

	for _, pr := range allPRs {
		assignments, err := u.assignmentRepo.GetByPRID(pr.ID)
		if err != nil {
			continue
		}

		for _, assignment := range assignments {
			if stats, exists := userStatsMap[assignment.ReviewerID]; exists {
				stats.Assignments++
			} else {
				user, err := u.userRepo.GetByID(assignment.ReviewerID)
				if err != nil {
					continue
				}
				userStatsMap[assignment.ReviewerID] = &UserAssignmentStats{
					UserID:      assignment.ReviewerID,
					UserName:    user.Name,
					Assignments: 1,
				}
			}
		}
	}

	stats := make([]*UserAssignmentStats, 0, len(userStatsMap))
	for _, stat := range userStatsMap {
		stats = append(stats, stat)
	}

	return stats, nil
}

func (u *StatisticsUsecase) GetPRAssignmentStats() ([]*PRAssignmentStats, error) {
	allPRs, err := u.prRepo.GetAll()
	if err != nil {
		return nil, err
	}

	stats := make([]*PRAssignmentStats, 0, len(allPRs))
	for _, pr := range allPRs {
		assignments, err := u.assignmentRepo.GetByPRID(pr.ID)
		if err != nil {
			continue
		}

		stats = append(stats, &PRAssignmentStats{
			PRID:        pr.ID,
			PRTitle:     pr.Title,
			Assignments: len(assignments),
		})
	}

	return stats, nil
}

