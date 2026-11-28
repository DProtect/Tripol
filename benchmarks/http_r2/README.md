# DIPO HTTP 2 — benchmark HTTP local

## Reproducción / Reproduction

Ejecutado en Linux x86-64, AMD EPYC 9V74 con tres CPU visibles. TRIPOL v0.29 y Go usan Go 1.23.12. Workers externos con CPython 3.12.14, Requests 2.32.5, HTTPX 0.28.1, Fetch de Node 24.22.3 y Reqwest 0.12.28 con Rust 1.85.0. El servidor local está incluido.

Run the following in this directory. Supply the DIPO HTTP 2 compiler from this branch and Python 3.12, Node and Rust. The harness controls the server and validates the body length and connection count.

```sh
python3 -m venv .venv
.venv/bin/python -m pip install -r requirements.txt
cargo build --manifest-path reqwest/Cargo.toml --release --locked
tripol build tripol.tri --target linux --profile turbo -o tripol-http
go build -trimpath -ldflags="-s -w" -o go-http net_http.go
python3 run.py --count 400 --runs 7
```

Cada muestra hace 400 GET secuenciales por HTTP/1.1 al mismo servidor loopback. Cada respuesta contiene 512 bytes de texto. El arnés alterna el orden de los clientes. Las cifras son medianas de siete muestras. La latencia a continuación es la media por petición de cada muestra, incluida la primera. RSS de TRIPOL incluye el lanzador y el programa generado. El pico suma los máximos de RSS de los procesos vivos al final de la muestra. La CPU incluye el arranque.

The table reports medians of seven samples of 400 sequential HTTP/1.1 GETs with identical 512-byte response bodies. TRIPOL memory includes the launcher and the generated program. The reported latency is each run's mean per request. The script also reports startup, idle RSS, transfer throughput and socket counts.

| Cliente | ms por GET | GET/s | RSS pico KiB | CPU ms | Conexiones |
| --- | ---: | ---: | ---: | ---: | ---: |
| TRIPOL | 0.333 | 3,001 | 13,420 | 122.486 | 1 |
| Requests | 0.547 | 1,830 | 21,168 | 304.175 | 1 |
| HTTPX | 0.444 | 2,253 | 23,040 | 307.948 | 1 |
| Fetch | 0.737 | 1,357 | 74,384 | 418.826 | 2 |
| Reqwest | 0.104 | 9,630 | 3,396 | 22.964 | 1 |
| Go net/http | 0.144 | 6,958 | 8,360 | 44.743 | 1 |
| TRIPOL antes de optimizar la lectura del resultado | 0.698 | 1,433 | 13,292 | 261.358 | 1 |

La optimización evita volver a recorrer el mensaje JSON del lanzador para cada campo de respuesta. Para repetir la comparación con la implementación anterior, compila el worker con un compilador que todavía usa búsquedas independientes en `parse_response` y pásalo a `--baseline`. [Resultados completos](results.json)

The optimization parses the response envelope once. Use `--baseline` with a worker built using the earlier per-field lookup implementation. Concurrent throughput, sustained streaming, Windows performance and live Discord round trips require separate workloads. [Full results](results.json)
