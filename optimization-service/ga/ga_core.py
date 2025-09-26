import random
from deap import base, creator, tools
from tqdm import tqdm

# Fitness Setup
creator.create("FitnessMin", base.Fitness, weights=(-1.0,))
creator.create("Individual", list, fitness=creator.FitnessMin)

# Toolbox Setup
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
    """A custom mutation operator that only mutates with a given probability."""
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


def run_ga(pop, hof, ngen, cxpb, mutpb, evaluate_func, label="GA"):
    """
    The main genetic algorithm loop.

    Args:
        pop: The initial population.
        hof: Hall of Fame object to store best individuals.
        ngen: The number of generations to run.
        cxpb: The probability of crossover.
        mutpb: The probability of mutation.
        evaluate_func: The function to use for evaluating fitness.
        label: A label for the progress bar.
    """
    toolbox.register("evaluate", evaluate_func)

    stats = tools.Statistics(lambda ind: ind.fitness.values[0])
    stats.register("avg", lambda fits: sum(fits) / len(fits))
    stats.register("min", min)

    logbook = tools.Logbook()
    logbook.header = ["gen", "nevals"] + stats.fields

    for gen in range(ngen + 1):
        if gen == 0:
            invalid_individuals = pop
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

            invalid_individuals = [ind for ind in offspring if not ind.fitness.valid]
            pop[:] = offspring

        with tqdm(total=len(invalid_individuals), desc=f"{label} Gen {gen}") as pbar:
            for i, ind in enumerate(invalid_individuals):
                ind.fitness.values = toolbox.evaluate(ind)
                pbar.update(1)

        hof.update(pop)
        record = stats.compile(pop)
        logbook.record(gen=gen, nevals=len(invalid_individuals), **record)
