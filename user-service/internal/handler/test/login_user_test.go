package test

import (
	"context"
	"errors"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestLoginUser_Success() {
	req := &userpb.LoginUserRequest{
		Email:    "valid@gmail.com",
		Password: "validpassword",
	}

	expectedToken := "jwt-token-12345"
	expectedExpiryTime := time.Now().Add(24 * time.Hour)

	ctx := context.Background()

	suite.service.On("LoginUser", ctx, req.Email, req.Password).
		Return(expectedToken, expectedExpiryTime, nil)

	result, err := suite.handler.LoginUser(ctx, req)

	suite.Require().NoError(err)
	suite.Equal(expectedToken, result.GetToken())
	suite.Equal(expectedExpiryTime.Unix(), result.GetExpiresAt().GetSeconds())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLoginUser_Failure() {
	req := &userpb.LoginUserRequest{
		Email:    "invalid@example.com",
		Password: "wrongpassword",
	}

	suite.service.On("LoginUser", mock.Anything, req.GetEmail(), req.GetPassword()).
		Return("", time.Time{}, errors.New("invalid credentials"))

	ctx := context.Background()

	result, err := suite.handler.LoginUser(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
