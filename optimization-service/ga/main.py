import os
import random
import json
from deap import tools

from ga.ga_core import toolbox, run_ga
from ga.evaluation import evaluate_waiting_and_travel, evaluate_safety_given_waiting
from ga.simulation_client import run_simulation


def main() -> dict:
    """Main function to run the genetic algorithm for optimising traffic light parameters.
    This function initialises the genetic algorithm, runs it in two phases (minimising waiting time and safety hazards),
    and saves the best parameters found to a JSON file. It also runs a final simulation with the best parameters
    and prints the results.

    Returns:
        dict: A dictionary containing the best parameters found during the optimisation process.
    """
    random.seed(1408)

    # GA parameters
    ngen_waiting = 2
    ngen_safety = 2
    pop_size = 5
    cxpb = 0.5
    mutpb = 0.3

    # Phase 1: Minimize waiting time
    print("\n--- Phase 1: Minimizing Waiting Time ---")
    pop1 = toolbox.population(n=pop_size)
    hof_wait = tools.HallOfFame(3)

    # Run GA with logging
    run_ga(
        pop1,
        hof_wait,
        ngen_waiting,
        cxpb,
        mutpb,
        evaluate_waiting_and_travel,
        label="WaitingTime",
    )

    # Phase 2: Minimize safety issues
    print("\n--- Phase 2: Minimizing Safety Hazards ---")
    pop2 = [toolbox.clone(ind) for ind in hof_wait]
    pop2 += toolbox.population(n=pop_size - len(pop2))

    hof_safety = tools.HallOfFame(1)

    # Run GA with the safety-focused evaluation
    run_ga(
        pop2,
        hof_safety,
        ngen_safety,
        cxpb,
        mutpb,
        evaluate_safety_given_waiting,
        label="Safety",
    )

    # Save and display best parameters
    os.makedirs("out", exist_ok=True)
    best = hof_safety[0]
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
