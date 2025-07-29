from contextlib import contextmanager
from pathlib import Path

RUN_COUNTER_FILE = Path("run_count.txt")


@contextmanager
def run_counter():
    """Context manager that yields the next run number and updates the counter."""
    count = 0
    if RUN_COUNTER_FILE.exists():
        count = int(RUN_COUNTER_FILE.read_text().strip())
    count += 1
    try:
        yield count
    finally:
        RUN_COUNTER_FILE.write_text(str(count))


def ensure_dir(path: Path) -> None:
    path.mkdir(parents=True, exist_ok=True)
