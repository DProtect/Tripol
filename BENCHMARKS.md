# TRIPOL v0.9 — ECOSYSTEM measurements

Linux x86-64, Go 1.23.2. 57 top-level test groups, 228 subtests passed; `go vet .` passed. Windows binaries cross-built and PE-inspected, not executed on Windows here. User validated v0.8 in Windows.

## Text prefix/suffix benchmark

`examples/bench_text.tri`, 200,000 iterations of both prefix and suffix checks. Identical source, native Linux nano profile, stdout discarded for timing, process startup included, one warm-up per version, 11 runs per version with alternating order. Both outputs verified as 200000.

| Version | Median ms | Min ms | Max ms | Linux bytes |
|---|---:|---:|---:|---:|
| v0.8 | 16.198 | 15.290 | 24.635 | 1262 |
| v0.9 | 3.656 | 3.167 | 4.351 | 815 |

Observed speedup 4.43x; old prefix/suffix allocate substrings, new native helpers compare bytes without allocating. A separate test performs one million pairs of checks successfully with no temporary strings. No claim that all programs are faster. No cross-language comparison in this release.

## Binary sizes

| Program | Linux bytes | Windows bytes |
|---|---:|---:|
| Hello | 373 | 1536 |
| Same hello importing all 15 standard modules | identical | identical |
| Ecosystem example | 3502 | 4608 |
| Atlas example (original v0.8 source) | not measured | 3584 |
| Vault sample | not measured | 2048 |

`tripol bench-build examples/ecosystem.tri --runs 10`: median 1.654 ms, min 1.211 ms, max 4.300 ms; deterministic 4608-byte Windows output. Includes source loading and native generation; excludes writing and compiler startup. PE inspection confirms GetTickCount64 and Sleep imports when needed, along with prior file/console/memory APIs. See test-results.txt.
