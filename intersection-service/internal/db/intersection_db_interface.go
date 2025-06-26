package db

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
)

// Repository interface defines the contract for user data operations
type IntersectionRepository interface {
	CreateIntersection(ctx context.Context, intersection *model.IntersectionResponse) (*model.IntersectionResponse, error)
	GetIntersectionByID(ctx context.Context, id string) (*model.IntersectionResponse, error)
	GetAllIntersections(ctx context.Context) ([]*model.IntersectionResponse, error)
}
