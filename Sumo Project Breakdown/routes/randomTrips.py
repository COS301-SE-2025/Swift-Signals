import os
import subprocess

SUMO_HOME = os.environ.get("SUMO_HOME")
TOOLS_PATH = os.path.join(SUMO_HOME, "tools")
NET_FILE = "../network/network.net.xml"
TRIPS_FILE = "trips.trips.xml"

cmd = [
    "python", os.path.join(TOOLS_PATH, "randomTrips.py"),
    "-n", NET_FILE,
    "-o", TRIPS_FILE,
    "--prefix", "veh",
    "--seed", "94",
    "--min-distance", "20",
    "--trip-attributes", 'departLane="best" departSpeed="max"'
]

subprocess.run(cmd, check=True)
print("Trips generated")