package service

import "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"

type IntersectionService struct{}

func NewIntersectionService() *IntersectionService {
	return &IntersectionService{}
}

func (s *IntersectionService) GetAllIntersections() (resp model.Intersections, err error) {
	return model.Intersections{}, nil
}

func (s *IntersectionService) GetIntersectionByID(id int) (resp model.Intersection, err error) {
	return model.Intersection{}, nil
}

func (s *IntersectionService) CreateIntersection(req model.CreateIntersectionRequest) (resp model.CreateIntersectionResponse, err error) {
	return model.CreateIntersectionResponse{}, nil
}

func (s *IntersectionService) UpdateIntersectionByID(id int, req model.UpdateIntersectionRequest) (resp model.Intersection, err error) {
	return model.Intersection{}, nil
}
