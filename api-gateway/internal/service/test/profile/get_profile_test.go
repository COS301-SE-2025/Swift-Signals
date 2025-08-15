package profile

import (
	"context"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestGetProfile_Success() {
	userID := "user-123"

	expectedUser := createTestUser(
		userID,
		"John Doe",
		"john@example.com",
		false,
		[]string{"intersection-1", "intersection-2"},
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(expectedUser, nil)

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal("John Doe", result.Username)
	suite.Equal("john@example.com", result.Email)
	suite.False(result.IsAdmin)
	suite.Equal([]string{"intersection-1", "intersection-2"}, result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_AdminUser() {
	userID := "admin-456"

	expectedUser := createTestUser(
		userID,
		"Jane Admin",
		"jane@admin.com",
		true,
		[]string{"intersection-1", "intersection-2", "intersection-3"},
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(expectedUser, nil)

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal("Jane Admin", result.Username)
	suite.Equal("jane@admin.com", result.Email)
	suite.True(result.IsAdmin)
	suite.Equal(
		[]string{"intersection-1", "intersection-2", "intersection-3"},
		result.IntersectionIDs,
	)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_UserWithNoIntersections() {
	userID := "user-789"

	expectedUser := createTestUser(
		userID,
		"Bob Wilson",
		"bob@example.com",
		false,
		[]string{},
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(expectedUser, nil)

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal("Bob Wilson", result.Username)
	suite.Equal("bob@example.com", result.Email)
	suite.False(result.IsAdmin)
	suite.Empty(result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_MissingUserID() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	// Note: Not setting userID in context to simulate missing userID

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("user ID missing inside of handler", svcError.Message)
}

func (suite *TestSuite) TestGetProfile_UserNotFound() {
	userID := "nonexistent-user"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{"userID": userID}))

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_InternalError() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_UnauthorizedError() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, errs.NewUnauthorizedError("invalid token", map[string]any{}))

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("invalid token", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_ServiceUnavailable() {
	userID := "user-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, errs.NewUnavailableError("user service is temporarily unavailable", map[string]any{}))

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)
	suite.Equal("user service is temporarily unavailable", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_EmptyUserResponse() {
	userID := "user-123"

	expectedUser := createTestUser(
		"",
		"",
		"",
		false,
		nil,
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("GetUserByID", ctx, userID).
		Return(expectedUser, nil)

	result, err := suite.service.GetProfile(ctx, userID)

	suite.Require().NoError(err)
	suite.Equal("", result.ID)
	suite.Equal("", result.Username)
	suite.Equal("", result.Email)
	suite.False(result.IsAdmin)
	suite.Nil(result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}
