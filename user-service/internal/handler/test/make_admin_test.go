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

func (suite *TestSuite) TestMakeAdmin_Success() {
	req := &userpb.AdminRequest{
		UserId:      "admin-user-id",
		AdminUserId: "target-user-id",
	}

	ctx := context.Background()

	suite.service.On("MakeAdmin", ctx, req.UserId, req.AdminUserId).
		Return(nil)

	result, err := suite.handler.MakeAdmin(ctx, req)

	suite.Require().NoError(err)
	suite.NotNil(result)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestMakeAdmin_Failure() {
	req := &userpb.AdminRequest{
		UserId:      "invalid-admin-id",
		AdminUserId: "target-user-id",
	}

	suite.service.On("MakeAdmin", mock.Anything, req.GetUserId(), req.GetAdminUserId()).
		Return(errors.New("insufficient permissions"))

	ctx := context.Background()

	result, err := suite.handler.MakeAdmin(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerMakeAdmin(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
