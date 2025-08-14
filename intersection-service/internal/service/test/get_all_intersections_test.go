package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestGetAllIntersections_Success() {
	// TODO: add test case implementation
	suite.True(true)
}

func TestServiceGetAllIntersections(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
