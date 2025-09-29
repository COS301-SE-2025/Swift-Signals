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
)

func (suite *TestSuite) TestLoginUser_Success() {
	// Arrange
	ctx := context.Background()
	email := "test@example.com"
	password := "testpassword123"

	expectedResponse := &userpb.LoginUserResponse{
		Token: "jwt-token-123",
	}

	suite.grpcClient.On("LoginUser",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.LoginUserRequest) bool {
			return req.Email == email && req.Password == password
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.LoginUser(ctx, email, password)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("LoginUser", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.LoginUserRequest) bool {
			return req.Email == "user@test.com" &&
				req.Password == "password123"
		})).Return(&userpb.LoginUserResponse{
		Token: "test-token",
	}, nil)

	// Act
	_, err := suite.client.LoginUser(ctx, "user@test.com", "password123")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_SetsTimeoutContext() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.On("LoginUser",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.AnythingOfType("*user.LoginUserRequest")).
		Return(&userpb.LoginUserResponse{}, nil)

	// Act
	_, err := suite.client.LoginUser(ctx, "test@example.com", "password")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_InvalidCredentials() {
	// Arrange
	ctx := context.Background()
	email := "wrong@example.com"
	password := "wrongpassword"

	grpcErr := status.Error(codes.Unauthenticated, "invalid credentials")

	suite.grpcClient.On("LoginUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.LoginUserRequest) bool {
			return req.Email == email && req.Password == password
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.LoginUser(ctx, email, password)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_EmptyCredentials() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "email and password required")

	suite.grpcClient.On("LoginUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.LoginUserRequest) bool {
			return req.Email == "" && req.Password == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.LoginUser(ctx, "", "")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_UserNotFound() {
	// Arrange
	ctx := context.Background()
	email := "nonexistent@example.com"
	password := "somepassword"

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("LoginUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.LoginUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.LoginUser(ctx, email, password)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_ContextTimeout() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	suite.grpcClient.On("LoginUser",
		mock.Anything,
		mock.AnythingOfType("*user.LoginUserRequest")).
		Return(nil, context.DeadlineExceeded)

	// Act
	result, err := suite.client.LoginUser(ctx, "test@example.com", "password")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_InternalServerError() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.Internal, "internal server error")

	suite.grpcClient.On("LoginUser",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.LoginUserRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.LoginUser(ctx, "test@example.com", "password")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestClientLoginUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
