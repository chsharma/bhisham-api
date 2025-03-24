package handlers

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type HandheldHandler struct {
	HandheldService *services.HandheldService
}

func (c *HandheldHandler) GetBhishamID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendResponse(w, http.StatusMethodNotAllowed, "", false, "Method not allowed", nil)
		return
	}
	var bsm models.SerialNo
	err := json.NewDecoder(r.Body).Decode(&bsm)
	if err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Invalid input", nil)
		return
	}
	result, _ := c.HandheldService.GetBhishamID(bsm.SerialNo)
	helper.SendFinalResponse(w, result)
}

func (c *HandheldHandler) GetAllBhishamData(w http.ResponseWriter, r *http.Request) {
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
	result, _ := c.HandheldService.GetAllBhishamData(bhishamID)
	// Send final response
	helper.SendFinalResponse(w, result)
}
