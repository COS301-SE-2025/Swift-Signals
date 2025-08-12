package user

import (
	"context"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ConstructorTestSuite for testing constructors separately
type ConstructorTestSuite struct {
	suite.Suite
}

// Constructor Tests
func (suite *ConstructorTestSuite) TestNewUserClient_Success() {
	// Arrange
	mockGrpcClient := new(mocks.MockUserServiceClient)

	// Act
	client := client.NewUserClient(mockGrpcClient)

	// Assert
	suite.Require().NotNil(client)
	// We can't directly access the internal client field, but we can test that methods work
}

func (suite *ConstructorTestSuite) TestNewUserClient_NilClient() {
	// Act
	client := client.NewUserClient(nil)

	// Assert
	suite.Require().NotNil(client)
	// Client should still be created even with nil grpc client
	// This tests defensive programming - the constructor doesn't validate input
}

func (suite *ConstructorTestSuite) TestNewUserClientFromConn_Success() {
	// Note: This test would require a real grpc.ClientConn which is complex to mock
	// In a real scenario, you might want to use a test server or integration test
	// For now, we'll document that this should be tested in integration tests
	suite.T().Skip("NewUserClientFromConn requires integration testing with real gRPC connection")
}

// Error Handling Tests
func (suite *TestSuite) TestRegisterUser_GrpcError() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.AlreadyExists, "user already exists")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(
		ctx,
		"Existing User",
		"existing@example.com",
		"password123",
	)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_InvalidArguments() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "invalid email format")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RegisterUserRequest) bool {
			return req.Email == "invalid-email"
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "invalid-email", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_WeakPassword() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "password does not meet requirements")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RegisterUserRequest) bool {
			return req.Password == "123"
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "test@example.com", "123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_EmptyName() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "name is required")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RegisterUserRequest) bool {
			return req.Name == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(ctx, "", "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_ContextTimeout() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	suite.grpcClient.On("RegisterUser",
		mock.Anything,
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(nil, context.DeadlineExceeded)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_InternalServerError() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.Internal, "internal server error")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Unavailable() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.Unavailable, "service unavailable")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// Context-related error tests
func (suite *TestSuite) TestContextCancellation_RegisterUser() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	suite.grpcClient.On("RegisterUser",
		mock.Anything,
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(nil, context.Canceled)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestContextCancellation_LoginUser() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	suite.grpcClient.On("LoginUser",
		mock.Anything,
		mock.AnythingOfType("*user.LoginUserRequest")).
		Return(nil, context.Canceled)

	// Act
	result, err := suite.client.LoginUser(ctx, "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestContextCancellation_LogoutUser() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	suite.grpcClient.On("LogoutUser",
		mock.Anything,
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(nil, context.Canceled)

	// Act
	result, err := suite.client.LogoutUser(ctx, "user-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestContextCancellation_GetUserByID() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	suite.grpcClient.On("GetUserByID",
		mock.Anything,
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(nil, context.Canceled)

	// Act
	result, err := suite.client.GetUserByID(ctx, "user-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestContextCancellation_UpdateUser() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	suite.grpcClient.On("UpdateUser",
		mock.Anything,
		mock.AnythingOfType("*user.UpdateUserRequest")).
		Return(nil, context.Canceled)

	// Act
	result, err := suite.client.UpdateUser(ctx, "user-123", "New Name", "new@example.com")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestContextCancellation_DeleteUser() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	suite.grpcClient.On("DeleteUser",
		mock.Anything,
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(nil, context.Canceled)

	// Act
	result, err := suite.client.DeleteUser(ctx, "user-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestContextCancellation_ChangePassword() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	suite.grpcClient.On("ChangePassword",
		mock.Anything,
		mock.AnythingOfType("*user.ChangePasswordRequest")).
		Return(nil, context.Canceled)

	// Act
	result, err := suite.client.ChangePassword(ctx, "user-123", "oldpass", "newpass")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// Network-related error tests
func (suite *TestSuite) TestNetworkError_RegisterUser() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.DeadlineExceeded, "deadline exceeded")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestResourceExhausted_RegisterUser() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.ResourceExhausted, "rate limit exceeded")

	suite.grpcClient.On("RegisterUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RegisterUser(ctx, "Test User", "test@example.com", "password123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// Comprehensive error handling for empty response tests
func (suite *TestSuite) TestLogoutUser_NilResponse() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.On("LogoutUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UserIDRequest")).
		Return((*emptypb.Empty)(nil), nil)

	// Act
	result, err := suite.client.LogoutUser(ctx, "user-123")

	// Assert
	suite.Require().NoError(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_NilResponse() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.On("DeleteUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UserIDRequest")).
		Return((*emptypb.Empty)(nil), nil)

	// Act
	result, err := suite.client.DeleteUser(ctx, "user-123")

	// Assert
	suite.Require().NoError(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestConstructorAndErrorHandling(t *testing.T) {
	// Run constructor tests
	constructorSuite := new(ConstructorTestSuite)
	suite.Run(t, constructorSuite)

	// Run error handling tests
	errorSuite := new(TestSuite)
	suite.Run(t, errorSuite)
}
