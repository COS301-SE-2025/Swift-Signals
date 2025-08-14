package test

import (
	mocks "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/mocks/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	repo    *mocks.MockIntersectionRepository
	service service.IntersectionService
}

func (suite *TestSuite) SetupTest() {
	suite.repo = new(mocks.MockIntersectionRepository)
	suite.service = service.NewIntersectionService(suite.repo)
}
