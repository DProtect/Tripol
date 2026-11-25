# TRIPOL v0.10 — CONNECT measurements

Linux x86-64, Go 1.23.2. 64 test groups and 234 subtests passed; `go vet .` passed. Native loopback TCP tests executed on Linux. Windows binaries cross-built and PE-inspected, not executed on Windows here. User reported v0.9 working in Windows (83.144 ms median Fibonacci; 4.860 ms ecosystem compilation; 4608-byte ecosystem executable).

## Native substring search vs manual byte-based implementation

`examples/bench_find_manual.tri` built by v0.9, and `examples/bench_find.tri` built by v0.10. Both search for `xyz0` in the same 36-byte text 200,000 times and return 200000. The manual implementation uses checked byte_at operations; the new intrinsic scans bytes directly and avoids repeated length scans. Neither allocates temporary strings. One warm-up and nine alternating measured executions each, stdout discarded, process startup included.

| Program | Median ms | Min ms | Max ms | Linux bytes |
|---|---:|---:|---:|---:|
| v0.9 manual substring | 120.250 | 116.211 | 122.866 | 1144 |
| v0.10 text.contains | 5.515 | 4.605 | 5.807 | 755 |

21.8x observed speedup for this specific implementation comparison, not a general compiler speedup or cross-language comparison. One million native searches also passed without temporary allocations. The original Fibonacci example compiled to identical 593-byte Linux executables on v0.9 and v0.10, so no speedup claimed.

## Binary sizes

| Program | Linux bytes | Windows bytes |
|---|---:|---:|
| Hello world | 373 | 1536 |
| Original Fibonacci example (no heading) | 593 | 1536 |
| HTTP client demo | 3896 | 5632 |
| Ecosystem | not measured | 4608 |

`tripol bench-build examples/ecosystem.tri --runs 10`: median 1.628 ms, min 1.114 ms, max 2.370 ms; deterministic 4608-byte Windows output, excludes compiler startup and executable write. Native TCP tests cover partial reads, zero-length response, 65535-byte limit, NUL, refused connection, bad types and ports. The demo server and generated Linux client yielded status 200 and expected body. PE inspection confirms KERNEL32.dll and WS2_32.dll imports, unused networking does not alter nano binaries. No TLS/UDP/DNS/IPv6 or native server yet. See README.md.
