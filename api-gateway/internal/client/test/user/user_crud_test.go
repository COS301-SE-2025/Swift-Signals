package user

import (
	"context"
	"testing"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetUserByID Tests
func (suite *TestSuite) TestGetUserByID_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"

	expectedResponse := &userpb.UserResponse{
		Id:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	suite.grpcClient.On("GetUserByID",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == userID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.GetUserByID(ctx, userID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("GetUserByID", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == "test-user-id"
		})).Return(&userpb.UserResponse{Id: "test-user-id"}, nil)

	// Act
	_, err := suite.client.GetUserByID(ctx, "test-user-id")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_UserNotFound() {
	// Arrange
	ctx := context.Background()
	userID := "nonexistent-user"

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("GetUserByID",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetUserByID(ctx, userID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_EmptyUserID() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "user ID is required")

	suite.grpcClient.On("GetUserByID",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetUserByID(ctx, "")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// GetUserByEmail Tests
func (suite *TestSuite) TestGetUserByEmail_Success() {
	// Arrange
	ctx := context.Background()
	email := "test@example.com"

	expectedResponse := &userpb.UserResponse{
		Id:    "user-123",
		Name:  "Test User",
		Email: email,
	}

	suite.grpcClient.On("GetUserByEmail",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.GetUserByEmailRequest) bool {
			return req.Email == email
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.GetUserByEmail(ctx, email)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("GetUserByEmail", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.GetUserByEmailRequest) bool {
			return req.Email == "user@test.com"
		})).Return(&userpb.UserResponse{Email: "user@test.com"}, nil)

	// Act
	_, err := suite.client.GetUserByEmail(ctx, "user@test.com")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_UserNotFound() {
	// Arrange
	ctx := context.Background()
	email := "nonexistent@example.com"

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("GetUserByEmail",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.GetUserByEmailRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetUserByEmail(ctx, email)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_InvalidEmail() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "invalid email format")

	suite.grpcClient.On("GetUserByEmail",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.GetUserByEmailRequest) bool {
			return req.Email == "invalid-email"
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetUserByEmail(ctx, "invalid-email")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// UpdateUser Tests
func (suite *TestSuite) TestUpdateUser_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"
	name := "Updated Name"
	email := "updated@example.com"

	expectedResponse := &userpb.UserResponse{
		Id:    userID,
		Name:  name,
		Email: email,
	}

	suite.grpcClient.On("UpdateUser",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.UpdateUserRequest) bool {
			return req.UserId == userID && req.Name == name && req.Email == email
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.UpdateUser(ctx, userID, name, email)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("UpdateUser", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UpdateUserRequest) bool {
			return req.UserId == "user-id" &&
				req.Name == "New Name" &&
				req.Email == "new@example.com"
		})).Return(&userpb.UserResponse{}, nil)

	// Act
	_, err := suite.client.UpdateUser(ctx, "user-id", "New Name", "new@example.com")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("UpdateUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UpdateUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.UpdateUser(ctx, "nonexistent", "Name", "email@test.com")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmailAlreadyExists() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.AlreadyExists, "email already exists")

	suite.grpcClient.On("UpdateUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UpdateUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.UpdateUser(ctx, "user-123", "Name", "existing@example.com")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmptyFields() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.On("UpdateUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UpdateUserRequest) bool {
			return req.UserId == "user-123" && req.Name == "" && req.Email == ""
		})).Return(&userpb.UserResponse{Id: "user-123"}, nil)

	// Act
	result, err := suite.client.UpdateUser(ctx, "user-123", "", "")

	// Assert
	suite.Require().NoError(err)
	suite.Equal("user-123", result.Id)
	suite.grpcClient.AssertExpectations(suite.T())
}

// DeleteUser Tests
func (suite *TestSuite) TestDeleteUser_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("DeleteUser",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == userID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.DeleteUser(ctx, userID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("DeleteUser", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == "delete-user-id"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.DeleteUser(ctx, "delete-user-id")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("DeleteUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteUser(ctx, "nonexistent")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_EmptyUserID() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "user ID is required")

	suite.grpcClient.On("DeleteUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteUser(ctx, "")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUser_Unauthorized() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.PermissionDenied, "insufficient permissions")

	suite.grpcClient.On("DeleteUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteUser(ctx, "user-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestClientUserCRUD(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
