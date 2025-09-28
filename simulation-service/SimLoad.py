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
    2: "roundabout",  # possibly removing
    3: "fourwaystop",  # possibly removing
    4: "tjunction",  # possibly removing
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
        parser.add_argument(
            "--params", help="Path to parameter JSON file", required=False
        )
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
    raw_density = data["intersection"].get("traffic density", 1)
    traffic_density = TRAFFIC_DENSITY.get(raw_density, "medium")
    raw_type = sim_params.get("intersection_type", 1)

    try:
        raw_type = int(raw_type)
    except (ValueError, TypeError):
        raw_type = 0

    intersection_type_str = INTERSECTION_TYPES.get(raw_type, "trafficlight")

    mapped = {
        "traffic_density": traffic_density,
        "intersection_type": "trafficlight",  # intersection_type_str,
        "speed": sim_params.get("Speed", 40),
        "seed": sim_params.get("Seed", 42),
    }

    if intersection_type_str == "trafficlight":
        mapped["green"] = sim_params.get("Green", 25)
        mapped["yellow"] = sim_params.get("Yellow", 3)
        mapped["red"] = sim_params.get("Red", 30)

    output_path = data["intersection"].get("output_path", None)
    return {
        "mapped": mapped,
        "raw": {
            "traffic_density": raw_density,
            "intersection_type": raw_type,
            "speed": sim_params.get("Speed", 40),
            "seed": sim_params.get("Seed", 42),
        },
        "output_path": output_path,
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
        return {"green": 25, "yellow": 3, "red": 30}
    elif speed <= 60:
        return {"green": 25, "yellow": 4, "red": 30}
    elif speed <= 80:
        return {"green": 30, "yellow": 5, "red": 35}
    else:
        print(
            "Speed exceeds reccomended safety for traffic lights, using default for 80km/h"
        )
        return {"green": 30, "yellow": 5, "red": 35}


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
                timings = {"green": green, "yellow": yellow, "red": red}
            except ValueError:
                print("Invalid input. Using default for 60 km/h.")
                timings = getDefaultTimingsBySpeed(60)

        return {
            "traffic_density": trafficDensity,
            "green": timings["green"],
            "yellow": timings["yellow"],
            "red": timings["red"],
            "speed": speed,
        }
    else:
        return {"traffic_density": trafficDensity, "speed": speed}


def saveParams(params, intersection_type_str, simName):
    density_map = {"low": 0, "medium": 1, "high": 2}
    density_num = density_map.get(params.get("traffic_density", "medium").lower(), 1)

    reverse_intersection_types = {v: k for k, v in INTERSECTION_TYPES.items()}
    intersection_type_num = reverse_intersection_types.get(intersection_type_str, 0)

    results = {}

    raw = {
        "speed": params.get("speed", 40),
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
            "traffic_density": density_num,
            "status": 0,
            "run_count": 0,
            "parameters": {
                "intersection_type": intersection_type_num,
                "speed": raw["speed"],
                "seed": raw["seed"],
            },
            "results": results,
        },
    }

    with open(fileName, "w") as f:
        json.dump(output, f, indent=4)
    # print(f"Saved parameters to {fileName}")


def main(param_dict=None) -> dict:
    base_dir = os.path.dirname(os.path.abspath(__file__))
    params = loadParams(param_dict)
    mapped = params["mapped"]
    raw = params["raw"]

    intersection_type = mapped.get("intersection_type", "unknown")
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

    parameters = {"intersection_type": raw.get("intersection_type")}

    if intersection_type == "trafficlight":
        parameters["green"] = mapped.get("Green")
        parameters["yellow"] = mapped.get("Yellow")
        parameters["red"] = mapped.get("Red")

    parameters["seed"] = raw.get("seed")

    output = {
        "_id": {"$oid": uuid.uuid4().hex[:24]},
        "intersection": {
            "id": simId,
            "name": simName,
            "owner": owner,
            "created_at": created,
            "last_run_at": nowIso,
            "traffic_density": raw.get("traffic_density"),
            "status": 0,
            "run_count": runCount,
            "parameters": parameters,
            "results": results,
        },
    }

    timestamp = time.strftime("%Y%m%d-%H%M%S")
    intersection_type = mapped.get("intersection_type", "unknown")

    custom_result_path = params.get("output_path", None)

    if custom_result_path:
        results_path = os.path.abspath(custom_result_path)
        output_path = os.path.join(
            base_dir,
            "out/simulationOut",
            f"simulation_output_{intersection_type}_{timestamp}.json",
        )
    else:
        results_filename = f"simulation_results_{intersection_type}_{timestamp}.json"
        output_filename = f"simulation_output_{intersection_type}_{timestamp}.json"
        results_path = os.path.join(base_dir, "out/results", results_filename)
        output_path = os.path.join(base_dir, "out/simulationOut", output_filename)

    os.makedirs(os.path.dirname(results_path), exist_ok=True)
    os.makedirs(os.path.dirname(output_path), exist_ok=True)

    with open(results_path, "w") as f:
        json.dump(output, f, indent=2)
    # print(f"Simulation saved to {results_path}")

    with open(output_path, "w") as jf:
        json.dump(fullOut, jf, indent=2)

    # print(f"Simulation output saved to {output_path}")

    try:
        os.remove("run_count.txt")
        # print("Cleaned up run_count.txt")
    except FileNotFoundError:
        pass
    except Exception as e:
        print(f"Warning: Could not delete run_count.txt - {e}")

    # print(output)
    return output, fullOut


if __name__ == "__main__":
    main()
