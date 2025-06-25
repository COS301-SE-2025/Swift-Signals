package handler

import (
	"context"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (h *Handler) GetIntersection(ctx context.Context, req *intersectionpb.IntersectionIDRequest) (*intersectionpb.IntersectionResponse, error) {
	// Validate request
	if err := h.validateIntersectionIDRequest(req); err != nil {
		return nil, h.handleValidationError(err)
	}

	// Call service layer
	intersection, err := h.service.GetIntersection(
		ctx,
		req.GetId(),
	)
	if err != nil {
		return nil, h.handleServiceError(err)
	}

	// Convert model response to protobuf response
	return h.mapToIntersectionResponse(intersection), nil
}

func (h *Handler) GetAllIntersections(req *intersectionpb.GetAllIntersectionsRequest, stream intersectionpb.IntersectionService_GetAllIntersectionsServer) error {
	ctx := stream.Context()

	// Validate request
	if err := h.validateGetAllIntersectionsRequest(req); err != nil {
		return h.handleValidationError(err)
	}

	// Get all intersections from service
	intersections, err := h.service.GetAllIntersections(ctx)
	if err != nil {
		return h.handleServiceError(err)
	}

	// Stream each intersection back to client
	for _, intersection := range intersections {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Convert model to protobuf response
		response := h.mapToIntersectionResponse(intersection)
		if response == nil {
			continue // Skip nil intersections
		}

		// Send the intersection to the client
		if err := stream.Send(response); err != nil {
			return status.Errorf(codes.Internal, "failed to send intersection: %v", err)
		}
	}

	return nil
}

func (h *Handler) UpdateIntersection(ctx context.Context, req *intersectionpb.UpdateIntersectionRequest) (*intersectionpb.IntersectionResponse, error) {
	// Validate request
	if err := h.validateUpdateIntersectionRequest(req); err != nil {
		return nil, h.handleValidationError(err)
	}

	// Convert protobuf request to model objects
	intersectionDetails := h.mapIntersectionDetails(req.GetDetails())

	// Call service layer
	intersection, err := h.service.UpdateIntersection(
		ctx,
		req.GetId(),
		req.GetName(),
		intersectionDetails,
	)
	if err != nil {
		return nil, h.handleServiceError(err)
	}

	// Convert model response to protobuf response
	return h.mapToIntersectionResponse(intersection), nil
}

func (h *Handler) DeleteIntersection(ctx context.Context, req *intersectionpb.IntersectionIDRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := h.validateIntersectionIDRequest(req); err != nil {
		return nil, h.handleValidationError(err)
	}

	// Call service layer to delete intersection
	err := h.service.DeleteIntersection(ctx, req.GetId())
	if err != nil {
		return nil, h.handleServiceError(err)
	}

	// Return empty response on success
	return &emptypb.Empty{}, nil
}

func (h *Handler) PutOptimisation(ctx context.Context, req *intersectionpb.PutOptimisationRequest) (*intersectionpb.PutOptimisationResponse, error) {
	// Validate request
	if err := h.validatePutOptimisationRequest(req); err != nil {
		return nil, h.handleValidationError(err)
	}

	// Convert protobuf request to model objects
	optimisationParams := h.mapOptimisationParameters(req.GetParameters())

	// Call service layer
	optimisationResponse, err := h.service.PutOptimisation(
		ctx,
		req.GetId(),
		optimisationParams,
	)
	if err != nil {
		return nil, h.handleServiceError(err)
	}

	// Convert model response to protobuf response
	return h.mapToPutOptimisationResponse(optimisationResponse), nil
}
