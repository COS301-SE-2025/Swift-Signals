from intersections import circle, stopStreet, tJunction, trafficLight
from datetime import datetime
import json
import time
import uuid
import os
import warnings


warnings.filterwarnings("ignore", category=DeprecationWarning)


def showMenu():
    print("Select an instersection type:")
    print("1. Traffic circle")
    print("2. Stop street")
    print("3. T-Junction")
    print("4. Traffic Light")
    choice = input("Enter choice (1-4): ").strip()
    return choice


def loadParams():
    filePath = input("Enter path to parameter JSON file: ").strip()

    if not os.path.exists(filePath):
        print("File not found. Exiting.")
        exit(1)

    with open(filePath, "r") as f:
        data = json.load(f)

    return data["simulation"]["parameters"]


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
        print("Speed exceeds reccomended safety for traffic lights, using default for 80km/h")
        return {"Green": 30, "Yellow": 5, "Red": 35}


def getParams(tL: bool):
    trafficDensity = input("Enter traffic density (low/medium/high): ").strip().lower()

    try:
        speed = int(input("Enter road speed limit in km/h (e.g. 40, 60, 80): ").strip())
    except ValueError:
        print("Invalid speed. Falling back to default (40 km/h).")
        speed = 40

    if tL:
        use_default = input("Use default light timings based on road speed? (y/n): ").strip().lower()
        if use_default == 'y':
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
            "Speed": speed
        }
    else:
        return {
            "Traffic Density": trafficDensity,
            "Speed": speed
        }


def saveParams(params, intersectionType, simName):
    timestamp = time.strftime("%Y-%m-%dT%H:%M:%SZ")
    time_for_file = time.strftime("%Y%m%d-%H%M%S")
    fileName = f"params_{simName}_{time_for_file}.json"
    fake_oid = str(uuid.uuid4()).hex[:24]

    simulationData = {
        "_id": {
            "$oid": fake_oid
        },
        "simulation": {
            "id": "simId",
            "name": simName,
            "owner": "username",
            "created_at": timestamp,
            "last_run_at": timestamp,
            "status": "completed",
            "run_count": 0,
            "parameters": {
                "Intersection Type": intersectionType,
                **params,
                "seed": 13
            }
        }
    }

    with open(fileName, "w") as f:
        json.dump(simulationData, f, indent=4)
    print(f"Saved parameters to {fileName}")


def main():
    params = loadParams()
    intersection_type = params.get("Intersection Type", "unknown")
    simName = intersection_type.capitalize()
    simId = "simId"
    owner = "username"
    created = params.get("created_at", "unknown")
    nowIso = datetime.utcnow().isoformat() + "Z"

    runCount = loadRunCount()
    runCount += 1
    saveRunCount(runCount)

    '''Run correct generator based on type'''
    if intersection_type == "trafficlight":
        results = trafficLight.generate(params)
    elif intersection_type == "roundabout":
        results = circle.generate(params)
    elif intersection_type == "fourwaystop":
        results = stopStreet.generate(params)
    elif intersection_type == "tjunction":
        results = tJunction.generate(params)
    else:
        print("Invalid intersection type in parameters.")
        return

    output = {
        "_id": {
            "$oid": str(uuid.uuid4())[:24].replace("-", "0")  
        },
        "simulation": {
            "id": simId,
            "name": simName,
            "owner": owner,
            "created_at": created,
            "last_run_at": nowIso,
            "status": "completed",
            "run_count": runCount,
            "parameters": params,
            "results": results
        }
    }

    '''Save the output to a file'''
    with open("simulation_output.json", "w") as f:
        json.dump(output, f, indent=4)

    print("Simulation saved to simulation_output.json")


if __name__ == "__main__":
    main()
