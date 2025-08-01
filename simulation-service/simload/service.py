import uuid
from datetime import datetime
from simload import generator, params as param_mod
from simload.constants import IntersectionType

_RUN_COUNTER = 0


def run_simulation(input_data: dict) -> dict:
    global _RUN_COUNTER
    _RUN_COUNTER += 1

    mapped, raw = param_mod.load_from_dict(input_data)

    intersection_type = IntersectionType(raw["intersection_type"])
    sim_name = intersection_type.name.title()
    sim_id = uuid.uuid4().hex[:24]
    created_at = datetime.utcnow().isoformat() + "Z"

    results, full_out = generator.run(intersection_type, mapped)

    response = param_mod.build_output_doc(
        sim_id=sim_id,
        sim_name=sim_name,
        owner=input_data.get("owner", "grpc-user"),
        run_count=_RUN_COUNTER,
        created_at=created_at,
        raw_params=raw,
        results=results,
    )

    return {
        "summary": response,
        "full_output": full_out,
    }
