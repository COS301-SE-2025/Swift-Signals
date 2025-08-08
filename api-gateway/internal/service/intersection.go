package service

import (
	"context"
	"io"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

type IntersectionService struct {
	intrClient *client.IntersectionClient
}

func NewIntersectionService(ic *client.IntersectionClient) IntersectionServiceInterface {
	return &IntersectionService{
		intrClient: ic,
	}
}

func (s *IntersectionService) GetAllIntersections(
	ctx context.Context,
) (model.Intersections, error) {
	logger := util.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("starting grpc stream")
	stream, err := s.intrClient.GetAllIntersections(ctx)
	if err != nil {
		return model.Intersections{}, err
	}

	intersections := []model.Intersection{}
	for {
		rpcIntersection, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return model.Intersections{}, errs.NewInternalError(
				"unable to get all intersections",
				err,
				map[string]any{},
			)
		}
		intersection := util.RPCIntersectionToIntersection(rpcIntersection)
		intersections = append(intersections, intersection)
	}

	resp := model.Intersections{Intersections: intersections}
	return resp, nil
}

func (s *IntersectionService) GetIntersectionByID(
	ctx context.Context,
	id string,
) (model.Intersection, error) {
	logger := util.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling intersection client to get intersection")
	pbResp, err := s.intrClient.GetIntersection(ctx, id)
	if err != nil {
		return model.Intersection{}, err
	}

	resp := util.RPCIntersectionToIntersection(pbResp)
	return resp, nil
}

func (s *IntersectionService) CreateIntersection(
	ctx context.Context,
	req model.CreateIntersectionRequest,
) (model.CreateIntersectionResponse, error) {
	logger := util.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling intersection client to create intersection")
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
		return model.CreateIntersectionResponse{}, err
	}

	resp := model.CreateIntersectionResponse{
		Id: pbResp.Id,
	}
	return resp, nil
}

func (s *IntersectionService) UpdateIntersectionByID(
	ctx context.Context,
	id string,
	req model.UpdateIntersectionRequest,
) (model.Intersection, error) {
	logger := util.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling intersection client to update intersection")
	pbResp, err := s.intrClient.UpdateIntersection(ctx, id, req.Name, req.Details)
	if err != nil {
		return model.Intersection{}, err
	}

	resp := util.RPCIntersectionToIntersection(pbResp)
	return resp, nil
}

func (s *IntersectionService) DeleteIntersectionByID(
	ctx context.Context,
	id string,
) error {
	logger := util.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling intersection client to update intersection")
	_, err := s.intrClient.DeleteIntersection(ctx, id)
	return err
}

// IntersectionServiceInterface creates stub for testing
type IntersectionServiceInterface interface {
	GetAllIntersections(ctx context.Context) (model.Intersections, error)
	GetIntersectionByID(ctx context.Context, id string) (model.Intersection, error)
	CreateIntersection(
		ctx context.Context,
		req model.CreateIntersectionRequest,
	) (model.CreateIntersectionResponse, error)
	UpdateIntersectionByID(
		ctx context.Context,
		id string,
		req model.UpdateIntersectionRequest,
	) (model.Intersection, error)
	DeleteIntersectionByID(ctx context.Context, id string) error
}

// Note: Asserts Interface Implementation
