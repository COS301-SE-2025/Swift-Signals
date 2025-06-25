package service

import (
	"context"
	"errors"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
	"time"
)

func TestService_CreateIntersection(t *testing.T) {
	tests := []struct {
		name           string
		inputName      string
		inputDetails   model.IntersectionDetails
		inputDensity   model.TrafficDensity
		inputParams    model.OptimisationParameters
		setupMock      func(*db.MockRepository)
		expectedError  string
		validateResult func(*testing.T, *model.IntersectionResponse)
	}{
		{
			name:      "successful creation",
			inputName: "Main Street & 1st Avenue",
			inputDetails: model.IntersectionDetails{
				Address:  "123 Main Street",
				City:     "Johannesburg",
				Province: "Gauteng",
			},
			inputDensity: model.TrafficHigh,
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptGridSearch,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionTrafficLight,
					Green:            30,
					Yellow:           3,
					Red:              25,
					Speed:            50,
					Seed:             12345,
				},
			},
			setupMock: func(mockRepo *db.MockRepository) {
				mockRepo.On("CreateIntersection", mock.Anything, mock.AnythingOfType("*model.IntersectionResponse")).
					Return(&model.IntersectionResponse{
						ID:   "generated-id",
						Name: "Main Street & 1st Avenue",
						Details: model.IntersectionDetails{
							Address:  "123 Main Street",
							City:     "Johannesburg",
							Province: "Gauteng",
						},
						CreatedAt:      time.Now(),
						LastRunAt:      time.Now(),
						Status:         model.Unoptimised,
						RunCount:       0,
						TrafficDensity: model.TrafficHigh,
						DefaultParameters: model.OptimisationParameters{
							OptimisationType: model.OptGridSearch,
							Parameters: model.SimulationParameters{
								IntersectionType: model.IntersectionTrafficLight,
								Green:            30,
								Yellow:           3,
								Red:              25,
								Speed:            50,
								Seed:             12345,
							},
						},
						BestParameters: model.OptimisationParameters{
							OptimisationType: model.OptGridSearch,
							Parameters: model.SimulationParameters{
								IntersectionType: model.IntersectionTrafficLight,
								Green:            30,
								Yellow:           3,
								Red:              25,
								Speed:            50,
								Seed:             12345,
							},
						},
						CurrentParameters: model.OptimisationParameters{
							OptimisationType: model.OptGridSearch,
							Parameters: model.SimulationParameters{
								IntersectionType: model.IntersectionTrafficLight,
								Green:            30,
								Yellow:           3,
								Red:              25,
								Speed:            50,
								Seed:             12345,
							},
						},
					}, nil)
			},
			expectedError: "",
			validateResult: func(t *testing.T, result *model.IntersectionResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, "Main Street & 1st Avenue", result.Name)
				assert.NotEmpty(t, result.ID)
				assert.Equal(t, model.Unoptimised, result.Status)
				assert.Equal(t, int32(0), result.RunCount)
				assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), result.LastRunAt, time.Second)
				assert.Equal(t, model.TrafficHigh, result.TrafficDensity)

				// Verify that all three parameter fields are set to the same value
				assert.Equal(t, result.DefaultParameters, result.BestParameters)
				assert.Equal(t, result.DefaultParameters, result.CurrentParameters)
				assert.Equal(t, model.OptGridSearch, result.DefaultParameters.OptimisationType)
			},
		},
		{
			name:      "repository error",
			inputName: "Test Intersection",
			inputDetails: model.IntersectionDetails{
				Address:  "456 Test Street",
				City:     "Cape Town",
				Province: "Western Cape",
			},
			inputDensity: model.TrafficMedium,
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptNone,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionRoundabout,
					Green:            20,
					Yellow:           2,
					Red:              15,
					Speed:            40,
					Seed:             54321,
				},
			},
			setupMock: func(mockRepo *db.MockRepository) {
				mockRepo.On("CreateIntersection", mock.Anything, mock.AnythingOfType("*model.IntersectionResponse")).
					Return(nil, errors.New("database connection failed"))
			},
			expectedError: "failed to create intersection: database connection failed",
			validateResult: func(t *testing.T, result *model.IntersectionResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "validation error - empty name",
			inputName: "",
			inputDetails: model.IntersectionDetails{
				Address:  "789 Empty Street",
				City:     "Durban",
				Province: "KwaZulu-Natal",
			},
			inputDensity: model.TrafficLow,
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptNone,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionStopSign,
					Green:            10,
					Yellow:           2,
					Red:              8,
					Speed:            30,
					Seed:             11111,
				},
			},
			setupMock: func(mockRepo *db.MockRepository) {
				// No repository call expected due to validation failure
			},
			expectedError: "validation failed: intersection name is required and cannot be empty",
			validateResult: func(t *testing.T, result *model.IntersectionResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "validation error - missing address",
			inputName: "Valid Name",
			inputDetails: model.IntersectionDetails{
				Address:  "",
				City:     "Pretoria",
				Province: "Gauteng",
			},
			inputDensity: model.TrafficMedium,
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptNone,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionTJunction,
					Green:            15,
					Yellow:           3,
					Red:              12,
					Speed:            35,
					Seed:             22222,
				},
			},
			setupMock: func(mockRepo *db.MockRepository) {
				// No repository call expected due to validation failure
			},
			expectedError: "validation failed: address is required",
			validateResult: func(t *testing.T, result *model.IntersectionResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "validation error - invalid timing parameters",
			inputName: "Invalid Timing",
			inputDetails: model.IntersectionDetails{
				Address:  "111 Invalid Street",
				City:     "Bloemfontein",
				Province: "Free State",
			},
			inputDensity: model.TrafficHigh,
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptGeneticEvaluation,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionTrafficLight,
					Green:            0, // Invalid - must be positive
					Yellow:           3,
					Red:              -5,  // Invalid - must be positive
					Speed:            500, // Invalid - exceeds max
					Seed:             33333,
				},
			},
			setupMock: func(mockRepo *db.MockRepository) {
				// No repository call expected due to validation failure
			},
			expectedError: "validation failed",
			validateResult: func(t *testing.T, result *model.IntersectionResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "validation error - multiple issues",
			inputName: strings.Repeat("a", 300), // Too long
			inputDetails: model.IntersectionDetails{
				Address:  "   ", // Empty after trim
				City:     "",    // Empty
				Province: "Valid Province",
			},
			inputDensity: model.TrafficDensity(999), // Invalid enum
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptimisationType(999), // Invalid enum
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionType(999), // Invalid enum
					Green:            -10,                         // Invalid
					Yellow:           0,                           // Invalid
					Red:              400,                         // Too high
					Speed:            -1,                          // Invalid
					Seed:             44444,
				},
			},
			setupMock: func(mockRepo *db.MockRepository) {
				// No repository call expected due to validation failure
			},
			expectedError: "validation failed",
			validateResult: func(t *testing.T, result *model.IntersectionResponse) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(db.MockRepository)
			tt.setupMock(mockRepo)

			service := &Service{
				repo: mockRepo,
			}

			ctx := context.Background()

			// Execute
			result, err := service.CreateIntersection(
				ctx,
				tt.inputName,
				tt.inputDetails,
				tt.inputDensity,
				tt.inputParams,
			)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			tt.validateResult(t, result)

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_CreateIntersection_ContextCancellation(t *testing.T) {
	// Test context cancellation
	mockRepo := new(db.MockRepository)
	service := &Service{
		repo: mockRepo,
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mockRepo.On("CreateIntersection", mock.Anything, mock.AnythingOfType("*model.IntersectionResponse")).
		Return(nil, context.Canceled)

	validDetails := model.IntersectionDetails{
		Address:  "123 Context Street",
		City:     "Johannesburg",
		Province: "Gauteng",
	}
	validParams := model.OptimisationParameters{
		OptimisationType: model.OptNone,
		Parameters: model.SimulationParameters{
			IntersectionType: model.IntersectionTrafficLight,
			Green:            30,
			Yellow:           3,
			Red:              25,
			Speed:            50,
			Seed:             12345,
		},
	}

	result, err := service.CreateIntersection(
		ctx,
		"Context Test",
		validDetails,
		model.TrafficMedium,
		validParams,
	)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create intersection")

	mockRepo.AssertExpectations(t)
}

// Test the validation helper function directly
func TestValidateCreateIntersectionInput(t *testing.T) {
	validDetails := model.IntersectionDetails{
		Address:  "123 Valid Street",
		City:     "Johannesburg",
		Province: "Gauteng",
	}
	validParams := model.OptimisationParameters{
		OptimisationType: model.OptGridSearch,
		Parameters: model.SimulationParameters{
			IntersectionType: model.IntersectionTrafficLight,
			Green:            30,
			Yellow:           3,
			Red:              25,
			Speed:            50,
			Seed:             12345,
		},
	}

	tests := []struct {
		name          string
		inputName     string
		inputDetails  model.IntersectionDetails
		inputDensity  model.TrafficDensity
		inputParams   model.OptimisationParameters
		expectError   bool
		errorContains string
	}{
		{
			name:         "valid input",
			inputName:    "Valid Intersection",
			inputDetails: validDetails,
			inputDensity: model.TrafficMedium,
			inputParams:  validParams,
			expectError:  false,
		},
		{
			name:          "empty name",
			inputName:     "",
			inputDetails:  validDetails,
			inputDensity:  model.TrafficMedium,
			inputParams:   validParams,
			expectError:   true,
			errorContains: "intersection name is required",
		},
		{
			name:          "name too long",
			inputName:     strings.Repeat("a", 300),
			inputDetails:  validDetails,
			inputDensity:  model.TrafficMedium,
			inputParams:   validParams,
			expectError:   true,
			errorContains: "cannot exceed 255 characters",
		},
		{
			name:          "invalid traffic density",
			inputName:     "Valid Name",
			inputDetails:  validDetails,
			inputDensity:  model.TrafficDensity(999),
			inputParams:   validParams,
			expectError:   true,
			errorContains: "invalid traffic density",
		},
		{
			name:         "negative green time",
			inputName:    "Valid Name",
			inputDetails: validDetails,
			inputDensity: model.TrafficMedium,
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptNone,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionTrafficLight,
					Green:            -5,
					Yellow:           3,
					Red:              25,
					Speed:            50,
					Seed:             12345,
				},
			},
			expectError:   true,
			errorContains: "green light duration must be positive",
		},
		{
			name:         "excessive speed",
			inputName:    "Valid Name",
			inputDetails: validDetails,
			inputDensity: model.TrafficMedium,
			inputParams: model.OptimisationParameters{
				OptimisationType: model.OptNone,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionTrafficLight,
					Green:            30,
					Yellow:           3,
					Red:              25,
					Speed:            500,
					Seed:             12345,
				},
			},
			expectError:   true,
			errorContains: "speed cannot exceed 200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateIntersectionInput(tt.inputName, tt.inputDetails, tt.inputDensity, tt.inputParams)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
