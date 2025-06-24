import os
import sys

script_dir = os.path.dirname(__file__)
project_root = os.path.abspath(os.path.join(script_dir, '..', '..', 'protos', 'gen', 'simulation'))

if project_root not in sys.path:
    sys.path.insert(0, project_root)

try:
    import simulation_pb2 as pb
except ImportError:
    # This block handles cases where the import path might not be set up
    # correctly, for example, if running this script directly without
    # the proper package structure. In a well-configured project,
    # the `sys.path.insert` above should prevent this.
    print("ERROR: Could not import protobuf messages. "
          "Please ensure 'protos/gen/simulation' is accessible and `protoc` was run correctly.")
    print("Current sys.path:", sys.path)
    sys.exit(1)


from google.protobuf.json_format import MessageToDict, ParseDict

# --- Functions for converting Protobuf messages to Python Dictionaries ---

def request_proto_to_dict(request_proto: pb.SimulationRequest) -> dict:
    """
    Converts a SimulationRequest protobuf message into a Python dictionary.

    Args:
        request_proto (pb.SimulationRequest): The protobuf message.

    Returns:
        dict: The dictionary representation of the request.
    """
    # `preserving_proto_field_name=True` ensures that field names in the dictionary
    # match the snake_case names defined in the .proto file (e.g., "intersection_id").
    # If set to False (default), it converts them to camelCase (e.g., "intersectionId").
    return MessageToDict(request_proto, preserving_proto_field_name=True)

# --- Functions for converting Python Dictionaries to Protobuf messages ---

def simulation_results_dict_to_proto(response_dict: dict) -> pb.SimulationResultsResponse:
    """
    Converts a Python dictionary into a SimulationResultsResponse protobuf message.

    Args:
        response_dict (dict): The dictionary representation of the response.

    Returns:
        pb.SimulationResultsResponse: The populated SimulationResultsResponse protobuf message.
    """
    response = pb.SimulationResultsResponse()
    ParseDict(response_dict, response, ignore_unknown_fields=False)
    return response

def simulation_output_dict_to_proto(response_dict: dict) -> pb.SimulationOutputResponse:
    """
    Converts a Python dictionary into a SimulationOutputResponse protobuf message.

    Args:
        response_dict (dict): The dictionary representation of the response.

    Returns:
        pb.SimulationOutputResponse: The populated SimulationOutputResponse protobuf message.
    """
    response = pb.SimulationOutputResponse()
    ParseDict(response_dict, response, ignore_unknown_fields=False)
    return response

# --- Example Usage (for testing the functions directly) ---
if __name__ == "__main__":
    print("--- Testing Conversion Functions ---")

    # Example 1: Converting a SimulationRequest Protobuf to dictionary
    print("\n--- Example: SimulationRequest Proto to Dict ---")
    sample_proto_request = pb.SimulationRequest(
        intersection_id="traffic_junction_202",
        parameters=pb.SimulationParameters(
            intersection_type=pb.IntersectionType.INTERSECTION_TYPE_ROUNDABOUT,
            green=30, yellow=3, red=1, speed=20, seed=12345
        )
    )
    print(f"Input Protobuf Request:\n{sample_proto_request}")
    try:
        dict_request = request_proto_to_dict(sample_proto_request)
        print(f"\nConverted Dictionary:\n{dict_request}")
        assert dict_request["intersection_id"] == "traffic_junction_202"
        assert dict_request["parameters"]["intersection_type"] == "INTERSECTION_TYPE_ROUNDABOUT"
        print("\nSUCCESS: SimulationRequest Proto to Dict conversion completed.")
    except Exception as e:
        print(f"\nERROR: Failed to convert SimulationRequest Proto to Dict: {e}")

    # Example 2: Converting a dictionary to SimulationResultsResponse Protobuf
    print("\n--- Example: Dict to SimulationResultsResponse Proto ---")
    sample_results_dict = {
        "total_vehicles": 800,
        "average_travel_time": 85,
        "total_travel_time": 68000,
        "average_speed": 50,
        "average_waiting_time": 8,
        "total_waiting_time": 60,
        "generated_vehicles": 780,
        "emergency_brakes": 3,
        "emergency_stops": 1,
        "near_collisions": 0,
    }
    print(f"Input Dictionary:\n{sample_results_dict}")
    try:
        proto_results_response = simulation_results_dict_to_proto(sample_results_dict)
        print(f"\nConverted Protobuf Response:\n{proto_results_response}")
        assert proto_results_response.total_vehicles == 800
        assert proto_results_response.average_travel_time == 85
        assert proto_results_response.total_travel_time == 68000
        assert proto_results_response.average_speed == 50
        assert proto_results_response.average_waiting_time == 8
        assert proto_results_response.total_waiting_time == 60
        assert proto_results_response.generated_vehicles == 780
        assert proto_results_response.emergency_brakes == 3
        assert proto_results_response.emergency_stops == 1
        assert proto_results_response.near_collisions == 0
        print("\nSUCCESS: Dictionary to SimulationResultsResponse Proto conversion completed.")
    except Exception as e:
        print(f"\nERROR: Failed to convert dictionary to SimulationResultsResponse Proto: {e}")

    # Example 3: Converting a dictionary to SimulationOutputResponse Protobuf
    print("\n--- Example: Dict to SimulationOutputResponse Proto ---")
    sample_output_dict = {
        "intersection": {
            "nodes": [
                {"id": "n_alpha", "x": 5.0, "y": 5.0, "type": "PRIORITY"},
            ],
            "edges": [
                {"id": "e_xy", "from": "n_alpha", "to": "n_beta", "speed": 70.0, "lanes": 3},
            ],
            "connections": [],
            "traffic_lights": []
        },
        "vehicles": [
            {"id": "v_red", "positions": [{"time": 0, "x": 0.0, "y": 0.0, "speed": 0.0}]},
            {"id": "v_green", "positions": [{"time": 0, "x": 0.0, "y": 0.0, "speed": 0.0}]}
        ]
    }
    print(f"Input Dictionary:\n{sample_output_dict}")
    try:
        proto_output_response = simulation_output_dict_to_proto(sample_output_dict)
        print(f"\nConverted Protobuf Response:\n{proto_output_response}")
        assert proto_output_response.intersection.nodes[0].id == "n_alpha"
        assert proto_output_response.vehicles[0].id == "v_red"
        print("\nSUCCESS: Dictionary to SimulationOutputResponse Proto conversion completed.")
    except Exception as e:
        print(f"\nERROR: Failed to convert dictionary to SimulationOutputResponse Proto: {e}")

    print("\n--- All conversion tests concluded. ---")
