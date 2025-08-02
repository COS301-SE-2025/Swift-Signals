package test

import (
	"database/sql"
	"log"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo db.UserRepository
}

func (suite *TestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	suite.Require().NoError(err)

	suite.repo = db.NewPostgresUserRepo(suite.db)
}

func (suite *TestSuite) TearDownTest() {
	if err := suite.db.Close(); err != nil {
		log.Printf("Failed to close suite: %v", err)
	}
}
