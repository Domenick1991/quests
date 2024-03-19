package services

import (
	"quests/internal"
	"quests/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

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

type Service struct {
	Authorization
	Quest
	History
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Quest:         NewQuestService(repos.Quest),
		History:       NewHistoryService(repos.History),
	}
}
