# TRIPOL v0.3.0 — NANO REVOLUTION

Apache-2.0, free experimental native programming language. The standalone Go compiler emits x86-64 PE32+ (Windows) and ELF64 (Linux) directly. No Python, GCC, Clang or Go installation is required for end users. Go 1.23+ is needed only to rebuild the compiler.

## Windows x64

Run `TRIPOL-v0.3-Setup-x64.exe` to update the per-user installation at `%LOCALAPPDATA%\Programs\TRIPOL`. The installer preserves existing `hello.tri` and adds the directory to the user PATH without administrator privileges. Close and reopen CMD; run `tripol --version`. The installer uses Windows PowerShell only to modify PATH. It is unsigned. v0.2 was confirmed working by the project owner on Windows; v0.3 was cross-compiled and PE-inspected, not executed in Windows yet. Alternatively use the portable ZIP and run `tripol.exe` directly.

```cmd
tripol --version
tripol new MyProject
cd MyProject
tripol run main.tri
tripol build main.tri
main.exe
```

`tripol check FILE.tri`, `tripol run FILE.tri [--profile nano|debug]`, `tripol build FILE.tri [-o OUTPUT] [--target windows|linux] [--profile nano|debug]`. Default profile nano. `debug` disables the new optimization passes (not a source-level debugger). `--target` is x86-64 only. `tripol new` creates a new folder and refuses to overwrite an existing one.

## Language additions

```tripol
values = [10, 20, 30]
for i in range(0, len(values)) {
    values[i] = values[i] * 2
    show values[i]
}
```

Output: 20, 40, 60 on separate lines. Arrays hold signed 64-bit integers and are fixed-length, local, zero-based and stack-allocated; 1–400 elements subject to the total 400-slot function limit. Index reads/writes have runtime bounds checks. No array copying, passing, returning, resizing or nesting. `len(array)` returns length (a user-defined `len` takes precedence). `const` arrays cannot be changed. `for i in range(start,end)` evaluates integer bounds once, iterates ascending with exclusive end, and has an immutable scoped index. `break`/`continue` apply to the nearest loop. UTF-8 BOM is accepted for Windows editors.

All v0.2 language features remain: int64, bool, immutable string literals, functions, recursion, checked arithmetic, if/else, while, short-circuiting, string equality. Functions cannot capture top-level locals. No floats, dynamic strings, input, modules/imports, networking or heap allocator yet. Not a sandbox. No production-grade ABI or complete memory safety; ordinary recursion may exhaust stack. Generated PE uses Windows system DLL functions; Linux uses syscalls.

## Optimization R1

- Velocity: safe constant folding, smaller immediates and local loads/stores, immediate arithmetic, direct conditional branches on integer constants. Checked arithmetic traps are preserved.
- ZeroMemory: lexical stack slot reuse; fixed arrays use inline stack storage, no garbage collector. This is not general heap management.
- MicroBinary: reachability-based unused function removal and omission of unused runtime helpers. All functions are checked before removal.
- Self-tail-call elimination: exact `return current_function(...)` evaluates arguments first, then reuses its frame. General recursion is unchanged.
- Native backend remains direct machine code. No TIR yet; these are AST and native code generation optimizations, not a new scientific algorithm.

## Build from source

```sh
go test -v .
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o tripol .
mkdir -p installer/payload
cp examples/hello.tri installer/payload/hello.tri
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o installer/payload/tripol.exe .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o TRIPOL-v0.3-Setup-x64.exe ./installer
```

See BENCHMARKS.md and test-results.txt. v0.1 Python/C source may remain in the Git repository as historical reference; it is not used by v0.3.
