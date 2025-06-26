//go:build integration
// +build integration

package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
)

func setupMongo(t *testing.T) (*mongo.Collection, func()) {
	ctx := context.Background()

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Use a test DB
	db := client.Database("integration_test_db")
	coll := db.Collection("intersections")

	// Clean up before and after test
	if err := coll.Drop(ctx); err != nil {
		t.Fatalf("Failed to clean up collection: %v", err)
	}

	cleanup := func() {
		_ = db.Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return coll, cleanup
}

func TestMongoIntersectionRepository_CreateIntersection_Integration(t *testing.T) {
	coll, cleanup := setupMongo(t)
	defer cleanup()

	repo := NewMongoIntersectionRepository(coll)
	ctx := context.Background()

	intersection := &model.IntersectionResponse{
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
	}

	result, err := repo.CreateIntersection(ctx, intersection)

	assert.NoError(t, err)
	assert.Equal(t, intersection.ID, result.ID)
	assert.Equal(t, intersection.Name, result.Name)
}
