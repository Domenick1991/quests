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
	CreateQuest(quest internal.NewQuest) error
	CreateQuestStep(newQuestDB internal.NewQuestStepDB) error
	CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error
	UpdateQuestSteps(updateQuestSteps internal.UpdateQuestSteps) error
}
type History interface {
	CompleteSteps(сompleteSteps internal.NewCompleteSteps) error
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
