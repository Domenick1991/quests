package services

import (
	"quests/internal"
	"quests/repository"
)

type HistoryService struct {
	repo repository.History
}

func NewHistoryService(repo repository.History) *HistoryService {
	return &HistoryService{repo: repo}
}

func (s *HistoryService) CompleteSteps(сompleteSteps internal.NewCompleteSteps) []internal.ErrorList {
	return s.repo.CompleteSteps(сompleteSteps)
}

func (s *HistoryService) GetHistory(userid int) (internal.UserBonus, error) {
	return s.repo.GetHistory(userid)
}
