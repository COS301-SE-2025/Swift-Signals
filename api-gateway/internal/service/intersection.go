package service

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
)

type IntersectionService struct {
	intrClient *client.IntersectionClient
}

func NewIntersectionService(ic *client.IntersectionClient) *IntersectionService {
	return &IntersectionService{
		intrClient: ic,
	}
}

func (s *IntersectionService) GetAllIntersections(ctx context.Context) (resp model.Intersections, err error) {
	return model.Intersections{}, nil
}

func (s *IntersectionService) GetIntersectionByID(ctx context.Context, id int) (resp model.Intersection, err error) {
	return model.Intersection{}, nil
}

func (s *IntersectionService) CreateIntersection(ctx context.Context, req model.CreateIntersectionRequest) (resp model.CreateIntersectionResponse, err error) {
	return model.CreateIntersectionResponse{}, nil
}

func (s *IntersectionService) UpdateIntersectionByID(ctx context.Context, id int, req model.UpdateIntersectionRequest) (resp model.Intersection, err error) {
	return model.Intersection{}, nil
}
