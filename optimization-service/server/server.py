from concurrent import futures
import os
import logging
from pprint import pformat

import grpc
from grpc_reflection.v1alpha import reflection
import optimisation_pb2 as pb
import optimisation_pb2_grpc as pb_grpc

from ga.main import main as run_optimisation

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format="%(asctime)s | %(levelname)s | %(name)s | %(message)s",
)
logger = logging.getLogger("optimisation-service")


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


class OptimisationServicer(pb_grpc.OptimisationServiceServicer):
    def RunOptimisation(self, request, context):
        logger.info(
            f"Received RunOptimisation request with intersection_type: {request.parameters.intersection_type}",
        )
        pretty_log("Request Parameters", request.parameters)

        # Call the main optimisation function
        result = run_optimisation()
        logger.info("Optimisation completed successfully.")

        # Convert the result to a dictionary
        response = pb.OptimisationParameters(
            optimisation_type=pb.OptimisationType.OPTIMISATION_TYPE_GENETIC_EVALUATION,
            parameters=pb.SimulationParameters(
                intersection_type=pb.IntersectionType.INTERSECTION_TYPE_TRAFFICLIGHT,
                green=result["Green"],
                yellow=result["Yellow"],
                red=result["Red"],
                speed=result["Speed"],
                seed=result["Seed"],
            ),
        )
        pretty_log("Response Parameters", response.parameters)

        return response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb_grpc.add_OptimisationServiceServicer_to_server(OptimisationServicer(), server)

    SERVICE_NAMES = (
        pb.DESCRIPTOR.services_by_name["OptimisationService"].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    port = os.environ.get("APP_PORT", "50054")
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    logger.info(f"Optimisation Service listening on port :{port}")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
