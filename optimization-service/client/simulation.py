import grpc
from google.protobuf.json_format import MessageToDict
import os
import sys

# Add current directory to path for protobuf imports
current_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(current_dir)

import simulation_pb2 as pb
import simulation_pb2_grpc as pb_grpc


class SimulationClient:
    def __init__(self, server_address=None):
        if server_address is None:
            server_address = os.environ.get("SIMULATION_SERVER_ADDRESS", "localhost:50053")
        self.channel = grpc.insecure_channel(server_address)
        self.stub = pb_grpc.SimulationServiceStub(self.channel)

    def get_simulation_results(self, green: int, yellow: int, red: int, speed: int, seed: int, intersection_id: str = ""):
        """
        Get simulation results using gRPC call

        Returns:
            dict: Simulation results in dictionary format
        """
        try:
            request = pb.SimulationRequest(
                intersection_id=intersection_id,
                simulation_parameters=pb.SimulationParameters(
                    intersection_type=pb.INTERSECTION_TYPE_TRAFFICLIGHT,
                    green=green,
                    yellow=yellow,
                    red=red,
                    speed=speed,
                    seed=seed,
                ),
            )

            response = self.stub.GetSimulationResults(request)

            # Convert protobuf response to dictionary
            result_dict = MessageToDict(
                response,
                preserving_proto_field_name=True,
                use_integers_for_enums=True
            )

            # Transform to match the expected format from the original file-based approach
            return {
                "Total Vehicles": result_dict.get("total_vehicles", 0),
                "Average Travel Time": result_dict.get("average_travel_time", 0.0),
                "Total Travel Time": result_dict.get("total_travel_time", 0.0),
                "Average Speed": result_dict.get("average_speed", 0.0),
                "Average Waiting Time": result_dict.get("average_waiting_time", 0.0),
                "Total Waiting Time": result_dict.get("total_waiting_time", 0.0),
                "Generated Vehicles": result_dict.get("generated_vehicles", 0),
                "Emergency Brakes": result_dict.get("emergency_brakes", 0),
                "Emergency Stops": result_dict.get("emergency_stops", 0),
                "Near collisions": result_dict.get("near_collisions", 0),
            }

        except grpc.RpcError as e:
            print(f"gRPC error occurred: {e}")
            return None
        except Exception as e:
            print(f"Error in simulation request: {e}")
            return None

    def close(self):
        """Close the gRPC channel"""
        if hasattr(self, 'channel'):
            self.channel.close()


# Legacy function for backward compatibility
def GetSimulationResults(green: int, yellow: int, red: int, speed: int, seed: int):
    """
    Legacy function for backward compatibility
    """
    client = SimulationClient()
    try:
        return client.get_simulation_results(green, yellow, red, speed, seed)
    finally:
        client.close()
