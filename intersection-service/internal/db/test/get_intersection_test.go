package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	db "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
)

func TestFindIntersection(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("find intersection", func(mt *mtest.T) {
		expected := &model.Intersection{
			ID:   "123",
			Name: "Test Intersection",
		}

		first := mtest.CreateCursorResponse(1, "db.collection", mtest.FirstBatch, bson.D{
			{Key: "id", Value: expected.ID},
			{Key: "name", Value: expected.Name},
		})
		killCursors := mtest.CreateCursorResponse(0, "db.collection", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll) // ✅ use db.

		intersection, err := repo.GetIntersectionByID(context.Background(), "123")

		require.NoError(t, err)
		assert.Equal(t, expected.ID, intersection.ID)
		assert.Equal(t, expected.Name, intersection.Name)
	})
}

func TestFindIntersection_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("not found", func(mt *mtest.T) {
		// No documents returned
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "db.collection", mtest.FirstBatch))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		intersection, err := repo.GetIntersectionByID(context.Background(), "nonexistent")

		assert.Nil(t, intersection)
		assert.ErrorContains(t, err, "intersection ID not found")
	})
}

func TestFindIntersection_DatabaseError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    123,
			Message: "some mongo error",
		}))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		intersection, err := repo.GetIntersectionByID(context.Background(), "anyid")

		assert.Nil(t, intersection)
		assert.ErrorContains(t, err, "failed to find intersection")
	})
}
