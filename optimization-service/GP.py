import os
import random
import subprocess
import json
import pickle
from datetime import datetime
from deap import base, creator, tools, algorithms

"""Paths"""
SIM_SCRIPT = "../simulation-service/SimLoad.py"
PARAMS_FOLDER = "../simulation-service/parameters"
RESULTS_FOLDER = "../simulation-service/out/results"
RESULT_FILE_TEMPLATE = os.path.join(RESULTS_FOLDER, "simulation_results_{}.json")
RESULT_FILE = "out/total_result.json"
BEST_PARAM_OUTPUT = "out/best_parameters.json"
REFERENCE_RESULT = "../simulation-service/out/results/simulation_results.json"

"""Track all generated parameter files for cleanup"""
generated_param_files = []
generated_result_files = []

ALL_RESULTS_CSV = "ga_results/all_individuas_log.csv"

"""Fitness Function"""
def evaluate(individual):
    green, yellow, red, speed, seed = individual
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S-%f")
    param_file = os.path.join(PARAMS_FOLDER, f"params_Trafficlight_{timestamp}.json")
    result_file = RESULT_FILE_TEMPLATE.format(timestamp)
    generated_result_files.append(result_file)
    generated_param_files.append(param_file)

    """Build simulation input"""
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
            },
            "output_path": result_file
        }
    }

    """Save parameter file"""
    os.makedirs(PARAMS_FOLDER, exist_ok=True)
    with open(param_file, "w") as f:
        json.dump(params, f)

    """Run simulation"""
    try:
        subprocess.run(
            ["python3", SIM_SCRIPT],
            input=param_file.encode(),
            check=True
        )
    except subprocess.CalledProcessError as e:
        print(f"[Penalty] Simulation subprocess failed for {param_file}: {e}")
        return 1e6,

    """Read simulation results"""
    try:
        with open(result_file, "r") as f:
            result_data = json.load(f)
        metrics = result_data["intersection"]["results"]

        waiting = metrics.get("Total Waiting Time")
        travel = metrics.get("Total Travel Time")
        brakes = metrics.get("Emergency Brakes")
        stops = metrics.get("Emergency Stops")
        collisions = metrics.get("Near collisions")

        """If any metric is missing or None, apply penalty"""
        if None in (waiting, travel, brakes, stops, collisions):
            print(f"[Penalty] Missing metrics in result file: {result_file}")
            return 1e6,

        """Compute weighted fitness (lower is better)"""
        PENALTY_BRAKES = 1000
        PENALTY_STOPS = 1000
        PENALTY_COLLISIONS = 20000

        fitness = (
            waiting * 2 +
            travel +
            PENALTY_BRAKES * brakes +
            PENALTY_STOPS * stops +
            PENALTY_COLLISIONS * collisions
        )

        return fitness,

    except Exception as e:
        print(f"[Penalty] Exception during evaluation: {e}")
        return 1e6,

def log_individual_to_file(individual, generation, ind_id):
    with open(ALL_RESULTS_CSV, "a") as f:
        f.write(f"{generation},{ind_id},{individual[0]},{individual[1]},{individual[2]},"
                f"{individual[3]},{individual[4]},{individual.fitness.values[0]}\n")

"""DEAP Setup"""
creator.create("FitnessMin", base.Fitness, weights=(-1.0,))
creator.create("Individual", list, fitness=creator.FitnessMin)

toolbox = base.Toolbox()
toolbox.register("green", random.randint, 10, 60)
toolbox.register("yellow", random.randint, 3, 8)
toolbox.register("red", random.randint, 10, 60)
toolbox.register("speed", lambda: random.choice([40, 60, 80, 100]))
toolbox.register("seed", lambda: 1408)

toolbox.register("individual", tools.initCycle, creator.Individual,
                 (toolbox.green, toolbox.yellow, toolbox.red, toolbox.speed, toolbox.seed), n=1)
toolbox.register("population", tools.initRepeat, list, toolbox.individual)

toolbox.register("evaluate", evaluate)
toolbox.register("mate", tools.cxTwoPoint)
toolbox.register("select", tools.selTournament, tournsize=3)

def custom_mutate(individual, indpb=0.2):
    if random.random() < indpb:
        individual[0] = random.randint(10, 60)  #Green
    if random.random() < indpb:
        individual[1] = random.randint(3, 8)    #Yellow
    if random.random() < indpb:
        individual[2] = random.randint(10, 60)  #Red
    if random.random() < indpb:
        individual[3] = random.choice([40, 60, 80, 100])  #Speed (restricted)
    return individual,

toolbox.register("mutate", custom_mutate)

"""Final Comparison"""
def run_final_simulation_and_compare(best_params):
    """Save to temporary file"""
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S-%f")
    final_param_file = os.path.join(PARAMS_FOLDER, f"final_params_{timestamp}.json")
    final_result_file = os.path.join(RESULTS_FOLDER, f"final_simulation_result_{timestamp}.json")

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
            },
            "output_path": final_result_file
        }
    }

    with open(final_param_file, "w") as f:
        json.dump(params, f)

    try:
        subprocess.run(
            ["python3", SIM_SCRIPT, "--params", final_param_file],
            check=True
        )
    except subprocess.CalledProcessError as e:
        print(f"Final simulation failed: {e}")
        return

    """Load final simulation results"""
    try:
        with open(final_result_file, "r") as f:
            final_results = json.load(f)["intersection"]["results"]
    except Exception as e:
        print(f"Failed to read final results: {e}")
        return

    """Load reference results"""
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

    """Cleanup"""
    try:
        os.remove(final_result_file)
    except:
        pass

"""Run GA"""
from tqdm import tqdm

def main():
    random.seed(42)
    ngen = 5
    pop_size = 10
    cxpb = 0.5
    mutpb = 0.3

    pop = toolbox.population(n=pop_size)
    hof = tools.HallOfFame(1)
    stats = tools.Statistics(lambda ind: ind.fitness.values[0])
    stats.register("avg", lambda fits: sum(fits) / len(fits))
    stats.register("min", min)

    logbook = tools.Logbook()
    logbook.header = ["gen", "nevals"] + stats.fields

    with open(ALL_RESULTS_CSV, "w") as f:
        f.write("generation,individual_id,green,yellow,red,speed,seed,fitness\n")

    """Evaluate the initial population"""
    print("Evaluating initial population...")
    with tqdm(total=len(pop), desc="Gen 0") as pbar:
        for i, ind in enumerate(pop):
            ind.fitness.values = toolbox.evaluate(ind)
            log_individual_to_file(ind, generation=0, ind_id=i)
            pbar.update(1)

    record = stats.compile(pop)
    logbook.record(gen=0, nevals=len(pop), **record)
    hof.update(pop)

    """Run each generation"""
    for gen in range(1, ngen + 1):
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

        """Evaluate individuals"""
        invalid_ind = [ind for ind in offspring if not ind.fitness.valid]
        print(f"Evaluating Gen {gen}...")
        with tqdm(total=len(invalid_ind), desc=f"Gen {gen}") as pbar:
            for i, ind in enumerate(offspring):
                ind.fitness.values = toolbox.evaluate(ind)
                log_individual_to_file(ind, generation=gen, ind_id=i)
                pbar.update(1)

        pop[:] = offspring
        hof.update(pop)

        record = stats.compile(pop)
        logbook.record(gen=gen, nevals=len(invalid_ind), **record)

    """Save results"""
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

    for fpath in generated_result_files:
        try:
            os.remove(fpath)
        except FileNotFoundError:
            pass
        except Exception as e:
            print(f"Warning: Could not delete result file {fpath}: {e}")

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

    run_final_simulation_and_compare(best_params)

if __name__ == "__main__":
    main()
