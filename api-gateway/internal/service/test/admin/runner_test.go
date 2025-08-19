package admin

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
