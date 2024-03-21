package repository

import (
	"errors"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"quests/internal"
	"strconv"
)

type QuestRepo struct {
	db *dbx.DB
}

func NewQuestRepo(db *dbx.DB) *QuestRepo {
	return &QuestRepo{db: db}
}

func (quest *QuestRepo) GetQuests() ([]internal.Quests, error) {
	var quests []internal.Quests
	err := quest.db.Select().From("quests").All(&quests)
	if err != nil {
		return quests, errors.New("Ошибка при получении данных о заданиях")
	}
	for i, q := range quests {
		var steps []internal.Steps
		err = quest.db.Select().From("queststeps").Where(dbx.HashExp{"questid": q.Id}).All(&steps)
		if err != nil {
			return quests, errors.New("Ошибка при получении данных о шагах заданий")
		}
		quests[i].Steps = steps
	}
	return quests, nil
}

func (quest *QuestRepo) CreateQuest(newquest internal.NewQuest) error {
	questDB, err := newquest.ConvertToDB()
	if err != nil {
		return err
	}

	var oldquestDB = internal.NewQuestDB{}
	err = quest.db.Select().From(questDB.TableName()).Where(dbx.HashExp{"questname": questDB.Name}).One(&oldquestDB)
	if err != nil {
		err = quest.db.Model(&questDB).Insert("Name", "Cost")
		if err != nil {
			return err
		} else {
			//Если передавалась информация о шагах - добавляем и шаги
			if newquest.QuestSteps != nil {
				for _, questStep := range newquest.QuestSteps {
					questStep.QuestId = questDB.Id
					questStepDB, err := questStep.ConvertToDB()
					if err != nil {
						return err
					}
					err = quest.CreateQuestStep(questStepDB)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}
	} else {
		return fmt.Errorf("Задание с таким именем существует, id :%s", strconv.Itoa(oldquestDB.Id))
	}
}

func (quest *QuestRepo) CreateQuestStep(newQuestStepDB internal.NewQuestStepDB) error {
	var oldquestStepDB = internal.NewQuestStepDB{}
	err := quest.db.Select().From(newQuestStepDB.TableName()).Where(dbx.HashExp{"stepname": newQuestStepDB.StepName, "questid": newQuestStepDB.QuestId}).One(&oldquestStepDB)
	if err != nil {
		err = quest.db.Model(&newQuestStepDB).Insert("QuestId", "StepName", "Bonus", "IsMulti")
		if err != nil {
			return errors.New("Не удалось добавить задание")
		}
	} else {
		return errors.New(fmt.Sprint("Не удалось добавить шаг '", newQuestStepDB.StepName, "' т.к. шаг с именем '", newQuestStepDB.StepName, "' уже сщуествует"))
	}
	return nil
}

func (quest *QuestRepo) CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error {
	for _, questStep := range newQuestSteps.QuestSteps {
		questStepDB, err := questStep.ConvertToDB()
		if err != nil {
			return err
		}
		err = quest.CreateQuestStep(questStepDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// @Summary Обновляет информацию о шагах заданий
func (quest *QuestRepo) UpdateQuestSteps(updateQuestSteps internal.UpdateQuestSteps) error {
	for _, questStep := range updateQuestSteps.QuestSteps {
		questStepDB, err := questStep.ConvertToDB()
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
	}
	return nil
}
