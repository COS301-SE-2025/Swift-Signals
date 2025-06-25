package service

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
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
	//TODO:Implement CreateIntersection
	return nil, nil
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
