package admin

import (
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestUpdateUserByID_Success() {
	userID := "user-123"
	name := "Updated Name"
	email := "updated@example.com"

	expectedUser := createTestUser(
		userID,
		name,
		email,
		false,
		[]string{"intersection-1", "intersection-2"},
	)

	ctx := createAdminContext()

	suite.client.On("UpdateUser", ctx, userID, name, email).
		Return(expectedUser, nil)

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal(name, result.Username)
	suite.Equal(email, result.Email)
	suite.False(result.IsAdmin)
	suite.Equal([]string{"intersection-1", "intersection-2"}, result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUserByID_AdminUser() {
	userID := "admin-456"
	name := "Updated Admin Name"
	email := "updated.admin@example.com"

	expectedUser := createTestUser(
		userID,
		name,
		email,
		true,
		[]string{"intersection-1", "intersection-2", "intersection-3"},
	)

	ctx := createAdminContext()

	suite.client.On("UpdateUser", ctx, userID, name, email).
		Return(expectedUser, nil)

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal(name, result.Username)
	suite.Equal(email, result.Email)
	suite.True(result.IsAdmin)
	suite.Equal(
		[]string{"intersection-1", "intersection-2", "intersection-3"},
		result.IntersectionIDs,
	)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUserByID_Forbidden_NonAdmin() {
	userID := "user-123"
	name := "Updated Name"
	email := "updated@example.com"

	ctx := createUserContext()

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestUpdateUserByID_Forbidden_NoRole() {
	userID := "user-123"
	name := "Updated Name"
	email := "updated@example.com"

	ctx := createContextWithoutRole()

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestUpdateUserByID_UserNotFound() {
	userID := "nonexistent-user"
	name := "Valid Name"
	email := "valid@example.com"

	ctx := createAdminContext()

	suite.client.On("UpdateUser", ctx, userID, name, email).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{"userID": userID}))

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUserByID_ValidationError() {
	userID := "user-123"
	name := "Valid Name"
	email := "invalid-email-format"

	ctx := createAdminContext()

	suite.client.On("UpdateUser", ctx, userID, name, email).
		Return(nil, errs.NewValidationError("invalid email format", map[string]any{"email": email}))

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid email format", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUserByID_EmptyName() {
	userID := "user-123"
	name := ""
	email := "valid@example.com"

	ctx := createAdminContext()

	suite.client.On("UpdateUser", ctx, userID, name, email).
		Return(nil, errs.NewValidationError("name cannot be empty", map[string]any{"name": name}))

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("name cannot be empty", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUserByID_EmptyEmail() {
	userID := "user-123"
	name := "Valid Name"
	email := ""

	ctx := createAdminContext()

	suite.client.On("UpdateUser", ctx, userID, name, email).
		Return(nil, errs.NewValidationError("email cannot be empty", map[string]any{"email": email}))

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("email cannot be empty", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUserByID_InternalError() {
	userID := "user-123"
	name := "Valid Name"
	email := "valid@example.com"

	ctx := createAdminContext()

	suite.client.On("UpdateUser", ctx, userID, name, email).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	result, err := suite.service.UpdateUserByID(ctx, userID, name, email)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}
