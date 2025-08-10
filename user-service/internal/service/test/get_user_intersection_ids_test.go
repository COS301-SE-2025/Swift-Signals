package test

import (
	"context"
	"errors"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestGetUserIntersectionIDs_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	expectedIntersectionIDs := []string{"intersection-1", "intersection-2", "intersection-3"}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(expectedIntersectionIDs, nil)

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedIntersectionIDs, result)
	suite.Len(result, 3)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_EmptyIntersections() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	expectedIntersectionIDs := []string{}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(expectedIntersectionIDs, nil)

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedIntersectionIDs, result)
	suite.Empty(result)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_GetUserByIDError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_GetUserByIDServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_GetIntersectionsError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	repoError := errors.New("failed to fetch intersections")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to fetch user's intersection IDs", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_GetIntersectionsServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	serviceError := errs.NewInternalError(
		"failed to fetch intersections",
		errors.New("db error"),
		nil,
	)

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to fetch intersections", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_EmptyUserID() {
	userID := ""

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_InvalidUUIDFormat() {
	userID := "invalid-uuid-format"

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_SingleIntersection() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	expectedIntersectionIDs := []string{"intersection-1"}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(expectedIntersectionIDs, nil)

	ctx := context.Background()

	result, err := suite.service.GetUserIntersectionIDs(ctx, userID)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedIntersectionIDs, result)
	suite.Len(result, 1)
	suite.Equal("intersection-1", result[0])

	suite.repo.AssertExpectations(suite.T())
}

func TestServiceGetUserIntersectionIDs(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
