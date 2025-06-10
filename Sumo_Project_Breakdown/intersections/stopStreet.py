import os
import subprocess

def generate(params):
    base = "stop_intersection"
    netFile = f"{base}.net.xml"
    routeFile = f"{base}.rou.xml"
    configFile = f"{base}.sumocfg"

    nodeFile = "stopInt.nod.xml"
    edgeFile = "stopInt.edg.xml"
    conFile = "stopInt.con.xml"

    writeNodeFile(nodeFile)
    writeEdgeFile(edgeFile)
    writeConnectionFile(conFile)

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