from concurrent import futures
import os
import sys
import traceback

from google.protobuf.json_format import MessageToDict

import grpc
from grpc_reflection.v1alpha import reflection
import optimisation_pb2 as pb
import optimisation_pb2_grpc as pb_grpc

# Add parent directory to path to import GP module
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from GP import OptimizationEngine


class OptimisationServicer(pb_grpc.OptimisationServiceServicer):
    def __init__(self):
        self.simulation_server_address = os.environ.get("SIMULATION_SERVER_ADDRESS", "localhost:50053")

    def RunOptimisation(self, request, context):
        """
        Run optimization using the refactored OptimizationEngine
        """
        print(
            "RunOptimisation request received with intersection_type:",
            request.parameters.intersection_type,
        )
        print(
            MessageToDict(
                request, preserving_proto_field_name=True, use_integers_for_enums=True
            )
        )

        try:
            # Extract optimization parameters from request
            optimization_type = request.optimisation_type
            params = request.parameters

            # For now, we'll use the genetic algorithm optimization
            # In the future, you could add different optimization types based on the request
            if optimization_type == pb.OPTIMISATION_TYPE_GENETIC_EVALUATION:
                # Initialize the optimization engine
                engine = OptimizationEngine(simulation_server_address=self.simulation_server_address)

                # Run optimization with default parameters
                # You could extract these from the request if needed
                optimization_results = engine.run_optimization(
                    ngen_waiting=30,
                    ngen_safety=10,
                    pop_size=30,
                    cxpb=0.5,
                    mutpb=0.3,
                    random_seed=params.seed if params.seed else 1408
                )

                # Check if optimization was successful
                if "error" in optimization_results:
                    context.set_code(grpc.StatusCode.INTERNAL)
                    context.set_details(f"Optimization failed: {optimization_results['error']}")
                    return pb.OptimisationParameters()

                # Extract best parameters from optimization results
                best_params = optimization_results.get("best_parameters", {})

                # Create response with optimized parameters
                response = pb.OptimisationParameters(
                    optimisation_type=optimization_type,
                    parameters=pb.SimulationParameters(
                        intersection_type=params.intersection_type,
                        green=best_params.get("Green", params.green),
                        yellow=best_params.get("Yellow", params.yellow),
                        red=best_params.get("Red", params.red),
                        speed=best_params.get("Speed", params.speed),
                        seed=best_params.get("Seed", params.seed),
                    )
                )

                print("Optimization completed successfully")
                print(f"Best parameters: {best_params}")
                return response

            else:
                # For other optimization types, return the original parameters
                print(f"Optimization type {optimization_type} not implemented, returning original parameters")
                return request

        except Exception as e:
            print(f"Error during optimization: {e}")
            print(traceback.format_exc())
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal server error: {str(e)}")
            return pb.OptimisationParameters()


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
    print(f"Optimisation Service listening on port :{port}")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
