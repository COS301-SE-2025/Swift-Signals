from concurrent import futures
import logging
import datetime
import grpc
import sys, os

sys.path.append(os.path.abspath(os.path.join('..', '..', 'protos', 'gen', 'simulation')))

import simulation_pb2 as pb
import simulation_pb2_grpc as pb_grpc

class SimulationServicer(pb_grpc.SimulationServiceServicer):
    def GetSimulationResults(self, request, context):
        """
        Runs the simulation and returns the result.
        """
        print(f"Received GetSimulationResults request for intersection_id: {request.intersection_id} with parameters: {request.parameters}")

        # TODO: Implement your actual simulation logic here.
        # This is where you would process the request.parameters,
        # run your simulation model, and calculate the results.
        # For demonstration, we'll return some dummy data.

        current_time = datetime.datetime.now()
        results = pb.SimulationResultsResponse(
            id=f"sim-result-{request.intersection_id}-{current_time.strftime('%Y%m%d-%H%M%S')}",
            total_vehicles=1500,
            avg_travel_time=90,
            total_travel_time=135000,
            avg_speed=45,
            avg_waiting_time=10,
            waiting_time=15000,
            generated_vehicles=1400,
            emergency_brakes=7,
            emergency_stops=3,
            near_collisions=2,
            date_run=current_time.isoformat(),
        )
        print(f"Returning GetSimulationResults: {results.id}")
        return results

    def GetSimulationOutput(self, request, context):
        """
        Runs the simulation with detailed output.
        """
        print(f"Received GetSimulationOutput request for intersection_id: {request.intersection_id} with parameters: {request.parameters}")

        # TODO: Implement your actual detailed simulation logic here.
        # This would involve simulating the intersection and vehicle movements
        # step-by-step and capturing the detailed state.
        # For demonstration, we'll return some dummy data.

        # Dummy Intersection data
        intersection = pb.Intersection(
            nodes=[
                pb.Node(id="n1", x=0.0, y=0.0, type=pb.NodeType.TRAFFIC_LIGHT),
                pb.Node(id="n2", x=100.0, y=0.0, type=pb.NodeType.PRIORITY),
                pb.Node(id="n3", x=0.0, y=100.0, type=pb.NodeType.PRIORITY),
            ],
            edges=[
                pb.Edge(id="e1", to="n2", speed=50.0, lanes=2),
                pb.Edge(id="e2", to="n3", speed=40.0, lanes=1),
            ],
            connections=[
                pb.Connection(to="e2", fromLane=0, toLane=0, tl=0),
            ],
            traffic_lights=[
                pb.TrafficLight(
                    id="tl_main",
                    type="standard",
                    phases=[
                        pb.Phase(duration=30, state="GGrr"),
                        pb.Phase(duration=5, state="yyRr"),
                        pb.Phase(duration=30, state="rrGG"),
                        pb.Phase(duration=5, state="rryy"),
                    ],
                )
            ],
        )

        # Dummy Vehicle data
        vehicles = [
            pb.Vehicle(
                id="car_A01",
                positions=[
                    pb.Position(time=0, x=0.0, y=0.0, speed=0.0),
                    pb.Position(time=5, x=10.0, y=0.0, speed=5.0),
                    pb.Position(time=10, x=25.0, y=0.0, speed=7.0),
                    pb.Position(time=15, x=45.0, y=0.0, speed=8.0),
                ],
            ),
            pb.Vehicle(
                id="truck_B02",
                positions=[
                    pb.Position(time=0, x=0.0, y=10.0, speed=0.0),
                    pb.Position(time=10, x=15.0, y=10.0, speed=3.0),
                    pb.Position(time=20, x=30.0, y=10.0, speed=4.0),
                ],
            ),
        ]

        output = pb.SimulationOutputResponse(
            intersection=intersection,
            vehicles=vehicles,
        )
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
