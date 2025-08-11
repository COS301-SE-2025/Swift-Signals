package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestAddIntersectionID_Success() {
	suite.mock.ExpectExec(addIntersectionIDQuery).
		WithArgs(testUser.ID, testIntersectionIDs[0]).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err := suite.repo.AddIntersectionID(ctx, testUser.ID, testIntersectionIDs[0])

	suite.Require().NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestAddIntersectionID_DatabaseError() {
	suite.mock.ExpectExec(addIntersectionIDQuery).
		WithArgs(testUser.ID, testIntersectionIDs[0]).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	err := suite.repo.AddIntersectionID(ctx, testUser.ID, testIntersectionIDs[0])

	suite.Require().Error(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestDBAddIntersectionID(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
