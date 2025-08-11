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

func (suite *TestSuite) TestMakeAdmin_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	existingUser := &model.User{
		ID:      userID,
		Name:    "Regular User",
		Email:   "user@example.com",
		IsAdmin: false,
	}

	updatedUser := &model.User{
		ID:      userID,
		Name:    "Regular User",
		Email:   "user@example.com",
		IsAdmin: true,
	}

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && u.IsAdmin == true
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_UserAlreadyAdmin() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	existingUser := &model.User{
		ID:      userID,
		Name:    "Already Admin User",
		Email:   "user@example.com",
		IsAdmin: true,
	}

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_InvalidUserID() {
	userID := "invalid-uuid"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

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

func (suite *TestSuite) TestMakeAdmin_InvalidAdminUserID() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "invalid-uuid"

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"adminuserid": "AdminUserID must be a valid UUID",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_EmptyUserID() {
	userID := ""
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

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

func (suite *TestSuite) TestMakeAdmin_EmptyAdminUserID() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := ""

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"adminuserid": "AdminUserID is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_AdminUserNotFound() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"
	noUserError := errs.NewNotFoundError("user not found", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(nil, noUserError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_AdminUserNotAdmin() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Regular User",
		Email:   "regular@example.com",
		IsAdmin: false,
	}

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("admin user does not have admin privileges", svcError.Message)
	suite.Equal(map[string]any{"adminUserID": adminUserID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_TargetUserNotFound() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"
	noUserError := errs.NewNotFoundError("user not found", map[string]any{})

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, noUserError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_RepositoryGetAdminUserError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find admin user", svcError.Message)
	suite.Equal(map[string]any{"adminUserID": adminUserID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_RepositoryGetTargetUserError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)
	suite.Equal(map[string]any{"userId": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_RepositoryUpdateError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	existingUser := &model.User{
		ID:      userID,
		Name:    "Regular User",
		Email:   "user@example.com",
		IsAdmin: false,
	}

	repoError := errors.New("database update failed")

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_RepositoryUpdateServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	existingUser := &model.User{
		ID:      userID,
		Name:    "Regular User",
		Email:   "user@example.com",
		IsAdmin: false,
	}

	repoError := errs.NewUnauthorizedError("failed to update user", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_ServiceErrorPropagationAdminUser() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"
	serviceError := errs.NewNotFoundError(
		"admin user not found",
		map[string]any{"admin_user_id": adminUserID},
	)

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("admin user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_ServiceErrorPropagationTargetUser() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	adminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      adminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	serviceError := errs.NewNotFoundError(
		"target user not found",
		map[string]any{"user_id": userID},
	)

	suite.repo.On("GetUserByID", mock.Anything, adminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("target user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_UserIDTrimming() {
	userID := "  550e8400-e29b-41d4-a716-446655440000  "
	adminUserID := "  550e8400-e29b-41d4-a716-446655440001  "
	trimmedUserID := "550e8400-e29b-41d4-a716-446655440000"
	trimmedAdminUserID := "550e8400-e29b-41d4-a716-446655440001"

	adminUser := &model.User{
		ID:      trimmedAdminUserID,
		Name:    "Admin User",
		Email:   "admin@example.com",
		IsAdmin: true,
	}

	existingUser := &model.User{
		ID:      trimmedUserID,
		Name:    "Regular User",
		Email:   "user@example.com",
		IsAdmin: false,
	}

	updatedUser := &model.User{
		ID:      trimmedUserID,
		Name:    "Regular User",
		Email:   "user@example.com",
		IsAdmin: true,
	}

	suite.repo.On("GetUserByID", mock.Anything, trimmedAdminUserID).Return(adminUser, nil)
	suite.repo.On("GetUserByID", mock.Anything, trimmedUserID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == trimmedUserID && u.IsAdmin == true
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.MakeAdmin(ctx, userID, adminUserID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func TestServiceMakeAdmin(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
