package services

import (
	"fmt"
	"quests/internal"
	"quests/repository"
	"strconv"
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

type QuestExsistErr struct {
	questID int
}

func (q QuestExsistErr) Error() string {
	return fmt.Sprintf("Задание с таким именем существует, id :%d", q.questID)
}

func (s *QuestService) CreateQuest(newQuest internal.NewQuest) (*internal.Quests, error) {

	questDB, err := newQuest.ConvertToDB()
	if err != nil {
		return nil, err
	}
	oldQuestId := s.repo.CheckQuest(questDB)
	if oldQuestId != 0 {
		return nil, QuestExsistErr{oldQuestId}
	}
	questID, err := s.repo.CreateQuest(questDB)
	if err != nil {
		return nil, err
	}

	var quest = internal.Quests{
		Id:        strconv.Itoa(questID),
		QuestName: newQuest.Name,
		Steps:     make([]internal.Steps, len(newQuest.QuestSteps)),
	}

	if newQuest.QuestSteps != nil {

		//Если передавалась информация о шагах - добавляем и шаги
		for i, questStep := range newQuest.QuestSteps {
			questStep.QuestId = questID
			stepID, err := s.repo.CreateQuestStep(questStep)
			if err != nil {
				return nil, err
			}
			step := internal.Steps{Id: stepID, StepName: questStep.StepName, Bonus: questStep.Bonus, IsMulti: questStep.IsMulti}
			quest.Steps[i] = step
		}
	}
	return &quest, nil
}

func (s *QuestService) CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error {
	return s.repo.CreateQuestSteps(newQuestSteps)
}

func (s *QuestService) UpdateQuestSteps(updateQuestSteps internal.UpdateQuestSteps) (*internal.UpdateQuestSteps, error) {
	for _, questStep := range updateQuestSteps.QuestSteps {
		err := s.repo.UpdateQuestSteps(questStep)
		if err != nil {
			return nil, err
		}
	}
	return &updateQuestSteps, nil
}
