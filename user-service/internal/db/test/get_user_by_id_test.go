package test

import (
	"context"
	"database/sql"
	"errors"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/DATA-DOG/go-sqlmock"
)

func (suite *TestSuite) TestGetUserByID_Success() {
	userRows := sqlmock.NewRows([]string{"uuid", "name", "email", "password", "is_admin", "created_at", "updated_at"}).
		AddRow(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.CreatedAt, testUser.UpdatedAt)

	suite.mock.ExpectQuery(getUserByIDQuery).
		WithArgs(testUser.ID).
		WillReturnRows(userRows)

	intersectionRows := sqlmock.NewRows([]string{"intersection_id"})
	for _, intID := range testIntersectionIDs {
		intersectionRows.AddRow(intID)
	}

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUser.ID).
		WillReturnRows(intersectionRows)

	ctx := context.Background()
	result, err := suite.repo.GetUserByID(ctx, testUser.ID)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(testUser.ID, result.ID)
	suite.Equal(testUser.Name, result.Name)
	suite.Equal(testUser.Email, result.Email)
	suite.Equal(testUser.Password, result.Password)
	suite.Equal(testUser.IsAdmin, result.IsAdmin)
	suite.Equal(testIntersectionIDs, result.IntersectionIDs)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestGetUserByID_UserNotFound() {
	suite.mock.ExpectQuery(getUserByIDQuery).
		WithArgs(testUser.ID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	result, err := suite.repo.GetUserByID(ctx, testUser.ID)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("query execution failed", svcError.Message)
	suite.Nil(svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestGetUserByID() {
	userRows := sqlmock.NewRows([]string{"uuid", "name", "email", "password", "is_admin", "created_at", "updated_at"}).
		AddRow(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.CreatedAt, testUser.UpdatedAt)

	suite.mock.ExpectQuery(getUserByIDQuery).
		WithArgs(testUser.ID).
		WillReturnRows(userRows)

	expectedError := errors.New("database connection failed")

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUser.ID).
		WillReturnError(expectedError)

	ctx := context.Background()
	result, err := suite.repo.GetUserByID(ctx, testUser.ID)

	suite.Require().Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "failed to get intersections by user id")
	suite.NoError(suite.mock.ExpectationsWereMet())
}
