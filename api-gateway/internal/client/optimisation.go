package client

import (
	"context"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	optimisationpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/optimisation"
	"google.golang.org/grpc"
)

type OptimisationClient struct {
	client optimisationpb.OptimisationServiceClient
}

func NewOptimisationClient(client optimisationpb.OptimisationServiceClient) *OptimisationClient {
	return &OptimisationClient{
		client: client,
	}
}

func NewOptimisationClientFromConn(conn *grpc.ClientConn) *OptimisationClient {
	return NewOptimisationClient(optimisationpb.NewOptimisationServiceClient(conn))
}

func (oc *OptimisationClient) RunOptimisation(
	ctx context.Context,
	params model.OptimisationParameters,
) (*optimisationpb.OptimisationParameters, error) {
	req := &optimisationpb.OptimisationParameters{
		OptimisationType: optimisationpb.OptimisationType(
			optimisationpb.OptimisationType_value[params.OptimisationType],
		),
		Parameters: &optimisationpb.SimulationParameters{
			IntersectionType: optimisationpb.IntersectionType(
				optimisationpb.IntersectionType_value[params.SimulationParameters.IntersectionType],
			),
			Green:  int32(params.SimulationParameters.Green),
			Yellow: int32(params.SimulationParameters.Yellow),
			Red:    int32(params.SimulationParameters.Red),
			Speed:  int32(params.SimulationParameters.Speed),
			Seed:   int32(params.SimulationParameters.Seed),
		},
	}
	ctx, cancel := context.WithTimeout(
		ctx,
		5*time.Hour,
	) // NOTE: This time depends on the hardware and should be adjusted accordingly
	defer cancel()

	resp, err := oc.client.RunOptimisation(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

// NOTE: Creates stub for testing
type OptimisationClientInterface interface {
	RunOptimisation(
		ctx context.Context,
		params model.OptimisationParameters,
	) (*optimisationpb.OptimisationParameters, error)
}

// NOTE: Asserts Interface Implementation
var _ OptimisationClientInterface = (*OptimisationClient)(nil)
