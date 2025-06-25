package model

import "time"

type Intersection struct {
	ID      string `json:"id" example:"1"`
	Name    string `json:"name" example:"My Intersection"`
	Details struct {
		Address  string `json:"address"  example:"Corner of Foo and Bar"`
		City     string `json:"city"     example:"Pretoria"`
		Province string `json:"province" example:"Gauteng"`
	} `json:"details"`
	CreatedAt         time.Time              `json:"created_at"      example:"2025-06-24T15:04:05Z"`
	LastRunAt         time.Time              `json:"last_run_at"     example:"2025-06-24T15:04:05Z"`
	Status            string                 `json:"status"          example:"unoptimised"`
	RunCount          int                    `json:"run_count"       example:"7"`
	TrafficDensity    string                 `json:"traffic_density" example:"high"`
	DefaultParameters OptimisationParameters `json:"default_parameters"`
	BestParameters    OptimisationParameters `json:"best_parameters"`
	CurrentParameters OptimisationParameters `json:"current_parameters"`
}

type Intersections struct {
	Intersections []Intersection `json:"intersections"`
}

type OptimisationParameters struct {
	OptimisationType     string               `json:"optimisation_type" example:"grid_search"`
	SimulationParameters SimulationParameters `json:"simulation_parameters"`
}

type SimulationParameters struct {
	IntersectionType string `json:"intersection_type" example:"t-junction"`
	Green            int    `json:"green"             example:"10"`
	Yellow           int    `json:"yellow"            example:"2"`
	Red              int    `json:"red"               example:"6"`
	Speed            int    `json:"speed"             example:"60"`
	Seed             int    `json:"seed"              example:"3247128304"`
}
