package model

type CreateIntersectionRequest struct {
	Name    string `json:"name"               example:"My Intersection" binding:"required" validate:"required,max=256"`
	Details struct {
		Address  string `json:"address"  example:"Corner of Foo and Bar"`
		City     string `json:"city"     example:"Pretoria"`
		Province string `json:"province" example:"Gauteng"`
	} `json:"details"`
	TrafficDensity    string               `json:"traffic_density"    example:"high"`
	DefaultParameters SimulationParameters `json:"default_parameters"                           binding:"required" validate:"required"`
}

type CreateIntersectionResponse struct {
	Id string `json:"id" example:"2"`
}

type UpdateIntersectionRequest struct {
	Name    string `json:"name"    example:"My Updated Intersection"`
	Details struct {
		Address  string `json:"address"  example:"Corner of Foo and Bar"`
		City     string `json:"city"     example:"Pretoria"`
		Province string `json:"province" example:"Gauteng"`
	} `json:"details"`
}
