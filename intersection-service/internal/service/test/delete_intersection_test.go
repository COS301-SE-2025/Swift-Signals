package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestDeleteIntersection_Success() {
	// TODO: add test case implementation
	suite.True(true)
}

func TestServiceDeleteIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
