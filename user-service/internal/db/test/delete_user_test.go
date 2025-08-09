package test

import (
	"context"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestDeleteUser_Success() {
	suite.mock.ExpectExec(deleteUserQuery).
		WithArgs(testUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	err := suite.repo.DeleteUser(ctx, testUser.ID)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestDeleteUser_Success_UserNotExists() {
	suite.mock.ExpectExec(deleteUserQuery).
		WithArgs(testUser.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()

	err := suite.repo.DeleteUser(ctx, testUser.ID)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestDeleteUser_ForeignKeyConstraint() {
	pqError := &pq.Error{
		Code:   "23503",
		Detail: "Key (uuid)=(test-id) is still referenced from table \"user_intersections\"",
	}

	suite.mock.ExpectExec(deleteUserQuery).
		WithArgs(testUser.ID).
		WillReturnError(pqError)

	ctx := context.Background()

	err := suite.repo.DeleteUser(ctx, testUser.ID)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrDatabase, svcError.Code)
	suite.Equal("cannot delete record that is referenced by other records", svcError.Message)
	suite.Nil(svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestDeleteUser_InvalidParameterFormat() {
	pqError := &pq.Error{
		Code:   "22P02",
		Detail: "invalid input syntax for type uuid",
	}

	suite.mock.ExpectExec(deleteUserQuery).
		WithArgs(testUser.ID).
		WillReturnError(pqError)

	ctx := context.Background()

	err := suite.repo.DeleteUser(ctx, testUser.ID)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid query parameter format", svcError.Message)
	suite.Equal(map[string]any{"detail": "invalid input syntax for type uuid"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestDBDeleteUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
