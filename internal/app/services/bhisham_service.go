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

func (s *BhishamService) CreateBhishamData(BhishamID int, UserID string) (map[string]interface{}, error) {
	return s.GameRepo.CreateBhishamData(BhishamID, UserID)
}

func (s *BhishamService) UpdateBhishamData(obj models.UpdateBhishamData, UserID string) (map[string]interface{}, error) {
	return s.GameRepo.UpdateBhishamData(obj, UserID)
}

func (s *BhishamService) MarkUpdateBhishamData(obj models.UpdateBhishamData, UserID string) (map[string]interface{}, error) {
	return s.GameRepo.MarkUpdateBhishamData(obj, UserID)
}

func (s *BhishamService) UpdateBhishamMapping(obj models.UpdateBhishamData, UserID string) (map[string]interface{}, error) {
	return s.GameRepo.UpdateBhishamMapping(obj, UserID)
}

func (s *BhishamService) MarkUpdateBhishamMapping(obj models.UpdateBhishamData, UserID string) (map[string]interface{}, error) {
	return s.GameRepo.MarkUpdateBhishamMapping(obj, UserID)
}
