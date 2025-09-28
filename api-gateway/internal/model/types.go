package model

import "time"

type Intersection struct {
	ID                string                 `json:"id"                 example:"1"`
	Name              string                 `json:"name"               example:"My Intersection"`
	Details           Details                `json:"details"`
	CreatedAt         time.Time              `json:"created_at"         example:"2025-06-24T15:04:05Z"`
	LastRunAt         time.Time              `json:"last_run_at"        example:"2025-06-24T15:04:05Z"`
	Status            string                 `json:"status"             example:"unoptimised"`
	RunCount          int                    `json:"run_count"          example:"7"`
	TrafficDensity    string                 `json:"traffic_density"    example:"high"`
	DefaultParameters OptimisationParameters `json:"default_parameters"`
	BestParameters    OptimisationParameters `json:"best_parameters"`
	CurrentParameters OptimisationParameters `json:"current_parameters"`
}

type Intersections struct {
	Intersections []Intersection `json:"intersections"`
}

type Details struct {
	Address  string `json:"address"  example:"Corner of Foo and Bar"`
	City     string `json:"city"     example:"Pretoria"`
	Province string `json:"province" example:"Gauteng"`
}

type OptimisationParameters struct {
	OptimisationType     string               `json:"optimisation_type"     example:"grid_search"`
	SimulationParameters SimulationParameters `json:"simulation_parameters"`
}

type SimulationParameters struct {
	IntersectionType string `json:"intersection_type" example:"t-junction" validate:"required"`
	Green            int    `json:"green"             example:"10"         validate:"required,min=1"`
	Yellow           int    `json:"yellow"            example:"2"          validate:"required,min=1"`
	Red              int    `json:"red"               example:"6"          validate:"required,min=1"`
	Speed            int    `json:"speed"             example:"60"         validate:"required,min=1"`
	Seed             int    `json:"seed"              example:"3247128304" validate:"required"`
}

type User struct {
	ID              string   `json:"id"               example:"1"`
	Username        string   `json:"username"         example:"johndoe"`
	Email           string   `json:"email"            example:"user@example.com"`
	IsAdmin         bool     `json:"is_admin"         example:"false"`
	IntersectionIDs []string `json:"intersection_ids" example:"[1,2,3]"`
}
