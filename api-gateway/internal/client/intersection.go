package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IntersectionClient struct {
	client intersectionpb.IntersectionServiceClient
}

func NewIntersectionClient(conn *grpc.ClientConn) *IntersectionClient {
	return &IntersectionClient{
		client: intersectionpb.NewIntersectionServiceClient(conn),
	}
}

func NewIntersectionClientWithGrpcClient(
	client intersectionpb.IntersectionServiceClient,
) *IntersectionClient {
	return &IntersectionClient{
		client: client,
	}
}

func (ic *IntersectionClient) CreateIntersection(
	ctx context.Context,
	intersection model.Intersection,
) (*intersectionpb.IntersectionResponse, error) {
	req := &intersectionpb.CreateIntersectionRequest{
		Name:              intersection.Name,
		Details:           convertDetailsToProto(intersection.Details),
		TrafficDensity:    StringToTrafficDensity(intersection.TrafficDensity),
		DefaultParameters: convertParametersToProto(intersection.DefaultParameters),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := ic.client.CreateIntersection(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (ic *IntersectionClient) GetIntersection(
	ctx context.Context,
	id string,
) (*intersectionpb.IntersectionResponse, error) {
	req := &intersectionpb.IntersectionIDRequest{
		Id: id,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := ic.client.GetIntersection(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (ic *IntersectionClient) GetAllIntersections(
	ctx context.Context,
) (intersectionpb.IntersectionService_GetAllIntersectionsClient, error) {
	req := &intersectionpb.GetAllIntersectionsRequest{}

	return ic.client.GetAllIntersections(ctx, req)
}

func (ic *IntersectionClient) UpdateIntersection(
	ctx context.Context,
	id, name string,
	details model.Details,
) (*intersectionpb.IntersectionResponse, error) {
	req := &intersectionpb.UpdateIntersectionRequest{
		Id:      id,
		Name:    name,
		Details: convertDetailsToProto(details),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := ic.client.UpdateIntersection(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (ic *IntersectionClient) DeleteIntersection(
	ctx context.Context,
	id string,
) (*emptypb.Empty, error) {
	req := &intersectionpb.IntersectionIDRequest{
		Id: id,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := ic.client.DeleteIntersection(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (ic *IntersectionClient) PutOptimisation(
	ctx context.Context,
	id string,
	parameters model.OptimisationParameters,
) (*intersectionpb.PutOptimisationResponse, error) {
	req := &intersectionpb.PutOptimisationRequest{
		Id:         id,
		Parameters: convertParametersToProto(parameters),
	}

	resp, err := ic.client.PutOptimisation(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

// NOTE: Creates stub for testing
type IntersectionClientInterface interface {
	CreateIntersection(
		ctx context.Context,
		intersection model.Intersection,
	) (*intersectionpb.IntersectionResponse, error)
	GetIntersection(ctx context.Context, id string) (*intersectionpb.IntersectionResponse, error)
	GetAllIntersections(
		ctx context.Context,
	) (intersectionpb.IntersectionService_GetAllIntersectionsClient, error)
	UpdateIntersection(
		ctx context.Context,
		id, name string,
		details model.Details,
	) (*intersectionpb.IntersectionResponse, error)
	DeleteIntersection(ctx context.Context, id string) (*emptypb.Empty, error)
	PutOptimisation(
		ctx context.Context,
		id string,
		parameters model.OptimisationParameters,
	) (*intersectionpb.PutOptimisationResponse, error)
}

// NOTE: Asserts Interface Implementation
var _ IntersectionClientInterface = (*IntersectionClient)(nil)

//////////////////////
// Helper functions //
//////////////////////

func convertDetailsToProto(details model.Details) *intersectionpb.IntersectionDetails {
	return &intersectionpb.IntersectionDetails{
		Address:  details.Address,
		City:     details.City,
		Province: details.Province,
	}
}

func convertParametersToProto(
	parameters model.OptimisationParameters,
) *intersectionpb.OptimisationParameters {
	return &intersectionpb.OptimisationParameters{
		OptimisationType: StringToOptimisationType(parameters.OptimisationType),
		Parameters: &intersectionpb.SimulationParameters{
			IntersectionType: StringToIntersectionType(
				parameters.SimulationParameters.IntersectionType,
			),
			Green:  int32(parameters.SimulationParameters.Green),
			Yellow: int32(parameters.SimulationParameters.Yellow),
			Red:    int32(parameters.SimulationParameters.Red),
			Speed:  int32(parameters.SimulationParameters.Speed),
			Seed:   int32(parameters.SimulationParameters.Seed),
		},
	}
}

func StringToOptimisationType(s string) intersectionpb.OptimisationType {
	switch strings.ToLower(s) {
	case "grid_search", "gridsearch":
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH
	case "genetic_evaluation", "genetic":
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_GENETIC_EVALUATION
	case "none", "":
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_NONE
	default:
		fmt.Printf("Warning: unknown optimisation type '%s', defaulting to GRIDSEARCH\n", s)
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH
	}
}

func StringToIntersectionType(s string) intersectionpb.IntersectionType {
	switch strings.ToLower(strings.ReplaceAll(s, "-", "")) {
	case "trafficlight", "traffic_light":
		return intersectionpb.IntersectionType_INTERSECTION_TYPE_TRAFFICLIGHT
	case "tjunction", "t_junction":
		return intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION
	case "roundabout":
		return intersectionpb.IntersectionType_INTERSECTION_TYPE_ROUNDABOUT
	case "stopsign", "stop_sign":
		return intersectionpb.IntersectionType_INTERSECTION_TYPE_STOP_SIGN
	case "unspecified", "":
		return intersectionpb.IntersectionType_INTERSECTION_TYPE_UNSPECIFIED
	default:
		fmt.Printf("Warning: unknown intersection type '%s', defaulting to UNSPECIFIED\n", s)
		return intersectionpb.IntersectionType_INTERSECTION_TYPE_UNSPECIFIED
	}
}

func StringToTrafficDensity(s string) intersectionpb.TrafficDensity {
	switch strings.ToLower(s) {
	case "high":
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH
	case "medium":
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM
	case "low":
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_LOW
	default:
		fmt.Printf("Warning: unknown traffic density '%s', defaulting to MEDIUM\n", s)
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM
	}
}

func StringToIntersectionStatus(s string) intersectionpb.IntersectionStatus {
	switch strings.ToLower(s) {
	case "unoptimised", "unoptimized":
		return intersectionpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED
	case "optimising", "optimizing":
		return intersectionpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISING
	case "optimised", "optimized":
		return intersectionpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISED
	case "failed":
		return intersectionpb.IntersectionStatus_INTERSECTION_STATUS_FAILED
	default:
		fmt.Printf("Warning: unknown intersection status '%s', defaulting to UNOPTIMISED\n", s)
		return intersectionpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED
	}
}
