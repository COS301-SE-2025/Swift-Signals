package handler

import (
	"log"
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
)

type IntersectionHandler struct{}

func NewIntersectionHandler() *IntersectionHandler {
	return &IntersectionHandler{}
}

// @Summary Get All Intersections
// @Description Retrieves all the intersections associated with the user.
// @Tags Intersections
// @Accept json
// @Produce json
// @Success 200 {object} model.Intersections "Successful intersections retrieval"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections [get]
func (h *IntersectionHandler) GetAllIntersections(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		util.SendErrorResponse(w, http.StatusUnauthorized, "Authorization token is missing")
		return
	}

	// TODO: Implement Service Logic
	log.Println("Getting all intersections...")
	resp := model.Intersections{}

	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary Get Intersection by ID
// @Description Retrieves a single intersection by its unique identifier.
// @Tags Intersections
// @Accept json
// @Produce json
// @Param id path string true "Intersection ID"
// @Success 200 {object} model.Intersection "Successful intersection retrieval"
// @Failure 400 {object} model.ErrorResponse "Bad Request: Invalid or missing ID parameter"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 404 {object} model.ErrorResponse "Not Found: Intersection does not exist"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections/{id} [get]
func (h *IntersectionHandler) GetIntersection(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		util.SendErrorResponse(w, http.StatusUnauthorized, "Authorization token is missing")
		return
	}

	util.SendJSONResponse(w, http.StatusOK, model.Intersection{})
}

// @Summary Create Intersection
// @Description Creates a new intersection with the given arguments
// @Tags Intersections
// @Accept json
// @Produce json
// @Param createIntersectionRequest body model.CreateIntersectionRequest true "intersection information"
// @Success 201 {object} model.AuthResponse "User successfully registered"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload or missing fields"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /intersections [post]
func (h *IntersectionHandler) CreateIntersection(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		util.SendErrorResponse(w, http.StatusUnauthorized, "Authorization token is missing")
		return
	}

	util.SendJSONResponse(w, http.StatusOK, model.CreateIntersectionRequest{})
}
