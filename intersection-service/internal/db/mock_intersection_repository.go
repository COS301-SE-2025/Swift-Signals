package db

import (
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
