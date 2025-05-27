package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var port = ":9090"

func main() {

	mux := http.NewServeMux()

	// Simulation Service Routes
	mux.HandleFunc("GET /simulations", getAllSimulations)
	mux.HandleFunc("GET /simulations/{id}", getSimulationByID)
	mux.HandleFunc("POST /simulations", createSimulation)
	mux.HandleFunc("DELETE /simulations/{id}", deleteSimulation)

	log.Println("API Gateway running on :9090")
	log.Fatal(http.ListenAndServe(":9090", mux))

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
