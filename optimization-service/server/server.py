from concurrent import futures

from google.protobuf.json_format import MessageToDict, ParseDict

import grpc
from grpc_reflection.v1alpha import reflection
import optimisation_pb2 as pb
import optimisation_pb2_grpc as pb_grpc


class OptimisationServicer(pb_grpc.OptimisationServiceServicer):
    def RunOptimisation(self, request, context):
        print("RunOptimisationResult request received with intersection_type:",
              request.parameters.intersection_type)
        print(MessageToDict(request, preserving_proto_field_name=True,
              use_integers_for_enums=True))
        return request


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb_grpc.add_OptimisationServiceServicer_to_server(
        OptimisationServicer(), server)

    SERVICE_NAMES = (
        pb.DESCRIPTOR.services_by_name["OptimisationService"].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    server.add_insecure_port("[::]:50054")
    server.start()
    print("Optimisation Service listening on port :50054")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
