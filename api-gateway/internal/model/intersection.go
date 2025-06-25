package model

type CreateIntersectionRequest struct {
	Name    string `json:"name" example:"My Intersection"`
	Details struct {
		Address  string `json:"address"  example:"Corner of Foo and Bar"`
		City     string `json:"city"     example:"Pretoria"`
		Province string `json:"province" example:"Gauteng"`
	} `json:"details"`
	TrafficDensity    string               `json:"traffic_density" example:"high"`
	DefaultParameters SimulationParameters `json:"default_parameters"`
}

type CreateIntersectionResponse struct {
	Id string `json:"id" example:"2"`
}

type UpdateIntersectionRequest struct {
	Name    string `json:"name" example:"My Updated Intersection"`
	Details struct {
		Address  string `json:"address"  example:"Corner of Foo and Bar"`
		City     string `json:"city"     example:"Pretoria"`
		Province string `json:"province" example:"Gauteng"`
	} `json:"details"`
}
