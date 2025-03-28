package handlers

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/services"
	"net/http"
	"strconv"
)

type DashboardHandler struct {
	DashboardService *services.DashboardService
}

func (c *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}
	// Call service with parsed values
	result, _ := c.DashboardService.GetDashboardStats()
	helper.SendFinalResponse(w, result)
}

func (c *DashboardHandler) GetBhisham(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}
	// Call service with parsed values
	result, _ := c.DashboardService.GetBhisham()
	helper.SendFinalResponse(w, result)
}

func (c *DashboardHandler) GetChildCube(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}

	// Extract query parameters
	bhishamIDStr := r.URL.Query().Get("bhishamid")
	motherCubeIDStr := r.URL.Query().Get("mcno")

	// Validate parameters
	if bhishamIDStr == "" || motherCubeIDStr == "" {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Missing required parameters", nil)
		return
	}

	// Convert parameters to integers
	bhishamID, err1 := strconv.Atoi(bhishamIDStr)
	motherCubeID, err2 := strconv.Atoi(motherCubeIDStr)

	if err1 != nil || err2 != nil {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Invalid parameter format", nil)
		return
	}

	// Call service layer
	result, _ := c.DashboardService.GetChildCube(bhishamID, motherCubeID)
	// Send final response
	helper.SendFinalResponse(w, result)
}

func (c *DashboardHandler) GetChildKits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}

	// Extract query parameters
	bhishamIDStr := r.URL.Query().Get("bhishamid")
	motherCubeIDStr := r.URL.Query().Get("mcno")
	cubeNumberStr := r.URL.Query().Get("ccno")

	// Validate required parameters
	if bhishamIDStr == "" || motherCubeIDStr == "" || cubeNumberStr == "" {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Missing required parameters", nil)
		return
	}

	// Convert parameters to integers
	bhishamID, err1 := strconv.Atoi(bhishamIDStr)
	motherCubeID, err2 := strconv.Atoi(motherCubeIDStr)
	cubeNumber, err3 := strconv.Atoi(cubeNumberStr)

	if err1 != nil || err2 != nil || err3 != nil {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Invalid parameter format", nil)
		return
	}

	// Call service layer
	result, _ := c.DashboardService.GetChildKits(bhishamID, motherCubeID, cubeNumber)

	// Send final response
	helper.SendFinalResponse(w, result)
}

func (c *DashboardHandler) GetKitItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}

	// Extract query parameters
	bhishamIDStr := r.URL.Query().Get("bhishamid")
	motherCubeIDStr := r.URL.Query().Get("mcno")
	cubeNumberStr := r.URL.Query().Get("ccno")
	cubeKitNameStr := r.URL.Query().Get("kitslug")

	// Validate required parameters
	if bhishamIDStr == "" || motherCubeIDStr == "" || cubeNumberStr == "" || cubeKitNameStr == "" {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Missing required parameters", nil)
		return
	}

	// Convert parameters to integers
	bhishamID, err1 := strconv.Atoi(bhishamIDStr)
	motherCubeID, err2 := strconv.Atoi(motherCubeIDStr)
	cubeNumber, err3 := strconv.Atoi(cubeNumberStr)

	if err1 != nil || err2 != nil || err3 != nil {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Invalid parameter format", nil)
		return
	}

	// Call service layer
	result, _ := c.DashboardService.GetKitItems(bhishamID, motherCubeID, cubeNumber, cubeKitNameStr)
	// Send final response
	helper.SendFinalResponse(w, result)
}

func (c *DashboardHandler) GetMappingKitItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}

	// Extract query parameters
	bhishamIDStr := r.URL.Query().Get("bhishamid")
	motherCubeIDStr := r.URL.Query().Get("mcno")
	cubeNumberStr := r.URL.Query().Get("ccno")
	cubeKitNameStr := r.URL.Query().Get("kitslug")

	// Validate required parameters
	if bhishamIDStr == "" || motherCubeIDStr == "" || cubeNumberStr == "" || cubeKitNameStr == "" {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Missing required parameters", nil)
		return
	}

	// Convert parameters to integers
	bhishamID, err1 := strconv.Atoi(bhishamIDStr)
	motherCubeID, err2 := strconv.Atoi(motherCubeIDStr)
	cubeNumber, err3 := strconv.Atoi(cubeNumberStr)

	if err1 != nil || err2 != nil || err3 != nil {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Invalid parameter format", nil)
		return
	}

	// Call service layer
	result, _ := c.DashboardService.GetMappingKitItems(bhishamID, motherCubeID, cubeNumber, cubeKitNameStr)
	helper.SendFinalResponse(w, result)
}

func (c *DashboardHandler) GetUpdateType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}
	// Call service with parsed values
	result, _ := c.DashboardService.GetUpdateType()
	helper.SendFinalResponse(w, result)
}

func (c *DashboardHandler) GetAllMappingBhishamData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendResponse(w, http.StatusMethodNotAllowed, nil, false, "Method not allowed", nil)
		return
	}

	// Extract query parameters
	bhishamIDStr := r.URL.Query().Get("bhishamid")

	// Validate required parameters
	if bhishamIDStr == "" {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Missing required parameters", nil)
		return
	}

	// Convert parameters to integers
	bhishamID, err1 := strconv.Atoi(bhishamIDStr)
	if err1 != nil {
		helper.SendResponse(w, http.StatusBadRequest, nil, false, "Invalid parameter format", nil)
		return
	}

	// Call service layer
	result, _ := c.DashboardService.GetAllMappingBhishamData(bhishamID)
	// Send final response
	helper.SendFinalResponse(w, result)
}
