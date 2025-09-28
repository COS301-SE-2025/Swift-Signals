from concurrent import futures
import logging
import os
from pprint import pformat

from google.protobuf.json_format import MessageToDict, ParseDict
import grpc
from grpc_reflection.v1alpha import reflection
from swiftsignals.simulation.v1 import simulation_pb2 as pb
from swiftsignals.simulation.v1 import simulation_pb2_grpc as pb_grpc

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
            "%s: type=dict, keys=%d -> %s", name, len(obj), list(obj.keys())[:5]
        )
    else:
        logger.debug("%s: type=%s, value=%s", name, type(obj).__name__, str(obj)[:100])


def pretty_log(name, obj, max_len=1000):
    """Log a pretty representation of an object, truncated to max_len."""
    if obj is None:
        logger.debug("%s: None", name)
        return
    try:
        text = pformat(obj, width=80, compact=True)
    except Exception as e:
        text = f"<Failed to format object: {e}>"
    if len(text) > max_len:
        text = text[:max_len] + "\n...[truncated]"
    logger.debug("%s:\n%s", name, text)


class SimulationServicer(pb_grpc.SimulationServiceServicer):
    def GetSimulationResults(self, request: pb.SimulationRequest, context):
        logger.info(
            "Received GetSimulationResults request",
            extra={"intersection_id": request.intersection_id},
        )

        req_dict = {
            "intersection": MessageToDict(
                request, preserving_proto_field_name=True, use_integers_for_enums=True
            )
        }
        req_dict["intersection"]["traffic density"] = 1
        req_dict["intersection"]["simulation_parameters"]["Green"] = req_dict[
            "intersection"
        ]["simulation_parameters"]["green"]
        req_dict["intersection"]["simulation_parameters"]["Yellow"] = req_dict[
            "intersection"
        ]["simulation_parameters"]["yellow"]
        req_dict["intersection"]["simulation_parameters"]["Red"] = req_dict[
            "intersection"
        ]["simulation_parameters"]["red"]
        req_dict["intersection"]["simulation_parameters"]["Speed"] = req_dict[
            "intersection"
        ]["simulation_parameters"]["speed"]
        req_dict["intersection"]["simulation_parameters"]["Seed"] = req_dict[
            "intersection"
        ]["simulation_parameters"]["seed"]

        pretty_log("Request dict", req_dict)

        try:
            sim_output = SimLoad.main(req_dict)
            pretty_log("Raw simulation output", sim_output)

            if not sim_output or sim_output[0] is None:
                msg = f"Simulation returned no results for intersection_id={request.intersection_id}"
                logger.error(msg)
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(msg)
                return pb.SimulationResultsResponse()

            results = sim_output[0]["intersection"]["results"]
            pretty_log("Parsed results", results)

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
        pretty_log("Request dict", req_dict)

        try:
            sim_output = SimLoad.main(req_dict)
            pretty_log("Raw simulation output", sim_output)

            if not sim_output or len(sim_output) < 2:
                msg = f"Simulation returned no output for intersection_id={request.intersection_id}"
                logger.error(msg)
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(msg)
                return pb.SimulationOutputResponse()

            results = sim_output[1]
            pretty_log("Parsed output", results)

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
    logger.info(f"Simulation Service listening on port :{port}")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
