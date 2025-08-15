package test

import (
	"context"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func (suite *TestSuite) TestUpdateUser_Success() {
	suite.mock.ExpectExec(updateUserQuery).
		WithArgs(testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	result, err := suite.repo.UpdateUser(ctx, testUser)

	suite.Require().NoError(err)
	suite.Equal(testUser, result)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestUpdateUser_Email_Already_In_Use() {
	pqError := &pq.Error{
		Code:   "23505",
		Table:  "users",
		Detail: `Key (email)=(test@gmail.com) already exists.`,
	}

	suite.mock.ExpectExec(updateUserQuery).
		WithArgs(testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.ID).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.UpdateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("email already in use by another user", svcError.Message)
	suite.Equal(map[string]any{"email": "test@gmail.com"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestUpdateUser_Record_Does_Not_Exist() {
	pqError := &pq.Error{
		Code:   "23503",
		Detail: `detail from error`,
	}

	suite.mock.ExpectExec(updateUserQuery).
		WithArgs(testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.ID).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.UpdateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("cannot update: referenced record does not exist", svcError.Message)
	suite.Equal(map[string]any{"detail": "detail from error"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestUpdateUser_Cannot_Set_To_Empty() {
	pqError := &pq.Error{
		Code:   "23502",
		Column: `column that cannot be empty`,
	}

	suite.mock.ExpectExec(updateUserQuery).
		WithArgs(testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.ID).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.UpdateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("cannot set required field to empty", svcError.Message)
	suite.Equal(map[string]any{"column": "column that cannot be empty"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestUpdateUser_Column_Does_Not_Exist() {
	pqError := &pq.Error{
		Code:   "42703",
		Column: `column that does not exist`,
	}

	suite.mock.ExpectExec(updateUserQuery).
		WithArgs(testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin, testUser.ID).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.UpdateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("cannot update: column does not exist", svcError.Message)
	suite.Equal(map[string]any{"column": "column that does not exist"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}
