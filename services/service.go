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
	CreateQuest(quest internal.NewQuest) (*internal.Quests, error)
	CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error
	UpdateQuestSteps(updateQuestSteps internal.UpdateQuestSteps) (*internal.UpdateQuestSteps, error)
}

type History interface {
	CompleteSteps(сompleteSteps internal.NewCompleteSteps) error
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
