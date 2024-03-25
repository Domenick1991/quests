package services

import (
	"errors"
	"quests/internal"
	"quests/repository"
)

type HistoryService struct {
	repo repository.History
}

func NewHistoryService(repo repository.History) *HistoryService {
	return &HistoryService{repo: repo}
}

func (s *HistoryService) CompleteSteps(сompleteSteps internal.NewCompleteSteps) error {
	for _, сompleteStep := range сompleteSteps.CompleteSteps {
		err := s.repo.CompleteSteps(сompleteStep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *HistoryService) GetHistory(userid int) (internal.UserBonus, error) {
	quests := s.repo.GetCompletedQuest(userid)
	userBonus := internal.UserBonus{}
	if len(quests) > 0 {
		for _, quest := range quests {
			CompletedQuest := s.repo.GetCompletedQuestForUser(userid, quest)
			userBonus.CompletedQuests = append(userBonus.CompletedQuests, CompletedQuest)
			userBonus.TotalBonus += CompletedQuest.Bonus
		}
		return userBonus, nil
	}
	return userBonus, errors.New("Пользователь еще не выполнял задания")
}
