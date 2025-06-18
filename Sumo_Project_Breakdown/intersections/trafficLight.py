import os
import subprocess
import xml.etree.ElementTree as ET


def generate(params):
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
    writeEdgeFile(edgeFile)
    writeConnectionFile(conFile)

    print("Generating traffic light intersection with params:", params)

    writeTrafficLightLogic(tllFile, params["Green"], params["Red"])

    subprocess.run([
        "netconvert",
        f"--node-files={nodeFile}",
        f"--edge-files={edgeFile}",
        f"--connection-files={conFile}",
        "-o", netFile
    ], check=True)

    generateTrips(netFile, routeFile, params["Traffic Density"])

    with open(configFile, "w") as cfg:
        cfg.write(f"""<configuration>
        <input>
            <net-file value="{netFile}"/>
            <route-files value="{routeFile}"/>
            <additional-files value="{tllFile}"/>
        </input>
        <time>
            <begin value="0"/>
            <end value="3600"/>
        </time>
    </configuration>""")

    print("Running SUMO simulation...")
    #run GUI by using sumo-gui
    logfile = f"{base}_warnings.log"
    with open(logfile, "w") as log:
        subprocess.run([
            "sumo-gui",
            "-c", configFile,
            "--tripinfo-output", tripinfoFile,
            "--no-warnings", "false", #print warnings
            "--message-log", logfile
        ], check=True, stdout=log, stderr=log)

    print("Simulation finished. Parsing results...")

    emergency_brakes = 0
    emergency_stops = 0

    with open(logfile, "r") as f:
        for line in f:
            if "performs emergency braking" in line:
                emergency_brakes += 1
            elif "performs emergency stop" in line:
                emergency_stops += 1

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
    avg_waiting_time = total_waiting_time/total_vehicles if total_vehicles > 0 else 0
    avg_travel_time = total_travel_time/total_vehicles if total_vehicles > 0 else 0

    results = {
        "Total Vehicles": total_vehicles,
        "Average Travel Time": avg_travel_time,
        "Total Travel Time": total_travel_time,
        "Average Speed": avg_speed,
        "Average Waiting Time": avg_waiting_time,
        "Total Waiting Time": total_waiting_time,
        "Generated Vehicles": total_vehicles,
        "Emergency Brakes": emergency_brakes,
        "Emergency Stops": emergency_stops
    }

    print("Results:", results)
    return results

def writeNodeFile(filename):
    #four nodes + center
    content = """<nodes>
    <node id="1" x="0" y="0" type="traffic_light"/>
    <node id="n2" x="0" y="100" type="priority"/>
    <node id="n3" x="0" y="-100" type="priority"/>
    <node id="n4" x="-100" y="0" type="priority"/>
    <node id="n5" x="100" y="0" type="priority"/>
</nodes>"""
    with open(filename, "w") as f:
        f.write(content)


def writeEdgeFile(filename):
    content = """<edges>
    <edge id="in_n2_1" from="n2" to="1" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_n3_1" from="n3" to="1" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_n4_1" from="n4" to="1" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_n5_1" from="n5" to="1" priority="1" numLanes="1" speed="13.9"/>

    <edge id="out_1_n2" from="1" to="n2" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_1_n3" from="1" to="n3" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_1_n4" from="1" to="n4" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_1_n5" from="1" to="n5" priority="1" numLanes="1" speed="13.9"/>
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


def writeTrafficLightLogic(filename, greenDuration, redDuration):
    phase1_state = list("r" * 12)
    for i in [0,1,2,6,7,8]:
        phase1_state[i] = "G"
    phase1_state = "".join(phase1_state)

    phase2_state = list("r" * 12)
    for i in [3,4,5,9,10,11]:
        phase2_state[i] = "G"
    phase2_state = "".join(phase2_state)

    with open(filename, "w") as tl:
        tl.write(f"""<additional>
    <tlLogic id="1" type="static" programID="custom" offset="0">
        <phase duration="{greenDuration}" state="{phase1_state}"/>
        <phase duration="{redDuration}" state="{phase2_state}"/>
    </tlLogic>
</additional>""")


def generateTrips(netFile, tripFile, density):
    import os
    import subprocess

    SUMO_HOME = os.environ.get("SUMO_HOME")
    TOOLS_PATH = os.path.join(SUMO_HOME, "tools")

    if density == "low":
        period = "12"
        #300 vehicles p/h
    elif density == "medium":
        period = "6"
        #600 vehicles p/h
    elif density == "high":
        period = "3"
        #1200 vehicles p/h
    else:
        period = "6"

    tripDir = os.path.dirname(tripFile)
    if tripDir:
        os.makedirs(tripDir, exist_ok=True)

    cmd = [
        "python", os.path.join(TOOLS_PATH, "randomTrips.py"),
        "-n", netFile,
        "-o", tripFile,
        "--prefix", "veh",
        "--seed", "13",
        "--min-distance", "20",
        "--trip-attributes", 'departLane="best" departSpeed="max"',
        "--period", period
    ]

    subprocess.run(cmd, check=True)
    print("Trips generated.")
