package handler

import (
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
)

type SimulationHandler struct {
	service service.SimulationServiceInterface
}

func NewSimulationHandler(s service.SimulationServiceInterface) *SimulationHandler {
	return &SimulationHandler{
		service: s,
	}
}

// @Summary Get Simulation Data
// @Description Generates and returns simulation data for a specific intersection.
// @Tags Simulation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Intersection ID"
// @Success 200 {object} model.SimulationResponse "Successful simulation data retrieval"
// @Failure 400 {object} model.ErrorResponse "Bad Request: Invalid input parameters"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 404 {object} model.ErrorResponse "Not Found: Intersection not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections/{id}/simulate [get]
func (h *SimulationHandler) GetSimulation(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "simulation",
		"action", "getSimulation",
	)
	logger.Info("processing getSimulation request")

	intersectionID := r.PathValue("id")

	resp, err := h.service.GetSimulationData(r.Context(), intersectionID)
	if err != nil {
		logger.Error("request failed",
			"error", err.Error(),
		)
		util.SendErrorResponse(w, err)
		return
	}

	logger.Info("request successful")
	logger.Debug("simulation response",
		"response", resp,
	)
	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary Get Optimised Simulation Data
// @Description Generates and returns optimised simulation data for a specific intersection.
// @Tags Simulation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Intersection ID"
// @Success 200 {object} model.SimulationResponse "Successful optimised simulation data retrieval"
// @Failure 400 {object} model.ErrorResponse "Bad Request: Invalid input parameters"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 404 {object} model.ErrorResponse "Not Found: Intersection not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections/{id}/optimise [get]
func (h *SimulationHandler) GetOptimisedSimulation(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "simulation",
		"action", "getOptimisedSimulation",
	)
	logger.Info("processing getOptimisedSimulation request")

	intersectionID := r.PathValue("id")

	resp, err := h.service.GetOptimisedData(r.Context(), intersectionID)
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

// @Summary Run Optimisation
// @Description Runs optimisation for a specific intersection and returns if there was an improvement.
// @Tags Simulation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Intersection ID"
// @Success 200 {object} model.OptimisationResponse "Successful optimisation run"
// @Failure 400 {object} model.ErrorResponse "Bad Request: Invalid input parameters"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 404 {object} model.ErrorResponse "Not Found: Intersection not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections/{id}/optimise [post]
func (h *SimulationHandler) RunOptimisation(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context()).With(
		"handler", "simulation",
		"action", "runOptimisation",
	)
	logger.Info("processing runOptimisation request")

	intersectionID := r.PathValue("id")

	resp, err := h.service.OptimiseIntersection(r.Context(), intersectionID)
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
