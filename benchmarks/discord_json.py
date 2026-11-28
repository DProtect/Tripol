"""Compare native JSON lookup with two TRIPOL compilers on Linux x86-64."""
import argparse
import json
import os
from pathlib import Path
from statistics import median
import subprocess
import tempfile
import time


def measure(program, env):
    before = time.perf_counter_ns()
    with tempfile.TemporaryFile() as stdout:
        pid = os.posix_spawn(
            str(program), [str(program)], env,
            file_actions=[(os.POSIX_SPAWN_DUP2, stdout.fileno(), 1)],
        )
        _, status, usage = os.wait4(pid, 0)
        elapsed = (time.perf_counter_ns() - before) / 1e9
        stdout.seek(0)
        output = stdout.read()
    if status != 0 or output != b"2100000\n":
        raise RuntimeError(f"{program.name} produced {output!r} with wait status {status}")
    return elapsed, usage.ru_utime + usage.ru_stime


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("before", type=Path, help="TRIPOL executable built from the baseline")
    parser.add_argument("after", type=Path, help="TRIPOL executable built from the candidate")
    parser.add_argument("--runs", type=int, default=11)
    args = parser.parse_args()
    if args.runs < 1:
        parser.error("--runs must be positive")
    env = {**os.environ, "TRIPOL_NO_UPDATE": "1"}
    event = {
        "op": 0,
        "t": "INTERACTION_CREATE",
        "s": 42,
        "d": {
            "id": 42,
            "type": 2,
            "data": {
                "type": 1,
                "name": "rating",
                "options": [{"name": "notes", "type": 3, "value": "preview " * 100}],
            },
        },
    }
    payload = json.dumps(json.dumps(event))
    source = f"""import tripol.json
payload={payload}
total=0
repeat 50000 {{scratch {{total += parse_int(json.get_raw_large(payload,"s"))}}}}
show total
"""
    with tempfile.TemporaryDirectory() as directory:
        root = Path(directory)
        test = root / "lookup.tri"
        test.write_text(source, encoding="utf-8")
        programs = {}
        results = {"before": [], "after": []}
        for label, compiler in (("before", args.before), ("after", args.after)):
            program = root / label
            subprocess.run(
                [str(compiler.resolve()), "build", str(test), "--target", "linux",
                 "--profile", "turbo", "-o", str(program)],
                env=env, check=True, stdout=subprocess.DEVNULL,
            )
            programs[label] = program
        for run in range(args.runs):
            order = ("before", "after") if run % 2 == 0 else ("after", "before")
            for label in order:
                results[label].append(measure(programs[label], env))
        report = {
            label: {
                "wall_s": round(median(wall for wall, _ in samples), 4),
                "cpu_s": round(median(cpu for _, cpu in samples), 4),
                "elf_bytes": programs[label].stat().st_size,
            }
            for label, samples in results.items()
        }
        print(json.dumps({"runs": args.runs, **report}, indent=2))


if __name__ == "__main__":
    main()
