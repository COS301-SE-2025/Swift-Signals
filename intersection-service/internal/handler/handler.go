package handler

import (
	"context"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
)

type Handler struct {
	intersectionpb.UnimplementedIntersectionServiceServer
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateIntersection(ctx context.Context, req *intersectionpb.CreateIntersectionRequest) (*intersectionpb.IntersectionResponse, error) {
	// Validate request
	if err := h.validateCreateIntersectionRequest(req); err != nil {
		return nil, h.handleValidationError(err)
	}

	// Convert protobuf request to model objects
	intersectionDetails := h.mapIntersectionDetails(req.GetDetails())
	trafficDensity := model.TrafficDensity(req.GetTrafficDensity())
	optimisationParams := h.mapOptimisationParameters(req.GetDefaultParameters())

	// Call service layer
	intersection, err := h.service.CreateIntersection(
		ctx,
		req.GetName(),
		intersectionDetails,
		trafficDensity,
		optimisationParams,
	)
	if err != nil {
		return nil, h.handleServiceError(err)
	}

	// Convert model response to protobuf response
	return h.mapToIntersectionResponse(intersection), nil
}
