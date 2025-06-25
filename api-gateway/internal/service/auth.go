package service

import "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) RegisterUser(req model.RegisterRequest) (resp model.AuthResponse, err error) {
	return model.AuthResponse{}, nil
}
