import os
import subprocess


def generate(params):
    base = "tl_intersection"
    netFile = f"{base}.net.xml"
    routeFile = f"{base}.rou.xml"
    configFile = f"{base}.sumocfg"

    nodeFile = "tlInt.nod.xml"
    edgeFile = "tlInt.edg.xml"
    conFile = "tlInt.con.xml"

    writeNodeFile(nodeFile)
    writeEdgeFile(edgeFile)
    writeConnectionFile(conFile)

    print("Generating traffic light intersection with params:", params)

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
    <node id="r1" x="0" y="20" type="priority"/>
    <node id="r2" x="20" y="0" type="priority"/>
    <node id="r3" x="0" y="-20" type="priority"/>
    <node id="r4" x="-20" y="0" type="priority"/>

    <node id="nN" x="0" y="100" type="priority"/>
    <node id="nE" x="100" y="0" type="priority"/>
    <node id="nS" x="0" y="-100" type="priority"/>
    <node id="nW" x="-100" y="0" type="priority"/>
</nodes>"""
    with open(filename, "w") as f:
        f.write(content)


def writeEdgeFile(filename):
    content = """<edges>
    <edge id="in_N_r1" from="nN" to="r1" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_E_r2" from="nE" to="r2" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_S_r3" from="nS" to="r3" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_W_r4" from="nW" to="r4" priority="1" numLanes="1" speed="13.9"/>

    <edge id="out_r1_N" from="r1" to="nN" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_r2_E" from="r2" to="nE" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_r3_S" from="r3" to="nS" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_r4_W" from="r4" to="nW" priority="1" numLanes="1" speed="13.9"/>

    <edge id="r1_r2" from="r1" to="r2" priority="2" numLanes="1" speed="8.3"/>
    <edge id="r2_r3" from="r2" to="r3" priority="2" numLanes="1" speed="8.3"/>
    <edge id="r3_r4" from="r3" to="r4" priority="2" numLanes="1" speed="8.3"/>
    <edge id="r4_r1" from="r4" to="r1" priority="2" numLanes="1" speed="8.3"/>
</edges>"""
    with open(filename, "w") as f:
        f.write(content)


def writeConnectionFile(filename):
    content = """<connections>
    <!-- Incoming edges yield to circulating edges -->
    <connection from="in_N_r1" to="r1_r2" fromLane="0" toLane="0" priority="1"/>
    <connection from="in_E_r2" to="r2_r3" fromLane="0" toLane="0" priority="1"/>
    <connection from="in_S_r3" to="r3_r4" fromLane="0" toLane="0" priority="1"/>
    <connection from="in_W_r4" to="r4_r1" fromLane="0" toLane="0" priority="1"/>

    <!-- Circulating edges have higher priority -->
    <connection from="r1_r2" to="r2_r3" fromLane="0" toLane="0" priority="2"/>
    <connection from="r2_r3" to="r3_r4" fromLane="0" toLane="0" priority="2"/>
    <connection from="r3_r4" to="r4_r1" fromLane="0" toLane="0" priority="2"/>
    <connection from="r4_r1" to="r1_r2" fromLane="0" toLane="0" priority="2"/>

    <!-- Exiting edges yield to circulating edges -->
    <connection from="r1_r2" to="out_r2_E" fromLane="0" toLane="0" priority="1"/>
    <connection from="r2_r3" to="out_r3_S" fromLane="0" toLane="0" priority="1"/>
    <connection from="r3_r4" to="out_r4_W" fromLane="0" toLane="0" priority="1"/>
    <connection from="r4_r1" to="out_r1_N" fromLane="0" toLane="0" priority="1"/>
</connections>"""
    with open(filename, "w") as f:
        f.write(content)


def generateTrips(netFile, routeFile, density):
    import subprocess

    SUMO_HOME = os.environ.get("SUMO_HOME")
    TOOLS_PATH = os.path.join(SUMO_HOME, "tools")

    if density == "low":
        period = "10"
    elif density == "medium":
        period = "5"
    elif density == "high":
        period = "2"
    else:
        period = "5"

    tripsFile = routeFile.replace(".rou.xml", ".trips.xml")

    cmd_trips = [
        "python", os.path.join(TOOLS_PATH, "randomTrips.py"),
        "-n", netFile,
        "-o", tripsFile,
        "--prefix", "veh",
        "--seed", "13",
        "--min-distance", "50",
        "--trip-attributes", 'departLane="best" departSpeed="max"',
        "--period", period,
        "--start-edges", "in_N_r1,in_E_r2,in_S_r3,in_W_r4",
        "--end-edges", "out_r1_N,out_r2_E,out_r3_S,out_r4_W"
    ]

    subprocess.run(cmd_trips, check=True)
    print("Trips generated.")

    cmd_route = [
        "duarouter",
        "-n", netFile,
        "-t", tripsFile,
        "-o", routeFile,
        "--ignore-errors"
    ]

    subprocess.run(cmd_route, check=True)
    print("Routes generated by duarouter.")
