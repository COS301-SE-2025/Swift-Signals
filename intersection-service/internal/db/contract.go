package db

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
)

type IntersectionRepository interface {
	CreateIntersection(
		ctx context.Context,
		intersection *model.Intersection,
	) (*model.Intersection, error)
	GetIntersectionByID(ctx context.Context, id string) (*model.Intersection, error)
	GetAllIntersections(
		ctx context.Context,
		limit, offset int,
		filter string,
	) ([]*model.Intersection, error)
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
