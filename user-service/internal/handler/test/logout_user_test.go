package test

import (
	"context"
	"errors"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestLogoutUser_Success() {
	req := &userpb.UserIDRequest{
		UserId: "valid-user-id",
	}

	ctx := context.Background()

	suite.service.On("LogoutUser", ctx, req.UserId).
		Return(nil)

	result, err := suite.handler.LogoutUser(ctx, req)

	suite.Require().NoError(err)
	suite.NotNil(result)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogoutUser_Failure() {
	req := &userpb.UserIDRequest{
		UserId: "invalid-user-id",
	}

	suite.service.On("LogoutUser", mock.Anything, req.GetUserId()).
		Return(errors.New("user not found"))

	ctx := context.Background()

	result, err := suite.handler.LogoutUser(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
