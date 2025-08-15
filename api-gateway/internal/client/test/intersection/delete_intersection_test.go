package intersection

import (
	"context"
	"testing"
	"time"

	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (suite *TestSuite) TestDeleteIntersection_Success() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("DeleteIntersection",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_NotFound() {
	// Arrange
	ctx := context.Background()
	intersectionID := "non-existent-id"

	grpcErr := status.Error(codes.NotFound, "intersection not found")

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_EmptyID() {
	// Arrange
	ctx := context.Background()
	intersectionID := ""

	grpcErr := status.Error(codes.InvalidArgument, "intersection id is required")

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_Unauthorized() {
	// Arrange
	ctx := context.Background()
	intersectionID := "unauthorized-intersection-id"

	grpcErr := status.Error(codes.PermissionDenied, "access denied")

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_InternalError() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"

	grpcErr := status.Error(codes.Internal, "internal server error")

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_ContextTimeout() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	intersectionID := "test-intersection-id"

	suite.grpcClient.On("DeleteIntersection",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, context.DeadlineExceeded)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_ConflictError() {
	// Test case where intersection cannot be deleted due to dependencies
	// Arrange
	ctx := context.Background()
	intersectionID := "intersection-with-dependencies"

	grpcErr := status.Error(
		codes.FailedPrecondition,
		"cannot delete intersection with active optimizations",
	)

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_ValidatesRequestStructure() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id-12345"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			// Validate that the request is properly structured
			return req != nil && req.Id == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_ServiceUnavailable() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"

	grpcErr := status.Error(codes.Unavailable, "service temporarily unavailable")

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_LongID() {
	// Test with a very long intersection ID
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id-with-very-long-name-that-might-cause-issues-in-the-system-and-should-be-handled-properly"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_SpecialCharactersInID() {
	// Test with special characters in intersection ID
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id-with-special-chars-@#$%^&*()"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_UUIDFormat() {
	// Test with UUID-formatted intersection ID
	// Arrange
	ctx := context.Background()
	intersectionID := "550e8400-e29b-41d4-a716-446655440000"

	expectedResponse := &emptypb.Empty{}

	suite.grpcClient.On("DeleteIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.DeleteIntersection(ctx, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestIntersectionClientDeleteIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
