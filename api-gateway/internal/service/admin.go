package service

import (
	"context"
	"io"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

type AdminService struct {
	userClient client.UserClientInterface
}

func NewAdminService(uc client.UserClientInterface) AdminServiceInterface {
	return &AdminService{
		userClient: uc,
	}
}

func (s *AdminService) GetAllUsers(ctx context.Context, pgNo, pgSz int) ([]model.User, error) {
	role, ok := middleware.GetRole(ctx)
	if !ok || role != "admin" {
		return nil, errs.NewForbiddenError(
			"only admins can access this endpoint",
			map[string]any{"role": role},
		)
	}

	stream, err := s.userClient.GetAllUsers(ctx, int32(pgNo), int32(pgSz), "")
	if err != nil {
		return nil, err
	}

	users := []model.User{}
	for {
		rpcUser, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errs.NewInternalError(
				"unable to get all users",
				err,
				map[string]any{},
			)
		}
		users = append(users, model.User{
			ID:              rpcUser.Id,
			Username:        rpcUser.Name,
			Email:           rpcUser.Email,
			IsAdmin:         rpcUser.IsAdmin,
			IntersectionIDs: rpcUser.IntersectionIds,
		})
	}

	return users, nil
}

func (s *AdminService) GetUserByID(ctx context.Context, userID string) (model.User, error) {
	role, ok := middleware.GetRole(ctx)
	if !ok || role != "admin" {
		return model.User{}, errs.NewForbiddenError(
			"only admins can access this endpoint",
			map[string]any{"role": role},
		)
	}

	rpcUser, err := s.userClient.GetUserByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:              rpcUser.Id,
		Username:        rpcUser.Name,
		Email:           rpcUser.Email,
		IsAdmin:         rpcUser.IsAdmin,
		IntersectionIDs: rpcUser.IntersectionIds,
	}, nil
}

func (s *AdminService) UpdateUserByID(
	ctx context.Context,
	userID, name, email string,
) (model.User, error) {
	role, ok := middleware.GetRole(ctx)
	if !ok || role != "admin" {
		return model.User{}, errs.NewForbiddenError(
			"only admins can access this endpoint",
			map[string]any{"role": role},
		)
	}

	rpcUser, err := s.userClient.UpdateUser(ctx, userID, name, email)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:              rpcUser.Id,
		Username:        rpcUser.Name,
		Email:           rpcUser.Email,
		IsAdmin:         rpcUser.IsAdmin,
		IntersectionIDs: rpcUser.IntersectionIds,
	}, nil
}

func (s *AdminService) DeleteUserByID(ctx context.Context, userID string) error {
	role, ok := middleware.GetRole(ctx)
	if !ok || role != "admin" {
		return errs.NewForbiddenError(
			"only admins can access this endpoint",
			map[string]any{"role": role},
		)
	}

	_, err := s.userClient.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// AuthServiceInterface creates stub for testing
type AdminServiceInterface interface {
	GetAllUsers(ctx context.Context, pgNo, pgSz int) ([]model.User, error)
	GetUserByID(ctx context.Context, userID string) (model.User, error)
	UpdateUserByID(ctx context.Context, userID, name, email string) (model.User, error)
	DeleteUserByID(ctx context.Context, userID string) error
}

// NOTE: Asserts Interface Implementation
var _ AdminServiceInterface = (*AdminService)(nil)
