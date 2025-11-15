package usecase

import (
	"math/rand"
	"pr-service-task/internal/domain"
	"time"
)

type ReassignmentUsecase struct {
	prRepo         domain.PullRequestRepository
	userRepo       domain.UserRepository
	assignmentRepo domain.ReviewerAssignmentRepository
}

func NewReassignmentUsecase(
	prRepo domain.PullRequestRepository,
	userRepo domain.UserRepository,
	assignmentRepo domain.ReviewerAssignmentRepository,
) *ReassignmentUsecase {
	return &ReassignmentUsecase{
		prRepo:         prRepo,
		userRepo:       userRepo,
		assignmentRepo: assignmentRepo,
	}
}

func (u *ReassignmentUsecase) ReassignDeactivatedReviewers(teamID string, deactivatedUserIDs []string) error {
	if len(deactivatedUserIDs) == 0 {
		return nil
	}

	allPRs, err := u.prRepo.GetAll()
	if err != nil {
		return err
	}

	openPRs := make([]*domain.PullRequest, 0)
	for _, pr := range allPRs {
		if pr.Status == domain.PRStatusOpen {
			openPRs = append(openPRs, pr)
		}
	}

	for _, pr := range openPRs {
		assignments, err := u.assignmentRepo.GetByPRID(pr.ID)
		if err != nil {
			continue
		}

		deactivatedReviewers := make([]string, 0)
		for _, assignment := range assignments {
			for _, deactivatedID := range deactivatedUserIDs {
				if assignment.ReviewerID == deactivatedID {
					deactivatedReviewers = append(deactivatedReviewers, assignment.ReviewerID)
					break
				}
			}
		}

		for _, deactivatedReviewerID := range deactivatedReviewers {
			if err := u.reassignReviewerForPR(pr, deactivatedReviewerID); err != nil {
				continue
			}
		}
	}

	return nil
}

func (u *ReassignmentUsecase) reassignReviewerForPR(pr *domain.PullRequest, oldReviewerID string) error {
	author, err := u.userRepo.GetByID(pr.AuthorID)
	if err != nil {
		return err
	}

	teamUsers, err := u.userRepo.GetActiveByTeamID(author.TeamID)
	if err != nil {
		return err
	}

	candidates := make([]*domain.User, 0)
	excludedIDs := make(map[string]bool)
	excludedIDs[pr.AuthorID] = true
	excludedIDs[oldReviewerID] = true

	assignments, err := u.assignmentRepo.GetByPRID(pr.ID)
	if err != nil {
		return err
	}
	for _, assignment := range assignments {
		excludedIDs[assignment.ReviewerID] = true
	}

	for _, user := range teamUsers {
		if !excludedIDs[user.ID] {
			candidates = append(candidates, user)
		}
	}

	if len(candidates) == 0 {
		return u.assignmentRepo.Delete(pr.ID, oldReviewerID)
	}

	rand.Seed(time.Now().UnixNano())
	newReviewer := candidates[rand.Intn(len(candidates))]

	if err := u.assignmentRepo.Delete(pr.ID, oldReviewerID); err != nil {
		return err
	}

	newAssignment := &domain.ReviewerAssignment{
		PRID:       pr.ID,
		ReviewerID: newReviewer.ID,
	}
	return u.assignmentRepo.Create(newAssignment)
}
