package profile

import (
	"context"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestDeleteProfile_Success() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, nil)

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().NoError(err)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_MissingUserID() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	// Note: Not setting userID in context to simulate missing userID

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("user ID missing inside of handler", svcError.Message)
}

func (suite *TestSuite) TestDeleteProfile_UserNotFound() {
	userID := "nonexistent-user"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{"userID": userID}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_InternalError() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_UnauthorizedError() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewUnauthorizedError("invalid token", map[string]any{}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("invalid token", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_ForbiddenError() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewForbiddenError("cannot delete account with active intersections", map[string]any{"userID": userID}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("cannot delete account with active intersections", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_ConflictError() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewConflictError("user has pending operations", map[string]any{"userID": userID}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrConflict, svcError.Code)
	suite.Equal("user has pending operations", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_ServiceUnavailable() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewUnavailableError("user service is temporarily unavailable", map[string]any{}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)
	suite.Equal("user service is temporarily unavailable", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_ValidationError() {
	userID := ""

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewValidationError("user ID cannot be empty", map[string]any{"userID": userID}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("user ID cannot be empty", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_AdminUserDeletion() {
	userID := "admin-456"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, nil)

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().NoError(err)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_UserWithIntersections() {
	userID := "user-with-intersections"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewForbiddenError("cannot delete user with assigned intersections", map[string]any{"userID": userID}))

	err := suite.service.DeleteProfile(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("cannot delete user with assigned intersections", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}
