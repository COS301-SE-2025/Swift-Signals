package service

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
)

type AuthService struct {
	userClient *client.UserClient
}

func NewAuthService(uc *client.UserClient) *AuthService {
	return &AuthService{
		userClient: uc,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, req model.RegisterRequest) (resp model.AuthResponse, err error) {
	return model.AuthResponse{}, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req model.LoginRequest) (resp model.AuthResponse, err error) {
	return model.AuthResponse{}, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, token string) (resp model.LogoutResponse, err error) {
	return model.LogoutResponse{}, nil
}
