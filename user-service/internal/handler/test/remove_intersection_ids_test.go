package test

import (
	"context"
	"errors"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestRemoveIntersectionIDs_Success() {
	req := &userpb.RemoveIntersectionIDRequest{
		UserId:         "valid-user-id",
		IntersectionId: []string{"intersection-123"},
	}

	ctx := context.Background()

	suite.service.On("RemoveIntersectionIDs", ctx, req.UserId, req.IntersectionId).
		Return(nil)

	result, err := suite.handler.RemoveIntersectionIDs(ctx, req)

	suite.Require().NoError(err)
	suite.NotNil(result)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_Failure() {
	req := &userpb.RemoveIntersectionIDRequest{
		UserId:         "invalid-user-id",
		IntersectionId: []string{"intersection-123"},
	}

	suite.service.On("RemoveIntersectionIDs", mock.Anything, req.GetUserId(), req.GetIntersectionId()).
		Return(errors.New("user not found"))

	ctx := context.Background()

	result, err := suite.handler.RemoveIntersectionIDs(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
