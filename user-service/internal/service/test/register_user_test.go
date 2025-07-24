package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestRegisterUser() {
	// TODO: Implement the unit test for TestRegisterUser
	suite.True(true)
}

func TestServiceRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
