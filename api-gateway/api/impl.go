package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/client"
)

type AuthHandler struct {
	UserClient client.UserClientInterface
}

func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.UserClient.LoginUser(ctx, string(req.Email), req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("login failed %v", err), http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{Token: res.Token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.UserClient.LoginUser(ctx, string(req.Email), req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("registration failed %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	// TODO: Get userID from JWT
	userID := "44"

	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.UserClient.LogoutUser(ctx, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("logout failed %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
