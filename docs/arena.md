# APEX Benchmark Arena: TRIPOL vs Rust

Install TRIPOL 0.19 and a native Rust toolchain from https://rustup.rs/. Windows MSVC Rust also needs the MSVC C++ linker and Windows SDK (Visual Studio Build Tools if absent). Reopen CMD and confirm `rustc --version` and `tripol --version`.

```cmd
tripol arena --runs 10 --profile turbo
tripol arena --case registers --runs 20 --profile balanced
tripol arena --tripol-only --runs 10
```

No Python, Cargo or Internet access is required by the arena itself. It writes paired source files into a temporary directory, invokes the installed compilers and runs the resulting executables. Rust command uses `-C opt-level=3 -C overflow-checks=yes -C debuginfo=0 -C strip=symbols`. Both use signed i64 checked arithmetic and the same input read before the timed region. Workloads and both source programs are embedded in `arena.go`, available for review. Three full CLI builds per contestant include startup, linking and writes. One program warm-up is excluded. Execution order alternates and each output must match the expected result. A build times out after 120 seconds and a run after 30 seconds. Temporary artifacts are removed on normal completion and diagnostic unwind.

- `division`: 12 million iterations, output -750000.
- `registers`: 6 million iterations, four accumulators and a branch, output 45000000 21000000 42000000 46500000.
- `fibonacci`: recursive Fibonacci(30), output 832040.
- `arrays`: sum one million initialized i64 values, output 511370976; allocation and initialization occur before timer.

The report prints median full build, internal calculation, end-to-end process time and on-disk binary bytes. Rust's standard library/runtime and OS linkage differ from TRIPOL's minimal runtime, so binary size does not measure the same feature footprint or RAM. Both compilers may optimize the algorithm, including loop transformations. Overflow checks are enabled in both, but other language/runtime checks may differ. Compare on the same PC with low background load. The results apply only to these workloads and do not prove general superiority. `--case all|division|registers|fibonacci|arrays`, `--profile nano|balanced|turbo|debug`, `--runs 1..30` are supported. Rust was not available in the release environment; TRIPOL-only was exercised, Rust integration requires local validation.

Use `tripol arena --runs 20 --profile turbo --json arena-apex.json` to save all raw build/internal/process samples and verified outputs. The named JSON file is created or overwritten. The original three workload sources have not changed.

## ECLIPSE

`--case copy` measures allocation and copy of one million int64 elements, with checksum after timing. `--case search` measures scanning one million integers for a missing value. `--case all` now includes six workloads. The original four workloads are unchanged. Rust uses checked i64 arithmetic; these tests are not a comprehensive language comparison.

## ASCENSION
`--timeout 300` configures per-process timeouts (default 120 seconds; 1–900 allowed). Only the Rust version probe retries once on a timeout; failed timed trials abort instead of silently dropping samples. `--case copy-reuse` compares preallocated equal-sized buffers on both languages. JSON includes nearest-rank internal p95. It is a separate workload, not a direct speedup of the allocating `copy` test.

OVERDRIVE retains the existing Arena workloads. Interactive update checks are disabled for Arena and benchmarks to avoid measurement interference. No new cross-language performance results claimed.
