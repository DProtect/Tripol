"""Measure equivalent offline Discord serializers on Linux."""
import argparse
import json
import os
from pathlib import Path
import platform
from statistics import median
import time

ROOT = Path(__file__).resolve().parent


def memory(pid, field):
    for line in Path(f"/proc/{pid}/status").read_text().splitlines():
        if line.startswith(field + ":"):
            return int(line.split()[1])
    raise RuntimeError(f"{field} missing from Linux process status")


def projection(mode, data):
    if mode:
        container, = data["components"]
        item, = container["components"]
        return (data["flags"], container["type"], container["accent_color"], item["type"], item["content"],
                data["allowed_mentions"]["parse"])
    embed, rows = data
    row, = rows
    button, = row["components"]
    field, = embed["fields"]
    return (embed["title"], embed["description"], embed["color"], field["name"], field["value"], field["inline"],
            row["type"], button["type"], button["style"], button["label"], button["custom_id"], button.get("disabled", False))


def measure(command, count, mode):
    read, write = os.pipe()
    reply, emit = os.pipe()
    actions = [(os.POSIX_SPAWN_DUP2, read, 0), (os.POSIX_SPAWN_DUP2, emit, 1)]
    start = time.perf_counter_ns()
    pid = os.posix_spawnp(command[0], command, os.environ, file_actions=actions)
    os.close(read)
    os.close(emit)
    with os.fdopen(write, "w") as send, os.fdopen(reply) as receive:
        try:
            if receive.readline().strip() != "READY":
                raise RuntimeError(f"serializer did not start: {command}")
            ready = time.perf_counter_ns() - start
            idle = memory(pid, "VmRSS")
            send.write(f"{count}\n{mode}\n")
            send.flush()
            line = receive.readline()
            while line and not line.strip():
                line = receive.readline()
            elapsed, total = map(int, line.split())
            payload = json.loads(receive.readline())
            peak = memory(pid, "VmHWM")
            send.write("0\n")
            send.flush()
        finally:
            send.close()
            _, status, usage = os.wait4(pid, 0)
            if status != 0:
                raise RuntimeError(f"{command} exited with wait status {status}")
    want = (32768, 17, 0x00A0FF, 10, "TRIPOL", []) if mode else (
        "TRIPOL", "Discord", 0x00A0FF, "Version", "0.28", False, 1, 2, 3, "Confirm", "confirm", False)
    if projection(mode, payload) != want or total <= 0:
        raise ValueError(f"{command} returned a non-equivalent payload")
    return {
        "startup_ms": ready / 1e6,
        "idle_rss_kib": idle,
        "peak_rss_kib": peak,
        "cpu_ms_with_startup": (usage.ru_utime + usage.ru_stime) * 1000,
        "serialize_ms": elapsed / 1e6,
        "payloads_per_second": count * 1e9 / elapsed,
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--tripol", type=Path, default=ROOT / "tripol-serializer")
    parser.add_argument("--python", type=Path, default=ROOT / ".venv/bin/python")
    parser.add_argument("--rust", type=Path, default=ROOT / "serenity/target/release/tripol-discord-serializer")
    parser.add_argument("--node", default="node")
    parser.add_argument("--runs", type=int, default=7)
    parser.add_argument("--count", type=int, default=20000)
    args = parser.parse_args()
    if args.runs < 1 or args.count < 1:
        parser.error("--runs and --count must be positive")
    commands = {
        "TRIPOL": [str(args.tripol.absolute())],
        "Discord.js": [args.node, str(ROOT / "serialize.cjs")],
        "Discord.py": [str(args.python.absolute()), str(ROOT / "serialize.py")],
        "Serenity": [str(args.rust.absolute())],
    }
    results = {}
    for mode, case in enumerate(("embed_button", "components_v2")):
        available = list(commands) if mode == 0 else list(commands)[:-1]
        samples = {label: [] for label in available}
        for run in range(args.runs):
            rotate = run % len(available)
            for label in available[rotate:] + available[:rotate]:
                samples[label].append(measure(commands[label], args.count, mode))
        results[case] = {
            label: {key: round(median(v[key] for v in rows), 3) for key in rows[0]}
            for label, rows in samples.items()
        }
    print(json.dumps({"platform": platform.platform(), "runs": args.runs, "count": args.count, **results}, indent=2))


if __name__ == "__main__":
    main()
