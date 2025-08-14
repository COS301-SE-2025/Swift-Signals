package admin

import (
	"context"
	"io"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	grpcmocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

// Helper functions
func createAdminContext() context.Context {
	return middleware.SetRole(context.Background(), "admin")
}

func createUserContext() context.Context {
	return middleware.SetRole(context.Background(), "user")
}

func createContextWithoutRole() context.Context {
	return context.Background()
}

// This helper function is moved to individual test methods as needed

func (suite *TestSuite) TestGetAllUsers_Success() {
	pgNo := 1
	pgSz := 10

	expectedUsers := []*userpb.UserResponse{
		createTestUser(
			"user-1",
			"John Doe",
			"john@example.com",
			false,
			[]string{"intersection-1", "intersection-2"},
		),
		createTestUser(
			"user-2",
			"Jane Smith",
			"jane@example.com",
			true,
			[]string{"intersection-3"},
		),
		createTestUser("user-3", "Bob Wilson", "bob@example.com", false, []string{}),
	}

	ctx := createAdminContext()

	mockStream := grpcmocks.NewMockUserService_GetAllUsersClient[userpb.UserResponse](suite.T())
	for _, user := range expectedUsers {
		mockStream.On("Recv").Return(user, nil).Once()
	}
	mockStream.On("Recv").Return(nil, io.EOF).Once()

	suite.client.On("GetAllUsers", ctx, int32(pgNo), int32(pgSz), "").
		Return(mockStream, nil)

	result, err := suite.service.GetAllUsers(ctx, pgNo, pgSz)

	suite.Require().NoError(err)
	suite.Len(result, 3)

	// Verify first user
	suite.Equal("user-1", result[0].ID)
	suite.Equal("John Doe", result[0].Username)
	suite.Equal("john@example.com", result[0].Email)
	suite.False(result[0].IsAdmin)
	suite.Equal([]string{"intersection-1", "intersection-2"}, result[0].IntersectionIDs)

	// Verify second user (admin)
	suite.Equal("user-2", result[1].ID)
	suite.Equal("Jane Smith", result[1].Username)
	suite.Equal("jane@example.com", result[1].Email)
	suite.True(result[1].IsAdmin)
	suite.Equal([]string{"intersection-3"}, result[1].IntersectionIDs)

	// Verify third user (no intersections)
	suite.Equal("user-3", result[2].ID)
	suite.Equal("Bob Wilson", result[2].Username)
	suite.Equal("bob@example.com", result[2].Email)
	suite.False(result[2].IsAdmin)
	suite.Empty(result[2].IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_EmptyResult() {
	pgNo := 1
	pgSz := 10

	ctx := createAdminContext()

	mockStream := grpcmocks.NewMockUserService_GetAllUsersClient[userpb.UserResponse](suite.T())
	mockStream.On("Recv").Return(nil, io.EOF)

	suite.client.On("GetAllUsers", ctx, int32(pgNo), int32(pgSz), "").
		Return(mockStream, nil)

	result, err := suite.service.GetAllUsers(ctx, pgNo, pgSz)

	suite.Require().NoError(err)
	suite.Empty(result)

	suite.client.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_Forbidden_NonAdmin() {
	pgNo := 1
	pgSz := 10

	ctx := createUserContext()

	result, err := suite.service.GetAllUsers(ctx, pgNo, pgSz)

	suite.Require().Error(err)
	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestGetAllUsers_Forbidden_NoRole() {
	pgNo := 1
	pgSz := 10

	ctx := createContextWithoutRole()

	result, err := suite.service.GetAllUsers(ctx, pgNo, pgSz)

	suite.Require().Error(err)
	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestGetAllUsers_UserServiceError() {
	pgNo := 1
	pgSz := 10

	ctx := createAdminContext()

	suite.client.On("GetAllUsers", ctx, int32(pgNo), int32(pgSz), "").
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	result, err := suite.service.GetAllUsers(ctx, pgNo, pgSz)

	suite.Require().Error(err)
	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_StreamError() {
	pgNo := 1
	pgSz := 10

	ctx := createAdminContext()

	mockStream := grpcmocks.NewMockUserService_GetAllUsersClient[userpb.UserResponse](suite.T())
	mockStream.On("Recv").Return(nil, errs.NewInternalError("stream error", nil, map[string]any{}))

	suite.client.On("GetAllUsers", ctx, int32(pgNo), int32(pgSz), "").
		Return(mockStream, nil)

	result, err := suite.service.GetAllUsers(ctx, pgNo, pgSz)

	suite.Require().Error(err)
	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("unable to get all users", svcError.Message)

	suite.client.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}
