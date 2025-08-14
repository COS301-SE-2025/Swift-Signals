package auth

import (
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/handler"
	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/service"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	service *mocks.MockAuthServiceInterface
	handler *handler.AuthHandler
}

func (suite *TestSuite) SetupTest() {
	suite.service = new(mocks.MockAuthServiceInterface)
	suite.handler = handler.NewAuthHandler(suite.service)
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
