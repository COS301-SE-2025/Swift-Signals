package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/api"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/client"
	"google.golang.org/grpc"
)

var REST_PORT = ":9090"

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // adjust for TLS if needed
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}

	userClient := client.NewUserClient(conn)

	handler := &api.AuthHandler{UserClient: userClient}
	mux := http.NewServeMux()
	api.HandlerFromMux(handler, mux)

	// Simulation Service Routes
	mux.HandleFunc("GET /simulations", getAllSimulations)
	mux.HandleFunc("GET /simulations/{id}", getSimulationByID)
	mux.HandleFunc("POST /simulations", createSimulation)
	mux.HandleFunc("DELETE /simulations/{id}", deleteSimulation)

	log.Printf("API Gateway running on %s\n", REST_PORT)
	log.Fatal(http.ListenAndServe(REST_PORT, mux))

}

// ---- Handlers ----
func getAllSimulations(rw http.ResponseWriter, r *http.Request) {
	dummy := []map[string]string{
		{"id": "1", "name": "Sim A"},
		{"id": "2", "name": "Sim B"},
	}
	writeJSON(rw, dummy)
}

func getSimulationByID(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	dummy := map[string]string{"id": id, "name": "Very real simulation"}
	writeJSON(rw, dummy)
}

func createSimulation(rw http.ResponseWriter, r *http.Request) {
	type Input struct {
		Name string `json:"name"`
	}
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}
	dummy := map[string]string{"id": "123", "name": input.Name}
	rw.WriteHeader(http.StatusCreated)
	writeJSON(rw, dummy)
}

func deleteSimulation(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	resp := map[string]string{"status": "deleted", "id": id}
	writeJSON(rw, resp)
}

// ---- Utility ----
func writeJSON(rw http.ResponseWriter, data any) {
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(data)
}
