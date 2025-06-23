package main

import (
    "context"
    "encoding/json"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "os"
    "os/exec"
)

func runSimulation(simID string) error {
    // 1. Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://Admin:Capstone2025@simulationservice.dnm1bu7.mongodb.net/"))
    if err != nil {
        return err
    }
    defer client.Disconnect(context.TODO())

    db := client.Database("SimulationService")
    collection := db.Collection("Intersections")

    // 2. Fetch simulation parameters
    var result map[string]interface{}
    err = collection.FindOne(context.TODO(), bson.M{"intersection.id": simID}).Decode(&result)
    if err != nil {
        return err
    }

    simParams, err := json.Marshal(result)
    if err != nil {
        return err
    }

    // 3. Save to temp file
    err = os.WriteFile("temp_params.json", simParams, 0644)
    if err != nil {
        return err
    }

    // 4. Run Python code (use py if python isn't recognized)
    cmd := exec.Command("py", "SimLoad.py", "--params", "temp_params.json")
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Simulation error: %s", string(output))
        return err
    }
    log.Printf("Simulation output: %s", string(output))

    // 5. Load results
    resultData, err := os.ReadFile("out/results/simulation_results.json")
    if err != nil {
        return err
    }

    // 6. Update MongoDB with results
    _, err = collection.UpdateOne(context.TODO(),
        bson.M{"intersection.id": simID},
        bson.M{"$set": bson.M{"intersection.simulation_results": json.RawMessage(resultData)}},
    )
    return err
}

func main() {
    simID := "simId" // replace with a real simId from your DB
    err := runSimulation(simID)
    if err != nil {
        log.Fatalf("Simulation failed: %v", err)
    } else {
        log.Println("Simulation completed and results updated.")
    }
}
