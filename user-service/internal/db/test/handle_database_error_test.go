package test

import (
	"errors"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestHandleDatabaseError_Unique_Constraint_Violation() {
	pqError := &pq.Error{
		Code:   "23505",
		Detail: "uncaught unique constraint violation",
	}

	ctx := db.ErrorContext{Operation: db.OpCreate}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("duplicate value violates unique constraint", svcError.Message)
	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal(map[string]any{"detail": "uncaught unique constraint violation"}, svcError.Context)
}

func (suite *TestSuite) TestHandlerDatabaseError_Conflict() {
	pqError := &pq.Error{
		Code:   "23505",
		Detail: `detail from conflict`,
	}

	ctx := db.ErrorContext{Operation: db.OpUpdate}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrDatabase, svcError.Code)
	suite.Equal("value conflicts with existing record", svcError.Message)
	suite.Equal(map[string]any{"detail": "detail from conflict"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestHandleDatabaseError_Err_Is_Nil() {
	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(nil, ctx)
	suite.NoError(err)
}

func (suite *TestSuite) TestHandleDatabaseError_Incorrect_Query() {
	pqError := errors.New("any error")
	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("query execution failed", svcError.Message)
	suite.Equal(errs.ErrInternal, svcError.Code)
}

func (suite *TestSuite) TestHandleDatabaseError_Connection_Lost() {
	pqError := &pq.Error{Code: "08003"}
	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("database connection lost", svcError.Message)
	suite.Equal(errs.ErrDatabase, svcError.Code)
}

func (suite *TestSuite) TestHandleDatabaseError_Connection_Limit_Reached() {
	pqError := &pq.Error{Code: "53300"}
	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("database connection limit reached", svcError.Message)
	suite.Equal(errs.ErrDatabase, svcError.Code)
}

func (suite *TestSuite) TestHandleDatabaseError_Configuration_Limit_Exceeded() {
	pqError := &pq.Error{Code: "53400"}
	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("database configuration limit exceeded", svcError.Message)
	suite.Equal(errs.ErrDatabase, svcError.Code)
}

func (suite *TestSuite) TestHandleDatabaseError_Query_Cancellation() {
	pqError := &pq.Error{Code: "57014"}
	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("query was canceled", svcError.Message)
	suite.Equal(errs.ErrDatabase, svcError.Code)
}

func (suite *TestSuite) TestHandleDatabaseError_Transaction_Conflict() {
	pqError := &pq.Error{Code: "40001"}
	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("transaction conflict, please retry", svcError.Message)
	suite.Equal(errs.ErrDatabase, svcError.Code)
}

func (suite *TestSuite) TestHandleDatabaseError_Uncaught_Code_and_Operation() {
	pqError := &pq.Error{
		Code:    "0",
		Message: "unknown postgres code and operation type",
		Detail:  "detail of unknown postgres error",
	}
	var (
		databaseOp        db.DatabaseOperation
		databaseErrorCode pq.ErrorCode
	)
	databaseErrorCode = "0"

	ctx := db.ErrorContext{}
	err := db.HandleDatabaseError(pqError, ctx)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("postgres error", svcError.Message)
	suite.Equal(errs.ErrInternal, svcError.Code)

	suite.Equal(
		map[string]any{
			"postgresErrCode":    databaseErrorCode,
			"postgresErrMessage": "unknown postgres code and operation type",
			"postgresErrDetail":  "detail of unknown postgres error",
			"operation":          databaseOp,
			"table":              "",
		},
		svcError.Context,
	)
}

func (suite *TestSuite) TestHandleDatabaseError_Uncaught_Code_Create_Operation() {
	pqError := &pq.Error{}
	ctx := db.ErrorContext{Operation: db.OpCreate}
	err := db.HandleDatabaseError(pqError, ctx)

	suite.NoError(err)
}

func (suite *TestSuite) TestHandleDatabaseError_Uncaught_Code_Read_Operation() {
	pqError := &pq.Error{}
	ctx := db.ErrorContext{Operation: db.OpRead}
	err := db.HandleDatabaseError(pqError, ctx)

	suite.NoError(err)
}

func (suite *TestSuite) TestHandleDatabaseError_Uncaught_Code_Update_Operation() {
	pqError := &pq.Error{}
	ctx := db.ErrorContext{Operation: db.OpUpdate}
	err := db.HandleDatabaseError(pqError, ctx)

	suite.NoError(err)
}

func (suite *TestSuite) TestHandleDatabaseError_Uncaught_Code_Delete_Operation() {
	pqError := &pq.Error{}
	ctx := db.ErrorContext{Operation: db.OpDelete}
	err := db.HandleDatabaseError(pqError, ctx)

	suite.NoError(err)
}

func TestDBHandleDatabaseError(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
