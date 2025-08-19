import os
import csv


def log_individual_to_file(
    individual, generation, ind_id, filename="ga_results/all_individuals.csv"
):
    """
    Logs the individual's parameters and fitness to a CSV file.

    Args:
        individual: The individual to log.
        generation: The current generation number.
        ind_id: The individual's ID in the population.
        filename: The path of the log file.
    """
    os.makedirs(os.path.dirname(filename), exist_ok=True)

    with open(filename, "a", newline="") as f:
        writer = csv.writer(f)
        if f.tell() == 0:
            writer.writerow(
                [
                    "generation",
                    "individual_id",
                    "green",
                    "yellow",
                    "red",
                    "speed",
                    "seed",
                    "fitness",
                ]
            )

        writer.writerow(
            [
                generation,
                ind_id,
                individual[0],  # green
                individual[1],  # yellow
                individual[2],  # red
                individual[3],  # speed
                individual[4],  # seed
                individual.fitness.values[0],  # fitness
            ]
        )
