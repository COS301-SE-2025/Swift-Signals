package auth

import (
	"testing"

	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	client  *mocks.MockUserClientInterface
	service service.AuthServiceInterface
}

func (suite *TestSuite) SetupTest() {
	suite.client = new(mocks.MockUserClientInterface)
	suite.service = service.NewAuthService(suite.client)
}

func TestService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
