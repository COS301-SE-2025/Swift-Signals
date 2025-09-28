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

// Note: GetAllUsers and GetUserIntersectionIDs streaming tests are skipped
// as they require additional mock generation for the streaming interfaces.
// These methods return streaming clients that need proper gRPC streaming mocks.

// AddIntersectionID Tests
func (suite *TestSuite) TestAddIntersectionID_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"
	intersectionID := "intersection-456"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("AddIntersectionID",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.AddIntersectionIDRequest) bool {
			return req.UserId == userID && req.IntersectionId == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.AddIntersectionID(ctx, userID, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("AddIntersectionID", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.AddIntersectionIDRequest) bool {
			return req.UserId == "test-user" && req.IntersectionId == "test-intersection"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.AddIntersectionID(ctx, "test-user", "test-intersection")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("AddIntersectionID",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AddIntersectionIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.AddIntersectionID(ctx, "nonexistent", "intersection-123")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_IntersectionAlreadyAdded() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.AlreadyExists, "intersection already added to user")

	suite.grpcClient.On("AddIntersectionID",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.AddIntersectionIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.AddIntersectionID(ctx, "user-123", "existing-intersection")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_EmptyFields() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "user ID and intersection ID are required")

	suite.grpcClient.On("AddIntersectionID",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.AddIntersectionIDRequest) bool {
			return req.UserId == "" || req.IntersectionId == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.AddIntersectionID(ctx, "", "")

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

// Note: GetUserIntersectionIDs streaming tests are skipped
// as they require additional mock generation for the streaming interface.

// RemoveIntersectionID Tests (single)
func (suite *TestSuite) TestRemoveIntersectionID_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"
	intersectionID := "intersection-456"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("RemoveIntersectionIDs",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.RemoveIntersectionIDRequest) bool {
			return req.UserId == userID &&
				len(req.IntersectionId) == 1 &&
				req.IntersectionId[0] == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.RemoveIntersectionID(ctx, userID, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionID_CallsRemoveIntersectionIDs() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("RemoveIntersectionIDs", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RemoveIntersectionIDRequest) bool {
			return req.UserId == "test-user" &&
				len(req.IntersectionId) == 1 &&
				req.IntersectionId[0] == "single-intersection"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.RemoveIntersectionID(ctx, "test-user", "single-intersection")

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

// RemoveIntersectionIDs Tests (multiple)
func (suite *TestSuite) TestRemoveIntersectionIDs_Success() {
	// Arrange
	ctx := context.Background()
	userID := "user-123"
	intersectionIDs := []string{"intersection-1", "intersection-2", "intersection-3"}

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("RemoveIntersectionIDs",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *userpb.RemoveIntersectionIDRequest) bool {
			return req.UserId == userID &&
				len(req.IntersectionId) == 3 &&
				req.IntersectionId[0] == "intersection-1" &&
				req.IntersectionId[1] == "intersection-2" &&
				req.IntersectionId[2] == "intersection-3"
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_BuildsCorrectRequest() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.Mock.On("RemoveIntersectionIDs", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RemoveIntersectionIDRequest) bool {
			return req.UserId == "bulk-user" &&
				len(req.IntersectionId) == 2 &&
				req.IntersectionId[0] == "bulk-1" &&
				req.IntersectionId[1] == "bulk-2"
		})).Return(&emptypb.Empty{}, nil)

	// Act
	_, err := suite.client.RemoveIntersectionIDs(ctx, "bulk-user", []string{"bulk-1", "bulk-2"})

	// Assert
	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_EmptyList() {
	// Arrange
	ctx := context.Background()

	suite.grpcClient.On("RemoveIntersectionIDs",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RemoveIntersectionIDRequest) bool {
			return req.UserId == "user-123" && len(req.IntersectionId) == 0
		})).Return(&emptypb.Empty{}, nil)

	// Act
	result, err := suite.client.RemoveIntersectionIDs(ctx, "user-123", []string{})

	// Assert
	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_UserNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "user not found")

	suite.grpcClient.On("RemoveIntersectionIDs",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.RemoveIntersectionIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RemoveIntersectionIDs(
		ctx,
		"nonexistent",
		[]string{"intersection-1"},
	)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_IntersectionNotFound() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.NotFound, "intersection not found in user's list")

	suite.grpcClient.On("RemoveIntersectionIDs",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*user.RemoveIntersectionIDRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RemoveIntersectionIDs(
		ctx,
		"user-123",
		[]string{"nonexistent-intersection"},
	)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_EmptyUserID() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.InvalidArgument, "user ID is required")

	suite.grpcClient.On("RemoveIntersectionIDs",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RemoveIntersectionIDRequest) bool {
			return req.UserId == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.RemoveIntersectionIDs(ctx, "", []string{"intersection-1"})

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestClientIntersectionManagement(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
