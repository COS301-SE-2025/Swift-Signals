package service

import (
	"context"
	"io"
	"slices"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

type SimulationService struct {
	intrClient client.IntersectionClientInterface
	optiClient client.OptimisationClientInterface
	userClient client.UserClientInterface
	simClient  client.SimulationClientInterface
}

func NewSimulationService(
	intrClient client.IntersectionClientInterface,
	optiClient client.OptimisationClientInterface,
	userClient client.UserClientInterface,
	simClient client.SimulationClientInterface,
) SimulationServiceInterface {
	return &SimulationService{
		intrClient: intrClient,
		optiClient: optiClient,
		userClient: userClient,
		simClient:  simClient,
	}
}

func (s *SimulationService) GetSimulationData(
	ctx context.Context,
	intersectionID string,
) (model.SimulationResponse, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "simulation",
	)

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return model.SimulationResponse{}, errs.NewInternalError(
			"user ID missing inside of handler",
			nil,
			map[string]any{},
		)
	}

	logger.Debug("calling user service to retrieve user's intersection IDs")
	intersectionIDs, err := s.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	if !slices.Contains(intersectionIDs, intersectionID) {
		return model.SimulationResponse{}, errs.NewForbiddenError(
			"you do not have access to this intersection",
			map[string]any{"intersectionID": intersectionID},
		)
	}

	logger.Debug("calling intersection service to get simulation parameters")
	intersection, err := s.intrClient.GetIntersection(ctx, intersectionID)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	defaultParams := util.RPCSimParamToSimParam(intersection.DefaultParameters.Parameters)

	logger.Debug("calling simulation service to get simulation results")
	simulationResults, err := s.simClient.GetSimulationResults(ctx, intersection.Id, defaultParams)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	logger.Debug("calling simulation service to get simulation output")
	simulationOutput, err := s.simClient.GetSimulationOutput(ctx, intersection.Id, defaultParams)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	return model.SimulationResponse{
		Results: util.RPCSimResultsToSimResults(simulationResults),
		Output:  util.RPCSimOutputToSimOutput(simulationOutput),
	}, nil
}

func (s *SimulationService) GetOptimisedData(
	ctx context.Context,
	intersectionID string,
) (model.SimulationResponse, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "simulation",
	)

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return model.SimulationResponse{}, errs.NewInternalError(
			"user ID missing inside of handler",
			nil,
			map[string]any{},
		)
	}

	logger.Debug("calling user service to retrieve user's intersection IDs")
	intersectionIDs, err := s.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	if !slices.Contains(intersectionIDs, intersectionID) {
		return model.SimulationResponse{}, errs.NewForbiddenError(
			"you do not have access to this intersection",
			map[string]any{"intersectionID": intersectionID},
		)
	}

	logger.Debug("calling intersection service to get simulation parameters")
	intersection, err := s.intrClient.GetIntersection(ctx, intersectionID)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	if intersection.BestParameters == nil {
		return model.SimulationResponse{}, errs.NewNotFoundError(
			"no optimised parameters found for this intersection",
			map[string]any{"intersectionID": intersectionID},
		)
	}

	bestParams := util.RPCSimParamToSimParam(intersection.BestParameters.Parameters)

	logger.Debug("calling simulation service to get simulation results")
	simulationResults, err := s.simClient.GetSimulationResults(ctx, intersection.Id, bestParams)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	logger.Debug("calling simulation service to get simulation output")
	simulationOutput, err := s.simClient.GetSimulationOutput(ctx, intersection.Id, bestParams)
	if err != nil {
		return model.SimulationResponse{}, err
	}

	return model.SimulationResponse{
		Results: util.RPCSimResultsToSimResults(simulationResults),
		Output:  util.RPCSimOutputToSimOutput(simulationOutput),
	}, nil
}

func (s *SimulationService) OptimiseIntersection(
	ctx context.Context,
	intersectionID string,
) (model.OptimisationResponse, error) {
	logger := middleware.LoggerFromContext(ctx).With(
		"service", "simulation",
	)

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return model.OptimisationResponse{}, errs.NewInternalError(
			"user ID missing inside of handler",
			nil,
			map[string]any{},
		)
	}

	logger.Debug("calling user service to retrieve user's intersection IDs")
	intersectionIDs, err := s.GetUserIntersectionIDs(ctx, userID)
	if err != nil {
		return model.OptimisationResponse{}, err
	}

	if !slices.Contains(intersectionIDs, intersectionID) {
		return model.OptimisationResponse{}, errs.NewForbiddenError(
			"you do not have access to this intersection",
			map[string]any{"intersectionID": intersectionID},
		)
	}

	logger.Debug("calling intersection service to get intersection details")
	intersection, err := s.intrClient.GetIntersection(ctx, intersectionID)
	if err != nil {
		return model.OptimisationResponse{}, err
	}

	logger.Debug(
		"calling intersection service to change status to 'INTERSECTION_STATUS_OPTIMISING'",
	)
	_, err = s.intrClient.UpdateIntersectionStatus(
		ctx,
		intersection.Id,
		intersection.Name,
		util.RPCDetailsToDetails(intersection.Details),
		commonpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISING,
	)
	if err != nil {
		logger.Warn("Could not update intersection status to 'INTERSECTION_STATUS_OPTIMISING'")
	}

	resultChan := make(chan model.OptimisationResponse, 1)
	errChan := make(chan error, 1)

	go func() {
		ctx := context.Background()
		logger.Debug("calling optimisation service to optimise intersection")
		response, err := s.optiClient.RunOptimisation(
			ctx,
			util.RPCOptiParamToOptiParam(intersection.DefaultParameters),
		)
		if err != nil {
			errChan <- err
			return
		}

		logger.Debug("updating intersection with optimised parameters")
		resp, err := s.intrClient.PutOptimisation(
			ctx,
			intersectionID,
			util.RPCOptiParamToOptiParamOp(response),
		)
		if err != nil {
			errChan <- err
			return
		}
		logger.Debug(
			"calling intersection service to change status to 'INTERSECTION_STATUS_OPTIMISED'",
		)
		_, err = s.intrClient.UpdateIntersectionStatus(
			ctx,
			intersection.Id,
			intersection.Name,
			util.RPCDetailsToDetails(intersection.Details),
			commonpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISED,
		)
		if err != nil {
			logger.Warn("Could not update intersection status to 'INTERSECTION_STATUS_OPTIMISED'")
		}
		resultChan <- model.OptimisationResponse{Improved: resp.Improved}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return model.OptimisationResponse{}, err
	case <-ctx.Done():
		return model.OptimisationResponse{}, ctx.Err()
	}
}

func (s *SimulationService) GetUserIntersectionIDs(
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

// SimulationServiceInterface creates stub for testing
type SimulationServiceInterface interface {
	GetSimulationData(ctx context.Context, intersectionID string) (model.SimulationResponse, error)
	GetOptimisedData(ctx context.Context, intersectionID string) (model.SimulationResponse, error)
	OptimiseIntersection(
		ctx context.Context,
		intersectionID string,
	) (model.OptimisationResponse, error)
}

// NOTE: Asserts the SimulationService implements the SimulationServiceInterface
var _ SimulationServiceInterface = (*SimulationService)(nil)
