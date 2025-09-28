package auth

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
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

	suite.client.Mock.On("RegisterUser", mock.Anything, "Valid Name", "valid@gmail.com", "8characters").
		Return(&userpb.UserResponse{Id: "created-id"}, nil)

	ctx := context.Background()
	result, err := suite.service.RegisterUser(ctx, expectedRequest)

	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)

	suite.client.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Failure() {
	expectedRequest := model.RegisterRequest{
		Username: "Valid Name",
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	suite.client.Mock.On("RegisterUser", mock.Anything, "Valid Name", "valid@gmail.com", "8characters").
		Return(nil, errs.NewAlreadyExistsError("user already exists", map[string]any{}))

	ctx := context.Background()
	_, err := suite.service.RegisterUser(ctx, expectedRequest)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal("user already exists", svcError.Message)

	suite.client.AssertExpectations(suite.T())
}
