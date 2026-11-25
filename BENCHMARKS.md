# TRIPOL v0.4 — Liberty measurements

Linux x86-64, Go 1.23.2. Compilers v0.3 and v0.4 built locally. Same fib(35) recursive source with header and result. One warmup, nine runs per version interleaved, stdout discarded, process startup included.

| Version | Median | Linux ELF | Windows PE |
|---|---:|---:|---:|
| 0.3 | 75.274 ms | 700 B | 2048 B |
| 0.4 | 75.370 ms | 700 B | 2048 B |

The generated Linux machine-code binaries are identical for this integer workload (verify with cmp). Timing difference is noise; no new speedup claimed.

User-reported Windows v0.3 five runs: 323.1765, 87.0602, 88.7783, 88.654, 92.7706 ms; median 88.7783 ms. v0.2 single reported run 23350.0547 ms, not a controlled comparison. Windows v0.4 runtime not executed here.

| Example | Linux bytes | Windows bytes |
|---|---:|---:|
| hello | 380 | 1536 |
| fibonacci(35) with header | 700 | 2048 |
| interactive calculator | 1615 | 3072 |

18 Go test groups, 109 nested subtests, all passed on Linux. PE imports inspected using objdump and debug/pe; Windows test suite cross-compiled, not run. No Windows runtime, memory or power measurements. No comparative Rust/C benchmarks.
