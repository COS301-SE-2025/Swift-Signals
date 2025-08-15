#!/usr/bin/env python3
"""
Example usage of the refactored optimization system

This script demonstrates how to use the OptimizationEngine class
to run genetic algorithm optimization using gRPC calls to the simulation service.
"""

import os
import sys
import json
import time

# Add current directory to path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from GP import OptimizationEngine


def example_basic_optimization():
    """
    Example 1: Basic optimization with default parameters
    """
    print("=" * 60)
    print("Example 1: Basic Optimization")
    print("=" * 60)

    # Initialize the optimization engine
    # The simulation service should be running on localhost:50053
    engine = OptimizationEngine()

    try:
        # Run optimization with default parameters
        results = engine.run_optimization()

        if "error" in results:
            print(f"Optimization failed: {results['error']}")
            return None

        print("\nOptimization Results:")
        print(f"Best Parameters: {json.dumps(results['best_parameters'], indent=2)}")

        # Show final comparison if available
        if results.get('final_comparison', {}).get('comparison'):
            print("\nFinal Comparison:")
            for metric, data in results['final_comparison']['comparison'].items():
                opt_val = data.get('optimized', 'N/A')
                ref_val = data.get('reference', 'N/A')
                improvement = data.get('improvement', 'N/A')
                print(f"{metric}: Optimized={opt_val}, Reference={ref_val}, Improvement={improvement}")

        return results

    except Exception as e:
        print(f"Error during optimization: {e}")
        return None


def example_custom_optimization():
    """
    Example 2: Custom optimization with specific parameters
    """
    print("\n" + "=" * 60)
    print("Example 2: Custom Optimization Parameters")
    print("=" * 60)

    # Initialize with custom simulation server address
    simulation_server = os.environ.get("SIMULATION_SERVER_ADDRESS", "localhost:50053")
    engine = OptimizationEngine(simulation_server_address=simulation_server)

    try:
        # Run optimization with custom parameters
        results = engine.run_optimization(
            ngen_waiting=10,     # Fewer generations for faster execution
            ngen_safety=5,       # Fewer safety generations
            pop_size=15,         # Smaller population
            cxpb=0.7,           # Higher crossover probability
            mutpb=0.2,          # Lower mutation probability
            random_seed=42      # Different random seed
        )

        if "error" in results:
            print(f"Optimization failed: {results['error']}")
            return None

        print("\nCustom Optimization Results:")
        print(f"Parameters used: {json.dumps(results['parameters'], indent=2)}")
        print(f"Best Parameters: {json.dumps(results['best_parameters'], indent=2)}")

        return results

    except Exception as e:
        print(f"Error during custom optimization: {e}")
        return None


def example_single_simulation():
    """
    Example 3: Single simulation call
    """
    print("\n" + "=" * 60)
    print("Example 3: Single Simulation Call")
    print("=" * 60)

    engine = OptimizationEngine()

    try:
        # Test a single simulation with specific parameters
        individual = [30, 5, 25, 60, 1408]  # green, yellow, red, speed, seed
        print(f"Testing individual: Green={individual[0]}, Yellow={individual[1]}, "
              f"Red={individual[2]}, Speed={individual[3]}, Seed={individual[4]}")

        result = engine.run_simulation(individual)

        if result is None:
            print("Simulation failed - check if simulation service is running")
            return None

        print("\nSimulation Results:")
        print(json.dumps(result, indent=2))

        # Test fitness evaluations
        waiting_fitness = engine.evaluate_waiting_and_travel(individual)
        safety_fitness = engine.evaluate_safety_given_waiting(individual)

        print(f"\nFitness Evaluations:")
        print(f"Waiting & Travel Fitness: {waiting_fitness[0]:.2f}")
        print(f"Safety Fitness: {safety_fitness[0]:.2f}")

        return result

    except Exception as e:
        print(f"Error during single simulation: {e}")
        return None


def example_batch_evaluation():
    """
    Example 4: Batch evaluation of multiple parameter sets
    """
    print("\n" + "=" * 60)
    print("Example 4: Batch Parameter Evaluation")
    print("=" * 60)

    engine = OptimizationEngine()

    # Define multiple parameter sets to test
    parameter_sets = [
        {"name": "Conservative", "params": [45, 4, 40, 40, 1408]},
        {"name": "Moderate", "params": [30, 5, 25, 60, 1408]},
        {"name": "Aggressive", "params": [20, 3, 15, 80, 1408]},
        {"name": "High Speed", "params": [25, 4, 20, 100, 1408]},
    ]

    results = []

    try:
        for param_set in parameter_sets:
            name = param_set["name"]
            params = param_set["params"]

            print(f"\nTesting {name} configuration: {params}")

            result = engine.run_simulation(params)
            if result is None:
                print(f"  Failed to get results for {name}")
                continue

            # Calculate fitness scores
            waiting_fitness = engine.evaluate_waiting_and_travel(params)[0]
            safety_fitness = engine.evaluate_safety_given_waiting(params)[0]

            evaluation = {
                "name": name,
                "parameters": {
                    "Green": params[0],
                    "Yellow": params[1],
                    "Red": params[2],
                    "Speed": params[3],
                    "Seed": params[4]
                },
                "results": result,
                "fitness": {
                    "waiting_travel": waiting_fitness,
                    "safety": safety_fitness
                }
            }

            results.append(evaluation)

            print(f"  Total Waiting Time: {result.get('Total Waiting Time', 'N/A')}")
            print(f"  Emergency Brakes: {result.get('Emergency Brakes', 'N/A')}")
            print(f"  Fitness Score: {waiting_fitness:.2f}")

        # Summary
        print(f"\n{'Configuration':<15}{'Waiting Time':<15}{'Emergency Brakes':<18}{'Fitness Score':<15}")
        print("-" * 65)

        for r in results:
            name = r["name"]
            waiting = r["results"].get("Total Waiting Time", "N/A")
            brakes = r["results"].get("Emergency Brakes", "N/A")
            fitness = r["fitness"]["waiting_travel"]
            print(f"{name:<15}{str(waiting):<15}{str(brakes):<18}{fitness:<15.2f}")

        return results

    except Exception as e:
        print(f"Error during batch evaluation: {e}")
        return None


def main():
    """
    Main function to run all examples
    """
    print("Optimization Service Integration Examples")
    print("========================================")
    print("\nNote: These examples require the simulation service to be running")
    print("on localhost:50053 (or the address specified in SIMULATION_SERVER_ADDRESS)")
    print("\nYou can run specific examples by uncommenting them below.")

    # Example 1: Basic optimization (uncomment to run)
    # basic_results = example_basic_optimization()

    # Example 2: Custom optimization (uncomment to run)
    # custom_results = example_custom_optimization()

    # Example 3: Single simulation (recommended for testing)
    single_results = example_single_simulation()

    # Example 4: Batch evaluation (uncomment to run)
    # batch_results = example_batch_evaluation()

    print("\n" + "=" * 60)
    print("Examples completed!")
    print("\nTo run the full optimization, uncomment the desired examples in main()")
    print("and ensure the simulation service is running.")


if __name__ == "__main__":
    main()
