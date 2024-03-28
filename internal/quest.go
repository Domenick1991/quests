package internal

import "errors"

// Quests model info
// @Description Quests json информация о заданиях и их шагов
type Quests struct {
	Id        string  `json:"Id" db:"id"`               //ИД задания
	QuestName string  `json:"QuestName" db:"questname"` //Имя выполненного задания пользователем
	Steps     []Steps `json:"Steps" db:"-`              //Шаги задания
}

func (quest *Quests) TableName() string {
	return "quests"
}

type Steps struct {
	StepName string `json:"StepName" db:"stepname"` //Имя шага
	Id       int    `json:"Id" db:"id"`             //ИД шага
	Bonus    int    `json:"Bonus" db:"bonus"`       //Бонус за выполнение шага
	IsMulti  bool   `json:"isMulti" db:"ismulti"`   //Признак того, что шаг можно выполнять повторно
}

// region типы для выполнения шагов

// NewCompleteSteps model info
// @Description NewCompleteSteps  json для отметки о выполнении шага задания пользователем
type NewCompleteSteps struct {
	CompleteSteps []CompletedStep `json:"CompleteSteps"` //Идентификатор задания
}

type CompletedStep struct {
	Stepid int `json:"stepid"` //Идентификатор шага
	Userid int `json:"userid"` //Идентификатор пользователя выполневшего шаг
}

func (complete *CompletedStep) ConvertToDB() (CompletedStepDB, error) {
	completedDB := CompletedStepDB{}
	if complete.Stepid == 0 {
		return completedDB, errors.New("Идентификатор шага должен быть больше 0")
	}
	if complete.Userid == 0 {
		return completedDB, errors.New("Идентификатор пользователя должен быть больше 0")
	}
	completedDB.Stepid = complete.Stepid
	completedDB.Userid = complete.Userid
	return completedDB, nil
}

type CompletedStepDB struct {
	Stepid int `json:"stepid"` //Идентификатор шага
	Userid int `json:"userid"` //Идентификатор пользователя выполневшего шаг
}

func (quest *CompletedStepDB) TableName() string {
	return "history"
}

//endregion типы для выполнения шагов

//region типы для создания новых заданий

// NewQuest model info
// @Description NewQuest json для создания задания с шагами
type NewQuest struct {
	Id         int            `json:"id"`         //Идентификатор задания
	Name       string         `json:"Name"`       //Имя задания
	QuestSteps []NewQuestStep `json:"QuestSteps"` //Шаги задания
}

func (quest *NewQuest) ConvertToDB() (NewQuestDB, error) {
	questdb := NewQuestDB{}
	questdb.Id = quest.Id

	if quest.Name == "" {
		return questdb, errors.New("Имя задания должно содержать от 1 до 200 символов")
	}
	questdb.Name = quest.Name
	return questdb, nil
}

type NewQuestDB struct {
	Id   int    `json:"id" db:"id"`          //идентификатор задания
	Name string `json:"Name" db:"questname"` //Имя задания
}

func (quest *NewQuestDB) TableName() string {
	return "quests"
}

//endregion типы для создания новых заданий

//region типы для создания и обновления шагов заданий

// NewQuestSteps model info
// @Description NewQuestStep json для создания шага задания
type NewQuestSteps struct {
	QuestSteps []NewQuestStep `json:"QuestSteps"` //Идентификатор задания
}

type NewQuestStep struct {
	Id       int    `json:"id"`       //Идентификатор задания
	QuestId  int    `json:"QuestId"`  //Идентификатор задания. При создании методом CreateQuest, значение будет проигнорировано, т.к. будет подставляться идентификатор создаваемого задания
	StepName string `json:"StepName"` //Описание шага
	Bonus    int    `json:"Bonus"`    //Бонус за задание
	IsMulti  *bool  `json:"IsMulti"`  //Признак того, что шаг можно выполнять несколько раз
}

func (questStep *NewQuestStep) ConvertToDB() (NewQuestStepDB, error) {
	questStepDB := NewQuestStepDB{}
	questStepDB.Id = questStep.Id

	if questStep.StepName == "" {
		return questStepDB, errors.New("Не указано описание шага")
	}
	if questStep.QuestId <= 0 {
		return questStepDB, errors.New("Не указан идентификатор задания, к которому относится шаг")
	}

	if questStep.Bonus < 0 {
		return questStepDB, errors.New("Бонус не может быть меньше 0")
	}

	questStepDB.IsMulti = *questStep.IsMulti
	if questStep.IsMulti == nil {
		questStepDB.IsMulti = false
	}

	questStepDB.QuestId = questStep.QuestId
	questStepDB.Bonus = questStep.Bonus
	questStepDB.StepName = questStep.StepName

	return questStepDB, nil

}

// UpdateQuestSteps model info
// @Description UpdateQuestSteps json для обновления шагов заданий
type UpdateQuestSteps struct {
	QuestSteps []UpdateQuestStep `json:"QuestSteps"` //Идентификатор задания
}

type UpdateQuestStep struct {
	Id      int   `json:"id"`      //Идентификатор задания
	Bonus   int   `json:"Bonus"`   //Бонус за задание
	IsMulti *bool `json:"IsMulti"` //Признак того, что шаг можно выполнять несколько раз
}

func (questStep *UpdateQuestStep) ConvertToDB() (NewQuestStepDB, error) {
	questStepDB := NewQuestStepDB{}
	questStepDB.Id = questStep.Id

	if questStep.Id == 0 {
		return questStepDB, errors.New("Не указан идентификатор шага, который необходимо обновить")
	}
	if questStep.IsMulti == nil {
		return questStepDB, errors.New("кажите признак многократного выполнения")
	}
	questStepDB.Bonus = questStep.Bonus
	questStepDB.IsMulti = *questStep.IsMulti
	return questStepDB, nil

}

type NewQuestStepDB struct {
	Id       int    `json:"id" db:"id"`
	QuestId  int    `json:"QuestId" db:"questid"`
	StepName string `json:"StepName" db:"stepname"`
	Bonus    int    `json:"Bonus" db:"bonus"`
	IsMulti  bool   `json:"IsMulti" db:"ismulti"`
}

func (quest *NewQuestStepDB) TableName() string {
	return "queststeps"
}

// GetUpdatesData Функция возвращает Мар содержащую только измененные поля, которые необходимо записать в БД. Если поле не было передано для обновления, то и записываться в базу оно не будет
func (questStep *NewQuestStepDB) GetUpdatesData() map[string]interface{} {
	var data = make(map[string]interface{})
	if questStep.Bonus > 0 {
		data["bonus"] = questStep.Bonus
	}
	data["ismulti"] = questStep.IsMulti
	return data
}

//endregion типы для создания и обновления шагов заданий

type ErrorList struct {
	Error string
}
