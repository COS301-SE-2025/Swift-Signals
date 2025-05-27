package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAllSimulations(t *testing.T) {
	// Create a GET request to "/simulations"
	req := httptest.NewRequest("GET", "/simulations", nil)

	// Create a response recorder to capture handler output
	rr := httptest.NewRecorder()

	// Call the handler function directly
	getAllSimulations(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	// Check content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	// Check response body (optional deep validation)
	var data []map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Errorf("response is not valid JSON: %v", err)
	}
	if len(data) != 2 {
		t.Errorf("expected 2 simulations, got %d", len(data))
	}
}

func TestGetSimulationByID(t *testing.T) {
	req := httptest.NewRequest("GET", "/simulations/42", nil)

	// Inject a path parameter manually (Go 1.22+ supports r.PathValue)
	req.SetPathValue("id", "42")

	rr := httptest.NewRecorder()
	getSimulationByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}

	var data map[string]string
	json.Unmarshal(rr.Body.Bytes(), &data)

	if data["id"] != "42" {
		t.Errorf("expected id '42', got '%s'", data["id"])
	}
}

func TestCreateSimulation(t *testing.T) {
	payload := `{"name":"Test Sim"}`
	req := httptest.NewRequest("POST", "/simulations", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	createSimulation(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", rr.Code)
	}

	var data map[string]string
	json.Unmarshal(rr.Body.Bytes(), &data)

	if data["name"] != "Test Sim" {
		t.Errorf("expected name 'Test Sim', got '%s'", data["name"])
	}
	if data["id"] != "123" {
		t.Errorf("expected id '123', got '%s'", data["id"])
	}
}

func TestDeleteSimulation(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/simulations/55", nil)
	req.SetPathValue("id", "55")

	rr := httptest.NewRecorder()
	deleteSimulation(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}

	var data map[string]string
	json.Unmarshal(rr.Body.Bytes(), &data)

	if data["status"] != "deleted" || data["id"] != "55" {
		t.Errorf("unexpected response: %v", data)
	}
}
