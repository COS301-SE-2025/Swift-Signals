package db

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	AddIntersectionID(ctx context.Context, userID string, intID string) error
	GetIntersectionsByUserID(ctx context.Context, userID string) ([]string, error)
}
