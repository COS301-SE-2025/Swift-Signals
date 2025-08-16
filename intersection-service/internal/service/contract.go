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
	GetAllIntersections(
		ctx context.Context,
		page, pageSize int,
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

type CreateIntersectionRequest struct {
	Name          string                       `validate:"required,min=2,max=100" json:"name"`
	Details       model.IntersectionDetails    `validate:"required"               json:"details"`
	Density       model.TrafficDensity         `validate:"required"               json:"density"`
	DefaultParams model.OptimisationParameters `validate:"required"               json:"default_params"`
}

type GetIntersectionRequest struct {
	ID string `validate:"required,uuid4" json:"id"`
}

type GetAllIntersectionsRequest struct {
	Page     int    `validate:"min=1"         json:"page"`
	PageSize int    `validate:"min=1,max=100" json:"page_size"`
	Filter   string `validate:"max=255"       json:"filter"`
}

type UpdateIntersectionRequest struct {
	ID      string                    `validate:"required,uuid4"         json:"id"`
	Name    string                    `validate:"required,min=2,max=100" json:"name"`
	Details model.IntersectionDetails `validate:"required"               json:"details"`
}

type DeleteIntersectionRequest struct {
	ID string `validate:"required,uuid4" json:"id"`
}

type PutOptimisationRequest struct {
	ID     string                       `validate:"required,uuid4" json:"id"`
	Params model.OptimisationParameters `validate:"required"       json:"params"`
}
