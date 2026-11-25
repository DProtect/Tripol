<div align="center">

# TRIPOL

**Simple by design. Native by nature.**

A free, experimental programming language that turns readable code directly into native x86-64 executables.

![Version](https://img.shields.io/badge/version-0.14.0-111827) ![License](https://img.shields.io/badge/license-Apache%202.0-2563eb) ![Targets](https://img.shields.io/badge/targets-Windows%20%7C%20Linux-64748b) ![Status](https://img.shields.io/badge/status-experimental-orange)

[Quick start](#quick-start) · [Language](#language) · [Libraries](#standard-library) · [Performance](#performance) · [Documentation](#documentation)

</div>

## Why TRIPOL?

TRIPOL aims to combine approachable syntax, static type checking, fast compilation and small standalone binaries. Its compiler emits x86-64 instructions directly into PE32+ or ELF64 images. It is implemented in Go, but generated programs do not require Go, Python, a C compiler, LLVM, or a virtual machine. Native operating-system APIs provide I/O, networking and Windows dialogs.

| Current release | Built-in modules | Callable library functions* | Hello world (Windows / Linux) |
|---|---:|---:|---:|
| **0.14.0 — QUANTUM** | 40 | 225 | 1,536 / 373 bytes |

*Includes JSON scanner helpers. Sizes are measured for `examples/hello.tri` in nano mode, not a promise for arbitrary programs. TRIPOL is pre-1.0 and is not yet a replacement for mature production languages.*

## Quick start

Install the Windows x64 distribution (`TRIPOL-v0.14-Setup-x64.exe`) when provided by the project maintainer, or [build from source](docs/build.md). The installer is unsigned, per-user and configures PATH. Binary release assets are not currently published in this repository.

```cmd
tripol --version
tripol new Hello
cd Hello
tripol run
tripol build
main.exe
```

Write `main.tri`:

```tripol
name = input("What is your name?")
show "Hello", name

numbers = [10, 20, 30]
show "Total:", sum(numbers)
```

Run `tripol learn es` for the offline Spanish tutorial. See [Inicio rápido en español](docs/inicio.md).

## Language

```tripol
import tripol.collections

func greet(name = "Guest", score = 10) -> string {
    return f"Hello {name}, score: {score}"
}
show greet()
show greet("Ana", 42)

numbers = [10, 20, 30]
expanded = collections.append(numbers, 40)
show expanded
```

Infinity adds opt-in interpolation (string/int/bool bindings), scalar default arguments and copy-based integer collection helpers. The original array remains unchanged. See [language reference](docs/language.md).

### Recoverable errors

```tripol
import tripol.json

func divide(a, b) {
    return a / b
}

try {
    show divide(10, 0)
} catch code {
    show "Recovered error:", code
} finally {
    show "Finished"
}

show json.get_raw("{\"score\":42}", "score")
```

- Inferred local types, typed functions and records, integer arrays, Float64, immutable text, loops and short-circuit logic.
- Checked integer arithmetic and array bounds; structured diagnostics, optional language-frame traces and recoverable exceptions for selected traps.
- Explicit temporary memory regions (`scratch`) with conservative escape restrictions. Ordinary allocations use a 256 MiB virtual reservation with demand-enabled pages; general reclamation remains pending.
- Direct native compilation, unused-code elimination, constant folding, direct self-tail recursion and SSE2 text scanning.
- Offline embedded standard library and a local package manager, Vault.

### QUANTUM 2D (experimental)

```tripol
import tripol.game
import tripol.canvas

pixels = canvas.create(320, 180)
game.window("QUANTUM", 960, 540)
while game.running() {
    if game.key(27) { break }
    canvas.clear(pixels, 1052688)
    canvas.rectangle(pixels, 320, 180, 140, 70, 40, 40, 16777215)
    game.present(pixels, 320, 180)
    game.wait(16)
}
game.close()
```

The software framebuffer (RGB24, up to 65,535 pixels) is written in TRIPOL; its fast clear and Windows presentation use native x86-64 templates and GDI. Windows window and keyboard execution still require validation. `tripol new Demo --template game` creates a movable-square demo. See [graphics guide](docs/graphics.md).

### Native Windows forms (experimental)

```tripol
import tripol.ui
name = ui.form("TRIPOL", "Your name:", "Say hello")
if len(name)>0 { show ui.message("TRIPOL", "Hello " + name) }
```

This is a modal form with a text field and submit/cancel buttons, not a general GUI toolkit. The v0.12 form was successfully demonstrated by the project owner on Windows; HTTPS and the new message-box variants still require Windows runtime validation.

## Standard library

| Category | Modules |
|---|---|
| Data and algorithms | `buffer`, `arrays`, `math`, `stats`, `sort`, `search`, `json`, `random`, `collections`, `matrix`, `sequences`, `bits` |
| Text and conversion | `text`, `validate`, `convert`, `encoding`, `compression`, `checksum` |
| System and I/O | `console`, `files`, `system`, `time`, `errors`, `datetime`, `path`, `memory` |
| Network | `network`, `http`, `https`, `discord` |
| UI and geometry | `ui`, `colors`, `geometry`, `canvas`, `game`, `events`, `image`, `animation`, `physics`, `vector` |

Use `import tripol.math`, then `math.square(7)`. Run `tripol libs` or `tripol docs json` to inspect exact available functions offline. See [library guide](docs/libraries.md) for supported operations and restrictions. HTTPS and the Discord public gateway discovery helper are Windows-only; authenticated bots are not implemented.

### Vault — local packages

```cmd
tripol add mylib mylib.tri --version 1.0.0
tripol packages
tripol verify
```

Import `vault.mylib`. The lockfile pins the source hash and version. Vault does not yet have a public online registry. Only install trusted code.

## Performance

QUANTUM framebuffer fill benchmark (Linux, 1,000 fills of 57,600 RGB pixels): manual loop 104.712 ms; native `canvas.clear` 10.552 ms (9.92x in this workload). Both return the same sum. The new [Benchmark Arena](benchmarks/arena.py) compares naive Fibonacci(30) against available compilers/interpreters; C, C++ and Go were faster in this run, while TRIPOL emitted the smallest native binary. [Raw results](benchmarks/quantum-linux.json).


| Workload (Linux) | v0.12 | v0.13 |
|---|---:|---:|
| Repeated length of a constant 4 KiB literal (200,000 calls) | 36.487 ms | 0.463 ms |
| Executable for the literal-length workload | 4,776 B | 595 B |
| Square-function loop (2 million calls) | 6.727 ms | 6.866 ms |

The first workload is intentionally a constant-folding microbenchmark: length is resolved at compile time in nano builds when the built-in `len` receives a string literal. The square workload did not demonstrate a reliable improvement; inline TIR is experimental. See [BENCHMARKS.md](BENCHMARKS.md) for measurement details.

### Previous measurements

| Workload | v0.11 | v0.12 |
|---|---:|---:|
| 20,000 length scans of a 32 KiB string (Linux) | 218.09 ms | 26.89 ms |
| Hello world executable (Windows) | 1,536 B | 1,536 B |

The length benchmark compares identical source and output on the same Linux environment (seven interleaved timed executions per version; process startup included). The speedup is workload-specific. No claim of universal superiority over C, Rust or Go. See [BENCHMARKS.md](BENCHMARKS.md) and `tripol bench` / `tripol bench-build`.

## Compiler architecture

```text
.tri → lexer → parser → static checks → TIR subset / native code generator → PE / ELF
             modules / Vault ─────────────┘
```

Go is the compiler's implementation language, not an intermediate language for user programs. The maintained assembly templates and generated Go tables encode the native helpers. TIR R2 lowers small pure int64 return expressions to a typed tree; eligible zero/one-parameter functions inline in nano, and up to four parameters in the experimental turbo profile, with overflow checks. All other code uses the existing backend. `tripol ir FILE.tri` inspects this subset. General control-flow TIR and full self-hosting remain future work; `bootstrap/demo.tri` exercises an experimental ASCII lexer written in TRIPOL.

## Documentation

| Guide | Description |
|---|---|
| [Inicio rápido](docs/inicio.md) | Spanish installation and first program |
| [Language reference](docs/language.md) | Types, statements, exceptions and memory |
| [Standard library](docs/libraries.md) | Modules, limits and Windows-specific features |
| [Compiler and build](docs/build.md) | Toolchain, native targets and tests |
| [Graphics](docs/graphics.md) | Software framebuffer and native game window |
| [Guard diagnostics](docs/errors.md) | Errors, recovery, suggestions and tracing |
| [Benchmarks](BENCHMARKS.md) | Reproducible measurements |
| [Changelog](CHANGELOG.md) | Current release changes |
| [Roadmap](ROADMAP.md) | Future work, not guaranteed release dates |
| [Contributing](CONTRIBUTING.md) | Development and review guidelines |

**Current limitations:** experimental x86-64 targets only; ordinary arena reservation capped at 256 MiB (no general GC); no threads, FFI or full GUI toolkit; limited networking and JSON extraction; no general exceptions for native I/O or memory failures. New Windows game windows, dialog variants and HTTPS require execution testing. See the linked guides before using TRIPOL for long-running or production software.

## License

TRIPOL is free and open-source under the [Apache License 2.0](LICENSE). Contributions are welcome subject to repository access. The project intends to keep the language, compiler and official libraries free.
