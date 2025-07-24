package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestCreateUser() {
	// TODO: Implement the unit test for TestCreateUser
	suite.True(true)
}

func TestDBRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
