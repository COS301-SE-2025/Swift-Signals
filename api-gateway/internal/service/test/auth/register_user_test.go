package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestRegisterUser_Success() {
	expectedRequest := model.RegisterRequest{
		Username: "Valid Name",
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	expectedResponse := model.RegisterResponse{
		UserID: "created-id",
	}

	suite.client.Mock.On("RegisterUser", mock.Anything, "Valid Name", "valid@gmail.com", "8characters").Return(&userpb.UserResponse{Id: "created-id"}, nil)

	ctx := context.Background()
	result, err := suite.service.RegisterUser(ctx, expectedRequest)

	suite.NoError(err)
	suite.Equal(expectedResponse, result)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Failure() {
	expectedRequest := model.RegisterRequest{
		Username: "Valid Name",
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	suite.client.Mock.On("RegisterUser", mock.Anything, "Valid Name", "valid@gmail.com", "8characters").Return(nil, errors.New("internal error caused"))

	ctx := context.Background()
	_, err := suite.service.RegisterUser(ctx, expectedRequest)

	svcError, ok := err.(*errs.ServiceError)

	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)

	expectedMessage := "unable to register"
	suite.Equal(expectedMessage, svcError.Message)

	expectedErrorMessage := "internal error caused"
	suite.Equal(expectedErrorMessage, svcError.Cause.Error())
	suite.client.AssertExpectations(suite.T())
}
func TestServiceRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
