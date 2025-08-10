package test

import (
	"context"
	"errors"
	"testing"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestResetPassword_Success() {
	req := &userpb.ResetPasswordRequest{
		Email: "valid@example.com",
	}

	ctx := context.Background()

	suite.service.On("ResetPassword", ctx, req.Email).
		Return(nil)

	result, err := suite.handler.ResetPassword(ctx, req)

	suite.Require().NoError(err)
	suite.NotNil(result)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestResetPassword_Failure() {
	req := &userpb.ResetPasswordRequest{
		Email: "nonexistent@example.com",
	}

	suite.service.On("ResetPassword", mock.Anything, req.GetEmail()).
		Return(errors.New("user not found"))

	ctx := context.Background()

	result, err := suite.handler.ResetPassword(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerResetPassword(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
