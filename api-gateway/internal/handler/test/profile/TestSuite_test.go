package profile

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/service"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	service *mocks.MockProfileServiceInterface
	handler *handler.ProfileHandler
	ctx     context.Context
}

func (suite *TestSuite) SetupSuite() {
	slogger := slog.NewTextHandler(os.NewFile(0, os.DevNull), nil)
	slog.SetDefault(slog.New(slogger))
}

func (suite *TestSuite) SetupTest() {
	suite.service = new(mocks.MockProfileServiceInterface)
	suite.handler = handler.NewProfileHandler(suite.service)
	suite.ctx = middleware.SetUserID(context.Background(), "test-user-id")
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
