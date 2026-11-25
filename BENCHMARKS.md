# TRIPOL v0.8 — ATLAS measurements

Linux x86-64, Go 1.23.2. All 41 test groups / 213 subtests passed; go vet passed. Windows binaries are cross-built and PE-inspected, not executed on Windows in this environment. v0.7 was confirmed by the project owner in Windows.

| Program | Linux bytes | Windows bytes |
|---|---:|---:|
| Hello (ASCII) | 372 | 1536 |
| Same hello, imports all 6 standard modules | 372 | 1536 |
| Atlas example: math, statistics, text | 2972 | 4096 |

The two hello binaries are byte-for-byte identical per target (nano profile). The library source is embedded in the compiler but unused routines are removed from generated executables. The test suite also successfully runs 400,000 fixed-array total calls without exhausting the 8 MiB arena (borrowed stack view rather than heap copies). These are size/correctness observations, not a speed comparison with other languages.
