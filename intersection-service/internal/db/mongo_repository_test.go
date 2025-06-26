package db

import (
	"context"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoIntersectionRepository_CreateIntersection(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	tests := []struct {
		name          string
		setup         func(*mtest.T)
		ctx           context.Context
		intersection  *model.IntersectionResponse
		expectError   bool
		errorContains string
	}{
		{
			name: "success",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			ctx: context.Background(),
			intersection: &model.IntersectionResponse{
				ID:   "test-uuid-123",
				Name: "Test Intersection",
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
				}},
			expectError: false,
		},
		{
			name: "insert error",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   0,
					Code:    11000, // Duplicate key error
					Message: "duplicate key error",
				}))
			},
			ctx: context.Background(),
			intersection: &model.IntersectionResponse{
				ID:   "test-uuid-123",
				Name: "Test Intersection",
			},
			expectError:   true,
			errorContains: "duplicate key error",
		},
		{
			name: "context cancelled",
			setup: func(mt *mtest.T) {
				// No mock responses needed for context cancellation
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel context immediately
				return ctx
			}(),
			intersection: &model.IntersectionResponse{
				ID:   "test-uuid-123",
				Name: "Test Intersection",
			},
			expectError:   true,
			errorContains: "context canceled",
		},
	}

	for _, tc := range tests {
		mt.Run(tc.name, func(mt *mtest.T) {
			// Arrange
			repo := NewMongoIntersectionRepository(mt.DB.Collection("intersections"))
			tc.setup(mt)

			// Act
			result, err := repo.CreateIntersection(tc.ctx, tc.intersection)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.intersection, result)
			}
		})
	}
}

func TestMongoIntersectionRepository_GetIntersectionByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	tests := []struct {
		name          string
		id            string
		setup         func(mt *mtest.T)
		ctx           context.Context
		expectNil     bool
		expectError   bool
		errorContains string
	}{
		{
			name: "success with string ID",
			id:   "test-uuid-123",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.intersections", mtest.FirstBatch, bson.D{
					{Key: "id", Value: "test-uuid-123"},
					{Key: "name", Value: "Test Intersection"},
					{Key: "status", Value: model.Unoptimised},
					{Key: "runcount", Value: int32(0)},
				}))
			},
			ctx: context.Background(),
		},
		{
			name: "not found",
			id:   "non-existent-id",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.intersections", mtest.FirstBatch))
			},
			ctx:       context.Background(),
			expectNil: true,
		},
		{
			name: "database error",
			id:   "test-uuid-123",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    2,
					Message: "database connection error",
				}))
			},
			ctx:           context.Background(),
			expectError:   true,
			expectNil:     true,
			errorContains: "database connection error",
		},
		{
			name:  "context timeout",
			id:    "test-uuid-123",
			setup: func(mt *mtest.T) {}, // no mock needed
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(1 * time.Millisecond)
				defer cancel()
				return ctx
			}(),
			expectError:   true,
			expectNil:     true,
			errorContains: "context deadline exceeded",
		},
	}

	for _, tc := range tests {
		mt.Run(tc.name, func(mt *mtest.T) {
			repo := NewMongoIntersectionRepository(mt.DB.Collection("intersections"))
			tc.setup(mt)

			result, err := repo.GetIntersectionByID(tc.ctx, tc.id)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			if tc.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.id, result.ID)
			}
		})
	}
}

func TestMongoIntersectionRepository_GetAllIntersections(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	tests := []struct {
		name          string
		setup         func(mt *mtest.T)
		ctx           context.Context
		expectedLen   int
		expectNil     bool
		expectError   bool
		errorContains string
	}{
		{
			name: "success with multiple intersections",
			setup: func(mt *mtest.T) {
				first := mtest.CreateCursorResponse(1, "test.intersections", mtest.FirstBatch, bson.D{
					{Key: "id", Value: "test-uuid-1"},
					{Key: "name", Value: "Intersection 1"},
					{Key: "status", Value: model.Unoptimised},
					{Key: "runcount", Value: int32(0)},
				})
				second := mtest.CreateCursorResponse(2, "test.intersections", mtest.NextBatch, bson.D{
					{Key: "id", Value: "test-uuid-2"},
					{Key: "name", Value: "Intersection 2"},
					{Key: "status", Value: model.Optimised},
					{Key: "runcount", Value: int32(5)},
				})
				end := mtest.CreateCursorResponse(0, "test.intersections", mtest.NextBatch)
				mt.AddMockResponses(first, second, end)
			},
			ctx:         context.Background(),
			expectedLen: 2,
		},
		{
			name: "success with empty result",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.intersections", mtest.FirstBatch))
			},
			ctx:         context.Background(),
			expectedLen: 0,
		},
		{
			name: "cursor decode error",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.intersections", mtest.FirstBatch, bson.D{
					{Key: "id", Value: "test-uuid-1"},
					{Key: "invalid_field", Value: make(chan int)}, // invalid BSON type
				}))
			},
			ctx:         context.Background(),
			expectError: true,
			expectNil:   true,
		},
		{
			name: "find error",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    2,
					Message: "collection not found",
				}))
			},
			ctx:           context.Background(),
			expectError:   true,
			expectNil:     true,
			errorContains: "collection not found",
		},
		{
			name: "context cancelled during iteration",
			setup: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.intersections", mtest.FirstBatch, bson.D{
					{Key: "id", Value: "test-uuid-1"},
					{Key: "name", Value: "Intersection 1"},
				}))
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectError:   true,
			expectNil:     true,
			errorContains: "context canceled",
		},
	}

	for _, tc := range tests {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			repo := NewMongoIntersectionRepository(mt.DB.Collection("intersections"))
			tc.setup(mt)

			result, err := repo.GetAllIntersections(tc.ctx)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			if tc.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Len(t, result, tc.expectedLen)
			}
		})
	}
}
