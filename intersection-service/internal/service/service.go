package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/google/uuid"
)

type Service struct {
	repo db.IntersectionRepository
}

func NewService(r db.IntersectionRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateIntersection(
	ctx context.Context,
	name string,
	details model.IntersectionDetails,
	density model.TrafficDensity,
	defaultParams model.OptimisationParameters,
) (*model.IntersectionResponse, error) {
	// Validate Input
	if err := validateCreateIntersectionInput(name, details, density, defaultParams); err != nil {
		return nil, err
	}

	// Create Intersection object
	id := uuid.New().String()
	createdAt := time.Now()
	lastRunAt := time.Now()
	status := model.Unoptimised
	runCount := int32(0)

	intersection := &model.IntersectionResponse{
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

	// Store Intersection in repo
	createdIntersection, err := s.repo.CreateIntersection(ctx, intersection)
	if err != nil {
		return nil, fmt.Errorf("failed to create intersection: %w", err)
	}

	// Return Created Intersection
	return createdIntersection, nil
}

func (s *Service) GetIntersection(ctx context.Context, id string) (*model.IntersectionResponse, error) {
	// TODO: Implement GetIntersection
	return nil, nil
}

func (s *Service) GetAllIntersections(ctx context.Context) ([]*model.IntersectionResponse, error) {
	//TODO: Implement GetAllIntersections
	return nil, nil
}

func (s *Service) UpdateIntersection(
	ctx context.Context,
	id string,
	name string,
	details model.IntersectionDetails,
) (*model.IntersectionResponse, error) {
	//TODO:Implement UpdateIntersection
	return nil, nil
}

func (s *Service) DeleteIntersection(ctx context.Context, id string) error {
	//TODO: DeleteIntersection
	return nil
}

func (s *Service) PutOptimisation(
	ctx context.Context,
	id string,
	params model.OptimisationParameters,
) (*model.PutOptimisationResponse, error) {
	//TODO:Implement PutOptimisation
	return nil, nil
}

////////////////////////
// Validation Helpers //
////////////////////////

func validateCreateIntersectionInput(name string, details model.IntersectionDetails, density model.TrafficDensity, params model.OptimisationParameters) error {
	var validationErrors []string

	// Validate name
	if strings.TrimSpace(name) == "" {
		validationErrors = append(validationErrors, "intersection name is required and cannot be empty")
	}
	if len(name) > 255 {
		validationErrors = append(validationErrors, "intersection name cannot exceed 255 characters")
	}

	// Validate intersection details
	if strings.TrimSpace(details.Address) == "" {
		validationErrors = append(validationErrors, "address is required")
	}
	if strings.TrimSpace(details.City) == "" {
		validationErrors = append(validationErrors, "city is required")
	}
	if strings.TrimSpace(details.Province) == "" {
		validationErrors = append(validationErrors, "province is required")
	}

	// Validate traffic density enum
	if density < model.TrafficLow || density > model.TrafficHigh {
		validationErrors = append(validationErrors, "invalid traffic density value")
	}

	// Validate optimisation parameters
	if params.OptimisationType < model.OptNone || params.OptimisationType > model.OptGeneticEvaluation {
		validationErrors = append(validationErrors, "invalid optimisation type")
	}

	// Validate simulation parameters
	if params.Parameters.IntersectionType < model.IntersectionUnspecified || params.Parameters.IntersectionType > model.IntersectionStopSign {
		validationErrors = append(validationErrors, "invalid intersection type")
	}

	// Validate timing parameters (must be positive)
	if params.Parameters.Green <= 0 {
		validationErrors = append(validationErrors, "green light duration must be positive")
	}
	if params.Parameters.Yellow <= 0 {
		validationErrors = append(validationErrors, "yellow light duration must be positive")
	}
	if params.Parameters.Red <= 0 {
		validationErrors = append(validationErrors, "red light duration must be positive")
	}
	if params.Parameters.Speed <= 0 {
		validationErrors = append(validationErrors, "speed must be positive")
	}

	// Validate reasonable timing ranges
	if params.Parameters.Green > 300 { // 5 minutes max
		validationErrors = append(validationErrors, "green light duration cannot exceed 300 seconds")
	}
	if params.Parameters.Yellow > 10 { // 10 seconds max
		validationErrors = append(validationErrors, "yellow light duration cannot exceed 10 seconds")
	}
	if params.Parameters.Red > 300 { // 5 minutes max
		validationErrors = append(validationErrors, "red light duration cannot exceed 300 seconds")
	}
	if params.Parameters.Speed > 200 { // 200 km/h max (reasonable for simulation)
		validationErrors = append(validationErrors, "speed cannot exceed 200 km/h")
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(validationErrors, "; "))
	}

	return nil
}
