from intersections import circle, stopStreet, tJunction, trafficLight
from simload.constants import IntersectionType

_GENERATOR_DISPATCH = {
    IntersectionType.TRAFFIC_LIGHT: trafficLight.generate,
    IntersectionType.ROUNDABOUT: circle.generate,
    IntersectionType.FOUR_WAY_STOP: stopStreet.generate,
    IntersectionType.T_JUNCTION: tJunction.generate,
}


def run(intersection_type: IntersectionType, params: dict):
    try:
        gen_fn = _GENERATOR_DISPATCH[intersection_type]
    except KeyError:
        raise ValueError(f"No generator registered for {intersection_type}")
    return gen_fn(params)
