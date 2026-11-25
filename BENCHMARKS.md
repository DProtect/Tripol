# TRIPOL v0.11 — VISION & GUARD verification

Linux x86-64, Go 1.23.2. 70 test groups and 238 subtests passed; `go vet ./...` passed. Windows compiler and installer cross-built, generated PE images inspected using debug/pe and objdump. Actual Windows UI execution pending. User verified v0.10 client build (5632 B) and demo server, Fibonacci median 88.051 ms and ecosystem compilation median 3.708 ms.

| Example | Linux bytes | Windows bytes |
|---|---:|---:|
| Hello | 373 | 1536 |
| Repository Fibonacci | 593 | 1536 |
| Safe recovery | 2260 | 3584 |
| UI demo | 1351 (unsupported-platform diagnostic) | 3072 |

The repository Fibonacci example's v0.10 and v0.11 Linux binaries differ by one byte, the arithmetic error exit code. No performance speedup claimed. All unused library imports still yield identical output to no imports in nano. The Windows UI demo imports USER32.MessageBoxA and KERNEL32; mixed network/UI import table verified by tests. Native error exit codes and diagnostic JSON were verified locally. This is not a Windows execution claim.
