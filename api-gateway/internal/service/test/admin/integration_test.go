package admin

import (
	grpcmocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestCompleteUserManagementLifecycle() {
	userID := "integration-user-id"
	initialName := "Integration User"
	initialEmail := "integration@example.com"
	updatedName := "Updated Integration User"
	updatedEmail := "updated.integration@example.com"

	ctx := createAdminContext()

	// Step 1: Get user by ID (initial state)
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

	getUserResult, err := suite.service.GetUserByID(ctx, userID)
	suite.Require().NoError(err)
	suite.Equal(userID, getUserResult.ID)
	suite.Equal(initialName, getUserResult.Username)
	suite.Equal(initialEmail, getUserResult.Email)
	suite.False(getUserResult.IsAdmin)

	// Step 2: Update the user
	updatedUser := createTestUser(
		userID,
		updatedName,
		updatedEmail,
		false,
		[]string{"intersection-1"},
	)

	suite.client.On("UpdateUser", ctx, userID, updatedName, updatedEmail).
		Return(updatedUser, nil).
		Once()

	updateResult, err := suite.service.UpdateUserByID(ctx, userID, updatedName, updatedEmail)
	suite.Require().NoError(err)
	suite.Equal(userID, updateResult.ID)
	suite.Equal(updatedName, updateResult.Username)
	suite.Equal(updatedEmail, updateResult.Email)
	suite.False(updateResult.IsAdmin)

	// Step 3: Delete the user
	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, nil).
		Once()

	err = suite.service.DeleteUserByID(ctx, userID)
	suite.Require().NoError(err)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPermissionConsistencyAcrossOperations() {
	userID := "permission-test-user"

	// Test all operations with non-admin role
	ctx := createUserContext()

	// Test GetAllUsers - should be forbidden
	_, err := suite.service.GetAllUsers(ctx, 1, 10)
	suite.Require().Error(err)
	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test GetUserByID - should be forbidden
	_, err = suite.service.GetUserByID(ctx, userID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test UpdateUserByID - should be forbidden
	_, err = suite.service.UpdateUserByID(ctx, userID, "New Name", "new@example.com")
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test DeleteUserByID - should be forbidden
	err = suite.service.DeleteUserByID(ctx, userID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
}

func (suite *TestSuite) TestPermissionConsistencyWithNoRole() {
	userID := "permission-test-user"

	// Test all operations with no role context
	ctx := createContextWithoutRole()

	// Test GetAllUsers - should be forbidden
	_, err := suite.service.GetAllUsers(ctx, 1, 10)
	suite.Require().Error(err)
	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test GetUserByID - should be forbidden
	_, err = suite.service.GetUserByID(ctx, userID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test UpdateUserByID - should be forbidden
	_, err = suite.service.UpdateUserByID(ctx, userID, "New Name", "new@example.com")
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test DeleteUserByID - should be forbidden
	err = suite.service.DeleteUserByID(ctx, userID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
}

func (suite *TestSuite) TestErrorPropagationConsistency() {
	userID := "test-user"
	ctx := createAdminContext()

	// Test that internal errors are properly propagated across different operations
	internalError := errs.NewInternalError("database connection failed", nil, map[string]any{})

	// Test GetUserByID error propagation
	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, internalError).
		Once()

	_, err := suite.service.GetUserByID(ctx, userID)
	suite.Require().Error(err)
	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	// Test UpdateUser error propagation
	suite.client.On("UpdateUser", ctx, userID, "Test Name", "test@example.com").
		Return(nil, internalError).
		Once()

	_, err = suite.service.UpdateUserByID(ctx, userID, "Test Name", "test@example.com")
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	// Test DeleteUser error propagation
	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, internalError).
		Once()

	err = suite.service.DeleteUserByID(ctx, userID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestBulkUserOperationsWorkflow() {
	ctx := createAdminContext()

	// Step 1: Get all users (should return empty initially)
	mockStream := suite.NewMockGetAllUsersStream()
	mockStream.On("Recv").Return(nil, errs.NewNotFoundError("no users found", map[string]any{}))

	suite.client.On("GetAllUsers", ctx, int32(1), int32(10), "").
		Return(mockStream, nil).
		Once()

	users, err := suite.service.GetAllUsers(ctx, 1, 10)
	suite.Require().Error(err) // Should error when no users found
	suite.Empty(users)

	// Step 2: Create scenario where users exist and update multiple users
	userIDs := []string{"user-1", "user-2", "user-3"}

	for i, userID := range userIDs {
		originalUser := createTestUser(
			userID,
			"Original User "+string(rune(i+1)),
			"original"+string(rune(i+1))+"@example.com",
			false,
			[]string{},
		)

		updatedUser := createTestUser(
			userID,
			"Updated User "+string(rune(i+1)),
			"updated"+string(rune(i+1))+"@example.com",
			false,
			[]string{},
		)

		// Get original user
		suite.client.On("GetUserByID", ctx, userID).
			Return(originalUser, nil).
			Once()

		// Update user
		suite.client.On("UpdateUser", ctx, userID, "Updated User "+string(rune(i+1)), "updated"+string(rune(i+1))+"@example.com").
			Return(updatedUser, nil).
			Once()
	}

	// Execute the workflow
	for i, userID := range userIDs {
		// Get user
		user, err := suite.service.GetUserByID(ctx, userID)
		suite.Require().NoError(err)
		suite.Equal("Original User "+string(rune(i+1)), user.Username)

		// Update user
		updatedUser, err := suite.service.UpdateUserByID(
			ctx,
			userID,
			"Updated User "+string(rune(i+1)),
			"updated"+string(rune(i+1))+"@example.com",
		)
		suite.Require().NoError(err)
		suite.Equal("Updated User "+string(rune(i+1)), updatedUser.Username)
		suite.Equal("updated"+string(rune(i+1))+"@example.com", updatedUser.Email)
	}

	suite.client.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestConcurrentAccessSimulation() {
	userID := "concurrent-user"
	ctx := createAdminContext()

	// Simulate concurrent read and update operations
	user := createTestUser(userID, "Concurrent User", "concurrent@example.com", false, []string{})

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
	suite.client.On("UpdateUser", ctx, userID, "Updated Concurrent User", "updated.concurrent@example.com").
		Return(updatedUser, nil).
		Once()

	// Simulate multiple reads
	for i := 0; i < 3; i++ {
		result, err := suite.service.GetUserByID(ctx, userID)
		suite.Require().NoError(err)
		suite.Equal(userID, result.ID)
		suite.Equal("Concurrent User", result.Username)
	}

	// Simulate update
	result, err := suite.service.UpdateUserByID(
		ctx,
		userID,
		"Updated Concurrent User",
		"updated.concurrent@example.com",
	)
	suite.Require().NoError(err)
	suite.Equal("Updated Concurrent User", result.Username)

	suite.client.AssertExpectations(suite.T())
}

// Helper method for creating mock streams in integration tests
func (suite *TestSuite) NewMockGetAllUsersStream() *grpcmocks.MockUserService_GetAllUsersClient[userpb.UserResponse] {
	return grpcmocks.NewMockUserService_GetAllUsersClient[userpb.UserResponse](suite.T())
}
