# TRIPOL v0.5 — VELOCITY measurements

Linux x86-64, Go 1.23.2. Same recursive Fibonacci(35) with heading and result. Both versions compiled to native code locally; one warmup, 11 timed executions each, alternating order, stdout discarded, process startup included. Correct output 9227465 was verified for both versions.

| | v0.4 | v0.5 |
|---|---:|---:|
| Fibonacci Linux binary | 700 B | 677 B |
| Fibonacci Linux median | 74.508 ms | 75.295 ms |
| Fibonacci Linux min | 71.545 ms | 74.711 ms |
| Fibonacci Linux max | 79.975 ms | 86.274 ms |
| Fibonacci Windows binary | 2048 B | 2048 B |
| Hello Linux | 380 B | 373 B |
| Hello Windows | 1536 B | 1536 B |

No demonstrated Fibonacci speed improvement. Linux code size decreased, while Windows PE alignment preserved the size. Windows sizes measured on cross-generated files, not runtime performance. v0.5 files example: Windows 3072 B, Linux 1846 B. v0.5 runtime array example: Windows 3072 B. The owner reported v0.4 Windows Fibonacci median 85.961 ms (10 runs, one warmup, 84.631–89.684 ms); this is a separate baseline, not directly comparable to Linux. No comparative Rust/C/Go/Python benchmark yet.
