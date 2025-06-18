package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthHandler struct{}

func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// Validate credentials and generate JWT
	fmt.Println("Generating JWT")
	token := "mocked-jwt-token"
	// -------------------------------------

	resp := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// Create the user
	fmt.Printf("Creating User %#v\n", req)
	// ---------------
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	// Clear token from db
	fmt.Println("Deleting JWT")
	// -------------------
	w.WriteHeader(http.StatusNoContent)
}
