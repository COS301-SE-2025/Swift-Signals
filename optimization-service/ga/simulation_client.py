import os

import grpc
from simulation_pb2 import (
    SimulationRequest,
    SimulationParameters,
    IntersectionType,
    SimulationResultsResponse,
)
from simulation_pb2_grpc import SimulationServiceStub

# Configuration
GRPC_SIM_SERVER_ADDRESS = os.environ.get("GRPC_SIM_SERVER_ADDRESS", "localhost:50053")

# Client
channel = grpc.insecure_channel(GRPC_SIM_SERVER_ADDRESS)
stub = SimulationServiceStub(channel)


def run_simulation(individual_params: list) -> dict | None:
    """
    Runs a simulation using the gRPC service and returns a dictionary of results.

    Args:
        individual_params: A list representing the individual's parameters:
                           [green, yellow, red, speed, seed]

    Returns:
        A dictionary with simulation results or None if the simulation fails.
    """
    green, yellow, red, speed, seed = individual_params

    params = SimulationParameters(
        intersection_type=IntersectionType.INTERSECTION_TYPE_TRAFFICLIGHT,
        green=green,
        yellow=yellow,
        red=red,
        speed=speed,
        seed=seed,
    )
    request = SimulationRequest(
        intersection_id="TrafficLight",
        simulation_parameters=params,
    )

    try:
        response: SimulationResultsResponse = stub.GetSimulationResults(request)
        results_dict = {
            "total_waiting_time": response.total_waiting_time,
            "total_travel_time": response.total_travel_time,
            "emergency_brakes": response.emergency_brakes,
            "emergency_stops": response.emergency_stops,
            "near_collisions": response.near_collisions,
        }
        return results_dict
    except grpc.RpcError as e:
        print(f"gRPC call failed: {e}")
        return None
