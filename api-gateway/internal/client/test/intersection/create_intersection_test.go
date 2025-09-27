package intersection

import (
	"context"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestCreateIntersection_Success() {
	// Arrange
	ctx := context.Background()

	intersection := model.Intersection{
		Name: "Test Intersection",
		Details: model.Details{
			Address:  "123 Test St",
			City:     "Test City",
			Province: "Test Province",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			OptimisationType: "grid_search",
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "t_junction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   "test-intersection-id",
		Name: "Test Intersection",
		Details: &intersectionpb.IntersectionDetails{
			Address:  "123 Test St",
			City:     "Test City",
			Province: "Test Province",
		},
		TrafficDensity: commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
		DefaultParameters: &commonpb.OptimisationParameters{
			OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.grpcClient.On("CreateIntersection",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *intersectionpb.CreateIntersectionRequest) bool {
			return req.Name == "Test Intersection" &&
				req.Details.Address == "123 Test St" &&
				req.Details.City == "Test City" &&
				req.Details.Province == "Test Province" &&
				req.TrafficDensity == commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH &&
				req.DefaultParameters.OptimisationType == commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH &&
				req.DefaultParameters.Parameters.IntersectionType == commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION &&
				req.DefaultParameters.Parameters.Green == 10 &&
				req.DefaultParameters.Parameters.Yellow == 3 &&
				req.DefaultParameters.Parameters.Red == 7 &&
				req.DefaultParameters.Parameters.Speed == 60 &&
				req.DefaultParameters.Parameters.Seed == 12345
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.CreateIntersection(ctx, intersection)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_GrpcError() {
	// Arrange
	ctx := context.Background()

	intersection := model.Intersection{
		Name: "Test Intersection",
		Details: model.Details{
			Address:  "123 Test St",
			City:     "Test City",
			Province: "Test Province",
		},
		TrafficDensity: "medium",
		DefaultParameters: model.OptimisationParameters{
			OptimisationType: "genetic_evaluation",
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "roundabout",
				Green:            15,
				Yellow:           2,
				Red:              8,
				Speed:            50,
				Seed:             54321,
			},
		},
	}

	grpcErr := status.Error(codes.AlreadyExists, "intersection already exists")

	suite.grpcClient.On("CreateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*intersection.CreateIntersectionRequest")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.CreateIntersection(ctx, intersection)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_ContextTimeout() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	intersection := model.Intersection{
		Name: "Test Intersection",
		Details: model.Details{
			Address:  "123 Test St",
			City:     "Test City",
			Province: "Test Province",
		},
		TrafficDensity: "low",
		DefaultParameters: model.OptimisationParameters{
			OptimisationType: "none",
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "stop_sign",
				Green:            5,
				Yellow:           1,
				Red:              3,
				Speed:            30,
				Seed:             98765,
			},
		},
	}

	suite.grpcClient.On("CreateIntersection",
		mock.Anything,
		mock.AnythingOfType("*intersection.CreateIntersectionRequest")).
		Return(nil, context.DeadlineExceeded)

	// Act
	result, err := suite.client.CreateIntersection(ctx, intersection)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_StringToEnumConversions() {
	// Test various string to enum conversions
	testCases := []struct {
		name                     string
		trafficDensity           string
		optimisationType         string
		intersectionType         string
		expectedTrafficDensity   commonpb.TrafficDensity
		expectedOptimisationType commonpb.OptimisationType
		expectedIntersectionType commonpb.IntersectionType
	}{
		{
			name:                     "High traffic density, grid search, traffic light",
			trafficDensity:           "high",
			optimisationType:         "grid_search",
			intersectionType:         "traffic_light",
			expectedTrafficDensity:   commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
			expectedOptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			expectedIntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TRAFFICLIGHT,
		},
		{
			name:                     "Medium traffic density, genetic, t-junction",
			trafficDensity:           "medium",
			optimisationType:         "genetic",
			intersectionType:         "t-junction",
			expectedTrafficDensity:   commonpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM,
			expectedOptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GENETIC_EVALUATION,
			expectedIntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
		},
		{
			name:                     "Low traffic density, none, roundabout",
			trafficDensity:           "low",
			optimisationType:         "none",
			intersectionType:         "roundabout",
			expectedTrafficDensity:   commonpb.TrafficDensity_TRAFFIC_DENSITY_LOW,
			expectedOptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_NONE,
			expectedIntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_ROUNDABOUT,
		},
		{
			name:                     "Invalid values default correctly",
			trafficDensity:           "invalid",
			optimisationType:         "invalid",
			intersectionType:         "invalid",
			expectedTrafficDensity:   commonpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM,
			expectedOptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			expectedIntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_UNSPECIFIED,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			ctx := context.Background()

			intersection := model.Intersection{
				Name: "Test Intersection",
				Details: model.Details{
					Address:  "123 Test St",
					City:     "Test City",
					Province: "Test Province",
				},
				TrafficDensity: tc.trafficDensity,
				DefaultParameters: model.OptimisationParameters{
					OptimisationType: tc.optimisationType,
					SimulationParameters: model.SimulationParameters{
						IntersectionType: tc.intersectionType,
						Green:            10,
						Yellow:           3,
						Red:              7,
						Speed:            60,
						Seed:             12345,
					},
				},
			}

			expectedResponse := &intersectionpb.IntersectionResponse{
				Id:   "test-intersection-id",
				Name: "Test Intersection",
			}

			suite.grpcClient.On("CreateIntersection",
				mock.AnythingOfType("*context.timerCtx"),
				mock.MatchedBy(func(req *intersectionpb.CreateIntersectionRequest) bool {
					return req.TrafficDensity == tc.expectedTrafficDensity &&
						req.DefaultParameters.OptimisationType == tc.expectedOptimisationType &&
						req.DefaultParameters.Parameters.IntersectionType == tc.expectedIntersectionType
				})).Return(expectedResponse, nil)

			// Act
			result, err := suite.client.CreateIntersection(ctx, intersection)

			// Assert
			suite.Require().NoError(err)
			suite.Equal(expectedResponse, result)
		})
	}
}

func (suite *TestSuite) TestCreateIntersection_EmptyFields() {
	// Arrange
	ctx := context.Background()

	intersection := model.Intersection{
		Name: "",
		Details: model.Details{
			Address:  "",
			City:     "",
			Province: "",
		},
		TrafficDensity: "",
		DefaultParameters: model.OptimisationParameters{
			OptimisationType: "",
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "",
				Green:            0,
				Yellow:           0,
				Red:              0,
				Speed:            0,
				Seed:             0,
			},
		},
	}

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   "test-intersection-id",
		Name: "",
	}

	suite.grpcClient.On("CreateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.CreateIntersectionRequest) bool {
			return req.Name == "" &&
				req.Details.Address == "" &&
				req.Details.City == "" &&
				req.Details.Province == "" &&
				req.DefaultParameters.Parameters.Green == 0 &&
				req.DefaultParameters.Parameters.Yellow == 0 &&
				req.DefaultParameters.Parameters.Red == 0 &&
				req.DefaultParameters.Parameters.Speed == 0 &&
				req.DefaultParameters.Parameters.Seed == 0
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.CreateIntersection(ctx, intersection)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestIntersectionClientCreateIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
