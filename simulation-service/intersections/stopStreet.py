import os
import subprocess
import xml.etree.ElementTree as ET


def generate(params):
    allowedSpeeds = [40, 60, 80, 100, 120]
    speedKm = params.get("Speed", 40)
    if speedKm not in allowedSpeeds:
        print(f"Warning: Speed {speedKm}km/h not allowed. Using default 40km/h.")
        speedKm = 40
    speedInMs = speedKm * (1000 / 3600)

    base = "stop_intersection"
    netFile = f"{base}.net.xml"
    routeFile = f"{base}.rou.xml"
    configFile = f"{base}.sumocfg"
    tripinfoFile = f"{base}_tripinfo.xml"

    nodeFile = "stopInt.nod.xml"
    edgeFile = "stopInt.edg.xml"
    conFile = "stopInt.con.xml"

    writeNodeFile(nodeFile)
    writeEdgeFile(edgeFile, speedInMs)
    writeConnectionFile(conFile)

    print("Generating stop-controlled intersection with params:", params)

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

    generateTrips(netFile, routeFile, params["Traffic Density"], params)

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
        "Total Vehicles": total_vehicles,
        "Average Travel Time": avg_travel_time,
        "Total Travel Time": total_travel_time,
        "Average Speed": avg_speed,
        "Average Waiting Time": avg_waiting_time,
        "Total Waiting Time": total_waiting_time,
        "Generated Vehicles": total_vehicles,
        "Emergency Brakes": emergency_brakes,
        "Emergency Stops": emergency_stops,
        "Near collisions": len(near_collisions),
    }

    trajectories = extractTrajectories(fcdOutputFile)

    fullOutput = {
        "intersection": {
            "nodes": parseNodes(nodeFile),
            "edges": parseEdges(edgeFile),
            "connections": parseConnections(conFile),
            "trafficLights": [],
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
        f"{base}.add.xml",
    ]
    for file in tempFiles:
        try:
            os.remove(file)
        except OSError as e:
            print(f"Warning: Could not delete {file} - {e}")

    return results, fullOutput


def writeNodeFile(filename):
    content = """<nodes>
    <node id="1" x="0" y="0" type="priority"/>
    <node id="n2" x="0" y="100" type="priority"/>
    <node id="n3" x="0" y="-100" type="priority"/>
    <node id="n4" x="-100" y="0" type="priority"/>
    <node id="n5" x="100" y="0" type="priority"/>
</nodes>"""
    with open(filename, "w") as f:
        f.write(content)


def writeEdgeFile(filename, speed=11.11):
    content = f"""<edges>
    <edge id="in_n2_1" from="n2" to="1" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="in_n3_1" from="n3" to="1" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="in_n4_1" from="n4" to="1" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="in_n5_1" from="n5" to="1" priority="1" numLanes="1" speed="{speed}"/>

    <edge id="out_1_n2" from="1" to="n2" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="out_1_n3" from="1" to="n3" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="out_1_n4" from="1" to="n4" priority="1" numLanes="1" speed="{speed}"/>
    <edge id="out_1_n5" from="1" to="n5" priority="1" numLanes="1" speed="{speed}"/>
</edges>"""
    with open(filename, "w") as f:
        f.write(content)


def writeConnectionFile(filename):
    content = """<connections>
    <connection from="in_n2_1" to="out_1_n3" fromLane="0" toLane="0"/>
    <connection from="in_n2_1" to="out_1_n5" fromLane="0" toLane="0"/>
    <connection from="in_n2_1" to="out_1_n4" fromLane="0" toLane="0"/>

    <connection from="in_n3_1" to="out_1_n4" fromLane="0" toLane="0"/>
    <connection from="in_n3_1" to="out_1_n5" fromLane="0" toLane="0"/>
    <connection from="in_n3_1" to="out_1_n2" fromLane="0" toLane="0"/>

    <connection from="in_n4_1" to="out_1_n5" fromLane="0" toLane="0"/>
    <connection from="in_n4_1" to="out_1_n2" fromLane="0" toLane="0"/>
    <connection from="in_n4_1" to="out_1_n3" fromLane="0" toLane="0"/>

    <connection from="in_n5_1" to="out_1_n2" fromLane="0" toLane="0"/>
    <connection from="in_n5_1" to="out_1_n3" fromLane="0" toLane="0"/>
    <connection from="in_n5_1" to="out_1_n4" fromLane="0" toLane="0"/>
</connections>"""
    with open(filename, "w") as f:
        f.write(content)


def parseNodes(filename):
    tree = ET.parse(filename)
    root = tree.getroot()
    return [
        {
            "id": n.get("id"),
            "x": float(n.get("x")),
            "y": float(n.get("y")),
            "type": n.get("type"),
        }
        for n in root.findall("node")
    ]


def parseEdges(filename):
    tree = ET.parse(filename)
    root = tree.getroot()
    return [
        {
            "id": e.get("id"),
            "from": e.get("from"),
            "to": e.get("to"),
            "speed": float(e.get("speed")),
            "lanes": int(e.get("numLanes")),
        }
        for e in root.findall("edge")
    ]


def parseConnections(filename):
    tree = ET.parse(filename)
    root = tree.getroot()
    return [
        {
            "from": c.get("from"),
            "to": c.get("to"),
            "fromLane": int(c.get("fromLane")),
            "toLane": int(c.get("toLane")),
        }
        for c in root.findall("connection")
    ]


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


def generateTrips(netFile, tripFile, density, params):
    SUMO_HOME = os.environ.get("SUMO_HOME")
    TOOLS_PATH = os.path.join(SUMO_HOME, "tools")

    if density == "low":
        period = "12"
    elif density == "medium":
        period = "6"
    elif density == "high":
        period = "3"
    else:
        period = "6"

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
