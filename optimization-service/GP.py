import os
import random
import subprocess
import json
import pickle
from datetime import datetime
from deap import base, creator, tools
from tqdm import tqdm

# Paths
SIM_SCRIPT = "../simulation-service/SimLoad.py"
PARAMS_FOLDER = "../simulation-service/parameters"
RESULTS_FOLDER = "../simulation-service/out/results"
RESULT_FILE_TEMPLATE = os.path.join(RESULTS_FOLDER, "simulation_results_{}.json")
REFERENCE_RESULT = "../simulation-service/out/results/simulation_results.json"
BEST_PARAM_OUTPUT = "out/best_parameters.json"
ALL_RESULTS_CSV = "ga_results/all_individuas_log.csv"

# Track files for cleanup
generated_param_files = []
generated_result_files = []

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
toolbox.register("individual", tools.initCycle, creator.Individual,
                 (toolbox.green, toolbox.yellow, toolbox.red, toolbox.speed, toolbox.seed), n=1)
toolbox.register("population", tools.initRepeat, list, toolbox.individual)
toolbox.register("mate", tools.cxTwoPoint)
toolbox.register("select", tools.selTournament, tournsize=3)

def custom_mutate(individual, indpb=0.2, min_speed=40):
    if random.random() < indpb:
        individual[0] = random.randint(10, 60)  # Green
    if random.random() < indpb:
        individual[1] = random.randint(3, 8)    # Yellow
    if random.random() < indpb:
        individual[2] = random.randint(10, 60)  # Red
    if random.random() < indpb:
        allowed_speeds = [s for s in [40, 60, 80, 100] if s >= min_speed]
        individual[3] = random.choice(allowed_speeds)
    return individual,

toolbox.register("mutate", custom_mutate)

def run_simulation(individual):
    green, yellow, red, speed, seed = individual
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S-%f")
    param_file = os.path.join(PARAMS_FOLDER, f"params_Trafficlight_{timestamp}.json")
    result_file = RESULT_FILE_TEMPLATE.format(timestamp)

    generated_param_files.append(param_file)
    generated_result_files.append(result_file)

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

    os.makedirs(PARAMS_FOLDER, exist_ok=True)
    with open(param_file, "w") as f:
        json.dump(params, f)

    try:
        subprocess.run(["python3", SIM_SCRIPT], input=param_file.encode(), check=True)
    except subprocess.CalledProcessError:
        return None

    try:
        with open(result_file, "r") as f:
            result_data = json.load(f)
        return result_data["intersection"]["results"]
    except:
        return None

def evaluate_waiting_and_travel(individual):
    result = run_simulation(individual)
    if result is None:
        return 1e6,
    waiting = result.get("Total Waiting Time", 1e6)
    travel = result.get("Total Travel Time", 1e6)
    return 0.9 * waiting + 0.3 * travel,  # Weighted objective

def evaluate_safety_given_waiting(individual):
    if individual[3] < 60:
        return 1e6,  # Penalize unsafe speeds below 60

    result = run_simulation(individual)
    if result is None:
        return 1e6,

    brakes = result.get("Emergency Brakes", 0)
    stops = result.get("Emergency Stops", 0)
    collisions = result.get("Near collisions", 0)
    waiting = result.get("Total Waiting Time", 0)

    fitness = (
        1000 * brakes +
        1000 * stops +
        20000 * collisions +
        0.9 * waiting
    )
    return fitness,

def log_individual_to_file(individual, generation, ind_id):
    os.makedirs("ga_results", exist_ok=True)
    with open(ALL_RESULTS_CSV, "a") as f:
        f.write(f"{generation},{ind_id},{individual[0]},{individual[1]},{individual[2]},"
                f"{individual[3]},{individual[4]},{individual.fitness.values[0]}\n")

def run_ga(pop, hof, ngen, cxpb, mutpb, label="GA"):
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
                log_individual_to_file(ind, generation=gen, ind_id=i)
                pbar.update(1)

        hof.update(pop)
        record = stats.compile(pop)
        logbook.record(gen=gen, nevals=len(invalid_ind), **record)

def run_final_simulation_and_compare(best_params):
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
        subprocess.run(["python3", SIM_SCRIPT, "--params", final_param_file], check=True)
    except subprocess.CalledProcessError as e:
        print(f"Final simulation failed: {e}")
        return

    try:
        with open(final_result_file, "r") as f:
            final_results = json.load(f)["intersection"]["results"]
        with open(REFERENCE_RESULT, "r") as f:
            reference_results = json.load(f)["intersection"]["results"]
    except Exception as e:
        print(f"Failed to read final or reference results: {e}")
        return

    print("\n--- Final Comparison ---")
    print(f"{'Metric':<25}{'Optimized':>15}{'Reference':>15}")
    for metric in ["Total Waiting Time", "Total Travel Time", "Emergency Brakes", "Emergency Stops", "Near collisions"]:
        opt = final_results.get(metric, "N/A")
        ref = reference_results.get(metric, "N/A")
        print(f"{metric:<25}{str(opt):>15}{str(ref):>15}")

def cleanup_files():
    for fpath in generated_param_files + generated_result_files:
        try:
            os.remove(fpath)
        except Exception:
            pass

def main():
    random.seed(1408)
    ngen_waiting = 30
    ngen_safety = 10
    pop_size = 30
    cxpb = 0.5
    mutpb = 0.3

    # Phase 1: Minimize waiting time
    print("\n--- Phase 1: Minimizing Waiting Time ---")
    toolbox.register("evaluate", evaluate_waiting_and_travel)
    toolbox.register("mutate", lambda ind: custom_mutate(ind, indpb=0.2, min_speed=60))
    pop = toolbox.population(n=pop_size)
    hof_wait = tools.HallOfFame(3)
    with open(ALL_RESULTS_CSV, "w") as f:
        f.write("generation,individual_id,green,yellow,red,speed,seed,fitness\n")
    run_ga(pop, hof_wait, ngen_waiting, cxpb, mutpb, label="WaitingTime")

    # Phase 2: Minimize safety issues
    print("\n--- Phase 2: Minimizing Safety Hazards ---")
    toolbox.register("evaluate", evaluate_safety_given_waiting)
    toolbox.register("mutate", lambda ind: custom_mutate(ind, indpb=0.2, min_speed=60))
    pop2 = [toolbox.clone(ind) for ind in hof_wait]
    pop2 += toolbox.population(n=pop_size - len(pop2))
    hof_safety = tools.HallOfFame(1)
    run_ga(pop2, hof_safety, ngen_safety, cxpb, mutpb, label="Safety")

    # Save and compare best
    os.makedirs("out", exist_ok=True)
    best = hof_safety[0]
    best_params = {
        "Green": best[0],
        "Yellow": best[1],
        "Red": best[2],
        "Speed": best[3],
        "Seed": best[4],
        "Fitness": best.fitness.values[0]
    }
    with open(BEST_PARAM_OUTPUT, "w") as f:
        json.dump(best_params, f, indent=2)

    print("\nBest Parameters Found:")
    print(json.dumps(best_params, indent=2))

    run_final_simulation_and_compare(best_params)
    cleanup_files()

if __name__ == "__main__":
    main()
