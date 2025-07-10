package db

import (
	"context"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
)

// Repository interface defines the contract for user data operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	AddIntersectionID(ctx context.Context, userID int, intID int) error
	GetIntersectionsByUserID(ctx context.Context, userID int) ([]int, error)
}
