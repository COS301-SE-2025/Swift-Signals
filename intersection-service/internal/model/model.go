package model

import (
	"errors"
)

// Common errors used throughout the application
var (
	ErrIntersectionNotFound = errors.New("intersection not found")
)

// User represents a user in the system
type Intersection struct {
}
