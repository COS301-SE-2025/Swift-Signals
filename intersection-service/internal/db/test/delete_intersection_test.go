package test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/integration/mtest"

    db "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
)

// Successful deletion
func TestDeleteIntersection_Success(t *testing.T) {
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

    mt.Run("successful deletion", func(mt *mtest.T) {
        mt.AddMockResponses(mtest.CreateSuccessResponse(
            bson.E{"n", int32(1)},
        ))

        repo := db.NewMongoIntersectionRepo(mt.Coll)
        err := repo.DeleteIntersection(context.Background(), "1")

        assert.NoError(t, err)
    })
}

// Intersection ID not found
func TestDeleteIntersection_NotFound(t *testing.T) {
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

    mt.Run("intersection not found", func(mt *mtest.T) {
        // DeletedCount = 0 triggers NotFoundError
        mt.AddMockResponses(mtest.CreateSuccessResponse(
            bson.E{"n", int32(0)},
        ))

        repo := db.NewMongoIntersectionRepo(mt.Coll)
        err := repo.DeleteIntersection(context.Background(), "missing-id")

        assert.Error(t, err)
        assert.Contains(t, err.Error(), "intersection ID not found")
    })
}

// Database error occurs
func TestDeleteIntersection_DatabaseError(t *testing.T) {
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

    mt.Run("database error on delete", func(mt *mtest.T) {
        mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
            mtest.WriteError{Code: 123, Message: "mongo deletion failure"},
        ))

        repo := db.NewMongoIntersectionRepo(mt.Coll)
        err := repo.DeleteIntersection(context.Background(), "1")

        assert.Error(t, err)
    })
}
