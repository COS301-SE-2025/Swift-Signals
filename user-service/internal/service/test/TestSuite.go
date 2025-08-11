package test

import (
	mocks "github.com/COS301-SE-2025/Swift-Signals/user-service/internal/mocks/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/service"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	repo    *mocks.MockUserRepository
	service service.UserService
}

func (suite *TestSuite) SetupTest() {
	suite.repo = new(mocks.MockUserRepository)
	suite.service = service.NewUserService(suite.repo)
}
