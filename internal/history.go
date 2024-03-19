package internal

// UserBonus model info
// @Description UserBonus json для получения история выполнения заданий и их шагов
type UserBonus struct {
	TotalBonus      int                  `json:"TotalBonus"`      //Общий бонусный счет пользователя
	CompletedQuests []UserCompletedQuest `json:"ComplitedQuests"` //Список заданий в которых участвовал пользователь
}

type UserCompletedQuest struct {
	QuestId             string               `json:"QuestId"`             //ИД задания
	QuestName           string               `json:"QuestName"`           //Имя выполненного задания пользователем
	Bonus               int                  `json:"Bonus"`               //Сумма Бонусов за выполненные задания
	CompletedStepsCount int                  `json:"CompletedStepsCount"` //Кол-во выполненных шагов заданий пользователем
	AllStepsCount       int                  `json:"AllStepsCount"`       //Кол-во шагов, доступное в задании
	CompletedSteps      []UserCompletedSteps `json:"CompletedSteps"`      //Выполненные шаги пользователем
}
type UserCompletedSteps struct {
	StepName      string `json:"StepName"`      //Имя выполненного шага
	Count         int    `json:"Count"`         //Кол-во выполнений шага
	UserBonusStep int    `json:"UserBonusStep"` //Бонус пользователя за выполнение шага
}
