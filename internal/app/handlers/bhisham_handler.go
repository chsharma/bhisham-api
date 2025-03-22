package handlers

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/services"
	"encoding/json"
	"net/http"
)

type BhishamHandler struct {
	BhishamService *services.BhishamService
}

// CreateGame handles the creation of a new game.
func (h *BhishamHandler) CreateBhisham(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendResponse(w, http.StatusMethodNotAllowed, "", false, "Method not allowed", nil)
		return
	}
	var bsm models.Bhisham
	err := json.NewDecoder(r.Body).Decode(&bsm)
	if err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Invalid input", nil)
		return
	}
	userID := r.Header.Get("User-ID")
	if userID == "" {
		helper.SendResponse(w, http.StatusBadRequest, []interface{}{}, false, "User-ID nof faound", nil)
		return
	}
	bsm.CreatedBy = &userID
	result, _ := h.BhishamService.CreateBhisham(bsm)
	helper.SendFinalResponse(w, result)
}

// CreateGame handles the creation of a new game.
func (h *BhishamHandler) CreateBhishamData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendResponse(w, http.StatusMethodNotAllowed, "", false, "Method not allowed", nil)
		return
	}
	var bsm models.GetBhishamID
	err := json.NewDecoder(r.Body).Decode(&bsm)
	if err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Invalid input", nil)
		return
	}
	userID := r.Header.Get("User-ID")
	if userID == "" {
		helper.SendResponse(w, http.StatusBadRequest, []interface{}{}, false, "User-ID nof faound", nil)
		return
	}

	result, _ := h.BhishamService.CreateBhishamData(bsm.BhishamID, userID)
	helper.SendFinalResponse(w, result)
}

func (h *BhishamHandler) UpdateBhishamData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendResponse(w, http.StatusMethodNotAllowed, "", false, "Method not allowed", nil)
		return
	}
	var bsm models.UpdateBhishamData
	err := json.NewDecoder(r.Body).Decode(&bsm)
	if err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Invalid input", nil)
		return
	}
	userID := r.Header.Get("User-ID")
	if userID == "" {
		helper.SendResponse(w, http.StatusBadRequest, []interface{}{}, false, "User-ID nof faound", nil)
		return
	}

	result, _ := h.BhishamService.UpdateBhishamData(bsm, userID)
	helper.SendFinalResponse(w, result)
}

func (h *BhishamHandler) UpdateBhishamMapping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendResponse(w, http.StatusMethodNotAllowed, "", false, "Method not allowed", nil)
		return
	}
	var bsm models.UpdateBhishamData
	err := json.NewDecoder(r.Body).Decode(&bsm)
	if err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Invalid input", nil)
		return
	}
	userID := r.Header.Get("User-ID")
	if userID == "" {
		helper.SendResponse(w, http.StatusBadRequest, []interface{}{}, false, "User-ID nof faound", nil)
		return
	}

	result, _ := h.BhishamService.UpdateBhishamMapping(bsm, userID)
	helper.SendFinalResponse(w, result)
}
