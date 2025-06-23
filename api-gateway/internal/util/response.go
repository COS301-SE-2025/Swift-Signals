package util

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
)

// Helper function for sending json responses
func SendJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			// Fallback error response if encoding fails, to avoid leaving the client hanging
			http.Error(w, `{"message":"Internal server error encoding response"}`, http.StatusInternalServerError)
		}
	}
}

// Send consistent error messages using model.ErrorResponse
func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errResp := model.ErrorResponse{
		Message: message,
	}
	switch statusCode {
	case http.StatusBadRequest:
		errResp.Code = "BAD_REQUEST"
	case http.StatusNotFound:
		errResp.Code = "NOT_FOUND"
	case http.StatusUnauthorized:
		errResp.Code = "UNAUTHORIZED"
	case http.StatusForbidden:
		errResp.Code = "FORBIDDEN"
	case http.StatusInternalServerError:
		errResp.Code = "INTERNAL_SERVER_ERROR"
	default:
		errResp.Code = "UNKNOWN_ERROR"
	}
	SendJSONResponse(w, statusCode, errResp)
}
