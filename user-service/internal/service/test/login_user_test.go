package test

import (
	"context"
	"errors"
	"testing"
	"time"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

func (suite *TestSuite) TestLoginUser_Success_RegularUser() {
	email := "valid@gmail.com"
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	user := &model.User{
		ID:       "user-id-123",
		Name:     "Valid User",
		Email:    email,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)

	ctx := context.Background()

	token, expiryDate, err := suite.service.LoginUser(ctx, email, plainPassword)

	suite.Require().NoError(err)
	suite.NotEmpty(token)
	suite.True(expiryDate.After(time.Now()))
	suite.True(expiryDate.Before(time.Now().Add(time.Hour * 73)))

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_Success_AdminUser() {
	email := "admin@gmail.com"
	plainPassword := "adminpass123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	user := &model.User{
		ID:       "admin-id-123",
		Name:     "Admin User",
		Email:    email,
		Password: string(hashedPassword),
		IsAdmin:  true,
	}

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)

	ctx := context.Background()

	token, expiryDate, err := suite.service.LoginUser(ctx, email, plainPassword)

	suite.Require().NoError(err)
	suite.NotEmpty(token)
	suite.True(expiryDate.After(time.Now()))
	suite.True(expiryDate.Before(time.Now().Add(time.Hour * 73)))

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_Invalid_Input() {
	email := "invalid-email"
	plainPassword := ""

	ctx := context.Background()

	token, expiryDate, err := suite.service.LoginUser(ctx, email, plainPassword)

	suite.Empty(token)
	suite.Equal(time.Time{}, expiryDate)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"email":    "Invalid email format",
		"password": "Password is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_Repository_Error() {
	email := "valid@gmail.com"
	plainPassword := "password123"

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, errors.New("database error"))

	ctx := context.Background()

	token, expiryDate, err := suite.service.LoginUser(ctx, email, plainPassword)

	suite.Empty(token)
	suite.Equal(time.Time{}, expiryDate)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to check existing user", svcError.Message)
	suite.Equal(map[string]any{"email": email}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_User_Not_Found() {
	email := "nonexistent@gmail.com"
	plainPassword := "password123"

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)

	ctx := context.Background()

	token, expiryDate, err := suite.service.LoginUser(ctx, email, plainPassword)

	suite.Empty(token)
	suite.Equal(time.Time{}, expiryDate)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("user does not exist", svcError.Message)
	suite.Equal(map[string]any{"email": email}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_Incorrect_Password() {
	email := "valid@gmail.com"
	plainPassword := "wrongpassword"
	correctPassword := "correctpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	user := &model.User{
		ID:       "user-id-123",
		Name:     "Valid User",
		Email:    email,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)

	ctx := context.Background()

	token, expiryDate, err := suite.service.LoginUser(ctx, email, plainPassword)

	suite.Empty(token)
	suite.Equal(time.Time{}, expiryDate)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("password is incorrect", svcError.Message)
	suite.Equal(map[string]any{"user": user.PublicUser()}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_Email_Normalization() {
	email := "  VALID@Gmail.COM  "
	normalizedEmail := "valid@gmail.com"
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	user := &model.User{
		ID:       "user-id-123",
		Name:     "Valid User",
		Email:    normalizedEmail,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	suite.repo.On("GetUserByEmail", mock.Anything, normalizedEmail).Return(user, nil)

	ctx := context.Background()

	token, expiryDate, err := suite.service.LoginUser(ctx, email, plainPassword)

	suite.Require().NoError(err)
	suite.NotEmpty(token)
	suite.True(expiryDate.After(time.Now()))

	suite.repo.AssertExpectations(suite.T())
}

func TestServiceLoginUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
