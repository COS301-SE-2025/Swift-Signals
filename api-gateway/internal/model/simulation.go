package model

type NodeType string

const (
	NodeTypePriority     NodeType = "PRIORITY"
	NodeTypeTrafficLight NodeType = "TRAFFIC_LIGHT"
)

type SimulationResponse struct {
	Results SimulationResults `json:"results"`
	Output  SimulationOutput  `json:"output"`
}

type SimulationResults struct {
	TotalVehicles      int     `json:"total_vehicles"       example:"100"`
	AverageTravelTime  float64 `json:"average_travel_time"  example:"300"`
	TotalTravelTime    float64 `json:"total_travel_time"    example:"30000"`
	AverageSpeed       float64 `json:"average_speed"        example:"50.0"`
	AverageWaitingTime float64 `json:"average_waiting_time" example:"60.0"`
	TotalWaitingTime   float64 `json:"total_waiting_time"   example:"6000"`
	GeneratedVehicles  int     `json:"generated_vehicles"   example:"120"`
	EmergencyBrakes    int     `json:"emergency_brakes"     example:"5"`
	EmergencyStops     int     `json:"emergency_stops"      example:"3"`
	NearCollisions     int     `json:"near_collisions"      example:"2"`
}

type SimulationOutput struct {
	Intersection SimulationIntersection `json:"intersection"`
	Vehicles     []SimulationVehicle    `json:"vehicles"`
}

type SimulationIntersection struct {
	Nodes         []SimulationNode         `json:"nodes"`
	Edges         []SimulationEdge         `json:"edges"`
	Connections   []SimulationConn         `json:"connections"`
	TrafficLights []SimulationTrafficLight `json:"trafficLights"`
}

type SimulationNode struct {
	ID   string   `json:"id"   example:"node1"`
	X    float64  `json:"x"    example:"100.0"`
	Y    float64  `json:"y"    example:"200.0"`
	Type NodeType `json:"type" example:"PRIORITY"`
}

type SimulationEdge struct {
	ID    string  `json:"id"    example:"edge1"`
	From  string  `json:"from"  example:"node1"`
	To    string  `json:"to"    example:"node2"`
	Speed float64 `json:"speed" example:"50.0"`
	Lanes int     `json:"lanes" example:"2"`
}

type SimulationConn struct {
	From     string `json:"from"     example:"node1"`
	To       string `json:"to"       example:"node2"`
	FromLane int    `json:"fromLane" example:"1"`
	ToLane   int    `json:"toLane"   example:"1"`
	TL       int    `json:"tl"`
}

type SimulationTrafficLight struct {
	ID     string            `json:"id"     example:"tl1"`
	Type   string            `json:"type"   example:"FIXED"`
	Phases []SimulationPhase `json:"phases"`
}

type SimulationPhase struct {
	Duration int    `json:"duration" example:"30"`
	State    string `json:"state"    example:"GREEN"`
}

type SimulationVehicle struct {
	ID        string     `json:"id"        example:"vehicle1"`
	Positions []Position `json:"positions"`
}

type Position struct {
	Time  int     `json:"time"  example:"0"`
	X     float64 `json:"x"     example:"100.0"`
	Y     float64 `json:"y"     example:"200.0"`
	Speed float64 `json:"speed" example:"50.0"`
}
