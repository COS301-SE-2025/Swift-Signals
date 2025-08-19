package test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handlerPkg "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/handler"
	handlerMocks "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/mocks/service"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateIntersection_Success(t *testing.T) {
	mockService := new(handlerMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	req := &intersectionpb.UpdateIntersectionRequest{
		Id:   "int-001",
		Name: "Updated Intersection",
	}

	mockService.On("UpdateIntersection", mock.Anything, req.GetId(), req.GetName(), mock.Anything).
		Return(&model.Intersection{
			ID:   req.GetId(),
			Name: req.GetName(),
		}, nil)

	resp, err := h.UpdateIntersection(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, req.GetId(), resp.GetId())
	assert.Equal(t, req.GetName(), resp.GetName())

	mockService.AssertExpectations(t)
}

func TestUpdateIntersection_ServiceError(t *testing.T) {
	mockService := new(handlerMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	req := &intersectionpb.UpdateIntersectionRequest{
		Id:   "int-001",
		Name: "Updated Intersection",
	}

	mockService.On("UpdateIntersection", mock.Anything, req.GetId(), req.GetName(), mock.Anything).
		Return(nil, errors.New("database failure"))

	resp, err := h.UpdateIntersection(context.Background(), req)

	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "internal server error")

	mockService.AssertExpectations(t)
}

func TestUpdateIntersection_ContextCancelled(t *testing.T) {
	mockService := new(handlerMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel context

	req := &intersectionpb.UpdateIntersectionRequest{
		Id:   "int-001",
		Name: "Updated Intersection",
	}

	mockService.On("UpdateIntersection", mock.Anything, req.GetId(), req.GetName(), mock.Anything).
		Return(nil, context.Canceled)

	resp, err := h.UpdateIntersection(ctx, req)

	assert.Nil(t, resp)

	_, ok := status.FromError(err)
	assert.True(t, ok)

	mockService.AssertExpectations(t)
}
