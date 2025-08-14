package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestPutOptimisation_Success() {
	// TODO: add test case implementation
	suite.True(true)
}

func TestServicePutOptimisation(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
