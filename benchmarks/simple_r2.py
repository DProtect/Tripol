#!/usr/bin/env python3
"""Compare two TRIPOL compilers on the same Linux programs and verify output parity."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import statistics
import subprocess
import tempfile
import time


CASES = (
    "examples/bench_compiler_revolution.tri",
    "examples/academy/simple-r1.tri",
    "examples/discord_new.tri",
)
ROOT = Path(__file__).resolve().parent.parent


def measure(command, env, input_data=None):
    start = time.perf_counter_ns()
    result = subprocess.run(
        command, env=env, input=input_data, capture_output=True, cwd=ROOT, check=True
    )
    return result.stdout, (time.perf_counter_ns() - start) / 1e6


def compare(path, compilers, binaries, samples, env):
    times = [[], []]
    for iteration in range(samples + 3):
        for i in ((0, 1) if iteration % 2 == 0 else (1, 0)):
            _, elapsed = measure(
                [compilers[i], "build", path, "--target", "linux", "--profile", "turbo",
                 "-o", binaries[i]], env
            )
            if iteration >= 3:
                times[i].append(elapsed)
    runtime = [[], []]
    expected = None
    for iteration in range(samples + 3):
        for i in ((0, 1) if iteration % 2 == 0 else (1, 0)):
            output, elapsed = measure([binaries[i]], env, b"7\n")
            if expected is None:
                expected = output
            if output != expected:
                raise SystemExit(f"Output differs for {path}")
            if iteration >= 3:
                runtime[i].append(elapsed)
    images = [binary.read_bytes() for binary in binaries]
    return {
        "build_ms": [round(statistics.median(v), 3) for v in times],
        "run_ms": [round(statistics.median(v), 3) for v in runtime],
        "bytes": [len(image) for image in images],
        "sha256": [hashlib.sha256(image).hexdigest() for image in images],
        "same_binary": images[0] == images[1],
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--before", required=True, type=Path)
    parser.add_argument("--after", required=True, type=Path)
    parser.add_argument("--runs", type=int, default=31)
    args = parser.parse_args()
    if args.runs < 1:
        parser.error("--runs must be positive")
    compilers = [args.before.resolve(), args.after.resolve()]
    env = dict(os.environ, TRIPOL_NO_UPDATE="1", GOMAXPROCS="2")
    with tempfile.TemporaryDirectory(prefix="tripol-simple-r2-") as folder:
        binaries = [Path(folder) / "before", Path(folder) / "after"]
        results = {path: compare(path, compilers, binaries, args.runs, env) for path in CASES}
    print(json.dumps({
        "columns": ["before", "after"], "runs": args.runs, "profile": "turbo",
        "GOMAXPROCS": 2, "cases": results,
    }, indent=2))


if __name__ == "__main__":
    main()
