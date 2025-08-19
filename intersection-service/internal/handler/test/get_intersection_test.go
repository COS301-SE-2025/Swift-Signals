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
	"github.com/stretchr/testify/require"
)

func TestGetIntersection_Success(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handler.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.IntersectionIDRequest{Id: "int-001"}

	expectedModel := &model.Intersection{
		ID:   "int-001",
		Name: "Main & First",
		Details: model.IntersectionDetails{
			Address:  "123 Main St",
			City:     "Pretoria",
			Province: "Gauteng",
		},
		TrafficDensity: model.TrafficMedium,
	}

	mockService.On("GetIntersection",
		mock.Anything,
		req.Id,
	).Return(expectedModel, nil)

	resp, err := h.GetIntersection(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedModel.ID, resp.GetId())
	mockService.AssertExpectations(t)
}

func TestGetIntersection_NotFound(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handler.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.IntersectionIDRequest{Id: "nonexistent"}

	mockService.On("GetIntersection",
		mock.Anything,
		req.Id,
	).Return(nil, errors.New("intersection not found"))

	resp, err := h.GetIntersection(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestGetIntersection_DatabaseError(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handler.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.IntersectionIDRequest{Id: "int-002"}

	mockService.On("GetIntersection",
		mock.Anything,
		req.Id,
	).Return(nil, errors.New("database failure"))

	resp, err := h.GetIntersection(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}
