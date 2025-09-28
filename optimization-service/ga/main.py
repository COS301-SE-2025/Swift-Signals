import os
import random
import json
import time
from deap import tools

from ga.ga_core import toolbox, run_ga
from ga.evaluation import evaluate_balanced
from ga.simulation_client import run_simulation


def main(traffic_density: int = 2) -> dict:
    """Main function to run the genetic algorithm for optimising traffic light parameters.
    This function initialises the genetic algorithm, runs it in two phases (minimising waiting time and safety hazards),
    and saves the best parameters found to a JSON file. It also runs a final simulation with the best parameters
    and prints the results.

    Returns:
        dict: A dictionary containing the best parameters found during the optimisation process.
    """
    random.seed(1408)

    # GA parameters
    ngen = 5
    pop_size = 10
    cxpb = 0.5
    mutpb = 0.3

    pop = toolbox.population(n=pop_size)
    hof = tools.HallOfFame(1)

    start_ga = time.time()
    run_ga(
        pop,
        hof,
        ngen,
        cxpb,
        mutpb,
        lambda ind: evaluate_balanced(ind, traffic_density=traffic_density),
        label="BalancedGA",
    )
    end_ga = time.time()

    total_time_run = end_ga - start_ga
    mins = int(total_time_run // 60)
    secs = int(total_time_run % 60)

    print(f"GA ran for {mins} minutes and {secs} seconds")

    # Save and display best parameters
    os.makedirs("out", exist_ok=True)
    best = hof[0]
    best_params = {
        "Green": best[0],
        "Yellow": best[1],
        "Red": best[2],
        "Speed": best[3],
        "Seed": best[4],
        "Fitness": best.fitness.values[0],
    }

    with open("out/best_params.json", "w") as f:
        json.dump(best_params, f, indent=2)

    print("\nBest Parameters Found:")
    print(json.dumps(best_params, indent=2))

    # Run a final simulation with the best parameters and print results
    final_results = run_simulation(list(best_params.values())[:5])

    if final_results:
        print("\n--- Final Simulation Results with Optimised Parameters ---")
        for metric, value in final_results.items():
            print(f"{metric}: {value}")

    # Return the final best parameters
    return best_params


if __name__ == "__main__":
    main()
