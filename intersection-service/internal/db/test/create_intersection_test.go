package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestCreateIntersection_Success() {
	// TODO: Add test case implementation
}

func TestDBCreateIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
