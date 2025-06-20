package api

import (
	"context"

	userpb "github.com/COS301-SE-2025/Swift-Signals/api-gateway/protos/user"
)

type MockUserClient struct {
	LoginUserFunc    func(ctx context.Context, email, password string) (*userpb.AuthResponse, error)
	RegisterUserFunc func(ctx context.Context, name, email, password string) (*userpb.UserResponse, error)
	LogoutUserFunc   func(ctx context.Context, userID string) (*userpb.LogoutResponse, error)
}

func (m *MockUserClient) LoginUser(ctx context.Context, email, password string) (*userpb.AuthResponse, error) {
	return m.LoginUserFunc(ctx, email, password)
}

func (m *MockUserClient) RegisterUser(ctx context.Context, name, email, password string) (*userpb.UserResponse, error) {
	return m.RegisterUserFunc(ctx, name, email, password)
}

func (m *MockUserClient) LogoutUser(ctx context.Context, userID string) (*userpb.LogoutResponse, error) {
	return m.LogoutUserFunc(ctx, userID)
}
