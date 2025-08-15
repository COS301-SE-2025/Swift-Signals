import os
import random
import json
from datetime import datetime
from deap import base, creator, tools
from tqdm import tqdm

from client.simulation import SimulationClient

# Local paths for optimization service
BEST_PARAM_OUTPUT = "out/best_parameters.json"
ALL_RESULTS_CSV = "ga_results/all_individuals_log.csv"

# Fitness setup
creator.create("FitnessMin", base.Fitness, weights=(-1.0,))
creator.create("Individual", list, fitness=creator.FitnessMin)

# Toolbox setup
toolbox = base.Toolbox()
toolbox.register("green", random.randint, 10, 60)
toolbox.register("yellow", random.randint, 3, 8)
toolbox.register("red", random.randint, 10, 60)
toolbox.register("speed", lambda: random.choice([40, 60, 80, 100]))
toolbox.register("seed", lambda: 1408)
toolbox.register(
    "individual",
    tools.initCycle,
    creator.Individual,
    (toolbox.green, toolbox.yellow, toolbox.red, toolbox.speed, toolbox.seed),
    n=1,
)
toolbox.register(
    "population",
    tools.initRepeat,
    list,
    toolbox.individual,
)
toolbox.register("mate", tools.cxTwoPoint)
toolbox.register("select", tools.selTournament, tournsize=3)


def custom_mutate(individual, indpb=0.2, min_speed=40):
    if random.random() < indpb:
        individual[0] = random.randint(10, 60)  # Green
    if random.random() < indpb:
        individual[1] = random.randint(3, 8)  # Yellow
    if random.random() < indpb:
        individual[2] = random.randint(10, 60)  # Red
    if random.random() < indpb:
        allowed_speeds = [s for s in [40, 60, 80, 100] if s >= min_speed]
        individual[3] = random.choice(allowed_speeds)
    return (individual,)


toolbox.register("mutate", custom_mutate)


class OptimizationEngine:
    def __init__(self, simulation_server_address=None):
        """
        Initialize the optimization engine with gRPC simulation client

        Args:
            simulation_server_address: Address of the simulation service (default: localhost:50053)
        """
        self.simulation_client = SimulationClient(simulation_server_address)
        self.generated_logs = []

    def run_simulation(self, individual):
        """
        Run simulation using gRPC client instead of subprocess

        Args:
            individual: List containing [green, yellow, red, speed, seed]

        Returns:
            dict or None: Simulation results or None if failed
        """
        green, yellow, red, speed, seed = individual

        try:
            result = self.simulation_client.get_simulation_results(
                green=green,
                yellow=yellow,
                red=red,
                speed=speed,
                seed=seed
            )
            return result
        except Exception as e:
            print(f"Simulation failed for individual {individual}: {e}")
            return None

    def evaluate_waiting_and_travel(self, individual):
        """
        Evaluate individual based on waiting and travel time
        """
        result = self.run_simulation(individual)
        if result is None:
            return (1e6,)

        waiting = result.get("Total Waiting Time", 1e6)
        travel = result.get("Total Travel Time", 1e6)
        return (0.9 * waiting + 0.3 * travel,)  # Weighted objective

    def evaluate_safety_given_waiting(self, individual):
        """
        Evaluate individual based on safety metrics with speed constraints
        """
        if individual[3] < 60:
            return (1e6,)  # Penalize unsafe speeds below 60

        result = self.run_simulation(individual)
        if result is None:
            return (1e6,)

        brakes = result.get("Emergency Brakes", 0)
        stops = result.get("Emergency Stops", 0)
        collisions = result.get("Near collisions", 0)
        waiting = result.get("Total Waiting Time", 0)

        fitness = 1000 * brakes + 1000 * stops + 20000 * collisions + 0.9 * waiting
        return (fitness,)

    def log_individual_to_file(self, individual, generation, ind_id):
        """
        Log individual results to CSV file
        """
        os.makedirs("ga_results", exist_ok=True)
        with open(ALL_RESULTS_CSV, "a") as f:
            f.write(
                f"{generation},{ind_id},{individual[0]},{individual[1]},{individual[2]},"
                f"{individual[3]},{individual[4]},{individual.fitness.values[0]}\n"
            )

    def run_ga(self, pop, hof, ngen, cxpb, mutpb, label="GA"):
        """
        Run genetic algorithm with progress tracking
        """
        stats = tools.Statistics(lambda ind: ind.fitness.values[0])
        stats.register("avg", lambda fits: sum(fits) / len(fits))
        stats.register("min", min)

        logbook = tools.Logbook()
        logbook.header = ["gen", "nevals"] + stats.fields

        for gen in range(ngen + 1):
            if gen == 0:
                invalid_ind = pop
            else:
                offspring = toolbox.select(pop, len(pop))
                offspring = list(map(toolbox.clone, offspring))

                for child1, child2 in zip(offspring[::2], offspring[1::2]):
                    if random.random() < cxpb:
                        toolbox.mate(child1, child2)
                        del child1.fitness.values
                        del child2.fitness.values

                for mutant in offspring:
                    if random.random() < mutpb:
                        toolbox.mutate(mutant)
                        del mutant.fitness.values

                invalid_ind = [ind for ind in offspring if not ind.fitness.valid]
                pop[:] = offspring

            with tqdm(total=len(invalid_ind), desc=f"{label} Gen {gen}") as pbar:
                for i, ind in enumerate(invalid_ind):
                    ind.fitness.values = toolbox.evaluate(ind)
                    self.log_individual_to_file(ind, generation=gen, ind_id=i)
                    pbar.update(1)

            hof.update(pop)
            record = stats.compile(pop)
            logbook.record(gen=gen, nevals=len(invalid_ind), **record)

        return logbook

    def get_final_comparison(self, best_params, reference_results=None):
        """
        Get final simulation results for comparison

        Args:
            best_params: Dictionary with optimized parameters
            reference_results: Optional reference results for comparison

        Returns:
            dict: Final simulation results and comparison
        """
        try:
            final_results = self.simulation_client.get_simulation_results(
                green=best_params["Green"],
                yellow=best_params["Yellow"],
                red=best_params["Red"],
                speed=best_params["Speed"],
                seed=best_params["Seed"]
            )

            comparison = {
                "optimized_results": final_results,
                "reference_results": reference_results,
                "comparison": {}
            }

            if reference_results:
                for metric in [
                    "Total Waiting Time",
                    "Total Travel Time",
                    "Emergency Brakes",
                    "Emergency Stops",
                    "Near collisions",
                ]:
                    opt_val = final_results.get(metric, "N/A")
                    ref_val = reference_results.get(metric, "N/A")
                    comparison["comparison"][metric] = {
                        "optimized": opt_val,
                        "reference": ref_val,
                        "improvement": ref_val - opt_val if isinstance(opt_val, (int, float)) and isinstance(ref_val, (int, float)) else "N/A"
                    }

            return comparison

        except Exception as e:
            print(f"Failed to get final comparison: {e}")
            return {"error": str(e)}

    def run_optimization(self, ngen_waiting=30, ngen_safety=10, pop_size=30, cxpb=0.5, mutpb=0.3, random_seed=1408):
        """
        Run the complete optimization process

        Args:
            ngen_waiting: Generations for waiting time optimization
            ngen_safety: Generations for safety optimization
            pop_size: Population size
            cxpb: Crossover probability
            mutpb: Mutation probability
            random_seed: Random seed for reproducibility

        Returns:
            dict: Complete optimization results
        """
        # Set random seed
        random.seed(random_seed)

        optimization_results = {
            "parameters": {
                "ngen_waiting": ngen_waiting,
                "ngen_safety": ngen_safety,
                "pop_size": pop_size,
                "cxpb": cxpb,
                "mutpb": mutpb,
                "random_seed": random_seed
            },
            "phases": {},
            "best_parameters": {},
            "final_comparison": {}
        }

        try:
            # Phase 1: Minimize waiting time
            print("\n--- Phase 1: Minimizing Waiting Time ---")
            toolbox.register("evaluate", self.evaluate_waiting_and_travel)
            toolbox.register("mutate", lambda ind: custom_mutate(ind, indpb=0.2, min_speed=60))

            pop = toolbox.population(n=pop_size)
            hof_wait = tools.HallOfFame(3)

            # Initialize CSV log
            os.makedirs("ga_results", exist_ok=True)
            with open(ALL_RESULTS_CSV, "w") as f:
                f.write("generation,individual_id,green,yellow,red,speed,seed,fitness\n")

            logbook_waiting = self.run_ga(pop, hof_wait, ngen_waiting, cxpb, mutpb, label="WaitingTime")
            optimization_results["phases"]["waiting_time"] = {
                "hall_of_fame": [
                    {
                        "parameters": {"Green": ind[0], "Yellow": ind[1], "Red": ind[2], "Speed": ind[3], "Seed": ind[4]},
                        "fitness": ind.fitness.values[0]
                    } for ind in hof_wait
                ],
                "logbook": str(logbook_waiting)
            }

            # Phase 2: Minimize safety issues
            print("\n--- Phase 2: Minimizing Safety Hazards ---")
            toolbox.register("evaluate", self.evaluate_safety_given_waiting)
            toolbox.register("mutate", lambda ind: custom_mutate(ind, indpb=0.2, min_speed=60))

            pop2 = [toolbox.clone(ind) for ind in hof_wait]
            pop2 += toolbox.population(n=pop_size - len(pop2))
            hof_safety = tools.HallOfFame(1)

            logbook_safety = self.run_ga(pop2, hof_safety, ngen_safety, cxpb, mutpb, label="Safety")
            optimization_results["phases"]["safety"] = {
                "hall_of_fame": [
                    {
                        "parameters": {"Green": ind[0], "Yellow": ind[1], "Red": ind[2], "Speed": ind[3], "Seed": ind[4]},
                        "fitness": ind.fitness.values[0]
                    } for ind in hof_safety
                ],
                "logbook": str(logbook_safety)
            }

            # Extract best parameters
            best = hof_safety[0]
            best_params = {
                "Green": best[0],
                "Yellow": best[1],
                "Red": best[2],
                "Speed": best[3],
                "Seed": best[4],
                "Fitness": best.fitness.values[0],
            }

            optimization_results["best_parameters"] = best_params

            # Optionally save to file for debugging/logging
            os.makedirs("out", exist_ok=True)
            with open(BEST_PARAM_OUTPUT, "w") as f:
                json.dump(best_params, f, indent=2)

            # Get final comparison
            final_comparison = self.get_final_comparison(best_params)
            optimization_results["final_comparison"] = final_comparison

            print("\nOptimization completed successfully!")
            print("Best Parameters Found:")
            print(json.dumps(best_params, indent=2))

            return optimization_results

        except Exception as e:
            print(f"Optimization failed: {e}")
            optimization_results["error"] = str(e)
            return optimization_results

        finally:
            # Clean up gRPC connection
            self.simulation_client.close()


def main():
    """
    Main function for standalone execution
    """
    engine = OptimizationEngine()
    results = engine.run_optimization()
    return results


if __name__ == "__main__":
    main()
