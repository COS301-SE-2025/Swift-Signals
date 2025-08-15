# Optimization Service Refactoring

This document describes the refactoring of the optimization service from file-based subprocess calls to gRPC-based communication.

## Overview

The optimization service has been refactored to improve modularity, eliminate file I/O dependencies, and use gRPC for communication with the simulation service. The main changes include:

1. **Replaced subprocess calls with gRPC client**
2. **Eliminated file-based parameter passing**
3. **Encapsulated functionality in OptimizationEngine class**
4. **Return structured data instead of writing to files**
5. **Local file storage for logs and outputs**

## Key Changes

### Before (Original GP.py)
- Used `subprocess.run()` to call `SimLoad.py` directly
- Created temporary parameter files in `../simulation-service/parameters/`
- Read results from files in `../simulation-service/out/results/`
- Wrote best parameters to files
- Cleanup required for temporary files

### After (Refactored GP.py)
- Uses `SimulationClient` for gRPC communication
- No temporary files created
- Direct data exchange via protobuf messages
- Returns structured Python dictionaries
- All logs stored locally in `optimization-service/`

## Architecture

### Class Structure

```
OptimizationEngine
├── __init__(simulation_server_address)
├── run_simulation(individual)
├── evaluate_waiting_and_travel(individual)
├── evaluate_safety_given_waiting(individual)
├── run_ga(pop, hof, ngen, cxpb, mutpb, label)
├── get_final_comparison(best_params, reference_results)
└── run_optimization(**kwargs) -> dict
```

### Communication Flow

```
OptimizationEngine -> SimulationClient -> gRPC -> SimulationService
                 <-                   <-      <-
```

## Usage

### Basic Usage

```python
from GP import OptimizationEngine

# Initialize engine
engine = OptimizationEngine("localhost:50053")

# Run optimization
results = engine.run_optimization()

# Access best parameters
best_params = results["best_parameters"]
print(f"Green: {best_params['Green']}")
print(f"Yellow: {best_params['Yellow']}")
print(f"Red: {best_params['Red']}")
print(f"Speed: {best_params['Speed']}")
```

### Custom Parameters

```python
# Custom optimization parameters
results = engine.run_optimization(
    ngen_waiting=20,      # Generations for waiting time optimization
    ngen_safety=8,        # Generations for safety optimization
    pop_size=25,          # Population size
    cxpb=0.6,            # Crossover probability
    mutpb=0.25,          # Mutation probability
    random_seed=42       # Random seed for reproducibility
)
```

### gRPC Server Integration

```python
# In server/server.py
from GP import OptimizationEngine

class OptimisationServicer(pb_grpc.OptimisationServiceServicer):
    def RunOptimisation(self, request, context):
        engine = OptimizationEngine()
        results = engine.run_optimization()
        
        # Convert to protobuf response
        best_params = results["best_parameters"]
        return pb.OptimisationParameters(
            parameters=pb.SimulationParameters(
                green=best_params["Green"],
                yellow=best_params["Yellow"],
                red=best_params["Red"],
                speed=best_params["Speed"],
                seed=best_params["Seed"]
            )
        )
```

## API Reference

### OptimizationEngine

#### Constructor
```python
OptimizationEngine(simulation_server_address=None)
```
- `simulation_server_address`: gRPC server address (default: "localhost:50053")

#### Methods

##### run_optimization()
```python
run_optimization(
    ngen_waiting=30,
    ngen_safety=10,
    pop_size=30,
    cxpb=0.5,
    mutpb=0.3,
    random_seed=1408
) -> dict
```

Returns a dictionary with:
```python
{
    "parameters": {...},          # Optimization parameters used
    "phases": {
        "waiting_time": {...},    # Phase 1 results
        "safety": {...}           # Phase 2 results
    },
    "best_parameters": {
        "Green": int,
        "Yellow": int,
        "Red": int,
        "Speed": int,
        "Seed": int,
        "Fitness": float
    },
    "final_comparison": {...}     # Final simulation results
}
```

##### run_simulation()
```python
run_simulation(individual) -> dict
```
- `individual`: List [green, yellow, red, speed, seed]
- Returns simulation results dictionary or None if failed

##### evaluate_waiting_and_travel()
```python
evaluate_waiting_and_travel(individual) -> tuple
```
Returns fitness tuple for waiting time and travel time objective.

##### evaluate_safety_given_waiting()
```python
evaluate_safety_given_waiting(individual) -> tuple
```
Returns fitness tuple for safety objective (with speed constraints).

### SimulationClient

#### Constructor
```python
SimulationClient(server_address=None)
```

#### Methods

##### get_simulation_results()
```python
get_simulation_results(green, yellow, red, speed, seed, intersection_id="") -> dict
```

Returns:
```python
{
    "Total Vehicles": int,
    "Average Travel Time": float,
    "Total Travel Time": float,
    "Average Speed": float,
    "Average Waiting Time": float,
    "Total Waiting Time": float,
    "Generated Vehicles": int,
    "Emergency Brakes": int,
    "Emergency Stops": int,
    "Near collisions": int
}
```

## File Structure

```
optimization-service/
├── GP.py                       # Refactored optimization engine
├── client/
│   ├── __init__.py
│   ├── simulation.py           # gRPC simulation client
│   ├── simulation_pb2.py       # Generated protobuf
│   └── simulation_pb2_grpc.py  # Generated gRPC stubs
├── server/
│   ├── server.py               # Updated gRPC server
│   ├── optimisation_pb2.py     # Generated protobuf
│   └── optimisation_pb2_grpc.py
├── out/                        # Local output directory
├── ga_results/                 # Local GA logs
├── test_refactored_gp.py      # Test suite
├── example_usage.py           # Usage examples
└── REFACTORING_README.md      # This file
```

## Environment Variables

- `SIMULATION_SERVER_ADDRESS`: gRPC simulation service address (default: "localhost:50053")
- `APP_PORT`: Optimization service port (default: "50054")

## Testing

Run the test suite to verify functionality:

```bash
cd optimization-service
python test_refactored_gp.py
```

The test suite includes:
- SimulationClient unit tests with mocked responses
- OptimizationEngine tests with mocked simulations
- Small-scale genetic algorithm test
- Real simulation service connectivity test

## Migration Guide

### For External Callers

**Before:**
```python
# Old way - file-based
import subprocess
subprocess.run(["python3", "GP.py"])

# Read results from file
with open("out/best_parameters.json") as f:
    results = json.load(f)
```

**After:**
```python
# New way - programmatic
from GP import OptimizationEngine

engine = OptimizationEngine()
results = engine.run_optimization()
best_params = results["best_parameters"]
```

### For gRPC Integration

**Before:**
```python
# Manual file handling required
# Complex subprocess management
# Error handling via exit codes
```

**After:**
```python
# Direct class instantiation
# Exception-based error handling
# Structured return data
engine = OptimizationEngine()
try:
    results = engine.run_optimization()
    if "error" not in results:
        # Success - use results["best_parameters"]
except Exception as e:
    # Handle error
```

## Dependencies

The refactored system requires:

```txt
grpcio>=1.74.0
grpcio-tools>=1.74.0
protobuf>=6.31.1
deap>=1.4.3
tqdm>=4.67.1
```

## Performance Improvements

1. **Eliminated file I/O overhead** - Direct memory communication
2. **Removed subprocess spawning** - Native Python function calls
3. **Better error handling** - Immediate exception propagation
4. **Parallel potential** - gRPC supports concurrent requests
5. **Resource management** - Explicit connection lifecycle

## Backward Compatibility

The main `GP.py` file retains a `main()` function for standalone execution:

```bash
python GP.py  # Still works for command-line usage
```

However, the recommended approach is to use the `OptimizationEngine` class directly.

## Future Enhancements

1. **Connection pooling** for gRPC clients
2. **Async optimization** for concurrent parameter evaluation
3. **Optimization result caching** to avoid redundant simulations
4. **Real-time progress reporting** via callbacks
5. **Custom fitness function support** for different optimization objectives

## Troubleshooting

### Common Issues

1. **"Connection refused"** - Ensure simulation service is running on specified port
2. **"Module not found"** - Check protobuf files are generated and imports are correct
3. **"Fitness evaluation failed"** - Check individual parameter constraints
4. **"gRPC error"** - Verify network connectivity and service availability

### Debug Mode

Enable detailed logging:

```python
import logging
logging.basicConfig(level=logging.DEBUG)

engine = OptimizationEngine()
results = engine.run_optimization()
```

### Simulation Service Health Check

```python
from client.simulation import SimulationClient

client = SimulationClient()
try:
    result = client.get_simulation_results(30, 5, 25, 60, 1408)
    if result:
        print("Simulation service is healthy")
    else:
        print("Simulation service returned None")
except Exception as e:
    print(f"Simulation service error: {e}")
finally:
    client.close()
```
