package test

import (
	"context"
	"errors"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TestSuite) TestUpdateUser_Success() {
	req := &userpb.UpdateUserRequest{
		UserId: "valid-user-id",
		Name:   "Updated Name",
		Email:  "updated@example.com",
	}

	expectedUser := &model.User{
		ID:              "valid-user-id",
		Name:            "Updated Name",
		Email:           "updated@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"intersection1"},
		CreatedAt:       time.Now().Add(-24 * time.Hour),
		UpdatedAt:       time.Now(),
	}

	ctx := context.Background()

	suite.service.On("UpdateUser", ctx, req.UserId, req.Name, req.Email).
		Return(expectedUser, nil)

	result, err := suite.handler.UpdateUser(ctx, req)

	suite.Require().NoError(err)
	suite.Equal(expectedUser.ID, result.GetId())
	suite.Equal(expectedUser.Name, result.GetName())
	suite.Equal(expectedUser.Email, result.GetEmail())
	suite.Equal(expectedUser.IsAdmin, result.GetIsAdmin())
	suite.Equal(expectedUser.IntersectionIDs, result.GetIntersectionIds())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_Failure() {
	req := &userpb.UpdateUserRequest{
		UserId: "invalid-user-id",
		Name:   "Updated Name",
		Email:  "updated@example.com",
	}

	suite.service.On("UpdateUser", mock.Anything, req.GetUserId(), req.GetName(), req.GetEmail()).
		Return(nil, errors.New("user not found"))

	ctx := context.Background()

	result, err := suite.handler.UpdateUser(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
