package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestUpdateIntersection_Success() {
	// TODO: add test case implementation
	suite.True(true)
}

func TestServiceUpdateIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
