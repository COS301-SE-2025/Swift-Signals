package test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/mongo/integration/mtest"

    db "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
    "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
)

func TestCreateIntersection(t *testing.T) {
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

    mt.Run("insert intersection", func(mt *mtest.T) {
        mt.AddMockResponses(mtest.CreateSuccessResponse())

        repo := db.NewMongoIntersectionRepo(mt.Coll) // ✅ use db.

        intersection := &model.Intersection{
            ID:   "123",
            Name: "Test Intersection",
        }

        result, err := repo.CreateIntersection(context.Background(), intersection)

        assert.NoError(t, err)
        assert.Equal(t, "123", result.ID)
    })
}

func TestCreateIntersection_Error(t *testing.T) {
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

    mt.Run("insert intersection error", func(mt *mtest.T) {
        // Simulate an insert error
        mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
            mtest.WriteError{
                Index:   0,
                Code:    12345,
                Message: "simulated insert error",
            },
        ))

        repo := db.NewMongoIntersectionRepo(mt.Coll)

        intersection := &model.Intersection{
            ID:   "999",
            Name: "Error Intersection",
        }

        result, err := repo.CreateIntersection(context.Background(), intersection)

        // Assert that error is returned and result is nil
        assert.Nil(t, result)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "failed to insert intersection")
    })
}
