from enum import IntEnum

class IntersectionType(IntEnum):
    UNSPECIFIED = 0
    TRAFFIC_LIGHT = 1
    ROUNDABOUT = 2
    FOUR_WAY_STOP = 3
    T_JUNCTION = 4

class TrafficDensity(IntEnum):
    LOW = 0
    MEDIUM = 1
    HIGH = 2

STATUS_MAP = {0: "unoptimized", 1: "optimizing", 2: "optimized", 3: "failed"}
DEFAULT_SPEED = 40
DEFAULT_SEED = 42
