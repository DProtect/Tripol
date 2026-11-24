# TRIPOL v0.1.0 — ORIGIN

Experimental, free, Apache-2.0 licensed programming language bootstrap.
TRIPOL targets simple syntax and small native executables. This release is a
working foundation, not a production-ready language or a complete optimizer.

## Start

Requires Python 3.10+ and GCC or Clang (a working C11 toolchain, including
linker and C library). No third-party Python packages are needed. Windows
users can use a configured MinGW-w64 GCC toolchain. Linux x86-64 was tested;
Windows and macOS have not been tested. Native binaries do not require Python,
but they may depend on the platform C runtime.

```sh
python tripol.py --version
python tripol.py check examples/demo.tri
python tripol.py run examples/demo.tri
python tripol.py build examples/hello.tri --profile nano
python tripol.py emit-c examples/demo.tri -o demo.c
python -m unittest discover -s tests -v
```

On Windows, the build command defaults to `examples/hello.exe`; on Linux it
writes `examples/hello`. Use `--cc gcc` to select GCC explicitly. A C compiler
must be installed separately. A Windows binary is not included.

## Example

```tripol
func sumar(a, b) {
    return a + b
}
const nombre = "TRIPOL"
show nombre
show sumar(10, 20)
i = 0
while i < 3 {
    show i
    i = i + 1
}
```

## Design

Source -> lexer -> recursive-descent / precedence parser -> AST -> static type
checking and C generation -> native compiler and linker.

* Core R1: lexical analysis, parsing, typed functions, lexical blocks.
* Forge R1 / NativeCore prototype: C11 backend with GCC/Clang.
* SafeCore prototype: type checks, constant checks, checked signed arithmetic,
  division checks, required function return paths. Not full memory safety.
* MicroBinary prototype: size optimization and dead-section removal delegated
  to the C toolchain. Not a novel optimizer.
* CLI: check, emit-c, build, run. Profiles nano (-Os, section GC and stripping
  except macOS), balanced (-O2), turbo (-O3), debug (-O0 -g).

## Language rules and limitations

- Statements end with newline, semicolon, block close or EOF. Blocks use braces.
- `let x = ...` introduces a mutable local, `const x = ...` an immutable local.
  Bare assignment declares a mutable local if not already declared, otherwise
  assigns to an existing mutable variable. Shadowing is currently rejected.
- `int` is a signed 64-bit integer, `bool` is a boolean, `string` is immutable
  UTF-8 text backed by static C strings. Strings support JSON-style escapes,
  equality and printing; embedded NUL is rejected. No concatenation yet.
- Functions are top-level, can call each other and recurse. Parameter and
  return types default to `int` if omitted. Functions cannot capture top-level
  variables. Function argument and non-short-circuit operand evaluation order
  is unspecified in this bootstrap (inherits C evaluation order).
- Conditions must be boolean; no implicit conversions. Integer division
  truncates toward zero. Arithmetic overflow traps at runtime. The literal
  -9223372036854775808 is not directly accepted; use -9223372036854775807 - 1.
- No arrays, floats, objects, modules, imports, heap management, IO other than
  console output, networking, GPU support, or automatic parallelization yet.
- Recursion can exhaust the native stack. This is not a sandbox: never compile
  or execute untrusted source without process isolation.
- The current frontend is Python, not Rust. The intended future Rust frontend,
  custom IR, ZeroMemory, Evolution and resource budgets are not implemented.
- No claim of outperforming Rust or C. Current executable size primarily comes
  from the C toolchain and platform ABI.

## Validation

34 automated tests passed in Linux x86-64. See BENCHMARKS.md.
Includes syntax/type rejection, checked integer operations,
control flow, recursive Fibonacci, and profile checks. The tests use GCC;
example binaries were also built using Clang 17. No Windows binary is included.

Licensed under Apache-2.0. See LICENSE.
