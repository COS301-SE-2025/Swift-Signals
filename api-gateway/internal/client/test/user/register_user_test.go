package user

import (
	"context"
	"testing"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestRegisterUser_BuildsCorrectRequest() {
	suite.grpcClient.Mock.On("RegisterUser", mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *userpb.RegisterUserRequest) bool {
			return req.Name == "Valid Name" &&
				req.Email == "valid@gmail.com" &&
				req.Password == "8characters"
		})).Return(&userpb.UserResponse{Id: "any-id"}, nil)

	ctx := context.Background()
	_, err := suite.client.RegisterUser(ctx, "Valid Name", "valid@gmail.com", "8characters")

	suite.Require().NoError(err)

	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegisterUser_SetsTimeoutContext() {
	suite.grpcClient.On("RegisterUser",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.AnythingOfType("*user.RegisterUserRequest")).
		Return(&userpb.UserResponse{}, nil)

	ctx := context.Background()
	_, err := suite.client.RegisterUser(ctx, "Test Name", "test@gmail.com", "testpassword")

	suite.Require().NoError(err)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestClientRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
