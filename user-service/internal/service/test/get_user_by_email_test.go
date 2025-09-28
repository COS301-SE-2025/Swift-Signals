package test

import (
	"context"
	"errors"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestGetUserByEmail_Success() {
	email := "test@example.com"
	expectedUser := &model.User{
		ID:      "550e8400-e29b-41d4-a716-446655440000",
		Name:    "Test User",
		Email:   email,
		IsAdmin: false,
	}

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(expectedUser, nil)

	ctx := context.Background()

	result, err := suite.service.GetUserByEmail(ctx, email)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedUser.ID, result.ID)
	suite.Equal(expectedUser.Name, result.Name)
	suite.Equal(expectedUser.Email, result.Email)
	suite.Equal(expectedUser.IsAdmin, result.IsAdmin)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_UserNotFound() {
	email := "nonexistent@example.com"

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)

	ctx := context.Background()

	result, err := suite.service.GetUserByEmail(ctx, email)

	suite.Require().NoError(err)
	suite.Nil(result)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_InvalidEmailFormat() {
	email := "invalid-email-format"

	ctx := context.Background()

	result, err := suite.service.GetUserByEmail(ctx, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"email": "Invalid email format",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_EmptyEmail() {
	email := ""

	ctx := context.Background()

	result, err := suite.service.GetUserByEmail(ctx, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"email": "Email is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_EmailTooLong() {
	email := "thisisaverylongemailaddressthatexceedsthemaximumallowedlengthof255charactersandshouldtriggeravalidationerrorbutthisisstilltooshortsoneedtostillsaythingstomakeupfortheextracharactersiamnotcountingallofthesecharactersbutihopethatthisissufficient@verylongdomainname.com"

	ctx := context.Background()

	result, err := suite.service.GetUserByEmail(ctx, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"email": "Email must be at most 255 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_RepositoryError() {
	email := "test@example.com"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, repoError)

	ctx := context.Background()

	result, err := suite.service.GetUserByEmail(ctx, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)
	suite.Equal(map[string]any{"email": email}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_ServiceErrorPropagation() {
	email := "test@example.com"
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"email": email})

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, serviceError)

	ctx := context.Background()

	result, err := suite.service.GetUserByEmail(ctx, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}
