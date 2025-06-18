from intersections import circle, stopStreet, tJunction, trafficLight
import json
import time
import uuid


print("Welcome to Sumo!")

def showMenu():
    print("Select an instersection type:")
    print("1. Traffic circle")
    print("2. Stop street")
    print("3. T-Junction")
    print("4. Traffic Light")
    choice = input("Enter choice (1-4): ").strip()
    return choice


def getDefaultTimingsBySpeed(speed):
    if speed <= 40:
        return {"Green":25, "Yellow":3, "Red":30}
    elif speed <= 60:
        return {"Green":25, "Yellow":4, "Red":30}
    elif speed <= 80:
        return {"Green":30, "Yellow":5, "Red":35}
    else:
        print("Speed exceeds reccomended safety for traffic lights, using default for 80km/h")
        return {"Green":30, "Yellow":5, "Red":35}


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
    fake_oid = uuid.uuid4().hex[:24]

    simulationData = {
        "_id": {
            "$oid": fake_oid
        },
        "simulation": {
            "id": "simId",                      # Placeholder
            "name": simName,
            "owner": "username",                # Placeholder
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
    tl = False
    choice = showMenu()

    if choice == '1':
        params = getParams(tl)
        saveParams(params, "roundabout", "Roundabout")
        circle.generate(params)
    elif choice == '2':
        params = getParams(tl)
        saveParams(params, "fourwaystop", "Fourwaystop")
        stopStreet.generate(params)
    elif choice == '3':
        params = getParams(tl)
        saveParams(params, "tjunction", "Tjunction")
        tJunction.generate(params)
    elif choice == '4':
        tl = True
        params = getParams(tl)
        saveParams(params, "trafficlight", "Trafficlight")
        trafficLight.generate(params)
    else:
        print("Invalid choice.")
        main()

if __name__ == "__main__":
    main()
