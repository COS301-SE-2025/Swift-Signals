package admin

import (
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestGetUserByID_Success() {
	userID := "user-123"

	expectedUser := createTestUser(
		userID,
		"John Doe",
		"john@example.com",
		false,
		[]string{"intersection-1", "intersection-2"},
	)

	ctx := createAdminContext()

	suite.client.On("GetUserByID", ctx, userID).
		Return(expectedUser, nil)

	result, err := suite.service.GetUserByID(ctx, userID)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal("John Doe", result.Username)
	suite.Equal("john@example.com", result.Email)
	suite.False(result.IsAdmin)
	suite.Equal([]string{"intersection-1", "intersection-2"}, result.IntersectionIDs)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_AdminUser() {
	userID := "admin-456"

	expectedUser := createTestUser(
		userID,
		"Jane Admin",
		"jane@admin.com",
		true,
		[]string{"intersection-1", "intersection-2", "intersection-3"},
	)

	ctx := createAdminContext()

	suite.client.On("GetUserByID", ctx, userID).
		Return(expectedUser, nil)

	result, err := suite.service.GetUserByID(ctx, userID)

	suite.Require().NoError(err)
	suite.Equal(userID, result.ID)
	suite.Equal("Jane Admin", result.Username)
	suite.Equal("jane@admin.com", result.Email)
	suite.True(result.IsAdmin)
	suite.Equal(
		[]string{"intersection-1", "intersection-2", "intersection-3"},
		result.IntersectionIDs,
	)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_Forbidden_NonAdmin() {
	userID := "user-123"
	ctx := createUserContext()

	result, err := suite.service.GetUserByID(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestGetUserByID_Forbidden_NoRole() {
	userID := "user-123"
	ctx := createContextWithoutRole()

	result, err := suite.service.GetUserByID(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("only admins can access this endpoint", svcError.Message)
}

func (suite *TestSuite) TestGetUserByID_UserNotFound() {
	userID := "nonexistent-user"
	ctx := createAdminContext()

	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{"userID": userID}))

	result, err := suite.service.GetUserByID(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_InternalError() {
	userID := "user-123"
	ctx := createAdminContext()

	suite.client.On("GetUserByID", ctx, userID).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	result, err := suite.service.GetUserByID(ctx, userID)

	suite.Require().Error(err)
	suite.Equal(model.User{}, result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}
