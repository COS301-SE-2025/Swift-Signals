package user

import (
	"context"
	"testing"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// LogoutUser Tests
func (suite *TestSuite) TestLogoutUser_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("LogoutUser",
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
	result, err := suite.client.LogoutUser(ctx, userID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("LogoutUser", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == "logout-user-id"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.LogoutUser(ctx, "logout-user-id")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_SetsTimeoutContext() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.On("LogoutUser",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.LogoutUser(ctx, "test-user")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("LogoutUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.UserIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.LogoutUser(ctx, "nonexistent-user")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_EmptyUserID() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "user ID is required")

	suite.grpcClient.On("LogoutUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.UserIDRequest) bool {
			return req.UserId == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.LogoutUser(ctx, "")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// ChangePassword Tests
func (suite *TestSuite) TestChangePassword_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"
	currentPassword := "oldpassword123"
	newPassword := "newpassword456"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("ChangePassword",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.ChangePasswordRequest) bool {
			return req.UserId == userID &&
				req.CurrentPassword == currentPassword &&
				req.NewPassword == newPassword
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.ChangePassword(ctx, userID, currentPassword, newPassword)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("ChangePassword", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.ChangePasswordRequest) bool {
			return req.UserId == "test-user" &&
				req.CurrentPassword == "current123" &&
				req.NewPassword == "new456"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.ChangePassword(ctx, "test-user", "current123", "new456")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_IncorrectCurrentPassword() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.Unauthenticated, "current password is incorrect")

	suite.grpcClient.On("ChangePassword",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.ChangePasswordRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.ChangePassword(ctx, "user-123", "wrongpassword", "newpassword")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_WeakNewPassword() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "new password does not meet requirements")

	suite.grpcClient.On("ChangePassword",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.ChangePasswordRequest) bool {
			return req.NewPassword == "123"
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.ChangePassword(ctx, "user-123", "currentpassword", "123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("ChangePassword",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.ChangePasswordRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.ChangePassword(ctx, "nonexistent", "current", "new")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_EmptyFields() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "all fields are required")

	suite.grpcClient.On("ChangePassword",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.ChangePasswordRequest) bool {
			return req.UserId == "" || req.CurrentPassword == "" || req.NewPassword == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.ChangePassword(ctx, "", "", "")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// ResetPassword Tests
func (suite *TestSuite) TestResetPassword_Success() {
	// Arrange
	ctx := context.Background()
	email := "user@example.com"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("ResetPassword",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.ResetPasswordRequest) bool {
			return req.Email == email
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.ResetPassword(ctx, email)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestResetPassword_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("ResetPassword", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.ResetPasswordRequest) bool {
			return req.Email == "reset@test.com"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.ResetPassword(ctx, "reset@test.com")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestResetPassword_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user with email not found")

	suite.grpcClient.On("ResetPassword",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.ResetPasswordRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.ResetPassword(ctx, "nonexistent@example.com")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestResetPassword_InvalidEmail() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "invalid email format")

	suite.grpcClient.On("ResetPassword",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.ResetPasswordRequest) bool {
			return req.Email == "invalid-email"
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.ResetPassword(ctx, "invalid-email")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestResetPassword_EmptyEmail() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "email is required")

	suite.grpcClient.On("ResetPassword",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.ResetPasswordRequest) bool {
			return req.Email == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.ResetPassword(ctx, "")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// MakeAdmin Tests
func (suite *TestSuite) TestMakeAdmin_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"
	adminUserID := "admin-456"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("MakeAdmin",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.AdminRequest) bool {
			return req.UserId == userID && req.AdminUserId == adminUserID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.MakeAdmin(ctx, userID, adminUserID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("MakeAdmin", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.AdminRequest) bool {
			return req.UserId == "target-user" && req.AdminUserId == "admin-user"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.MakeAdmin(ctx, "target-user", "admin-user")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_Unauthorized() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.PermissionDenied, "insufficient permissions to make admin")

	suite.grpcClient.On("MakeAdmin",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AdminRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.MakeAdmin(ctx, "user-123", "non-admin")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("MakeAdmin",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AdminRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.MakeAdmin(ctx, "nonexistent", "admin-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_AlreadyAdmin() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.AlreadyExists, "user is already an admin")

	suite.grpcClient.On("MakeAdmin",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AdminRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.MakeAdmin(ctx, "admin-user", "super-admin")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// RemoveAdmin Tests
func (suite *TestSuite) TestRemoveAdmin_Success() {
	// Arrange
	ctx := context.Background()
	userID := "admin-user-123"
	adminUserID := "super-admin-456"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("RemoveAdmin",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.AdminRequest) bool {
			return req.UserId == userID && req.AdminUserId == adminUserID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.RemoveAdmin(ctx, userID, adminUserID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveAdmin_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("RemoveAdmin", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.AdminRequest) bool {
			return req.UserId == "remove-admin" && req.AdminUserId == "super-admin"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.RemoveAdmin(ctx, "remove-admin", "super-admin")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveAdmin_Unauthorized() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.PermissionDenied, "insufficient permissions to remove admin")

	suite.grpcClient.On("RemoveAdmin",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AdminRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RemoveAdmin(ctx, "admin-123", "regular-user")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveAdmin_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("RemoveAdmin",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AdminRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RemoveAdmin(ctx, "nonexistent", "admin-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveAdmin_NotAnAdmin() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "user is not an admin")

	suite.grpcClient.On("RemoveAdmin",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AdminRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RemoveAdmin(ctx, "regular-user", "admin-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveAdmin_CannotRemoveSelf() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "cannot remove admin privileges from yourself")

	suite.grpcClient.On("RemoveAdmin",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.AdminRequest) bool {
			return req.UserId == req.AdminUserId
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.RemoveAdmin(ctx, "admin-123", "admin-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestClientAuthAdmin(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
