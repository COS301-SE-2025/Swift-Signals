import os
import logging

import grpc
from simulation_pb2 import (
    SimulationRequest,
    SimulationParameters,
    IntersectionType,
    SimulationResultsResponse,
)
from simulation_pb2_grpc import SimulationServiceStub

# Configuration
SIMU_GRPC_ADDR = os.environ.get("SIMU_GRPC_ADDR", "localhost:50053")
logging.info(f"Connecting to gRPC server at {SIMU_GRPC_ADDR}")

# Client
channel = grpc.insecure_channel(SIMU_GRPC_ADDR)
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
