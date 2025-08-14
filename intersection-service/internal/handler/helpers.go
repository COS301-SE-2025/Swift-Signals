package handler

import (
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// =============================================================================
// MAPPING HELPERS - PROTOBUF TO MODEL
// =============================================================================

func (h *Handler) mapIntersectionDetails(
	pbDetails *intersectionpb.IntersectionDetails,
) model.IntersectionDetails {
	if pbDetails == nil {
		return model.IntersectionDetails{}
	}

	return model.IntersectionDetails{
		Address:  pbDetails.GetAddress(),
		City:     pbDetails.GetCity(),
		Province: pbDetails.GetProvince(),
	}
}

func (h *Handler) mapSimulationParameters(
	pbParams *intersectionpb.SimulationParameters,
) model.SimulationParameters {
	if pbParams == nil {
		return model.SimulationParameters{}
	}

	return model.SimulationParameters{
		IntersectionType: model.IntersectionType(pbParams.GetIntersectionType()),
		Green:            int(pbParams.GetGreen()),
		Yellow:           int(pbParams.GetYellow()),
		Red:              int(pbParams.GetRed()),
		Speed:            int(pbParams.GetSpeed()),
		Seed:             int(pbParams.GetSeed()),
	}
}

func (h *Handler) mapOptimisationParameters(
	pbOptParams *intersectionpb.OptimisationParameters,
) model.OptimisationParameters {
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

func (h *Handler) mapToProtoIntersectionDetails(
	details model.IntersectionDetails,
) *intersectionpb.IntersectionDetails {
	return &intersectionpb.IntersectionDetails{
		Address:  details.Address,
		City:     details.City,
		Province: details.Province,
	}
}

func (h *Handler) mapToProtoSimulationParameters(
	params model.SimulationParameters,
) *intersectionpb.SimulationParameters {
	intersectionType := intersectionpb.IntersectionType_INTERSECTION_TYPE_TRAFFICLIGHT
	if enumValue, ok := intersectionpb.IntersectionType_value[string(params.IntersectionType)]; ok {
		intersectionType = intersectionpb.IntersectionType(enumValue)
	}
	return &intersectionpb.SimulationParameters{
		IntersectionType: intersectionType,
		Green:            int32(params.Green),
		Yellow:           int32(params.Yellow),
		Red:              int32(params.Red),
		Speed:            int32(params.Speed),
		Seed:             int32(params.Seed),
	}
}

func (h *Handler) mapToProtoOptimisationParameters(
	optParams model.OptimisationParameters,
) *intersectionpb.OptimisationParameters {
	return &intersectionpb.OptimisationParameters{
		OptimisationType: intersectionpb.OptimisationType(
			intersectionpb.OptimisationType_value[string(optParams.OptimisationType)],
		),
		Parameters: h.mapToProtoSimulationParameters(optParams.Parameters),
	}
}

func (h *Handler) mapToIntersection(
	intersection *model.Intersection,
) *intersectionpb.IntersectionResponse {
	if intersection == nil {
		return nil
	}

	return &intersectionpb.IntersectionResponse{
		Id:        intersection.ID,
		Name:      intersection.Name,
		Details:   h.mapToProtoIntersectionDetails(intersection.Details),
		CreatedAt: timestamppb.New(intersection.CreatedAt),
		LastRunAt: timestamppb.New(intersection.LastRunAt),
		Status: intersectionpb.IntersectionStatus(
			intersectionpb.IntersectionStatus_value[string(intersection.Status)]),
		RunCount: int32(intersection.RunCount),
		TrafficDensity: intersectionpb.TrafficDensity(
			intersectionpb.TrafficDensity_value[string(intersection.TrafficDensity)]),
		DefaultParameters: h.mapToProtoOptimisationParameters(intersection.DefaultParameters),
		BestParameters:    h.mapToProtoOptimisationParameters(intersection.BestParameters),
		CurrentParameters: h.mapToProtoOptimisationParameters(intersection.CurrentParameters),
	}
}
