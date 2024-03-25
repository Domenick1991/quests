package repository

import (
	"errors"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"quests/internal"
)

type QuestRepo struct {
	db *dbx.DB
}

func NewQuestRepo(db *dbx.DB) *QuestRepo {
	return &QuestRepo{db: db}
}

func (quest *QuestRepo) GetQuestWithoutStep() ([]internal.Quests, error) {
	var quests []internal.Quests
	err := quest.db.Select().From("quests").All(&quests)
	if err != nil {
		return quests, errors.New("Ошибка при получении данных о заданиях")
	}
	return quests, nil
}

func (quest *QuestRepo) GetSteps(q internal.Quests) ([]internal.Steps, error) {
	var steps []internal.Steps
	err := quest.db.Select().From("queststeps").Where(dbx.HashExp{"questid": q.Id}).All(&steps)
	if err != nil {
		return steps, errors.New("Ошибка при получении данных о шагах заданий")
	}
	return steps, nil
}

func (quest *QuestRepo) CheckQuest(questDB internal.NewQuestDB) int {
	var oldquestDB = internal.NewQuestDB{}
	err := quest.db.Select().From(questDB.TableName()).Where(dbx.HashExp{"questname": questDB.Name}).One(&oldquestDB)
	if err != nil {
		return 0 //не существует
	}
	return oldquestDB.Id //существует
}

func (quest *QuestRepo) CreateQuest(questDB internal.NewQuestDB) (int, error) {
	err := quest.db.Model(&questDB).Insert("Name", "Cost")
	if err != nil {
		return 0, err
	}
	return questDB.Id, nil

}

func (quest *QuestRepo) CreateQuestStep(questStep internal.NewQuestStep) (int, error) {
	questStepDB, err := questStep.ConvertToDB()
	if err != nil {
		return 0, err
	}
	stepID, err := quest.createStep(questStepDB)
	if err != nil {
		return 0, err
	}
	return stepID, nil
}

func (quest *QuestRepo) createStep(newQuestStepDB internal.NewQuestStepDB) (int, error) {
	var oldquestStepDB = internal.NewQuestStepDB{}
	err := quest.db.Select().From(newQuestStepDB.TableName()).Where(dbx.HashExp{"stepname": newQuestStepDB.StepName, "questid": newQuestStepDB.QuestId}).One(&oldquestStepDB)
	if err != nil {
		err = quest.db.Model(&newQuestStepDB).Insert("QuestId", "StepName", "Bonus", "IsMulti")
		if err != nil {
			return 0, errors.New("Не удалось добавить задание")
		}
	} else {
		return 0, errors.New(fmt.Sprint("Не удалось добавить шаг '", newQuestStepDB.StepName, "' т.к. шаг уже существует"))
	}
	return newQuestStepDB.Id, nil
}

func (quest *QuestRepo) CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error {
	for _, questStep := range newQuestSteps.QuestSteps {
		_, err := quest.CreateQuestStep(questStep)
		if err != nil {
			return err
		}
	}
	return nil
}

// @Summary Обновляет информацию о шаге задания
func (quest *QuestRepo) UpdateQuestSteps(updateQuestStep internal.UpdateQuestStep) error {
	questStepDB, err := updateQuestStep.ConvertToDB()
	if err != nil {
		return err
	}
	params := questStepDB.GetUpdatesData()
	if len(params) > 0 {
		_, err = quest.db.Update(questStepDB.TableName(), params, dbx.HashExp{"id": questStepDB.Id}).Execute()
		if err != nil {
			return err
		}
	}
	return nil
}
