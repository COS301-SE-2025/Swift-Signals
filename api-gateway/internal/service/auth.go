package service

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
)

type AuthService struct {
	userClient client.UserClientInterface
}

func NewAuthService(uc client.UserClientInterface) AuthServiceInterface {
	return &AuthService{
		userClient: uc,
	}
}

func (s *AuthService) RegisterUser(
	ctx context.Context,
	req model.RegisterRequest,
) (model.RegisterResponse, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "auth",
	)

	logger.Debug("calling user client to register user",
		"username", req.Username,
	)
	registerResp, err := s.userClient.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	resp := model.RegisterResponse{
		UserID: registerResp.Id,
	}
	return resp, nil
}

func (s *AuthService) LoginUser(
	ctx context.Context,
	req model.LoginRequest,
) (model.LoginResponse, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "auth",
	)

	logger.Debug("calling user client to login user")
	loginResp, err := s.userClient.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		return model.LoginResponse{}, err
	}

	resp := model.LoginResponse{
		Message: "Login Successful",
		Token:   loginResp.Token,
	}
	return resp, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, token string) (model.LogoutResponse, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "auth",
	)

	logger.Debug("calling user client to logout user",
		"token", token,
	)
	_, err := s.userClient.LogoutUser(ctx, token)
	if err != nil {
		return model.LogoutResponse{}, err
	}

	resp := model.LogoutResponse{
		Message: "Logout Successful",
	}
	return resp, nil
}

// AuthServiceInterface creates stub for testing
type AuthServiceInterface interface {
	RegisterUser(ctx context.Context, req model.RegisterRequest) (model.RegisterResponse, error)
	LoginUser(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	LogoutUser(ctx context.Context, token string) (model.LogoutResponse, error)
}

// Note: Asserts Interface Implementation
var _ AuthServiceInterface = (*AuthService)(nil)
