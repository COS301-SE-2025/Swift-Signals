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

type AdminHandler struct {
	service   service.AdminServiceInterface
	validator *validator.Validate
}

func NewAdminHandler(s service.AdminServiceInterface) *AdminHandler {
	return &AdminHandler{
		service:   s,
		validator: validator.New(),
	}
}

// @Summary Get All Users
// @Description Retrieves a paginated list of all users. Only accessible by admins.
// @Tags Admin
// @Accept json
// @Produce json
// @Param getAllUsersRequest body model.GetAllUsersRequest true "Pagination options"
// @Success 200 {array} model.User "List of users"
// @Failure 403 {object} model.ErrorResponse "Forbidden - Only admins can access this endpoint"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /admin/users [get]
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "admin",
		"action", "get_all_users",
	)
	logger.Info("processing get all users request")

	var req model.GetAllUsersRequest
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
			errs.NewValidationError("Invalid request parameters", map[string]any{}),
		)
		return
	}

	users, err := h.service.GetAllUsers(r.Context(), req.Page, req.PageSize)
	if err != nil {
		logger.Error("failed to get all users",
			"error", err.Error(),
		)
		util.SendErrorResponse(
			w,
			err,
		)
		return
	}

	logger.Info("successfully retrieved all users",
		"count", len(users),
	)
	util.SendJSONResponse(w, http.StatusOK, users)
}

// @Summary Get User by ID
// @Description Retrieves a user by their ID. Only accessible by admins.
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User "User details"
// @Failure 403 {object} model.ErrorResponse "Forbidden - Only admins can access this endpoint"
// @Failure 404 {object} model.ErrorResponse "User not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "admin",
		"action", "get_user_by_id",
	)
	logger.Info("processing get user by ID request")

	userID := r.PathValue("id")
	if userID == "" {
		util.SendErrorResponse(
			w,
			errs.NewValidationError("User ID is required", map[string]any{}),
		)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		logger.Error("failed to get user by ID",
			"error", err.Error(),
			"user_id", userID,
		)
		util.SendErrorResponse(
			w,
			err,
		)
		return
	}

	logger.Info("successfully retrieved user by ID",
		"user_id", user.ID,
	)
	util.SendJSONResponse(w, http.StatusOK, user)
}

// @Summary Update User by ID
// @Description Updates a user's details by their ID. Only accessible by admins.
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param updateUserRequest body model.UpdateUserRequest true "Updated user details"
// @Success 200 {object} model.User "Updated user details"
// @Failure 403 {object} model.ErrorResponse "Forbidden - Only admins can access this endpoint"
// @Failure 404 {object} model.ErrorResponse "User not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /admin/users/{id} [patch]
func (h *AdminHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "admin",
		"action", "update_user_by_id",
	)
	logger.Info("processing update user by ID request")

	userID := r.PathValue("id")
	if userID == "" {
		util.SendErrorResponse(
			w,
			errs.NewValidationError("User ID is required", map[string]any{}),
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
			errs.NewValidationError("Invalid request parameters", map[string]any{}),
		)
		return
	}

	user, err := h.service.UpdateUserByID(r.Context(), userID, req.Username, req.Email)
	if err != nil {
		logger.Error("failed to update user by ID",
			"error", err.Error(),
			"user_id", userID,
		)
		util.SendErrorResponse(
			w,
			err,
		)
		return
	}

	logger.Info("successfully updated user by ID",
		"user_id", user.ID,
	)
	util.SendJSONResponse(w, http.StatusOK, user)
}

// @Summary Delete User by ID
// @Description Deletes a user by their ID. Only accessible by admins.
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 403 {object} model.ErrorResponse "Forbidden - Only admins can access this endpoint"
// @Failure 404 {object} model.ErrorResponse "User not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "admin",
		"action", "delete_user_by_id",
	)
	logger.Info("processing delete user by ID request")

	userID := r.PathValue("id")
	if userID == "" {
		util.SendErrorResponse(
			w,
			errs.NewValidationError("User ID is required", map[string]any{}),
		)
		return
	}

	if err := h.service.DeleteUserByID(r.Context(), userID); err != nil {
		logger.Error("failed to delete user by ID",
			"error", err.Error(),
			"user_id", userID,
		)
		util.SendErrorResponse(
			w,
			err,
		)
		return
	}

	logger.Info("successfully deleted user by ID",
		"user_id", userID,
	)
	w.WriteHeader(http.StatusNoContent)
}
