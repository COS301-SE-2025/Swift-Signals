import os
import subprocess


def generate(params):
    base = "stop_intersection"
    netFile = f"{base}.net.xml"
    routeFile = f"{base}.rou.xml"
    configFile = f"{base}.sumocfg"
    stopIntFile = f"{base}.add.xml"

    nodeFile = "stopInt.nod.xml"
    edgeFile = "stopInt.edg.xml"
    conFile = "stopInt.con.xml"

    writeNodeFile(nodeFile)
    writeEdgeFile(edgeFile)
    writeConnectionFile(conFile)
    writeStopLogic(stopIntFile)

    print("Generating stop street intersection with params:", params)

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
            <additional-files value="{stopIntFile}"/>
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
    <node id="n1" x="0" y="0" type="priority"/>
    <node id="n2" x="0" y="100" type="priority"/>
    <node id="n3" x="0" y="-100" type="priority"/>
    <node id="n4" x="-100" y="0" type="priority"/>
    <node id="n5" x="100" y="0" type="priority"/>
</nodes>"""
    with open(filename, "w") as f:
        f.write(content)


def writeEdgeFile(filename):
    content = """<edges>
    <edge id="in_n2_1" from="n2" to="n1" priority="0" numLanes="1" speed="13.9"/>
    <edge id="in_n3_1" from="n3" to="n1" priority="0" numLanes="1" speed="13.9"/>
    <edge id="in_n4_1" from="n4" to="n1" priority="0" numLanes="1" speed="13.9"/>
    <edge id="in_n5_1" from="n5" to="n1" priority="0" numLanes="1" speed="13.9"/>

    <edge id="out_n1_n2" from="n1" to="n2" priority="0" numLanes="1" speed="13.9"/>
    <edge id="out_n1_n3" from="n1" to="n3" priority="0" numLanes="1" speed="13.9"/>
    <edge id="out_n1_n4" from="n1" to="n4" priority="0" numLanes="1" speed="13.9"/>
    <edge id="out_n1_n5" from="n1" to="n5" priority="0" numLanes="1" speed="13.9"/>
</edges>"""
    with open(filename, "w") as f:
        f.write(content)


def writeConnectionFile(filename):
    content = """<connections>
    <connection from="in_n2_1" to="out_n1_n3" fromLane="0" toLane="0"/>
    <connection from="in_n2_1" to="out_n1_n5" fromLane="0" toLane="0"/>
    <connection from="in_n2_1" to="out_n1_n4" fromLane="0" toLane="0"/>

    <connection from="in_n3_1" to="out_n1_n4" fromLane="0" toLane="0"/>
    <connection from="in_n3_1" to="out_n1_n5" fromLane="0" toLane="0"/>
    <connection from="in_n3_1" to="out_n1_n2" fromLane="0" toLane="0"/>

    <connection from="in_n4_1" to="out_n1_n5" fromLane="0" toLane="0"/>
    <connection from="in_n4_1" to="out_n1_n2" fromLane="0" toLane="0"/>
    <connection from="in_n4_1" to="out_n1_n3" fromLane="0" toLane="0"/>

    <connection from="in_n5_1" to="out_n1_n2" fromLane="0" toLane="0"/>
    <connection from="in_n5_1" to="out_n1_n3" fromLane="0" toLane="0"/>
    <connection from="in_n5_1" to="out_n1_n4" fromLane="0" toLane="0"/>
</connections>"""
    with open(filename, "w") as f:
        f.write(content)


def writeStopLogic(filename):
    with open(filename, "w") as tl:
        tl.write("""<additional>
    <priority id="priority0" type="stop" lane="in_n2_1_0" startPos="0" endPos="5"/>
    <priority id="priority1" type="stop" lane="in_n3_1_0" startPos="0" endPos="5"/>
    <priority id="priority2" type="stop" lane="in_n4_1_0" startPos="0" endPos="5"/>
    <priority id="priority3" type="stop" lane="in_n5_1_0" startPos="0" endPos="5"/>
</additional>""")


def generateTrips(netFile, tripFile, density):
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
