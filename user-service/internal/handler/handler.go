package handler

import (
	"context"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/service"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	service service.UserService
}

func NewUserHandler(s service.UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.UserResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing register user request")

	user, err := h.service.RegisterUser(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		logger.Error("registration failed",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("registration successful",
		"user_id", user.ID,
	)
	return &userpb.UserResponse{
		Id:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		IsAdmin:         user.IsAdmin,
		IntersectionIds: user.IntersectionIDs,
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *Handler) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.LoginUserResponse, error) {
	token, expiryTime, err := h.service.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userpb.LoginUserResponse{
		Token:     token,
		ExpiresAt: timestamppb.New(expiryTime),
	}, nil
}

func (h *Handler) LogoutUser(ctx context.Context, req *userpb.UserIDRequest) (*emptypb.Empty, error) {
	err := h.service.LogoutUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) GetUserByID(ctx context.Context, req *userpb.UserIDRequest) (*userpb.UserResponse, error) {
	user, err := h.service.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &userpb.UserResponse{
		Id:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		IsAdmin:         user.IsAdmin,
		IntersectionIds: user.IntersectionIDs,
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *Handler) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.UserResponse, error) {
	user, err := h.service.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &userpb.UserResponse{
		Id:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		IsAdmin:         user.IsAdmin,
		IntersectionIds: user.IntersectionIDs,
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *Handler) GetAllUsers(req *userpb.GetAllUsersRequest, stream userpb.UserService_GetAllUsersServer) error {
	ctx := stream.Context()
	users, err := h.service.GetAllUsers(ctx, req.GetPage(), req.GetPageSize(), req.GetFilter())
	if err != nil {
		return err
	}

	for _, user := range users {
		userResponse := &userpb.UserResponse{
			Id:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			IsAdmin:         user.IsAdmin,
			IntersectionIds: user.IntersectionIDs,
			CreatedAt:       timestamppb.New(user.CreatedAt),
			UpdatedAt:       timestamppb.New(user.UpdatedAt),
		}
		if err := stream.Send(userResponse); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	user, err := h.service.UpdateUser(ctx, req.GetUserId(), req.GetName(), req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &userpb.UserResponse{
		Id:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		IsAdmin:         user.IsAdmin,
		IntersectionIds: user.IntersectionIDs,
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *userpb.UserIDRequest) (*emptypb.Empty, error) {
	err := h.service.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) GetUserIntersectionIDs(req *userpb.UserIDRequest, stream userpb.UserService_GetUserIntersectionIDsServer) error {
	ctx := stream.Context()
	intersectionIDs, err := h.service.GetUserIntersectionIDs(ctx, req.GetUserId())
	if err != nil {
		return err
	}

	for _, intersectionID := range intersectionIDs {
		response := &userpb.IntersectionIDResponse{
			IntersectionId: intersectionID,
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) AddIntersectionID(ctx context.Context, req *userpb.AddIntersectionIDRequest) (*emptypb.Empty, error) {
	err := h.service.AddIntersectionID(ctx, req.GetUserId(), req.GetIntersectionId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) RemoveIntersectionIDs(ctx context.Context, req *userpb.RemoveIntersectionIDRequest) (*emptypb.Empty, error) {
	err := h.service.RemoveIntersectionIDs(ctx, req.GetUserId(), req.GetIntersectionId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) ChangePassword(ctx context.Context, req *userpb.ChangePasswordRequest) (*emptypb.Empty, error) {
	err := h.service.ChangePassword(ctx, req.GetUserId(), req.GetCurrentPassword(), req.GetNewPassword())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) ResetPassword(ctx context.Context, req *userpb.ResetPasswordRequest) (*emptypb.Empty, error) {
	err := h.service.ResetPassword(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) MakeAdmin(ctx context.Context, req *userpb.AdminRequest) (*emptypb.Empty, error) {
	err := h.service.MakeAdmin(ctx, req.GetUserId(), req.GetAdminUserId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) RemoveAdmin(ctx context.Context, req *userpb.AdminRequest) (*emptypb.Empty, error) {
	err := h.service.RemoveAdmin(ctx, req.GetUserId(), req.GetAdminUserId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
