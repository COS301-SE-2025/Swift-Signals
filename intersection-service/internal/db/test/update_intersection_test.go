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

func TestUpdateIntersection_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful update", func(mt *mtest.T) {
		updated := model.Intersection{
			ID:   "1",
			Name: "Updated Name",
			Details: model.IntersectionDetails{
				Address:  "123 Main St",
				City:     "Testville",
				Province: "TestState",
			},
		}

		// Proper mock response
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "value", Value: bson.D{
				{Key: "id", Value: updated.ID},
				{Key: "name", Value: updated.Name},
				{Key: "details", Value: bson.D{
					{Key: "address", Value: updated.Details.Address},
					{Key: "city", Value: updated.Details.City},
					{Key: "province", Value: updated.Details.Province},
				}},
			}},
		))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		result, err := repo.UpdateIntersection(context.Background(), updated.ID, updated.Name, updated.Details)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updated.ID, result.ID)
		assert.Equal(t, updated.Name, result.Name)
		assert.Equal(t, updated.Details.Address, result.Details.Address)
		assert.Equal(t, updated.Details.City, result.Details.City)
		assert.Equal(t, updated.Details.Province, result.Details.Province)
	})
}

func TestUpdateIntersection_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("intersection not found", func(mt *mtest.T) {
		// Return empty response to trigger mongo.ErrNoDocuments
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{}))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		result, err := repo.UpdateIntersection(context.Background(), "missing-id", "name", model.IntersectionDetails{})

		assert.Nil(t, result)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "intersection ID not found") // ensures mongo.ErrNoDocuments branch executed
	})
}

func TestUpdateIntersection_DatabaseError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database error on update", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
			mtest.WriteError{Code: 123, Message: "mongo update failure"},
		))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		result, err := repo.UpdateIntersection(context.Background(), "1", "name", model.IntersectionDetails{})

		assert.Nil(t, result)
		require.Error(t, err)
	})
}

func TestUpdateIntersection_InvalidDecode(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("decode failure after update", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "decode failure",
		}))

		repo := db.NewMongoIntersectionRepo(mt.Coll)
		result, err := repo.UpdateIntersection(context.Background(), "1", "name", model.IntersectionDetails{})

		assert.Nil(t, result)
		require.Error(t, err)
	})
}
