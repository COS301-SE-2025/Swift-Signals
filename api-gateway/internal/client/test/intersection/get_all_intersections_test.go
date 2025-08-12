package intersection

import (
	"context"
	"testing"

	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestGetAllIntersections_Success() {
	// Arrange
	ctx := context.Background()

	mockStream := mocks.NewMockIntersectionService_GetAllIntersectionsClient[intersectionpb.IntersectionResponse](
		suite.T(),
	)

	suite.grpcClient.On("GetAllIntersections",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.GetAllIntersectionsRequest) bool {
			return req != nil
		})).Return(mockStream, nil)

	// Act
	stream, err := suite.client.GetAllIntersections(ctx)

	// Assert
	suite.Require().NoError(err)
	suite.NotNil(stream)
	suite.Equal(mockStream, stream)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_GrpcError() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.Internal, "internal server error")

	suite.grpcClient.On("GetAllIntersections",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.GetAllIntersectionsRequest) bool {
			return req != nil
		})).Return(nil, grpcErr)

	// Act
	stream, err := suite.client.GetAllIntersections(ctx)

	// Assert
	suite.Require().Error(err)
	suite.Nil(stream)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_Unauthorized() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.PermissionDenied, "access denied")

	suite.grpcClient.On("GetAllIntersections",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.GetAllIntersectionsRequest) bool {
			return req != nil
		})).Return(nil, grpcErr)

	// Act
	stream, err := suite.client.GetAllIntersections(ctx)

	// Assert
	suite.Require().Error(err)
	suite.Nil(stream)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ServiceUnavailable() {
	// Arrange
	ctx := context.Background()

	grpcErr := status.Error(codes.Unavailable, "service unavailable")

	suite.grpcClient.On("GetAllIntersections",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.GetAllIntersectionsRequest) bool {
			return req != nil
		})).Return(nil, grpcErr)

	// Act
	stream, err := suite.client.GetAllIntersections(ctx)

	// Assert
	suite.Require().Error(err)
	suite.Nil(stream)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ValidatesRequestStructure() {
	// Arrange
	ctx := context.Background()

	mockStream := mocks.NewMockIntersectionService_GetAllIntersectionsClient[intersectionpb.IntersectionResponse](
		suite.T(),
	)

	suite.grpcClient.On("GetAllIntersections",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.GetAllIntersectionsRequest) bool {
			return req != nil
		})).Return(mockStream, nil)

	// Act
	stream, err := suite.client.GetAllIntersections(ctx)

	// Assert
	suite.Require().NoError(err)
	suite.NotNil(stream)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ContextCancellation() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	grpcErr := status.Error(codes.Canceled, "context canceled")

	suite.grpcClient.On("GetAllIntersections",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.GetAllIntersectionsRequest) bool {
			return req != nil
		})).Return(nil, grpcErr)

	// Act
	stream, err := suite.client.GetAllIntersections(ctx)

	// Assert
	suite.Require().Error(err)
	suite.Nil(stream)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_NoTimeout() {
	// Test that GetAllIntersections doesn't set a timeout unlike other methods
	// Arrange
	ctx := context.Background()

	mockStream := mocks.NewMockIntersectionService_GetAllIntersectionsClient[intersectionpb.IntersectionResponse](
		suite.T(),
	)

	suite.grpcClient.On("GetAllIntersections",
		mock.MatchedBy(func(ctx context.Context) bool {
			// This method should not set a timeout, so deadline should not be set
			_, hasDeadline := ctx.Deadline()
			return !hasDeadline
		}),
		mock.AnythingOfType("*intersection.GetAllIntersectionsRequest")).Return(mockStream, nil)

	// Act
	stream, err := suite.client.GetAllIntersections(ctx)

	// Assert
	suite.Require().NoError(err)
	suite.NotNil(stream)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestIntersectionClientGetAllIntersections(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
