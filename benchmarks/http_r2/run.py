"""Compare identical HTTP/1.1 reads over loopback on Windows and Linux."""
import argparse
import hashlib
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
import os
import platform
import subprocess
from pathlib import Path
import socket
from statistics import median
from tempfile import TemporaryDirectory
import threading
import time

import psutil

ROOT = Path(__file__).parent.resolve()
SMALL = b"t" * 512
LARGE = b"t" * (4 << 20)
PROFILES = {"sequential": (1, SMALL), "concurrent": (4, SMALL), "stream": (1, LARGE)}


class Server(ThreadingHTTPServer):
    def __init__(self):
        self.observed = {}
        self.lock = threading.Lock()
        super().__init__(("127.0.0.1", 0), Handler)


class Handler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def setup(self):
        self.request.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)
        super().setup()

    def do_GET(self):
        body = LARGE if self.path.startswith("/stream/") else SMALL
        with self.server.lock:
            count, sockets = self.server.observed.setdefault(self.path, [0, set()])
            self.server.observed[self.path][0] = count + 1
            sockets.add(self.client_address[1])
        self.send_response(200)
        self.send_header("Content-Type", "text/plain; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *_):
        pass


def resources(pid):
    root = psutil.Process(pid)
    resident, peak, cpu = 0, 0, 0.0
    for process in [root, *root.children(recursive=True)]:
        usage = process.memory_info()
        resident += usage.rss
        if os.name == "nt":
            peak += usage.peak_wset
        else:
            status = Path(f"/proc/{process.pid}/status").read_text()
            peak += 1024 * int(next(row.split()[1] for row in status.splitlines() if row.startswith("VmHWM:")))
        counters = process.cpu_times()
        cpu += counters.user + counters.system
    return resident // 1024, peak // 1024, cpu * 1000


def percentile(sorted_latencies, percent):
    index = (len(sorted_latencies) * percent + 99) // 100 - 1
    return sorted_latencies[index] / 1e6


def measure(command, url, count, mode, width, body, directory, legacy=False):
    started = time.perf_counter_ns()
    with subprocess.Popen(command, stdin=subprocess.PIPE, stdout=subprocess.PIPE, text=True) as worker:
        def next_output():
            for line in worker.stdout:
                if line.strip():
                    return line.strip()
            raise RuntimeError(f"HTTP worker closed its output: {command}")

        def send(text):
            worker.stdin.write(text)
            worker.stdin.flush()

        try:
            configuration = f"{url}\n{count}\n"
            if not legacy:
                send(configuration + f"{mode}\n{width}\n{len(body)}\n{directory}\n")
            if next_output() != "READY":
                raise RuntimeError(f"HTTP worker did not start: {command}")
            startup = (time.perf_counter_ns() - started) / 1e6
            idle, _, _ = resources(worker.pid)
            send(configuration if legacy else "1\n")
            elapsed, total = map(int, next_output().split())
            if total != count * len(body):
                raise RuntimeError(f"unexpected transferred byte count: {total}")
            latencies = []
            if not legacy:
                while (line := next_output()) != "END":
                    latencies.extend(map(int, line.split()))
                if len(latencies) != count // width:
                    raise RuntimeError(f"unexpected number of latency samples: {len(latencies)}")
                latencies.sort()
            _, peak, cpu = resources(worker.pid)
            send("0\n")
            worker.stdin.close()
            if worker.wait(timeout=10):
                raise RuntimeError(f"HTTP worker failed: {command}")
        finally:
            if worker.poll() is None:
                worker.kill()
                worker.wait(timeout=10)
    if mode == "stream":
        expected = hashlib.sha256(body).digest()
        files = sorted(Path(directory).glob("*.bin"))
        if len(files) != count:
            raise RuntimeError(f"expected {count} downloaded files, found {len(files)}")
        for file in files:
            with file.open("rb") as data:
                if hashlib.file_digest(data, "sha256").digest() != expected:
                    raise RuntimeError(f"corrupt downloaded file: {file}")
    summary = {
        "startup_ms": startup,
        "idle_rss_kib": idle,
        "peak_rss_kib": peak,
        "cpu_ms_with_startup": cpu,
        "batch_ms_per_request": elapsed / count / 1e6,
        "requests_per_second": count * 1e9 / elapsed,
        "bytes_per_second": total * 1e9 / elapsed,
    }
    if latencies:
        summary.update({f"p{p}_ms": percentile(latencies, p) for p in (50, 95, 99)})
    return summary


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--tripol", default=str(ROOT / "tripol-http"))
    parser.add_argument("--baseline", help="Optional pre-optimization TRIPOL executable")
    parser.add_argument("--reference", help="Reference TRIPOL worker using the current measurement protocol")
    parser.add_argument("--go", default=str(ROOT / "go-http"))
    parser.add_argument("--rust", default=str(ROOT / "reqwest/target/release/tripol-http-fixture"))
    parser.add_argument("--python", default=str(ROOT / ".venv/bin/python"))
    parser.add_argument("--runs", type=int, default=7)
    parser.add_argument("--count", type=int, help="Override the number of requests in each sample")
    parser.add_argument("--mode", choices=PROFILES, default="sequential")
    parser.add_argument("--client", action="append", help="Select a client, repeat to select multiple")
    parser.add_argument("--launcher", help="TRIPOL CLI for launching generated Windows PE workers")
    args = parser.parse_args()
    width, body = PROFILES[args.mode]
    count = args.count if args.count is not None else (12 if args.mode == "stream" else 400)
    if args.runs < 1 or count < 1 or count % width:
        parser.error("runs and count must be positive and count divisible by the concurrency width")
    if args.baseline and args.mode != "sequential":
        parser.error("the previous executable supports only the sequential profile")
    commands = {
        "TRIPOL": [args.launcher, "net-run", args.tripol] if args.launcher else [args.tripol],
        "Requests": [args.python, str(ROOT / "python_clients.py"), "requests"],
        "HTTPX": [args.python, str(ROOT / "python_clients.py"), "httpx"],
        "Fetch": ["node", str(ROOT / "fetch.cjs")],
        "Reqwest": [args.rust],
        "Go": [args.go],
    }
    if args.baseline:
        commands["TRIPOL before"] = [args.baseline]
    if args.reference:
        commands["TRIPOL reference"] = [args.launcher, "net-run", args.reference] if args.launcher else [args.reference]
    if args.client:
        unknown = set(args.client) - commands.keys()
        if unknown:
            parser.error(f"unknown clients: {', '.join(sorted(unknown))}")
        commands = {name: command for name, command in commands.items() if name in args.client}
    samples = {name: [] for name in commands}
    with Server() as server:
        worker = threading.Thread(target=server.serve_forever, daemon=True)
        worker.start()
        try:
            names = list(commands)
            for run in range(args.runs):
                rotate = run % len(names)
                for name in names[rotate:] + names[:rotate]:
                    route = f"/{args.mode}/{name.replace(' ', '-')}/{run}"
                    url = f"http://127.0.0.1:{server.server_port}{route}"
                    with TemporaryDirectory(prefix="tripol-http-bench-") as directory:
                        result = measure(commands[name], url, count, args.mode, width, body,
                                         directory, legacy=name == "TRIPOL before")
                    with server.lock:
                        observed, sockets = server.observed[route]
                    if observed != count:
                        raise RuntimeError(f"{name} reached the server {observed} times")
                    result["connections"] = len(sockets)
                    samples[name].append(result)
        finally:
            server.shutdown()
            worker.join()
    output = {
        name: {key: round(median(row[key] for row in runs), 3) for key in runs[0]}
        for name, runs in samples.items()
    }
    print(json.dumps({"platform": platform.platform(), "mode": args.mode, "runs": args.runs, "count": count, "width": width,
                      "bytes": len(body), "results": output}, indent=2))


if __name__ == "__main__":
    main()
