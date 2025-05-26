import os
import traci
import sys
import json

sumo_binary = "sumo-gui" #sumo-gui for debugging
sumo_config = "config/simulation.sumocfg"

traci.start([sumo_binary, "-c", sumo_config])
print("Simulation started")

step = 0
vehicle_log = {}
traffic_log = {}

traffic_light_ids = traci.trafficlight.getIDList()
print("Traffic lights:", traffic_light_ids)

while step < 100:
    traci.simulationStep()

    vehicle_ids = traci.vehicle.getIDList()
    for vid in vehicle_ids:
        pos = traci.vehicle.getPosition(vid)
        speed = traci.vehicle.getSpeed(vid)
        vehicle_log.setdefault(vid, []).append((step, pos, speed))
        print(f"Step {step} - {vid}: Pos={pos}, Speed={speed:.2f}")

    for tl in traffic_light_ids:
        phase = traci.trafficlight.getPhase(tl)
        state = traci.trafficlight.getRedYellowGreenState(tl)
        traffic_log.setdefault(tl, []).append((step, phase, state))
        print(f"Step {step} - TLS {tl}: Phase={phase}, State={state}")

    step += 1

traci.close()
print("Simulation complete")

with open("vehicle_log.json", "w") as f:
    json.dump(vehicle_log, f, indent=2)

with open("traffic_log.json", "w") as f:
    json.dump(traffic_log, f, indent=2)