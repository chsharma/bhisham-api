package services

import (
	"bhisham-api/internal/app/repositories"
)

type HandheldService struct {
	HandheldRepo *repositories.HandheldRepository
}

func (s *HandheldService) GetBhishamID(SerialNo string) (map[string]interface{}, error) {
	return s.HandheldRepo.GetBhishamID(SerialNo)
}

func (s *HandheldService) GetAllBhishamData(BhishamID int) (map[string]interface{}, error) {
	return s.HandheldRepo.GetAllBhishamData(BhishamID)
}
