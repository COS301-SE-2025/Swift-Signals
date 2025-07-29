from SimLoad import main

param_dict = {
    "intersection": {
        "simulation_parameters": {"intersection_type": 4, "speed": 60, "seed": 3012},
        "traffic_density": 2,
    }
}

results_dict, full_output = main(param_dict)
