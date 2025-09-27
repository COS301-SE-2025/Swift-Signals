package user

import (
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	grpcClient *mocks.MockUserServiceClient
	client     *client.UserClient
}

func (suite *TestSuite) SetupTest() {
	suite.grpcClient = new(mocks.MockUserServiceClient)
	suite.client = client.NewUserClient(suite.grpcClient)
}
