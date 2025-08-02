package test

import (
	"context"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestCreateUser_Success() {
	testUser := &model.User{
		ID:       uuid.NewString(),
		Name:     "Valid Name",
		Email:    "valid@gmail.com",
		Password: "8characters",
		IsAdmin:  false,
	}

	expectedQuery := `INSERT INTO users \(uuid, name, email, password, is_admin, created_at, updated_at\) 
	          VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\), NOW\(\)\)`

	suite.mock.ExpectExec(expectedQuery).
		WithArgs(testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.IsAdmin).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.NoError(err)
	suite.Equal(testUser, result)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateUser_Duplicate_Email() {
	testUser := &model.User{}

	suite.mock.ExpectExec("^INSERT (.+)").WillReturnError(&pq.Error{
		Code:   "23505",
		Table:  "users",
		Detail: `Key (email)=(duplicate@gmail.com) already exists.`,
	})

	ctx := context.Background()

	result, err := suite.repo.CreateUser(ctx, testUser)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal("email already exists", svcError.Message)
	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal(map[string]any{"email": "duplicate@gmail.com"}, svcError.Context)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateUser_Duplicate_UUID() {
	testUser := &model.User{}

	suite.mock.ExpectExec("^INSERT (.+)").WillReturnError(&pq.Error{
		Code:   "23505",
		Table:  "users",
		Detail: "uuid",
	})

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

func TestDBRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
