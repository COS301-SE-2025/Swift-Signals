package test

import (
	"context"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type TestSuite struct {
	suite.Suite
	mt   *mtest.T
	repo db.IntersectionRepository
	ctx  context.Context
}

func (suite *TestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mt = mtest.New(suite.T(), mtest.NewOptions().ClientType(mtest.Mock))

	suite.mt.Run("setup", func(mt *mtest.T) {
		suite.repo = db.NewMongoIntersectionRepo(mt.DB.Collection("intersections"))
	})
}

func TestDB(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
