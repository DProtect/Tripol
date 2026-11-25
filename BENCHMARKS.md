# TRIPOL v0.13 — INFINITY verification

Linux x86-64, Go 1.23.2. 89 top-level test groups and 285 subtests passed; `go vet .` passed. Windows compiler and installer cross-compiled; native Windows test examples PE-inspected, not executed here. Previous v0.12 modal form was demonstrated by the project owner on Windows. No third-party language speed claims.

| Workload | v0.12 | v0.13 | Notes |
|---|---:|---:|---|
| `examples/bench_literal_len.tri` | 36.487 ms | 0.463 ms | 11 alternating runs per binary; identical output 819200000; process startup included |
| Linux bytes (same workload) | 4776 | 595 | Literal length resolved at compile time |
| `examples/bench_inline.tri` | 6.727 ms | 6.866 ms | 21 alternating runs; identical output 6567000000; no stable improvement |
| Linux bytes (inline workload) | 684 | 684 | Both nano |
| Hello world Windows/Linux | 1536/373 B | 1536/373 B | All 30 unused imports produce identical output bytes |
| Infinity demo Windows | — | 6656 B | New syntax and five libraries |

Use the same source with the two compilers, `tripol build FILE --target linux -o PATH`, then run the binaries repeatedly interleaved. Wall-clock medians include process startup. The literal test isolates a compile-time optimization; do not extrapolate to general string processing. TIR inlining preserves checked arithmetic and argument evaluation, and is disabled under trace/debug. General memory and platform limitations remain.

## Previous release verification

# TRIPOL 0.12 — SENTINEL verification

Linux x86-64, Go 1.23.2. `go test -count=1 -v .` and `go vet .` pass; 80 top-level tests and 274 subtests. Windows compiler and installer cross-built; generated PE binaries inspected with debug/pe and objdump. New Windows UI and HTTPS execution not tested. No cross-language performance claim.

| Example | Linux bytes | Windows bytes |
|---|---:|---:|
| Hello | 373 | 1536 |
| Repository Fibonacci | 593 | 1536 |
| Sentinel (try, JSON, scratch) | 8696 | 9728 |
| Form | unsupported at runtime | 6144 |
| HTTPS | unsupported at runtime | 6656 |

SSE2 strlen comparison: generate `bench_strlen.tri` with `python examples/bench_strlen.py`, build identical source with v0.11 and v0.12 Linux nano, verify output `655360000`, warm up once and measure seven interleaved runs including process startup with stdout discarded. On this host: v0.11 median 218.091972 ms, v0.12 26.885522 ms (8.11x faster for this workload); executables 33396 and 33456 bytes. Long-string workload only; no general speedup claimed. The user measured v0.10 Fibonacci median 88.051 ms and ecosystem compilation 3.708 ms in Windows; no comparable v0.12 Windows figures yet.

Scratch test: 2000 arrays of 65535 int64 elements (1,048,576,000 cumulative allocated bytes including length headers) reused within the 8 MiB arena; output 4000. Cumulative allocation is not peak RSS. Escape rejection and recovery cleanup verified. All prior regression tests remain in the Go suite. See test-results.txt for exact logs.
