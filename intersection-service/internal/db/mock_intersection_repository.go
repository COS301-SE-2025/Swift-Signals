package db

import (
	"context"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the UserRepository interface
type MockRepository struct {
	mock.Mock
}

func NewIntersectionRepository() IntersectionRepository {
	return &MockRepository{
		Mock: mock.Mock{},
	}
}

func (m *MockRepository) CreateIntersection(ctx context.Context, intersection *model.IntersectionResponse) (*model.IntersectionResponse, error) {
	args := m.Called(ctx, intersection)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.IntersectionResponse), args.Error(1)
}

func (m *MockRepository) GetIntersectionByID(ctx context.Context, id string) (*model.IntersectionResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.IntersectionResponse), args.Error(1)
}

func (m *MockRepository) GetAllIntersections(ctx context.Context) ([]*model.IntersectionResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.IntersectionResponse), args.Error(1)
}
