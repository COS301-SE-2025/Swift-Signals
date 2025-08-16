package util

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

// Helper function for sending json responses
func SendJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			// Fallback error response if encoding fails, to avoid leaving the client hanging
			http.Error(
				w,
				`{"message":"Internal server error encoding response"}`,
				http.StatusInternalServerError,
			)
		}
	}
}

func SendErrorResponse(w http.ResponseWriter, err error) {
	if err == nil {
		logger := slog.Default()
		logger.Warn("returning nil error to client")
	}

	errResp := model.ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "something went wrong",
	}

	var svcErr *errs.ServiceError
	if errors.As(err, &svcErr) {
		errResp.Message = svcErr.Message
		switch svcErr.Code {
		case errs.ErrValidation:
			errResp.Code = http.StatusBadRequest
		case errs.ErrNotFound:
			errResp.Code = http.StatusNotFound
		case errs.ErrAlreadyExists:
			errResp.Code = http.StatusConflict
		case errs.ErrUnauthorized:
			errResp.Code = http.StatusUnauthorized
		case errs.ErrForbidden:
			errResp.Code = http.StatusForbidden
		default:
			errResp.Code = http.StatusInternalServerError
			errResp.Message = "something went wrong"
		}
	}

	SendJSONResponse(w, errResp.Code, errResp)
}
