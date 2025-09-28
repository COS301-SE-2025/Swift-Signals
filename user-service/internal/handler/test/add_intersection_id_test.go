package test

import (
	"context"
	"errors"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestAddIntersectionID_Success() {
	req := &userpb.AddIntersectionIDRequest{
		UserId:         "valid-user-id",
		IntersectionId: "intersection-123",
	}

	ctx := context.Background()

	suite.service.On("AddIntersectionID", ctx, req.UserId, req.IntersectionId).
		Return(nil)

	result, err := suite.handler.AddIntersectionID(ctx, req)

	suite.Require().NoError(err)
	suite.NotNil(result)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_Failure() {
	req := &userpb.AddIntersectionIDRequest{
		UserId:         "invalid-user-id",
		IntersectionId: "intersection-123",
	}

	suite.service.On("AddIntersectionID", mock.Anything, req.GetUserId(), req.GetIntersectionId()).
		Return(errors.New("user not found"))

	ctx := context.Background()

	result, err := suite.handler.AddIntersectionID(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
