import sys
import time

if sys.argv[1] == "requests":
    import requests
    client = requests.Session()
else:
    import httpx
    client = httpx.Client()

print("READY", flush=True)
url, count = input(), int(input())
start = time.perf_counter_ns()
total = 0
with client:
    for _ in range(count):
        response = client.get(url, timeout=10)
        if response.status_code != 200 or len(response.text) != 512:
            raise RuntimeError("unexpected response")
        total += len(response.text)
    print(time.perf_counter_ns() - start, total, flush=True)
    input()
