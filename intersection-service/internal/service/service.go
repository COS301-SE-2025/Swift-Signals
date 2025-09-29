package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/util"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service struct {
	repo      db.IntersectionRepository
	validator *validator.Validate
}

func NewIntersectionService(r db.IntersectionRepository) IntersectionService {
	return &Service{
		repo:      r,
		validator: validator.New(),
	}
}

func (s *Service) CreateIntersection(
	ctx context.Context,
	name string,
	details model.IntersectionDetails,
	density model.TrafficDensity,
	defaultParams model.OptimisationParameters,
) (*model.Intersection, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := CreateIntersectionRequest{
		Name:          strings.TrimSpace(name),
		Details:       details,
		Density:       density,
		DefaultParams: defaultParams,
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("creating intersection")
	id := uuid.New().String()
	createdAt := time.Now()
	lastRunAt := time.Now()
	status := model.Unoptimised
	runCount := 0

	intersection := &model.Intersection{
		ID:                id,
		Name:              name,
		Details:           details,
		CreatedAt:         createdAt,
		LastRunAt:         lastRunAt,
		Status:            status,
		RunCount:          runCount,
		TrafficDensity:    density,
		DefaultParameters: defaultParams,
		BestParameters:    defaultParams,
		CurrentParameters: defaultParams,
	}

	createdIntersection, err := s.repo.CreateIntersection(ctx, intersection)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError("failed to create intersection", err, map[string]any{})
	}

	return createdIntersection, nil
}

func (s *Service) GetIntersection(ctx context.Context, id string) (*model.Intersection, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := GetIntersectionRequest{
		ID: strings.TrimSpace(id),
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("finding intersection")
	intersection, err := s.repo.GetIntersectionByID(ctx, id)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError("failed to find intersection", err, map[string]any{})
	}
	return intersection, nil
}

func (s *Service) GetAllIntersections(
	ctx context.Context,
	page, pageSize int,
	filter string,
) ([]*model.Intersection, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := GetAllIntersectionsRequest{
		Page:     page,
		PageSize: pageSize,
		Filter:   strings.TrimSpace(filter),
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("finding all intersections")
	offset := (page - 1) * pageSize
	limit := pageSize
	intersections, err := s.repo.GetAllIntersections(ctx, limit, offset, filter)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError("failed to find all intersections", err, map[string]any{})
	}
	return intersections, nil
}

func (s *Service) UpdateIntersection(
	ctx context.Context,
	id string,
	name string,
	details model.IntersectionDetails,
	status model.IntersectionStatus,
) (*model.Intersection, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := UpdateIntersectionRequest{
		ID:      strings.TrimSpace(id),
		Name:    strings.TrimSpace(name),
		Details: details,
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("updating intersection")
	intersection, err := s.repo.UpdateIntersection(ctx, id, name, details, status)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError("failed to update intersection", err, map[string]any{})
	}
	return intersection, nil
}

func (s *Service) DeleteIntersection(ctx context.Context, id string) error {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := DeleteIntersectionRequest{
		ID: strings.TrimSpace(id),
	}
	if err := s.validator.Struct(req); err != nil {
		return handleValidationError(err)
	}

	logger.Debug("deleting intersection")
	err := s.repo.DeleteIntersection(ctx, id)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError("failed to delete intersection", err, map[string]any{})
	}
	return nil
}

func (s *Service) PutOptimisation(
	ctx context.Context,
	id string,
	params model.OptimisationParameters,
) (bool, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := PutOptimisationRequest{
		ID:     strings.TrimSpace(id),
		Params: params,
	}
	if err := s.validator.Struct(req); err != nil {
		return false, handleValidationError(err)
	}

	logger.Debug("finding existing current params")
	_, err := s.repo.GetIntersectionByID(ctx, id)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return false, err
		}
		return false, errs.NewInternalError(
			"failed to update current intersection",
			err,
			map[string]any{},
		)
	}

	logger.Debug("evaluating whether current params are better then best params")
	// TODO: Implement logic to decide whether current params are better than current
	// NOTE: For now always updating best params to current's
	better := true

	logger.Debug("updating best params")
	err = s.repo.UpdateBestParams(ctx, id, params)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return false, err
		}
		return false, errs.NewInternalError(
			"failed to update best params for intersection",
			err,
			map[string]any{},
		)
	}

	logger.Debug("updating current params")
	err = s.repo.UpdateCurrentParams(ctx, id, params)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return false, err
		}
		return false, errs.NewInternalError(
			"failed to update current params for intersection",
			err,
			map[string]any{},
		)
	}

	return better, nil
}
