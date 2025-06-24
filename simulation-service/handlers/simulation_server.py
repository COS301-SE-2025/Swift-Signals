from concurrent import futures
import logging
import grpc
import sys, os

import convert

sys.path.append(os.path.abspath(os.path.join('..', '..', 'protos', 'gen', 'simulation')))

import simulation_pb2 as pb
import simulation_pb2_grpc as pb_grpc

sys.path.append(os.path.abspath(os.path.join('..')))

import SimLoad as sim

class SimulationServicer(pb_grpc.SimulationServiceServicer):
    def GetSimulationResults(self, request, context):
        """
        Runs the simulation and returns the result.
        """
        print(f"Received GetSimulationResults request for intersection_id: {request.intersection_id} with parameters: {request.parameters}")

        results_dict = sim.main(convert.request_proto_to_dict(request))
        results = convert.proto_results_response(results_dict)

        print(f"Returning GetSimulationResults: {results}")
        return results

    def GetSimulationOutput(self, request, context):
        """
        Runs the simulation with detailed output.
        """
        print(f"Received GetSimulationOutput request for intersection_id: {request.intersection_id} with parameters: {request.parameters}")

        output_dict = sim.main(convert.request_proto_to_dict(request))
        output = convert.proto_results_response(output_dict)

        print("Returning GetSimulationOutput with dummy data.")
        return output

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb_grpc.add_SimulationServiceServicer_to_server(
        SimulationServicer(), server
    )
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
