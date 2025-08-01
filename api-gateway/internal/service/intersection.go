package service

import (
	"context"
	"errors"
	"io"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
)

type IntersectionServiceInterface interface {
	GetAllIntersections(ctx context.Context) (model.Intersections, error)
	GetIntersectionByID(ctx context.Context, id string) (model.Intersection, error)
	CreateIntersection(ctx context.Context, req model.CreateIntersectionRequest) (model.CreateIntersectionResponse, error)
	UpdateIntersectionByID(ctx context.Context, id string, req model.UpdateIntersectionRequest) (model.Intersection, error)
}

type IntersectionService struct {
	intrClient *client.IntersectionClient
}

func NewIntersectionService(ic *client.IntersectionClient) IntersectionServiceInterface {
	return &IntersectionService{
		intrClient: ic,
	}
}

func (s *IntersectionService) GetAllIntersections(ctx context.Context) (model.Intersections, error) {
	stream, err := s.intrClient.GetAllIntersections(ctx)
	if err != nil {
		return model.Intersections{}, errors.New("Unable to get all intersections")
	}

	intersections := []model.Intersection{}
	for {
		rpcIntersection, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return model.Intersections{}, errors.New("Unable to get all intersections")
		}
		intersection := util.RPCIntersectionToIntersection(rpcIntersection)
		intersections = append(intersections, intersection)
	}

	resp := model.Intersections{Intersections: intersections}
	return resp, nil
}

func (s *IntersectionService) GetIntersectionByID(ctx context.Context, id string) (model.Intersection, error) {
	pbResp, err := s.intrClient.GetIntersection(ctx, id)
	if err != nil {
		return model.Intersection{}, errors.New("Unable to get intersection by ID")
	}

	resp := util.RPCIntersectionToIntersection(pbResp)

	return resp, nil
}

func (s *IntersectionService) CreateIntersection(ctx context.Context, req model.CreateIntersectionRequest) (model.CreateIntersectionResponse, error) {
	intersection := model.Intersection{
		Name:           req.Name,
		Details:        req.Details,
		TrafficDensity: req.TrafficDensity,
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: req.DefaultParameters,
		},
	}
	pbResp, err := s.intrClient.CreateIntersection(ctx, intersection)
	if err != nil {
		return model.CreateIntersectionResponse{}, errors.New("Unable to create intersection")
	}

	resp := model.CreateIntersectionResponse{
		Id: pbResp.Id,
	}

	return resp, nil
}

func (s *IntersectionService) UpdateIntersectionByID(ctx context.Context, id string, req model.UpdateIntersectionRequest) (model.Intersection, error) {
	return model.Intersection{}, nil
}
