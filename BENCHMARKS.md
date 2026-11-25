## ORIGIN v0.23

Linux x86-64, 21 interleaved executions, Turbo, one million integer search (same source/output -1): NEXUS median 0.158765 ms; ORIGIN 0.115130 ms. 1.38x throughput improvement in this specific workload, not a general Rust comparison. See `benchmarks/origin.py` and `benchmarks/origin-linux.json`.

# TRIPOL v0.20 — OMEGA verification

Linux x86-64, 21 interleaved Turbo samples per version, same source, same runtime input, all outputs verified. Full CLI build includes startup and output write; internal time excludes these. Rust not installed.

| Workload | APEX internal ms | OMEGA internal ms | Speedup | Build ms old/new | ELF bytes old/new |
|---|---:|---:|---:|---:|---:|
| division | 9.996 | 10.149 | 0.98x | 2.134 / 2.122 | 1586 / 1586 |
| registers | 7.982 | 7.971 | 1.00x | 2.149 / 2.217 | 1679 / 1679 |
| fibonacci | 3.951 | 3.845 | 1.03x | 1.848 / 1.912 | 1491 / 1491 |
| arrays | 1.452 | 0.273 | 5.32x | 1.979 / 2.198 | 1763 / 1835 |

Reproduce with `python benchmarks/apex.py --old PATH_TO_019 --new PATH_TO_020 --runs 21 --out benchmarks/omega-linux.json`. 127 top-level test groups and 347 first-level subtests pass. Windows compiler, installer and Discord executable cross-built and PE imports checked, not executed on Windows. Hello world remains 1536/373 bytes; Discord bot template PE is 15872 bytes. No claim of live Discord connectivity or general Rust superiority.

# TRIPOL v0.17 — HYPER verification

## APEX 0.19 — same runtime-input workloads, TITAN vs APEX

Linux x86-64, 21 interleaved samples per version, Turbo profile, same source and input, output verified each run. Internal time excludes input and printing. Full CLI build time includes process startup and output writes. Raw samples: [apex-linux.json](benchmarks/apex-linux.json). Reproduce with `python benchmarks/apex.py --old PATH_TO_018 --new PATH_TO_019 --runs 21`.

| Workload | TITAN internal ms | APEX internal ms | Speedup | Full CLI build ms (old/new) | ELF bytes (old/new) |
|---|---:|---:|---:|---:|---:|
| division | 17.025 | 10.178 | 1.67x | 2.290 / 2.333 | 1586 / 1586 |
| registers | 11.733 | 7.948 | 1.48x | 1.856 / 1.863 | 1711 / 1679 |
| fibonacci | 4.862 | 3.798 | 1.28x | 1.893 / 1.802 | 1526 / 1491 |
| arrays | 1.412 | 1.404 | 1.01x | 1.931 / 1.813 | 1763 / 1763 |

These are version-to-version measurements, not fresh Rust measurements. Rust was not installed on this Linux host. The user's Windows TITAN Arena comparisons showed Rust faster internally on all three original workloads, with TRIPOL much smaller and faster to compile. Repeat the same named cases on Windows with APEX to assess the gap; cross-host timings cannot establish a win. No general compilation speedup is claimed. The arrays case is new to Arena, not a claim of SIMD support.


## TITAN 0.18 — checked integer workloads

Linux x86-64, same Turbo source, 11 interleaved executions per version after one warm-up, identical verified outputs, internal time measured with `time.now_ns()`.

| Workload | HYPER internal | TITAN internal | HYPER / TITAN build (full CLI) | HYPER / TITAN ELF bytes |
|---|---:|---:|---:|---:|
| 12M division/remainder | 17.435 ms | 15.096 ms | 2.277 / 2.469 ms | 1051 / 1030 |
| 6M four accumulators + branch | 22.882 ms | 12.435 ms | 2.360 / 2.342 ms | 1226 / 1169 |

Reproduce: `python benchmarks/titan.py --old /path/to/hyper --new /path/to/titan --runs 11 --out benchmarks/titan-linux.json`. Raw samples included. Compile time does not show a consistent improvement. `tripol arena --runs 10` is a separate built-in cross-language tool; Rust not installed in this environment, so no Rust measurement is claimed. `tripol arena --tripol-only --runs 3` verified all three outputs.


Linux x86-64. 115 top-level groups and 300 subtests; all three optimized profiles pass the suite (shared runtime test helper selects profile). 100 generated scalar expression cases, plus inherited 20,008 lexer cases. Go vet passed. Windows binaries cross-built, not run here. Both baseline and HYPER compiler binaries used default Go build flags.

| Workload | CORE turbo | HYPER balanced | HYPER turbo |
|---|---:|---:|---:|
| 12M division/remainder iterations | 23.145 ms | 16.311 ms | 15.926 ms |
| 12M mixed arithmetic iterations | 28.144 ms | 22.875 ms | 23.044 ms |
| 6M four-argument scalar calls | 17.997 ms | 16.157 ms | 14.126 ms |
| Fibonacci(30) | 5.283 ms | 5.338 ms | 5.332 ms |

11 interleaved randomized runs per profile, native Linux executable, process startup and stdout capture included; all numeric outputs identical. Full samples and nano comparisons in `benchmarks/hyper-linux.json`. End-to-end nano compile (Windows target, includes process startup and output write): 5,000 calls 23.424→21.675 ms; NOVA 3.105→3.315 ms. No across-the-board compiler throughput improvement established. Hello world remains 1,536 B Windows and 373 B Linux.

Reproduce: `python benchmarks/hyper.py --old /path/to/tripol16 --new /path/to/tripol17 --runs 11 --out results.json`. New internal timing example: `tripol bench examples/bench_hyper_internal.tri --profile turbo --runs 20 --internal`. Rust was not available in this environment.

# TRIPOL v0.16 — CORE verification

Linux x86-64, Go 1.23.2, both compilers built with -trimpath -ldflags="-s -w". 109 top-level test groups, 288 subtests; 20,008 deterministic lexer differential inputs inside one group. Go vet passed. Windows x64 cross-built and PE-inspected, not executed here.

| Workload | v0.15 | v0.16 | Median speed ratio |
|---|---:|---:|---:|
| Compile 5,000 calls (Windows nano output) | 46.225 ms | 17.737 ms | 2.61x |
| Compile NOVA game (Windows nano output) | 2.959 ms | 1.139 ms | 2.60x |
| 12 million signed /8 and %8 operations (12 million iterations, 24 million ops) | 49.481131 ms | 23.159472 ms | 2.14x |
| 12 million mixed arithmetic iterations | 51.639246 ms | 28.724725 ms | 1.80x |
| Recursive Fibonacci(30) | 5.604973 ms | 5.647236 ms | 0.99x |

Compilation: 11 runs per compiler, sequential batches, warmup excluded, compiler's internal `bench-build` includes source loading and compilation, excludes CLI startup and executable writing. Runtime: 11 interleaved runs per executable, native Linux nano, process startup included, checked identical stdout on every run. No general speed claim, no measurement on the user's Windows PC. Exact output bytes for Fibonacci 593/593, signed division 750/680, mixed arithmetic 767/701; game PE 13,312/12,800; hello remains 1,536 Windows and 373 Linux. The two compilers were rebuilt from the supplied source distributions.

Run `python benchmarks/core.py --old /path/to/tripol15 --new /path/to/tripol16 --runs 11 --out results.json` on Linux x64. Raw data: `benchmarks/core-linux.json`. The script does not require third-party Python packages; Python is only a maintainer benchmarking tool.

## Previous measurements

# TRIPOL v0.14 — QUANTUM verification

Linux x86-64, Go 1.23.2. 98 top-level test groups and 288 subtests passed; `go vet .` passed. Windows compiler and installer cross-built, game PE imported symbols inspected; actual game window and new arena Windows execution pending.

| Workload | Manual loop | Native canvas.clear |
|---|---:|---:|
| 1,000 fills × 57,600 RGB pixels | 104.711519 ms | 10.551845 ms |
| Linux executable | 1,169 B | 1,594 B |

11 interleaved timed runs per binary, process startup included, verified output 7111065600 every run. The manual fixture uses the same v0.14 compiler, array size, fill value, and count. The native fixture uses `rep stosq` and RGB validation. This isolates the operation, not an overall language speedup. Sizes: hello Windows/Linux 1,536/373 B, game Windows 10,240 B, memory example Windows 3,072 B. Test importing all 40 libraries without using any produces byte-identical hello binaries.

Memory: 32 arrays of 65,535 ints use exactly 16,777,216 arena bytes; first array remains valid. Repeated scratch restores used bytes. A 513th full-size array exceeds 256 MiB and traps. Committed/enabled pages are not RSS; no general reclamation.

Benchmark Arena: `python benchmarks/arena.py --compiler ./tripol --runs 5` (optional maintainer script). Naive recursive Fibonacci(30), verified 832040 each run, startup included, 5 alternated/reversed runs; C/C++ -O2, Go optimized stripped, Rust -O when installed. Checked TRIPOL vs unchecked arithmetic in C/C++/Go/Rust, not identical safety semantics. Python and JavaScript native file size is not applicable. No measured RSS. Results from this Linux environment in `benchmarks/quantum-linux.json`; Rust unavailable. C/C++ and Go faster on this workload, TRIPOL smallest native executable. Single build times are not statistically robust. No universal language performance claim.

## Previous measurements

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

Scratch test: 2000 arrays of 65535 int64 elements (1,048,576,000 cumulative allocated bytes including length headers) reused within the 8 MiB arena; output 4000. Cumulative allocation is not peak RSS. Escape rejection and recovery cleanup verified. All prior regression tests remain in the Go suite. See tests/history/test-results.txt for exact logs.

## NOVA: clipped full rectangle fill

1000 calls to canvas.rectangle on the same 320×180 buffer (color 7), identical source, output 403200. Eleven alternating executions per release, process startup included, Linux x64. Median v0.14 104.841032 ms (2392 B), v0.15 11.610016 ms (2407 B), approximately 9.03× faster in this workload. Raw samples: benchmarks/nova-linux.json. No inference about GPU performance, FPS or general language speed.

## ECLIPSE — 0.21.0

Linux x86-64, Turbo, 21 interleaved runs per compiler with one warmup. Both use runtime input of 1,000,000 and identical source. Internal elapsed time excludes setup, process startup and output. Copy includes allocation and copies all bytes; the result checksum is verified after the timer. Search uses a missing key to scan the whole array.

| Case | OMEGA ms | ECLIPSE ms | Verified output |
|---|---:|---:|---:|
| Copy | 4.298682 | 2.865491 | 511370976 |
| Search | 1.021374 | 0.581076 | -1 |

Reproduce: `python benchmarks/eclipse.py /path/to/omega /path/to/eclipse`. Raw samples: `benchmarks/eclipse-linux.json`. Rust comparison remains to be executed with `tripol arena --runs 20 --profile turbo --json arena-eclipse.json`; the six TRIPOL workloads were tested locally. No general execution or compiler speed improvement claimed.

## NEXUS (0.22)

21 interleaved warmed samples on Linux x86-64 (AMD EPYC 9V74, AVX2), Turbo, identical result -1 for search over one million int64 values. ECLIPSE median 0.574797 ms; NEXUS 0.156361 ms (3.68x, 72.8% lower time). `python benchmarks/nexus.py /path/to/eclipse /path/to/nexus` reproduces the comparison. See `benchmarks/nexus-linux.json` for raw times. All six Arena cases verified in TRIPOL-only mode (`benchmarks/nexus-arena-linux.json`). Rust unavailable here. No live Discord latency or general compiler speedup measured. Windows executable imports inspected, not run.
