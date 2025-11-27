package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// respondJSON sends a JSON response with the given status code and data.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// respondError sends a JSON error response with the given status code, error message, and details.
func respondError(w http.ResponseWriter, status int, errMsg string, details string) {
	respondJSON(w, status, ErrorResponse{
		Error:   errMsg,
		Message: details,
	})
}
