package admin

import (
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestDeleteUserByID_Success() {
	userID := "user-123"
	ctx := createAdminContext()

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, nil)

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().NoError(err)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_Forbidden_NonAdmin() {
	userID := "user-123"
	ctx := createUserContext()

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestDeleteUserByID_Forbidden_NoRole() {
	userID := "user-123"
	ctx := createContextWithoutRole()

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestDeleteUserByID_UserNotFound() {
	userID := "nonexistent-user"
	ctx := createAdminContext()

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{"userID": userID}))

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_CannotDeleteSelf() {
	userID := "admin-self"
	ctx := createAdminContext()

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewForbiddenError("cannot delete your own account", map[string]any{"userID": userID}))

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("cannot delete your own account", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_InternalError() {
	userID := "user-123"
	ctx := createAdminContext()

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_EmptyUserID() {
	userID := ""
	ctx := createAdminContext()

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewValidationError("user ID cannot be empty", map[string]any{"userID": userID}))

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("user ID cannot be empty", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_InvalidUserID() {
	userID := "invalid-user-id-format"
	ctx := createAdminContext()

	suite.client.On("DeleteUser", ctx, userID).
		Return(nil, errs.NewValidationError("invalid user ID format", map[string]any{"userID": userID}))

	err := suite.service.DeleteUserByID(ctx, userID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid user ID format", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}
