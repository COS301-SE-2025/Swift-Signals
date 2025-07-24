package model

import (
	"errors"
	"time"
)

// Common errors used throughout the application
var (
	ErrIntersectionNotFound = errors.New("intersection not found")
	ErrIntersectionExists   = errors.New("intersection already exists")
	ErrInvalidParameters    = errors.New("invalid parameters")
)

// IntersectionStatus enum
type IntersectionStatus int

const (
	Unoptimised IntersectionStatus = iota
	Optimising
	Optimised
	Failed
)

// TrafficDensity enum
type TrafficDensity int

const (
	TrafficLow TrafficDensity = iota
	TrafficMedium
	TrafficHigh
)

// OptimisationType enum
type OptimisationType int

const (
	OptNone OptimisationType = iota
	OptGridSearch
	OptGeneticEvaluation
)

// IntersectionType enum
type IntersectionType int

const (
	IntersectionUnspecified IntersectionType = iota
	IntersectionTrafficLight
	IntersectionTJunction
	IntersectionRoundabout
	IntersectionStopSign
)

type IntersectionResponse struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Details           IntersectionDetails    `json:"details"`
	CreatedAt         time.Time              `json:"created_at"`
	LastRunAt         time.Time              `json:"last_run_at"`
	Status            IntersectionStatus     `json:"status"`
	RunCount          int32                  `json:"run_count"`
	TrafficDensity    TrafficDensity         `json:"traffic_density"`
	DefaultParameters OptimisationParameters `json:"default_parameters"`
	BestParameters    OptimisationParameters `json:"best_parameters"`
	CurrentParameters OptimisationParameters `json:"current_parameters"`
}

type PutOptimisationResponse struct {
	Improved bool `json:"improved"`
}

type IntersectionDetails struct {
	Address  string `json:"address"`
	City     string `json:"city"`
	Province string `json:"province"`
}

type OptimisationParameters struct {
	OptimisationType OptimisationType     `json:"optimisation_type"`
	Parameters       SimulationParameters `json:"parameters"`
}

type SimulationParameters struct {
	IntersectionType IntersectionType `json:"intersection_type"`
	Green            int32            `json:"green"`
	Yellow           int32            `json:"yellow"`
	Red              int32            `json:"red"`
	Speed            int32            `json:"speed"`
	Seed             int32            `json:"seed"`
}

// To pretty print the enum IntersectionStatus
func (s IntersectionStatus) String() string {
	switch s {
	case Unoptimised:
		return "Unoptimised"
	case Optimising:
		return "Optimising"
	case Optimised:
		return "Optimised"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// To pretty print the enum TrafficDensity
func (s TrafficDensity) String() string {
	switch s {
	case TrafficLow:
		return "Traffic Density Low"
	case TrafficMedium:
		return "Traffic Density Medium"
	case TrafficHigh:
		return "Traffic Density High"
	default:
		return "Unknown"
	}
}

// To pretty print the enum OptimisationType
func (s OptimisationType) String() string {
	switch s {
	case OptNone:
		return "Optimisation Type None"
	case OptGridSearch:
		return "Optimisation Type Gridsearch"
	case OptGeneticEvaluation:
		return "Optimisation Type Genetic Evaluation"
	default:
		return "Unknown"
	}
}

// To pretty print the enum IntersectionType
func (s IntersectionType) String() string {
	switch s {
	case IntersectionUnspecified:
		return "Intersection Type Unspecified"
	case IntersectionTrafficLight:
		return "Intersection Type Traffic Light"
	case IntersectionTJunction:
		return "Intersection Type TJunction"
	case IntersectionRoundabout:
		return "Intersection Type Roundabout"
	case IntersectionStopSign:
		return "Intersection Type Stop Sign"
	default:
		return "Unknown"
	}
}
