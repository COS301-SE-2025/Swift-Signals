package auth

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestLoginUser_Success() {
	expectedRequest := model.LoginRequest{
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	expectedResponse := model.LoginResponse{
		Message: "Login Successful",
		Token:   "jwt.token.here",
	}

	suite.client.Mock.On("LoginUser", mock.Anything, "valid@gmail.com", "8characters").
		Return(&userpb.LoginUserResponse{Token: "jwt.token.here"}, nil)

	ctx := context.Background()
	result, err := suite.service.LoginUser(ctx, expectedRequest)

	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_InvalidCredentials() {
	expectedRequest := model.LoginRequest{
		Email:    "invalid@gmail.com",
		Password: "wrongpassword",
	}

	suite.client.Mock.On("LoginUser", mock.Anything, "invalid@gmail.com", "wrongpassword").
		Return(nil, errs.NewUnauthorizedError("invalid credentials", map[string]any{}))

	ctx := context.Background()
	_, err := suite.service.LoginUser(ctx, expectedRequest)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("invalid credentials", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_UserNotFound() {
	expectedRequest := model.LoginRequest{
		Email:    "notfound@gmail.com",
		Password: "8characters",
	}

	suite.client.Mock.On("LoginUser", mock.Anything, "notfound@gmail.com", "8characters").
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{}))

	ctx := context.Background()
	_, err := suite.service.LoginUser(ctx, expectedRequest)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_InternalError() {
	expectedRequest := model.LoginRequest{
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	suite.client.Mock.On("LoginUser", mock.Anything, "valid@gmail.com", "8characters").
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	ctx := context.Background()
	_, err := suite.service.LoginUser(ctx, expectedRequest)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_ValidationError() {
	expectedRequest := model.LoginRequest{
		Email:    "invalid-email",
		Password: "8characters",
	}

	suite.client.Mock.On("LoginUser", mock.Anything, "invalid-email", "8characters").
		Return(nil, errs.NewValidationError("invalid email format", map[string]any{}))

	ctx := context.Background()
	_, err := suite.service.LoginUser(ctx, expectedRequest)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid email format", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_EmptyResponse() {
	expectedRequest := model.LoginRequest{
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	expectedResponse := model.LoginResponse{
		Message: "Login Successful",
		Token:   "",
	}

	suite.client.Mock.On("LoginUser", mock.Anything, "valid@gmail.com", "8characters").
		Return(&userpb.LoginUserResponse{Token: ""}, nil)

	ctx := context.Background()
	result, err := suite.service.LoginUser(ctx, expectedRequest)

	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_AccountLocked() {
	expectedRequest := model.LoginRequest{
		Email:    "locked@gmail.com",
		Password: "8characters",
	}

	suite.client.Mock.On("LoginUser", mock.Anything, "locked@gmail.com", "8characters").
		Return(nil, errs.NewForbiddenError("account is locked", map[string]any{}))

	ctx := context.Background()
	_, err := suite.service.LoginUser(ctx, expectedRequest)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("account is locked", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}
