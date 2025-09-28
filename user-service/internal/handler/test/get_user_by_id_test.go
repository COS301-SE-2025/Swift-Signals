package test

import (
	"context"
	"errors"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestGetUserByID_Success() {
	req := &userpb.UserIDRequest{
		UserId: "valid-user-id",
	}

	expectedUser := &model.User{
		ID:              "valid-user-id",
		Name:            "John Doe",
		Email:           "john@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"intersection1", "intersection2"},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx := context.Background()

	suite.service.On("GetUserByID", ctx, req.UserId).
		Return(expectedUser, nil)

	result, err := suite.handler.GetUserByID(ctx, req)

	suite.Require().NoError(err)
	suite.Equal(expectedUser.ID, result.GetId())
	suite.Equal(expectedUser.Name, result.GetName())
	suite.Equal(expectedUser.Email, result.GetEmail())
	suite.Equal(expectedUser.IsAdmin, result.GetIsAdmin())
	suite.Equal(expectedUser.IntersectionIDs, result.GetIntersectionIds())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_Failure() {
	req := &userpb.UserIDRequest{
		UserId: "invalid-user-id",
	}

	suite.service.On("GetUserByID", mock.Anything, req.GetUserId()).
		Return(nil, errors.New("user not found"))

	ctx := context.Background()

	result, err := suite.handler.GetUserByID(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
