package usecase

import (
	"errors"
	"math/rand"
	"pr-service-task/internal/domain"
	"time"
)

type PRUsecase struct {
	prRepo              domain.PullRequestRepository
	userRepo            domain.UserRepository
	assignmentRepo      domain.ReviewerAssignmentRepository
	reviewerService     *ReviewerService
}

func NewPRUsecase(
	prRepo domain.PullRequestRepository,
	userRepo domain.UserRepository,
	assignmentRepo domain.ReviewerAssignmentRepository,
) *PRUsecase {
	return &PRUsecase{
		prRepo:          prRepo,
		userRepo:        userRepo,
		assignmentRepo:  assignmentRepo,
		reviewerService: NewReviewerService(userRepo, assignmentRepo),
	}
}

func (u *PRUsecase) CreatePR(pr *domain.PullRequest) error {
	author, err := u.userRepo.GetByID(pr.AuthorID)
	if err != nil {
		return errors.New("author not found")
	}

	teamUsers, err := u.userRepo.GetActiveByTeamID(author.TeamID)
	if err != nil {
		return err
	}

	candidates := make([]*domain.User, 0)
	for _, user := range teamUsers {
		if user.ID != pr.AuthorID {
			candidates = append(candidates, user)
		}
	}

	reviewerCount := 2
	if len(candidates) < reviewerCount {
		reviewerCount = len(candidates)
	}

	rand.Seed(time.Now().UnixNano())
	selectedReviewers := make([]*domain.User, 0, reviewerCount)
	availableCandidates := make([]*domain.User, len(candidates))
	copy(availableCandidates, candidates)

	for i := 0; i < reviewerCount && len(availableCandidates) > 0; i++ {
		idx := rand.Intn(len(availableCandidates))
		selectedReviewers = append(selectedReviewers, availableCandidates[idx])
		availableCandidates = append(availableCandidates[:idx], availableCandidates[idx+1:]...)
	}

	if err := u.prRepo.Create(pr); err != nil {
		return err
	}

	pr.ReviewerIDs = make([]string, 0, len(selectedReviewers))
	for _, reviewer := range selectedReviewers {
		assignment := &domain.ReviewerAssignment{
			PRID:       pr.ID,
			ReviewerID: reviewer.ID,
		}
		if err := u.assignmentRepo.Create(assignment); err != nil {
			return err
		}
		pr.ReviewerIDs = append(pr.ReviewerIDs, reviewer.ID)
	}

	return nil
}

func (u *PRUsecase) GetPRByID(id string) (*domain.PullRequest, error) {
	pr, err := u.prRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	assignments, err := u.assignmentRepo.GetByPRID(id)
	if err != nil {
		return nil, err
	}

	pr.ReviewerIDs = make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		pr.ReviewerIDs = append(pr.ReviewerIDs, assignment.ReviewerID)
	}

	return pr, nil
}

func (u *PRUsecase) GetPRsByAuthorID(authorID string) ([]*domain.PullRequest, error) {
	prs, err := u.prRepo.GetByAuthorID(authorID)
	if err != nil {
		return nil, err
	}

	for _, pr := range prs {
		assignments, err := u.assignmentRepo.GetByPRID(pr.ID)
		if err != nil {
			return nil, err
		}
		pr.ReviewerIDs = make([]string, 0, len(assignments))
		for _, assignment := range assignments {
			pr.ReviewerIDs = append(pr.ReviewerIDs, assignment.ReviewerID)
		}
	}

	return prs, nil
}

func (u *PRUsecase) GetPRsByReviewerID(reviewerID string) ([]*domain.PullRequest, error) {
	assignments, err := u.assignmentRepo.GetByReviewerID(reviewerID)
	if err != nil {
		return nil, err
	}

	prs := make([]*domain.PullRequest, 0, len(assignments))
	for _, assignment := range assignments {
		pr, err := u.prRepo.GetByID(assignment.PRID)
		if err != nil {
			continue
		}

		allAssignments, err := u.assignmentRepo.GetByPRID(pr.ID)
		if err != nil {
			continue
		}
		pr.ReviewerIDs = make([]string, 0, len(allAssignments))
		for _, a := range allAssignments {
			pr.ReviewerIDs = append(pr.ReviewerIDs, a.ReviewerID)
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

func (u *PRUsecase) ReassignReviewer(prID string, oldReviewerID string) error {
	pr, err := u.prRepo.GetByID(prID)
	if err != nil {
		return errors.New("PR not found")
	}

	if pr.Status == domain.PRStatusMerged {
		return errors.New("cannot reassign reviewers for merged PR")
	}

	oldReviewer, err := u.userRepo.GetByID(oldReviewerID)
	if err != nil {
		return errors.New("old reviewer not found")
	}

	teamUsers, err := u.userRepo.GetActiveByTeamID(oldReviewer.TeamID)
	if err != nil {
		return err
	}

	candidates := make([]*domain.User, 0)
	excludedIDs := make(map[string]bool)
	excludedIDs[pr.AuthorID] = true
	excludedIDs[oldReviewerID] = true

	assignments, err := u.assignmentRepo.GetByPRID(prID)
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
		return errors.New("no available reviewers")
	}

	rand.Seed(time.Now().UnixNano())
	newReviewer := candidates[rand.Intn(len(candidates))]

	if err := u.assignmentRepo.Delete(prID, oldReviewerID); err != nil {
		return err
	}

	newAssignment := &domain.ReviewerAssignment{
		PRID:       prID,
		ReviewerID: newReviewer.ID,
	}
	return u.assignmentRepo.Create(newAssignment)
}

func (u *PRUsecase) MergePR(prID string) error {
	pr, err := u.prRepo.GetByID(prID)
	if err != nil {
		return errors.New("PR not found")
	}

	// Если уже merged, просто возвращаем успех (идемпотентность)
	if pr.Status == domain.PRStatusMerged {
		return nil
	}

	pr.Status = domain.PRStatusMerged
	return u.prRepo.Update(pr)
}

type ReviewerService struct {
	userRepo       domain.UserRepository
	assignmentRepo domain.ReviewerAssignmentRepository
}

func NewReviewerService(userRepo domain.UserRepository, assignmentRepo domain.ReviewerAssignmentRepository) *ReviewerService {
	return &ReviewerService{
		userRepo:       userRepo,
		assignmentRepo: assignmentRepo,
	}
}

