# TRIPOL v0.2.0 — FORGE

Free, experimental Apache-2.0 programming language. A standalone Go compiler generates x86-64 machine instructions, ELF64 (Linux) and PE32+ (Windows) directly. No Python, GCC, Clang, Go installation, or interpreter is needed to use the distributed compiler. Go is required only to build the compiler from source. The v0.1 Python/C bootstrap remains in the repository for reference.

## Windows installation (x64)

Download and run `TRIPOL-v0.2-Setup-x64.exe`. The unsigned console installer copies the compiler and a sample to `%LOCALAPPDATA%\Programs\TRIPOL`, then updates the current user's PATH via Windows PowerShell. No administrator rights required. Close and reopen your terminal (or sign out/in if PATH has not refreshed). The installer can also be run with `--quiet`, or `--uninstall` to remove the compiler and its PATH entry without deleting other files. The installation and generated Windows programs have not been executed on a Windows host yet.

```powershell
tripol --version
tripol run "$env:LOCALAPPDATA\Programs\TRIPOL\hello.tri"
tripol build main.tri
.\main.exe
```

## Language

```tripol
func square(n) {
    return n * n
}
show square(8)
```

Supports signed 64-bit integers, bool, immutable UTF-8 string literals, variables and constants, typed functions, recursion, if/else, while, short-circuit logic, checked integer arithmetic, and string equality. Code generation uses a simple stack-based internal calling convention and direct native instructions. The Windows runtime imports GetStdHandle, WriteFile and ExitProcess from KERNEL32.dll; Linux uses system calls. String output on Windows is UTF-8 bytes (console encoding may affect display). No ASLR/unwind metadata in generated PE files yet. Not a sandbox; do not run untrusted source.

`tripol check FILE.tri`, `tripol run FILE.tri`, `tripol build FILE.tri [-o OUTPUT] [--target windows|linux]`. The default target is the host OS; only x86-64 output is supported. Only `--profile nano` is currently accepted by the new backend (no mature optimization passes yet). Legacy `emit-c` and other profiles remain in v0.1 only. No arrays, floats, dynamic strings, modules, input, heap management, GPU, networking, or memory budget system yet. No production-grade ABI or optimization claims.

## Build from source

Go 1.23+; standard library only. Linux/macOS shell:

```sh
go test -v .
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o tripol .
mkdir -p installer/payload
cp examples/hello.tri installer/payload/hello.tri
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o installer/payload/tripol.exe .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o TRIPOL-v0.2-Setup-x64.exe ./installer
```

## Measurements

Linux x86-64: native hello executable 662 bytes; recursive fib(30) executable 827 bytes (output 832040). Windows x64 hello PE 2048 bytes, parsed with objdump; no Windows execution yet. Compiler and installer are larger (about 2.5 MiB and 4.3 MiB). No meaningful RSS or comparative performance claim yet. 44 language test cases plus format and deterministic output checks pass on Linux. See test-results.txt in the source distribution.
