package db

import (
	"context"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the UserRepository interface
type MockRepository struct {
	mock.Mock
}

func NewUserRepository() UserRepository {
	return &MockRepository{
		Mock: mock.Mock{},
	}
}

// CreateUser mocks the CreateUser method
func (m *MockRepository) CreateUser(ctx context.Context, user *model.UserResponse) (*model.UserResponse, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

// GetUserByID mocks the GetUserByID method
func (m *MockRepository) GetUserByID(ctx context.Context, id string) (*model.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

// GetUserByEmail mocks the GetUserByEmail method
func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*model.UserResponse, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockRepository) UpdateUser(ctx context.Context, user *model.UserResponse) (*model.UserResponse, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

// DeleteUser mocks the DeleteUser method
func (m *MockRepository) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListUsers mocks the ListUsers method
func (m *MockRepository) ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.UserResponse), args.Error(1)
}

// Verify MockRepository implements UserRepository interface at compile time
var _ UserRepository = (*MockRepository)(nil)
