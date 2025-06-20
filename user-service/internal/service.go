package user

import (
	"context"
	"errors"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/models"
	"github.com/google/uuid"
)

type Service struct {
	repo db.UserRepository
}

func NewService(r db.UserRepository) *Service {
	return &Service{repo: r}
}

// Mocked business logic

func (s *Service) RegisterUser(ctx context.Context, name, email, password string) (*models.User, error) {
	// TODO: Implement RegisterUser
	id := uuid.New().String()
	user := &models.User{ID: id, Name: name, Email: email, Password: password}
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (string, error) {
	// TODO: Implement RegisterUser
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user.Password != password {
		return "", errors.New("invalid credentials")
	}
	return "mock-token-for-" + user.ID, nil
}

func (s *Service) LogoutUser(ctx context.Context, userID string) bool {
	// TODO: Implement RegisterUser
	return true
}
