import os
import subprocess


def generate(params):
    base = "t_junction"
    netFile = f"{base}.net.xml"
    routeFile = f"{base}.rou.xml"
    configFile = f"{base}.sumocfg"

    nodeFile = "tjInt.nod.xml"
    edgeFile = "tjInt.edg.xml"
    conFile = "tjInt.con.xml"

    writeNodeFile(nodeFile)
    writeEdgeFile(edgeFile)
    writeConnectionFile(conFile)

    print("Generating T-junction with params:", params)

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
    </input>
    <time>
        <begin value="0"/>
        <end value="1000"/>
    </time>
</configuration>""")

    print("Running SUMO simulation...")
    subprocess.run(["sumo-gui", "-c", configFile])


def writeNodeFile(filename):
    content = """<nodes>
    <node id="J1" x="0" y="0" type="priority"/> <!-- T-junction -->
    <node id="nNorth" x="0" y="100" type="priority"/>
    <node id="nEast" x="100" y="0" type="priority"/>
    <node id="nWest" x="-100" y="0" type="priority"/>
</nodes>"""
    with open(filename, "w") as f:
        f.write(content)


def writeEdgeFile(filename):
    content = """<edges>
    <edge id="in_nNorth_J1" from="nNorth" to="J1" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_nEast_J1" from="nEast" to="J1" priority="3" numLanes="1" speed="13.9"/>
    <edge id="in_nWest_J1" from="nWest" to="J1" priority="3" numLanes="1" speed="13.9"/>

    <edge id="out_J1_nNorth" from="J1" to="nNorth" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_J1_nEast" from="J1" to="nEast" priority="3" numLanes="1" speed="13.9"/>
    <edge id="out_J1_nWest" from="J1" to="nWest" priority="3" numLanes="1" speed="13.9"/>
</edges>"""
    with open(filename, "w") as f:
        f.write(content)


def writeConnectionFile(filename):
    content = """<connections>
    <connection from="in_nNorth_J1" to="out_J1_nEast" fromLane="0" toLane="0"/>
    <connection from="in_nNorth_J1" to="out_J1_nWest" fromLane="0" toLane="0"/>

    <connection from="in_nEast_J1" to="out_J1_nNorth" fromLane="0" toLane="0"/>
    <connection from="in_nEast_J1" to="out_J1_nWest" fromLane="0" toLane="0"/>

    <connection from="in_nWest_J1" to="out_J1_nNorth" fromLane="0" toLane="0"/>
    <connection from="in_nWest_J1" to="out_J1_nEast" fromLane="0" toLane="0"/>
</connections>"""
    with open(filename, "w") as f:
        f.write(content)


def generateTrips(netFile, tripFile, density):
    SUMO_HOME = os.environ.get("SUMO_HOME")
    TOOLS_PATH = os.path.join(SUMO_HOME, "tools")

    period = {
        "low": "10",
        "medium": "5",
        "high": "2"
    }.get(density, "5")

    tripDir = os.path.dirname(tripFile)
    if tripDir:
        os.makedirs(tripDir, exist_ok=True)

    cmd = [
        "python", os.path.join(TOOLS_PATH, "randomTrips.py"),
        "-n", netFile,
        "-o", tripFile,
        "--prefix", "veh",
        "--seed", "42",
        "--min-distance", "20",
        "--trip-attributes", 'departLane="best" departSpeed="max"',
        "--period", period
    ]

    subprocess.run(cmd, check=True)
    print("Trips generated.")
