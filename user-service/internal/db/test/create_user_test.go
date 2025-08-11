package test

import (
	"context"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestCreateUser_Success() {
	suite.mock.ExpectExec(insertUserQuery).
		WithArgs(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.Require().NoError(err)
	suite.Equal(testUser, result)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateUser_Duplicate_Email() {
	pqError := &pq.Error{
		Code:   "23505",
		Table:  "users",
		Detail: `Key (email)=(test@gmail.com) already exists.`,
	}

	suite.mock.ExpectExec(insertUserQuery).
		WithArgs(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin).
		WillReturnError(pqError)

	ctx := context.Background()

	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("email already exists", svcError.Message)
	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal(map[string]any{"email": "test@gmail.com"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateUser_Duplicate_UUID() {
	pqError := &pq.Error{
		Code:   "23505",
		Table:  "users",
		Detail: "uuid",
	}
	suite.mock.ExpectExec(insertUserQuery).
		WithArgs(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin).
		WillReturnError(pqError)

	ctx := context.Background()

	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal("user ID already exists", svcError.Message)
	suite.Nil(svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateUser_Invalid_Reference() {
	pqError := &pq.Error{
		Code: "23503",
	}

	suite.mock.ExpectExec(insertUserQuery).
		WithArgs(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("invalid reference to related resource", svcError.Message)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateUser_Missing_Required_Field() {
	pqError := &pq.Error{
		Code:   "23502",
		Column: `column where field is missing`,
	}

	suite.mock.ExpectExec(insertUserQuery).
		WithArgs(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("missing required field", svcError.Message)
	suite.Equal(map[string]any{"column": "column where field is missing"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateUser_Field_Too_Long() {
	pqError := &pq.Error{
		Code:   "22001",
		Column: `column where field is too long`,
	}

	suite.mock.ExpectExec(insertUserQuery).
		WithArgs(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin).
		WillReturnError(pqError)

	ctx := context.Background()
	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("field value too long", svcError.Message)
	suite.Equal(map[string]any{"column": "column where field is too long"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestDBRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
