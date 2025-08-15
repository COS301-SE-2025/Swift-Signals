package service

import (
	"context"
	"io"
	"slices"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

type IntersectionService struct {
	intrClient client.IntersectionClientInterface
	optiClient client.OptimisationClientInterface
	userClient client.UserClientInterface
}

func NewIntersectionService(
	ic client.IntersectionClientInterface,
	oc client.OptimisationClientInterface,
	uc client.UserClientInterface,
) IntersectionServiceInterface {
	return &IntersectionService{
		intrClient: ic,
		optiClient: oc,
		userClient: uc,
	}
}

func (s *IntersectionService) GetAllIntersections(
	ctx context.Context,
	userID string,
) (model.Intersections, error) {
	logger := middleware.LoggerFromContext(ctx).With(
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
	userID string,
	intersectionID string,
) (model.Intersection, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling user service to retrieve user's intersection list")
	ids, err := s.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return model.Intersection{}, err
	}

	if !slices.Contains(ids, intersectionID) {
		return model.Intersection{}, errs.NewForbiddenError(
			"intersection not in user's intersection list",
			map[string]any{},
		)
	}

	logger.Debug("calling intersection client to get intersection")
	pbResp, err := s.intrClient.GetIntersection(ctx, intersectionID)
	if err != nil {
		return model.Intersection{}, err
	}

	resp := util.RPCIntersectionToIntersection(pbResp)
	return resp, nil
}

func (s *IntersectionService) CreateIntersection(
	ctx context.Context,
	userId string,
	req model.CreateIntersectionRequest,
) (model.CreateIntersectionResponse, error) {
	logger := middleware.LoggerFromContext(ctx).With(
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
	intrResp, err := s.intrClient.CreateIntersection(ctx, intersection)
	if err != nil {
		return model.CreateIntersectionResponse{}, err
	}

	logger.Debug("calling user client to add intersection id")
	_, err = s.userClient.AddIntersectionID(ctx, userId, intrResp.Id)
	if err != nil {
		return model.CreateIntersectionResponse{}, err
	}

	resp := model.CreateIntersectionResponse{
		Id: intrResp.Id,
	}
	return resp, nil
}

func (s *IntersectionService) UpdateIntersectionByID(
	ctx context.Context,
	userID string,
	intersectionID string,
	req model.UpdateIntersectionRequest,
) (model.Intersection, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling user service to retrieve user's intersection list")
	ids, err := s.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return model.Intersection{}, err
	}

	if !slices.Contains(ids, intersectionID) {
		return model.Intersection{}, errs.NewForbiddenError(
			"intersection not in user's intersection list",
			map[string]any{},
		)
	}

	logger.Debug("calling intersection client to update intersection")
	pbResp, err := s.intrClient.UpdateIntersection(ctx, intersectionID, req.Name, req.Details)
	if err != nil {
		return model.Intersection{}, err
	}

	resp := util.RPCIntersectionToIntersection(pbResp)
	return resp, nil
}

func (s *IntersectionService) DeleteIntersectionByID(
	ctx context.Context,
	userID string,
	intersectionID string,
) error {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling user service to retrieve user's intersection list")
	ids, err := s.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return err
	}

	if !slices.Contains(ids, intersectionID) {
		return errs.NewForbiddenError(
			"intersection not in user's intersection list",
			map[string]any{},
		)
	}

	logger.Debug("calling user client to remove intersection id")
	_, err = s.userClient.RemoveIntersectionID(ctx, userID, intersectionID)
	if err != nil {
		return err
	}

	logger.Debug("calling intersection client to delete intersection")
	_, err = s.intrClient.DeleteIntersection(ctx, intersectionID)
	return err
}

func (s *IntersectionService) OptimiseIntersectionByID(
	ctx context.Context,
	userID string,
	intersectionID string,
) error {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "intersection",
	)

	logger.Debug("calling user service to retrieve user's intersection list")
	ids, err := s.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return err
	}

	if !slices.Contains(ids, intersectionID) {
		return errs.NewForbiddenError(
			"intersection not in user's intersection list",
			map[string]any{},
		)
	}

	logger.Debug("calling intersection client to retrieve parameters")
	intersection, err := s.intrClient.GetIntersection(ctx, intersectionID)
	if err != nil {
		return err
	}

	logger.Debug("calling optimisation client to optimise intersection")
	optimisedParams, err := s.optiClient.RunOptimisation(
		ctx,
		util.RPCOptiParamToOptiParam(intersection.DefaultParameters),
	)
	if err != nil {
		return err
	}

	logger.Debug("calling intersection client to update optimised parameters")
	_, err = s.intrClient.PutOptimisation(
		ctx,
		intersectionID,
		util.RPCOptiParamToOptiParamOp(optimisedParams),
	)

	return err
}

/******************/
/* Helper Methods */
/******************/
func (s *IntersectionService) GetUserIntersectionIDs(
	ctx context.Context,
	userID string,
) ([]string, error) {
	stream, err := s.userClient.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	intersectionIDs := []string{}
	for {
		intID, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errs.NewInternalError(
				"unable to retrieve intersection IDs",
				err,
				map[string]any{},
			)
		}
		intersectionIDs = append(intersectionIDs, intID.IntersectionId)
	}
	return intersectionIDs, nil
}

/******************/

// IntersectionServiceInterface creates stub for testing
type IntersectionServiceInterface interface {
	GetAllIntersections(ctx context.Context, userID string) (model.Intersections, error)
	GetIntersectionByID(
		ctx context.Context,
		userID string,
		intersectionID string,
	) (model.Intersection, error)
	CreateIntersection(
		ctx context.Context,
		userId string,
		req model.CreateIntersectionRequest,
	) (model.CreateIntersectionResponse, error)
	UpdateIntersectionByID(
		ctx context.Context,
		userID string,
		intersectionID string,
		req model.UpdateIntersectionRequest,
	) (model.Intersection, error)
	DeleteIntersectionByID(ctx context.Context, userID string, intersectionID string) error
	OptimiseIntersectionByID(ctx context.Context, userID string, intersectionID string) error
}

// NOTE: Asserts Interface Implementation
var _ IntersectionServiceInterface = (*IntersectionService)(nil)
