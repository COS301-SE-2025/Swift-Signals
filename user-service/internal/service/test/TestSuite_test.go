package test

import (
	"testing"

	mocks "github.com/COS301-SE-2025/Swift-Signals/user-service/internal/mocks/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	repo    *mocks.MockUserRepository
	service service.UserService
}

func (suite *TestSuite) SetupTest() {
	suite.repo = new(mocks.MockUserRepository)
	suite.repo.On("AdminExists", mock.Anything).Return(true, nil)
	suite.service = service.NewUserService(suite.repo)
}

func TestService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
