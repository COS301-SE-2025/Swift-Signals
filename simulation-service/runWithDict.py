from SimLoad import main

param_dict = {
    "intersection": {
        "simulation_parameters": {
            "Intersection Type": 1,
            "Speed": 60,
            "seed": 3012,
            "Green": 25,
            "Yellow": 4,
            "Red": 35
        },
        "Traffic Density": 2
    }
}

main(param_dict)