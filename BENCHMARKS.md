# TRIPOL v0.7 — HORIZON measurements

Linux x86-64, Go 1.23.2. Same 3,000 distinct string literal print statements compiled with v0.6 and v0.7; one warm-up per version, 11 timed compilations each alternating order, includes process startup and disk write. Both produced byte-for-byte identical 95,235-byte Linux executables and identical stdout.

| Compiler | Median | Min | Max |
|---|---:|---:|---:|
| v0.6 | 264.935 ms | 254.472 ms | 284.503 ms |
| v0.7 | 24.650 ms | 20.491 ms | 37.622 ms |

About 10.75x faster on this particular literal-heavy compilation. Not a general speed comparison or application runtime speedup. Reproduce with `python examples/bench_literals.py /tmp/literals.tri` (maintainer-only generator), then `tripol build /tmp/literals.tri -o /tmp/literals` with both versions. Timing depends on hardware and system load.

The experimental change to for-loop code generation was reverted after measurements did not show a reliable benefit. No general application runtime speedup is claimed. No comparative benchmark with Rust, C or Python has been run.

Windows PE sizes (cross-built and inspected, execution pending): hello 1,536 B, interactive record example 2,560 B. Tests execute in Linux. User-confirmed v0.6 Windows calculator template and offline tutorial.
