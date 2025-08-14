package handler

import (
	"encoding/json"
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/go-playground/validator/v10"
)

type ProfileHandler struct {
	service   service.ProfileServiceInterface
	validator *validator.Validate
}

func NewProfileHandler(s service.ProfileServiceInterface) *ProfileHandler {
	return &ProfileHandler{
		service:   s,
		validator: validator.New(),
	}
}

// @Summary Get User Profile
// @Description Retrieves the profile of the currently authenticated user.
// @Tags Profile
// @Produce json
// @Success 200 {object} model.User "User profile retrieved successfully"
// @Failure 401 {object} model.ErrorResponse "Unauthorized access"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /me [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "profile",
		"action", "getProfile",
	)
	logger.Info("processing get user profile request")

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		logger.Error("user ID missing inside of handler")
		util.SendErrorResponse(
			w,
			errs.NewInternalError(
				"user ID missing inside of handler",
				nil,
				map[string]any{},
			),
		)
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		logger.Error("failed to get user profile",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	util.SendJSONResponse(w, http.StatusOK, user)
}

// @Summary Update User Profile
// @Description Updates the profile of the currently authenticated user.
// @Tags Profile
// @Accept json
// @Produce json
// @Param user body model.UpdateUserRequest true "User profile data"
// @Success 200 {object} model.User "User profile updated successfully"
// @Failure 400 {object} model.ErrorResponse "Bad request: Validation error"
// @Failure 401 {object} model.ErrorResponse "Unauthorized access"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /me [patch]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "profile",
		"action", "updateProfile",
	)
	logger.Info("processing update user profile request")

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		logger.Error("user ID missing inside of handler")
		util.SendErrorResponse(
			w,
			errs.NewInternalError(
				"user ID missing inside of handler",
				nil,
				map[string]any{},
			),
		)
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("failed to decode request body",
			"error", err.Error(),
		)
		util.SendErrorResponse(
			w,
			errs.NewValidationError("Invalid request payload", map[string]any{}),
		)
		return
	}

	logger.Debug("validating request")
	if err := h.validator.Struct(req); err != nil {
		logger.Warn("validation failed",
			"error", err.Error(),
		)
		util.SendErrorResponse(
			w,
			errs.NewValidationError("Invalid input data", map[string]any{}),
		)
		return
	}

	user, err := h.service.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		logger.Error("failed to update user profile",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("user profile updated successfully",
		"userID", user.ID,
		"username", user.Username,
	)
	util.SendJSONResponse(w, http.StatusOK, user)
}

// @Summary Delete User Profile
// @Description Deletes the profile of the currently authenticated user.
// @Tags Profile
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} model.ErrorResponse "Unauthorized access"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /me [delete]
func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "profile",
		"action", "deleteProfile",
	)
	logger.Info("processing delete user profile request")

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		logger.Error("user ID missing inside of handler")
		util.SendErrorResponse(
			w,
			errs.NewInternalError(
				"user ID missing inside of handler",
				nil,
				map[string]any{},
			),
		)
		return
	}

	err := h.service.DeleteProfile(r.Context(), userID)
	if err != nil {
		logger.Error("failed to delete user profile",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("user profile deleted successfully",
		"userID", userID,
	)
	w.WriteHeader(http.StatusNoContent)
}
