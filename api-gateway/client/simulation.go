package client

import (
	"context"
	"errors"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	simulationpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/simulation"
	"google.golang.org/grpc"
)

type SimulationClient struct {
	client simulationpb.SimulationServiceClient
}

func NewSimulationClient(conn *grpc.ClientConn) *SimulationClient {
	return &SimulationClient{
		client: simulationpb.NewSimulationServiceClient(conn),
	}
}

func (sc *SimulationClient) GetSimulationResults(ctx context.Context, id string, simulation_parameters model.SimulationParameters) (*simulationpb.SimulationResultsResponse, error) {

	intersection, ok := simulationpb.IntersectionType_value[simulation_parameters.IntersectionType]

	if !ok {
		return nil, errors.New("invalid intersection type")
	}

	req := &simulationpb.SimulationRequest{
		IntersectionId: id,
		SimulationParameters: &simulationpb.SimulationParameters{
			IntersectionType: simulationpb.IntersectionType(intersection),
			Green:            int32(simulation_parameters.Green),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return sc.client.GetSimulationResults(ctx, req)
}

func (sc *SimulationClient) GetSimulationOutput(ctx context.Context, id string, simulation_parameters model.SimulationParameters) (*simulationpb.SimulationOutputResponse, error) {

	intersection, ok := simulationpb.IntersectionType_value[simulation_parameters.IntersectionType]

	if !ok {
		return nil, errors.New("invalid intersection type")
	}

	req := &simulationpb.SimulationRequest{
		IntersectionId: id,
		SimulationParameters: &simulationpb.SimulationParameters{
			IntersectionType: simulationpb.IntersectionType(intersection),
			Green:            int32(simulation_parameters.Green),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return sc.client.GetSimulationOutput(ctx, req)
}
