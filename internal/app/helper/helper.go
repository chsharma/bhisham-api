package helper

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, statusCode int, data interface{}, status bool, message string, additionalFields map[string]interface{}) {
	if data == nil {
		data = []interface{}{}
	}

	response := map[string]interface{}{
		"data":       data,
		"status":     status,
		"success":    status,
		"message":    message,
		"statusCode": statusCode,
	}

	// Add additional fields to the response if present
	for key, value := range additionalFields {
		response[key] = value
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
func CreateDynamicResponse(message string, success bool, data interface{}, status_code int, additionalFields map[string]interface{}) map[string]interface{} {
	// Ensure data is an empty slice if nil
	if data == nil {
		data = []interface{}{}
	}

	response := map[string]interface{}{
		"data":       data,
		"message":    message,
		"success":    success,
		"statusCode": status_code,
	}

	// Add additional fields to the response if present
	for key, value := range additionalFields {
		response[key] = value
	}

	return response
}

func SendFinalResponse(w http.ResponseWriter, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// Ensure type assertion for status code as int
	if statusCode, ok := response["statusCode"].(int); ok {
		w.WriteHeader(statusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError) // Fallback for invalid status code
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
