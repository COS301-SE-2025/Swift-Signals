package test

import (
	"context"
	"errors"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestRegisterUser_Success() {
	req := &userpb.RegisterUserRequest{
		Name:     "Valid User",
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	expectedUser := &model.User{
		ID:              "generated id",
		Name:            "Valid User",
		Email:           "valid@gmail.com",
		IsAdmin:         false,
		IntersectionIDs: nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx := context.Background()

	suite.service.On("RegisterUser", ctx, req.Name, req.Email, req.Password).
		Return(expectedUser, nil)

	result, err := suite.handler.RegisterUser(ctx, req)

	suite.Require().NoError(err)
	suite.Equal(expectedUser.ID, result.GetId())
	suite.Equal(expectedUser.Name, result.GetName())
	suite.Equal(expectedUser.Email, result.GetEmail())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_Failure() {
	req := &userpb.RegisterUserRequest{
		Name:     "Invalid",
		Email:    "fail@example.com",
		Password: "weak",
	}

	suite.service.On("RegisterUser", mock.Anything, req.GetName(), req.GetEmail(), req.GetPassword()).
		Return(nil, errors.New("any error"))

	ctx := context.Background()

	result, err := suite.handler.RegisterUser(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
