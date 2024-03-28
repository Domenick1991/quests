package repository

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"quests/internal"
)

type Authorization interface {
	CreateUser(user internal.User) (int, error)
	GetAllUser() []internal.User
	DeleteUser(id int) error
	GetUser(username, password string) (internal.User, error)
}

type Quest interface {
	GetQuestWithoutStep() ([]internal.Quests, error)
	GetSteps(q internal.Quests) ([]internal.Steps, error)
	CheckQuest(questDB internal.NewQuestDB) int
	CreateQuest(questDB internal.NewQuestDB) (int, error)
	CreateQuestStep(newQuest internal.NewQuestStep) (int, error)
	CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error
	UpdateQuestSteps(updateQuestStep internal.UpdateQuestStep) error
}
type History interface {
	CompleteSteps(сompletedStep internal.CompletedStep) error
	GetCompletedQuest(userId int) []internal.Quests
	GetCompletedQuestForUser(userId int, quest internal.Quests) internal.UserCompletedQuest
}

type Repository struct {
	Authorization
	Quest
	History
}

func NewRepository(db *dbx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepo(db),
		Quest:         NewQuestRepo(db),
		History:       NewHistoryRepo(db),
	}
}
