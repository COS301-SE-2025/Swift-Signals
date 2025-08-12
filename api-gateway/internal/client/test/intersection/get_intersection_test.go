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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (suite *TestSuite) TestGetIntersection_Success() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   intersectionID,
		Name: "Test Intersection",
		Details: &intersectionpb.IntersectionDetails{
			Address:  "123 Test St",
			City:     "Test City",
			Province: "Test Province",
		},
		CreatedAt:      timestamppb.Now(),
		LastRunAt:      timestamppb.Now(),
		Status:         intersectionpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		RunCount:       5,
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
		DefaultParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
		BestParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            12,
				Yellow:           3,
				Red:              5,
				Speed:            60,
				Seed:             12345,
			},
		},
		CurrentParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            11,
				Yellow:           3,
				Red:              6,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.grpcClient.On("GetIntersection",
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
	result, err := suite.client.GetIntersection(ctx, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_NotFound() {
	// Arrange
	ctx := context.Background()
	intersectionID := "non-existent-id"

	grpcErr := status.Error(codes.NotFound, "intersection not found")

	suite.grpcClient.On("GetIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_EmptyID() {
	// Arrange
	ctx := context.Background()
	intersectionID := ""

	grpcErr := status.Error(codes.InvalidArgument, "intersection id is required")

	suite.grpcClient.On("GetIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_InternalError() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"

	grpcErr := status.Error(codes.Internal, "internal server error")

	suite.grpcClient.On("GetIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_ContextTimeout() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	intersectionID := "test-intersection-id"

	suite.grpcClient.On("GetIntersection",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, context.DeadlineExceeded)

	// Act
	result, err := suite.client.GetIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_ValidatesRequestStructure() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id-12345"

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:             intersectionID,
		Name:           "Complex Test Intersection",
		Status:         intersectionpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISED,
		RunCount:       10,
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM,
	}

	suite.grpcClient.On("GetIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			// Validate that the request is properly structured
			return req != nil && req.Id == intersectionID
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.GetIntersection(ctx, intersectionID)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(intersectionID, result.Id)
	suite.Equal("Complex Test Intersection", result.Name)
	suite.Equal(intersectionpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISED, result.Status)
	suite.Equal(int32(10), result.RunCount)
	suite.Equal(intersectionpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM, result.TrafficDensity)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_Unauthorized() {
	// Arrange
	ctx := context.Background()
	intersectionID := "unauthorized-intersection-id"

	grpcErr := status.Error(codes.PermissionDenied, "access denied")

	suite.grpcClient.On("GetIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.IntersectionIDRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.GetIntersection(ctx, intersectionID)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestIntersectionClientGetIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
