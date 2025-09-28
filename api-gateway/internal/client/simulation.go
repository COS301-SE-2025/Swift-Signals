package client

import (
	"context"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	simulationpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/simulation/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"google.golang.org/grpc"
)

type SimulationClient struct {
	client simulationpb.SimulationServiceClient
}

func NewSimulationClient(client simulationpb.SimulationServiceClient) *SimulationClient {
	return &SimulationClient{
		client: client,
	}
}

func NewSimulationClientFromConn(conn *grpc.ClientConn) *SimulationClient {
	return NewSimulationClient(simulationpb.NewSimulationServiceClient(conn))
}

func (sc *SimulationClient) GetSimulationResults(
	ctx context.Context,
	id string,
	simulation_parameters model.SimulationParameters,
) (*simulationpb.SimulationResultsResponse, error) {
	intersection, ok := commonpb.IntersectionType_value[simulation_parameters.IntersectionType]

	if !ok {
		return nil, errs.NewValidationError("invalid simulation parameters", map[string]any{})
	}

	req := &simulationpb.SimulationRequest{
		IntersectionId: id,
		SimulationParameters: &commonpb.SimulationParameters{
			IntersectionType: commonpb.IntersectionType(intersection),
			Green:            int32(simulation_parameters.Green),
			Yellow:           int32(simulation_parameters.Yellow),
			Red:              int32(simulation_parameters.Red),
			Speed:            int32(simulation_parameters.Speed),
			Seed:             int32(simulation_parameters.Seed),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := sc.client.GetSimulationResults(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (sc *SimulationClient) GetSimulationOutput(
	ctx context.Context,
	id string,
	simulation_parameters model.SimulationParameters,
) (*simulationpb.SimulationOutputResponse, error) {
	intersection, ok := commonpb.IntersectionType_value[simulation_parameters.IntersectionType]

	if !ok {
		return nil, errs.NewValidationError("invalid simulation parameters", map[string]any{})
	}

	req := &simulationpb.SimulationRequest{
		IntersectionId: id,
		SimulationParameters: &commonpb.SimulationParameters{
			IntersectionType: commonpb.IntersectionType(intersection),
			Green:            int32(simulation_parameters.Green),
			Yellow:           int32(simulation_parameters.Yellow),
			Red:              int32(simulation_parameters.Red),
			Speed:            int32(simulation_parameters.Speed),
			Seed:             int32(simulation_parameters.Seed),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := sc.client.GetSimulationOutput(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

// NOTE: Creates stub for testing
type SimulationClientInterface interface {
	GetSimulationResults(
		ctx context.Context,
		id string,
		simulation_parameters model.SimulationParameters,
	) (*simulationpb.SimulationResultsResponse, error)
	GetSimulationOutput(
		ctx context.Context,
		id string,
		simulation_parameters model.SimulationParameters,
	) (*simulationpb.SimulationOutputResponse, error)
}

// NOTE: Asserts the SimulationClient implements the SimulationClientInterface
var _ SimulationClientInterface = (*SimulationClient)(nil)
