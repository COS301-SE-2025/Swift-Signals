package test

import (
	"context"
	"errors"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	grpcmocks "github.com/COS301-SE-2025/Swift-Signals/user-service/internal/mocks/grpc"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestGetUserIntersectionIDs_Success() {
	req := &userpb.UserIDRequest{
		UserId: "valid-user-id",
	}

	expectedIntersectionIDs := []string{
		"intersection1",
		"intersection2",
		"intersection3",
	}

	ctx := context.Background()

	mockStream := grpcmocks.NewMockUserService_GetUserIntersectionIDsServer[userpb.IntersectionIDResponse](
		suite.T(),
	)

	mockStream.On("Context").Return(ctx)

	suite.service.On("GetUserIntersectionIDs", ctx, req.UserId).
		Return(expectedIntersectionIDs, nil)

	for _, intersectionID := range expectedIntersectionIDs {
		expectedResponse := &userpb.IntersectionIDResponse{
			IntersectionId: intersectionID,
		}
		mockStream.On("Send", mock.MatchedBy(func(resp *userpb.IntersectionIDResponse) bool {
			return resp.IntersectionId == expectedResponse.IntersectionId
		})).Return(nil)
	}

	err := suite.handler.GetUserIntersectionIDs(req, mockStream)

	suite.Require().NoError(err)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_ServiceFailure() {
	req := &userpb.UserIDRequest{
		UserId: "invalid-user-id",
	}

	ctx := context.Background()

	mockStream := grpcmocks.NewMockUserService_GetUserIntersectionIDsServer[userpb.IntersectionIDResponse](
		suite.T(),
	)

	mockStream.On("Context").Return(ctx)

	suite.service.On("GetUserIntersectionIDs", ctx, req.UserId).
		Return(nil, errors.New("user not found"))

	err := suite.handler.GetUserIntersectionIDs(req, mockStream)

	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_StreamSendFailure() {
	req := &userpb.UserIDRequest{
		UserId: "valid-user-id",
	}

	expectedIntersectionIDs := []string{"intersection1"}

	ctx := context.Background()

	mockStream := grpcmocks.NewMockUserService_GetUserIntersectionIDsServer[userpb.IntersectionIDResponse](
		suite.T(),
	)

	mockStream.On("Context").Return(ctx)

	suite.service.On("GetUserIntersectionIDs", ctx, req.UserId).
		Return(expectedIntersectionIDs, nil)

	mockStream.On("Send", mock.Anything).Return(errors.New("stream send error"))

	err := suite.handler.GetUserIntersectionIDs(req, mockStream)

	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserIntersectionIDs_EmptyList() {
	req := &userpb.UserIDRequest{
		UserId: "user-with-no-intersections",
	}

	expectedIntersectionIDs := []string{}

	ctx := context.Background()

	mockStream := grpcmocks.NewMockUserService_GetUserIntersectionIDsServer[userpb.IntersectionIDResponse](
		suite.T(),
	)

	mockStream.On("Context").Return(ctx)

	suite.service.On("GetUserIntersectionIDs", ctx, req.UserId).
		Return(expectedIntersectionIDs, nil)

	err := suite.handler.GetUserIntersectionIDs(req, mockStream)

	suite.Require().NoError(err)

	suite.service.AssertExpectations(suite.T())
}
