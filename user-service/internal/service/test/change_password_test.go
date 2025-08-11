package test

import (
	"context"
	"errors"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

func (suite *TestSuite) TestChangePassword_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"

	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(currentPassword),
		bcrypt.DefaultCost,
	)

	existingUser := &model.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedCurrentPassword),
		IsAdmin:  false,
	}

	updatedUser := &model.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "new-hashed-password",
		IsAdmin:  false,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && u.Password != string(hashedCurrentPassword)
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_InvalidUserID() {
	userID := "invalid-uuid"
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID must be a valid UUID",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_EmptyUserID() {
	userID := ""
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_EmptyCurrentPassword() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := ""
	newPassword := "newpassword456"

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"currentpassword": "CurrentPassword is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_EmptyNewPassword() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := ""

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"newpassword": "NewPassword is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_NewPasswordTooShort() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "short"

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"newpassword": "NewPassword must be at least 8 characters long",
	}

	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_NewPasswordTooLong() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "this-is-a-very-long-password-that-exceeds-the-maximum-allowed-length-of-128-characters-and-should-trigger-validation-error-when-implemented-properly-with-validation-logic"

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"newpassword": "NewPassword must be at most 128 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_UserNotFound() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"

	suite.repo.On("GetUserByID", mock.Anything, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{}))

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_IncorrectCurrentPassword() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "wrongpassword"
	newPassword := "newpassword456"

	correctPassword := "correctpassword"
	hashedCorrectPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(correctPassword),
		bcrypt.DefaultCost,
	)

	existingUser := &model.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedCorrectPassword),
		IsAdmin:  false,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("current password is incorrect", svcError.Message)
	suite.Equal(map[string]any{"user": existingUser.PublicUser()}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_RepositoryGetUserError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_RepositoryUpdateError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"

	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(currentPassword),
		bcrypt.DefaultCost,
	)

	existingUser := &model.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedCurrentPassword),
		IsAdmin:  false,
	}

	repoError := errors.New("database update failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_RepositoryUpdateServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"

	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(currentPassword),
		bcrypt.DefaultCost,
	)

	existingUser := &model.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedCurrentPassword),
		IsAdmin:  false,
	}

	repoError := errs.NewUnauthorizedError("failed to update user", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_ServiceErrorPropagation() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "currentpassword123"
	newPassword := "newpassword456"
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_AdminUserPasswordChange() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	currentPassword := "admincurrentpass"
	newPassword := "adminnewpass123"

	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(currentPassword),
		bcrypt.DefaultCost,
	)

	existingUser := &model.User{
		ID:       userID,
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: string(hashedCurrentPassword),
		IsAdmin:  true,
	}

	updatedUser := &model.User{
		ID:       userID,
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: "new-hashed-password",
		IsAdmin:  true,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && u.IsAdmin && u.Password != string(hashedCurrentPassword)
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.ChangePassword(ctx, userID, currentPassword, newPassword)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func TestServiceChangePassword(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
