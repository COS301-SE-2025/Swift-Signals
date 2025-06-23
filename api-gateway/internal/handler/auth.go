package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// @Summary User Login
// @Description Authenticates a user and returns an authentication token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginRequest body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.AuthResponse "Successful login"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload or credentials"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// TODO: Implement Service Logic
	log.Println("Attempting to log in user")

	resp := &model.AuthResponse{
		Message: "Login Example Successful",
		Token:   "example-token",
	}
	util.SendJSONResponse(w, http.StatusOK, resp)
}
