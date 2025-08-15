import os
import random
import subprocess
import json
import pickle
from datetime import datetime
from deap import base, creator, tools, algorithms

# --- Paths ---
SIM_SCRIPT = "../simulation-service/SimLoad.py"
PARAMS_FOLDER = "../simulation-service/parameters"
RESULT_FILE = "../simulation-service/out/results/simulation_results.json"
BEST_PARAM_OUTPUT = "out/best_parameters.json"
REFERENCE_RESULT = "../simulation-service/out/results/simulation_results.json"

# Track all generated parameter files for cleanup
generated_param_files = []

# --- Fitness Function ---
def evaluate(individual):
    green, yellow, red, speed, seed = individual
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S-%f")
    param_file = os.path.join(PARAMS_FOLDER, f"params_Trafficlight_{timestamp}.json")
    generated_param_files.append(param_file)

    # Build simulation input
    params = {
        "intersection": {
            "name": "Trafficlight",
            "Traffic Density": 2,
            "simulation_parameters": {
                "Intersection Type": 1,
                "Green": green,
                "Yellow": yellow,
                "Red": red,
                "Speed": speed,
                "seed": seed
            }
        }
    }

    # Save parameter file
    os.makedirs(PARAMS_FOLDER, exist_ok=True)
    with open(param_file, "w") as f:
        json.dump(params, f)

    # Run simulation
    try:
        subprocess.run(
            ["python3", SIM_SCRIPT],
            input=param_file.encode(),  # Pass parameter file path via stdin
            check=True
        )
    except subprocess.CalledProcessError as e:
        print(f"Simulation failed: {e}")
        return 1e6,  # Heavy penalty

    # Read simulation results
    try:
        with open(RESULT_FILE, "r") as f:
            result_data = json.load(f)
        metrics = result_data["intersection"]["results"]

        waiting = metrics.get("Total Waiting Time", 1e6)
        travel = metrics.get("Total Travel Time", 1e6)
        brakes = metrics.get("Emergency Brakes", 0)
        stops = metrics.get("Emergency Stops", 0)
        collisions = metrics.get("Near collisions", 0)

        # Compute weighted fitness (lower is better)
        PENALTY_BRAKES = 1000
        PENALTY_STOPS = 1000
        PENALTY_COLLISIONS = 2000

        fitness = waiting + travel + PENALTY_BRAKES * brakes + PENALTY_STOPS * stops + PENALTY_COLLISIONS * collisions
        return fitness,

    except Exception as e:
        print(f"Error reading results: {e}")
        return 1e6,

# --- DEAP Setup ---
creator.create("FitnessMin", base.Fitness, weights=(-1.0,))
creator.create("Individual", list, fitness=creator.FitnessMin)

toolbox = base.Toolbox()
toolbox.register("green", random.randint, 10, 60)
toolbox.register("yellow", random.randint, 3, 8)
toolbox.register("red", random.randint, 10, 60)
toolbox.register("speed", random.randint, 40, 120)
toolbox.register("seed", random.randint, 0, 10000)

toolbox.register("individual", tools.initCycle, creator.Individual,
                 (toolbox.green, toolbox.yellow, toolbox.red, toolbox.speed, toolbox.seed), n=1)
toolbox.register("population", tools.initRepeat, list, toolbox.individual)

toolbox.register("evaluate", evaluate)
toolbox.register("mate", tools.cxTwoPoint)
toolbox.register("mutate", tools.mutUniformInt,
                 low=[10, 3, 10, 40, 0],
                 up=[60, 8, 60, 120, 10000],
                 indpb=0.2)
toolbox.register("select", tools.selTournament, tournsize=3)

# --- Final Comparison ---
def run_final_simulation_and_compare(best_params):
    # Save to temporary file
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S-%f")
    final_param_file = os.path.join(PARAMS_FOLDER, f"final_params_{timestamp}.json")
    params = {
        "intersection": {
            "name": "Trafficlight",
            "Traffic Density": 2,
            "simulation_parameters": {
                "Intersection Type": 1,
                "Green": best_params["Green"],
                "Yellow": best_params["Yellow"],
                "Red": best_params["Red"],
                "Speed": best_params["Speed"],
                "seed": best_params["Seed"]
            }
        }
    }

    with open(final_param_file, "w") as f:
        json.dump(params, f)

    try:
        subprocess.run(
            ["python3", SIM_SCRIPT],
            input=final_param_file.encode(),
            check=True
        )
    except subprocess.CalledProcessError as e:
        print(f"Final simulation failed: {e}")
        return

    # Load final simulation results
    try:
        with open(RESULT_FILE, "r") as f:
            final_results = json.load(f)["intersection"]["results"]
    except Exception as e:
        print(f"Failed to read final results: {e}")
        return

    # Load reference results
    try:
        with open(REFERENCE_RESULT, "r") as f:
            reference_results = json.load(f)["intersection"]["results"]
    except Exception as e:
        print(f"Failed to read reference results: {e}")
        return

    print("\n--- Final Comparison ---")
    print(f"{'Metric':<25}{'Optimized':>15}{'Reference':>15}")
    for metric in ["Total Waiting Time", "Total Travel Time", "Emergency Brakes", "Emergency Stops", "Near collisions"]:
        opt = final_results.get(metric, "N/A")
        ref = reference_results.get(metric, "N/A")
        print(f"{metric:<25}{str(opt):>15}{str(ref):>15}")

    # Cleanup
    try:
        os.remove(final_param_file)
    except:
        pass

# --- Run GA ---
def main():
    random.seed(42)
    pop = toolbox.population(n=10)
    hof = tools.HallOfFame(1)
    stats = tools.Statistics(lambda ind: ind.fitness.values[0])
    stats.register("avg", lambda fits: sum(fits) / len(fits))
    stats.register("min", min)

    pop, logbook = algorithms.eaSimple(pop, toolbox,
                                       cxpb=0.5, mutpb=0.3,
                                       ngen=5, stats=stats,
                                       halloffame=hof, verbose=True)

    os.makedirs("ga_results", exist_ok=True)
    with open("ga_results/best_result.pkl", "wb") as f:
        pickle.dump((pop, hof, logbook), f)

    for fpath in generated_param_files:
        try:
            os.remove(fpath)
        except FileNotFoundError:
            pass
        except Exception as e:
            print(f"Warning: Could not delete {fpath}: {e}")

    os.makedirs("out", exist_ok=True)
    best_params = {
        "Green": hof[0][0],
        "Yellow": hof[0][1],
        "Red": hof[0][2],
        "Speed": hof[0][3],
        "Seed": hof[0][4],
        "Fitness": hof[0].fitness.values[0]
    }
    with open(BEST_PARAM_OUTPUT, "w") as f:
        json.dump(best_params, f, indent=2)

    print("\nBest Parameters Found:")
    print(json.dumps(best_params, indent=2))

    # Compare final simulation against reference
    run_final_simulation_and_compare(best_params)

if __name__ == "__main__":
    main()
