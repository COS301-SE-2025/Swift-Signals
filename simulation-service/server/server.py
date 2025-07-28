from multiprocessing import context
from concurrent import futures

from google.protobuf.json_format import MessageToDict, ParseDict
import grpc
from grpc_reflection.v1alpha import reflection
import simulation_pb2 as pb
import simulation_pb2_grpc as pb_grpc

import SimLoad


class SimulationServicer(pb_grpc.SimulationServiceServicer):
    def GetSimulationResults(self, request, context):
        print("GetSimulationResults request received with id:",
              request.intersection_id)

        req_dict = {
            "intersection": MessageToDict(request,
                                          preserving_proto_field_name=True,
                                          use_integers_for_enums=True)
        }
        results = SimLoad.main(
            req_dict)[0]["intersection"]["results"]
        msg_results = pb.SimulationResultsResponse()
        ParseDict(results, msg_results)

        return msg_results

    def GetSimulationOutput(self, request, context):
        print("GetSimulationOutput request received with id:",
              request.intersection_id)

        req_dict = {
            "intersection": MessageToDict(request,
                                          preserving_proto_field_name=True,
                                          use_integers_for_enums=True)
        }
        results = SimLoad.main(
            req_dict)[1]
        msg_results = pb.SimulationOutputResponse()
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
