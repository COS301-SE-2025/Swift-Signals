package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestDeleteIntersection_Success() {
	// TODO: Add test case implementation
}

func TestDBDeleteIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
