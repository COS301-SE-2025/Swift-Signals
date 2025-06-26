import sys, os

sys.path.append(os.path.abspath(os.path.join('..')))
import SimLoad as sim

script_dir = os.path.dirname(__file__)
project_root = os.path.abspath(os.path.join(script_dir, '..', '..', 'protos', 'gen', 'simulation'))
if project_root not in sys.path:
    sys.path.insert(0, project_root)
import simulation_pb2 as pb

import convert


params = sample_proto_request = pb.SimulationRequest(
            intersection_id="traffic_junction_202",
            parameters=pb.SimulationParameters(
                intersection_type=pb.IntersectionType.INTERSECTION_TYPE_ROUNDABOUT,
                green=30, yellow=3, red=1, speed=20, seed=12345
            )
        )

paramsDict = convert.request_proto_to_dict(params)

data = {"intersection": paramsDict}

data["intersection"]["parameters"]["intersection_type"] = 2

# print(data)
# print(data["intersection"]["parameters"])

data = {
    "intersection": {
        "simulation_parameters": {
            "intersection_type": 1,
            "green": 25,
            "yellow": 3,
            "red": 30,
            "speed": 80,
            "seed": 1234,
        }
    }
}

# print("\n============================\n")
#
# print(data)

print(sim.main(data))
