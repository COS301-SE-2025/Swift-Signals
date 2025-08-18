package test

import (
	"context"
	"errors"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/handler"
	serviceMocks "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/mocks/service"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateIntersection_Success(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handler.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.CreateIntersectionRequest{
		Name: "Main & First",
		Details: &intersectionpb.IntersectionDetails{
			Address:  "123 Main St",
			City:     "Pretoria",
			Province: "Gauteng",
		},
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM,
		DefaultParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
		},
	}

	expectedModel := &model.Intersection{
		ID:   "int-001",
		Name: req.Name,
		Details: model.IntersectionDetails{
			Address:  req.Details.Address,
			City:     req.Details.City,
			Province: req.Details.Province,
		},
		TrafficDensity: model.TrafficMedium,
	}

	// Mock the service call
	mockService.On("CreateIntersection",
		mock.Anything,
		req.Name,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(expectedModel, nil)

	resp, err := h.CreateIntersection(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedModel.ID, resp.GetId())
	mockService.AssertExpectations(t)
}

func TestCreateIntersection_ServiceError(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handler.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.CreateIntersectionRequest{
		Name: "Bad Intersection",
		Details: &intersectionpb.IntersectionDetails{
			Address: "Unknown",
			City:    "Nowhere",
		},
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_LOW,
	}

	mockService.On("CreateIntersection",
		mock.Anything,
		req.Name,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, errors.New("database failure"))

	resp, err := h.CreateIntersection(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestCreateIntersection_EmptyName(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handler.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.CreateIntersectionRequest{
		Name: "",
		Details: &intersectionpb.IntersectionDetails{
			Address: "456 Side St",
			City:    "Cape Town",
		},
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	mockService.On("CreateIntersection",
		mock.Anything,
		req.Name,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, errors.New("invalid intersection name"))

	resp, err := h.CreateIntersection(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}
