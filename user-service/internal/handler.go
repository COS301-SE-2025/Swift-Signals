package user

import (
	"context"

	userpb "github.com/COS301-SE-2025/Swift-Signals/user-service/proto"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.UserResponse, error) {
	user, err := h.service.RegisterUser(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userpb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (h *Handler) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.AuthResponse, error) {
	token, err := h.service.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userpb.AuthResponse{Token: token}, nil
}

func (h *Handler) LogoutUser(ctx context.Context, req *userpb.LogoutUserRequest) (*userpb.LogoutResponse, error) {
	success := h.service.LogoutUser(ctx, req.GetUserId())
	return &userpb.LogoutResponse{Success: success}, nil
}
