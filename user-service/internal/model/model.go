package model

import (
	"time"
)

type User struct {
	ID              string    `json:"id"               db:"id"`
	Name            string    `json:"name"             db:"name"`
	Email           string    `json:"email"            db:"email"`
	Password        string    `json:"-"                db:"password"`
	IsAdmin         bool      `json:"is_admin"         db:"is_admin"`
	IntersectionIDs []string  `json:"intersection_ids" db:"intersection_ids"`
	CreatedAt       time.Time `json:"created_at"       db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"       db:"updated_at"`
}

func (u *User) PublicUser() *User {
	return &User{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		IsAdmin:         u.IsAdmin,
		IntersectionIDs: u.IntersectionIDs,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
