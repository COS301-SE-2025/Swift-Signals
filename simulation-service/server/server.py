from concurrent import futures

from google.protobuf.json_format import MessageToDict, ParseDict
import grpc
from grpc_reflection.v1alpha import reflection
import simulation_pb2 as pb
import simulation_pb2_grpc as pb_grpc

from simload import service


class SimulationServicer(pb_grpc.SimulationServiceServicer):
    def GetSimulationResults(self, request, context):
        print("GetSimulationResults request received with id:",
              request.intersection_id)

        req_dict = MessageToDict(request)
        print(req_dict)
        print("-------------------------------------------")
        print(service.run_simulation(req_dict))
        print("===========================================")
        results = service.run_simulation(
            req_dict)["summary"]["intersection"]["results"]
        msg_results = pb.SimulationResultsResponse()
        ParseDict(results, msg_results)

        return msg_results


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb_grpc.add_SimulationServiceServicer_to_server(
        SimulationServicer(), server)

    SERVICE_NAMES = (
        pb.DESCRIPTOR.services_by_name["SimulationService"].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    server.add_insecure_port("[::]:50053")
    server.start()
    print("Simulation Service listening on port :50053")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
