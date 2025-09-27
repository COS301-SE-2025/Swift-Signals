package util

import (
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
	simulationpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/simulation/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RPCIntersectionToIntersection(rpc *intersectionpb.IntersectionResponse) model.Intersection {
	return model.Intersection{
		ID:                rpc.Id,
		Name:              rpc.Name,
		Details:           RPCDetailsToDetails(rpc.Details),
		CreatedAt:         rpc.CreatedAt.AsTime(),
		LastRunAt:         rpc.LastRunAt.AsTime(),
		Status:            rpc.Status.String(),
		RunCount:          int(rpc.RunCount),
		TrafficDensity:    rpc.TrafficDensity.String(),
		DefaultParameters: RPCOptiParamToOptiParam(rpc.DefaultParameters),
		BestParameters:    RPCOptiParamToOptiParam(rpc.BestParameters),
		CurrentParameters: RPCOptiParamToOptiParam(rpc.CurrentParameters),
	}
}

func RPCDetailsToDetails(rpc *intersectionpb.IntersectionDetails) model.Details {
	return model.Details{
		Address:  rpc.Address,
		City:     rpc.City,
		Province: rpc.Province,
	}
}

func RPCOptiParamToOptiParam(
	rpc *commonpb.OptimisationParameters,
) model.OptimisationParameters {
	return model.OptimisationParameters{
		OptimisationType:     rpc.OptimisationType.String(),
		SimulationParameters: RPCSimParamToSimParam(rpc.Parameters),
	}
}

func RPCOptiParamToOptiParamOp(
	rpc *commonpb.OptimisationParameters,
) model.OptimisationParameters {
	return model.OptimisationParameters{
		OptimisationType:     rpc.OptimisationType.String(),
		SimulationParameters: RPCSimParamToSimParamOp(rpc.Parameters),
	}
}

func RPCSimParamToSimParam(rpc *commonpb.SimulationParameters) model.SimulationParameters {
	return model.SimulationParameters{
		IntersectionType: rpc.IntersectionType.String(),
		Red:              int(rpc.Red),
		Yellow:           int(rpc.Yellow),
		Green:            int(rpc.Green),
		Speed:            int(rpc.Speed),
		Seed:             int(rpc.Seed),
	}
}

func RPCSimParamToSimParamOp(rpc *commonpb.SimulationParameters) model.SimulationParameters {
	return model.SimulationParameters{
		IntersectionType: rpc.IntersectionType.String(),
		Red:              int(rpc.Red),
		Yellow:           int(rpc.Yellow),
		Green:            int(rpc.Green),
		Speed:            int(rpc.Speed),
		Seed:             int(rpc.Seed),
	}
}

func RPCSimResultsToSimResults(
	rpc *simulationpb.SimulationResultsResponse,
) model.SimulationResults {
	return model.SimulationResults{
		TotalVehicles:      int(rpc.TotalVehicles),
		AverageTravelTime:  float64(rpc.AverageTravelTime),
		TotalTravelTime:    float64(rpc.TotalTravelTime),
		AverageSpeed:       float64(rpc.AverageSpeed),
		AverageWaitingTime: float64(rpc.AverageWaitingTime),
		TotalWaitingTime:   float64(rpc.TotalWaitingTime),
		GeneratedVehicles:  int(rpc.GeneratedVehicles),
		EmergencyBrakes:    int(rpc.EmergencyBrakes),
		EmergencyStops:     int(rpc.EmergencyStops),
		NearCollisions:     int(rpc.NearCollisions),
	}
}

func RPCSimOutputToSimOutput(rpc *simulationpb.SimulationOutputResponse) model.SimulationOutput {
	return model.SimulationOutput{
		Intersection: RPCSimIntersectionToSimIntersection(rpc.Intersection),
		Vehicles:     RPCSimVehiclesToSimVehicles(rpc.Vehicles),
	}
}

func RPCSimIntersectionToSimIntersection(
	rpc *simulationpb.Intersection,
) model.SimulationIntersection {
	return model.SimulationIntersection{
		Nodes:         RPCSimNodesToSimNodes(rpc.Nodes),
		Edges:         RPCSimEdgesToSimEdges(rpc.Edges),
		Connections:   RPCSimConnsToSimConns(rpc.Connections),
		TrafficLights: RPCSimTrafficLightsToSimTrafficLights(rpc.TrafficLights),
	}
}

func RPCSimNodesToSimNodes(rpc []*simulationpb.Node) []model.SimulationNode {
	nodes := make([]model.SimulationNode, len(rpc))
	for i, n := range rpc {
		nodes[i] = model.SimulationNode{
			ID:   n.Id,
			X:    float64(n.X),
			Y:    float64(n.Y),
			Type: model.NodeType(n.Type.String()),
		}
	}
	return nodes
}

func RPCSimEdgesToSimEdges(rpc []*simulationpb.Edge) []model.SimulationEdge {
	edges := make([]model.SimulationEdge, len(rpc))
	for i, e := range rpc {
		edges[i] = model.SimulationEdge{
			ID:    e.Id,
			From:  e.From,
			To:    e.To,
			Speed: float64(e.Speed),
			Lanes: int(e.Lanes),
		}
	}
	return edges
}

func RPCSimConnsToSimConns(rpc []*simulationpb.Connection) []model.SimulationConn {
	conns := make([]model.SimulationConn, len(rpc))
	for i, c := range rpc {
		conns[i] = model.SimulationConn{
			From:     c.From,
			To:       c.To,
			FromLane: int(c.FromLane),
			ToLane:   int(c.ToLane),
			TL:       int(c.Tl),
		}
	}
	return conns
}

func RPCSimTrafficLightsToSimTrafficLights(
	rpc []*simulationpb.TrafficLight,
) []model.SimulationTrafficLight {
	tls := make([]model.SimulationTrafficLight, len(rpc))
	for i, tl := range rpc {
		phases := make([]model.SimulationPhase, len(tl.Phases))
		for j, p := range tl.Phases {
			phases[j] = model.SimulationPhase{
				Duration: int(p.Duration),
				State:    p.State,
			}
		}
		tls[i] = model.SimulationTrafficLight{
			ID:     tl.Id,
			Type:   tl.Type,
			Phases: phases,
		}
	}
	return tls
}

func RPCSimVehiclesToSimVehicles(rpc []*simulationpb.Vehicle) []model.SimulationVehicle {
	vehicles := make([]model.SimulationVehicle, len(rpc))
	for i, v := range rpc {
		positions := make([]model.Position, len(v.Positions))
		for j, p := range v.Positions {
			positions[j] = model.Position{
				Time:  int(p.Time),
				X:     float64(p.X),
				Y:     float64(p.Y),
				Speed: float64(p.Speed),
			}
		}
		vehicles[i] = model.SimulationVehicle{
			ID:        v.Id,
			Positions: positions,
		}
	}
	return vehicles
}

func GrpcErrorToErr(err error) *errs.ServiceError {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return errs.NewValidationError(err.Error(), map[string]any{})
	case codes.NotFound:
		return errs.NewNotFoundError(err.Error(), map[string]any{})
	case codes.AlreadyExists:
		return errs.NewAlreadyExistsError(err.Error(), map[string]any{})
	case codes.Unauthenticated:
		return errs.NewUnauthorizedError(err.Error(), map[string]any{})
	case codes.PermissionDenied:
		return errs.NewForbiddenError(err.Error(), map[string]any{})
	default:
		return errs.NewInternalError(err.Error(), err, map[string]any{})
	}
}
