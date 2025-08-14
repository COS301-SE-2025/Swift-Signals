package test

import (
	"context"
	"database/sql"
	"errors"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func (suite *TestSuite) TestGetUserByEmail_Success() {
	userRows := sqlmock.NewRows([]string{"uuid", "name", "email", "password", "is_admin", "created_at", "updated_at"}).
		AddRow(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.CreatedAt, testUser.UpdatedAt)

	suite.mock.ExpectQuery(getUserByEmailQuery).
		WithArgs(testUser.Email).
		WillReturnRows(userRows)

	intersectionRows := sqlmock.NewRows([]string{"intersection_id"})
	for _, intID := range testIntersectionIDs {
		intersectionRows.AddRow(intID)
	}

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUser.ID).
		WillReturnRows(intersectionRows)

	ctx := context.Background()
	result, err := suite.repo.GetUserByEmail(ctx, testUser.Email)

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

func (suite *TestSuite) TestGetUserByEmail_UserNotFound() {
	suite.mock.ExpectQuery(getUserByEmailQuery).
		WithArgs(testUser.Email).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	result, err := suite.repo.GetUserByEmail(ctx, testUser.Email)

	suite.Require().NoError(err)
	suite.Nil(result)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestGetUserByEmail_DatabaseError_InvalidParameter() {
	pqError := &pq.Error{
		Code:   "22P02",
		Detail: "invalid input syntax for type uuid",
	}

	suite.mock.ExpectQuery(getUserByEmailQuery).
		WithArgs(testUser.Email).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.GetUserByEmail(ctx, testUser.Email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid query parameter format", svcError.Message)
	suite.Equal(map[string]any{"detail": "invalid input syntax for type uuid"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestGetUserByEmail_DatabaseError_UndefinedColumn() {
	pqError := &pq.Error{
		Code:   "42703",
		Column: "nonexistent_column",
	}

	suite.mock.ExpectQuery(getUserByEmailQuery).
		WithArgs(testUser.Email).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.GetUserByEmail(ctx, testUser.Email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("query references undefined column", svcError.Message)
	suite.Equal(map[string]any{"column": "nonexistent_column"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestGetUserByEmail_DatabaseError_UndefinedTable() {
	pqError := &pq.Error{
		Code: "42P01",
	}

	suite.mock.ExpectQuery(getUserByEmailQuery).
		WithArgs(testUser.Email).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.GetUserByEmail(ctx, testUser.Email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("query references undefined table", svcError.Message)
	suite.Nil(svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestGetUserByEmail_GetIntersectionsError() {
	userRows := sqlmock.NewRows([]string{"uuid", "name", "email", "password", "is_admin", "created_at", "updated_at"}).
		AddRow(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.CreatedAt, testUser.UpdatedAt)

	suite.mock.ExpectQuery(getUserByEmailQuery).
		WithArgs(testUser.Email).
		WillReturnRows(userRows)

	expectedError := errors.New("database connection failed")

	suite.mock.ExpectQuery(getIntersectionIDQuery).
		WithArgs(testUser.ID).
		WillReturnError(expectedError)

	ctx := context.Background()
	result, err := suite.repo.GetUserByEmail(ctx, testUser.Email)

	suite.Require().Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "failed to get intersections by user id")
	suite.NoError(suite.mock.ExpectationsWereMet())
}
