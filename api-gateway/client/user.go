package client

import (
	"context"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/api-gateway/protos/user"
	"google.golang.org/grpc"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		client: userpb.NewUserServiceClient(conn),
	}
}

func (uc *UserClient) RegisterUser(ctx context.Context, name, email, password string) (*userpb.UserResponse, error) {
	req := &userpb.RegisterUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.RegisterUser(ctx, req)
}

func (uc *UserClient) LoginUser(ctx context.Context, email, password string) (*userpb.AuthResponse, error) {
	req := &userpb.LoginUserRequest{
		Email:    email,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.LoginUser(ctx, req)
}

func (uc *UserClient) LogoutUser(ctx context.Context, userID string) (*userpb.LogoutResponse, error) {
	req := &userpb.LogoutUserRequest{
		UserId: userID,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.LogoutUser(ctx, req)
}

// NOTE: Creates stub for testing
type UserClientInterface interface {
	RegisterUser(ctx context.Context, name, email, password string) (*userpb.UserResponse, error)
	LoginUser(ctx context.Context, email, password string) (*userpb.AuthResponse, error)
	LogoutUser(ctx context.Context, userID string) (*userpb.LogoutResponse, error)
}

// NOTE: Asserts Interface Implementation
var _ UserClientInterface = (*UserClient)(nil)
