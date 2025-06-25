package service

import (
	"context"
	"errors"

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

func (s *AuthService) RegisterUser(ctx context.Context, req model.RegisterRequest) (model.RegisterResponse, error) {
	registerResp, err := s.userClient.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return model.RegisterResponse{}, errors.New("unable to register user")
	}
	resp := model.RegisterResponse{
		UserID: registerResp.Id,
	}

	return resp, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	loginResp, err := s.userClient.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		return model.LoginResponse{}, errors.New("unable to login user")
	}
	resp := model.LoginResponse{
		Message: "Login Successful",
		Token:   loginResp.Token,
	}
	return resp, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, token string) (model.LogoutResponse, error) {
	_, err := s.userClient.LogoutUser(ctx, token)
	if err != nil {
		return model.LogoutResponse{}, errors.New("unable to logout user")
	}
	resp := model.LogoutResponse{
		Message: "Logout Successful",
	}
	return resp, nil
}
