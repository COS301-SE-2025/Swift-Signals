package service

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
)

type ProfileService struct {
	userClient client.UserClientInterface
}

func NewProfileService(uc client.UserClientInterface) ProfileServiceInterface {
	return &ProfileService{
		userClient: uc,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (model.User, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "profile",
	)

	logger.Debug("calling user client to get user profile")
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

func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	userID string,
	req model.UpdateUserRequest,
) (model.User, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "profile",
	)

	logger.Debug("calling user client to update user profile")
	rpcUser, err := s.userClient.UpdateUser(ctx, userID, req.Username, req.Email)
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

func (s *ProfileService) DeleteProfile(ctx context.Context, userID string) error {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "profile",
	)

	logger.Debug("calling user client to delete user profile")
	_, err := s.userClient.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// ProfileServiceInterface creates stub for testing
type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, userID string) (model.User, error)
	UpdateProfile(
		ctx context.Context,
		userID string,
		req model.UpdateUserRequest,
	) (model.User, error)
	DeleteProfile(ctx context.Context, userID string) error
}

// NOTE: Asserts Interface Implementation
var _ ProfileServiceInterface = (*ProfileService)(nil)
