package profile

import (
	"context"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestCompleteProfileLifecycle() {
	userID := "integration-user-id"
	initialName := "Integration User"
	initialEmail := "integration@example.com"
	updatedName := "Updated Integration User"
	updatedEmail := "updated.integration@example.com"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	// Step 1: Get profile (initial state)
	initialUser := createTestUser(
		userID,
		initialName,
		initialEmail,
		false,
		[]string{"intersection-1"},
	)

	suite.client.On("GetUserByID", ctx, userID).
		Return(initialUser, nil).
		Once()

	getProfileResult, err := suite.service.GetProfile(ctx, userID)
	suite.Require().NoError(err)
	suite.Equal(userID, getProfileResult.ID)
	suite.Equal(initialName, getProfileResult.Username)
	suite.Equal(initialEmail, getProfileResult.Email)
	suite.False(getProfileResult.IsAdmin)

	// Step 2: Update the profile
	updatedUser := createTestUser(
		userID,
		updatedName,
		updatedEmail,
		false,
		[]string{"intersection-1"},
	)

	updateRequest := model.UpdateUserRequest{
		Username: updatedName,
		Email:    updatedEmail,
	}

	suite.client.On("UpdateUser", ctx, userID, updatedName, updatedEmail).
		Return(updatedUser, nil).
		Once()

	updateResult, err := suite.service.UpdateProfile(ctx, userID, updateRequest)
	suite.Require().NoError(err)
	suite.Equal(userID, updateResult.ID)
	suite.Equal(updatedName, updateResult.Username)
	suite.Equal(updatedEmail, updateResult.Email)
	suite.False(updateResult.IsAdmin)

	// Step 3: Get profile again to verify update
	suite.client.On("GetUserByID", ctx, userID).
		Return(updatedUser, nil).
		Once()

	getUpdatedResult, err := suite.service.GetProfile(ctx, userID)
	suite.Require().NoError(err)
	suite.Equal(userID, getUpdatedResult.ID)
	suite.Equal(updatedName, getUpdatedResult.Username)
	suite.Equal(updatedEmail, getUpdatedResult.Email)

	// Step 4: Delete the profile
	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, nil).
		Once()

	err = suite.service.DeleteProfile(ctx, userID)
	suite.Require().NoError(err)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestProfileErrorPropagationConsistency() {
	userID := "test-user"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	// Test that internal errors are properly propagated across different operations
	internalError := errs.NewInternalError("database connection failed", nil, map[string]any{})

	// Test GetProfile error propagation
	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, internalError).
		Once()

	_, err := suite.service.GetProfile(ctx, userID)
	suite.Require().Error(err)
	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	// Test UpdateProfile error propagation
	updateRequest := model.UpdateUserRequest{
		Username: "Test Name",
		Email:    "test@example.com",
	}

	suite.client.On("UpdateUser", ctx, userID, "Test Name", "test@example.com").
		Return(nil, internalError).
		Once()

	_, err = suite.service.UpdateProfile(ctx, userID, updateRequest)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	// Test DeleteProfile error propagation
	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, internalError).
		Once()

	err = suite.service.DeleteProfile(ctx, userID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestProfileOperationsWithMissingUserID() {
	userID := "test-user"

	// Test all operations with missing userID in context
	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	// Note: Not setting userID in context

	// Test GetProfile - should fail with internal error
	_, err := suite.service.GetProfile(ctx, userID)
	suite.Require().Error(err)
	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("user ID missing inside of handler", svcError.Message)

	// Test UpdateProfile - should fail with internal error
	updateRequest := model.UpdateUserRequest{
		Username: "Test Name",
		Email:    "test@example.com",
	}
	_, err = suite.service.UpdateProfile(ctx, userID, updateRequest)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("user ID missing inside of handler", svcError.Message)

	// Test DeleteProfile - should fail with internal error
	err = suite.service.DeleteProfile(ctx, userID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("user ID missing inside of handler", svcError.Message)
}

func (suite *TestSuite) TestProfileOperationsWithDifferentUserTypes() {
	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Test 1: Regular user profile operations
	regularUserID := "regular-user"
	ctx1 := middleware.SetUserID(ctx, regularUserID)

	regularUser := createTestUser(
		regularUserID,
		"Regular User",
		"regular@example.com",
		false,
		[]string{"intersection-1"},
	)

	suite.client.On("GetUserByID", ctx1, regularUserID).
		Return(regularUser, nil).
		Once()

	result, err := suite.service.GetProfile(ctx1, regularUserID)
	suite.Require().NoError(err)
	suite.False(result.IsAdmin)
	suite.Equal([]string{"intersection-1"}, result.IntersectionIDs)

	// Test 2: Admin user profile operations
	adminUserID := "admin-user"
	ctx2 := middleware.SetUserID(ctx, adminUserID)

	adminUser := createTestUser(
		adminUserID,
		"Admin User",
		"admin@example.com",
		true,
		[]string{"intersection-1", "intersection-2", "intersection-3"},
	)

	suite.client.On("GetUserByID", ctx2, adminUserID).
		Return(adminUser, nil).
		Once()

	result, err = suite.service.GetProfile(ctx2, adminUserID)
	suite.Require().NoError(err)
	suite.True(result.IsAdmin)
	suite.Equal(
		[]string{"intersection-1", "intersection-2", "intersection-3"},
		result.IntersectionIDs,
	)

	// Test 3: User with no intersections
	noIntersectionUserID := "no-intersection-user"
	ctx3 := middleware.SetUserID(ctx, noIntersectionUserID)

	noIntersectionUser := createTestUser(
		noIntersectionUserID,
		"No Intersection User",
		"nointersection@example.com",
		false,
		[]string{},
	)

	suite.client.On("GetUserByID", ctx3, noIntersectionUserID).
		Return(noIntersectionUser, nil).
		Once()

	result, err = suite.service.GetProfile(ctx3, noIntersectionUserID)
	suite.Require().NoError(err)
	suite.False(result.IsAdmin)
	suite.Empty(result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMultipleProfileUpdates() {
	userID := "multi-update-user"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	updates := []struct {
		username string
		email    string
	}{
		{"First Update", "first@example.com"},
		{"Second Update", "second@example.com"},
		{"Final Update", "final@example.com"},
	}

	for i, update := range updates {
		expectedUser := createTestUser(
			userID,
			update.username,
			update.email,
			false,
			[]string{"intersection-1"},
		)

		updateRequest := model.UpdateUserRequest{
			Username: update.username,
			Email:    update.email,
		}

		suite.client.On("UpdateUser", ctx, userID, update.username, update.email).
			Return(expectedUser, nil).
			Once()

		result, err := suite.service.UpdateProfile(ctx, userID, updateRequest)
		suite.Require().NoError(err, "Failed for update %d", i+1)
		suite.Equal(update.username, result.Username)
		suite.Equal(update.email, result.Email)
	}

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestProfileDeletionScenarios() {
	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	testCases := []struct {
		name        string
		userID      string
		setupError  *errs.ServiceError
		expectError bool
	}{
		{
			name:        "Successful deletion",
			userID:      "success-user",
			setupError:  nil,
			expectError: false,
		},
		{
			name:        "User not found",
			userID:      "notfound-user",
			setupError:  errs.NewNotFoundError("user not found", map[string]any{}),
			expectError: true,
		},
		{
			name:   "User with active intersections",
			userID: "active-user",
			setupError: errs.NewForbiddenError(
				"cannot delete user with active intersections",
				map[string]any{},
			),
			expectError: true,
		},
		{
			name:        "Database error",
			userID:      "db-error-user",
			setupError:  errs.NewInternalError("database error", nil, map[string]any{}),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		userCtx := middleware.SetUserID(ctx, tc.userID)

		if tc.setupError != nil {
			suite.client.On("DeleteUser", userCtx, tc.userID).
				Return(nil, tc.setupError).
				Once()
		} else {
			suite.client.On("DeleteUser", userCtx, tc.userID).
				Return(nil, nil).
				Once()
		}

		err := suite.service.DeleteProfile(userCtx, tc.userID)

		if tc.expectError {
			suite.Require().Error(err, "Expected error for test case: %s", tc.name)
		} else {
			suite.Require().NoError(err, "Expected no error for test case: %s", tc.name)
		}
	}

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestConcurrentProfileOperations() {
	userID := "concurrent-user"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	// Simulate concurrent read and update operations
	user := createTestUser(
		userID,
		"Concurrent User",
		"concurrent@example.com",
		false,
		[]string{},
	)

	// Setup expectations for concurrent reads
	suite.client.On("GetUserByID", ctx, userID).
		Return(user, nil).
		Times(3)

	// Setup expectations for update
	updatedUser := createTestUser(
		userID,
		"Updated Concurrent User",
		"updated.concurrent@example.com",
		false,
		[]string{},
	)

	updateRequest := model.UpdateUserRequest{
		Username: "Updated Concurrent User",
		Email:    "updated.concurrent@example.com",
	}

	suite.client.On("UpdateUser", ctx, userID, "Updated Concurrent User", "updated.concurrent@example.com").
		Return(updatedUser, nil).
		Once()

	// Simulate multiple reads
	for i := 0; i < 3; i++ {
		result, err := suite.service.GetProfile(ctx, userID)
		suite.Require().NoError(err)
		suite.Equal(userID, result.ID)
		suite.Equal("Concurrent User", result.Username)
	}

	// Simulate update
	result, err := suite.service.UpdateProfile(ctx, userID, updateRequest)
	suite.Require().NoError(err)
	suite.Equal("Updated Concurrent User", result.Username)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestProfileServiceResilience() {
	userID := "resilience-user"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	// Test service behavior with various network/service issues
	testCases := []struct {
		name      string
		operation string
		error     *errs.ServiceError
	}{
		{
			name:      "Service unavailable",
			operation: "get",
			error: errs.NewUnavailableError(
				"service temporarily unavailable",
				map[string]any{},
			),
		},
		{
			name:      "Timeout error",
			operation: "update",
			error:     errs.NewUnavailableError("request timeout", map[string]any{}),
		},
		{
			name:      "Rate limit error",
			operation: "delete",
			error:     errs.NewForbiddenError("rate limit exceeded", map[string]any{}),
		},
	}

	for _, tc := range testCases {
		switch tc.operation {
		case "get":
			suite.client.On("GetUserByID", ctx, userID).
				Return(nil, tc.error).
				Once()

			_, err := suite.service.GetProfile(ctx, userID)
			suite.Require().Error(err, "Expected error for test case: %s", tc.name)

		case "update":
			updateRequest := model.UpdateUserRequest{
				Username: "Test User",
				Email:    "test@example.com",
			}

			suite.client.On("UpdateUser", ctx, userID, "Test User", "test@example.com").
				Return(nil, tc.error).
				Once()

			_, err := suite.service.UpdateProfile(ctx, userID, updateRequest)
			suite.Require().Error(err, "Expected error for test case: %s", tc.name)

		case "delete":
			suite.client.On("DeleteUser", ctx, userID).
				Return(nil, tc.error).
				Once()

			err := suite.service.DeleteProfile(ctx, userID)
			suite.Require().Error(err, "Expected error for test case: %s", tc.name)
		}
	}

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestContextPropagation() {
	userID := "context-user"

	// Test that context is properly propagated through all service calls
	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)
	ctx = middleware.SetUserID(ctx, userID)

	user := createTestUser(
		userID,
		"Context User",
		"context@example.com",
		false,
		[]string{},
	)

	// The key test here is that the exact context is passed through
	suite.client.On("GetUserByID", ctx, userID).Return(user, nil)

	_, err := suite.service.GetProfile(ctx, userID)
	suite.Require().NoError(err)

	suite.client.AssertExpectations(suite.T())
}
