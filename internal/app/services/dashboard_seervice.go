package services

import (
	"bhisham-api/internal/app/repositories"
)

type DashboardService struct {
	DashboardRepo *repositories.DashboardRepository
}

func (s *DashboardService) GetDashboardStats() (map[string]interface{}, error) {
	return s.DashboardRepo.GetDashboardStats()
}

func (s *DashboardService) GetBhisham() (map[string]interface{}, error) {
	return s.DashboardRepo.GetBhisham()
}

func (s *DashboardService) GetChildCube(BhishamID, MotherCubeID int) (map[string]interface{}, error) {
	return s.DashboardRepo.GetChildCube(BhishamID, MotherCubeID)
}

func (s *DashboardService) GetChildKits(BhishamID, MotherCubeID, CCNo int) (map[string]interface{}, error) {
	return s.DashboardRepo.GetChildKits(BhishamID, MotherCubeID, CCNo)
}

func (s *DashboardService) GetKitItems(BhishamID, MotherCubeID, CCNo int, KitName string) (map[string]interface{}, error) {
	return s.DashboardRepo.GetKitItems(BhishamID, MotherCubeID, CCNo, KitName)
}

func (s *DashboardService) GetMappingKitItems(BhishamID, MotherCubeID, CCNo int, KitName string) (map[string]interface{}, error) {
	return s.DashboardRepo.GetMappingKitItems(BhishamID, MotherCubeID, CCNo, KitName)
}

func (s *DashboardService) GetAllBhishamData(BhishamID int) (map[string]interface{}, error) {
	return s.DashboardRepo.GetAllBhishamData(BhishamID)
}

func (s *DashboardService) GetUpdateType() (map[string]interface{}, error) {
	return s.DashboardRepo.GetUpdateType()
}

func (s *DashboardService) GetBhishamID(SerialNo string) (map[string]interface{}, error) {
	return s.DashboardRepo.GetBhishamID(SerialNo)
}
