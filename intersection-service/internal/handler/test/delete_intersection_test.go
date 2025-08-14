package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestDeleteIntersection_Success() {
	// TODO: Implement test case
	suite.True(true)
}

func TestHandlerDeleteIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
