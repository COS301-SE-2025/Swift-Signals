package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongo(t *testing.T) (*mongo.Collection, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	dbName := fmt.Sprintf("integration_test_db_%d", time.Now().UnixNano())
	db := client.Database(dbName)
	coll := db.Collection("intersections")

	// Clean up before test
	if err := coll.Drop(ctx); err != nil {
		t.Fatalf("Failed to clean up collection: %v", err)
	}

	cleanup := func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()

		if err := db.Drop(cleanupCtx); err != nil {
			t.Logf("Warning: failed to drop test database: %v", err)
		}
		if err := client.Disconnect(cleanupCtx); err != nil {
			t.Logf("Warning: failed to disconnect from MongoDB: %v", err)
		}
	}

	return coll, cleanup
}

func TestMongoIntersectionRepository_CreateIntersection_Integration(t *testing.T) {
	testCases := []struct {
		name         string
		intersection *model.IntersectionResponse
		expectError  bool
	}{
		{
			name: "Valid intersection creation",
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
			},
			expectError: false,
		},
		{
			name: "Intersection with minimal data",
			intersection: &model.IntersectionResponse{
				ID:        "test-uuid-456",
				Name:      "Minimal Intersection",
				CreatedAt: time.Now(),
				Status:    model.Unoptimised,
				RunCount:  0,
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			coll, cleanup := setupMongo(t)
			defer cleanup()

			repo := NewMongoIntersectionRepository(coll)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create intersection
			result, err := repo.CreateIntersection(ctx, tc.intersection)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.intersection.ID, result.ID)
			assert.Equal(t, tc.intersection.Name, result.Name)
			assert.Equal(t, tc.intersection.Status, result.Status)
			assert.Equal(t, tc.intersection.RunCount, result.RunCount)

			// Verify it was actually saved to MongoDB by retrieving it
			saved, err := repo.GetIntersectionByID(ctx, tc.intersection.ID)
			assert.NoError(t, err)
			assert.NotNil(t, saved)
			assert.Equal(t, tc.intersection.ID, saved.ID)
			assert.Equal(t, tc.intersection.Name, saved.Name)
			assert.Equal(t, tc.intersection.Status, saved.Status)
			assert.Equal(t, tc.intersection.RunCount, saved.RunCount)

			// Verify specific fields if they were set
			if tc.intersection.Details.Address != "" {
				assert.Equal(t, tc.intersection.Details.Address, saved.Details.Address)
				assert.Equal(t, tc.intersection.Details.City, saved.Details.City)
				assert.Equal(t, tc.intersection.Details.Province, saved.Details.Province)
			}

			if tc.intersection.TrafficDensity.String() != "" {
				assert.Equal(t, tc.intersection.TrafficDensity, saved.TrafficDensity)
			}

			if tc.intersection.DefaultParameters.OptimisationType.String() != "" {
				assert.Equal(t, tc.intersection.DefaultParameters.OptimisationType, saved.DefaultParameters.OptimisationType)
				assert.Equal(t, tc.intersection.DefaultParameters.Parameters.Green, saved.DefaultParameters.Parameters.Green)
				assert.Equal(t, tc.intersection.DefaultParameters.Parameters.Yellow, saved.DefaultParameters.Parameters.Yellow)
				assert.Equal(t, tc.intersection.DefaultParameters.Parameters.Red, saved.DefaultParameters.Parameters.Red)
			}
		})
	}
}

func TestMongoIntersectionRepository_GetIntersectionByID_Integration(t *testing.T) {
	testCases := []struct {
		name        string
		setupData   *model.IntersectionResponse
		searchID    string
		expectFound bool
		expectError bool
	}{
		{
			name: "Get existing intersection",
			setupData: &model.IntersectionResponse{
				ID:   "existing-intersection-123",
				Name: "Existing Test Intersection",
				Details: model.IntersectionDetails{
					Address:  "456 Test Avenue",
					City:     "Cape Town",
					Province: "Western Cape",
				},
				CreatedAt:      time.Now(),
				LastRunAt:      time.Now(),
				Status:         model.Unoptimised,
				RunCount:       5,
				TrafficDensity: model.TrafficMedium,
				DefaultParameters: model.OptimisationParameters{
					OptimisationType: model.OptGridSearch,
					Parameters: model.SimulationParameters{
						IntersectionType: model.IntersectionTrafficLight,
						Green:            25,
						Yellow:           4,
						Red:              30,
						Speed:            60,
						Seed:             67890,
					},
				},
			},
			searchID:    "existing-intersection-123",
			expectFound: true,
			expectError: false,
		},
		{
			name:        "Get non-existent intersection",
			setupData:   nil,
			searchID:    "non-existent-id",
			expectFound: false,
			expectError: false,
		},
		{
			name: "Get intersection with minimal data",
			setupData: &model.IntersectionResponse{
				ID:        "minimal-intersection-456",
				Name:      "Minimal Intersection",
				CreatedAt: time.Now(),
				Status:    model.Unoptimised,
				RunCount:  0,
			},
			searchID:    "minimal-intersection-456",
			expectFound: true,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			coll, cleanup := setupMongo(t)
			defer cleanup()

			repo := NewMongoIntersectionRepository(coll)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Setup test data if provided
			if tc.setupData != nil {
				_, err := repo.CreateIntersection(ctx, tc.setupData)
				assert.NoError(t, err)
			}

			// Test GetIntersectionByID
			result, err := repo.GetIntersectionByID(ctx, tc.searchID)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tc.expectFound {
				assert.NotNil(t, result)
				assert.Equal(t, tc.setupData.ID, result.ID)
				assert.Equal(t, tc.setupData.Name, result.Name)
				assert.Equal(t, tc.setupData.Status, result.Status)
				assert.Equal(t, tc.setupData.RunCount, result.RunCount)

				// Verify detailed fields if they exist
				if tc.setupData.Details.Address != "" {
					assert.Equal(t, tc.setupData.Details.Address, result.Details.Address)
					assert.Equal(t, tc.setupData.Details.City, result.Details.City)
					assert.Equal(t, tc.setupData.Details.Province, result.Details.Province)
				}

				if tc.setupData.TrafficDensity.String() != "" {
					assert.Equal(t, tc.setupData.TrafficDensity, result.TrafficDensity)
				}
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestMongoIntersectionRepository_GetAllIntersections_Integration(t *testing.T) {
	testCases := []struct {
		name          string
		setupData     []*model.IntersectionResponse
		expectedCount int
	}{
		{
			name:          "Get all when collection is empty",
			setupData:     []*model.IntersectionResponse{},
			expectedCount: 0,
		},
		{
			name: "Get all with single intersection",
			setupData: []*model.IntersectionResponse{
				{
					ID:        "single-intersection-1",
					Name:      "Single Test Intersection",
					CreatedAt: time.Now(),
					Status:    model.Unoptimised,
					RunCount:  0,
				},
			},
			expectedCount: 1,
		},
		{
			name: "Get all with multiple intersections",
			setupData: []*model.IntersectionResponse{
				{
					ID:   "multi-intersection-1",
					Name: "First Intersection",
					Details: model.IntersectionDetails{
						Address:  "123 First Street",
						City:     "Johannesburg",
						Province: "Gauteng",
					},
					CreatedAt:      time.Now(),
					Status:         model.Unoptimised,
					RunCount:       0,
					TrafficDensity: model.TrafficHigh,
				},
				{
					ID:   "multi-intersection-2",
					Name: "Second Intersection",
					Details: model.IntersectionDetails{
						Address:  "456 Second Avenue",
						City:     "Durban",
						Province: "KwaZulu-Natal",
					},
					CreatedAt:      time.Now(),
					Status:         model.Optimised,
					RunCount:       3,
					TrafficDensity: model.TrafficMedium,
				},
				{
					ID:        "multi-intersection-3",
					Name:      "Third Intersection",
					CreatedAt: time.Now(),
					Status:    model.Unoptimised,
					RunCount:  1,
				},
			},
			expectedCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			coll, cleanup := setupMongo(t)
			defer cleanup()

			repo := NewMongoIntersectionRepository(coll)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Setup test data
			createdIDs := make(map[string]bool)
			for _, intersection := range tc.setupData {
				_, err := repo.CreateIntersection(ctx, intersection)
				assert.NoError(t, err)
				createdIDs[intersection.ID] = true
			}

			// Test GetAllIntersections
			results, err := repo.GetAllIntersections(ctx)
			assert.NoError(t, err)
			assert.Len(t, results, tc.expectedCount)

			// Verify all created intersections are returned
			if tc.expectedCount > 0 {
				returnedIDs := make(map[string]bool)
				for _, result := range results {
					returnedIDs[result.ID] = true

					// Verify the intersection exists in our setup data
					assert.True(t, createdIDs[result.ID], "Returned intersection should exist in setup data")
				}

				// Verify all setup intersections are returned
				assert.Equal(t, len(createdIDs), len(returnedIDs), "All created intersections should be returned")
				for id := range createdIDs {
					assert.True(t, returnedIDs[id], "Created intersection %s should be in results", id)
				}
			}
		})
	}
}

func TestMongoIntersectionRepository_GetAllIntersections_VerifyData_Integration(t *testing.T) {
	coll, cleanup := setupMongo(t)
	defer cleanup()

	repo := NewMongoIntersectionRepository(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create test intersections with different data completeness
	testIntersections := []*model.IntersectionResponse{
		{
			ID:   "full-data-intersection",
			Name: "Full Data Intersection",
			Details: model.IntersectionDetails{
				Address:  "789 Full Data Street",
				City:     "Pretoria",
				Province: "Gauteng",
			},
			CreatedAt:      time.Now(),
			LastRunAt:      time.Now(),
			Status:         model.Optimised,
			RunCount:       10,
			TrafficDensity: model.TrafficLow,
			DefaultParameters: model.OptimisationParameters{
				OptimisationType: model.OptGridSearch,
				Parameters: model.SimulationParameters{
					IntersectionType: model.IntersectionTrafficLight,
					Green:            35,
					Yellow:           5,
					Red:              20,
					Speed:            40,
					Seed:             11111,
				},
			},
		},
		{
			ID:        "minimal-data-intersection",
			Name:      "Minimal Data Intersection",
			CreatedAt: time.Now(),
			Status:    model.Unoptimised,
			RunCount:  0,
		},
	}

	// Setup data
	for _, intersection := range testIntersections {
		_, err := repo.CreateIntersection(ctx, intersection)
		assert.NoError(t, err)
	}

	// Get all intersections
	results, err := repo.GetAllIntersections(ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	resultMap := make(map[string]*model.IntersectionResponse)
	for _, result := range results {
		resultMap[result.ID] = result
	}

	// Verify full data intersection
	fullDataResult := resultMap["full-data-intersection"]
	assert.NotNil(t, fullDataResult)
	assert.Equal(t, "Full Data Intersection", fullDataResult.Name)
	assert.Equal(t, "789 Full Data Street", fullDataResult.Details.Address)
	assert.Equal(t, "Pretoria", fullDataResult.Details.City)
	assert.Equal(t, "Gauteng", fullDataResult.Details.Province)
	assert.Equal(t, model.Optimised, fullDataResult.Status)
	assert.Equal(t, int32(10), fullDataResult.RunCount)
	assert.Equal(t, model.TrafficLow, fullDataResult.TrafficDensity)
	assert.Equal(t, model.OptGridSearch, fullDataResult.DefaultParameters.OptimisationType)
	assert.Equal(t, int32(35), fullDataResult.DefaultParameters.Parameters.Green)

	// Verify minimal data intersection
	minimalDataResult := resultMap["minimal-data-intersection"]
	assert.NotNil(t, minimalDataResult)
	assert.Equal(t, "Minimal Data Intersection", minimalDataResult.Name)
	assert.Equal(t, model.Unoptimised, minimalDataResult.Status)
	assert.Equal(t, int32(0), minimalDataResult.RunCount)
}
