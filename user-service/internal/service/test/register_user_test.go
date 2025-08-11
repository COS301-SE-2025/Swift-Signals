package test

import (
	"context"
	"errors"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestRegisterUser_Success() {
	email := "valid@gmail.com"
	name := "Valid Name"
	plainPassword := "8characters"

	suite.repo.On("GetUserByEmail", mock.Anything, email).
		Return(nil, nil)

	suite.repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.Name == name && u.Email == email && len(u.ID) > 0 && len(u.Password) > 0
	})).Return(func(_ context.Context, u *model.User) *model.User {
		return u
	}, nil)

	ctx := context.Background()

	result, err := suite.service.RegisterUser(ctx, name, email, plainPassword)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(name, result.Name)
	suite.Equal(email, result.Email)
	suite.NotEmpty(result.ID)
	suite.NotEmpty(result.Password)
	suite.NotEqual(plainPassword, result.Password)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Invalid_Input() {
	name := ""
	email := "noatsign"
	plainPassword := "<8"

	ctx := context.Background()

	result, err := suite.service.RegisterUser(ctx, name, email, plainPassword)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)

	expectedMessage := "invalid input"
	suite.Equal(expectedMessage, svcError.Message)

	expectedErrors := map[string]string{
		"email":    "Invalid email format",
		"name":     "Name is required",
		"password": "Password must be at least 8 characters long",
	}

	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Fail_To_Check_Existing_User() {
	name := "Valid Name"
	email := "valid@gmail.com"
	plainPassword := "password"

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, errors.New("any error"))

	ctx := context.Background()

	result, err := suite.service.RegisterUser(ctx, name, email, plainPassword)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)

	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to check existing user", svcError.Message)
	suite.Equal(map[string]any{"email": "valid@gmail.com"}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Email_Exists() {
	name := "Valid Name"
	email := "valid@gmail.com"
	plainPassword := "password"

	existingUser := &model.User{
		ID:   "Made up",
		Name: "Vaild Name",
	}

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(existingUser, nil)

	ctx := context.Background()
	result, err := suite.service.RegisterUser(ctx, name, email, plainPassword)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)

	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal("email already exists", svcError.Message)
	suite.Equal(map[string]any{"email": email, "user": existingUser}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_CreateUser_Error_Propagated() {
	email := "valid@gmail.com"
	name := "Valid Name"
	plainPassword := "8characters"

	suite.repo.On("GetUserByEmail", mock.Anything, email).
		Return(nil, nil)

	propError := errs.NewInternalError("testing error propagation", errors.New("any error"), nil)

	suite.repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, propError)

	ctx := context.Background()

	result, err := suite.service.RegisterUser(ctx, name, email, plainPassword)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)

	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("testing error propagation", svcError.Message)
	suite.Nil(svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Fail_To_Create_User() {
	email := "valid@gmail.com"
	name := "Valid Name"
	plainPassword := "8characters"
	anyError := errors.New("any error")

	suite.repo.On("GetUserByEmail", mock.Anything, email).
		Return(nil, nil)

	suite.repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, anyError)

	ctx := context.Background()

	result, err := suite.service.RegisterUser(ctx, name, email, plainPassword)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)

	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to create user", svcError.Message)

	expectedError := errs.NewInternalError("failed to create user", anyError, map[string]any{}).
		Error()

	suite.Equal(expectedError, svcError.Error())
	suite.Equal(map[string]any{}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func TestServiceRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
