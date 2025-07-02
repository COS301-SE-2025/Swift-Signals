package util

import (
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
)

func RPCIntersectionToIntersection(rpc *intersection.IntersectionResponse) model.Intersection {
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

func RPCDetailsToDetails(rpc *intersection.IntersectionDetails) model.Details {
	return model.Details{
		Address:  rpc.Address,
		City:     rpc.City,
		Province: rpc.Province,
	}
}

func RPCOptiParamToOptiParam(rpc *intersection.OptimisationParameters) model.OptimisationParameters {
	return model.OptimisationParameters{
		OptimisationType:     rpc.OptimisationType.String(),
		SimulationParameters: RPCSimParamToSimParam(rpc.Parameters),
	}
}

func RPCSimParamToSimParam(rpc *intersection.SimulationParameters) model.SimulationParameters {
	return model.SimulationParameters{
		IntersectionType: rpc.IntersectionType.String(),
		Red:              int(rpc.Red),
		Yellow:           int(rpc.Yellow),
		Green:            int(rpc.Green),
		Speed:            int(rpc.Speed),
		Seed:             int(rpc.Seed),
	}
}
