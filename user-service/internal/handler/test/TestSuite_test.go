package test

import (
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/handler"
	mocks "github.com/COS301-SE-2025/Swift-Signals/user-service/internal/mocks/service"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	service *mocks.MockUserService
	handler *handler.Handler
}

func (suite *TestSuite) SetupTest() {
	suite.service = new(mocks.MockUserService)
	suite.handler = handler.NewUserHandler(suite.service)
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
