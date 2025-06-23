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
