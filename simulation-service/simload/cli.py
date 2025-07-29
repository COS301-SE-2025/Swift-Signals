import json
import uuid
import logging
from pathlib import Path
from datetime import datetime

import argparse
from simload import generator, io_utils, params as param_mod
from simload.constants import IntersectionType

LOG = logging.getLogger(__name__)


def parse_cli(argv: list[str] | None = None) -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Run a SUMO intersection simulation.")
    p.add_argument(
        "--params", type=Path, required=True, help="Path to parameter JSON file"
    )
    p.add_argument("--owner", default="username")
    p.add_argument(
        "--outdir",
        type=Path,
        default=Path("out"),
        help="Directory where output files are stored",
    )
    return p.parse_args(argv)


def main(argv: list[str] | None = None) -> None:
    args = parse_cli(argv)
    mapped, raw = param_mod.load_from_file(args.params)

    intersection_type = IntersectionType(raw["intersection_type"])
    sim_name = intersection_type.name.title()

    with io_utils.run_counter() as run_count:
        results, full_out = generator.run(intersection_type, mapped)

        sim_id = uuid.uuid4().hex[:24]
        created = datetime.utcnow().isoformat() + "Z"

        output_doc = param_mod.build_output_doc(
            sim_id, sim_name, args.owner, run_count, created, raw, results
        )

        # ---------- write files ----------
        results_dir = args.outdir / "results"
        simout_dir = args.outdir / "simulationOut"
        io_utils.ensure_dir(results_dir)
        io_utils.ensure_dir(simout_dir)

        (results_dir / "simulation_results.json").write_text(
            json.dumps(output_doc, indent=2)
        )
        (simout_dir / "simulation_output.json").write_text(
            json.dumps(full_out, indent=2)
        )

        print(full_out)

        LOG.info("Saved artefacts to %s", args.outdir.resolve())


if __name__ == "__main__":
    main()
