package service

import (
	"context"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
)

type UserService interface {
	RegisterUser(ctx context.Context, name, email, password string) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (string, time.Time, error)
	LogoutUser(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllUsers(ctx context.Context, page, pageSize int32, filter string) ([]*model.User, error)
	UpdateUser(ctx context.Context, userID, name, email string) (*model.User, error)
	DeleteUser(ctx context.Context, userID string) error
	GetUserIntersectionIDs(ctx context.Context, userID string) ([]string, error)
	AddIntersectionID(ctx context.Context, userID string, intersectionID string) error
	RemoveIntersectionIDs(ctx context.Context, userID string, intersectionIDs []string) error
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) error
	MakeAdmin(ctx context.Context, userID, adminUserID string) error
	RemoveAdmin(ctx context.Context, userID, adminUserID string) error
}

type RegisterUserRequest struct {
	Name     string `validate:"required,min=2,max=100" json:"name"`
	Email    string `validate:"required,email,max=255" json:"email"`
	Password string `validate:"required,min=8,max=128" json:"password"`
}

type LoginUserRequest struct {
	Email    string `validate:"required,email,max=255" json:"email"`
	Password string `validate:"required,min=8"         json:"password"`
}

type UpdateUserRequest struct {
	UserID string `validate:"required,uuid4"         json:"user_id"`
	Name   string `validate:"required,min=2,max=100" json:"name"`
	Email  string `validate:"required,email,max=255" json:"email"`
}

type ChangePasswordRequest struct {
	UserID          string `validate:"required,uuid4"         json:"user_id"`
	CurrentPassword string `validate:"required,min=1"         json:"current_password"`
	NewPassword     string `validate:"required,min=8,max=128" json:"new_password"`
}

type GetUserByIDRequest struct {
	UserID string `validate:"required,uuid4" json:"user_id"`
}

type GetUserByEmailRequest struct {
	Email string `validate:"required,email,max=255" json:"email"`
}

type GetAllUsersRequest struct {
	Page     int32  `validate:"min=1"         json:"page"`
	PageSize int32  `validate:"min=1,max=100" json:"page_size"`
	Filter   string `validate:"max=255"       json:"filter"`
}

type DeleteUserRequest struct {
	UserID string `validate:"required,uuid4" json:"user_id"`
}

type LogoutUserRequest struct {
	UserID string `validate:"required,uuid4" json:"user_id"`
}

type GetUserIntersectionIDsRequest struct {
	UserID string `validate:"required,uuid4" json:"user_id"`
}

type AddIntersectionIDRequest struct {
	UserID         string `validate:"required,uuid4"         json:"user_id"`
	IntersectionID string `validate:"required,min=1,max=100" json:"intersection_id"`
}

type RemoveIntersectionIDsRequest struct {
	UserID          string   `validate:"required,uuid4"                             json:"user_id"`
	IntersectionIDs []string `validate:"required,min=1,dive,required,min=1,max=100" json:"intersection_ids"`
}
