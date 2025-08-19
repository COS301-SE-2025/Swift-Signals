package test

import (
	"context"
	"errors"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestRemoveAdmin_Success() {
	req := &userpb.AdminRequest{
		UserId:      "admin-user-id",
		AdminUserId: "target-admin-id",
	}

	ctx := context.Background()

	suite.service.On("RemoveAdmin", ctx, req.UserId, req.AdminUserId).
		Return(nil)

	result, err := suite.handler.RemoveAdmin(ctx, req)

	suite.Require().NoError(err)
	suite.NotNil(result)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveAdmin_Failure() {
	req := &userpb.AdminRequest{
		UserId:      "invalid-admin-id",
		AdminUserId: "target-admin-id",
	}

	suite.service.On("RemoveAdmin", mock.Anything, req.GetUserId(), req.GetAdminUserId()).
		Return(errors.New("insufficient permissions"))

	ctx := context.Background()

	result, err := suite.handler.RemoveAdmin(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
