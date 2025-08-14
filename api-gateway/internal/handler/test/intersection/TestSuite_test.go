package intersection

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
	service *mocks.MockIntersectionServiceInterface
	handler *handler.IntersectionHandler
	ctx     context.Context
}

func (suite *TestSuite) SetupSuite() {
	slogger := slog.NewTextHandler(os.NewFile(0, os.DevNull), nil)
	slog.SetDefault(slog.New(slogger))
}

func (suite *TestSuite) SetupTest() {
	suite.service = new(mocks.MockIntersectionServiceInterface)
	suite.handler = handler.NewIntersectionHandler(suite.service)
	suite.ctx = middleware.SetUserID(context.Background(), "test-user-id")
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
