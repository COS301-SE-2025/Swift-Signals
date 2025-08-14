package test

import (
	"context"
	"errors"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestDeleteUser_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	existingUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("DeleteUser", mock.Anything, userID).Return(nil)

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_InvalidUserID() {
	userID := "invalid-uuid"

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

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

func (suite *TestSuite) TestDeleteUser_EmptyUserID() {
	userID := ""

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

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

func (suite *TestSuite) TestDeleteUser_UserNotFound() {
	userID := "550e8400-e29b-41d4-a716-446655440000"

	suite.repo.On("GetUserByID", mock.Anything, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{}))

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_RepositoryGetUserError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_RepositoryDeleteError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	existingUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	repoError := errors.New("database delete failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("DeleteUser", mock.Anything, userID).Return(repoError)

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to delete user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_ServiceErrorPropagationGetUser() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_ServiceErrorPropagationDelete() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	existingUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	serviceError := errs.NewForbiddenError("cannot delete user", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("DeleteUser", mock.Anything, userID).Return(serviceError)

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("cannot delete user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_AdminUser() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	existingUser := &model.User{
		ID:      userID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("DeleteUser", mock.Anything, userID).Return(nil)

	ctx := context.Background()

	err := suite.service.DeleteUser(ctx, userID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}
