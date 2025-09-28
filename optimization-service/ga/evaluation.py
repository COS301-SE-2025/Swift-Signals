from ga.simulation_client import run_simulation


def evaluate_waiting_and_travel(individual: list, traffic_density: int = 2) -> tuple:
    """Evaluates an individual based on a weighted sum of waiting and travel times."""
    results = run_simulation(individual)
    if results is None:
        return (1e6,)  # Returns a high penalty if the simulation fails
    waiting = results.get("total_waiting_time", 1e6)
    travel = results.get("total_travel_time", 1e6)

    fitness = 0.9 * waiting + 0.3 * travel  # Weighted objective
    return (fitness,)


def evaluate_safety_given_waiting(individual: list, traffic_density: int = 2) -> tuple:
    """
    Evaluates an individual based on safety metrics and waiting time,
    with a penalty for low speeds.
    """
    if individual[3] < 60:
        return (1e6,)  # Penalize unsafe speeds below 60

    results = run_simulation(individual)
    if results is None:
        return (1e6,)  # Returns a high penalty if the simulation fails

    brakes = results.get("emergency_brakes", 0)
    stops = results.get("emergency_stops", 0)
    collisions = results.get("near_collisions", 0)
    waiting = results.get("total_waiting_time", 0)

    fitness = 1000 * brakes + 1000 * stops + 20000 * collisions + 0.9 * waiting
    return (fitness,)


def evaluate_balanced(individual: list, traffic_density: int = 2) -> tuple:
    if individual[3] < 60:  # Penalize very unsafe speeds
        return (1e6,)

    result = run_simulation(individual, traffic_density)
    if result is None:
        return (1e6,)  # Returns a high penalty if the simulation fails

    # Efficiency metrics
    waiting = result.get("Total Waiting Time", 1e6)
    travel = result.get("Total Travel Time", 1e6)

    # Safety metrics
    brakes = result.get("Emergency Brakes", 0)
    stops = result.get("Emergency Stops", 0)
    collisions = result.get("Near collisions", 0)

    # Weighted combination
    efficiency_score = 0.7 * (0.9 * waiting + 0.3 * travel)
    safety_score = 0.3 * (1000 * brakes + 1000 * stops + 20000 * collisions)

    return (efficiency_score + safety_score,)
