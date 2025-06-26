from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class IntersectionType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    INTERSECTION_TYPE_UNSPECIFIED: _ClassVar[IntersectionType]
    INTERSECTION_TYPE_TRAFFICLIGHT: _ClassVar[IntersectionType]
    INTERSECTION_TYPE_ROUNDABOUT: _ClassVar[IntersectionType]
    INTERSECTION_TYPE_STOP_SIGN: _ClassVar[IntersectionType]

class NodeType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PRIORITY: _ClassVar[NodeType]
    TRAFFIC_LIGHT: _ClassVar[NodeType]
INTERSECTION_TYPE_UNSPECIFIED: IntersectionType
INTERSECTION_TYPE_TRAFFICLIGHT: IntersectionType
INTERSECTION_TYPE_ROUNDABOUT: IntersectionType
INTERSECTION_TYPE_STOP_SIGN: IntersectionType
PRIORITY: NodeType
TRAFFIC_LIGHT: NodeType

class SimulationRequest(_message.Message):
    __slots__ = ("intersection_id", "simulation_parameters")
    INTERSECTION_ID_FIELD_NUMBER: _ClassVar[int]
    SIMULATION_PARAMETERS_FIELD_NUMBER: _ClassVar[int]
    intersection_id: str
    simulation_parameters: SimulationParameters
    def __init__(self, intersection_id: _Optional[str] = ..., simulation_parameters: _Optional[_Union[SimulationParameters, _Mapping]] = ...) -> None: ...

class SimulationParameters(_message.Message):
    __slots__ = ("intersection_type", "green", "yellow", "red", "speed", "seed")
    INTERSECTION_TYPE_FIELD_NUMBER: _ClassVar[int]
    GREEN_FIELD_NUMBER: _ClassVar[int]
    YELLOW_FIELD_NUMBER: _ClassVar[int]
    RED_FIELD_NUMBER: _ClassVar[int]
    SPEED_FIELD_NUMBER: _ClassVar[int]
    SEED_FIELD_NUMBER: _ClassVar[int]
    intersection_type: IntersectionType
    green: int
    yellow: int
    red: int
    speed: int
    seed: int
    def __init__(self, intersection_type: _Optional[_Union[IntersectionType, str]] = ..., green: _Optional[int] = ..., yellow: _Optional[int] = ..., red: _Optional[int] = ..., speed: _Optional[int] = ..., seed: _Optional[int] = ...) -> None: ...

class SimulationResultsResponse(_message.Message):
    __slots__ = ("total_vehicles", "average_travel_time", "total_travel_time", "average_speed", "average_waiting_time", "total_waiting_time", "generated_vehicles", "emergency_brakes", "emergency_stops", "near_collisions")
    TOTAL_VEHICLES_FIELD_NUMBER: _ClassVar[int]
    AVERAGE_TRAVEL_TIME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_TRAVEL_TIME_FIELD_NUMBER: _ClassVar[int]
    AVERAGE_SPEED_FIELD_NUMBER: _ClassVar[int]
    AVERAGE_WAITING_TIME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_WAITING_TIME_FIELD_NUMBER: _ClassVar[int]
    GENERATED_VEHICLES_FIELD_NUMBER: _ClassVar[int]
    EMERGENCY_BRAKES_FIELD_NUMBER: _ClassVar[int]
    EMERGENCY_STOPS_FIELD_NUMBER: _ClassVar[int]
    NEAR_COLLISIONS_FIELD_NUMBER: _ClassVar[int]
    total_vehicles: int
    average_travel_time: float
    total_travel_time: float
    average_speed: float
    average_waiting_time: float
    total_waiting_time: float
    generated_vehicles: int
    emergency_brakes: int
    emergency_stops: int
    near_collisions: int
    def __init__(self, total_vehicles: _Optional[int] = ..., average_travel_time: _Optional[float] = ..., total_travel_time: _Optional[float] = ..., average_speed: _Optional[float] = ..., average_waiting_time: _Optional[float] = ..., total_waiting_time: _Optional[float] = ..., generated_vehicles: _Optional[int] = ..., emergency_brakes: _Optional[int] = ..., emergency_stops: _Optional[int] = ..., near_collisions: _Optional[int] = ...) -> None: ...

class SimulationOutputResponse(_message.Message):
    __slots__ = ("intersection", "vehicles")
    INTERSECTION_FIELD_NUMBER: _ClassVar[int]
    VEHICLES_FIELD_NUMBER: _ClassVar[int]
    intersection: Intersection
    vehicles: _containers.RepeatedCompositeFieldContainer[Vehicle]
    def __init__(self, intersection: _Optional[_Union[Intersection, _Mapping]] = ..., vehicles: _Optional[_Iterable[_Union[Vehicle, _Mapping]]] = ...) -> None: ...

class Intersection(_message.Message):
    __slots__ = ("nodes", "edges", "connections", "traffic_lights")
    NODES_FIELD_NUMBER: _ClassVar[int]
    EDGES_FIELD_NUMBER: _ClassVar[int]
    CONNECTIONS_FIELD_NUMBER: _ClassVar[int]
    TRAFFIC_LIGHTS_FIELD_NUMBER: _ClassVar[int]
    nodes: _containers.RepeatedCompositeFieldContainer[Node]
    edges: _containers.RepeatedCompositeFieldContainer[Edge]
    connections: _containers.RepeatedCompositeFieldContainer[Connection]
    traffic_lights: _containers.RepeatedCompositeFieldContainer[TrafficLight]
    def __init__(self, nodes: _Optional[_Iterable[_Union[Node, _Mapping]]] = ..., edges: _Optional[_Iterable[_Union[Edge, _Mapping]]] = ..., connections: _Optional[_Iterable[_Union[Connection, _Mapping]]] = ..., traffic_lights: _Optional[_Iterable[_Union[TrafficLight, _Mapping]]] = ...) -> None: ...

class Node(_message.Message):
    __slots__ = ("id", "x", "y", "type")
    ID_FIELD_NUMBER: _ClassVar[int]
    X_FIELD_NUMBER: _ClassVar[int]
    Y_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    id: str
    x: float
    y: float
    type: NodeType
    def __init__(self, id: _Optional[str] = ..., x: _Optional[float] = ..., y: _Optional[float] = ..., type: _Optional[_Union[NodeType, str]] = ...) -> None: ...

class Edge(_message.Message):
    __slots__ = ("id", "to", "speed", "lanes")
    ID_FIELD_NUMBER: _ClassVar[int]
    FROM_FIELD_NUMBER: _ClassVar[int]
    TO_FIELD_NUMBER: _ClassVar[int]
    SPEED_FIELD_NUMBER: _ClassVar[int]
    LANES_FIELD_NUMBER: _ClassVar[int]
    id: str
    to: str
    speed: float
    lanes: int
    def __init__(self, id: _Optional[str] = ..., to: _Optional[str] = ..., speed: _Optional[float] = ..., lanes: _Optional[int] = ..., **kwargs) -> None: ...

class Connection(_message.Message):
    __slots__ = ("to", "fromLane", "toLane", "tl")
    FROM_FIELD_NUMBER: _ClassVar[int]
    TO_FIELD_NUMBER: _ClassVar[int]
    FROMLANE_FIELD_NUMBER: _ClassVar[int]
    TOLANE_FIELD_NUMBER: _ClassVar[int]
    TL_FIELD_NUMBER: _ClassVar[int]
    to: str
    fromLane: int
    toLane: int
    tl: int
    def __init__(self, to: _Optional[str] = ..., fromLane: _Optional[int] = ..., toLane: _Optional[int] = ..., tl: _Optional[int] = ..., **kwargs) -> None: ...

class TrafficLight(_message.Message):
    __slots__ = ("id", "type", "phases")
    ID_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PHASES_FIELD_NUMBER: _ClassVar[int]
    id: str
    type: str
    phases: _containers.RepeatedCompositeFieldContainer[Phase]
    def __init__(self, id: _Optional[str] = ..., type: _Optional[str] = ..., phases: _Optional[_Iterable[_Union[Phase, _Mapping]]] = ...) -> None: ...

class Phase(_message.Message):
    __slots__ = ("duration", "state")
    DURATION_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    duration: int
    state: str
    def __init__(self, duration: _Optional[int] = ..., state: _Optional[str] = ...) -> None: ...

class Vehicle(_message.Message):
    __slots__ = ("id", "positions")
    ID_FIELD_NUMBER: _ClassVar[int]
    POSITIONS_FIELD_NUMBER: _ClassVar[int]
    id: str
    positions: _containers.RepeatedCompositeFieldContainer[Position]
    def __init__(self, id: _Optional[str] = ..., positions: _Optional[_Iterable[_Union[Position, _Mapping]]] = ...) -> None: ...

class Position(_message.Message):
    __slots__ = ("time", "x", "y", "speed")
    TIME_FIELD_NUMBER: _ClassVar[int]
    X_FIELD_NUMBER: _ClassVar[int]
    Y_FIELD_NUMBER: _ClassVar[int]
    SPEED_FIELD_NUMBER: _ClassVar[int]
    time: int
    x: float
    y: float
    speed: float
    def __init__(self, time: _Optional[int] = ..., x: _Optional[float] = ..., y: _Optional[float] = ..., speed: _Optional[float] = ...) -> None: ...
