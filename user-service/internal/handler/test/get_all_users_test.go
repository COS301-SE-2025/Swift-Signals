package test

import (
	"context"
	"errors"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	grpcmocks "github.com/COS301-SE-2025/Swift-Signals/user-service/internal/mocks/grpc"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestGetAllUsers_Success() {
	req := &userpb.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter:   "",
	}

	expectedUsers := []*model.User{
		{
			ID:              "user1",
			Name:            "John Doe",
			Email:           "john@example.com",
			IsAdmin:         false,
			IntersectionIDs: []string{"intersection1"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              "user2",
			Name:            "Jane Smith",
			Email:           "jane@example.com",
			IsAdmin:         true,
			IntersectionIDs: []string{"intersection2", "intersection3"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	ctx := context.Background()

	mockStream := grpcmocks.NewMockUserService_GetAllUsersServer[userpb.UserResponse](suite.T())

	mockStream.On("Context").Return(ctx)

	suite.service.On("GetAllUsers", ctx, req.Page, req.PageSize, req.Filter).
		Return(expectedUsers, nil)

	for _, user := range expectedUsers {
		expectedResponse := &userpb.UserResponse{
			Id:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			IsAdmin:         user.IsAdmin,
			IntersectionIds: user.IntersectionIDs,
		}
		mockStream.On("Send", mock.MatchedBy(func(resp *userpb.UserResponse) bool {
			return resp.Id == expectedResponse.Id &&
				resp.Name == expectedResponse.Name &&
				resp.Email == expectedResponse.Email &&
				resp.IsAdmin == expectedResponse.IsAdmin &&
				len(resp.IntersectionIds) == len(expectedResponse.IntersectionIds)
		})).Return(nil)
	}

	err := suite.handler.GetAllUsers(req, mockStream)

	suite.Require().NoError(err)

	suite.service.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_ServiceFailure() {
	req := &userpb.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter:   "",
	}

	ctx := context.Background()

	mockStream := grpcmocks.NewMockUserService_GetAllUsersServer[userpb.UserResponse](suite.T())

	mockStream.On("Context").Return(ctx)

	suite.service.On("GetAllUsers", ctx, req.Page, req.PageSize, req.Filter).
		Return(nil, errors.New("database error"))

	err := suite.handler.GetAllUsers(req, mockStream)

	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_StreamSendFailure() {
	req := &userpb.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter:   "",
	}

	expectedUsers := []*model.User{
		{
			ID:              "user1",
			Name:            "John Doe",
			Email:           "john@example.com",
			IsAdmin:         false,
			IntersectionIDs: []string{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	ctx := context.Background()

	mockStream := grpcmocks.NewMockUserService_GetAllUsersServer[userpb.UserResponse](suite.T())

	mockStream.On("Context").Return(ctx)

	suite.service.On("GetAllUsers", ctx, req.Page, req.PageSize, req.Filter).
		Return(expectedUsers, nil)

	mockStream.On("Send", mock.Anything).Return(errors.New("stream send error"))

	err := suite.handler.GetAllUsers(req, mockStream)

	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}
