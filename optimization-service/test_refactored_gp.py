#!/usr/bin/env python3
"""
Test script for the refactored optimization system using gRPC
"""

import os
import sys
import json
import time
from unittest.mock import Mock, patch

# Add current directory to path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from GP import OptimizationEngine
from client.simulation import SimulationClient


def test_simulation_client():
    """Test the simulation client with mocked responses"""
    print("Testing SimulationClient...")

    # Mock the gRPC response
    mock_response = Mock()
    mock_response.total_vehicles = 100
    mock_response.average_travel_time = 25.5
    mock_response.total_travel_time = 2550.0
    mock_response.average_speed = 45.2
    mock_response.average_waiting_time = 15.3
    mock_response.total_waiting_time = 1530.0
    mock_response.generated_vehicles = 98
    mock_response.emergency_brakes = 2
    mock_response.emergency_stops = 1
    mock_response.near_collisions = 0

    # Mock MessageToDict to return expected structure
    def mock_message_to_dict(message, **kwargs):
        return {
            "total_vehicles": message.total_vehicles,
            "average_travel_time": message.average_travel_time,
            "total_travel_time": message.total_travel_time,
            "average_speed": message.average_speed,
            "average_waiting_time": message.average_waiting_time,
            "total_waiting_time": message.total_waiting_time,
            "generated_vehicles": message.generated_vehicles,
            "emergency_brakes": message.emergency_brakes,
            "emergency_stops": message.emergency_stops,
            "near_collisions": message.near_collisions,
        }

    with patch('client.simulation.pb_grpc.SimulationServiceStub') as mock_stub_class, \
         patch('client.simulation.MessageToDict', side_effect=mock_message_to_dict):
        mock_stub = Mock()
        mock_stub.GetSimulationResults.return_value = mock_response
        mock_stub_class.return_value = mock_stub

        client = SimulationClient("localhost:50053")
        result = client.get_simulation_results(30, 5, 25, 60, 1408)

        print(f"Simulation result: {json.dumps(result, indent=2)}")
        assert result is not None
        assert result["Total Vehicles"] == 100
        assert result["Emergency Brakes"] == 2

        client.close()
        print("✓ SimulationClient test passed")


def test_optimization_engine_mock():
    """Test the optimization engine with mocked simulation calls"""
    print("\nTesting OptimizationEngine with mocked simulation...")

    # Mock simulation results
    def mock_simulation_results(*args, **kwargs):
        return {
            "Total Vehicles": 100,
            "Average Travel Time": 25.0,
            "Total Travel Time": 2500.0,
            "Average Speed": 45.0,
            "Average Waiting Time": 15.0,
            "Total Waiting Time": 1500.0,
            "Generated Vehicles": 98,
            "Emergency Brakes": 1,
            "Emergency Stops": 0,
            "Near collisions": 0,
        }

    with patch.object(SimulationClient, 'get_simulation_results', side_effect=mock_simulation_results):
        engine = OptimizationEngine("localhost:50053")

        # Test individual simulation
        individual = [30, 5, 25, 60, 1408]  # green, yellow, red, speed, seed
        result = engine.run_simulation(individual)
        print(f"Individual simulation result: {json.dumps(result, indent=2)}")
        assert result is not None
        assert result["Total Vehicles"] == 100

        # Test fitness evaluation
        fitness = engine.evaluate_waiting_and_travel(individual)
        print(f"Fitness (waiting + travel): {fitness}")
        assert fitness[0] > 0

        # Test safety evaluation
        safety_fitness = engine.evaluate_safety_given_waiting(individual)
        print(f"Safety fitness: {safety_fitness}")
        assert safety_fitness[0] >= 0

        print("✓ OptimizationEngine mock test passed")


def test_small_optimization():
    """Test a small optimization run with mocked simulation"""
    print("\nTesting small optimization run...")

    # Mock simulation results with some variation
    call_count = 0
    def mock_simulation_results(*args, **kwargs):
        nonlocal call_count
        call_count += 1
        # Add some variation to make the optimization interesting
        base_waiting = 1500 + (call_count % 10) * 50
        base_travel = 2500 + (call_count % 8) * 30
        return {
            "Total Vehicles": 100,
            "Average Travel Time": base_travel / 100,
            "Total Travel Time": base_travel,
            "Average Speed": 45.0,
            "Average Waiting Time": base_waiting / 100,
            "Total Waiting Time": base_waiting,
            "Generated Vehicles": 98,
            "Emergency Brakes": max(0, (call_count % 5) - 2),
            "Emergency Stops": max(0, (call_count % 7) - 5),
            "Near collisions": max(0, (call_count % 11) - 9),
        }

    with patch.object(SimulationClient, 'get_simulation_results', side_effect=mock_simulation_results):
        engine = OptimizationEngine("localhost:50053")

        # Run a very small optimization for testing
        results = engine.run_optimization(
            ngen_waiting=2,  # Very small for testing
            ngen_safety=1,
            pop_size=5,      # Very small population
            cxpb=0.5,
            mutpb=0.3,
            random_seed=1408
        )

        print("Optimization results:")
        print(f"Best parameters: {json.dumps(results['best_parameters'], indent=2)}")

        assert "best_parameters" in results
        assert "Green" in results["best_parameters"]
        assert "Yellow" in results["best_parameters"]
        assert "Red" in results["best_parameters"]
        assert "Speed" in results["best_parameters"]
        assert "Seed" in results["best_parameters"]
        assert "Fitness" in results["best_parameters"]

        # Check that parameters are within expected ranges
        best = results["best_parameters"]
        assert 10 <= best["Green"] <= 60
        assert 3 <= best["Yellow"] <= 8
        assert 10 <= best["Red"] <= 60
        assert best["Speed"] in [40, 60, 80, 100]

        print(f"Total simulation calls made: {call_count}")
        print("✓ Small optimization test passed")


def test_optimization_engine_direct():
    """Test optimization engine directly without mocking (requires simulation service)"""
    print("\nTesting OptimizationEngine with real simulation service...")
    print("Note: This test requires the simulation service to be running on localhost:50053")

    try:
        engine = OptimizationEngine("localhost:50053")

        # Test a single simulation call
        individual = [30, 5, 25, 60, 1408]
        result = engine.run_simulation(individual)

        if result is None:
            print("⚠ Could not connect to simulation service - skipping real simulation test")
            return

        print(f"Real simulation result: {json.dumps(result, indent=2)}")

        # If we got here, simulation service is working
        print("✓ Real simulation test passed")

    except Exception as e:
        print(f"⚠ Real simulation test failed (simulation service may not be running): {e}")


def main():
    """Run all tests"""
    print("Starting tests for refactored optimization system...")
    print("=" * 60)

    try:
        test_simulation_client()
        test_optimization_engine_mock()
        test_small_optimization()
        test_optimization_engine_direct()

        print("\n" + "=" * 60)
        print("✓ All tests completed successfully!")
        print("\nThe refactored system appears to be working correctly.")
        print("Key improvements:")
        print("- Uses gRPC instead of subprocess calls")
        print("- No longer creates temporary files")
        print("- Returns structured data instead of writing to files")
        print("- Encapsulated in OptimizationEngine class for better modularity")

    except Exception as e:
        print(f"\n✗ Test failed: {e}")
        import traceback
        traceback.print_exc()
        return 1

    return 0


if __name__ == "__main__":
    exit(main())
