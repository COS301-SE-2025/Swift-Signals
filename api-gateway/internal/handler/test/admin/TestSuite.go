package admin

import (
	"context"
	"log/slog"
	"os"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/service"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	service *mocks.MockAdminServiceInterface
	handler *handler.AdminHandler
	ctx     context.Context
}

func (suite *TestSuite) SetupSuite() {
	slogger := slog.NewTextHandler(os.NewFile(0, os.DevNull), nil)
	slog.SetDefault(slog.New(slogger))
}

func (suite *TestSuite) SetupTest() {
	suite.service = new(mocks.MockAdminServiceInterface)
	suite.handler = handler.NewAdminHandler(suite.service)
	suite.ctx = middleware.SetRole(context.Background(), "admin")
}
