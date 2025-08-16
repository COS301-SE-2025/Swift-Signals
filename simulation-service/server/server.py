from concurrent import futures
import logging
import os

from google.protobuf.json_format import MessageToDict, ParseDict
import grpc
from grpc_reflection.v1alpha import reflection
import simulation_pb2 as pb
import simulation_pb2_grpc as pb_grpc

import SimLoad


# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format="%(asctime)s | %(levelname)s | %(name)s | %(message)s",
)
logger = logging.getLogger("simulation-service")


def log_object_stats(name, obj):
    """Log type + size stats of an object without dumping everything."""
    if obj is None:
        logger.debug("%s: None", name)
    elif isinstance(obj, (list, tuple)):
        logger.debug("%s: type=%s, length=%d", name, type(obj).__name__, len(obj))
    elif isinstance(obj, dict):
        logger.debug(
            "%s: type=dict, keys=%d -> %s", name, len(obj.keys()), list(obj.keys())[:5]
        )
    else:
        logger.debug("%s: type=%s, value=%s", name, type(obj).__name__, str(obj)[:100])


class SimulationServicer(pb_grpc.SimulationServiceServicer):
    def GetSimulationResults(self, request, context):
        logger.info(
            "Received GetSimulationResults request",
            extra={"intersection_id": request.intersection_id},
        )

        req_dict = {
            "intersection": MessageToDict(
                request, preserving_proto_field_name=True, use_integers_for_enums=True
            )
        }
        log_object_stats("Request dict", req_dict)

        try:
            sim_output = SimLoad.main(req_dict)
            log_object_stats("Raw simulation output", sim_output)

            if not sim_output or sim_output[0] is None:
                msg = f"Simulation returned no results for intersection_id={request.intersection_id}"
                logger.error(msg)
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(msg)
                return pb.SimulationResultsResponse()

            results = sim_output[0]["intersection"]["results"]
            log_object_stats("Parsed results", results)

            msg_results = pb.SimulationResultsResponse()
            ParseDict(results, msg_results)

            logger.info(
                "Returning simulation results",
                extra={"intersection_id": request.intersection_id},
            )
            return msg_results

        except Exception as e:
            logger.exception("Error while processing GetSimulationResults")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return pb.SimulationResultsResponse()

    def GetSimulationOutput(self, request, context):
        logger.info(
            "Received GetSimulationOutput request",
            extra={"intersection_id": request.intersection_id},
        )

        req_dict = {
            "intersection": MessageToDict(
                request, preserving_proto_field_name=True, use_integers_for_enums=True
            )
        }
        log_object_stats("Request dict", req_dict)

        try:
            sim_output = SimLoad.main(req_dict)
            log_object_stats("Raw simulation output", sim_output)

            if not sim_output or len(sim_output) < 2:
                msg = f"Simulation returned no output for intersection_id={request.intersection_id}"
                logger.error(msg)
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(msg)
                return pb.SimulationOutputResponse()

            results = sim_output[1]
            log_object_stats("Parsed output", results)

            msg_results = pb.SimulationOutputResponse()
            ParseDict(results, msg_results)

            logger.info(
                "Returning simulation output",
                extra={"intersection_id": request.intersection_id},
            )
            return msg_results

        except Exception as e:
            logger.exception("Error while processing GetSimulationOutput")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return pb.SimulationOutputResponse()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb_grpc.add_SimulationServiceServicer_to_server(SimulationServicer(), server)

    SERVICE_NAMES = (
        pb.DESCRIPTOR.services_by_name["SimulationService"].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    port = os.environ.get("APP_PORT", "50053")
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    logger.info("Simulation Service listening on port :50053")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
