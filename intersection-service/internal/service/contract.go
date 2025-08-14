package service

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
)

type IntersectionService interface {
	CreateIntersection(
		ctx context.Context,
		name string,
		details model.IntersectionDetails,
		density model.TrafficDensity,
		defaultParams model.OptimisationParameters,
	) (*model.Intersection, error)
	GetIntersection(ctx context.Context, id string) (*model.Intersection, error)
	GetAllIntersections(ctx context.Context) ([]*model.Intersection, error)
	UpdateIntersection(
		ctx context.Context,
		id string,
		name string,
		details model.IntersectionDetails,
	) (*model.Intersection, error)
	DeleteIntersection(ctx context.Context, id string) error
	PutOptimisation(
		ctx context.Context,
		id string,
		params model.OptimisationParameters,
	) (bool, error)
}
