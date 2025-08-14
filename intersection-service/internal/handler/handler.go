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
	service service.IntersectionService
}

func NewIntersectionHandler(s service.IntersectionService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateIntersection(
	ctx context.Context,
	req *intersectionpb.CreateIntersectionRequest,
) (*intersectionpb.IntersectionResponse, error) {
	// TODO: Add logger

	intersectionDetails := h.mapIntersectionDetails(req.GetDetails())
	trafficDensity := model.TrafficDensity(req.GetTrafficDensity())
	optimisationParams := h.mapOptimisationParameters(req.GetDefaultParameters())

	intersection, err := h.service.CreateIntersection(
		ctx,
		req.GetName(),
		intersectionDetails,
		trafficDensity,
		optimisationParams,
	)
	// TODO: Add custom errors
	if err != nil {
		return nil, err
	}

	return h.mapToIntersection(intersection), nil
}

func (h *Handler) GetIntersection(
	ctx context.Context,
	req *intersectionpb.IntersectionIDRequest,
) (*intersectionpb.IntersectionResponse, error) {
	// TODO: Add logger
	intersection, err := h.service.GetIntersection(
		ctx,
		req.GetId(),
	)
	// TODO: Add custom errors
	if err != nil {
		return nil, err
	}

	return h.mapToIntersection(intersection), nil
}

func (h *Handler) GetAllIntersections(
	req *intersectionpb.GetAllIntersectionsRequest,
	stream intersectionpb.IntersectionService_GetAllIntersectionsServer,
) error {
	// TODO: Add logger
	ctx := stream.Context()

	intersections, err := h.service.GetAllIntersections(ctx)
	// TODO: Add custom errors
	if err != nil {
		return err
	}

	for _, intersection := range intersections {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		response := h.mapToIntersection(intersection)
		if response == nil {
			continue
		}

		// TODO: Add custom errors
		if err := stream.Send(response); err != nil {
			return status.Errorf(codes.Internal, "failed to send intersection: %v", err)
		}
	}

	return nil
}

func (h *Handler) UpdateIntersection(
	ctx context.Context,
	req *intersectionpb.UpdateIntersectionRequest,
) (*intersectionpb.IntersectionResponse, error) {
	// TODO: Add logger
	intersectionDetails := h.mapIntersectionDetails(req.GetDetails())

	intersection, err := h.service.UpdateIntersection(
		ctx,
		req.GetId(),
		req.GetName(),
		intersectionDetails,
	)
	// TODO: Add custom errors
	if err != nil {
		return nil, err
	}

	return h.mapToIntersection(intersection), nil
}

func (h *Handler) DeleteIntersection(
	ctx context.Context,
	req *intersectionpb.IntersectionIDRequest,
) (*emptypb.Empty, error) {
	// TODO: Add logger
	err := h.service.DeleteIntersection(ctx, req.GetId())
	// TODO: Add custom errors
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) PutOptimisation(
	ctx context.Context,
	req *intersectionpb.PutOptimisationRequest,
) (*intersectionpb.PutOptimisationResponse, error) {
	// TODO: Add logger
	optimisationParams := h.mapOptimisationParameters(req.GetParameters())

	optimisationResponse, err := h.service.PutOptimisation(
		ctx,
		req.GetId(),
		optimisationParams,
	)
	// TODO: Add custom errors
	if err != nil {
		return nil, err
	}

	return &intersectionpb.PutOptimisationResponse{Improved: optimisationResponse}, nil
}
