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
	GetQuests() ([]internal.Quests, error)
	CreateQuest(quest internal.NewQuest) []internal.ErrorList
	CreateQuestStep(newQuestDB internal.NewQuestStepDB) error
	CreateQuestSteps(newQuestSteps internal.NewQuestSteps) []internal.ErrorList
	UpdateQuestSteps(updateQuestSteps internal.UpdateQuestSteps) []internal.ErrorList
}
type History interface {
	CompleteSteps(сompleteSteps internal.NewCompleteSteps) []internal.ErrorList
	GetHistory(userid int) (internal.UserBonus, error)
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
