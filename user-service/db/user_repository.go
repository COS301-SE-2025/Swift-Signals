package db

import (
	"context"
	"errors"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/models"
)

// Database Interface
type UserRepository interface {
	CreateUser(context.Context, *models.User) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
}

// Mock implementation
type inMemoryUserRepo struct {
	users map[string]*models.User
}

func NewUserRepository() UserRepository {
	return &inMemoryUserRepo{users: make(map[string]*models.User)}
}

func (r *inMemoryUserRepo) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	r.users[u.Email] = u
	return u, nil
}

func (r *inMemoryUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
