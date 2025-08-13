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

type IntersectionHandler struct {
	service   service.IntersectionServiceInterface
	validator *validator.Validate
}

func NewIntersectionHandler(s service.IntersectionServiceInterface) *IntersectionHandler {
	return &IntersectionHandler{
		service:   s,
		validator: validator.New(),
	}
}

// @Summary Get All Intersections
// @Description Retrieves all the intersections associated with the user.
// @Tags Intersections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Intersections "Successful intersections retrieval"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections [get]
func (h *IntersectionHandler) GetAllIntersections(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "intersection",
		"action", "getAllIntersections",
	)
	logger.Info("processing getAllIntersections request")

	userID, ok := middleware.GetUserID(r)
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

	resp, err := h.service.GetAllIntersections(r.Context(), userID)
	if err != nil {
		logger.Error("request failed",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("request successful")
	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary Get Intersection by ID
// @Description Retrieves a single intersection by its unique identifier.
// @Tags Intersections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Intersection ID"
// @Success 200 {object} model.Intersection "Successful intersection retrieval"
// @Failure 400 {object} model.ErrorResponse "Bad Request: Invalid or missing ID parameter"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 404 {object} model.ErrorResponse "Not Found: Intersection does not exist"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections/{id} [get]
func (h *IntersectionHandler) GetIntersection(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "intersection",
		"action", "getIntersection",
	)
	logger.Info("processing getIntersection request")

	userID, ok := middleware.GetUserID(r)
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

	intersectionID := r.PathValue("id")

	resp, err := h.service.GetIntersectionByID(r.Context(), userID, intersectionID)
	if err != nil {
		logger.Error("request failed",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("request successful",
		"intersection_id", resp.ID,
	)
	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary Create Intersection
// @Description Creates a new intersection with the given arguments
// @Tags Intersections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param createIntersectionRequest body model.CreateIntersectionRequest true "Intersection information"
// @Success 201 {object} model.CreateIntersectionResponse "Intersection successfully created"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload or missing fields"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections [post]
func (h *IntersectionHandler) CreateIntersection(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "intersection",
		"action", "createIntersection",
	)
	logger.Info("processing createIntersection request")

	var req model.CreateIntersectionRequest
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
			errs.NewValidationError("name, density, and parameters are required", map[string]any{}),
		)
		return
	}

	userID, ok := middleware.GetUserID(r)
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

	resp, err := h.service.CreateIntersection(r.Context(), userID, req)
	if err != nil {
		logger.Error("request failed",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("request successful",
		"intersection_id", resp.Id,
	)
	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary Update Intersection
// @Description Partially updates fields of an existing intersection by ID.
// @Tags Intersections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Intersection ID"
// @Param body body model.UpdateIntersectionRequest true "Fields to update"
// @Success 200 {object} model.Intersection "Successful update"
// @Failure 400 {object} model.ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 404 {object} model.ErrorResponse "Not Found: Intersection does not exist"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections/{id} [patch]
func (h *IntersectionHandler) UpdateIntersection(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "intersection",
		"action", "updateIntersection",
	)
	logger.Info("processing updateIntersection request")

	var req model.UpdateIntersectionRequest
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
			errs.NewValidationError("validation failed", map[string]any{}),
		)
		return
	}

	userID, ok := middleware.GetUserID(r)
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

	intersectionID := r.PathValue("id")

	resp, err := h.service.UpdateIntersectionByID(r.Context(), userID, intersectionID, req)
	if err != nil {
		logger.Error("request failed",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("request successful",
		"intersection_id", resp.ID,
	)
	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary Delete Intersection
// @Description Deletes the intersection with the given ID.
// @Tags Intersections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Intersection ID"
// @Success 204 "No Content"
// @Failure 400 {object} model.ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 404 {object} model.ErrorResponse "Not Found: Intersection does not exist"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections/{id} [delete]
func (h *IntersectionHandler) DeleteIntersection(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "intersection",
		"action", "deleteIntersection",
	)
	logger.Info("processing deleteIntersection request")

	userID, ok := middleware.GetUserID(r)
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

	intersectionID := r.PathValue("id")

	err := h.service.DeleteIntersectionByID(r.Context(), userID, intersectionID)
	if err != nil {
		logger.Error("request failed",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("request successful",
		"intersection_id", intersectionID,
	)
	w.WriteHeader(http.StatusNoContent)
}
