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

// @Summary User Registration
// @Description Registers a new user and returns an authentication token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param registerRequest body model.RegisterRequest true "User registration details"
// @Success 201 {object} model.AuthResponse "User successfully registered"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload or missing fields"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		util.SendErrorResponse(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// TODO: Implement Service Logic
	log.Println("Attempting to register user...")

	resp := &model.AuthResponse{
		Message: "Registration Example Successful",
		Token:   "example.token",
	}
	util.SendJSONResponse(w, http.StatusCreated, resp)
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
	log.Println("Attempting to log in user...")

	resp := &model.AuthResponse{
		Message: "Login Example Successful",
		Token:   "example.token",
	}
	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary User Logout
// @Description Invalidates the user's session or token on the server-side.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.LogoutResponse "Successful logout"
// @Failure 401 {object} model.ErrorResponse "Unauthorized: Token missing or invalid"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		util.SendErrorResponse(w, http.StatusUnauthorized, "Authorization token is missing")
		return
	}

	// TODO: Implement Service Logic
	log.Println("Attempting to logout user...")

	resp := &model.LogoutResponse{
		Message: "Logout Successful",
	}
	util.SendJSONResponse(w, http.StatusOK, resp)
}

// @Summary Reset Password
// @Description Reset's a user's password in case they forgot it.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param resetPasswordRequest body model.ResetPasswordRequest true "User Email"
// @Success 200 {object} model.ResetPasswordResponse "Successful password reset"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload or email"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req model.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// TODO: Implement Service Logic

	resp := &model.ResetPasswordResponse{
		Message: "Reset Example Successful",
	}
	util.SendJSONResponse(w, http.StatusOK, resp)
}
