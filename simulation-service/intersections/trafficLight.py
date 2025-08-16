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

    base = "tl_intersection"
    netFile = f"{base}.net.xml"
    routeFile = f"{base}.rou.xml"
    configFile = f"{base}.sumocfg"
    tllFile = f"{base}.tll.xml"
    tripinfoFile = f"{base}_tripinfo.xml"

    nodeFile = "tlInt.nod.xml"
    edgeFile = "tlInt.edg.xml"
    conFile = "tlInt.con.xml"

    writeNodeFile(nodeFile)
    writeEdgeFile(edgeFile, speedInMs)
    writeConnectionFile(conFile)

    print("Generating traffic light intersection with params:", params)

    writeTrafficLightLogic(tllFile, params["green"], params["yellow"], params["red"])

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
            <additional-files value="{tllFile}"/>
        </input>
        <time>
            <begin value="0"/>
            <end value="3600"/>
        </time>
    </configuration>"""
        )

    """run GUI by using sumo-gui"""
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

            if (
                i + 1 < len(lines)
                and vehicle_id in lines[i + 1]
                and "because of a red traffic light" in lines[i + 1]
            ):
                continue
            else:
                near_collisions.append(line)

        elif "performs emergency stop" in line:
            vehicle_id = line.split("'")[1]
            emergency_stops += 1

            if (
                i + 1 < len(lines)
                and vehicle_id in lines[i + 1]
                and "because of a red traffic light" in lines[i + 1]
            ):
                continue
            else:
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
            "trafficLights": parseTrafficLights(tllFile),
        },
        "vehicles": trajectories,
    }

    tempFiles = [
        netFile,
        routeFile,
        configFile,
        tllFile,
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


def parseNodes(filename):
    tree = ET.parse(filename)
    root = tree.getroot()
    return [
        {
            "id": n.get("id"),
            "x": float(n.get("x")),
            "y": float(n.get("y")),
            "type": n.get("type", "").upper(),
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
            "tl": c.get("tl"),
        }
        for c in root.findall("connection")
    ]


def parseTrafficLights(filename):
    tree = ET.parse(filename)
    root = tree.getroot()
    result = []
    for logic in root.findall("tlLogic"):
        result.append(
            {
                "id": logic.get("id"),
                "type": logic.get("type"),
                "phases": [
                    {
                        "duration": int(phase.get("duration")),
                        "state": phase.get("state"),
                    }
                    for phase in logic.findall("phase")
                ],
            }
        )
    return result


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


def writeNodeFile(filename):
    content = """<nodes>
    <node id="1" x="0" y="0" type="traffic_light"/>
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
    <connection from="in_n2_1" to="out_1_n3" fromLane="0" toLane="0" tl="0"/>
    <connection from="in_n2_1" to="out_1_n5" fromLane="0" toLane="0" tl="1"/>
    <connection from="in_n2_1" to="out_1_n4" fromLane="0" toLane="0" tl="2"/>

    <connection from="in_n3_1" to="out_1_n4" fromLane="0" toLane="0" tl="3"/>
    <connection from="in_n3_1" to="out_1_n5" fromLane="0" toLane="0" tl="4"/>
    <connection from="in_n3_1" to="out_1_n2" fromLane="0" toLane="0" tl="5"/>

    <connection from="in_n4_1" to="out_1_n5" fromLane="0" toLane="0" tl="6"/>
    <connection from="in_n4_1" to="out_1_n2" fromLane="0" toLane="0" tl="7"/>
    <connection from="in_n4_1" to="out_1_n3" fromLane="0" toLane="0" tl="8"/>

    <connection from="in_n5_1" to="out_1_n2" fromLane="0" toLane="0" tl="9"/>
    <connection from="in_n5_1" to="out_1_n3" fromLane="0" toLane="0" tl="10"/>
    <connection from="in_n5_1" to="out_1_n4" fromLane="0" toLane="0" tl="11"/>
</connections>"""
    with open(filename, "w") as f:
        f.write(content)


def writeTrafficLightLogic(filename, greenDuration, yellowDuration, redDuration):
    """Phase 1: green for some lanes"""
    phase1_state = list("r" * 12)
    for i in [0, 1, 2, 6, 7, 8]:
        phase1_state[i] = "G"
    phase1_state = "".join(phase1_state)

    """Phase 2: yellow for same lanes (transitional)"""
    phase2_state = list("r" * 12)
    for i in [0, 1, 2, 6, 7, 8]:
        phase2_state[i] = "y"
    phase2_state = "".join(phase2_state)

    """Phase 3: green for other lanes"""
    phase3_state = list("r" * 12)
    for i in [3, 4, 5, 9, 10, 11]:
        phase3_state[i] = "G"
    phase3_state = "".join(phase3_state)

    """Phase 4: yellow for other lanes"""
    phase4_state = list("r" * 12)
    for i in [3, 4, 5, 9, 10, 11]:
        phase4_state[i] = "y"
    phase4_state = "".join(phase4_state)

    with open(filename, "w") as tl:
        tl.write(
            f"""<additional>
    <tlLogic id="1" type="static" programID="custom" offset="0">
        <phase duration="{greenDuration}" state="{phase1_state}"/>
        <phase duration="{yellowDuration}" state="{phase2_state}"/>
        <phase duration="{redDuration}" state="{phase3_state}"/>
        <phase duration="{yellowDuration}" state="{phase4_state}"/>
    </tlLogic>
</additional>"""
        )


def generateTrips(netFile, tripFile, density, params):
    import subprocess

    SUMO_HOME = os.environ.get("SUMO_HOME")
    TOOLS_PATH = os.path.join(SUMO_HOME, "tools")

    if density == "low":
        period = "12"
        """300 vehicles p/h"""
    elif density == "medium":
        period = "6"
        """600 vehicles p/h"""
    elif density == "high":
        period = "3"
        """1200 vehicles p/h"""
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
