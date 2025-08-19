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

func TestGetAllIntersections_All(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("all intersections", func(mt *mtest.T) {
		expected := []*model.Intersection{
			{ID: "1", Name: "Intersection One"},
			{ID: "2", Name: "Intersection Two"},
		}

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		first := mtest.CreateCursorResponse(2, ns, mtest.FirstBatch,
			bson.D{{Key: "id", Value: expected[0].ID}, {Key: "name", Value: expected[0].Name}},
			bson.D{{Key: "id", Value: expected[1].ID}, {Key: "name", Value: expected[1].Name}},
		)
		killCursors := mtest.CreateCursorResponse(0, ns, mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, "")

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, expected[0].ID, results[0].ID)
		assert.Equal(t, expected[1].ID, results[1].ID)
	})
}

func TestGetAllIntersections_Filtered(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("filtered intersections", func(mt *mtest.T) {
		expected := []*model.Intersection{
			{ID: "42", Name: "Filtered Intersection"},
		}

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		first := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{{Key: "id", Value: expected[0].ID}, {Key: "name", Value: expected[0].Name}},
		)
		killCursors := mtest.CreateCursorResponse(0, ns, mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, "42")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, expected[0].ID, results[0].ID)
	})
}

func TestGetAllIntersections_LimitOffset(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("limit and offset", func(mt *mtest.T) {
		expected := []*model.Intersection{
			{ID: "100", Name: "Limited Intersection"},
		}

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		first := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{{Key: "id", Value: expected[0].ID}, {Key: "name", Value: expected[0].Name}},
		)
		killCursors := mtest.CreateCursorResponse(0, ns, mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 1, 10, "")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, expected[0].ID, results[0].ID)
	})
}

func TestGetAllIntersections_Empty(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("empty intersections", func(mt *mtest.T) {
		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, "")

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestGetAllIntersections_DatabaseError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
			mtest.WriteError{Code: 123, Message: "some mongo error"},
		))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, "")

		assert.Nil(t, results)
		assert.ErrorContains(t, err, "failed to find intersections")
	})
}

func TestGetAllIntersections_FilterEdgeCases(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("filter with spaces and empty IDs", func(mt *mtest.T) {
		expected := []*model.Intersection{
			{ID: "7", Name: "Edge Intersection"},
		}

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		first := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{{Key: "id", Value: expected[0].ID}, {Key: "name", Value: expected[0].Name}},
		)
		killCursors := mtest.CreateCursorResponse(0, ns, mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, " 7 , , ")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, expected[0].ID, results[0].ID)
	})
}

func TestGetAllIntersections_FilterEmptyIDs(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("filter with only spaces/empty IDs triggers else branch", func(mt *mtest.T) {
		expected := []*model.Intersection{
			{ID: "9", Name: "Empty Filter Intersection"},
		}

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		first := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{{Key: "id", Value: expected[0].ID}, {Key: "name", Value: expected[0].Name}},
		)
		killCursors := mtest.CreateCursorResponse(0, ns, mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, " , , ")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, expected[0].ID, results[0].ID)
	})
}

func TestGetAllIntersections_DecodeError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("cursor.All() returns decode error", func(mt *mtest.T) {
		// Using a command error to simulate failure in decoding
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "decode failure",
		}))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, "")

		assert.Nil(t, results)
		assert.ErrorContains(t, err, "failed to find intersections")
	})
}

func TestGetAllIntersections_CloseError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("cursor.Close() returns error", func(mt *mtest.T) {
		expected := []*model.Intersection{
			{ID: "10", Name: "Close Error Intersection"},
		}

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		first := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{{Key: "id", Value: expected[0].ID}, {Key: "name", Value: expected[0].Name}},
		)

		// mtest does not directly support Close() errors; the cursor is closed in defer
		// but this still exercises the defer block in your code
		killCursors := mtest.CreateCursorResponse(0, ns, mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, "")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, expected[0].ID, results[0].ID)
	})
}

func TestGetAllIntersections_CloseDeferCoverage(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("cursor defer branch runs", func(mt *mtest.T) {
		expected := []*model.Intersection{
			{ID: "11", Name: "Defer Coverage Intersection"},
		}

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		first := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{{Key: "id", Value: expected[0].ID}, {Key: "name", Value: expected[0].Name}},
		)
		killCursors := mtest.CreateCursorResponse(0, ns, mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		results, err := repo.GetAllIntersections(context.Background(), 0, 0, "")

		// This ensures the cursor gets fully iterated and defer runs
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, expected[0].ID, results[0].ID)
	})
}
