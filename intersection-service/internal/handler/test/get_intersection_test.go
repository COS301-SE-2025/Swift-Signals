package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestGetIntersection_Success() {
	// TODO: Implement test case
	suite.True(true)
}

func TestHandlerGetIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
