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

// Successful update of current parameters
func TestUpdateCurrentParams_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful update", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: int32(1)}, // MatchedCount = 1
		))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		params := model.OptimisationParameters{} // use valid fields
		err := repo.UpdateCurrentParams(context.Background(), "1", params)

		require.NoError(t, err)
	})
}

// Intersection ID not found
func TestUpdateCurrentParams_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("intersection not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: int32(0)}, // MatchedCount = 0 triggers NotFoundError
		))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		params := model.OptimisationParameters{} // use valid fields
		err := repo.UpdateCurrentParams(context.Background(), "missing-id", params)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "intersection ID not found for optimisation update")
	})
}

// Database error occurs
func TestUpdateCurrentParams_DatabaseError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database error on update", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
			mtest.WriteError{Code: 123, Message: "mongo update failure"},
		))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		params := model.OptimisationParameters{} // use valid fields
		err := repo.UpdateCurrentParams(context.Background(), "1", params)

		require.Error(t, err)
	})
}
