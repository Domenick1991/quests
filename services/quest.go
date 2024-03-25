package services

import (
	"fmt"
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
	var quests []internal.Quests
	quests, err := s.repo.GetQuestWithoutStep()
	if err != nil {
		return quests, err
	}
	for i, q := range quests {
		steps, err := s.repo.GetSteps(q)
		if err == nil {
			quests[i].Steps = steps
		}

	}
	return quests, nil
}

func (s *QuestService) CreateQuest(newQuest internal.NewQuest) error {
	questDB, err := newQuest.ConvertToDB()
	if err != nil {
		return err
	}
	oldQuestId := s.repo.CheckQuest(questDB)
	if oldQuestId != 0 {
		return fmt.Errorf("Задание с таким именем существует, id :%d", oldQuestId)
	}
	questid, err := s.repo.CreateQuest(questDB)
	if err != nil {
		return err
	}
	if newQuest.QuestSteps != nil {
		//Если передавалась информация о шагах - добавляем и шаги
		for _, questStep := range newQuest.QuestSteps {
			questStep.QuestId = questid
			err = s.repo.CreateQuestStep(questStep)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *QuestService) CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error {
	return s.repo.CreateQuestSteps(newQuestSteps)
}
func (s *QuestService) UpdateQuestSteps(updateQuestSteps internal.UpdateQuestSteps) error {
	for _, questStep := range updateQuestSteps.QuestSteps {
		err := s.repo.UpdateQuestSteps(questStep)
		if err != nil {
			return err
		}
	}
	return nil
}
