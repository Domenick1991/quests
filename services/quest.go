package services

import (
	"quests/internal"
	"quests/repository"
)

type QuestService struct {
	repo repository.Quest
}

func NewQuestService(repo repository.Quest) *QuestService {
	return &QuestService{repo: repo}
}

func (s *QuestService) GetQuests() ([]internal.Quests, error) {
	return s.repo.GetQuests()
}

func (s *QuestService) CreateQuest(newquest internal.NewQuest) []internal.ErrorList {
	return s.repo.CreateQuest(newquest)
}

func (s *QuestService) CreateQuestStep(newQuestDB internal.NewQuestStepDB) error {
	return s.repo.CreateQuestStep(newQuestDB)
}

func (s *QuestService) CreateQuestSteps(newQuestSteps internal.NewQuestSteps) []internal.ErrorList {
	return s.repo.CreateQuestSteps(newQuestSteps)
}
func (s *QuestService) UpdateQuestSteps(updateQuestSteps internal.UpdateQuestSteps) []internal.ErrorList {
	return s.repo.UpdateQuestSteps(updateQuestSteps)
}
