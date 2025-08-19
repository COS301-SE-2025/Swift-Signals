package test

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
)

func (suite *TestSuite) TestGetIntersectionsByUserID_Success() {
	rows := sqlmock.NewRows([]string{"intersection_id"})
	for _, id := range testIntersectionIDs {
		rows.AddRow(id)
	}

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUser.ID).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := suite.repo.GetIntersectionsByUserID(ctx, testUser.ID)

	suite.Require().NoError(err)
	suite.Equal(testIntersectionIDs, result)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestGetIntersectionsByUserID_ScanError() {
	rows := sqlmock.NewRows([]string{"intersection_id", "extra_column"}).
		AddRow("int-1", "extra-value")

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUser.ID).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := suite.repo.GetIntersectionsByUserID(ctx, testUser.ID)

	suite.Require().Error(err)
	suite.Nil(result)
	suite.NoError(suite.mock.ExpectationsWereMet())
}
