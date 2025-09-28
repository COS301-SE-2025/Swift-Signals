package test

import (
	"context"
	"errors"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestChangePassword_Success() {
	req := &userpb.ChangePasswordRequest{
		UserId:          "valid-user-id",
		CurrentPassword: "oldpassword123",
		NewPassword:     "newpassword456",
	}

	ctx := context.Background()

	suite.service.On("ChangePassword", ctx, req.UserId, req.CurrentPassword, req.NewPassword).
		Return(nil)

	result, err := suite.handler.ChangePassword(ctx, req)

	suite.Require().NoError(err)
	suite.NotNil(result)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestChangePassword_Failure() {
	req := &userpb.ChangePasswordRequest{
		UserId:          "invalid-user-id",
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword456",
	}

	suite.service.On("ChangePassword", mock.Anything, req.GetUserId(), req.GetCurrentPassword(), req.GetNewPassword()).
		Return(errors.New("invalid current password"))

	ctx := context.Background()

	result, err := suite.handler.ChangePassword(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
