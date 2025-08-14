package test

import (
	"context"
	"errors"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestLogoutUser_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)

	ctx := context.Background()

	err := suite.service.LogoutUser(ctx, userID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_InvalidUserID() {
	userID := "invalid-uuid"

	ctx := context.Background()

	err := suite.service.LogoutUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID must be a valid UUID",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_EmptyUserID() {
	userID := ""

	ctx := context.Background()

	err := suite.service.LogoutUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_UserNotFound() {
	userID := "550e8400-e29b-41d4-a716-446655440000"

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, nil)

	ctx := context.Background()

	err := suite.service.LogoutUser(ctx, userID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_RepositoryError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.LogoutUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_ServiceErrorPropagation() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.LogoutUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}
