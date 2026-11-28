from concurrent.futures import ThreadPoolExecutor
import os
from pathlib import Path
import sys
import time

REQUESTS = sys.argv[1] == "requests"
if REQUESTS:
    import requests
else:
    import httpx


def receive(client, url, size, destination=None):
    if destination is None:
        response = client.get(url, timeout=10)
        status, length = response.status_code, len(response.text)
    else:
        response = client.get(url, stream=True, timeout=10) if REQUESTS else client.stream("GET", url, timeout=10)
        with response as stream, destination.open("xb") as output:
            status, length = stream.status_code, 0
            chunks = stream.iter_content(32768) if REQUESTS else stream.iter_bytes(32768)
            for block in chunks:
                length += output.write(block)
            output.flush()
            os.fsync(output.fileno())
    if status != 200 or length != size:
        raise RuntimeError("unexpected response")
    return length


def main():
    url, count, mode = input(), int(input()), input()
    width, size, directory = int(input()), int(input()), Path(input())
    if REQUESTS:
        # Session state is mutable. Give each concurrent lane its own session.
        clients = [requests.Session() for _ in range(width)]
    else:
        client = httpx.Client(limits=httpx.Limits(max_connections=4, max_keepalive_connections=4))
        clients = [client] * width
    try:
        print("READY", flush=True)
        input()
        start = time.perf_counter_ns()
        total = 0
        if mode == "concurrent":
            with ThreadPoolExecutor(max_workers=width) as executor:
                for _ in range(0, count, width):
                    batch = [executor.submit(receive, client, url, size) for client in clients]
                    total += sum(job.result() for job in batch)
        else:
            for index in range(count):
                destination = directory / f"{index}.bin" if mode == "stream" else None
                total += receive(clients[0], url, size, destination)
        print(time.perf_counter_ns() - start, total, flush=True)
        input()
    finally:
        for client in set(clients):
            client.close()


if __name__ == "__main__":
    main()
