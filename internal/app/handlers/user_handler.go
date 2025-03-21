package handlers

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/services"
	"encoding/json"
	"net/http"
	"strings"
)

type UserHandler struct {
	UserService *services.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendResponse(w, http.StatusMethodNotAllowed, "", false, "Method not allowed", nil)
		return
	}

	var usr models.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Invalid input", nil)
		return
	}
	if strings.TrimSpace(usr.Name) == "" {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Name is required", nil)
		return
	}
	if strings.TrimSpace(usr.LoginID) == "" {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "LoginID is required", nil)
		return
	}
	if strings.TrimSpace(usr.Password) == "" || len(usr.Password) < 6 {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Password must be at least 6 characters long", nil)
		return
	}
	// Call Service to Create User
	result, _ := h.UserService.CreateUser(usr)

	// Success Response
	helper.SendFinalResponse(w, result)
}

func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendResponse(w, http.StatusMethodNotAllowed, "", false, "Method not allowed", nil)
		return
	}

	var creds struct {
		LoginID  string `json:"login_id"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(creds.LoginID) == "" {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "LoginID is required", nil)
		return
	}
	if strings.TrimSpace(creds.Password) == "" || len(creds.Password) < 6 {
		helper.SendResponse(w, http.StatusBadRequest, "", false, "Password must be at least 6 characters long", nil)
		return
	}
	// Call Service to Create User
	result, _ := h.UserService.AuthenticateUser(creds.LoginID, creds.Password)

	// Success Response
	helper.SendFinalResponse(w, result)
}
