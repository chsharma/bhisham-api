package services

import (
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/repositories"
)

type BhishamService struct {
	GameRepo *repositories.BhishamRepository
}

// CreateGame delegates the creation of a new game to the repository.
func (s *BhishamService) CreateBhisham(game models.Bhisham) (map[string]interface{}, error) {
	return s.GameRepo.CreateBhisham(game)
}
