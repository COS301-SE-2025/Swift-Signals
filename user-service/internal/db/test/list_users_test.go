package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestListUsers_Success() {
	limit := 10
	offset := 0

	rows := sqlmock.NewRows([]string{
		"uuid", "name", "email", "password", "is_admin", "created_at", "updated_at",
	})

	for _, user := range testUsers {
		rows.AddRow(
			user.ID,
			user.Name,
			user.Email,
			user.Password,
			user.IsAdmin,
			user.CreatedAt,
			user.UpdatedAt,
		)
	}

	suite.mock.ExpectQuery(listUsersQuery).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	testIntersections1 := []string{"int-1", "int-2"}
	testIntersections2 := []string{"int-3"}

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUsers[0].ID).
		WillReturnRows(sqlmock.NewRows([]string{"intersection_id"}).
			AddRow("int-1").
			AddRow("int-2"))

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUsers[1].ID).
		WillReturnRows(sqlmock.NewRows([]string{"intersection_id"}).
			AddRow("int-3"))

	ctx := context.Background()
	result, err := suite.repo.ListUsers(ctx, limit, offset)

	suite.Require().NoError(err)
	suite.Len(result, 2)

	suite.Equal(testUsers[0].ID, result[0].ID)
	suite.Equal(testUsers[0].Name, result[0].Name)
	suite.Equal(testUsers[0].Email, result[0].Email)
	suite.Equal(testUsers[0].IsAdmin, result[0].IsAdmin)
	suite.Equal(testIntersections1, result[0].IntersectionIDs)

	suite.Equal(testUsers[1].ID, result[1].ID)
	suite.Equal(testUsers[1].Name, result[1].Name)
	suite.Equal(testUsers[1].Email, result[1].Email)
	suite.Equal(testUsers[1].IsAdmin, result[1].IsAdmin)
	suite.Equal(testIntersections2, result[1].IntersectionIDs)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestListUsers_QueryError() {
	limit := 10
	offset := 0

	suite.mock.ExpectQuery(listUsersQuery).
		WithArgs(limit, offset).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	result, err := suite.repo.ListUsers(ctx, limit, offset)

	suite.Require().Error(err)
	suite.Nil(result)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestListUsers_ScanError() {
	limit := 10
	offset := 0

	rows := sqlmock.NewRows([]string{
		"uuid", "name", "email", "password", "is_admin", "created_at", "updated_at",
	}).AddRow("test user", "Test Name", "test@gmail.com", "8characters", "invalid_boolean", time.Now(), time.Now())

	suite.mock.ExpectQuery(listUsersQuery).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := suite.repo.ListUsers(ctx, limit, offset)

	suite.Require().Error(err)
	suite.Nil(result)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestListUsers_GetIntersectionsError() {
	limit := 10
	offset := 0

	rows := sqlmock.NewRows([]string{
		"uuid", "name", "email", "password", "is_admin", "created_at", "updated_at",
	}).AddRow(
		testUser.ID,
		testUser.Name,
		testUser.Email,
		testUser.Password,
		testUser.IsAdmin,
		testUser.CreatedAt,
		testUser.UpdatedAt,
	)

	suite.mock.ExpectQuery(listUsersQuery).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUser.ID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	result, err := suite.repo.ListUsers(ctx, limit, offset)

	suite.Require().Error(err)
	suite.Nil(result)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestDBListUsers(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
