from intersections import circle, stopStreet, tJunction, trafficLight
from datetime import datetime
import json
import time
import uuid
import os
import warnings
import argparse


warnings.filterwarnings("ignore", category=DeprecationWarning)


"""ENUM mappings"""
INTERSECTION_TYPES = {
    0: "unspecified",
    1: "trafficlight",
    2: "roundabout",
    3: "fourwaystop",
    4: "tjunction",
}


TRAFFIC_DENSITY = {0: "low", 1: "medium", 2: "high"}


INTERSECTION_STATUS = {0: "unoptimized", 1: "optimizing", 2: "optimized", 3: "failed"}


def showMenu():
    print("Select an instersection type:")
    print("1. Traffic circle")
    print("2. Stop street")
    print("3. T-Junction")
    print("4. Traffic Light")
    choice = input("Enter choice (1-4): ").strip()
    return choice


def loadParams(param_dict=None):
    if param_dict:
        data = param_dict
    else:
        parser = argparse.ArgumentParser()
        parser.add_argument("--params", help="Path to parameter JSON file", required=False)
        args = parser.parse_args()

        if args.params and os.path.exists(args.params):
            filePath = args.params
        else:
            filePath = input("Enter path to parameter JSON file: ").strip()

        if not os.path.exists(filePath):
            print("File not found. Exiting.")
            exit(1)

        with open(filePath, "r") as f:
            data = json.load(f)

    sim_params = data["intersection"]["simulation_parameters"]
    raw_density = data["intersection"].get("Traffic Density", 1)
    traffic_density = TRAFFIC_DENSITY.get(raw_density, "medium")
    raw_type = sim_params.get("Intersection Type", 0)

    try:
        raw_type = int(raw_type)
    except (ValueError, TypeError):
        raw_type = 0

    intersection_type_str = INTERSECTION_TYPES.get(raw_type, "unspecified")

    mapped = {
        "Traffic Density": traffic_density,
        "Intersection Type": intersection_type_str,
        "Speed": sim_params.get("Speed", 40),
        "seed": sim_params.get("seed", 42),
    }

    if intersection_type_str == "trafficlight":
        mapped["Green"] = sim_params.get("Green", 25)
        mapped["Yellow"] = sim_params.get("Yellow", 3)
        mapped["Red"] = sim_params.get("Red", 30)

    return {
        "mapped": mapped,
        "raw": {
            "Traffic Density": raw_density,
            "Intersection Type": raw_type,
            "Speed": sim_params.get("Speed", 40),
            "seed": sim_params.get("seed", 42),
        },
    }


def loadRunCount(counter_file="run_count.txt"):
    if os.path.exists(counter_file):
        with open(counter_file, "r") as f:
            return int(f.read().strip())
    return 0


def saveRunCount(count, counter_file="run_count.txt"):
    with open(counter_file, "w") as f:
        f.write(str(count))


def getDefaultTimingsBySpeed(speed):
    if speed <= 40:
        return {"Green": 25, "Yellow": 3, "Red": 30}
    elif speed <= 60:
        return {"Green": 25, "Yellow": 4, "Red": 30}
    elif speed <= 80:
        return {"Green": 30, "Yellow": 5, "Red": 35}
    else:
        print(
            "Speed exceeds reccomended safety for traffic lights, using default for 80km/h"
        )
        return {"Green": 30, "Yellow": 5, "Red": 35}


def getParams(tL: bool):
    trafficDensity = input("Enter traffic density (low/medium/high): ").strip().lower()

    try:
        speed = int(input("Enter road speed limit in km/h (e.g. 40, 60, 80): ").strip())
    except ValueError:
        print("Invalid speed. Falling back to default (40 km/h).")
        speed = 40

    if tL:
        use_default = (
            input("Use default light timings based on road speed? (y/n): ")
            .strip()
            .lower()
        )
        if use_default == "y":
            timings = getDefaultTimingsBySpeed(speed)
        else:
            try:
                green = int(input("Enter green light duration in seconds: ").strip())
                yellow = int(input("Enter yellow light duration in seconds: ").strip())
                red = int(input("Enter red light duration in seconds: ").strip())
                timings = {"Green": green, "Yellow": yellow, "Red": red}
            except ValueError:
                print("Invalid input. Using default for 60 km/h.")
                timings = getDefaultTimingsBySpeed(60)

        return {
            "Traffic Density": trafficDensity,
            "Green": timings["Green"],
            "Yellow": timings["Yellow"],
            "Red": timings["Red"],
            "Speed": speed,
        }
    else:
        return {"Traffic Density": trafficDensity, "Speed": speed}


def saveParams(params, intersection_type_str, simName):
    density_map = {"low": 0, "medium": 1, "high": 2}
    density_num = density_map.get(params.get("Traffic Density", "medium").lower(), 1)

    reverse_intersection_types = {v: k for k, v in INTERSECTION_TYPES.items()}
    intersection_type_num = reverse_intersection_types.get(intersection_type_str, 0)

    results = {}

    raw = {
        "Speed": params.get("Speed", 40),
        "seed": params.get("seed", 42),
    }

    timestamp = time.strftime("%Y-%m-%dT%H:%M:%SZ")
    time_for_file = time.strftime("%Y%m%d-%H%M%S")
    fileName = f"params_{simName}_{time_for_file}.json"

    output = {
        "_id": {"$oid": uuid.uuid4().hex[:24]},
        "intersection": {
            "id": "simId",
            "name": simName,
            "owner": "username",
            "created_at": timestamp,
            "last_run_at": datetime.utcnow().isoformat() + "Z",
            "Traffic Density": density_num,
            "status": 0,
            "run_count": 0,
            "parameters": {
                "Intersection Type": intersection_type_num,
                "Speed": raw["Speed"],
                "seed": raw["seed"],
            },
            "results": results,
        },
    }

    with open(fileName, "w") as f:
        json.dump(output, f, indent=4)
    print(f"Saved parameters to {fileName}")


def main(param_dict=None):
    params = loadParams(param_dict)
    mapped = params["mapped"]
    raw = params["raw"]

    intersection_type = mapped.get("Intersection Type", "unknown")
    simName = intersection_type.capitalize()
    simId = "simId"
    owner = "username"
    created = params.get("created_at", "unknown")
    nowIso = datetime.utcnow().isoformat() + "Z"

    runCount = loadRunCount()
    runCount += 1
    saveRunCount(runCount)

    """Run correct generator based on type"""
    if intersection_type == "trafficlight":
        results, fullOut = trafficLight.generate(mapped)
    elif intersection_type == "roundabout":
        results, fullOut = circle.generate(mapped)
    elif intersection_type == "fourwaystop":
        results, fullOut = stopStreet.generate(mapped)
    elif intersection_type == "tjunction":
        results, fullOut = tJunction.generate(mapped)
    else:
        print("Invalid intersection type in parameters.")
        return

    parameters = {"Intersection Type": raw.get("Intersection Type")}

    if intersection_type == "trafficlight":
        parameters["Green"] = mapped.get("Green")
        parameters["Yellow"] = mapped.get("Yellow")
        parameters["Red"] = mapped.get("Red")

    parameters["Seed"] = raw.get("seed")

    output = {
        "_id": {"$oid": uuid.uuid4().hex[:24]},
        "intersection": {
            "id": simId,
            "name": simName,
            "owner": owner,
            "created_at": created,
            "last_run_at": nowIso,
            "Traffic Density": raw.get("Traffic Density"),
            "status": 0,
            "run_count": runCount,
            "parameters": parameters,
            "results": results,
        },
    }

    """Save the output to a file"""
    os.makedirs("out/results", exist_ok=True)
    with open("out/results/simulation_results.json", "w") as f:
        json.dump(output, f, indent=2)

    print("Simulation saved to simulation_results.json")

    with open("out/simulationOut/simulation_output.json", "w") as jf:
        json.dump(fullOut, jf, indent=2)

    print("Simulation output saved to simulation_output.json")

    try:
        os.remove("run_count.txt")
        print("Cleaned up run_count.txt")
    except FileNotFoundError:
        pass
    except Exception as e:
        print(f"Warning: Could not delete run_count.txt - {e}")


if __name__ == "__main__":
    main()
