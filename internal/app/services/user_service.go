package services

import (
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/repositories"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

// CreateGame delegates the creation of a new game to the repository.
func (s *UserService) CreateUser(usr models.User) (map[string]interface{}, error) {
	return s.UserRepo.CreateUser(usr)
}

func (s *UserService) UpdateUser(usr models.User) (map[string]interface{}, error) {
	return s.UserRepo.UpdateUser(usr)
}

// CreateGame delegates the creation of a new game to the repository.
func (s *UserService) AuthenticateUser(loginID, password string) (map[string]interface{}, error) {
	return s.UserRepo.AuthenticateUser(loginID, password)
}

func (s *UserService) GetUsers() (map[string]interface{}, error) {
	return s.UserRepo.GetUsers()
}

func (s *UserService) UpdatePassword(loginID, password string) (map[string]interface{}, error) {
	return s.UserRepo.UpdatePassword(loginID, password)
}
