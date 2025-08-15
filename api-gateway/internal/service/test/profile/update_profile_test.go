package profile

import (
	"context"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestUpdateProfile_Success() {
	userID := "user-123"
	username := "Updated Name"
	email := "updated@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	expectedUser := createTestUser(
		userID,
		username,
		email,
		false,
		[]string{"intersection-1", "intersection-2"},
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(expectedUser, nil)

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal(username, result.Username)
	suite.Equal(email, result.Email)
	suite.False(result.IsAdmin)
	suite.Equal([]string{"intersection-1", "intersection-2"}, result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_AdminUser() {
	userID := "admin-456"
	username := "Updated Admin Name"
	email := "updated.admin@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	expectedUser := createTestUser(
		userID,
		username,
		email,
		true,
		[]string{"intersection-1", "intersection-2", "intersection-3"},
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(expectedUser, nil)

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal(username, result.Username)
	suite.Equal(email, result.Email)
	suite.True(result.IsAdmin)
	suite.Equal(
		[]string{"intersection-1", "intersection-2", "intersection-3"},
		result.IntersectionIDs,
	)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_MissingUserID() {
	userID := "user-123"
	request := model.UpdateUserRequest{
		Username: "Updated Name",
		Email:    "updated@example.com",
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	// Note: Not setting userID in context to simulate missing userID

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("user ID missing inside of handler", svcError.Message)
}

func (suite *TestSuite) TestUpdateProfile_UserNotFound() {
	userID := "nonexistent-user"
	username := "Updated Name"
	email := "updated@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{"userID": userID}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ValidationError() {
	userID := "user-123"
	username := "Valid Name"
	email := "invalid-email-format"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewValidationError("invalid email format", map[string]any{"email": email}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid email format", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_EmptyUsername() {
	userID := "user-123"
	username := ""
	email := "valid@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewValidationError("username cannot be empty", map[string]any{"username": username}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("username cannot be empty", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_EmptyEmail() {
	userID := "user-123"
	username := "Valid Name"
	email := ""

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewValidationError("email cannot be empty", map[string]any{"email": email}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("email cannot be empty", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_InternalError() {
	userID := "user-123"
	username := "Valid Name"
	email := "valid@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ConflictError() {
	userID := "user-123"
	username := "Valid Name"
	email := "existing@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewConflictError("email already exists", map[string]any{"email": email}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrConflict, svcError.Code)
	suite.Equal("email already exists", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_UnauthorizedError() {
	userID := "user-123"
	username := "Valid Name"
	email := "valid@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewUnauthorizedError("invalid token", map[string]any{}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("invalid token", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ServiceUnavailable() {
	userID := "user-123"
	username := "Valid Name"
	email := "valid@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(nil, errs.NewUnavailableError("user service is temporarily unavailable", map[string]any{}))

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)
	suite.Equal("user service is temporarily unavailable", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_PartialUpdate_UsernameOnly() {
	userID := "user-123"
	username := "Updated Name Only"
	email := "original@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	expectedUser := createTestUser(
		userID,
		username,
		email,
		false,
		[]string{"intersection-1"},
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(expectedUser, nil)

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal(username, result.Username)
	suite.Equal(email, result.Email)
	suite.False(result.IsAdmin)
	suite.Equal([]string{"intersection-1"}, result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_PartialUpdate_EmailOnly() {
	userID := "user-123"
	username := "Original Name"
	email := "newemail@example.com"

	request := model.UpdateUserRequest{
		Username: username,
		Email:    email,
	}

	expectedUser := createTestUser(
		userID,
		username,
		email,
		false,
		[]string{"intersection-1"},
	)

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	suite.client.On("UpdateUser", ctx, userID, username, email).
		Return(expectedUser, nil)

	result, err := suite.service.UpdateProfile(ctx, userID, request)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal(username, result.Username)
	suite.Equal(email, result.Email)
	suite.False(result.IsAdmin)
	suite.Equal([]string{"intersection-1"}, result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}
