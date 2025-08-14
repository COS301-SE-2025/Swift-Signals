package handler

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/util"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
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
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing CreateIntersection request")

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
	if err != nil {
		logger.Error("failed to create intersection",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("CreateIntersection successful")
	return h.mapToIntersection(intersection), nil
}

func (h *Handler) GetIntersection(
	ctx context.Context,
	req *intersectionpb.IntersectionIDRequest,
) (*intersectionpb.IntersectionResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing GetIntersection request")

	intersection, err := h.service.GetIntersection(
		ctx,
		req.GetId(),
	)
	if err != nil {
		logger.Error("failed to find intersection",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("GetIntersection successful")
	return h.mapToIntersection(intersection), nil
}

func (h *Handler) GetAllIntersections(
	req *intersectionpb.GetAllIntersectionsRequest,
	stream intersectionpb.IntersectionService_GetAllIntersectionsServer,
) error {
	ctx := stream.Context()
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing GetAllIntersections request")

	intersections, err := h.service.GetAllIntersections(
		ctx,
		int(req.GetPage()),
		int(req.GetPageSize()),
		req.GetFilter(),
	)
	if err != nil {
		logger.Error("failed to find all intersections",
			"error", err.Error(),
		)
		return errs.HandleServiceError(err)
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

		if err := stream.Send(response); err != nil {
			logger.Error("failed to send intersection",
				"error", err.Error(),
			)
			return errs.HandleServiceError(err)
		}
	}

	logger.Info("GetAllIntersections successful")
	return nil
}

func (h *Handler) UpdateIntersection(
	ctx context.Context,
	req *intersectionpb.UpdateIntersectionRequest,
) (*intersectionpb.IntersectionResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing UpdateIntersection request")

	intersectionDetails := h.mapIntersectionDetails(req.GetDetails())

	intersection, err := h.service.UpdateIntersection(
		ctx,
		req.GetId(),
		req.GetName(),
		intersectionDetails,
	)
	if err != nil {
		logger.Error("failed to update intersection",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("UpdateIntersection successful")
	return h.mapToIntersection(intersection), nil
}

func (h *Handler) DeleteIntersection(
	ctx context.Context,
	req *intersectionpb.IntersectionIDRequest,
) (*emptypb.Empty, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing DeleteIntersection request")

	err := h.service.DeleteIntersection(ctx, req.GetId())
	if err != nil {
		logger.Error("failed to delete intersection",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("DeleteIntersection successful")
	return &emptypb.Empty{}, nil
}

func (h *Handler) PutOptimisation(
	ctx context.Context,
	req *intersectionpb.PutOptimisationRequest,
) (*intersectionpb.PutOptimisationResponse, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("processing PutOptimisation request")

	optimisationParams := h.mapOptimisationParameters(req.GetParameters())

	optimisationResponse, err := h.service.PutOptimisation(
		ctx,
		req.GetId(),
		optimisationParams,
	)
	if err != nil {
		logger.Error("failed to update intersection optimisation params",
			"error", err.Error(),
		)
		return nil, errs.HandleServiceError(err)
	}

	logger.Info("PutOptimisation successful")
	return &intersectionpb.PutOptimisationResponse{Improved: optimisationResponse}, nil
}
