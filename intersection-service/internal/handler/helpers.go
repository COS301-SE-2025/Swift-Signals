package handler

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
)

// =============================================================================
// VALIDATION HELPERS
// =============================================================================

func (h *Handler) validateCreateIntersectionRequest(req *intersectionpb.CreateIntersectionRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.GetName() == "" {
		return fmt.Errorf("intersection name is required")
	}

	if req.GetDetails() == nil {
		return fmt.Errorf("intersection details are required")
	}

	if err := h.validateIntersectionDetails(req.GetDetails()); err != nil {
		return fmt.Errorf("invalid intersection details: %w", err)
	}

	if req.GetDefaultParameters() == nil {
		return fmt.Errorf("default parameters are required")
	}

	return nil
}

func (h *Handler) validateIntersectionDetails(details *intersectionpb.IntersectionDetails) error {
	if details.GetAddress() == "" {
		return fmt.Errorf("address is required")
	}

	if details.GetCity() == "" {
		return fmt.Errorf("city is required")
	}

	return nil
}

// =============================================================================
// MAPPING HELPERS - PROTOBUF TO MODEL
// =============================================================================

func (h *Handler) mapIntersectionDetails(pbDetails *intersectionpb.IntersectionDetails) model.IntersectionDetails {
	if pbDetails == nil {
		return model.IntersectionDetails{}
	}

	return model.IntersectionDetails{
		Address:  pbDetails.GetAddress(),
		City:     pbDetails.GetCity(),
		Province: pbDetails.GetProvince(),
	}
}

func (h *Handler) mapSimulationParameters(pbParams *intersectionpb.SimulationParameters) model.SimulationParameters {
	if pbParams == nil {
		return model.SimulationParameters{}
	}

	return model.SimulationParameters{
		IntersectionType: model.IntersectionType(pbParams.GetIntersectionType()),
		Green:            pbParams.GetGreen(),
		Yellow:           pbParams.GetYellow(),
		Red:              pbParams.GetRed(),
		Speed:            pbParams.GetSpeed(),
		Seed:             pbParams.GetSeed(),
	}
}

func (h *Handler) mapOptimisationParameters(pbOptParams *intersectionpb.OptimisationParameters) model.OptimisationParameters {
	if pbOptParams == nil {
		return model.OptimisationParameters{}
	}

	return model.OptimisationParameters{
		OptimisationType: model.OptimisationType(pbOptParams.GetOptimisationType()),
		Parameters:       h.mapSimulationParameters(pbOptParams.GetParameters()),
	}
}

// =============================================================================
// MAPPING HELPERS - MODEL TO PROTOBUF
// =============================================================================

func (h *Handler) mapToProtoIntersectionDetails(details model.IntersectionDetails) *intersectionpb.IntersectionDetails {
	return &intersectionpb.IntersectionDetails{
		Address:  details.Address,
		City:     details.City,
		Province: details.Province,
	}
}

func (h *Handler) mapToProtoSimulationParameters(params model.SimulationParameters) *intersectionpb.SimulationParameters {
	return &intersectionpb.SimulationParameters{
		IntersectionType: intersectionpb.IntersectionType(params.IntersectionType),
		Green:            params.Green,
		Yellow:           params.Yellow,
		Red:              params.Red,
		Speed:            params.Speed,
		Seed:             params.Seed,
	}
}

func (h *Handler) mapToProtoOptimisationParameters(optParams model.OptimisationParameters) *intersectionpb.OptimisationParameters {
	return &intersectionpb.OptimisationParameters{
		OptimisationType: intersectionpb.OptimisationType(optParams.OptimisationType),
		Parameters:       h.mapToProtoSimulationParameters(optParams.Parameters),
	}
}

func (h *Handler) mapToIntersectionResponse(intersection *model.IntersectionResponse) *intersectionpb.IntersectionResponse {
	if intersection == nil {
		return nil
	}

	return &intersectionpb.IntersectionResponse{
		Id:                intersection.ID,
		Name:              intersection.Name,
		Details:           h.mapToProtoIntersectionDetails(intersection.Details),
		CreatedAt:         h.formatTime(intersection.CreatedAt),
		LastRunAt:         h.formatTime(intersection.LastRunAt),
		Status:            intersectionpb.IntersectionStatus(intersection.Status),
		RunCount:          intersection.RunCount,
		TrafficDensity:    intersectionpb.TrafficDensity(intersection.TrafficDensity),
		DefaultParameters: h.mapToProtoOptimisationParameters(intersection.DefaultParameters),
		BestParameters:    h.mapToProtoOptimisationParameters(intersection.BestParameters),
		CurrentParameters: h.mapToProtoOptimisationParameters(intersection.CurrentParameters),
	}
}

// =============================================================================
// UTILITY HELPERS
// =============================================================================

func (h *Handler) formatTime(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// =============================================================================
// ERROR HANDLING HELPERS
// =============================================================================

func (h *Handler) handleValidationError(err error) error {
	return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
}

func (h *Handler) handleServiceError(err error) error {
	switch {
	case errors.Is(err, model.ErrIntersectionExists):
		return status.Errorf(codes.AlreadyExists, "intersection already exists")
	case errors.Is(err, model.ErrInvalidParameters):
		return status.Errorf(codes.InvalidArgument, "invalid parameters: %v", err)
	case errors.Is(err, model.ErrIntersectionNotFound):
		return status.Errorf(codes.NotFound, "intersection not found")
	default:
		return status.Errorf(codes.Internal, "internal server error")
	}
}
