# TRIPOL v0.6 — NEXUS measurements

Linux x86-64, Go compiler toolchain 1.23.x. All execution measurements were performed on Linux. Windows outputs were cross-compiled and inspected, not executed. User confirmed v0.5 on Windows using file and dynamic-array examples; no v0.5 Windows Fibonacci timing was supplied.

## Native array sum

Two programs in this version initialize an array of 65,535 integers (0 through 65,534), then sum it 500 times. One uses `sum(a)`, the other a checked, indexed `for` loop. Both print 1073692672500. Same compiler profile (nano). One warm-up each, nine alternating measured executions each, stdout discarded, process startup included. Timing from Python monotonic timer used only in development.

| Method | Median ms | Min ms | Max ms | ELF bytes |
|---|---:|---:|---:|---:|
| Built-in sum | 19.358 | 17.995 | 21.712 | 1107 |
| Indexed loop | 39.910 | 33.406 | 44.176 | 1176 |

This is a specialized built-in compared to a manually written loop, not an across-version or cross-language speed claim. Source: examples/bench_sum.tri and examples/bench_loop.tri.

## Size

- Windows hello (examples/hello.tri): 1536 bytes.
- Windows Fibonacci without heading (examples/fibonacci.tri): 1536 bytes.
- Windows interactive Nexus demo (examples/nexus.tri): 2560 bytes.
- Linux calculator template: 1650 bytes.

30 test groups / 174 subtests passed on Linux; `go vet .` passed. Native PE inspection and cross-build only for Windows. No new general allocation/reclamation or TIR has been implemented.
