package handler

import (
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
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
	pbParams *commonpb.SimulationParameters,
) model.SimulationParameters {
	if pbParams == nil {
		return model.SimulationParameters{}
	}

	return model.SimulationParameters{
		IntersectionType: model.IntersectionType(
			commonpb.IntersectionType_name[int32(pbParams.GetIntersectionType())],
		),
		Green:  int(pbParams.GetGreen()),
		Yellow: int(pbParams.GetYellow()),
		Red:    int(pbParams.GetRed()),
		Speed:  int(pbParams.GetSpeed()),
		Seed:   int(pbParams.GetSeed()),
	}
}

func (h *Handler) mapOptimisationParameters(
	pbOptParams *commonpb.OptimisationParameters,
) model.OptimisationParameters {
	if pbOptParams == nil {
		return model.OptimisationParameters{}
	}

	return model.OptimisationParameters{
		OptimisationType: model.OptimisationType(
			commonpb.OptimisationType_name[int32(pbOptParams.GetOptimisationType())],
		),
		Parameters: h.mapSimulationParameters(pbOptParams.GetParameters()),
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
) *commonpb.SimulationParameters {
	intersectionType := commonpb.IntersectionType_INTERSECTION_TYPE_TRAFFICLIGHT
	if enumValue, ok := commonpb.IntersectionType_value[string(params.IntersectionType)]; ok {
		intersectionType = commonpb.IntersectionType(enumValue)
	}
	return &commonpb.SimulationParameters{
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
) *commonpb.OptimisationParameters {
	return &commonpb.OptimisationParameters{
		OptimisationType: commonpb.OptimisationType(
			commonpb.OptimisationType_value[string(optParams.OptimisationType)],
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
		Status: commonpb.IntersectionStatus(
			commonpb.IntersectionStatus_value[string(intersection.Status)]),
		RunCount: int32(intersection.RunCount),
		TrafficDensity: commonpb.TrafficDensity(
			commonpb.TrafficDensity_value[string(intersection.TrafficDensity)]),
		DefaultParameters: h.mapToProtoOptimisationParameters(intersection.DefaultParameters),
		BestParameters:    h.mapToProtoOptimisationParameters(intersection.BestParameters),
		CurrentParameters: h.mapToProtoOptimisationParameters(intersection.CurrentParameters),
	}
}
