package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestGetAllIntersections_Success() {
	// TODO: Add test case implementation
}

func TestDBGetAllIntersections(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
