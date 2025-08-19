package intersection

import (
	"context"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestPutOptimisation_Success() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	parameters := model.OptimisationParameters{
		OptimisationType: "grid_search",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "t_junction",
			Green:            12,
			Yellow:           3,
			Red:              5,
			Speed:            60,
			Seed:             12345,
		},
	}

	expectedResponse := &intersectionpb.PutOptimisationResponse{
		Improved: true,
	}

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == intersectionID &&
				req.Parameters.OptimisationType == intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH &&
				req.Parameters.Parameters.IntersectionType == intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION &&
				req.Parameters.Parameters.Green == 12 &&
				req.Parameters.Parameters.Yellow == 3 &&
				req.Parameters.Parameters.Red == 5 &&
				req.Parameters.Parameters.Speed == 60 &&
				req.Parameters.Parameters.Seed == 12345
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_NotFound() {
	// Arrange
	ctx := context.Background()
	intersectionID := "non-existent-id"
	parameters := model.OptimisationParameters{
		OptimisationType: "genetic_evaluation",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "roundabout",
			Green:            15,
			Yellow:           2,
			Red:              8,
			Speed:            50,
			Seed:             54321,
		},
	}

	grpcErr := status.Error(codes.NotFound, "intersection not found")

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_EmptyID() {
	// Arrange
	ctx := context.Background()
	intersectionID := ""
	parameters := model.OptimisationParameters{
		OptimisationType: "none",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "stop_sign",
			Green:            5,
			Yellow:           1,
			Red:              3,
			Speed:            30,
			Seed:             98765,
		},
	}

	grpcErr := status.Error(codes.InvalidArgument, "intersection id is required")

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_Unauthorized() {
	// Arrange
	ctx := context.Background()
	intersectionID := "unauthorized-intersection-id"
	parameters := model.OptimisationParameters{
		OptimisationType: "grid_search",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "traffic_light",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	grpcErr := status.Error(codes.PermissionDenied, "access denied")

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_InternalError() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	parameters := model.OptimisationParameters{
		OptimisationType: "genetic_evaluation",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "t_junction",
			Green:            8,
			Yellow:           2,
			Red:              10,
			Speed:            45,
			Seed:             11111,
		},
	}

	grpcErr := status.Error(codes.Internal, "internal server error")

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_StringToEnumConversions() {
	// Test various string to enum conversions
	testCases := []struct {
		name                     string
		optimisationType         string
		intersectionType         string
		expectedOptimisationType intersectionpb.OptimisationType
		expectedIntersectionType intersectionpb.IntersectionType
	}{
		{
			name:                     "Grid search with traffic light",
			optimisationType:         "grid_search",
			intersectionType:         "traffic_light",
			expectedOptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			expectedIntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TRAFFICLIGHT,
		},
		{
			name:                     "Genetic with t-junction",
			optimisationType:         "genetic",
			intersectionType:         "t-junction",
			expectedOptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GENETIC_EVALUATION,
			expectedIntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
		},
		{
			name:                     "None with roundabout",
			optimisationType:         "none",
			intersectionType:         "roundabout",
			expectedOptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_NONE,
			expectedIntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_ROUNDABOUT,
		},
		{
			name:                     "Invalid values default correctly",
			optimisationType:         "invalid",
			intersectionType:         "invalid",
			expectedOptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			expectedIntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_UNSPECIFIED,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			ctx := context.Background()
			intersectionID := "test-intersection-id"
			parameters := model.OptimisationParameters{
				OptimisationType: tc.optimisationType,
				SimulationParameters: model.SimulationParameters{
					IntersectionType: tc.intersectionType,
					Green:            10,
					Yellow:           3,
					Red:              7,
					Speed:            60,
					Seed:             12345,
				},
			}

			expectedResponse := &intersectionpb.PutOptimisationResponse{
				Improved: true,
			}

			suite.grpcClient.On("PutOptimisation",
				mock.Anything,
				mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
					return req.Parameters.OptimisationType == tc.expectedOptimisationType &&
						req.Parameters.Parameters.IntersectionType == tc.expectedIntersectionType
				})).Return(expectedResponse, nil)

			// Act
			result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

			// Assert
			suite.Require().NoError(err)
			suite.Equal(expectedResponse, result)
		})
	}
}

func (suite *TestSuite) TestPutOptimisation_NoTimeout() {
	// Test that PutOptimisation doesn't set a timeout unlike other methods
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	parameters := model.OptimisationParameters{
		OptimisationType: "grid_search",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "t_junction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	expectedResponse := &intersectionpb.PutOptimisationResponse{
		Improved: false,
	}

	suite.grpcClient.On("PutOptimisation",
		mock.MatchedBy(func(ctx context.Context) bool {
			// This method should not set a timeout, so deadline should not be set
			_, hasDeadline := ctx.Deadline()
			return !hasDeadline
		}),
		mock.AnythingOfType("*intersection.PutOptimisationRequest")).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_ZeroValues() {
	// Test with zero values in simulation parameters
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	parameters := model.OptimisationParameters{
		OptimisationType: "",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "",
			Green:            0,
			Yellow:           0,
			Red:              0,
			Speed:            0,
			Seed:             0,
		},
	}

	expectedResponse := &intersectionpb.PutOptimisationResponse{
		Improved: false,
	}

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == intersectionID &&
				req.Parameters.Parameters.Green == 0 &&
				req.Parameters.Parameters.Yellow == 0 &&
				req.Parameters.Parameters.Red == 0 &&
				req.Parameters.Parameters.Speed == 0 &&
				req.Parameters.Parameters.Seed == 0
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_LargeValues() {
	// Test with large values in simulation parameters
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	parameters := model.OptimisationParameters{
		OptimisationType: "genetic_evaluation",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "traffic_light",
			Green:            999999,
			Yellow:           888888,
			Red:              777777,
			Speed:            666666,
			Seed:             555555,
		},
	}

	expectedResponse := &intersectionpb.PutOptimisationResponse{
		Improved: true,
	}

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == intersectionID &&
				req.Parameters.Parameters.Green == 999999 &&
				req.Parameters.Parameters.Yellow == 888888 &&
				req.Parameters.Parameters.Red == 777777 &&
				req.Parameters.Parameters.Speed == 666666 &&
				req.Parameters.Parameters.Seed == 555555
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestPutOptimisation_ConflictError() {
	// Test case where optimisation cannot be applied due to conflicts
	// Arrange
	ctx := context.Background()
	intersectionID := "intersection-being-optimized"
	parameters := model.OptimisationParameters{
		OptimisationType: "grid_search",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "t_junction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	grpcErr := status.Error(codes.FailedPrecondition, "optimization already in progress")

	suite.grpcClient.On("PutOptimisation",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.PutOptimisationRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.PutOptimisation(ctx, intersectionID, parameters)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestIntersectionClientPutOptimisation(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
