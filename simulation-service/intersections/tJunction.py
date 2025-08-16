import os
import subprocess
import xml.etree.ElementTree as ET


def generate(params):
    allowedSpeeds = [40, 60, 80, 100, 120]
    speedKm = params.get("speed", 40)
    if speedKm not in allowedSpeeds:
        print(f"Warning: Speed {speedKm}km/h not allowed. Using default 40km/h.")
        speedKm = 40
    speedInMs = speedKm * (1000 / 3600)

    base = "tjunction"
    netFile = f"{base}.net.xml"
    routeFile = f"{base}.rou.xml"
    configFile = f"{base}.sumocfg"
    tripinfoFile = f"{base}_tripinfo.xml"

    nodeFile = "tjInt.nod.xml"
    edgeFile = "tjInt.edg.xml"
    conFile = "tjInt.con.xml"

    writeNodeFile(nodeFile)
    writeEdgeFile(edgeFile, speedInMs)
    writeConnectionFile(conFile)

    print("Generating T-junction with params:", params)

    subprocess.run(
        [
            "netconvert",
            f"--node-files={nodeFile}",
            f"--edge-files={edgeFile}",
            f"--connection-files={conFile}",
            "-o",
            netFile,
        ],
        check=True,
    )

    generateTrips(netFile, routeFile, params["traffic_density"], params)

    with open(configFile, "w") as cfg:
        cfg.write(
            f"""<configuration>
        <input>
            <net-file value="{netFile}"/>
            <route-files value="{routeFile}"/>
        </input>
        <time>
            <begin value="0"/>
            <end value="3600"/>
        </time>
    </configuration>"""
        )

    logfile = f"{base}_warnings.log"
    fcdOutputFile = f"{base}_fcd.xml"

    with open(logfile, "w") as log:
        subprocess.run(
            [
                "sumo",
                "-c",
                configFile,
                "--tripinfo-output",
                tripinfoFile,
                "--fcd-output",
                fcdOutputFile,
                "--no-warnings",
                "false",
                "--message-log",
                logfile,
            ],
            check=True,
            stdout=log,
            stderr=log,
        )

    print("Simulation finished. Parsing results...")

    emergency_brakes = 0
    emergency_stops = 0
    near_collisions = []

    with open(logfile, "r") as f:
        lines = f.readlines()

    for i in range(len(lines)):
        line = lines[i].strip()
        if "performs emergency braking" in line:
            vehicle_id = line.split("'")[1]
            emergency_brakes += 1
            near_collisions.append(line)
        elif "performs emergency stop" in line:
            vehicle_id = line.split("'")[1]
            print(vehicle_id)
            emergency_stops += 1
            near_collisions.append(line)

    tree = ET.parse(tripinfoFile)
    root = tree.getroot()

    total_vehicles = 0
    total_travel_time = 0.0
    total_waiting_time = 0.0
    total_distance = 0.0
    speeds = []

    for trip in root.findall("tripinfo"):
        total_vehicles += 1
        travel_time = float(trip.get("duration"))
        waiting_time = float(trip.get("waitingTime"))
        distance = float(trip.get("routeLength"))
        total_travel_time += travel_time
        total_waiting_time += waiting_time
        total_distance += distance
        if travel_time > 0:
            speeds.append(distance / travel_time)

    avg_speed = sum(speeds) / len(speeds) if speeds else 0
    avg_waiting_time = total_waiting_time / total_vehicles if total_vehicles > 0 else 0
    avg_travel_time = total_travel_time / total_vehicles if total_vehicles > 0 else 0

    results = {
        "total_vehicles": total_vehicles,
        "average_travel_time": avg_travel_time,
        "total_travel_time": total_travel_time,
        "average_speed": avg_speed,
        "average_waiting_time": avg_waiting_time,
        "total_waiting_time": total_waiting_time,
        "generated_vehicles": total_vehicles,
        "emergency_brakes": emergency_brakes,
        "emergency_stops": emergency_stops,
        "near_collisions": len(near_collisions),
    }

    trajectories = extractTrajectories(fcdOutputFile)

    fullOutput = {
        "intersection": {
            "nodes": parseNodes(nodeFile),
            "edges": parseEdges(edgeFile),
            "connections": parseConnections(conFile),
        },
        "vehicles": trajectories,
    }

    tempFiles = [
        netFile,
        routeFile,
        configFile,
        tripinfoFile,
        nodeFile,
        edgeFile,
        conFile,
        fcdOutputFile,
        logfile,
    ]

    for file in tempFiles:
        try:
            os.remove(file)
        except OSError as e:
            print(f"Warning: Could not delete {file} - {e}")

    return results, fullOutput


def writeNodeFile(filename):
    content = """<nodes>
    <node id="center" x="0" y="0" type="priority"/>
    <node id="n1" x="-100" y="0" type="priority"/>
    <node id="n2" x="100" y="0" type="priority"/>
    <node id="n3" x="0" y="100" type="priority"/>
</nodes>"""
    with open(filename, "w") as f:
        f.write(content)


def writeEdgeFile(filename, speed=11.11):
    content = f"""<edges>
    <edge id="in_n1" from="n1" to="center" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="in_n2" from="n2" to="center" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="in_n3" from="n3" to="center" priority="1" numLanes="1" speed="{speed}"/>

    <edge id="out_center_n1" from="center" to="n1" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="out_center_n2" from="center" to="n2" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="out_center_n3" from="center" to="n3" priority="1" numLanes="1" speed="{speed}"/>
</edges>"""
    with open(filename, "w") as f:
        f.write(content)


def writeConnectionFile(filename):
    content = """<connections>
    <connection from="in_n1" to="out_center_n2" fromLane="0" toLane="0"/>
    <connection from="in_n1" to="out_center_n3" fromLane="0" toLane="0"/>

    <connection from="in_n2" to="out_center_n1" fromLane="0" toLane="0"/>
    <connection from="in_n2" to="out_center_n3" fromLane="0" toLane="0"/>

    <connection from="in_n3" to="out_center_n1" fromLane="0" toLane="0"/>
    <connection from="in_n3" to="out_center_n2" fromLane="0" toLane="0"/>
</connections>"""
    with open(filename, "w") as f:
        f.write(content)


def generateTrips(netFile, tripFile, density, params):
    SUMO_HOME = os.environ.get("SUMO_HOME")
    TOOLS_PATH = os.path.join(SUMO_HOME, "tools")

    period = {"low": "12", "medium": "6", "high": "3"}.get(density, "6")

    tripDir = os.path.dirname(tripFile)
    if tripDir:
        os.makedirs(tripDir, exist_ok=True)

    cmd = [
        "python3",
        os.path.join(TOOLS_PATH, "randomTrips.py"),
        "-n",
        netFile,
        "-o",
        tripFile,
        "--prefix",
        "veh",
        "--seed",
        str(params["seed"]),
        "--min-distance",
        "20",
        "--trip-attributes",
        'departLane="best" departSpeed="max"',
        "--period",
        period,
    ]

    with open(os.devnull, "w") as devnull:
        subprocess.run(cmd, check=True, stderr=devnull)
    print("Trips generated.")


def extractTrajectories(fcdOutputFile):
    tree = ET.parse(fcdOutputFile)
    root = tree.getroot()
    trajectories = {}

    for timestep in root.findall("timestep"):
        time = float(timestep.get("time"))
        for vehicle in timestep.findall("vehicle"):
            vid = vehicle.get("id")
            x = float(vehicle.get("x"))
            y = float(vehicle.get("y"))
            speed = float(vehicle.get("speed"))
            if vid not in trajectories:
                trajectories[vid] = {"id": vid, "positions": []}
            trajectories[vid]["positions"].append(
                {"time": time, "x": x, "y": y, "speed": speed}
            )

    return list(trajectories.values())


def parseNodes(filename):
    tree = ET.parse(filename)
    return [
        {
            "id": n.get("id"),
            "x": float(n.get("x")),
            "y": float(n.get("y")),
            "type": n.get("type").upper(),
        }
        for n in tree.getroot()
    ]


def parseEdges(filename):
    tree = ET.parse(filename)
    return [
        {
            "id": e.get("id"),
            "from": e.get("from"),
            "to": e.get("to"),
            "speed": float(e.get("speed")),
            "lanes": int(e.get("numLanes")),
        }
        for e in tree.getroot()
    ]


def parseConnections(filename):
    tree = ET.parse(filename)
    return [
        {
            "from": c.get("from"),
            "to": c.get("to"),
            "fromLane": int(c.get("fromLane")),
            "toLane": int(c.get("toLane")),
        }
        for c in tree.getroot()
    ]
