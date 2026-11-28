"""Compare HTTP/1.1 byte-identical sequential reads on a loopback server."""
import argparse
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
import os
from pathlib import Path
import socket
from statistics import median
import threading
import time

ROOT = Path(__file__).parent.resolve()
BODY = b"t" * 512


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
        with self.server.lock:
            count, sockets = self.server.observed.setdefault(self.path, [0, set()])
            self.server.observed[self.path][0] = count + 1
            sockets.add(self.client_address[1])
        self.send_response(200)
        self.send_header("Content-Type", "text/plain; charset=utf-8")
        self.send_header("Content-Length", str(len(BODY)))
        self.end_headers()
        self.wfile.write(BODY)

    def log_message(self, *_):
        pass


def memory(pid, field):
    process = Path(f"/proc/{pid}")
    line = next(x for x in (process / "status").read_text().splitlines() if x.startswith(field + ":"))
    children = []
    for entry in Path("/proc").iterdir():
        if not entry.name.isdecimal():
            continue
        try:
            parent = int((entry / "stat").read_text().rsplit(")", 1)[1].split()[1])
        except (FileNotFoundError, ProcessLookupError, PermissionError):
            continue
        if parent == pid:
            children.append(int(entry.name))
    return int(line.split()[1]) + sum(memory(child, field) for child in children)


def measure(command, url, count):
    read, send = os.pipe()
    receive, write = os.pipe()
    ops = [(os.POSIX_SPAWN_DUP2, read, 0), (os.POSIX_SPAWN_DUP2, write, 1)]
    started = time.perf_counter_ns()
    pid = os.posix_spawnp(command[0], command, os.environ, file_actions=ops)
    os.close(read)
    os.close(write)
    with os.fdopen(send, "w") as stdin, os.fdopen(receive) as stdout:
        try:
            if stdout.readline().strip() != "READY":
                raise RuntimeError(f"HTTP worker did not start: {command}")
            startup = (time.perf_counter_ns() - started) / 1e6
            idle = memory(pid, "VmRSS")
            stdin.write(f"{url}\n{count}\n")
            stdin.flush()
            line = stdout.readline()
            while line and not line.strip():
                line = stdout.readline()
            elapsed, total = map(int, line.split())
            if total != count * len(BODY):
                raise RuntimeError(f"unexpected transferred byte count: {total}")
            peak = memory(pid, "VmHWM")
            stdin.write("0\n")
            stdin.flush()
        finally:
            stdin.close()
            _, status, cpu = os.wait4(pid, 0)
            if status != 0:
                raise RuntimeError(f"{command} exited with status {status}")
    return {
        "startup_ms": startup,
        "idle_rss_kib": idle,
        "peak_rss_kib": peak,
        "cpu_ms_with_startup": (cpu.ru_utime + cpu.ru_stime) * 1000,
        "mean_request_ms": elapsed / count / 1e6,
        "requests_per_second": count * 1e9 / elapsed,
        "bytes_per_second": total * 1e9 / elapsed,
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--tripol", default=str(ROOT / "tripol-http"))
    parser.add_argument("--baseline", help="Optional pre-optimization TRIPOL executable")
    parser.add_argument("--go", default=str(ROOT / "go-http"))
    parser.add_argument("--rust", default=str(ROOT / "reqwest/target/release/tripol-http-fixture"))
    parser.add_argument("--python", default=str(ROOT / ".venv/bin/python"))
    parser.add_argument("--runs", type=int, default=7)
    parser.add_argument("--count", type=int, default=400)
    args = parser.parse_args()
    if args.runs < 1 or args.count < 1:
        parser.error("runs and count must be positive")
    commands = {
        "TRIPOL": [args.tripol],
        "Requests": [args.python, str(ROOT / "python_clients.py"), "requests"],
        "HTTPX": [args.python, str(ROOT / "python_clients.py"), "httpx"],
        "Fetch": ["node", str(ROOT / "fetch.cjs")],
        "Reqwest": [args.rust],
        "Go": [args.go],
    }
    if args.baseline:
        commands["TRIPOL before"] = [args.baseline]
    samples = {name: [] for name in commands}
    with Server() as server:
        worker = threading.Thread(target=server.serve_forever, daemon=True)
        worker.start()
        try:
            names = list(commands)
            for run in range(args.runs):
                rotate = run % len(names)
                for name in names[rotate:] + names[:rotate]:
                    route = f"/{name.replace(' ', '-')}/{run}"
                    result = measure(commands[name], f"http://127.0.0.1:{server.server_port}{route}", args.count)
                    with server.lock:
                        count, sockets = server.observed[route]
                    if count != args.count:
                        raise RuntimeError(f"{name} reached the server {count} times")
                    result["connections"] = len(sockets)
                    samples[name].append(result)
        finally:
            server.shutdown()
            worker.join()
    output = {
        name: {key: round(median(row[key] for row in runs), 3) for key in runs[0]}
        for name, runs in samples.items()
    }
    print(json.dumps({"runs": args.runs, "count": args.count, "bytes": len(BODY), "results": output}, indent=2))


if __name__ == "__main__":
    main()
