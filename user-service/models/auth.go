package models

import (
	"time"
)

type AuthResponse struct {
	Token     string    `json:"token"`
	User      *User     `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ValidateResponse represents the token validation response from the service layer
type ValidateResponse struct {
	IsValid bool   `json:"is_valid"`
	UserID  string `json:"user_id,omitempty"` // omitempty in case validation fails
}

// PublicAuthResponse returns an auth response with public user info only
func (a *AuthResponse) PublicAuthResponse() *AuthResponse {
	return &AuthResponse{
		Token:     a.Token,
		User:      a.User.PublicUser(), // Use the User's PublicUser method
		ExpiresAt: a.ExpiresAt,
	}
}
