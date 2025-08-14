package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestCreateIntersection_Success() {
	// TODO: Implement test case
	suite.True(true)
}

func TestHandlerCreateIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
