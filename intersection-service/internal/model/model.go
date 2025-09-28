package model

import (
	"time"
)

type Intersection struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Details           IntersectionDetails    `json:"details"`
	CreatedAt         time.Time              `json:"created_at"`
	LastRunAt         time.Time              `json:"last_run_at"`
	Status            IntersectionStatus     `json:"status"`
	RunCount          int                    `json:"run_count"`
	TrafficDensity    TrafficDensity         `json:"traffic_density"`
	DefaultParameters OptimisationParameters `json:"default_parameters"`
	BestParameters    OptimisationParameters `json:"best_parameters"`
	CurrentParameters OptimisationParameters `json:"current_parameters"`
}

type IntersectionDetails struct {
	Address  string `json:"address"`
	City     string `json:"city"`
	Province string `json:"province"`
}

type IntersectionStatus string

const (
	Unspecified IntersectionStatus = "INTERSECTION_STATUS_UNSPECIFIED"
	Unoptimised IntersectionStatus = "INTERSECTION_STATUS_UNOPTIMISED"
	Optimising  IntersectionStatus = "INTERSECTION_STATUS_OPTIMISING"
	Optimised   IntersectionStatus = "INTERSECTION_STATUS_OPTIMISED"
	Failed      IntersectionStatus = "INTERSECTION_STATUS_FAILED"
)

type TrafficDensity string

const (
	TrafficLow    TrafficDensity = "TRAFFIC_DENSITY_LOW"
	TrafficMedium TrafficDensity = "TRAFFIC_DENSITY_MEDIUM"
	TrafficHigh   TrafficDensity = "TRAFFIC_DENSITY_HIGH"
)

type OptimisationParameters struct {
	OptimisationType OptimisationType     `json:"optimisation_type"`
	Parameters       SimulationParameters `json:"parameters"`
}

type OptimisationType string

const (
	OptNone              OptimisationType = "OPTIMISATION_TYPE_NONE"
	OptGridSearch        OptimisationType = "OPTIMISATION_TYPE_GRIDSEARCH"
	OptGeneticEvaluation OptimisationType = "OPTIMISATION_TYPE_GENETIC_EVALUATION"
)

type SimulationParameters struct {
	IntersectionType IntersectionType `json:"intersection_type"`
	Green            int              `json:"green"`
	Yellow           int              `json:"yellow"`
	Red              int              `json:"red"`
	Speed            int              `json:"speed"`
	Seed             int              `json:"seed"`
}

type IntersectionType string

const (
	IntersectionTrafficLight IntersectionType = "INTERSECTION_TYPE_TRAFFICLIGHT"
	IntersectionTJunction    IntersectionType = "INTERSECTION_TYPE_TJUNCTION"
	IntersectionRoundabout   IntersectionType = "INTERSECTION_TYPE_ROUNDABOUT"
	IntersectionStopSign     IntersectionType = "INTERSECTION_TYPE_STOP_SIGN"
)
