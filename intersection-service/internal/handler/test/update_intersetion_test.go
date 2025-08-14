package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestUpdateIntersection_Success() {
	// TODO: Implement test case
	suite.True(true)
}

func TestHandlerUpdateIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
