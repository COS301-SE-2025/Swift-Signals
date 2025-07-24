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
