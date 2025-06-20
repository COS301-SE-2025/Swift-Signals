package user

import (
	"context"

	userpb "github.com/COS301-SE-2025/Swift-Signals/user-service/proto"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.UserResponse, error) {
	return &userpb.UserResponse{}, nil
}

func (h *Handler) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.AuthResponse, error) {
	return &userpb.AuthResponse{}, nil
}

func (h *Handler) LogoutUser(ctx context.Context, req *userpb.LogoutUserRequest) (*userpb.LogoutResponse, error) {
	return &userpb.LogoutResponse{}, nil
}
