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

func (h *Handler) RegisterUser(
	ctx context.Context,
	req *userpb.RegisterUserRequest,
) (*userpb.UserResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing RegisterUser request")

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

func (h *Handler) LoginUser(
	ctx context.Context,
	req *userpb.LoginUserRequest,
) (*userpb.LoginUserResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing LoginUser request")
	token, expiryTime, err := h.service.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		logger.Error("login failed",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("LoginUser successful",
		"token", token,
	)
	return &userpb.LoginUserResponse{
		Token:     token,
		ExpiresAt: timestamppb.New(expiryTime),
	}, nil
}

func (h *Handler) LogoutUser(
	ctx context.Context,
	req *userpb.UserIDRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing LogoutUser request")

	err := h.service.LogoutUser(ctx, req.GetUserId())
	if err != nil {
		logger.Error("logout failed",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("LogoutUser successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) GetUserByID(
	ctx context.Context,
	req *userpb.UserIDRequest,
) (*userpb.UserResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing GetUserByID request")

	user, err := h.service.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		logger.Error("failed to find user",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("GetUserByID successful",
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

func (h *Handler) GetUserByEmail(
	ctx context.Context,
	req *userpb.GetUserByEmailRequest,
) (*userpb.UserResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing GetUserByEmail request")

	user, err := h.service.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		logger.Error("failed to find user",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("GetUserByEmail successful",
		"email", user.Email,
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

func (h *Handler) GetAllUsers(
	req *userpb.GetAllUsersRequest,
	stream userpb.UserService_GetAllUsersServer,
) error {
	ctx := stream.Context()

	logger := util.LoggerFromContext(ctx)
	logger.Info("processing GetAllUsers request")

	users, err := h.service.GetAllUsers(ctx, req.GetPage(), req.GetPageSize(), req.GetFilter())
	if err != nil {
		logger.Error("failed to get all users",
			"error", err.Error(),
		)
		return errs.HandleServiceError(err)
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
			logger.Error("failed to send user",
				"error", err.Error(),
			)
			return errs.HandleServiceError(err)
		}
	}

	logger.Info("GetAllUsers successful")
	return nil
}

func (h *Handler) UpdateUser(
	ctx context.Context,
	req *userpb.UpdateUserRequest,
) (*userpb.UserResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing UpdateUser request")

	user, err := h.service.UpdateUser(ctx, req.GetUserId(), req.GetName(), req.GetEmail())
	if err != nil {
		logger.Error("failed to update user",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("UpdateUser successful",
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

func (h *Handler) DeleteUser(
	ctx context.Context,
	req *userpb.UserIDRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing DeleteUser request")

	err := h.service.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		logger.Error("failed to delete user",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("DeleteUser successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) GetUserIntersectionIDs(
	req *userpb.UserIDRequest,
	stream userpb.UserService_GetUserIntersectionIDsServer,
) error {
	ctx := stream.Context()

	logger := util.LoggerFromContext(ctx)
	logger.Info("processing GetUserIntersectionIDs request")

	intersectionIDs, err := h.service.GetUserIntersectionIDs(ctx, req.GetUserId())
	if err != nil {
		logger.Error("failed to get all user intersection ids",
			"error", err.Error(),
		)
		return errs.HandleServiceError(err)
	}

	for _, intersectionID := range intersectionIDs {
		response := &userpb.IntersectionIDResponse{
			IntersectionId: intersectionID,
		}
		if err := stream.Send(response); err != nil {
			logger.Error("failed to send intersection id",
				"error", err.Error(),
			)
			return errs.HandleServiceError(err)
		}
	}

	logger.Info("GetUserIntersectionIDs successful")
	return nil
}

func (h *Handler) AddIntersectionID(
	ctx context.Context,
	req *userpb.AddIntersectionIDRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing AddIntersectionID request")

	err := h.service.AddIntersectionID(ctx, req.GetUserId(), req.GetIntersectionId())
	if err != nil {
		logger.Error("failed to add intersection id",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("AddIntersectionID successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) RemoveIntersectionIDs(
	ctx context.Context,
	req *userpb.RemoveIntersectionIDRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing RemoveIntersectionIDs request")

	err := h.service.RemoveIntersectionIDs(ctx, req.GetUserId(), req.GetIntersectionId())
	if err != nil {
		logger.Error("failed to remove intersection ids",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("RemoveIntersectionIDs successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) ChangePassword(
	ctx context.Context,
	req *userpb.ChangePasswordRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing ChangePassword request")

	err := h.service.ChangePassword(
		ctx,
		req.GetUserId(),
		req.GetCurrentPassword(),
		req.GetNewPassword(),
	)
	if err != nil {
		logger.Error("failed to change password",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("ChangePassword successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) ResetPassword(
	ctx context.Context,
	req *userpb.ResetPasswordRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing ResetPassword request")

	err := h.service.ResetPassword(ctx, req.GetEmail())
	if err != nil {
		logger.Error("failed to reset password",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("ResetPassword successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) MakeAdmin(ctx context.Context, req *userpb.AdminRequest) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing MakeAdmin request")

	err := h.service.MakeAdmin(ctx, req.GetUserId(), req.GetAdminUserId())
	if err != nil {
		logger.Error("failed to make user admin",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("MakeAdmin successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) RemoveAdmin(
	ctx context.Context,
	req *userpb.AdminRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing RemoveAdmin request")

	err := h.service.RemoveAdmin(ctx, req.GetUserId(), req.GetAdminUserId())
	if err != nil {
		logger.Error("failed to remove user from admin",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("RemoveAdmin successful")
	return &emptypb.Empty{}, nil
}
