package optimisation

import (
	"context"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestRunOptimisation_Success() {
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
		OptimisationType: "OPTIMISATION_TYPE_GRIDSEARCH",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "t_junction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	expectedResponse := &commonpb.OptimisationParameters{
		OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
		Parameters: &commonpb.SimulationParameters{
			IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
			Green:            12,
			Yellow:           3,
			Red:              5,
			Speed:            60,
			Seed:             12345,
		},
	}

	suite.grpcClient.On("RunOptimisation",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *commonpb.OptimisationParameters) bool {
			return req.Parameters.Green == 10 &&
				req.Parameters.Yellow == 3 &&
				req.Parameters.Red == 7 &&
				req.Parameters.Speed == 60 &&
				req.Parameters.Seed == 12345
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_InternalError() {
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
		OptimisationType: "OPTIMISATION_TYPE_GENETIC_EVALUATION",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "roundabout",
			Green:            15,
			Yellow:           2,
			Red:              8,
			Speed:            50,
			Seed:             54321,
		},
	}

	grpcErr := status.Error(codes.Internal, "optimization service error")

	suite.grpcClient.On("RunOptimisation",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*optimisation.OptimisationParameters")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_InvalidParameters() {
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
		OptimisationType: "INVALID_TYPE",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "invalid_type",
			Green:            -1,
			Yellow:           -1,
			Red:              -1,
			Speed:            -1,
			Seed:             -1,
		},
	}

	grpcErr := status.Error(codes.InvalidArgument, "invalid optimization parameters")

	suite.grpcClient.On("RunOptimisation",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*optimisation.OptimisationParameters")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_ContextTimeout() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	params := model.OptimisationParameters{
		OptimisationType: "OPTIMISATION_TYPE_GRIDSEARCH",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "traffic_light",
			Green:            8,
			Yellow:           2,
			Red:              10,
			Speed:            45,
			Seed:             98765,
		},
	}

	suite.grpcClient.On("RunOptimisation",
		mock.Anything,
		mock.AnythingOfType("*optimisation.OptimisationParameters")).
		Return(nil, context.DeadlineExceeded)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_ServiceUnavailable() {
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
		OptimisationType: "OPTIMISATION_TYPE_GRIDSEARCH",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "stop_sign",
			Green:            5,
			Yellow:           1,
			Red:              3,
			Speed:            30,
			Seed:             11111,
		},
	}

	grpcErr := status.Error(codes.Unavailable, "optimization service unavailable")

	suite.grpcClient.On("RunOptimisation",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*optimisation.OptimisationParameters")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_ZeroValues() {
	// Test with zero values in parameters
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
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

	expectedResponse := &commonpb.OptimisationParameters{
		OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
		Parameters: &commonpb.SimulationParameters{
			IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_UNSPECIFIED,
			Green:            0,
			Yellow:           0,
			Red:              0,
			Speed:            0,
			Seed:             0,
		},
	}

	suite.grpcClient.On("RunOptimisation",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *commonpb.OptimisationParameters) bool {
			return req.Parameters.Green == 0 &&
				req.Parameters.Yellow == 0 &&
				req.Parameters.Red == 0 &&
				req.Parameters.Speed == 0 &&
				req.Parameters.Seed == 0
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_LargeValues() {
	// Test with large values in parameters
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
		OptimisationType: "OPTIMISATION_TYPE_GENETIC_EVALUATION",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "traffic_light",
			Green:            999999,
			Yellow:           888888,
			Red:              777777,
			Speed:            666666,
			Seed:             555555,
		},
	}

	expectedResponse := &commonpb.OptimisationParameters{
		OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GENETIC_EVALUATION,
		Parameters: &commonpb.SimulationParameters{
			IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TRAFFICLIGHT,
			Green:            999999,
			Yellow:           888888,
			Red:              777777,
			Speed:            666666,
			Seed:             555555,
		},
	}

	suite.grpcClient.On("RunOptimisation",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *commonpb.OptimisationParameters) bool {
			return req.Parameters.Green == 999999 &&
				req.Parameters.Yellow == 888888 &&
				req.Parameters.Red == 777777 &&
				req.Parameters.Speed == 666666 &&
				req.Parameters.Seed == 555555
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_OptimizationTypeMapping() {
	// Test different optimization type values and their mapping
	testCases := []struct {
		name                     string
		inputOptimisationType    string
		expectedOptimisationType commonpb.OptimisationType
	}{
		{
			name:                     "OPTIMISATION_TYPE_GRIDSEARCH",
			inputOptimisationType:    "OPTIMISATION_TYPE_GRIDSEARCH",
			expectedOptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
		},
		{
			name:                     "OPTIMISATION_TYPE_GENETIC_EVALUATION",
			inputOptimisationType:    "OPTIMISATION_TYPE_GENETIC_EVALUATION",
			expectedOptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GENETIC_EVALUATION,
		},
		{
			name:                     "OPTIMISATION_TYPE_UNSPECIFIED",
			inputOptimisationType:    "OPTIMISATION_TYPE_UNSPECIFIED",
			expectedOptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_UNSPECIFIED,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			ctx := context.Background()

			params := model.OptimisationParameters{
				OptimisationType: tc.inputOptimisationType,
				SimulationParameters: model.SimulationParameters{
					IntersectionType: "t_junction",
					Green:            10,
					Yellow:           3,
					Red:              7,
					Speed:            60,
					Seed:             12345,
				},
			}

			expectedResponse := &commonpb.OptimisationParameters{
				OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
				Parameters: &commonpb.SimulationParameters{
					IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
					Green:            12,
					Yellow:           3,
					Red:              5,
					Speed:            60,
					Seed:             12345,
				},
			}

			// Verify that the client sends the correct optimization type in the request
			suite.grpcClient.On("RunOptimisation",
				mock.AnythingOfType("*context.timerCtx"),
				mock.MatchedBy(func(req *commonpb.OptimisationParameters) bool {
					return req.OptimisationType == tc.expectedOptimisationType
				})).Return(expectedResponse, nil)

			// Act
			result, err := suite.client.RunOptimisation(ctx, params)

			// Assert
			suite.Require().NoError(err)
			suite.Equal(expectedResponse, result)
		})
	}
}

func (suite *TestSuite) TestRunOptimisation_ResourceExhausted() {
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
		OptimisationType: "OPTIMISATION_TYPE_GRIDSEARCH",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "traffic_light",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	grpcErr := status.Error(codes.ResourceExhausted, "optimization queue full")

	suite.grpcClient.On("RunOptimisation",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*optimisation.OptimisationParameters")).
		Return(nil, grpcErr)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRunOptimisation_ValidatesRequestStructure() {
	// Arrange
	ctx := context.Background()

	params := model.OptimisationParameters{
		OptimisationType: "OPTIMISATION_TYPE_GRIDSEARCH",
		SimulationParameters: model.SimulationParameters{
			IntersectionType: "t_junction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	expectedResponse := &commonpb.OptimisationParameters{
		OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
		Parameters: &commonpb.SimulationParameters{
			IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
			Green:            12,
			Yellow:           3,
			Red:              5,
			Speed:            60,
			Seed:             12345,
		},
	}

	suite.grpcClient.On("RunOptimisation",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *commonpb.OptimisationParameters) bool {
			// Validate that the request is properly structured
			return req != nil &&
				req.Parameters != nil &&
				req.Parameters.Green == 10 &&
				req.Parameters.Yellow == 3 &&
				req.Parameters.Red == 7 &&
				req.Parameters.Speed == 60 &&
				req.Parameters.Seed == 12345
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.RunOptimisation(ctx, params)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestOptimisationClientRunOptimisation(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
