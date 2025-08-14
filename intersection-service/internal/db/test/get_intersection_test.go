package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestGetIntersection_Success() {
	// TODO: Add test case implementation
}

func TestDBGetIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
