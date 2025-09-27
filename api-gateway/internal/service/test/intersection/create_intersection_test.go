package intersection

import (
	"context"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestCreateIntersection_Success() {
	userID := "valid-user-id"

	request := model.CreateIntersectionRequest{
		Name: "New Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "tjunction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	expectedCreateResponse := &intersectionpb.IntersectionResponse{
		Id: "new-intersection-id",
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	expectedIntersection := model.Intersection{
		Name: "New Intersection",
		Details: model.Details{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(expectedCreateResponse, nil)

	suite.userClient.On("AddIntersectionID", ctx, userID, "new-intersection-id").
		Return(nil, nil)

	result, err := suite.service.CreateIntersection(ctx, userID, request)

	suite.Require().NoError(err)
	suite.Equal("new-intersection-id", result.Id)

	suite.intrClient.AssertExpectations(suite.T())
	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_InvalidTrafficDensity() {
	userID := "valid-user-id"

	request := model.CreateIntersectionRequest{
		Name: "New Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "invalid-density",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "tjunction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	expectedIntersection := model.Intersection{
		Name: "New Intersection",
		Details: model.Details{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "invalid-density",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(nil, errs.NewValidationError("invalid traffic density", map[string]any{}))

	_, err := suite.service.CreateIntersection(ctx, userID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Contains(svcError.Message, "invalid traffic density")

	suite.intrClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_InvalidIntersectionType() {
	userID := "valid-user-id"

	request := model.CreateIntersectionRequest{
		Name: "New Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "invalid-type",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	expectedIntersection := model.Intersection{
		Name: "New Intersection",
		Details: model.Details{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "invalid-type",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(nil, errs.NewValidationError("invalid intersection type", map[string]any{}))

	_, err := suite.service.CreateIntersection(ctx, userID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Contains(svcError.Message, "invalid intersection type")

	suite.intrClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_IntersectionServiceError() {
	userID := "valid-user-id"

	request := model.CreateIntersectionRequest{
		Name: "New Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "tjunction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	expectedIntersection := model.Intersection{
		Name: "New Intersection",
		Details: model.Details{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	_, err := suite.service.CreateIntersection(ctx, userID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.intrClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_UserAssignmentError() {
	userID := "valid-user-id"

	request := model.CreateIntersectionRequest{
		Name: "New Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "tjunction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	expectedCreateResponse := &intersectionpb.IntersectionResponse{
		Id: "new-intersection-id",
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	expectedIntersection := model.Intersection{
		Name: "New Intersection",
		Details: model.Details{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(expectedCreateResponse, nil)

	suite.userClient.On("AddIntersectionID", ctx, userID, "new-intersection-id").
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{}))

	_, err := suite.service.CreateIntersection(ctx, userID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.intrClient.AssertExpectations(suite.T())
	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_EmptyName() {
	userID := "valid-user-id"

	request := model.CreateIntersectionRequest{
		Name: "",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "tjunction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	expectedIntersection := model.Intersection{
		Name: "",
		Details: model.Details{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(nil, errs.NewValidationError("intersection name is required", map[string]any{}))

	_, err := suite.service.CreateIntersection(ctx, userID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("intersection name is required", svcError.Message)

	suite.intrClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_DuplicateName() {
	userID := "valid-user-id"

	request := model.CreateIntersectionRequest{
		Name: "Existing Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "tjunction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	expectedIntersection := model.Intersection{
		Name: "Existing Intersection",
		Details: model.Details{
			Address:  "123 New Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(nil, errs.NewAlreadyExistsError("intersection with this name already exists", map[string]any{}))

	_, err := suite.service.CreateIntersection(ctx, userID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal("intersection with this name already exists", svcError.Message)

	suite.intrClient.AssertExpectations(suite.T())
}
