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

func (suite *TestSuite) TestGetUserByEmail_Success() {
	req := &userpb.GetUserByEmailRequest{
		Email: "john@example.com",
	}

	expectedUser := &model.User{
		ID:              "user-id-123",
		Name:            "John Doe",
		Email:           "john@example.com",
		IsAdmin:         true,
		IntersectionIDs: []string{"intersection1", "intersection2"},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx := context.Background()

	suite.service.On("GetUserByEmail", ctx, req.Email).
		Return(expectedUser, nil)

	result, err := suite.handler.GetUserByEmail(ctx, req)

	suite.Require().NoError(err)
	suite.Equal(expectedUser.ID, result.GetId())
	suite.Equal(expectedUser.Name, result.GetName())
	suite.Equal(expectedUser.Email, result.GetEmail())
	suite.Equal(expectedUser.IsAdmin, result.GetIsAdmin())
	suite.Equal(expectedUser.IntersectionIDs, result.GetIntersectionIds())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByEmail_Failure() {
	req := &userpb.GetUserByEmailRequest{
		Email: "nonexistent@example.com",
	}

	suite.service.On("GetUserByEmail", mock.Anything, req.GetEmail()).
		Return(nil, errors.New("user not found"))

	ctx := context.Background()

	result, err := suite.handler.GetUserByEmail(ctx, req)

	suite.Nil(result)
	suite.Require().Error(err)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, st.Code())
	suite.Equal("internal server error", st.Message())

	suite.service.AssertExpectations(suite.T())
}
