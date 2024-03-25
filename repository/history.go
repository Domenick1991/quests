package repository

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"quests/internal"
	"strconv"
)

type HistoryRepo struct {
	db *dbx.DB
}

func NewHistoryRepo(db *dbx.DB) *HistoryRepo {
	return &HistoryRepo{db: db}
}

func (history *HistoryRepo) CompleteSteps(сompleteStep internal.CompleteStep) error {
	сompleteStepDB, err := сompleteStep.ConvertToDB()
	if err != nil {
		return err
	}
	//Если задание можно выполнить - выполняем, если нет, то просто игнорируем
	if checkCompliteStep(history.db, сompleteStep) {
		err = history.db.Model(&сompleteStepDB).Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

// ExicuteCountSumQuery Обёртка предназначена для запросов, которые возвращают одно целое значение
func ExicuteCountSumQuery(db *dbx.DB, queryText string) int {
	result := 0
	query := db.NewQuery(queryText)
	rows, err := query.Rows()
	if err != nil {
		return result
	}
	for rows.Next() {
		rows.Scan(&result)
	}
	return result
}

// checkCompliteStep Устанавливает признак выполнения у шага
func checkCompliteStep(DB *dbx.DB, сompleteStep internal.CompleteStep) bool {
	//Получаем id всех выполненных заданий

	queryText := "SELECT s.id FROM public.queststeps as s where (s.ismulti = false and s.id not in (select h.stepid from history h where h.userid = " + strconv.Itoa(сompleteStep.Userid) + " )) or s.ismulti = true"
	query := DB.NewQuery(queryText)
	rows, err := query.Rows()
	if err != nil {

	}
	var stepIds map[int]bool = make(map[int]bool)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		stepIds[id] = true
	}
	_, inMap := stepIds[сompleteStep.Stepid]
	return inMap
}

// getCompletedQuestId возвращает ИД заданий в которых участвовал пользователь
func (history *HistoryRepo) GetCompletedQuest(userId int) []internal.Quests {
	var quests []internal.Quests

	queryText := `SELECT distinct q.id, q.questname
						FROM public.queststeps as s
						left join history as h on s.id = h.stepid
						left join quests as q on s.questid = q.id
						where h.userid = ` + strconv.Itoa(userId)
	query := history.db.NewQuery(queryText)
	rows, err := query.Rows()
	if err != nil {
		return quests
	}

	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		quests = append(quests, internal.Quests{id, name, nil})
	}
	return quests
}

// GetCompletedQuestForUser Возвращает информацию по заданию для пользователя
func (history *HistoryRepo) GetCompletedQuestForUser(userId int, quest internal.Quests) internal.UserCompletedQuest {
	UserCompletedQuest := internal.UserCompletedQuest{}

	UserCompletedQuest.QuestId = quest.Id
	UserCompletedQuest.QuestName = quest.QuestName

	//Всего заданий
	queryText := `	SELECT count(*)
					FROM public.queststeps
					where questid = ` + quest.Id
	UserCompletedQuest.AllStepsCount = ExicuteCountSumQuery(history.db, queryText)

	//Всего заданий выполненных пользователей
	queryText = `SELECT count(*)
					FROM public.queststeps as s
					left join history as h on s.id = h.stepid
					left join quests as q on s.questid = q.id
					where h.userid = ` + strconv.Itoa(userId) + ` and q.id = ` + quest.Id
	UserCompletedQuest.CompletedStepsCount = ExicuteCountSumQuery(history.db, queryText)

	//Сумма бонуса за выполненные шаги задания
	queryText = `SELECT Sum(s.bonus)
					FROM public.queststeps as s
					left join history as h on s.id = h.stepid
					left join quests as q on s.questid = q.id
					where h.userid = ` + strconv.Itoa(userId) + ` and q.id = ` + quest.Id
	UserCompletedQuest.Bonus = ExicuteCountSumQuery(history.db, queryText)

	//region пройдемся по каждому выполненному шагу пользователя и посчитаем сколько раз был выполнен каждый шаг и сумму бонусов за это
	queryText = `SELECT distinct (s.id)
					FROM public.queststeps as s
					left join history as h on s.id = h.stepid
					left join quests as q on s.questid = q.id
					where h.userid = ` + strconv.Itoa(userId) + ` and q.id = ` + quest.Id
	query := history.db.NewQuery(queryText)
	rowsSteps, _ := query.Rows()
	CompletedStepsCount := 0
	type stepInfo struct{ count, bonus int }
	for rowsSteps.Next() {
		var stepid string
		rowsSteps.Scan(&stepid)

		queryText = `SELECT s.stepname, s.bonus
					FROM public.queststeps as s
					left join history as h on s.id = h.stepid
					left join quests as q on s.questid = q.id
					where h.userid = ` + strconv.Itoa(userId) + ` and q.id = ` + quest.Id + ` and s.id =` + stepid
		query = history.db.NewQuery(queryText)
		rows, _ := query.Rows()

		var stepsInfo map[string]stepInfo = make(map[string]stepInfo)
		for rows.Next() {
			var stepname string
			var bonus int
			rows.Scan(&stepname, &bonus)
			if value, inMap := stepsInfo[stepname]; inMap {
				value.count++
				value.bonus += bonus
				stepsInfo[stepname] = value
			} else {
				stepsInfo[stepname] = stepInfo{1, bonus}
			}

		}
		for key, value := range stepsInfo {
			UserCompletedQuest.CompletedSteps = append(UserCompletedQuest.CompletedSteps, internal.UserCompletedSteps{key, value.count, value.bonus})
		}

		CompletedStepsCount++
	}
	//endregion

	UserCompletedQuest.CompletedStepsCount = CompletedStepsCount

	return UserCompletedQuest
}
