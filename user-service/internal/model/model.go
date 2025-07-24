package model

import (
	"errors"
	"time"
)

// Common errors used throughout the application
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidUserID   = errors.New("invalid user ID")
	ErrInvalidUserData = errors.New("invalid user data")
)

// User represents a user in the system
type User struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"-" db:"password"` // "-" ensures password is never serialized to JSON
	IsAdmin         bool      `json:"is_admin" db:"is_admin"`
	IntersectionIDs []string  `json:"intersection_ids" db:"intersection_ids"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// PublicUser returns a user struct without sensitive information
func (u *User) PublicUser() *User {
	return &User{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		IsAdmin:         u.IsAdmin,
		IntersectionIDs: u.IntersectionIDs,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		// Password is intentionally omitted
	}
}
