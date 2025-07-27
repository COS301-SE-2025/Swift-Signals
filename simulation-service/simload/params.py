import json
from pathlib import Path
from datetime import datetime
from simload.constants import (
    IntersectionType, TrafficDensity,
    DEFAULT_SEED, DEFAULT_SPEED
)


def _map_density(raw: int) -> TrafficDensity:
    try:
        return TrafficDensity(raw)
    except ValueError:
        return TrafficDensity.MEDIUM


def load_from_file(path: Path) -> tuple[dict, dict]:
    """Return (mapped_for_generator, raw_numeric)"""
    data = json.loads(path.read_text())
    return load_from_dict(data)


def load_from_dict(data: dict) -> tuple[dict, dict]:
    sim = data["intersection"]["simulation_parameters"]
    raw_density = data["intersection"].get("traffic_density", 1)

    raw_density = data["intersection"].get("traffic_density", 1)
    density = _map_density(raw_density)

    raw_type = int(sim.get("intersection_type", 0))
    intersection = IntersectionType(raw_type)

    mapped = {
        "traffic_density": density.name.lower(),
        "intersection_type": intersection.name.lower(),
        "speed": sim.get("speed", DEFAULT_SPEED),
        "seed": sim.get("seed", DEFAULT_SEED),
    }

    if intersection is IntersectionType.TRAFFIC_LIGHT:
        mapped.update({
            "green":   sim.get("green", 25),
            "yellow":  sim.get("yellow", 25),
            "red":     sim.get("red", 25),
        })

    raw = {
        "traffic_density": raw_density,
        "intersection_type": raw_type,
        "speed": sim.get("speed", DEFAULT_SPEED),
        "seed": sim.get("seed", DEFAULT_SEED),
    }

    return mapped, raw


def build_output_doc(
    sim_id: str,
    sim_name: str,
    owner: str,
    run_count: int,
    created_at: str,
    raw_params: dict,
    results: dict,
) -> dict:
    """Return the final Mongo‑style payload."""
    return {
        "_id": {"$oid": sim_id},
        "intersection": {
            "id": sim_id,
            "name": sim_name,
            "owner": owner,
            "created_at": created_at,
            "last_run_at": datetime.utcnow().isoformat() + "Z",
            "traffic_density": raw_params["traffic_density"],
            "status": 0,
            "run_count": run_count,
            "parameters": {
                "intersection_type": raw_params["intersection_type"]
            },
            "results": results,
        },
    }
