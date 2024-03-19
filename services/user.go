package services

import (
	"quests/internal"
	"quests/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user internal.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *AuthService) DeleteUser(userId int) error {
	return s.repo.DeleteUser(userId)
}

func (s *AuthService) GetAllUser() []internal.User {
	return s.repo.GetAllUser()
}

func (s *AuthService) GetUser(username, password string) (internal.User, error) {
	return s.repo.GetUser(username, password)
}
