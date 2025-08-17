import os

import grpc
from google.protobuf.json_format import MessageToDict

import simulation_pb2 as pb
import simulation_pb2_grpc as pb_grpc


server_address = os.environ.get("SIMULATION_SERVER_ADDRESS", "localhost:50053")
channel = grpc.insecure_channel(server_address)
stub = pb_grpc.SimulationServiceStub(channel)


def get_simulation_results(green: int, yellow: int, red: int, speed: int, seed: int):
    results = stub.GetSimulationResults(
        pb.SimulationRequest(
            intersection_id="",
            simulation_parameters=pb.SimulationParameters(
                intersection_type=pb.INTERSECTION_TYPE_TRAFFICLIGHT,
                green=green,
                yellow=yellow,
                red=red,
                speed=speed,
                seed=seed,
            ),
        )
    )

    return MessageToDict(
        results, preserving_proto_field_name=True, use_integers_for_enums=True
    )
