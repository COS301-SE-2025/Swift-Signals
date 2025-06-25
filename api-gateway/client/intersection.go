package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"google.golang.org/grpc"
)

type IntersectionClient struct {
	client intersectionpb.IntersectionServiceClient
}

func NewIntersectionClient(conn *grpc.ClientConn) *IntersectionClient {
	return &IntersectionClient{
		client: intersectionpb.NewIntersectionServiceClient(conn),
	}
}

func (ic *IntersectionClient) CreateIntersection(ctx context.Context, intersection model.Intersection) (*intersectionpb.IntersectionResponse, error) {
	req := &intersectionpb.CreateIntersectionRequest{
		Name:           intersection.Name,
		Details:        convertDetailsToProto(intersection.Details),
		TrafficDensity: StringToTrafficDensity(intersection.TrafficDensity),
		DefaultParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: StringToOptimisationType(intersection.DefaultParameters.OptimisationType),
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: StringToIntersectionType(intersection.DefaultParameters.SimulationParameters.IntersectionType),
				Green:            int32(intersection.DefaultParameters.SimulationParameters.Green),
				Yellow:           int32(intersection.DefaultParameters.SimulationParameters.Yellow),
				Red:              int32(intersection.DefaultParameters.SimulationParameters.Red),
				Speed:            int32(intersection.DefaultParameters.SimulationParameters.Speed),
				Seed:             int32(intersection.DefaultParameters.SimulationParameters.Seed),
			},
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return ic.client.CreateIntersection(ctx, req)
}

func convertDetailsToProto(details model.Details) *intersectionpb.IntersectionDetails {
	return &intersectionpb.IntersectionDetails{
		Address:  details.Address,
		City:     details.City,
		Province: details.Province,
	}
}

// StringToOptimisationType converts a string to protobuf OptimisationType enum
func StringToOptimisationType(s string) intersectionpb.OptimisationType {
	switch strings.ToLower(s) {
	case "grid_search", "gridsearch":
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH
	case "genetic_evaluation", "genetic":
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_GENETIC_EVALUATION
	case "none", "":
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_NONE
	default:
		// Log warning or handle error as needed
		fmt.Printf("Warning: unknown optimisation type '%s', defaulting to GRIDSEARCH\n", s)
		return intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH
	}
}

// StringToIntersectionType converts a string to protobuf IntersectionType enum
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
		// Log warning or handle error as needed
		fmt.Printf("Warning: unknown intersection type '%s', defaulting to UNSPECIFIED\n", s)
		return intersectionpb.IntersectionType_INTERSECTION_TYPE_UNSPECIFIED
	}
}

// StringToTrafficDensity converts a string to protobuf TrafficDensity enum
func StringToTrafficDensity(s string) intersectionpb.TrafficDensity {
	switch strings.ToLower(s) {
	case "high":
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH
	case "medium":
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM
	case "low":
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_LOW
	default:
		// Log warning or handle error as needed
		fmt.Printf("Warning: unknown traffic density '%s', defaulting to MEDIUM\n", s)
		return intersectionpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM
	}
}

// StringToIntersectionStatus converts a string to protobuf IntersectionStatus enum
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
		// Log warning or handle error as needed
		fmt.Printf("Warning: unknown intersection status '%s', defaulting to UNOPTIMISED\n", s)
		return intersectionpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED
	}
}
