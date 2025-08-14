package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
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
	// TODO: add validator logic

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
		return nil, fmt.Errorf("failed to create intersection: %w", err)
	}

	return createdIntersection, nil
}

func (s *Service) GetIntersection(ctx context.Context, id string) (*model.Intersection, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	intersection, err := s.repo.GetIntersectionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing intersection: %w", err)
	}
	return intersection, nil
}

func (s *Service) GetAllIntersections(ctx context.Context) ([]*model.Intersection, error) {
	intersections, err := s.repo.GetAllIntersections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch intersections: %w", err)
	}
	return intersections, nil
}

func (s *Service) UpdateIntersection(
	ctx context.Context,
	id string,
	name string,
	details model.IntersectionDetails,
) (*model.Intersection, error) {
	// TODO:Implement UpdateIntersection
	return nil, nil
}

func (s *Service) DeleteIntersection(ctx context.Context, id string) error {
	// TODO: DeleteIntersection
	return nil
}

func (s *Service) PutOptimisation(
	ctx context.Context,
	id string,
	params model.OptimisationParameters,
) (bool, error) {
	// TODO:Implement PutOptimisation
	return true, nil
}
