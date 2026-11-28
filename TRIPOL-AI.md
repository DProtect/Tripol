# TRIPOL AI Reference — v0.29.0 DIPO HTTP 2 (development)

**Referencia oficial para asistentes de IA / Official reference for AI assistants**

Experimental v0.29 development branch · Apache-2.0 · Windows and Linux x86-64. Browser support is a separate subset. [Release and downloads](https://github.com/Tripollas/Tripol/releases/tag/v0.28.0) · [Documentation](https://tripollas.github.io/Tripol/docs/)

## Cómo usar este archivo / How to use this file

Usa TRIPOL v0.29.0 de desarrollo. Prioriza SIMPLE R3, la guía de Discord de esta versión y las firmas del catálogo actual frente a explicaciones de versiones anteriores. No inventes módulos, sintaxis, métodos, sobrecargas ni capacidades. Escribe ejemplos sencillos y legibles, con tipos explícitos cuando el compilador los necesite. Indica qué requiere el compilador nativo, qué funciona sin conexión y qué necesita una aplicación y un servidor de Discord autorizados.

Use development TRIPOL v0.29.0. Prefer the SIMPLE R3 section, the current Discord guide and exact catalog signatures over historical descriptions. Do not invent modules, syntax, methods, overloads or capabilities. Prefer readable code; add annotations when inference is insufficient. Distinguish native, browser and offline execution from connected validation. Test with the actual compiler before claiming that a program works.

**Current catalog / Catálogo actual:** 64 libraries, 466 functions and 251 methods (717 entries).

**Version precedence / Prioridad de versión:** SIMPLE R2 extends older SIMPLE R1, Infinity and Quantum rules: contextual monomorphic inference, implicit final-expression returns, named and inline callbacks, scalar literal constant captures and static overloads are supported in their documented scope. Dynamic captures, general polymorphism and arbitrary named arguments are not established. Earlier limitations in historical sections describe earlier versions.

**Discord:** A returned interaction payload or an explicit `ctx.respond(payload)` sends the initial response; a discarded `ctx.reply(...)` expression does not send it. Use the shared launcher on Windows and Linux for advanced networking; directly executed Windows PE uses the older WinHTTP path. Keep DISCORD_TOKEN in the process environment, never in source, logs or shared binaries. Long live validation, connected-bot comparisons and channel lock recovery after launcher restart remain pending. No production-readiness claim follows from local tests.

Sources: release source `e635e7d26e6b347cb216e1476bffc193179c6236` (language, library and Discord guides); public README `9390b7f360bbff820878fb1836a6216595a452a1`; web catalog and categories `de3b2448bad3bcefd718ea8e1a786937bc6b0b18`. Remaining guide excerpts and 14 lessons retain their earlier snapshot from `ad41ca60704675fb68a100dfa79e2a2e8210c5be` and `279b0e2b4c4e5ab7c002a3ef79f21c082031a8d4`. Older version-specific instructions are historical, not current release requirements.

## SIMPLE R3 collections and typed decoding / Colecciones y decodificación


Use `list([Point(1), Point(2)])` to infer a homogeneous list of records. Declare an empty typed list with `let points:list[Point]=list()`. A nested type such as `list[list[Point]]` keeps the element type at each level. Use `rows.get(0).get(1)` to read a nested element and bind the inner list to a local variable before mutating it. Lists share a checked handle when assigned to another variable. They are allocated in the program arena. A `scratch` scope cannot extract an aggregate element from an outer list since it could let temporary storage escape. Search and sorting on object lists require an explicit field selection.

Usa `list([Point(1), Point(2)])` para inferir una lista homogénea de estructuras. Declara una lista vacía con `let points:list[Point]=list()`. Un tipo como `list[list[Point]]` conserva el tipo de cada nivel. Lee un elemento con `rows.get(0).get(1)` y asigna la lista interna a una variable local si vas a modificarla. Los alias comparten un identificador verificado. Las listas se asignan en la arena del programa. Un bloque `scratch` restringe el acceso a elementos compuestos de listas exteriores para que la memoria temporal no salga del bloque.

`json.decode(Report, raw)` uses the struct declaration at compile time. It handles user structs, nested structs, typed lists, text, booleans and signed 64-bit integers. It raises an assertion on malformed JSON, missing required fields, duplicate keys or incompatible values. Unknown fields are accepted. Fractional fields, null, omitted fields with defaults and function references need another schema. Use the same source on native targets and in the Playground. [Runnable example](../examples/simple_r3_objects.tri)

`json.decode(Report, raw)` usa la estructura en tiempo de compilación. Admite estructuras anidadas, listas tipadas, cadenas, booleanos y enteros de 64 bits. Genera un error de aserción ante JSON inválido, campos requeridos ausentes, claves duplicadas o valores incompatibles. Admite campos adicionales. Los decimales, valores nulos, campos omitidos con valor por defecto y referencias a funciones requieren otro esquema. [Ejemplo ejecutable](../examples/simple_r3_objects.tri)


## Index / Índice

- Project overview and downloads
- Native language and toolchain guides (including SIMPLE R2)
- Discord Revolution guide
- Complete standard library catalog (717 entries)
- Additional official guides and 14 lessons

---

# Project overview

<div align="center">

<img src="assets/tripol-emblem.png" alt="TRIPOL official logo" width="220" />

# TRIPOL

**Simple by design. Native by nature.**

A free, experimental programming language that turns readable code directly into native x86-64 executables.

![Version](https://img.shields.io/badge/version-0.28.0-111827) ![License](https://img.shields.io/badge/license-Apache%202.0-2563eb) ![Targets](https://img.shields.io/badge/targets-Windows%20%7C%20Linux-64748b) ![Status](https://img.shields.io/badge/status-experimental-orange)

[Official website](https://tripollas.github.io/Tripol/) · [Online docs](https://tripollas.github.io/Tripol/docs/) · [Quick start](#quick-start) · [Language](#language) · [Libraries](#standard-library) · [Performance](#performance) · [Documentation](#documentation)

</div>

## DISCORD REVOLUTION — v0.28.0 experimental release

**0.28.0 — DISCORD REVOLUTION** introduces SIMPLE R2 contextual inference, inline callbacks and static overloads, alongside slash commands, context menus, modals, embeds, Components V2, emojis, cooldowns, permissions and server administration. Windows and Linux support experimental Discord networking through the shared launcher, including Resume, compression, rate limits and multipart. Prolonged live validation remains pending. [Release and downloads](https://github.com/Tripollas/Tripol/releases/tag/v0.28.0). See the bilingual guide and validation matrix for supported features and open work. [Discord guide](docs/discord.md) · [Status and evidence](docs/discord-status.md)

## AI-ready reference / Referencia para IA

[TRIPOL-AI.md](TRIPOL-AI.md) includes v0.28 Discord additions and the historical v0.27.1 documentation snapshot. The current online catalog has 64 modules, 449 functions and 196 methods (645 entries). Download the Markdown and attach it to your preferred AI assistant. [Direct download](https://tripollas.github.io/Tripol/TRIPOL-AI.md).

## Why TRIPOL?

TRIPOL aims to combine approachable syntax, static type checking, fast compilation and small standalone binaries. Its compiler emits x86-64 instructions directly into PE32+ or ELF64 images. It is implemented in Go, but generated programs do not require Go, Python, a C compiler, LLVM, or a virtual machine. Native operating-system APIs provide I/O, networking and Windows dialogs.

| Current release | Built-in modules | Callable library functions* | Hello world (Windows / Linux) |
|---|---:|---:|---:|
| **0.28.0 — DISCORD REVOLUTION (experimental release)** | 64 | 449 functions + 196 methods (645 catalog entries) | see v0.24 baseline |

*Includes JSON scanner helpers. Sizes are measured for `examples/hello.tri` in nano mode, not a promise for arbitrary programs. TRIPOL is pre-1.0 and is not yet a replacement for mature production languages.*

## QUANTUM X — 0.27.1

- SIMPLE R1: scalar function inference, typed scalar lists, fixed-schema mixed-scalar dictionary literals and the native teaching shell (`tripol shell`).

- SIMPLE R1: expression interpolation, integer list initializers, readable error descriptions and nine offline bilingual [Academy lessons](docs/academy.md) (`tripol learn es list`).

- Experimental native WGL/OpenGL triangles, GDI fallback, XInput and focused mouse input, clipped tile maps, sprite transforms and particle bursts. The expanded pixel-art demo includes a menu, scrolling camera and optional audio. Actual hardware acceleration requires a compatible driver and local validation.

- Indexed iteration, list subscripts and trailing scalar defaults for record fields.
- Standalone `tripol lsp` editor protocol and `tripol fmt` indentation formatter with VS Code integration.

- 64 standard modules, including SHA-256, PCM WAV creation, Windows WAV playback, explicit option parsing, native process ID/exit and a joined native Windows worker for checked integer sum. General thread callbacks, process spawning and OS argument parsing are not yet included.

- MEMORY R3 extends automatic temporary reclamation to locally owned typed lists and outer allocations around nested loops whose inner iterations allocate no temporaries.
- TIR R7: turbo emits multi-block native SSA for eligible call-free integer and boolean functions, with phi edges, liveness and register allocation. The checked emitter handles the remaining language.
- INCREMENTAL R1 reuses parsed local and embedded modules in cached builds. VECTOR R2 dispatches long constant integer array fills to OS-enabled AVX2, with SSE2 and scalar-tail handling.
- Native checked x86-64 atomics (load, exchange, compare-and-swap, add); no general user thread API yet.
- UPDATE R2: transactional per-user installation with confirmation, rollback backups and a Windows helper; portable/system installs remain manual.

- Native byte I/O (`tripol.io`), bounded binary cursor (`tripol.stream`) and bounded literal/dot/star matching (`tripol.regex`).
- Growable integer lists: `list()`, `add`, `get`, `set`, `remove_at`, `clear`, `size`, `capacity`; stable reference semantics across growth.
- GENESIS R3 extends the experimental Linux ELF subset with typed integer functions, forward calls, mutual recursion, short-circuit boolean conditions, else-if, `for`/range, `break`, `continue`, repeated `show`, and direct binary output in `compiler/direct.tri`. Experimental PE32+ console backend also emits Windows executables; not self-hosting yet.

### 0.27.0 foundation

- GENESIS R3 experimental multiline ELF generation: checked runtime division/remainder, comparisons, `if`/`else`, `while`, and a file driver. `tripol genesis` removes the Python transport dependency but the host Go CLI still writes the output; no PE emission or self-hosting yet. [GENESIS](docs/genesis.md).
- SYNTAX R2: checked indexed compound assignments `+=`, `-=`, `*=`, `/=`, `%=`; index evaluated once.
- Three new standard modules: bounded `heap`, strict UTF-8 scalar helpers and single-record `csv` helpers. [Libraries](docs/libraries.md).
- UPDATE R2 first stage: verified download progress, `--download-only` and receipt history; installation remains manual.
- TIR R6 first stage: statement-level CFG inspection; existing checked optimizer remains the production backend. No full SSA/register allocation.
- VS Code completion, hover, local definition lookup and explicit compiler commands; Windows software-rendered game demo and [Validation Center](docs/testing-v0.27.md).

## OVERDRIVE (v0.26)

- UPDATE R1: interactive update notices, 24-hour cache, `tripol update --check` and verified official release downloads with `tripol update`. No automatic installation. [Update guide](docs/update.md).
- GENESIS R2: integer assignments and runtime x86-64 arithmetic with overflow checks, emitting Linux ELF. Still not self-hosting or a full compiler. [GENESIS guide](docs/genesis.md).
- `tripol.bitset` and `tripol.profiler`: compact integer bitsets and monotonic timers.

## EVOLUTION (v0.25)

GENESIS R1 is the first experimental native executable generator written in TRIPOL: it parses a constant integer `show` expression and emits Linux x86-64 ELF bytes. It is not self-hosting yet. See [GENESIS](docs/genesis.md). New bounded reusable `stack` and `queue` modules, plus checked scalar TIR identity simplification. The Go compiler remains the production implementation. Performance gains and binary sizes must be remeasured before release.

## ASCENSION (v0.24)

New: `arrays.copy_into(destination, source)` reuses equal-sized storage and returns false without modifying the destination on a size mismatch. Arena now supports `--timeout SECONDS`, `copy-reuse`, and p95 in JSON. The static bilingual website is maintained separately on the [`gh-pages` branch](https://github.com/Tripollas/Tripol/tree/gh-pages); see [website documentation](docs/website.md).

## Repository layout

- `core/`: production Go compiler and white-box regression tests.
- `compiler/`: GENESIS compiler and demos written in TRIPOL.
- `stdlib/`: the 64 standard-library modules, embedded at build time.
- `native/`: assembly source and code generator.
- `tests/`: CLI integration tests and archived test evidence.
- `examples/`, `benchmarks/`, `installer/`, `vscode/`, `docs/`: samples, reproducible measurements, distribution and documentation.

The small root `main.go` keeps `go build -o tripol .` compatible. Run the complete compiler and integration suite with `go test ./core ./tests`.

## Community and tools

[Releases](https://github.com/Tripollas/Tripol/releases) · [VS Code extension](vscode/) · [Practical projects](examples/projects/) · [Security policy](SECURITY.md) · [Release and community checklist](docs/ecosystem-release.md)

## Quick start

Install the [Windows x64 setup](https://github.com/Tripollas/Tripol/releases/download/v0.28.0/TRIPOL-v0.28-Setup-x64.exe), [Linux x64 archive](https://github.com/Tripollas/Tripol/releases/download/v0.28.0/TRIPOL-v0.28-Linux-x64.tar.gz), or [build from source](docs/build.md). The Windows installer is unsigned, per-user and configures PATH. See the [experimental v0.28.0 release](https://github.com/Tripollas/Tripol/releases/tag/v0.28.0) for all assets, including the VS Code VSIX and SHA-256 checksums. Physical Ryzen 5 3400G and Radeon Vega 11 / WebGPU validation remains pending.

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

## Discord Revolution (Windows and Linux)

```tripol
import tripol.discord

bot = discord.bot()

bot.slash("hola", "Saluda", func(ctx) {
    ctx.reply("Hola desde TRIPOL!")
})

bot.start()
```

Set `DISCORD_TOKEN` in the process environment and use `tripol run` in an authorized test server. The shared launcher supports automatic reconnect and Resume, REST rate limits, compression and multipart; directly launched Windows PE files use the older WinHTTP path. The Discord catalog contains 98 functions and 196 methods. Live networking remains experimental. See [Discord guide](docs/discord.md), [validation](docs/discord-validation.md) and [SDK benchmarks](benchmarks/discord_sdks/).

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

The software framebuffer (RGB24, up to 4,194,304 pixels (2048 × 2048)) is written in TRIPOL; its fast clear and Windows presentation use native x86-64 templates and GDI. The v0.14 window was reported working by the project owner; the expanded v0.15 path is PE-inspected and still requires Windows runtime validation. `tripol new Demo --template game` creates a 720p movable-sprite demo. See [graphics guide](docs/graphics.md).

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
| Data and algorithms | `buffer`, `bitset`, `arrays`, `math`, `stats`, `sort`, `search`, `json`, `random`, `collections`, `matrix`, `sequences`, `bits`, `simd` |
| Text and conversion | `text`, `validate`, `convert`, `encoding`, `compression`, `checksum` |
| System and I/O | `console`, `files`, `system`, `time`, `errors`, `datetime`, `path`, `memory`, `env` |
| Network | `network`, `http`, `https`, `discord` |
| UI and geometry | `ui`, `colors`, `geometry`, `canvas`, `game`, `events`, `image`, `animation`, `physics`, `vector`, `sprite`, `camera`, `particles` |
| Data and testing | `dictionary`, `stack`, `queue`, `testing`, `profiler` |

Use `import tripol.math`, then `math.square(7)`. Run `tripol libs` or `tripol docs json` to inspect exact available functions offline. See [library guide](docs/libraries.md) for supported operations and restrictions. Discord Gateway and REST support Windows and Linux through the experimental shared launcher; see [Discord setup](docs/discord.md). General HTTPS support is platform-dependent.

### Vault — local packages

```cmd
tripol add mylib mylib.tri --version 1.0.0
tripol packages
tripol verify
```

Import `vault.mylib`. The lockfile pins the source hash and version. Vault does not yet have a public online registry. Only install trusted code.

## Performance

OMEGA 0.20 specializes checked zero-based dynamic integer array reductions using four sequential checked additions per iteration and a bounded pointer traversal. The same Arena array workload improved from 1.452 to 0.273 ms internally (5.32x, Linux x64, 21 alternating samples). Other workloads were essentially unchanged; see [raw measurements](benchmarks/omega-linux.json). `simd.equal` implements SSE2 exact equality, not AVX2 or general auto-vectorization.


APEX 0.19 adds a private register ABI for structurally verified binary integer recurrences, loop-invariant power-of-two arithmetic hoisting, nonnegative-index specialization, direct modulo-zero branches and paired quotient/remainder lowering with separate checked accumulator updates. In 21 interleaved Linux Turbo samples the division, registers and recursive workloads fell from 17.025 to 10.178 ms, 11.733 to 7.948 ms and 4.862 to 3.798 ms respectively. The array workload was essentially unchanged. [Raw samples](benchmarks/apex-linux.json). Arena adds a checked array-sum workload and `--json FILE` for raw samples. Full SSA, general register allocation, AVX2 and memory reclamation are still pending.

TITAN 0.18 extends RegisterFlow to four outer integer accumulators and call-free conditional branches with break/continue. It saves/restores added registers and falls back when branches declare locals, calls, nested loops, guards or tracing require the generic backend. Signed power-of-two division/remainder now uses a shorter CQO-based lowering; overflow checks are retained. In Linux Turbo tests, internal division time was 17.435 → 15.096 ms, and four-accumulator time 22.882 → 12.435 ms. Compilation showed no consistent improvement. [Raw samples](benchmarks/titan-linux.json).

**TRIPOL vs Rust:** install Rust with [rustup](https://rustup.rs/), make sure `rustc --version` works (and the native linker is available), and run `tripol arena --runs 10`. It compiles four equivalent i64 programs, uses runtime inputs, enables Rust overflow checks, validates each output, interleaves runs, and reports full compile, internal compute, process and executable size. `--case division|registers|fibonacci|arrays` selects a workload; `--tripol-only` works without Rust. Rust was not installed in the release test environment, so no cross-language results are claimed. See [Arena guide](docs/arena.md).

HYPER 0.17 introduces register-backed straight-line integer loops and a private register calling convention for eligible pure scalar leaf functions. `balanced` enables register loops and register leaf calls; `turbo` also inlines up to four scalar arguments using expression registers. `nano` remains the default compact profile. `tripol bench-build` batches small builds and supports profile selection; `tripol bench --internal` reads one explicit `TRIPOL_NS: N` marker produced with `time.now_ns()` / `time.elapsed_ns()`. [HYPER benchmark data](benchmarks/hyper-linux.json) and [reproduction script](benchmarks/hyper.py) compare identical outputs; improvements are specific to the workloads.

CORE 0.16 replaces the regex lexer with a single-pass byte scanner, indexes native relocations, caches scalar TIR per function and avoids linking a discarded validation executable. It preserves full source validation before optimization. Checked direct memory arithmetic, signed constant power-of-two division/remainder, and immutable range-counter lowering improve generated code. In the Linux comparison with 0.15, a 5,000-call program compiled in 17.737 ms versus 46.225 ms, and the game example in 1.139 ms versus 2.959 ms (Windows x64 output, nano, source loading included; compiler startup and executable writes excluded). Signed division/remainder runtime improved from 49.481 to 23.159 ms; recursive Fibonacci remained effectively unchanged. [Reproduction script](benchmarks/core.py) · [Raw measurements](benchmarks/core-linux.json). These results are workload-specific.

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

Go is the compiler's implementation language, not an intermediate language for user programs. The maintained assembly templates and generated Go tables encode the native helpers. TIR R4 with CORE caching lowers small pure int64 return expressions to a typed tree; eligible zero/one-parameter functions inline in nano, and up to four parameters in the experimental turbo profile using registers; balanced uses a private register ABI for eligible multi-argument leaves, with overflow checks. All other code uses the existing backend. `tripol ir FILE.tri` inspects this subset. General control-flow TIR and full self-hosting remain future work; `compiler/demo.tri` exercises an experimental ASCII lexer written in TRIPOL; `compiler/parser_demo.tri` exercises its checked integer expression parser.

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

ECLIPSE 0.21 adds native checked independent array copy and bounded linear search. Linux interleaved 21-sample medians: copy 4.299 → 2.865 ms; search 1.021 → 0.581 ms, one million int64 values. Arena adds `copy` and `search`. These are workload-specific improvements; Rust was not installed on the release host.

NEXUS 0.22 adds CPU/OS-gated AVX2 first-match search with a scalar fallback. Linux 21-sample medians: 0.575 → 0.156 ms (one million int64 values, no match). See [raw samples](benchmarks/nexus-linux.json). No general compiler speedup claimed.

## Open ecosystem

The complete buildable compiler, standard library, examples and tests in this repository are Apache-2.0. Optional future private research belongs in separate repositories and is not required to build or use TRIPOL. No proprietary research component is distributed here. See [licensing and governance](docs/open-core.md), [security](SECURITY.md), [contributing](CONTRIBUTING.md) and the [public release checklist](docs/public-release.md).

# Official native guides


---

## Source: docs/language.md

# Language reference — QUANTUM X 0.27.1

## DISCORD REVOLUTION 0.28 development

A named function can be stored as a typed callback. `&format_name` takes a reference to the function. Declare the corresponding parameter or field as `func(string)->string` and call a local callback as a regular function.

```tripol
func format_name(name:string)->string {
    return "Hola, " + name
}

func welcome(render:func(string)->string,name:string)->string {
    return render(name)
}

show welcome(&format_name, "Ada")
```

Callback signatures must match the argument and result types. Pass all arguments when invoking a callback, including arguments whose defaults apply to direct calls. References work in native builds and WebAssembly. The current implementation takes references to named functions with ordinary parameter and result types.

Use `list([&callback_a, &callback_b])` to create a list with a single callback signature. Annotate fields and parameters with `list[func(string)->string]` when they store or receive such a list. A scratch block can call `get`, `size` and `capacity` on a surrounding list. These operations expose immutable values and measurements. Mutations of surrounding lists remain outside scratch blocks.

## QUANTUM X 0.27.1 reference

Source files are `.tri`. Imports must precede declarations and statements; use `import tripol.math`, `import vault.name`, or a quoted project-local `.tri` path. `//` and `#` start comments; newline or semicolon separates statements. Blocks use braces. Strings use JSON-style quoting and escape sequences (no embedded NUL). Identifiers are ASCII; text is byte-oriented. `let` declares mutable bindings, `const` immutable bindings; bare assignment creates or updates a mutable local. No shadowing. Types: signed checked int64 (`int`), `bool`, IEEE binary64 (`float`), immutable `string`, runtime-length `ints`, named `struct`. Integer and float operations require matching types. `float(n)` converts integers; `show` prints floats with six decimals. Struct fields support scalar, string and `ints` types, not nested records. Struct and dynamic array assignments alias references, and const is not deep immutability.

Functions use `func name(a:int,b:int)->int {return a+b}`; omitted types default to int. Functions have no globals or closures. `if/else if/else`, `while`, `for i in range(start,end)`, `for n in array`, `repeat`, `break`, `continue` are supported. `and/or/not` short circuit; `+=` etc. modify bindings. Fixed arrays of up to 400 integers use stack space; runtime-sized `ints(n)` has 0..4194304 elements and cannot grow directly; `tripol.buffer` implements an explicit growable integer backing array. The 401 stack-slot limit includes temporary locals and exception records. Recursion can exhaust native stack. `assert condition,"message"` checks runtime conditions.

## Exceptions

`try { ... } catch code { ... } finally { ... }` (finally optional). Catch receives integer 10 for arithmetic, 11 for bounds, 12 for assertions, or 100..255 from `throw n` (other values normalize to 100). Catches work across function calls and nested blocks. Finally runs after success, recovery, or a recoverable failure in catch; an exception in finally replaces a pending one. Native input, file, socket, float-format, UI and allocator failures remain fatal. Return inside protected blocks is rejected; break/continue may only target loops inside the same protected block. See [errors](errors.md). No exceptions for OS faults or native stack overflow.

## Temporary memory

```tripol
total=0
repeat 1000 {
    scratch {
        values=ints(65535)
        values[0]=7
        total+=values[0]
    }
}
show total
```

The process-lifetime arena reserves 256 MiB of contiguous virtual address space lazily, enabling 4 KiB pages as allocations advance. On Windows this uses MEM_RESERVE/MEM_COMMIT; on Linux mmap(PROT_NONE)/mprotect. The `memory` module exposes used, accessible/committed and capacity bytes (not RSS). Pages already enabled stay enabled after scratch restores the frontier. `scratch` records and restores the allocation frontier (and input scratch pointer) on normal exit and recoverable unwinding. Locals in the region cannot escape. Outer scalar bindings may be read and modified; outer strings are read-only and outer aggregate bindings cannot be accessed. `return` and loop control leaving a protected region are rejected. Eligible loop-local arrays, strings and lists are automatically reclaimed using proven scratch regions, including supported conditional scopes. General automatic memory reclamation remains pending. Scratch does not sandbox I/O or reverse external side effects. Buffers are reused, not securely wiped. Ordinary allocations outside scratch live until exit. Text concatenation and file/network response caps are 65535 bytes; input is 4096 bytes. See library guide.

## Infinity syntax

`f"Hello {name}, score {score}"` interpolates existing scalar bindings of type string, int or bool. Braces may be escaped as `{{` and `}}`; ordinary quoted strings remain unchanged, including JSON strings. Expressions of type int, bool or string are allowed inside braces, including arithmetic, calls, field access and indexing. Float formatting is not supported. Each field is evaluated once, left to right. Interpolation creates immutable strings using the same bounded arena, with the usual 65535-byte result cap. Use `scratch` for repeated temporary strings when safe.

Functions support trailing optional scalar parameters: `func greet(name="Guest", score:int=10)->string {return f"{name}: {score}"}`. Defaults are side-effect-free scalar constants (including checked foldable integer expressions), evaluated as literals at call sites; type inference for omitted parameter annotations is based on the default, otherwise int. Required parameters must precede optional parameters. Struct constructors may omit trailing fields with scalar defaults. No keyword arguments, dynamic defaults, overloads or callbacks.

`tripol ir FILE.tri` validates the project and prints the subset of small pure checked int64 expressions eligible for TIR lowering. This is not a full program IR. Nano may inline small zero/one-parameter functions; turbo can inline up to four parameters; trace/debug retain call boundaries. The inliner evaluates arguments exactly once and keeps overflow checks. The general backend remains unchanged.

## Quantum syntax and growable buffers

`func twice(n) => n * 2` is shorthand for a single returning expression. The default return type is int; annotate other return types, e.g. `func title()->string => "TRIPOL"`. No anonymous functions or struct methods.

```tripol
import tripol.buffer
b = buffer.create()
b = buffer.push(b, 10)
b = buffer.push(b, 20)
show buffer.values(b)
```

`buffer.push` returns a backing array and must be assigned back. It doubles capacity as necessary up to 4,194,303 values, retaining old allocations until exit or safe scratch reclamation. Existing aliases share in-place updates until growth replaces the backing array; this is not an opaque list object. Mutating library helpers require dynamic arrays rather than borrowed stack literals. See library guide.

## Nova collections

`dictionary.create(capacity)` provides a fixed-capacity byte-string-to-int64 map using an integer array, 1..4096 entries, keys 1..32 bytes. `dictionary.set` returns false when full; `get` has an optional fallback. `remove` clears a slot for reuse. This is not general dictionary literal syntax or a hash table. Avoid mutating its raw storage.

CORE preserves the 0.15 syntax and signed integer division semantics (truncation toward zero, remainder sign follows dividend). Positive literal powers of two can compile to shifts and bias correction. Arithmetic overflow checks remain active for potentially overflowing operations. No new public syntax or memory ownership model in this release.

HYPER adds no new syntax. `balanced` and `turbo` preserve checked integer semantics, argument evaluation order and all source validation. Register lowering is limited to pure scalar leaves and call-free straight-line integer range loops; guards and tracing disable loop register promotion.

TITAN adds no new syntax. Eligible call-free integer range loops may have up to four promoted outer accumulators and conditional branches with break/continue. All checked arithmetic and validation remain.

APEX adds no syntax. The optimizer may reuse a total invariant expression and specialize restricted binary recurrences; integer arithmetic remains checked.

OMEGA adds no grammar changes. ArrayCore specializes eligible dynamic integer array reduction loops without changing checked arithmetic or bounds semantics.

ECLIPSE adds the stateful Discord REST API (see [Discord guide](discord.md)) and native array copy/search; no new grammar.

Direct self tail calls with full argument lists reuse the frame in all profiles unless tracing. Other recursion remains subject to the OS stack limit.

## Indexed compound assignment (v0.27)

Mutable fixed and dynamic integer arrays support `a[i] += x`, `-=`, `*=`, `/=`, `%=`. The index is checked and evaluated exactly once before loading the old value, then the RHS is evaluated once. Checked arithmetic and recoverable traps are unchanged. Constant arrays cannot be modified. Dynamic integer lists also support these indexed assignments.

## Growable integer lists (0.27.1)

`nums=list()` or `list(initial_capacity)` creates a stable reference to a growable integer collection. `add(value)` returns new size; `get(index)` and `set(index,value)` check bounds; `remove_at(index)` returns removed value; `clear()` resets size; `size()`, `capacity()`, `len(nums)` and `show nums` are supported. Capacity is at most 1048576, subject to process arena limits. Aliases share one handle across growth; old backing allocations remain in the arena. Direct mutation through a constant binding is rejected; another mutable alias can still mutate the shared collection. This is not general memory reclamation. Lists are distinct from `ints`; methods and checked indexing are supported, but array-only functions do not accept lists.

### Indexed collection iteration and default record fields

```tripol
struct Player {name:string="Guest";score:int=0}
p=Player()
numbers=list()
numbers.add(10)
numbers.add(20)
numbers[0]+=5
for i,value in numbers {show i,value}
```
Default fields must be trailing scalar constants of the declared type. Constructors may omit those fields. Iteration supports fixed arrays, dynamic integer arrays and lists; the bound is captured once on entry and each access remains checked. List indexing and checked compound indexed assignment are now supported. A list index is evaluated once; assignment rechecks the current backing array after the RHS (which may mutate an alias). A shrinking list may cause a bounds trap during iteration rather than an invalid memory access.

## SIMPLE R1

`numbers=list([10,20,30])` initializes a growable integer list. Unlike `numbers=[10,20,30]` (fixed stack array), it can grow with `numbers.add(40)`. The initializer is evaluated left to right once, with up to 400 integer expressions. `list([])` creates an empty list. Existing list(capacity) and fixed arrays retain their behavior. Lists now also infer homogeneous string, bool and float elements from a nonempty initializer. Use `string_list()`, `bool_list()` or `float_list()` for an empty typed list. Mixed lists are rejected. Methods, iteration, display and checked indexing work for these scalar lists. Numeric `sum` and compound indexed assignment remain integer-only.

```tripol
import tripol.errors
numbers=list([10,20,30])
numbers.add(40)
show f"Total items: {len(numbers)}"
try {show numbers[99]} catch code {show errors.message(code)}
```

### Natural dictionaries

```tripol
player = {
    "name": "Alex",
    "health": 100,
    "active": true
}
player["health"] -= 20
show player
```
A dictionary literal is a statically shaped reference with up to 64 unique nonempty literal keys and independently inferred scalar types. It supports checked literal-key lookup, assignment, integer compound assignment, aliases, display and len. Keys and their types cannot change; no dynamic lookup or insertion in this literal type. Use `tripol.dictionary` for runtime string-to-int keys. These are native arena references, not a dynamically typed object model. General reclamation remains limited.

### Simple function inference

A function with a single return expression and no explicit return annotation can infer a scalar return type and omitted parameter types from local expression constraints and calls in the same source unit. For example `func greet(name) => "Hello "+name` and `func double(n) => n*2`. The function is monomorphic; inconsistent calls are rejected by normal type checking. Defaults and explicit signatures are preserved. Ambiguous parameters retain the integer default. Multiple statements, imported call sites and recursive return inference may need explicit annotations.

Integer, boolean and string lists support `contains(value)`, `remove(value)` (first match, returns bool) and `sort()` (in-place, stable, returns size; maximum 4096 items). Strings sort lexicographically by UTF-8 bytes, not locale collation. Float lists retain indexed access, add/remove_at and iteration; value-search and sorting are not included for floats in R1.

Compatibility note: a single-return function with no explicit `->int` can now infer a non-integer return type. Write `->int` if an integer return contract is required, including for exported Vault functions. Explicit mismatches remain compile errors even in unused functions.


## SIMPLE R2

TRIPOL infers a single signature for each function across a project. Calls and typed callback parameters supply context. For example, `apply(twice)` gives `twice` the signature required by `apply`. A function name is a function reference when a local variable does not shadow it. `&twice` remains valid.

```tripol
func apply(action:func(int)->int) { action(21) }
func twice(x) { x*2 }

show apply(twice)
show apply(func(x) { x+1 })
```

An inferred function returns the value of its last expression when it reaches the end. Final `if/else` branches can each yield a value. Explicit `return` remains available. An anonymous function can capture a scalar literal declared with `const`. The compiler copies that value into the function. Mutable captures, runtime values and aggregate captures require explicit arguments. Method overloads use distinct parameter types. The compiler selects one signature from the provided arguments and available defaults. Conflicting evidence requires a type annotation or separate functions.

TRIPOL infiere una firma por función en todo el proyecto. Las llamadas y los parámetros de callback tipados aportan contexto. En el ejemplo, `apply(twice)` aporta a `twice` la firma requerida por `apply`. El nombre de una función actúa como referencia si ninguna variable local lo oculta. `&twice` sigue siendo válido.

Las funciones con retorno inferido devuelven el valor de la última expresión al llegar al final. Las ramas terminales de un `if/else` pueden aportar un valor cada una. También puedes escribir `return`. Una función anónima puede capturar un literal escalar declarado con `const`. El compilador copia ese valor a la función. Pasa como argumentos explícitos las variables mutables, los valores de ejecución y las estructuras. Las sobrecargas de un método se distinguen por los tipos de sus parámetros. El compilador elige una firma a partir de los argumentos y los valores predeterminados. Si hay conflicto, anota el tipo o escribe funciones separadas.


---

## Source: docs/inicio.md

# Primeros pasos con TRIPOL

Instala la distribución Windows x64 proporcionada por el mantenedor (sin firma digital), o compila desde el código fuente siguiendo [build.md](build.md). Abre un CMD nuevo. No necesitas instalar Go, Python o GCC para usar el compilador distribuido.

```cmd
tripol --version
tripol new MiProyecto
cd MiProyecto
tripol run
tripol build
main.exe
```

Edita `main.tri` en tu editor preferido. En CMD, `cd` cambia de carpeta y `cls` limpia la consola.

```tripol
nombre = input("Tu nombre?")
show "Hola", nombre
```

`tripol learn es` muestra el tutorial sin Internet. `tripol libs` lista bibliotecas; `tripol docs math` muestra su código. `tripol check main.tri --lang=es` comprueba errores. `tripol run main.tri --trace` añade una traza de funciones si el programa falla. Las nuevas funciones Windows (formulario y HTTPS) necesitan prueba real en Windows.

### Novedades de Infinity

```tripol
func saludo(nombre="Usuario", puntos=10)->string {
    return f"Hola {nombre}, tienes {puntos} puntos"
}
show saludo()
show saludo("Ana", 25)
```

Prueba `tripol libs`, `tripol docs collections` y `tripol ir main.tri`.

## Primer juego 2D (Windows x64, experimental)

```cmd
tripol new MiJuego2D --template game
cd MiJuego2D
tripol run
```

Usa WASD o las flechas y Esc para salir. El motor dibuja en un arreglo RGB y presenta la imagen mediante GDI. La ejecución de la ventana todavía requiere validación real en Windows. `tripol run examples/quantum_memory.tri` desde el código fuente comprueba la reserva y la reutilización de memoria. `tripol fix main.tri` muestra una sugerencia de cambio sin modificar código.

---

## Source: docs/academy.md

# TRIPOL Academy — SIMPLE R1

Offline beginner course. No account, payment or network required. The compiler includes nine lessons; their exact source examples run in the native regression suite. Start with `tripol learn es list` or `tripol learn en list`, then `tripol learn es 1`. Each lesson displays source, expected output and an exercise. Copy it into `main.tri` and run `tripol run`. The older `tripol learn es` and `tripol learn en` overview commands still work. The lessons do not save progress or automatically grade answers yet.

## 1. Primer programa / First program

```tripol
show "Hola, TRIPOL!"
```

Expected:
```text
Hola, TRIPOL!
```

show muestra texto en la terminal. Cambia el saludo. / show prints text to the terminal. Change the greeting.

## 2. Variables / Variables

```tripol
name = "Alex"
score = 10
show f"{name}: {score * 2}"
```

Expected:
```text
Alex: 20
```

Las variables infieren su tipo. Cambia score y vuelve a ejecutar. / Variables infer their type. Change score and run again.

## 3. Condiciones / Conditions

```tripol
score = 80
if score >= 60 {
    show "Passed"
} else {
    show "Try again"
}
```

Expected:
```text
Passed
```

Cambia score a 40 para probar la otra rama. / Change score to 40 to test the other branch.

## 4. Bucles / Loops

```tripol
total = 0
for n in range(1, 6) {
    total += n
}
show total
```

Expected:
```text
15
```

range excluye el limite final. Prueba sumar del 1 al 10. / range excludes its end. Try summing 1 through 10.

## 5. Funciones / Functions

```tripol
func twice(n) => n * 2
show twice(7)
```

Expected:
```text
14
```

Una funcion puede devolver un resultado. Prueba otro numero. / A function can return a result. Try another number.

## 6. Listas / Lists

```tripol
numbers = list([10, 20, 30])
numbers.add(40)
show numbers
```

Expected:
```text
[10, 20, 30, 40]
```

Una lista puede crecer. Los indices comienzan en cero. / A list can grow. Indices start at zero.

## 7. Estructuras / Records

```tripol
struct Player {
    score: int = 0
}
p = Player()
p.score = 5
show p.score
```

Expected:
```text
5
```

Una estructura agrupa datos relacionados. / A record groups related data.

## 8. Bibliotecas / Libraries

```tripol
import tripol.math
show math.square(7)
```

Expected:
```text
49
```

tripol libs enumera bibliotecas; tripol docs math muestra sus funciones. / tripol libs lists libraries; tripol docs math shows their functions.

## 9. Errores / Errors

```tripol
import tripol.errors
try {
    show 10 / 0
} catch code {
    show errors.message(code)
}
```

Expected:
```text
Integer overflow or division by zero
```

try/catch recupera errores comprobados. No todos los fallos son recuperables. / try/catch handles checked errors. Not all failures are recoverable.

## Test your projects

Write `assert` checks in files ending in `.test.tri` and run `tripol test examples/academy` from the repository root. `tripol test` selects only immediate `*.test.tri` files in the current directory; `tripol test FILE.tri` explicitly selects one file. Every selected program is compiled with the balanced profile, run in a separate temporary working directory and checked for exit status. Ten seconds per program, two minutes total execution budget, at most 64 tests. Output from failed tests is limited to 64 KiB. A caught assertion is not a failure. This is not a sandbox; only run trusted project tests, which have ordinary native process permissions. No package installation or network access is performed by the runner itself.

## Interactive shell

Run `tripol shell` to try scalar expressions, assignments and interpolated text without making a project. Each accepted command compiles to a temporary native executable. Pure scalar assignments persist by deterministic replay; failed commands do not change session source. `:help`, `:code`, `:reset`, `:exit` are available. The shell rejects arbitrary calls, imports, loops, aggregate mutation and IO so replay cannot repeat observable effects. It is a teaching shell, not a full incremental REPL.

---

## Source: docs/libraries.md

# Standard library — 55 modules

Exact source and signatures: `tripol docs NAME`, `tripol libs`, or [stdlib](../stdlib). Counts include seven JSON scanner helper functions; use the documented high-level APIs for stable behavior. Modules are embedded and offline. In nano mode unused routines are eliminated; imports do not imply runtime dependencies.

| Module | Functions | Notes |
|---|---:|---|
| binary | 8 | Validated byte arrays and u16/u32 byte-order conversion |
| arrays | 7 | Integer arrays, copying and reversing |
| console | 3 | Prompts and lines |
| convert | 8 | Unit conversions |
| discord | 34 | Experimental Windows Gateway and REST bot, discovery and protocol helpers; see Discord guide |
| env | 1 | Windows environment access |
| simd | 1 | SSE2 exact int64 array equality |
| encoding | 4 | Hexadecimal encoding, not cryptography |
| errors | 5 | Checked input, division and array access |
| files | 5 | Text reads/writes; 65535-byte max read |
| http | 3 | Basic HTTP/1.0 GET, raw status and body |
| https | 1 | Windows WinHTTP GET, public endpoint only |
| json | 15 | Serialization, syntax validation and raw top-level value extraction |
| math | 12 | Checked integer math and Float64 helpers |
| network | 3 | Blocking TCP IPv4 text exchange |
| random | 3 | Seeded pseudo-random, not crypto |
| search | 4 | Linear and binary search |
| sort | 4 | Integer heapsort, returns copies |
| stats | 6 | Mean, population variance, deviation, etc. |
| system | 5 | Compiled target metadata |
| text | 9 | Byte offsets, prefix/suffix and substring search |
| time | 6 | Monotonic timer and sleep, not calendar dates |
| ui | 7 | Windows modal dialogs and text form |
| validate | 6 | ASCII validation |

`json.valid(s)` checks JSON grammar, nested arrays/objects, escapes and UTF-8 bytes up to 4096 bytes and nesting depth 32. `json.get_raw(s,key)` returns the first matching top-level value as raw JSON text, empty if invalid/missing; keys match the serializer's exact quoted spelling (no escaped-key normalization); key <=512 bytes. It does not deserialize arbitrary typed records or resolve JSON paths. Escaped surrogate code units are grammar-checked, not decoded. Internal helper signatures are subject to change. Serializer limits: quoted text <=512 bytes, integer arrays <=128 items.

`ui.form(title,label,button)` displays a fixed modal form with one edit field, submit and cancel. Empty return on cancellation/empty input; text uses Windows ANSI encoding (ASCII labels required), up to 4096 bytes, may truncate longer input. The existing `ui.message` and `ui.confirm` are also Windows-only. No custom layout/events or arbitrary windows. The project owner demonstrated the form on Windows; the new message-box variants are PE-verified and still need Windows execution testing.

`https.get(host,path)` uses Windows WinHTTP with OS certificate verification, DNS, port 443, 5-second per-operation timeouts, no redirects, only successful 2xx responses, ASCII host/path, body <=65535 bytes without NUL. Not an HTTP general-purpose client: no custom headers, auth, POST, async, total deadline or Linux support; certificate/network errors are fatal. The new path is PE-verified, not executed on Windows in the development environment. `discord.gateway()` requests the public discovery endpoint. OMEGA adds Discord-specific authenticated REST requests and a test bot, not general authenticated HTTPS or Gateway support. See [Discord guide](discord.md).

`network.exchange` is a blocking IPv4 TCP text exchange (half-close, read until EOF, <=65535 bytes); no DNS/UDP/IPv6/TLS here. `http.get_raw` is plaintext, no chunk decoding or redirects; use trusted test services, never send secrets. File operations can overwrite arbitrary accessible paths, are not atomic and do not create directories. Windows file paths are ANSI. No sandbox. For full signatures and examples consult the source directly.

## Infinity additions

| Module | Functions | Notes |
|---|---:|---|
| bits | 6 | Nonnegative int64 bit operations; indices 0..62 |
| checksum | 3 | Adler-32, sum8, sum16; not cryptography |
| collections | 6 | Append, concatenate, insert, remove, slice, equality; return copies, fixed length |
| colors | 8 | Packed RGB24, channels 0..255, no alpha |
| compression | 3 | Integer-array run-length encoding; not general file compression |
| datetime | 7 | Gregorian calendar 1..9999, ordinal and weekday Monday=1; no clock/time zone |
| geometry | 6 | Float64 shapes, distance and interpolation |
| matrix | 6 | Flat row-major integer matrices; explicit shapes, checked math |
| sequences | 4 | Integer sequences, Fibonacci, factorial and triangular numbers |

Collections return independent arrays capped at 4,194,304 elements. They are not amortized growable vectors or dictionaries. RLE uses `[value,count,...]` pairs, may expand data, limits to 32767 runs and 65535 decoded elements, and validates all counts. `checksum` is noncryptographic and unsuitable for authentication. `ui.warning`, `ui.error`, `ui.retry`, and `ui.choose` use native message boxes; choose returns Win32 IDs Yes=6, No=7, Cancel=2. No custom widgets yet; see Quantum canvas and game APIs below.

## Quantum additions

| Module | Functions | Notes |
|---|---:|---|
| animation | 6 | Interpolation, easing and frame selection |
| buffer | 8 | Growable int buffer, explicit reassignment on push; cap 4,194,303 values |
| canvas | 9 | RGB24 framebuffer, clipping, fill, lines, circles and pixel blend; 4,194,304 pixels max |
| events | 7 | Foreground key state (Windows) and common key codes |
| game | 6 | Single Win32 window, GDI presentation, pump, key polling, close, wait |
| image | 6 | Flip, crop, grayscale, small P3 PPM export (<=256 pixels) |
| memory | 4 | Virtual reservation, used and committed/enabled page counts, not RSS |
| path | 6 | Lexical separators/basename/dirname/extension/join; no normalization |
| physics | 5 | AABB, bounds, integration, gravity; not a rigid-body solver |
| vector | 7 | Float64 vector helpers and int array scaling |

For graphics limitations, see [graphics](graphics.md). `buffer.push` must be assigned back; it grows by copying, old backing storage stays allocated. `path` does not resolve symlinks or validate filesystem paths. The game window is Windows-only and PE-inspected, actual execution pending.

## Nova additions

| Module | Functions | Notes |
|---|---:|---|
| camera | 4 | Integer coordinates and camera bounds |
| dictionary | 8 | Fixed-capacity byte-string to int64 map, linear lookup, reusable slots |
| particles | 4 | Caller-owned integer position, velocity and life arrays |
| sprite | 3 | RGB array drawing with clipping and optional transparent color |
| testing | 4 | Checked expectations for scalar values |

Sprite source and destination must not alias; copy first if needed. The map is not a general heterogeneous collection. Testing failures use recoverable assertions. Camera world conversion truncates toward zero. Particles advance by integer ticks and have no autonomous emission.

`time.now_ns()` and `time.elapsed_ns(start)` expose a monotonic timer in nanoseconds, with platform-dependent precision. The origin is unspecified; only differences should be used.

ECLIPSE adds the stateful Discord REST API (see [Discord guide](discord.md)) and native array copy/search; no new grammar.

NEXUS: `search.linear` / `arrays.contains` use AVX2 for large arrays when supported by CPU and OS, scalar otherwise. `discord.run` uses experimental Gateway; `run_rest` preserves old behavior.

## stack / queue (EVOLUTION)

Both modules expose `create(limit=32)`, `capacity`, `size`, `push`, `peek`, `pop`, `clear`. They use bounded integer backing arrays. `push` returns false when full; `peek`/`pop` assert on empty. `clear` resets logical size without freeing arena storage. Do not mutate the internal headers. See `examples/projects/collections.tri`.

## bitset / profiler (OVERDRIVE)

`bitset.create(count)` allocates a compact bounded bitset (63 bits per word, at most 1,048,576 bits). `capacity`, `contains`, `add`, `remove`, `count` expose checked operations. `add` and `remove` report whether membership changed. `profiler.start()` and `profiler.elapsed(started)` measure elapsed monotonic nanoseconds, not CPU sampling.

## QUANTUM X modules

- `heap`: bounded signed integer min-heap, create/capacity/size/push/peek/pop. No allocation after creation during valid operations. Full push returns false; empty peek/pop assert.
- `utf8`: continuation, width_at, valid, count, scalar_at. Byte offsets, strict scalar validation (overlong encodings, surrogates and out-of-range code points rejected). NUL is not representable in current strings. No normalization or grapheme segmentation.
- `csv`: escape a single field, fields/valid for a single comma-separated record with quoted fields and doubled quotes. No file/stream parser yet.

See `examples/projects/quantum_x_libraries.tri` and `memory_stress.tri`.

## 0.27.1 additional modules

- `io`: `read(path,limit=65536)` and `write(path,data)`. Binary int arrays, values 0..255; limit 65536 bytes. Invalid sizes, bytes, oversized input and OS I/O failures are fatal, matching existing file I/O. Writes replace the destination after input validation; not transactional. Reads reserve storage for the requested limit and return actual length.
- `stream`: in-memory byte cursor `create`, `position`, `remaining`, `seek`, `read_byte`, `read_u16` (little endian). Cursor is mutable, checked, and allocates only at creation.
- `regex`: `valid`, `matches`. Bounded byte-oriented full-string matching supporting literal, `.` and `*`; pattern max 128 bytes, input max 4096. No groups, captures, alternation or Unicode matching.

- `atomic`: `load(data,index)`, `exchange(data,index,value)` returns old value, `compare_exchange(data,index,expected,value)` returns bool, `add(data,index,delta)` returns new value. Aligned x86-64 loads and locked RMW instructions; checked index/arithmetic. No general thread API or claim of thread-safe runtime allocation yet.

## QUANTUM X standard-library expansion

- `crypto`: `sha256_bytes(data)` (32 digest bytes), `sha256(data)` and `sha256_text(text)` (lowercase hex). Standard SHA-256, bounded input 65536 bytes, tested against Go's standard implementation. No encryption, KDF or independent security audit.
- `audio`: `wav_mono(samples,rate=22050)`, `square_wave(hz,milliseconds,rate=22050,amplitude=5000)`, `save_wav(path,samples,rate=22050)`, `playback_supported()`, `play_file(path)`. Signed 16-bit mono PCM, 8–48 kHz, at most 32746 samples. Windows synchronous WAV playback via WinMM; other targets return false for playback. Device availability is not guaranteed. No mixer or streaming audio.
- `cli`: `option_name`, `find_option`, `has`, `value`, `integer` operate on an explicitly supplied whitespace-separated option string with `--name` and `--name=value` tokens. Last occurrence wins; no quoting or shell interpretation. Does not read OS argv.
- `process`: `id()`, `platform()`, `exit(code=0)`. Current process ID and explicit termination, checked exit code 0–255. No child spawning.
- `thread`: `supported()`, `sum(data)`, `sleep(milliseconds)`. Windows starts and joins an OS thread for an ordered checked integer reduction; caller waits, callback cannot allocate or mutate the input. Linux reports unsupported for the worker. No general thread callbacks or claims of parallel speedup.

QUANTUM R3 additionally extends game/events/canvas/sprite/particles/audio with OpenGL triangles, XInput, focused mouse coordinates, software tile maps and transformations, reusable particle bursts and drawing, and optional system tones. See [graphics](graphics-r3.md) for contracts and platform limits.


---

## Source: docs/errors.md

# GUARD R2

`tripol check main.tri --lang=es` shows file, line and column where available, error code and hint. `--diagnostics=json` emits a structured diagnostic on stderr, including `suggestion` when an unknown variable/function has one close unambiguous spelling. Suggestions never rewrite files. Type diagnostics may only report statement lines. Reports first compilation failure, not all errors.

`tripol run main.tri --trace` or `tripol build main.tri --trace` instruments native language function frames; uncaught failures display the active call chain, not instruction-level debug symbols or native OS stack traces. Tail-call elimination can remove logical frames; a rethrow reports the active frames at propagation. Tracing adds overhead; use normal nano for benchmarks. `tripol explain TP2001` gives the meaning of a code.

| Code | Meaning | Recoverable via catch? |
|---|---|---|
| TP1000–TP1008 | Compilation / I/O diagnostics | Not applicable |
| TP2001 (10) | Arithmetic | Yes |
| TP2002 (11) | Bounds | Yes |
| TP2003 (12) | Assertion | Yes |
| TP2004 (13) | Input/text/arena | No |
| TP2005 (14) | Files | No |
| TP2006 (15) | Network / HTTPS | No |
| TP2007 (16) | Float formatting | No |
| TP2008 (17) | UI | No |
| TP2009 (100–255) | Uncaught explicit exception | Yes |
| TP9000 | Internal compiler error | Not applicable |

`try { show 10/0 } catch code { show code } finally { show "done" }` prints 10 then done. Catch receives the integer native error code. Finally executes on successful/recoverable paths, not on fatal native failure or OS termination. A failure inside finally propagates. `throw` accepts int; values outside 100–255 normalize to 100. The native runtime still has no general resource-unwinding machinery for I/O failures. See [language](language.md).

`tripol fix FILE.tri` validates a project and previews an unambiguous single identifier spelling correction on the diagnostic line when available. It never edits a source file; ambiguous occurrences and qualified names are not rewritten. JSON diagnostics include `fix_before` / `fix_after` if a preview exists. First-error reporting remains.

---

## Source: docs/build.md

# Build and architecture

Go 1.23+ is required for contributors rebuilding the compiler, not for users of distributed binaries. On a supported x86-64 host:

```sh
go test -count=1 -v ./core ./tests
go vet . ./core ./stdlib ./tests
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o tripol .
mkdir -p installer/payload
cp examples/hello.tri installer/payload/hello.tri
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o installer/payload/tripol.exe .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o TRIPOL-v0.27-Setup-x64.exe ./installer
```

The installer is unsigned, installs under `%LOCALAPPDATA%\Programs\TRIPOL`, uses PowerShell to update user PATH, and supports `--quiet` and `--uninstall`. Close and reopen CMD. Linux x86-64 ELF uses direct syscalls; Windows PE32+ uses OS DLLs for selected functionality. There is no external language runtime in generated executables.

`core/main.go` owns parser, checking and x86 emitter; feature files implement language/runtime helpers; `native/liberty.S` is maintained assembly; `core/liberty_code.go` contains checked-in generated byte templates. Maintainers can regenerate with `python native/generate.py` (requires GNU as/objdump/objcopy); not needed for regular compiler builds. `stdlib` is embedded; `installer` produces Windows setup. `vscode` provides syntax highlighting, not a full debugger. `core/tir.go` implements a typed scalar expression subset used for experimental inlining and `tripol ir` inspection. The experimental turbo profile can inline up to four int64 arguments; general TIR is still pending. `compiler/` contains GENESIS and its lexer/parser demos, not the production compiler. The root `main.go` delegates to `core.Run()`. Go white-box tests stay under `core/`; CLI integration tests live in `tests/`.

CLI: `tripol build FILE.tri --target windows|linux -o OUTPUT --profile nano|balanced|turbo|debug`, `tripol check`, `tripol run`, `tripol bench FILE.tri --runs N`, `tripol bench-build FILE.tri --runs N`. Default source main.tri, default nano; debug disables optional optimizations, not a native debugger; explicit self tail calls reuse their frame unless tracing. Turbo enables additional scalar inlining, not guaranteed speed improvement. `tripol fix FILE.tri` previews a unique spelling replacement without editing. Bench measures process startup and execution, one warm-up excluded, stdout discarded, stdin closed; use noninteractive programs. Bench-build measures source load and compilation excluding startup/write, checks deterministic output. No built-in timeout.

Vault: `tripol add NAME FILE.tri --version X.Y.Z`, `tripol packages`, `tripol verify`, `tripol remove NAME`. Lockfile pins version and SHA-256, cache `.tripol/packages/`; copy both for offline builds. No remote registry, scripts or semver solver; trust package code. See `core/vault.go` and tests for exact constraints.

CORE: `core/core.go` implements the byte scanner and selected checked native instruction lowering. The validation pass still checks every source body before reachability optimization; only the final pass links runtime helpers and serializes the executable. `core/tir.go` caches each scalar lowering for the duration of the optimized build. The optimized range induction increment cannot overflow because the immutable index is strictly below its signed bound.

## HYPER profiles and measurement

`nano` retains the existing compact backend. `balanced` promotes selected call-free integer loops and uses a register ABI for small scalar leaves. `turbo` also performs register-based multi-argument scalar inlining. General control-flow allocation remains pending. All profiles retain source validation and checked arithmetic.

`tripol bench-build FILE.tri --runs 20 --profile turbo` reports per-build averages within calibrated batches, includes deterministic output checking, and excludes CLI startup and disk writes. `tripol bench examples/bench_hyper_internal.tri --profile turbo --runs 20 --internal` reports both user-instrumented region time and process time. The benchmark program must emit exactly one `TRIPOL_NS: N` line with a nonnegative integer. Use noninteractive programs with small output. `time.now_ns()` uses CLOCK_MONOTONIC on Linux and QPC on Windows; units are nanoseconds, not a guarantee of one-nanosecond clock resolution.

TITAN extends register loops with multiple outer bindings and restricted conditionals. `tripol arena --runs 10` uses locally installed Rust, not an embedded or downloaded compiler; see [arena.md](arena.md).

APEX adds a restricted recursive register ABI and proof-driven loop optimizations; `tripol arena --json arena.json` exports raw samples.

OMEGA adds `tripol new MiBot --template discord`, native Windows REST polling and `tripol.simd`. See [Discord guide](discord.md) for security and limitations.

ECLIPSE adds the stateful Discord REST API (see [Discord guide](discord.md)) and native array copy/search; no new grammar.

Direct self tail calls with full argument lists reuse the frame in all profiles unless tracing. Other recursion remains subject to the OS stack limit.

NEXUS adds experimental native WinHTTP WebSocket Gateway and AVX2 search. Gateway requires Windows; offline tests and PE imports validated on Linux. See [Discord guide](discord.md).

ORIGIN adds optional Discord channel filtering, validated binary byte-array codecs and two-vector AVX2 search. See the CI workflow for proposed Windows and Linux checks.

GENESIS R2: [experimental compiler written in TRIPOL](genesis.md). The Go compiler remains required to bootstrap it. It emits Linux ELF from constant expressions or bounded integer assignments and runtime arithmetic.

See [UPDATE R1](update.md) for explicit update checks and downloads.

TIR R6 adds experimental statement-level CFG inspection to `tripol ir`; protected control flow is skipped. The production optimizer remains the existing scalar/register backend. See [GENESIS R3](genesis.md), [UPDATE R2](update.md), and [Windows validation](testing-v0.27.md).

### MEMORY R2 and TIR R6

The compiler proves a narrow subset of straight-line loop bodies cannot retain references beyond the iteration: local integer arrays, scalar arithmetic, local indexing and selected pure builtins. It inserts a scratch region with exception cleanup in every profile. Unproven bodies retain the previous arena behavior. General ownership/GC, string reclamation and arbitrary escape analysis remain pending. `examples/projects/auto_memory.tri` checks reuse. Constant pure conditions can eliminate unreachable branches only after complete source validation. `tripol ir` shows CFG and typed single-block SSA def/use inspection for eligible int64 functions; checked operations are retained even if unused. CFG phi nodes and a general SSA code generator remain pending.

## INTELLISENSE R1 editor commands

`tripol fmt main.tri` prints formatted source without modifying the file. `tripol fmt --stdin` reads source from standard input. Formatting adjusts indentation and verifies unchanged tokens, comments are preserved, and CRLF is retained. VS Code Format Document invokes this command in a trusted workspace.

`tripol lsp` starts a standalone JSON-RPC Language Server Protocol server over standard input/output (no sockets). It accepts full-buffer synchronization and provides syntax/type diagnostics for self-contained or standard-library programs, completion, library signature hover, same-document definitions, document symbols and formatting. Relative local imports can be resolved from currently open file buffers, with project-wide revalidation and cross-file definition/hover/completion. Missing open buffers and Vault imports produce an incomplete-validation warning. No filesystem access or subprocess execution by the LSP server. This is an initial LSP implementation, not full workspace indexing or incremental type analysis.

TIR R6 now lowers eligible multi-statement, call-free checked int64 functions in balanced/turbo to straight-line SSA machine code. Linear scan allocates six volatile value registers and spills excess live values into frame slots; the backend preserves checked dead computations and source evaluation order. Division/remainder, calls and control-flow joins retain the existing backend. No general CFG SSA or phi insertion yet.

MEMORY R2 additionally proves local immutable strings and growable integer list operations non-escaping inside eligible straight-line loop bodies. Outer strings may be read but not reassigned; outer aggregate aliases, arbitrary functions, branches and control transfer remain outside this proof. `examples/projects/memory_temporaries.tri` exercises one million string iterations and 100,000 list iterations. TIR R6 now reuses dominated identical pure scalar SSA expressions and propagates checked constant expressions; an overflow/divide-by-zero expression is not folded away. No general GC or cross-block SSA yet.

TIR R6 extends the native scalar SSA backend to signed checked division and remainder. Functions containing either operation reserve RCX/RDX as scratch and use four value registers plus spills; other functions retain six. Positive power-of-two constant divisors use the existing signed truncation-toward-zero optimization. Zero divisors and MinInt64/-1 continue to trap, including unused checked results.

INTELLISENSE R1 open-buffer mode resolves transitive relative `.tri` imports confined to the entry file's directory, checks cycles and imported top-level statements, and revalidates open documents when a buffer opens, changes, or closes. Navigation reads only currently open documents. This is not disk-wide workspace indexing; closed files and Vault packages still require normal `tripol check`.

MEMORY R2 now extends the region proof across if/else scopes: each branch is analyzed independently and new local names disappear at the join. Eligible conditional loop bodies reclaim all temporary allocations at iteration exit, including on recoverable traps. A possible outer reference assignment in either branch disables automatic reclamation for the whole loop. Arbitrary nested loops, calls and cross-region control transfer remain conservative fallbacks.

### TIR R6 compilation efficiency

Optimized profiles validate every body with the original checked emitter, but reuse the validation machine-code, relocation and literal buffers between bodies rather than retaining a throwaway image. The standalone LSP uses the same bounded validation path. Reachability skips only branches proven constant by the same fold used by optimized emission; all branches are still type-checked. Debug keeps its original full-image behavior and skips unnecessary reachability traversal. SSA CSE uses typed value keys, avoiding string formatting, and linear scan recycles spill slots after their last use. Measure the validation path with `go test ./core -run '^$' -bench BenchmarkValidationBuffer -benchmem -count=5`. This benchmark compares the old accumulated validation image with the new reused buffers, not total compiler peak RSS. No numerical speedup is asserted before measurement.

Use `tripol compile-report main.tri --profile turbo --json` for measured load/parse, validation, reachability, emission and linking nanoseconds; also reports selected functions, image size, SHA-256 and validation buffer capacity. It produces exactly the bytes the normal compiler would produce but does not write or execute the image. Timings exclude formatting and hashing, and buffer capacity is not process peak RAM. Debug omits the separate validation pass. The benchmark workflow publishes stage measurements and the old/new validation microbenchmark alongside existing benchmark artifacts.

## COMPILER REVOLUTION — CFG constants

Balanced and turbo run a post-validation forward scalar constant pass over structured control flow. At an if/else join, every branch that can continue must agree on the value. A proven condition selects its branch, and an edge that returns, breaks or continues contributes no value to the join. The pass tracks int/bool values already known at entry to the conditional, keeps the original AST for validation and debug, and treats scratch and try as barriers. Reports include optimization time, scalar substitutions, agreed joins and constant branches. General multi-block SSA and phi-node register allocation remain experimental goals.

COMPILER REVOLUTION also emits checked x86-64 immediate arithmetic for signed 8/32-bit constants in scalar SSA. Constants used exclusively as encodable right-hand operands of addition, subtraction or multiplication no longer consume value registers or spill slots. Wide constants, return values and division operands keep the existing allocation. Every arithmetic instruction still branches on signed overflow.

The CFG pass retains known scalar facts across a loop only when no statement in its body writes that variable. Range bounds may use entry constants because they are evaluated once; while conditions use invariant constants only. The body is analyzed with invariant entry facts and its outgoing facts are discarded at the backedge. Protected/opaque loop bodies disable this optimization. No trapping expression is hoisted or speculatively evaluated. `compile-report` includes an invariant-fact counter; this is not a count of dynamically saved instructions.

### TIR PARALLEL R1

Structured CFG optimization uses independent function-body workers for units with at least 64 declarations, bounded to four workers and GOMAXPROCS by default. Set `TRIPOL_JOBS=1` for serial operation or `TRIPOL_JOBS=2`, `3`, `4` for a bounded override. Smaller units remain serial. Validation, native emission and linking remain sequential. No concurrent modification of input AST, no shared mutable fact map, and stable output order; tests compare both AST and native PE/ELF bytes between one and four workers. `go test ./core -run '^$' -bench BenchmarkCFGWorkers -benchmem -count=3` measures the optimization pass, not overall compilation speed. Parallelism can be slower for small workloads and allocate more memory, so no speedup is assumed.

### INCREMENTAL R1 (experimental)

`tripol build main.tri --profile turbo --cache` enables parsed-module and native image caches. Local and embedded modules reuse independent parsed units by source bytes and compiler executable SHA-256. An edit invalidates the affected unit. Each build still resolves imports, checks Vault integrity and reads source bytes. Debug and traced builds bypass both caches. The image cache also keys the complete resolved AST, profile and target. When that key changes, semantic validation, optimization, emission and linking run for the complete project.

Module entries are limited to 8 MiB apiece under `TRIPOL/modules-r1`. Image entries are limited to 64 MiB apiece under `TRIPOL/compiler-r1`. Both locations live under the OS user cache directory. Entries are checksummed and written through temporary files. Use a cache directory controlled by your user account. Remove either directory to reclaim space.

### TIR LINK R1 and CFG R2

The native linker reserves the literal tail in one capacity operation, preserving byte layout and first-use order. Relative relocations are checked for valid write offsets and signed 32-bit displacement before patching. `BenchmarkLinkLiterals` compares incremental growth with preallocation. Balanced and turbo additionally remove unreachable statement tails after direct return/break/continue, and use validated scalar equality/boolean conditions to refine branch-local facts. Validation still runs before any pruning, and checked return expressions remain intact. These are not a general SSA or PGO backend. No new speedup is claimed before measurement.

### VECTOR R2, SSA R2 and PGO R1

Balanced and turbo recognize a zero-based integer array fill with a constant value. Balanced uses SSE2. Turbo checks CPUID and XGETBV before using AVX2 on fills of at least 32 elements. Short fills and CPUs without OS-enabled AVX2 use SSE2 and a scalar tail. Oversized fills keep the original bounds trap. Protected scopes and other loop bodies use the ordinary emitter.

Turbo lowers eligible call-free integer and boolean functions to native multi-block SSA. It supports conditionals, short-circuit logic, loops, break, continue and returns. The register allocator uses live-in and live-out sets across edges, reuses spill slots, and applies parallel copies on phi edges. Checked computations keep their original evaluation order. Calls, aggregate values and protected regions use the established native backend. `compile-report` counts functions lowered through multi-block SSA.

PGO R1 records physical function entries during a separate checked, unoptimized training execution on the local x86-64 host, then sorts reachable native function bodies in descending observed frequency, keeping equal-frequency source order. It does not change semantics, inline hot calls, or profile individual branches. Tail-restart jumps are not counted. Training executes the source program with the current user's permissions, no stdin, discarded stdout, a 30-second timeout and up to 8 MiB of stderr observations; use only trusted noninteractive programs. Existing profile files are not overwritten. Profiles are matched to the exact compiler executable, expanded parsed project and target. Rebuild/retrain if any changes. Profile data is a layout hint, not trusted executable code. No speedup is guaranteed; compare against ordinary turbo.

```sh
tripol pgo-train examples/bench_compiler_revolution.tri -o profile.json
tripol build examples/bench_compiler_revolution.tri --profile turbo --pgo profile.json -o optimized
```

The CI matrix checks ordinary debug, turbo and PGO output equality and cache byte-equivalence on native Windows/Linux. These are proposed automated acceptance gates until the corresponding CI runs pass. COMPILER REVOLUTION R1 is scoped to the implemented conservative optimization subsets; general vectorization, full CFG SSA, function-level incremental builds and richer PGO remain research work.

---

## Source: docs/testing-v0.27.md

# QUANTUM X — Windows validation

Run only after the v0.27.1 release is published. Save v0.26 results first, on the same PC and with the same profile, workload and power settings:

```cmd
tripol --version
tripol arena --runs 20 --profile turbo --tripol-only --json arena-v026.json
tripol update --check
tripol update
```

UPDATE R1 downloads and verifies the official installer. Close TRIPOL processes, run the verified unsigned installer from the printed path and reopen the terminal. Do not install if verification failed.

## Installed compiler

```cmd
where tripol
tripol --version
tripol doctor
tripol new Test027
cd Test027
tripol check
tripol run
tripol build
main.exe
```

Expected: TRIPOL 0.27.1 and Hello, TRIPOL! from the generated program. Check all four profiles: nano, balanced, turbo, debug. Each must give equivalent results.

## Validation Center

From a clone of the repository, in PowerShell:

```powershell
powershell -NoProfile -File .\tests\manual\test-v027.ps1
```

It verifies the installed version, native compiler execution, all four profiles, compound array operations, recoverable arithmetic/bounds failures, new libraries and invalid source rejection. It writes `tripol-validation-v027.json` in the invocation directory. Add `-Benchmark` for the full Arena run, or `-CheckUpdates` for a live network check. It never installs updates. If script execution is restricted, review the source and use a process-scoped execution policy if permitted, or run the commands manually.

## GENESIS R3

```cmd
tripol build compiler/file.tri --target windows -o genesis-stage.exe
tripol genesis genesis-stage.exe examples/projects/genesis_r3.tri -o genesis-demo
```

Expected: an ELF64 Linux x86-64 executable. On a Linux x86-64 system run `chmod +x genesis-demo && ./genesis-demo`; expected output 30. The stage is written in TRIPOL; Go CLI transports hexadecimal bytes to disk (no Python). Linux executes the emitted ELF. A separate native PE stage now produces Windows executables directly. Full self-hosting remains pending.

```cmd
tripol build compiler/windows.tri --target windows -o genesis-win.exe
genesis-win.exe
```

Enter `examples/projects/genesis_control.tri` and `genesis-control.exe` on separate lines, using a new output path. Then run `genesis-control.exe`. Expected four lines: 4, 42, 2, 1. This exercises forward calls, mutual recursion, else-if, and short-circuit evaluation.

## Libraries and memory

```cmd
tripol run examples/projects/quantum_x_libraries.tri
tripol run examples/projects/memory_stress.tri
```

Library output: 10 20 30, then 8, then 3. The stress example performs one million push/pop pairs on a bounded heap and verifies unchanged used arena bytes. Arena bytes are not total process RSS; use Task Manager separately. Also run `tripol run examples/projects/memory_temporaries.tri` to check one million string iterations and 100,000 list iterations. Expected equality checks are true. These examples demonstrate scoped automatic reclamation, not general garbage collection.

## Graphics

```cmd
tripol build examples/projects/quantum_x/main.tri -o quantum-x.exe
quantum-x.exe
```

Space/click/XInput A starts. Move with WASD/arrows or the left stick; collect golden lights; M toggles optional system tones; ESC exits. Check camera, particles, mouse, controls, stability and average paced FPS. The main demo uses CPU software rendering; the separate OpenGL demo below is experimental.

## Editor

Install the official v0.27.1 VSIX through VS Code's Extensions: Install from VSIX command. Open a .tri file, try import completion, function completion/hover and local function definition lookup. Save, then run TRIPOL: Check source. Try invalid code and inspect Problems. Explicit Build/Run commands require trusted workspace; check-on-save is optional and off by default. Use Format Document (or `tripol fmt main.tri`) and try the experimental `tripol lsp` stdio server. The experimental standalone LSP resolves imports across open buffers; test an unsaved change to an imported function and Go to Definition. Full disk-wide indexing is still pending.

## Updates and offline

Run `tripol update --check`, `tripol update --history`, and (when a newer release exists) `tripol update --download-only`. Disconnect the network and build a local program: compilation must remain independent of the network. Explicit online checks may fail cleanly. For supported per-user installations with a newer official compiler asset, UPDATE R2 offers `tripol update --install` and `tripol update --rollback` with confirmation. Older updaters and portable/system installations remain manual. Test rollback in an isolated installation, not on your only working compiler. Use local Go HTTP tests to exercise invalid metadata, tampered assets and interrupted downloads; do not publish fake releases.

## Repeat measurements and source tests

```cmd
tripol arena --runs 20 --profile turbo --tripol-only --json arena-v027.json
go test -count=1 -v ./core ./tests
go vet . ./core ./stdlib ./tests
go build -o tripol.exe .
tripol.exe --version
```

Compare medians and p95 for the same case names in both JSON files. Record hardware, versions, profile, output equivalence and sample counts. No universal speed claims. Keep any failure logs; do not disable tests.

For the 0.27.1 additions: `tripol run examples/projects/dynamic_list.tri`, `tripol run examples/projects/binary_io.tri` from a disposable working directory (writes quantum-bytes.bin), and `tripol build compiler/direct.tri -o genesis-stage.exe`. See [GENESIS](genesis.md) for the direct driver. Retain prior version reports for comparison.

### Experimental OpenGL and XInput

See [QUANTUM R3 graphics](graphics-r3.md). Compile the OpenGL demo, resize the window and press ESC. Record whether OpenGL succeeds, frames and elapsed time. The fallback is a reopened GDI window. Test physical gamepads separately; CI only checks imports, no-window behavior and validation.

---

## Source: docs/arena.md

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

QUANTUM X retains Arena's existing comparable workloads and JSON measurements. Capture v0.26 and v0.27 on the same host using `--runs 20 --profile turbo --tripol-only` with distinct `--json` paths. No new cross-language backends or measured improvements are claimed in this first implementation.

## ARENA R8 companion benchmarks

On Linux, `python benchmarks/compare_versions.py --old /path/to/v026 --new /path/to/v027 --runs 7 --out comparison.json` compares the same existing Arena workloads. `python benchmarks/quantum_x.py --compiler /path/to/v027 --runs 7 --out languages.json` measures checked Fibonacci(30) against available C, C++, Rust, Go, Zig and Python tools, verifies every result, records build/runtime medians and p95, native executable sizes, and a separate GNU time RSS probe when available. Missing or incompatible tools are explicitly recorded, not assigned synthetic results. The dedicated QUANTUM X benchmarks workflow uploads both reports. Hosted runners vary, so do not extrapolate to other hardware or workloads.

---

## Source: docs/graphics.md

# QUANTUM 2D R2

`tripol.canvas` creates software RGB24 framebuffers using dynamic integer arrays (row-major, no alpha). Dimensions must be 1..2048 each, product <=4194304. For example, 320×180 is 57,600 pixels. Use `canvas.clear`, `canvas.rectangle` (clipped), `canvas.pixel`, and `canvas.get`. NOVA adds native checked horizontal span fills, `canvas.line`, `canvas.circle`, and `canvas.blend` (constant-opacity single pixel). Mutable operations accept dynamic arrays; do not mutate raw shape metadata or pass a const alias if immutability matters. `image.flip_x`, `flip_y`, `crop`, and `grayscale` produce copies. PPM export is limited to 256 pixels. These operations are available on Windows and Linux.

`tripol.game` is Windows x64 only, one window and one native staging buffer. `game.window(title,w,h)` opens/replaces a window (ASCII title, window size 64..2048 each). `game.running()` pumps Win32 messages and tests window lifetime; `game.key(vk)` reads a virtual-key value only when the window is foreground; `game.present(pixels,w,h)` validates RGB and shape then scales the framebuffer into the client area with GDI StretchDIBits. `game.close()` destroys the window and releases its staging buffer. `game.wait(ms)` uses the existing native sleep routine; no vsync. `events` exposes arrow, space, escape codes. Escape in the demo exits. Closing via X also exits the message loop. Call running and present regularly; no retained drawing on WM_PAINT, no dedicated render thread. Resize can cause temporary repaint/flicker. RGB-array sprites support clipped drawing and a transparent color key. There is no PNG decoding, texture filtering, text/font rendering, mouse events, audio, GPU acceleration, general window controls or FPS guarantee. The owner reported the previous Windows demo working. The NOVA staging buffer is 16 MiB and uses direct GDI presentation; the new path was cross-built and PE-inspected, not executed on Windows.

```tripol
import tripol.game
import tripol.canvas
p=canvas.create(320,180)
game.window("QUANTUM",960,540)
while game.running() {
    if game.key(27) {break}
    canvas.clear(p,1052688)
    canvas.rectangle(p,320,180,140,70,40,40,16777215)
    game.present(p,320,180)
    game.wait(16)
}
game.close()
```

Use `tripol new Demo --template game` for a 720p moving-sprite demo, or `examples/quantum_game.tri`. Use `tripol build examples/quantum_game.tri --target windows -o quantum.exe` from the source tree. Generated executable calls USER32/GDI32/KERNEL32 directly; no Go runtime in user code. Platform-specific errors exit 17, reported as TP2008 by `tripol run`.

---

## Source: docs/graphics-r3.md

# QUANTUM R3 — experimental graphics backend
The existing Windows GDI framebuffer renderer is unchanged. A separate native WGL/OpenGL 1.1 path accepts an integer array of triangles: each triangle has nine elements (x, y, RGB per vertex), coordinates -1000..1000, colors 0..16777215, maximum 10000 triangles. It renders interpolated vertex colors, clears the background and swaps buffers. It retains the window's WGL context and releases it when game.close runs. GL failure returns false. Do not mix GDI and OpenGL in the same window: close/reopen the window before switching.

Windows OpenGL may be backed by hardware or Microsoft software rendering. No claim of actual hardware acceleration or measured FPS is made. This backend does not yet provide GPU textures, shaders or batched modern GPU commands. Linux returns false; no new runtime dependencies beyond Windows system DLLs.

```cmd
tripol build examples/projects/quantum_x/opengl.tri -o quantum-gl.exe
quantum-gl.exe
```

Expected: a colored triangle, or a cyan square on the software fallback window. Escape closes. Test window resize, repeat close/reopen and report frames and elapsed ns. The current demo uses a 16 ms sleep, so it measures paced presentation, not maximum throughput.

XInput 1.4 (Windows 8+) support is available via events.pad_connected, pad_buttons, pad_axis and pad_trigger. Controllers 0–3; axes 0–3 (left X/Y, right X/Y, -32768..32767), triggers 0–1 (0..255), buttons are the native XInput mask. -1 means unavailable, which overlaps a valid axis value; check connection separately. No vibration. Linux returns disconnected. Device availability and physical controller input require manual testing.

The software API now includes clipped palette tile maps (`canvas.tile_map`), clockwise 90° rotation and nearest-neighbor resizing (`sprite.rotate90`, `sprite.scale_nearest`), clipped scaled drawing with color key and uniform opacity (`sprite.draw_scaled`), pooled particle bursts and drawing, and game-focused logical mouse coordinates (`events.mouse_x`, `mouse_y`, `mouse_left`). `audio.beep` uses a synchronous Windows system tone without writing files. The updated main demo has a scrolling map, start menu, camera, XInput movement, optional audio and paced FPS report. GPU textures, shaders, general scene graph and asynchronous audio remain pending.

---

## Source: docs/web-engine.md

# TRIPOL WEB ENGINE R3 — compatibility contract

The web backend reuses the native parser and x86 validation pass, then emits WebAssembly with an explicit browser host ABI. It is not yet a universal shared IR or a full replacement for the native backend. The existing educational interpreter remains available.

## Implemented in source (subject to CI validation)
- int64, bool, string and float scalar lowering; string interpolation, concatenation and equality.
- Scalar functions, defaults, recursion, control flow.
- Fixed and dynamic integer arrays, indexed reads and writes, iteration, sum and length.
- Typed growable lists, collection methods and indexed access.
- Record constructors, defaults, typed field reads/writes; shaped dictionaries with literal-key reads/writes and len.
- Console input; selected pure math functions.
- Bounded browser host handles; checked integer arithmetic, array bounds, step and output limits.

## Remaining parity work
- Exceptions: native try/catch/finally/throw and recoverable faults across calls.
- Scratch memory, allocation frontier, native memory metrics, and identical lifetime behavior.
- Complete aggregate display, full library intrinsic coverage, precise float/UTF-8 behavior, and alias/lifetime edge cases.
- Virtual project files, package imports, virtual filesystem, permitted browser networking, graphics and audio adapters.
- Shared whole-program typed IR; native-equivalent optimization, SIMD and incremental cache.
- Windows differential testing and full platform-specific capability matrix.

## Verification
The workflow builds the Go compiler for wasm, runs actual generated WASM modules, and compares the 100 Playground sources to Linux native executables. The parity-report.json records matches, differences, and rejections. Required smoke tests and all 100 native differential cases gate subsequent publication. Earlier R3 builds may have a partial parity report. No claim of full language compatibility is made. The browser discovers the highest verified ABI and preserves explicit engine selection; it does not silently replace a WASM failure with the interpreter.

Generated modules require the matching host ABI. The browser is sandboxed; native OS access, arbitrary filesystem paths and raw sockets cannot be provided directly. Unsupported features must produce explicit diagnostics rather than fabricated output.

## WEB LIBRARIES — portable expansion
The R3 ABI now has source implementations for byte-oriented strings, byte slicing and search, integer conversion, input, array equality/fill/range fill, checked u32 operations and monotonic clocks. Library imports use the original embedded TRIPOL sources rather than JavaScript reimplementations of library algorithms.

The browser import allowlist includes 50 candidate modules: math, arrays, collections, datetime, bits, bitset, binary, buffer, checksum, colors, compression, console, convert, dictionary, encoding, geometry, json, matrix, path, physics, random, search, sequences, simd, sort, stats, text, validate, vector, stack, queue, animation, camera, particles, canvas, sprite, testing, errors, network, http, image, time, profiler, heap, utf8, csv, stream, regex, cli, crypto.

Allowlisted is not equivalent to fully compatible. Required smoke tests for datetime, text and JSON run before publication. A separate library-parity-report.json compares one real program per candidate module against native Linux; passing a probe does not certify every function. Browser limits remain intentionally smaller than native limits. Native-only operations such as TCP exchange, image file export, sleep and native device access may still be rejected even if their containing module is allowed. Unsupported modules require dedicated capability adapters. No simulated successful output.

## WEB ADAPTERS R1 (source implementation, CI pending)
- files / io: per-execution in-memory virtual filesystem (64 files, 64 KiB/file, 512 KiB total), relative paths only, no host filesystem access. Text and binary I/O reuse standard library wrappers. Files are returned on successful execution and downloadable by explicit user action in the Playground and Docs Live. No persistence across runs, no automatic uploads.
- system: browser target metadata (web / wasm32 / actual compiler version). is_windows and is_linux return false.
- memory: logical reserved collection slots and interned byte payloads, not actual browser process memory. Logical capacity 6 MiB; used and committed report the same logical charged bytes. Allocation caps still apply separately.
- atomic: load, exchange, add and compare_exchange within a single isolated execution; no shared memory/thread-safety claim.
- audio: portable PCM/WAV generation and virtual export only. Native playback/tone requests reject explicitly.
- time/thread: bounded sleep (<=1000 ms) requires an isolated worker with SharedArrayBuffer and Atomics.wait; no native worker sum or general parallel thread implementation.
- Remaining adapters: live game/window/input, native dialogs, browser-authorized HTTPS, Discord server-side credentials, shared-memory atomics and parallel workers. No raw TCP or secret-bearing Discord bot in a public static page. Browser-specific semantics must be tested separately from native parity.

## WEB GRAPHICS PREVIEW R1 (experimental, CI gated)
- The real browser compiler lowers game.window, game.running, game.present, game.close and game.gl_frame to host operations 82–86. Framebuffer values and triangle vertices originate in generated WASM, not a source-code simulation.
- Preview is bounded to 512×512 pixels, 120 frames, 16 MiB of frame data and 1000 triangles per frame. During an open preview window, sleep_ms records frame playback delay (up to 1000 ms), rather than blocking. Frames are replayed after execution; this is not a real-time interactive game loop. Use finite loops and explicit close. No native window or device access is exposed.
- Keyboard/mouse/controller polling deliberately fails with a diagnostic until a live input protocol is implemented. Graphics operations outside the supported contract may fail. WebGPU presentation is optional; a software Canvas 2D renderer handles devices without WebGPU. Native behavior is unchanged.
- Required smoke checks compile and execute real pixel and triangle programs and verify exact frame payloads. Existing 100-example native differential tests remain mandatory.

---

## Source: docs/discord.md

# DISCORD REVOLUTION — TRIPOL 0.28 experimental release

## SIMPLE R2 bots / Bots con SIMPLE R2

Use `bot.slash("info", "Informacion del bot", info)` to infer the signature of `func info(ctx)`. End the callback with `ctx.reply(value)` to return the payload. The accepted values are strings, a single embed, an embed batch and Components V2. Only a returned payload or an explicit `ctx.respond(payload)` sends the initial response. An earlier expression whose result is discarded builds a payload without sending it.

Usa `bot.slash("info", "Informacion del bot", info)` para inferir la firma de `func info(ctx)`. Termina el callback con `ctx.reply(valor)` para devolver el payload. Puedes pasar texto, un embed, un lote de embeds o un panel de Components V2. El envío inicial ocurre al devolver el payload o llamar explícitamente a `ctx.respond(payload)`. Una expresión anterior cuyo resultado se descarta solo construye el payload.

- [Bot con un embed / Embed bot](../examples/discord_simple_embeds.tri)
- [Bot con panel y botón inline / Components V2 bot](../examples/discord_simple_components.tri)

`discord.container("#FFFFFF")` wraps the layout in an accented container. Use button styles `"primary"`, `"secondary"`, `"success"` and `"danger"` or their numeric equivalents. `panel.premium(sku_id)` adds a purchase button for an application SKU. Discord handles the purchase. Use the typed button callback for ordinary buttons.

`discord.container("#FFFFFF")` envuelve el panel en un contenedor con color. Usa los estilos `"primary"`, `"secondary"`, `"success"` y `"danger"` o sus equivalentes numéricos. `panel.premium(sku_id)` añade un botón de compra de un SKU de la aplicación. Discord gestiona la compra. Los botones comunes admiten los callbacks tipados.

## Run an offline example / Prueba un ejemplo sin conexión

```tripol
import tripol.discord

bot = discord.bot()
bot.command("hola", "¡Hola, compañero!")
show bot.preview("hola")

embed = discord.embed()
embed.title("Hola")
embed.description("Creado con TRIPOL")
embed.color("#00A0FF")
embed.field("Versión", "0.28.0", true)
show embed.payload()
```

Run `tripol run examples/discord_new.tri` on Linux or Windows, or open the Discord library examples in the Playground. Neither offline operation uses credentials. On a new bot the prefix is empty. Set `bot.prefix("!")` to match `!hola`. Duplicate command names replace their previous response.

Ejecuta `tripol run examples/discord_new.tri` en Linux o Windows, o abre los ejemplos de la biblioteca Discord en el Playground. Ambas operaciones funcionan sin credenciales. El prefijo inicial es vacío. Usa `bot.prefix("!")` para responder a `!hola`. Un comando repetido sustituye la respuesta anterior.

## Permissions and decimal choices / Permisos y elecciones decimales

Specify a permission name such as `"manage_messages"` when registering a command. Use `ctx.can("manage_messages")` in the callback to check the effective member permissions Discord supplied, including Administrator. Numeric bitsets still work for registration. This check does not inspect a target’s role hierarchy.

Puedes indicar un permiso como `"manage_messages"` al registrar un comando. Usa `ctx.can("manage_messages")` en el callback para comprobar los permisos efectivos del miembro informados por Discord, incluido Administrador. El registro también admite máscaras numéricas. Para actuar sobre un objetivo, la jerarquía de sus roles se comprueba por separado.

Build fixed decimal choices with `choices.number("Half", "0.5")`. The decimal text becomes a JSON number. `ctx.target_username()` and `ctx.target_content()` read the resolved selection in context menus. `server.move_role(role_id, position)` reorders a role above @everyone.

Añade elecciones decimales con `choices.number("Mitad", "0.5")`. El texto decimal se convierte en un número JSON. `ctx.target_username()` y `ctx.target_content()` leen la selección de un menú contextual. `server.move_role(role_id, position)` reordena un rol por encima de @everyone.

## Connect / Conecta

Set `DISCORD_TOKEN` before starting a bot in an authorized private server. For text commands and message callbacks, enable Message Content Intent. Use `tripol new MiBot --template discord` and `tripol run` on Windows or Linux x86-64.

Configura `DISCORD_TOKEN` antes de iniciar el bot en un servidor privado autorizado. Para comandos de texto y callbacks de mensajes, activa Message Content Intent. Usa `tripol new MiBot --template discord` y `tripol run` en Windows o Linux x86-64.

On Linux x86-64, `tripol build` packages the native program together with the TRIPOL launcher when it uses Discord networking. The launcher extracts the inner ELF into a private temporary directory and uses two inherited pipes for transport. It uses system trust roots for TLS. A raw Linux ELF cross-compiled from Windows can be launched on Linux with `tripol discord-run FILE`. A Linux bundle includes the compiler launcher, so it is larger than the inner ELF.

En Linux x86-64, `tripol build` empaqueta el programa nativo con el lanzador de TRIPOL si usa la red de Discord. El lanzador extrae el ELF a un directorio temporal privado y usa dos tuberías heredadas para el transporte. TLS usa los certificados de confianza del sistema. Un ELF generado desde Windows se ejecuta en Linux con `tripol discord-run FILE`. El paquete de Linux incluye el lanzador del compilador y ocupa más que el ELF interior.

On Windows x64, `tripol run` and `tripol discord-run FILE.exe` use the same TLS launcher as Linux. It negotiates zlib-stream, resumes sessions, checks heartbeat acknowledgements, reuses HTTPS connections and handles global and per-bucket rate limits. The directly launched PE uses the standalone WinHTTP fallback. See the [live validation guide](discord-validation.md).

En Windows x64, `tripol run` y `tripol discord-run FILE.exe` usan el mismo lanzador TLS que Linux. Negocia zlib-stream, reanuda sesiones, verifica las confirmaciones de heartbeat, reutiliza HTTPS y controla los límites globales y por bucket. El PE iniciado directamente usa WinHTTP. Consulta la [guía de validación conectada](discord-validation.md).

## Message callbacks / Callbacks de mensajes

```tripol
import tripol.discord

func answer(message:discord.Message)->string {
    if message.matches("hola") {
        return message.reply("¡Hola!")
    }
    return ""
}

bot = discord.bot()
bot.on_message(&answer)
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"hola\",\"author\":{\"id\":\"456\"}}}")
```

Use `bot.start()` for the Gateway. Fixed text commands have precedence over the callback. Returning an empty string skips the reply. The context exposes `author_id`, `channel_id`, `content_json` and `event_json`. Bot and webhook messages are ignored.

Usa `bot.start()` para conectarte al Gateway. Los comandos de texto fijos tienen prioridad. Devuelve un texto vacío para omitir la respuesta. El contexto expone `author_id`, `channel_id`, `content_json` y `event_json`. Se omiten los mensajes de bots y webhooks.

## Gateway events / Eventos del Gateway

```tripol
import tripol.discord

func joined(event:discord.GatewayEvent)->bool {
    show "Nuevo miembro", event.guild_id, event.user_id
    return true
}

registry = discord.events()
registry.intent("GUILD_MEMBERS")
registry.on("member_join", &joined)
show registry.preview("{\"op\":0,\"t\":\"GUILD_MEMBER_ADD\",\"s\":9,\"d\":{\"guild_id\":\"321\",\"user\":{\"id\":\"567\"}}}")
```

Call `bot.start_with_events(registry)` to connect with the same callbacks while keeping existing text and slash commands. Each event uses a `func(discord.GatewayEvent)->bool` handler. The context exposes `name`, `sequence`, `guild_id`, `channel_id`, `user_id` and `data_json`. Use `event.field(name)` to read any event property as JSON. Callbacks execute in arrival order and `registry.preview(event_json)` returns the handler's boolean result for testing.

Usa `bot.start_with_events(registry)` para conectar los mismos callbacks y conservar los comandos de texto y slash. Cada evento usa un manejador `func(discord.GatewayEvent)->bool`. El contexto expone `name`, `sequence`, `guild_id`, `channel_id`, `user_id` y `data_json`. Usa `event.field(nombre)` para leer cualquier propiedad como JSON. Los callbacks se ejecutan en orden de llegada y `registry.preview(event_json)` devuelve el booleano del manejador durante las pruebas.

Register official dispatch names such as `READY`, `GUILD_ROLE_UPDATE`, `MESSAGE_DELETE`, `MESSAGE_REACTION_ADD`, `VOICE_STATE_UPDATE` and `RESUMED`. The three member aliases are `member_join`, `member_leave` and `member_update`. The registry requests `GUILDS` by default. Use `registry.intent("GUILD_MESSAGE_REACTIONS")` for guild reactions and enable `GUILD_MEMBERS` or `GUILD_PRESENCES` in the registry and Developer Portal before registering privileged callbacks. The dedicated slash, button, select and modal methods handle interactions. Voice server credentials remain reserved for voice transport.

Registra nombres oficiales como `READY`, `GUILD_ROLE_UPDATE`, `MESSAGE_DELETE`, `MESSAGE_REACTION_ADD`, `VOICE_STATE_UPDATE` y `RESUMED`. Los alias de miembros son `member_join`, `member_leave` y `member_update`. El registro solicita `GUILDS` de forma predeterminada. Usa `registry.intent("GUILD_MESSAGE_REACTIONS")` para las reacciones del servidor y activa `GUILD_MEMBERS` o `GUILD_PRESENCES` en el registro y en el portal de desarrolladores antes de registrar callbacks privilegiados. Los métodos slash, de botones, selectores y formularios manejan las interacciones. Las credenciales del servidor de voz quedan reservadas para el transporte de voz.

## Slash commands and interactions / Comandos slash e interacciones

```tripol
import tripol.discord

func hello(ctx:discord.Interaction)->string {
    return ctx.reply("Hola, " + ctx.option_text("name"), true)
}

options=discord.slash_options()
options.add("name", "Tu nombre", "string", true)
bot=discord.bot()
bot.slash_with_options("hello", "Saluda", &hello, options, "123")
show bot.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"data\":{\"type\":1,\"name\":\"hello\",\"guild_id\":\"123\",\"options\":[{\"type\":3,\"name\":\"name\",\"value\":\"Ana\"}]}}}")
```

Use `bot.slash(name, description, &handler, guild_id="")` for a command without options. An empty guild ID registers globally. A guild ID limits registration to one server. The same name can have a different callback in each registered scope. At startup the bot reads its application ID, compares remote commands with the local definitions and creates or updates the changed ones. Registration leaves unrelated commands in place. This version supports ASCII lowercase command names and string, integer, boolean, user, channel, role, mentionable, number and attachment options.

Usa `bot.slash(nombre, descripción, &manejador, guild_id="")` para un comando sin opciones. Un ID de servidor vacío hace el registro global. Con un ID, el registro se limita a ese servidor. Un mismo nombre puede tener un callback diferente en cada ámbito. Al iniciar, el bot consulta su aplicación, compara los comandos y crea o actualiza los que cambiaron. El registro conserva los demás comandos. Se admiten nombres ASCII en minúsculas y opciones de texto, entero, booleano, usuario, canal, rol, mencionable, número y adjunto.

Register component and modal callbacks with `bot.on_button(id, &handler)`, `bot.on_select(id, &handler)` and `bot.on_modal(id, &handler)`. The handler can return an initial response such as `ctx.reply(text, true)`, `ctx.defer(true)`, `ctx.ack()`, `ctx.reply_components(layout)`, `ctx.update_components(layout)`, `ctx.reply_embeds(batch)` or `ctx.open_modal(form)`. Discord requires an initial response within three seconds.

Registra callbacks con `bot.on_button(id, &manejador)`, `bot.on_select(id, &manejador)` y `bot.on_modal(id, &manejador)`. El manejador puede devolver la respuesta inicial. `ctx.reply(texto, true)` responde en privado, `ctx.defer(true)` difiere la respuesta privada, y `ctx.ack()` confirma una interacción sobre un mensaje. Los métodos de componentes, embeds y formularios devuelven sus respectivos payloads. Discord exige una respuesta inicial en tres segundos.

```tripol
import tripol.discord

func report(ctx:discord.Interaction)->string {
    if ctx.respond(ctx.defer(true)) {
        if ctx.edit_reply("Informe listo") {
            ctx.followup("Consulta el resultado arriba", true)
        }
    }
    return ""
}

bot = discord.bot()
bot.slash("informe", "Prepara un informe", &report)
bot.start()
```

Connected native handlers can use `ctx.respond(payload)` to send the initial response before doing more work. It attempts delivery once and reports whether Discord acknowledged it. Use `ctx.edit_reply(text)` to fill a deferred response or edit its original message. `ctx.edit_response(body)` accepts a Discord edit payload, including `layout.payload()` for an original Components V2 message. `ctx.followup(text, true)` sends a private follow-up and returns its message ID, or an empty string if the request fails. For a response deferred with `ctx.defer()`, edit the original before posting a separate follow-up. The interaction token stays internal. Linux end-to-end tests use a local TLS server. Windows builds use WinHTTP and still need live validation.

Los callbacks nativos conectados pueden usar `ctx.respond(payload)` para enviar la respuesta inicial antes de continuar. Intenta entregarla una vez e indica si Discord la confirmó. Usa `ctx.edit_reply(texto)` para completar una respuesta diferida o editar su mensaje original. `ctx.edit_response(body)` acepta un payload de edición de Discord, incluido `layout.payload()` para un mensaje original Components V2. `ctx.followup(texto, true)` envía un seguimiento privado y devuelve el ID del mensaje, o una cadena vacía si falla la solicitud. Para una respuesta diferida con `ctx.defer()`, edita el original antes de enviar un seguimiento separado. El token de interacción queda en el contexto interno. Las pruebas de extremo a extremo de Linux usan un servidor TLS local. Las compilaciones de Windows usan WinHTTP y siguen pendientes de validación en vivo.

`ctx.option(name)` returns raw JSON, `ctx.option_text(name)` decodes a JSON string, `ctx.values()` returns selected values and `ctx.field(id)` reads a submitted text input. `ctx.field_data(id)` returns the JSON array for select and file-upload inputs. The original interaction data is available as `ctx.data_json`, including resolved attachment metadata.

`ctx.option(nombre)` devuelve JSON, `ctx.option_text(nombre)` decodifica una cadena, `ctx.values()` devuelve los valores seleccionados y `ctx.field(id)` lee un texto enviado. `ctx.field_data(id)` devuelve el arreglo JSON de los selectores y campos de archivos. `ctx.data_json` contiene los datos originales de la interacción y los metadatos resueltos de adjuntos.

After `id = ctx.followup(text)`, use `ctx.edit_followup(id, text)` or `ctx.delete_followup(id)`. Use `ctx.edit_followup_payload(id, layout.payload())` to update an existing Components V2 follow-up. These methods return false on a failed request and require a successful initial response in the same handler.

Después de `id = ctx.followup(texto)`, usa `ctx.edit_followup(id, texto)` o `ctx.delete_followup(id)`. Usa `ctx.edit_followup_payload(id, layout.payload())` para actualizar un seguimiento existente de Components V2. Estos métodos devuelven false si falla la solicitud y requieren una respuesta inicial confirmada dentro del mismo callback.

Use `tripol new Feedback --template discord-slash` for a runnable example. Set `DISCORD_TOKEN`, optionally set `DISCORD_GUILD_ID`, and install the bot application with the `applications.commands` scope on a private test server.

Usa `tripol new Feedback --template discord-slash` para generar un ejemplo ejecutable. Configura `DISCORD_TOKEN`, opcionalmente `DISCORD_GUILD_ID`, e instala la aplicación con el alcance `applications.commands` en un servidor privado.

### Context menus / Menús contextuales

```tripol
import tripol.discord

func inspect(ctx) {
    ctx.reply("ID: " + ctx.target_id(), true)
}

bot = discord.bot()
bot.user_menu("Ver perfil", inspect, "123")
bot.message_menu("Ver mensaje", inspect, "123")
bot.start()
```

User and message menus can share a name because Discord assigns them different command types. `ctx.target_id()` returns the selected ID, and `ctx.target()` returns the user or message object from Discord's resolved data. Use the last argument of `user_menu`, `message_menu`, `slash` or `slash_with_options` for default permissions. For example, `"8192"` requires Manage Messages, `"0"` disables the command by default except for administrators, and `"everyone"` clears that default requirement. An empty value omits the permission field.

Los menús de usuario y mensaje pueden compartir nombre porque Discord les asigna tipos distintos. `ctx.target_id()` devuelve el ID seleccionado y `ctx.target()` devuelve el objeto de usuario o mensaje que resolvió Discord. Usa el último argumento de `user_menu`, `message_menu`, `slash` o `slash_with_options` para configurar los permisos predeterminados. Por ejemplo, `"8192"` requiere Manage Messages, `"0"` desactiva el comando de forma predeterminada salvo para administradores y `"everyone"` elimina ese requisito predeterminado. Una cadena vacía omite el campo.

### Subcommands, choices and autocomplete / Subcomandos, elecciones y autocompletado

```tripol
import tripol.discord

func search(ctx:discord.Interaction)->string {
    return ctx.reply("Busca " + ctx.option_text("query"), true)
}

func suggest(ctx:discord.Interaction)->string {
    choices = discord.command_choices()
    if ctx.focused()=="query" and ctx.option_text("query")=="ca" {
        choices.add("Café", "cafe")
    }
    return ctx.autocomplete(choices)
}

arguments = discord.slash_options()
arguments.add("query", "Qué buscar", "string", true, true)
limits = discord.command_choices()
limits.integer("Diez", 10)
arguments.choose("limit", "Resultados máximos", limits)
commands = discord.slash_options()
commands.subcommand("find", "Busca una bebida", arguments)
root = discord.slash_options()
root.group("drinks", "Bebidas", commands)
bot = discord.bot()
bot.slash_with_options("order", "Haz tu pedido", &search, root)
bot.on_autocomplete("order", &suggest)
bot.start()
```

The example registers `/order drinks find`. `subcommand(name, description, options)` creates a branch and `group(name, description, subcommands)` collects one or more subcommands. Each command or subcommand accepts up to 25 options at its level. Keep required arguments first. The callback reads `ctx.command_group()` and `ctx.subcommand()` when it needs to select a branch. Pass `true` as the fifth argument of `options.add` to enable autocomplete. Use `ctx.focused()` to find the active option and `ctx.option(name)` to read a raw numeric value. Return `ctx.autocomplete(choices)` even when the list is empty. `options.choose` defines a fixed choice list. The builder supports up to 25 text or integer choices, including integer suggestions for number options. Fractional choice values are pending.

El ejemplo registra `/order drinks find`. `subcommand(nombre, descripción, opciones)` crea un subcomando y `group(nombre, descripción, subcomandos)` agrupa uno o más. Cada comando o subcomando admite hasta 25 opciones por nivel. Coloca primero los argumentos obligatorios. El callback lee `ctx.command_group()` y `ctx.subcommand()` para identificar la rama elegida. Pasa `true` como quinto argumento de `options.add` para activar el autocompletado. Usa `ctx.focused()` para identificar la opción activa y `ctx.option(nombre)` para leer un valor numérico en JSON. Devuelve `ctx.autocomplete(elecciones)` aunque la lista esté vacía. `options.choose` define una lista de elecciones fijas. El constructor admite hasta 25 elecciones de texto o enteras, con sugerencias enteras para opciones de tipo número. Quedan pendientes los valores fraccionarios.

## Components V2

```tripol
import tripol.discord

panel = discord.components()
panel.text("## Bienvenido")
panel.separator()
panel.link("Documentación", "https://tripollas.github.io/Tripol/docs/")
message = discord.components()
message.container(panel, 41215)
show message.payload()
```

The payload includes flag 32768 and disables mentions. The builder also has `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` and `gallery(urls)`. It validates the 40-component limit, duplicate button IDs, container nesting and field lengths. Text Display content can total 4,000 characters across the message. The encoded body must fit in the TRIPOL string capacity of 65,535 bytes. Use `message.send(channel_id)` on Windows or Linux. Use `discord.button_row()` to group up to five buttons. `discord.select(id, kind)` supports string, user, role, mentionable and channel menus. `discord.modal(id, title)` supports text inputs, selects, help text, file-upload inputs, radio groups, checkbox groups and individual checkboxes. The Linux launcher can send multipart attachments with `send_files`.

El JSON incluye la bandera 32768 y desactiva las menciones. El constructor también ofrece `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` y `gallery(urls)`. Valida el límite de 40 componentes, los IDs de botones duplicados, el anidamiento y las longitudes de los campos. El contenido de Text Display admite 4.000 caracteres por mensaje. El JSON debe caber en la capacidad de cadena de TRIPOL de 65.535 bytes. Usa `message.send(channel_id)` en Windows o Linux. Usa `discord.button_row()` para agrupar hasta cinco botones. `discord.select(id, tipo)` admite texto, usuario, rol, mencionable y canal. `discord.modal(id, título)` admite texto, selectores, ayuda, campos de archivos, grupos de opción única, grupos de casillas y casillas individuales. Los lanzadores de Windows y Linux envían adjuntos multipart con `send_files`.

For modal choices, create a String Select with the options and pass it to `form.radio(label, menu)` or `form.checkbox_group(label, menu)`. Use `form.checkbox(id, label)` for a single checkbox. Read a radio selection with `ctx.field(id)`, a group or selector with `ctx.field_values(id)` and a checkbox with `ctx.checked(id)`.

Para los campos de opciones, crea un selector de texto con las opciones y pásalo a `form.radio(etiqueta, selector)` o `form.checkbox_group(etiqueta, selector)`. Usa `form.checkbox(id, etiqueta)` para una casilla individual. Lee la selección única con `ctx.field(id)`, el grupo o selector con `ctx.field_values(id)` y la casilla con `ctx.checked(id)`.

#Use `ctx.has_field(id)` to distinguish a submitted empty field from an absent one. `ctx.field_attachment(id, index=0)` returns the selected file’s `id`, `filename`, `url`, `size` and `content_type`. An unresolved file has an empty ID. The method reads the submission without opening or downloading the file.

Usa `ctx.has_field(id)` para distinguir un campo enviado vacío de uno ausente. `ctx.field_attachment(id, index=0)` devuelve `id`, `filename`, `url`, `size` y `content_type`. Un archivo sin resolver tiene ID vacío. El método lee la respuesta del formulario sin abrir ni descargar el archivo.

## Button and select customization / Personalización de botones y selectores

`panel.button("Aceptar", "accept", "success", "✅")` combines text and emoji. Use an empty label for an emoji-only button. Add `true` after the emoji to disable it. For a row, the existing disabled argument comes first, as in `row.button("Aceptar", "accept", "success", false, "✅")`. `row.disable("accept")` disables that button before the row is added to a new layout, and `row.disable("accept", false)` enables it. To update a sent Components V2 message, build the new layout and use `ctx.update_components(panel)` or `room.edit_components(message_id, panel)`.

`panel.button("Aceptar", "accept", "success", "✅")` combina texto y emoji. Usa una etiqueta vacía para mostrar solo el emoji. Añade `true` después del emoji para deshabilitarlo. En una fila se conserva primero el argumento de deshabilitación, como en `row.button("Aceptar", "accept", "success", false, "✅")`. `row.disable("accept")` deshabilita ese botón antes de añadir la fila a un panel nuevo. Usa `row.disable("accept", false)` para habilitarlo. Para actualizar un mensaje Components V2 ya enviado, crea el nuevo panel y usa `ctx.update_components(panel)` o `room.edit_components(message_id, panel)`.

Use Unicode, a custom emoji ID, `"<:check:123>"` or `"<a:party:123>"` for an animated custom emoji. Pass it as the last argument of `menu.option(label, value, description, selected, emoji)`. `menu.disable()` applies to a select in a message. Discord handles link and premium buttons directly. See the executable example for all six styles. [Buttons example](../examples/discord_buttons.tri)

Usa Unicode, un ID de emoji personalizado, `"<:check:123>"` o `"<a:party:123>"` para un emoji personalizado animado. Pásalo como último argumento de `menu.option(etiqueta, valor, descripción, seleccionado, emoji)`. `menu.disable()` se aplica a los selectores de mensajes. Discord gestiona los botones de enlace y de compra. El ejemplo ejecutable reúne los seis estilos. [Ejemplo de botones](../examples/discord_buttons.tri)

### Attachments / Adjuntos

```tripol
import tripol.discord

panel = discord.container()
panel.file("report.txt")
message_id = panel.send_files("123", list(["./report.txt"]))
show message_id
```

With the Linux or Windows launcher, `send_files` streams 1–10 local files into a multipart channel message and returns its ID or an empty string on failure. `ctx.reply_files(content, paths, true)` sends a private initial response immediately and returns whether it was accepted. It accepts text, an embed, a batch or Components V2. Use distinct ASCII filenames and use the exact basename in `discord.attachment_url("photo.png")` or `panel.file("report.txt")`. Follow-up uploads and replacement of existing attachments remain pending.

Con el lanzador de Linux o Windows, `send_files` adjunta entre 1 y 10 archivos locales a un canal y devuelve el ID o una cadena vacía si falla. `ctx.reply_files(contenido, rutas, true)` envía de inmediato la respuesta inicial privada e indica si se aceptó. Admite texto, un embed, un lote o Components V2. Usa nombres ASCII distintos y el nombre exacto del archivo en `discord.attachment_url("photo.png")` o `panel.file("report.txt")`. Quedan pendientes los adjuntos en seguimientos y la sustitución de archivos.

For uploads submitted by users, call `form.file("files", "Upload", list(["image", ".pdf"]), 0, 2, false)`. Discord applies up to ten extension filters. Check the actual contents before trusting an upload.

Para archivos enviados por usuarios, llama a `form.file("files", "Adjuntar", list(["image", ".pdf"]), 0, 2, false)`. Discord aplica hasta diez filtros por extensión. Comprueba el contenido del archivo antes de confiar en él.

## Discord Studio

```sh
tripol discord studio main.tri --metrics studio.jsonl
```

On Linux x86-64, Studio runs the bot and samples its RSS, CPU ticks, connection state, heartbeat latency, uptime, guild count at READY, queued dispatches, REST requests, errors and rate limits. The launcher and bot process are reported separately. The output file must be new. Without `--metrics`, snapshots are sent to standard error. The export contains counters and process measurements. Live command and interaction metrics, log filtering and Windows monitoring are pending.

En Linux x86-64, Studio ejecuta el bot y toma muestras de RSS, ticks de CPU, conexión, latencia de heartbeat, tiempo activo, servidores en READY, cola de eventos, peticiones REST, errores y límites. El lanzador y el bot se miden por separado. Usa un archivo de salida nuevo. Sin `--metrics`, las muestras se escriben en la salida de error. Quedan pendientes las métricas por comando e interacción, el filtrado de logs y el monitoreo en Windows.

Studio adds `events_by_type`, `interactions_by_type` and `commands_by_name` to its Linux samples. They count events admitted to the Gateway queue, not callback results. Command keys include their type, such as `user/Ver perfil`. The counters contain neither interaction tokens nor component IDs. Names beyond the tracked set are counted under `other`.

Studio añade `events_by_type`, `interactions_by_type` y `commands_by_name` a sus muestras de Linux. Cuentan los eventos admitidos en la cola del Gateway. Los resultados de los callbacks se registran por separado si el programa los mide. Las claves incluyen el tipo del comando, como `user/Ver perfil`. Los contadores mantienen fuera los tokens de interacción y los IDs de componentes. Los nombres que exceden el conjunto registrado se agrupan en `other`.

## Member moderation and welcome bot / Moderación de miembros y bienvenida

```tripol
import tripol.discord

server = discord.guild("123")
member = server.member("456")
show member.read()
```

Use the numeric IDs from an authorized server. `member.kick()`, `member.ban(delete_message_seconds=0)` and `member.unban()` remove a member, ban them with optional removal of up to seven days of messages, or revoke a ban. `member.nickname(text)` sets a nickname, and an empty string clears it. `member.timeout(until)` sets an ISO 8601 deadline, and `member.timeout()` clears it. Discord enforces the 28-day maximum. `member.role(id, false)` removes a role and `member.role(id)` adds it. These methods report success as a boolean so the bot can continue after a failed request. Discord enforces bot permissions and role hierarchy. `member.dm(text)` and `discord.post_message(channel_id, text)` return the sent message ID or an empty string if delivery fails.

Usa los IDs numéricos de un servidor autorizado. `member.kick()`, `member.ban(delete_message_seconds=0)` y `member.unban()` expulsan, banean con borrado opcional de hasta siete días de mensajes o revocan el baneo. `member.nickname(texto)` cambia el apodo y una cadena vacía lo elimina. `member.timeout(hasta)` recibe una fecha ISO 8601 y `member.timeout()` retira el aislamiento. Discord aplica el máximo de 28 días. `member.role(id, false)` quita un rol y `member.role(id)` lo añade. Estos métodos devuelven un booleano para que el bot pueda continuar si falla la solicitud. Discord aplica los permisos y la jerarquía de roles del bot. `member.dm(texto)` y `discord.post_message(channel_id, texto)` devuelven el ID del mensaje o una cadena vacía si falla el envío.

Use `tripol new Welcome --template discord-welcome` to send a welcome to a channel when a human joins the configured server. The generated README covers the required environment variables, Server Members Intent and channel permissions. The template skips bots and other servers, and continues receiving events if a send fails.

Usa `tripol new Bienvenida --template discord-welcome` para enviar un saludo a un canal cuando se une una persona al servidor configurado. El README generado indica las variables de entorno, Server Members Intent y los permisos del canal. La plantilla omite bots y otros servidores y continúa recibiendo eventos si falla un envío.

Create a restricted bot with `tripol new Moderator --template discord-moderation`. Set `DISCORD_TOKEN`, `DISCORD_GUILD_ID` and `DISCORD_OPERATOR_ID`. The `/kick`, `/ban`, `/unban`, `/timeout` and `/warn` commands check the operator, guild and effective permissions returned by Discord. Verified callers receive a private acknowledgment. `/warn` appends a JSON line to `discord-warnings.jsonl` or the path in `DISCORD_WARNINGS_PATH`. Keep it in a directory writable only by the bot account.

Crea un bot restringido con `tripol new Moderador --template discord-moderation`. Configura `DISCORD_TOKEN`, `DISCORD_GUILD_ID` y `DISCORD_OPERATOR_ID`. Los comandos `/kick`, `/ban`, `/unban`, `/timeout` y `/warn` comprueban operador, servidor y permisos efectivos devueltos por Discord. Se confirma en privado a los usuarios autorizados. `/warn` guarda una línea JSON en `discord-warnings.jsonl` o en la ruta de `DISCORD_WARNINGS_PATH`. Usa un directorio donde solo pueda escribir la cuenta del bot.

## Channel moderation / Moderación de canales

```tripol
import tripol.discord

room = discord.channel_ref("123")
show room.read()
```

`room.slowmode(15)` sets a 15-second delay and `room.slowmode(0)` clears it. Use `room.pin(message_id)` to pin and `room.pin(message_id, false)` to unpin. `room.delete_message(id)` removes one message. `room.purge(ids)` requires a list of 2–100 distinct message IDs. Discord enforces bulk deletion eligibility, including the 14-day age limit. `room.react(id, "🔥")` adds the bot's reaction. Use `false` as the last argument to remove it. Unicode emoji and custom `name:id` values are percent-encoded automatically.

`room.slowmode(15)` fija una espera de 15 segundos y `room.slowmode(0)` la retira. Usa `room.pin(message_id)` para fijar un mensaje y `room.pin(message_id, false)` para desfijarlo. `room.delete_message(id)` borra un mensaje. `room.purge(ids)` requiere entre 2 y 100 IDs distintos. Discord aplica los requisitos del borrado masivo, incluido el límite de 14 días de antigüedad. `room.react(id, "🔥")` añade la reacción del bot. Usa `false` como último argumento para retirarla. Los emoji Unicode y los valores personalizados `nombre:id` se codifican automáticamente.

Configure the token of an authorized bot. Discord checks the required permissions. These methods return a boolean on mutations and keep the native program running if the request fails.

Configura el token de un bot autorizado. Discord comprueba los permisos requeridos. Los métodos que modifican el canal o un mensaje devuelven un booleano y permiten que el programa siga ejecutándose si falla la solicitud.

### Channels, categories and audit / Canales, categorías y auditoría

Use `server.create_channel(name, kind, parent_id, topic)` with `"text"`, `"voice"`, `"category"` or `"announcement"`. The result is the new ID or an empty string on failure. `server.channels()` returns the REST response as JSON. Read the audit log with `server.audit(limit, before, action)`. On a channel, use `rename`, `topic`, `move` and `delete`. An empty parent in `move()` detaches the channel from its category. `edit_message` and `edit_components` update messages owned by the bot.

Usa `server.create_channel(nombre, tipo, parent_id, tema)` con `"text"`, `"voice"`, `"category"` o `"announcement"`. Devuelve el ID nuevo o una cadena vacía si falla. `server.channels()` devuelve la respuesta REST como JSON. Consulta auditoría con `server.audit(límite, antes, acción)`. En un canal están disponibles `rename`, `topic`, `move` y `delete`. Un padre vacío en `move()` saca el canal de su categoría. `edit_message` y `edit_components` actualizan mensajes del propio bot.

### Guild roles / Roles de servidor

Use `server.create_role("Reviewer", "0", "#FFFFFF")` to create a role with no permissions of its own. The method returns its ID. `server.roles()` lists the roles. Use `role = server.role_ref(id)` with `role.read()`, `role.edit(name, permissions, color)` and `role.delete()`. Empty edit arguments leave those properties as they are. `role.permissions()` reads the combined guild permissions of the role and @everyone as a decimal string. Effective permissions for a member in a channel also depend on other roles and overrides.

Usa `server.create_role("Revisor", "0", "#FFFFFF")` para crear un rol sin permisos propios. El método devuelve su ID. `server.roles()` lista los roles. Usa `role = server.role_ref(id)` con `role.read()`, `role.edit(nombre, permisos, color)` y `role.delete()`. Los argumentos vacíos de edición conservan esas propiedades. `role.permissions()` devuelve como cadena decimal los permisos combinados del rol y @everyone en el servidor. Los permisos efectivos de un miembro en un canal también dependen de los demás roles y las excepciones.

## API

| Method / Método | Behavior / Comportamiento |
| --- | --- |
| `discord.bot()` | Create a bot configuration / Crea una configuración de bot |
| `bot.prefix(value)` | Set a prefix of up to 32 UTF-8 bytes / Define un prefijo de hasta 32 bytes UTF-8 |
| `bot.status(value)` | Set `online`, `idle`, `dnd` or `invisible` / Define el estado indicado |
| `bot.activity(value)` | Set the activity text, up to 128 UTF-8 bytes / Define una actividad de hasta 128 bytes UTF-8 |
| `bot.command(name, response)` | Set one of up to 16 exact text commands / Configura uno de hasta 16 comandos de texto exactos |
| `bot.preview(content)` | Return the exact reply or an empty string / Devuelve la respuesta exacta o texto vacío |
| `bot.start()` | Open a blocking Gateway session / Abre una sesión Gateway bloqueante |
| `discord.embed()` | Create an empty embed / Crea un embed vacío |
| `embed.title(value)` | Set the title / Define el título |
| `embed.description(value)` | Set the description / Define la descripción |
| `embed.footer(value)` | Set the footer text / Define el pie |
| `embed.url(value)`, `embed.image(value)`, `embed.thumbnail(value)` | Set the title link, image or thumbnail / Define el enlace, la imagen o la miniatura |
| `embed.color(value)` | Set an RGB color as `#RRGGBB` / Define un color RGB como `#RRGGBB` |
| `embed.field(label, value, inline=false)` | Append one of up to 25 fields / Añade uno de hasta 25 campos |
| `embed.payload()` | Build JSON including disabled mentions / Genera el JSON con las menciones desactivadas |
| `embed.send(channel_id)` | Send the embed over Discord REST from Windows or Linux / Envía el embed por Discord REST desde Windows o Linux |

Text commands accept names of 1–96 bytes and replies of 1–512 bytes, both UTF-8. Use `embed.author(name, url="", icon_url="")` and `embed.timestamp(rfc3339)` for authors and timestamps. `discord.embeds()` groups up to 10 embeds with 6,000 characters of combined text. Individual limits are 256 for title, author and field name, 4,096 for description, 2,048 for footer, and 1,024 for field value. The optional `intents` parameter on legacy `discord.run()` accepts `33281` (default for guild prefix commands) or `1` (GUILDS only).

Los comandos aceptan nombres de 1–96 bytes y respuestas de 1–512 bytes, ambos en UTF-8. Usa `embed.author(nombre, url="", icon_url="")` y `embed.timestamp(rfc3339)` para autor y fecha. `discord.embeds()` agrupa hasta 10 embeds con 6.000 caracteres de texto combinado. Los límites son 256 para título, autor y nombre de campo, 4.096 para descripción, 2.048 para pie y 1.024 para valor de campo. El parámetro opcional `intents` de `discord.run()` acepta `33281` (predeterminado para comandos de texto) o `1` (solo GUILDS).

## Discord R5 — Windows WinHTTP backend

`discord.run()` now uses a persistent TLS WebSocket to `gateway.discord.gg/?v=10&encoding=json`, not REST polling. The native WinHTTP upgrade, send and receive calls are compiled directly into the executable. A native receiver thread reads one complete UTF-8 message (including fragments) into bounded private storage. The main thread waits locally, wakes on message arrival and schedules Gateway heartbeats; no 2.5-second network polling. Sending messages still uses Discord HTTPS REST POST. No Python, Node.js or external Discord SDK is required by the generated executable.

### Create a legacy bot

Create a bot at https://discord.com/developers/applications, enable Message Content Intent for prefix commands, install it on a small private test server and give it View Channel and Send Messages permissions. Read Message History is needed for the explicit REST fallback. A channel ID is optional: only select one if you want to restrict replies or send a greeting before any command. Use a bot token, not a user token.

```tripol
import tripol.discord

discord.login_env()
discord.activity("Programando en TRIPOL")
discord.reply("!hola", "Hola compa!")
discord.reply("!ping", "Pong!")
discord.run()
```

`setup` simply combines `login` and `channel`. These older methods continue to work. The optional `welcome` sends a greeting once on READY; `say` sends immediately by REST. `status("online")` is optional (default online); `idle`, `dnd`, `invisible` also supported in the initial Identify payload. `activity` sets a playing activity, up to 128 UTF-8 bytes. Up to 16 exact commands across accessible guild channels (or one selected channel); duplicate registration replaces the response. Bots and webhooks are ignored. Prefix commands match decoded text exactly, including UTF-8 characters represented with JSON escapes.

For sharing code or binaries, do not embed a token. Instead use:

```powershell
$secure = Read-Host 'Bot token' -AsSecureString
$env:DISCORD_TOKEN = [System.Net.NetworkCredential]::new('', $secure).Password
tripol new MiBot --template discord
cd MiBot
notepad main.tri
tripol run
```

The generated template does not need a channel ID. `login_env()` reads `DISCORD_TOKEN`; optional `channel(...)` restricts replies. Never publish or share the token; a literal is extractable from both source and compiled binary. Reset a leaked token in the developer portal. Remove the environment variable after stopping: `Remove-Item Env:DISCORD_TOKEN`.

### Protocol and limitations

Initial Hello (10), Identify (2), Dispatch (0), sequence tracking, periodic Heartbeat (1) with null before the first sequence, Heartbeat ACK (11), server heartbeat requests, READY and MESSAGE_CREATE are implemented. `connected()` reports local READY state; `ping()` reports last heartbeat round trip in milliseconds or -1 before ACK, not message latency. Both are primarily diagnostics while running; the simple `run` API blocks. `preview` and `gateway_reply` support offline matching tests. `identify_payload` is a protocol helper containing the secret; do not print or log it.

The standalone Windows PE uses WinHTTP and the system certificate store. It accepts uncompressed text frames up to 65,535 bytes and dispatches callbacks in order. Use the TRIPOL launcher on Windows for Resume, transport compression, HTTPS reuse, rate-limit handling and multipart uploads.

Explicit fallback: `discord.run_rest()` restores ECLIPSE's 2.5-second REST polling, which may miss messages. Legacy `listen`, `send`, `latest`, `matches`, `message_id`, `gateway` and token/URL helpers remain available. See `tripol docs discord` for the complete module, including the protocol helpers.

Official protocol: https://discord.com/developers/docs/events/gateway ; WinHTTP WebSocket API: https://learn.microsoft.com/windows/win32/api/winhttp/nf-winhttp-winhttpwebsocketcompleteupgrade ; rate limits: https://discord.com/developers/docs/topics/rate-limits .

### ORIGIN routing

`configured()` now means a token is configured, rather than a token and channel. `has_channel()` checks whether a channel filter exists. `all_channels()` clears the filter. `send_to(channel, message)` uses the configured bot token; `say()` and a nonempty `welcome()` require a selected channel. `run_rest()` still requires a channel. `event_channel(event)` validates the numeric channel ID and optional filter; `gateway_reply(event)` ignores bot/webhook authors as before. Guild message intents are unchanged; DMs are not enabled. Tests verify routing on Linux in all four profiles, but live Windows Gateway is not validated here.

## Cooldowns and command permissions / Cooldowns y permisos de comandos

Register the slash command before configuring it. `bot.cooldown("work", 60)` reserves one execution per user every 60 seconds. Choose `"channel"`, `"guild"` or `"global"` for a shared cooldown within the command’s registration scope. An optional message is followed by the remaining seconds. Use zero seconds to disable the rule. A command that needs a channel, server or user ID is denied when that context is absent.

Registra el comando slash antes de configurarlo. `bot.cooldown("work", 60)` reserva una ejecución por usuario cada 60 segundos. Elige `"channel"`, `"guild"` o `"global"` para compartir el cooldown dentro del ámbito de registro del comando. Un mensaje opcional precede a los segundos restantes. Usa cero segundos para desactivar la regla. Si falta el ID requerido por el ámbito, el comando se rechaza.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {ctx.reply("Done",true)})
show bot.cooldown("work",60,"user","Try again in")
```

`ctx.cooldown_remaining()` returns rounded-up seconds, or `-1` when unavailable. Use `ctx.reset_cooldown()` for the current bucket. For administrative tools, use `discord.cooldown_bucket("work","user",user_id,guild_id)` with `.remaining()` and `.reset()`. A reservation is made before the handler and remains after an error unless the handler resets it. `bot.persist_cooldowns("cooldowns.json")` keeps active cooldowns across restarts. Use a separate file per app in a directory writable by its account. In persistent mode each successful reservation or reset saves a snapshot. A file lock prevents two launchers from writing the same file. In-memory entries expire automatically.

`ctx.cooldown_remaining()` devuelve los segundos redondeados hacia arriba o `-1` cuando la consulta no está disponible. Usa `ctx.reset_cooldown()` para restablecer el registro actual. Para herramientas administrativas, usa `discord.cooldown_bucket("work","user",user_id,guild_id)` con `.remaining()` y `.reset()`. La reserva se hace antes del callback y se mantiene si ocurre un error. El callback puede restablecerla. `bot.persist_cooldowns("cooldowns.json")` conserva los cooldowns activos al reiniciar. Usa un archivo distinto por aplicación en un directorio donde pueda escribir su cuenta. Cuando la persistencia está activa, cada reserva o restablecimiento guarda una instantánea. Un candado impide la escritura por dos lanzadores. Los registros en memoria caducan automáticamente.

`bot.permission("ban","BAN_MEMBERS")` sets the registration default and checks the effective permissions Discord sends for the invoking member and bot in the current channel. For kick, ban and timeout, a single top-level user option is inferred as the member to protect. Specify the option as the third argument when the target cannot be inferred. Pass the registration guild ID as the fourth argument to `permission` and the fifth to `cooldown`. Run protected moderation and cooldowns through `tripol run` or `tripol discord-run`.

`bot.permission("ban","BAN_MEMBERS")` configura el permiso predeterminado y comprueba los permisos efectivos del miembro y del bot que Discord comunica para el canal. Para expulsar, banear y aislar se infiere el objetivo cuando hay una sola opción de usuario en el nivel superior. Indica la opción como tercer argumento si el objetivo no se puede deducir. El ID de servidor de registro es el cuarto argumento de `permission` y el quinto de `cooldown`. Ejecuta la moderación protegida y los cooldowns con `tripol run` o `tripol discord-run`.

The hierarchy check protects the owner, the operator and the bot itself. It verifies the highest roles for the member, target and bot, and denies timeout of an administrator. `ctx.can()`, `ctx.bot_can()` and `ctx.can_moderate()` are available for explicit checks. Lower-level REST methods use the bot’s authorization on Discord. Use a protected callback when accepting administrative requests from users.

La comprobación de jerarquía protege al propietario, al operador y al propio bot. Verifica los roles superiores del miembro, el objetivo y el bot, y rechaza el aislamiento de un administrador. `ctx.can()`, `ctx.bot_can()` y `ctx.can_moderate()` permiten comprobaciones explícitas. Los métodos REST de bajo nivel usan la autorización del bot en Discord. Usa un callback protegido para aceptar solicitudes administrativas de usuarios.

Use `bot.restrict("ban","role",role_id)` to restrict a command to a role. Use `"user"`, `"channel"` or `"guild"` for those scopes. IDs within a scope are alternatives. A command that configures several scopes must satisfy each of them. `ctx.has_role(role_id)` supports an explicit role check. The launcher evaluates restrictions before reserving a cooldown or starting a hierarchy check.

Usa `bot.restrict("ban","role",role_id)` para restringir un comando a un rol. Usa `"user"`, `"channel"` o `"guild"` para esos ámbitos. Los IDs de un mismo ámbito son alternativas. Un comando con varios ámbitos configurados debe satisfacer cada uno. `ctx.has_role(role_id)` permite comprobar un rol expresamente. El lanzador evalúa las restricciones antes de reservar el cooldown o verificar la jerarquía.

## Channel and resource management / Gestión de canales y recursos

`channel.lock()` saves the @everyone overwrite and denies Send Messages for that role. Explicit user and role overrides continue to apply. `channel.unlock()` restores exactly the saved overwrite in the same launcher, or removes the temporary overwrite if the channel did not have one. If the overwrite has changed externally, the operation returns false and preserves the external edit. Use `overwrite(id, allow, deny, "user")` or `overwrite(id, allow, deny, "role")` for explicit changes, and `clear_overwrite(id)` to remove one. `guild.private_channel(name, user_id)` creates a text channel that grants that member access and denies viewing to @everyone. `guild.move_channel(id, position)` updates the channel order.

`channel.lock()` guarda el permiso de @everyone y le deniega Enviar mensajes. Los permisos explícitos de otros roles y usuarios siguen aplicándose. `channel.unlock()` restaura la configuración guardada dentro del mismo lanzador o quita el permiso temporal si antes no existía. Si el permiso cambió externamente, devuelve false y conserva el cambio. Usa `overwrite(id, allow, deny, "user")` o `overwrite(id, allow, deny, "role")` para cambios explícitos, y `clear_overwrite(id)` para eliminar uno. `guild.private_channel(name, user_id)` crea un canal de texto que permite el acceso a ese miembro y deniega la vista a @everyone. `guild.move_channel(id, position)` actualiza el orden.

`channel.create_webhook(name)` returns a typed webhook. Use `.read()`, `.edit(name)` and `.delete()` to administer it. The returned read information omits credentials. `.send()` accepts text, embeds or Components V2. `.send_files()` streams local attachments through the launcher. Discord controls which components each webhook can use. `channel.invite(max_age, max_uses)` returns the code; list invitations with `.invites()`, inspect one with `discord.invite_info(code)` and remove it with `discord.revoke_invite(code)`.

`channel.create_webhook(name)` devuelve un webhook tipado. Usa `.read()`, `.edit(name)` y `.delete()` para administrarlo. La información devuelta al consultar omite las credenciales. `.send()` recibe texto, embeds o Components V2. `.send_files()` transmite adjuntos locales mediante el lanzador. Discord controla qué componentes puede usar cada webhook. `channel.invite(max_age, max_uses)` devuelve el código. Enumera invitaciones con `.invites()`, consulta una con `discord.invite_info(code)` y elimínala con `discord.revoke_invite(code)`.

Create a thread with `channel.thread(name)` or `channel.thread_from(message_id, name)`. Use `thread.thread_member(id)`, `thread.archive()` and `thread.lock_thread()`. For forums, use `guild.create_channel(name, "forum")`, `forum.post(title, text, tags)` and `post.apply_tags(tags)`. Use `discord.forum_tags()` to build the full set supplied to `forum.set_forum_tags(tags)`, including the IDs of existing labels that should stay. The [live validation guide](discord-validation.md) lists the required permissions and expected results.

Crea un hilo con `channel.thread(name)` o `channel.thread_from(message_id, name)`. Usa `thread.thread_member(id)`, `thread.archive()` y `thread.lock_thread()`. Para foros, usa `guild.create_channel(name, "forum")`, `forum.post(title, text, tags)` y `post.apply_tags(tags)`. Usa `discord.forum_tags()` para construir el conjunto completo que recibe `forum.set_forum_tags(tags)), incluidos los IDs de las etiquetas que quieras conservar. La [guía de validación conectada](discord-validation.md) recoge los permisos y resultados esperados.


---

# Complete standard library reference / Biblioteca estándar completa

Source: `playground/library-catalog-r6.1.js` at gh-pages `de3b2448bad3bcefd718ea8e1a786937bc6b0b18`. 64 modules, 449 functions and 196 object methods (645 catalog entries). Native flags, targets, examples and expected output are copied from the catalog; they are not proof of live execution or full browser support. Signatures identify documented entries, not necessarily every overload.

Module index: `arrays` (8), `cli` (5), `collections` (6), `console` (3), `convert` (8), `csv` (3), `datetime` (7), `errors` (7), `files` (5), `math` (12), `path` (6), `random` (3), `search` (4), `sort` (4), `text` (9), `validate` (6), `animation` (6), `buffer` (8), `camera` (4), `canvas` (10), `colors` (8), `dictionary` (8), `geometry` (6), `heap` (6), `image` (6), `io` (2), `json` (17), `matrix` (6), `particles` (6), `physics` (5), `queue` (7), `regex` (2), `sequences` (4), `sprite` (6), `stack` (7), `stats` (6), `stream` (6), `testing` (4), `time` (6), `utf8` (5), `vector` (7), `atomic` (4), `audio` (6), `binary` (8), `bits` (6), `bitset` (6), `checksum` (3), `compression` (3), `crypto` (3), `discord` (98 functions + 196 methods), `encoding` (4), `env` (1), `events` (14), `game` (7), `http` (3), `https` (1), `memory` (4), `network` (3), `process` (3), `profiler` (2), `simd` (1), `system` (5), `thread` (3), `ui` (7).

## tripol.arrays — 8 functions

Import: `import tripol.arrays`

### arrays.make

`func arrays.make(n:int, value:int) -> ints`

EN: Create an array filled with a single value.

ES: Crea un arreglo lleno del mismo valor.

Native example flag: false.

```tripol
import tripol.arrays

show arrays.make(3, 7)
```

Documented expected output:

```text
[7, 7, 7]
```

### arrays.total

`func arrays.total(a:ints) -> int`

EN: Sum every integer in an array.

ES: Suma todos los enteros de un arreglo.

Native example flag: false.

```tripol
import tripol.arrays

a = [1, 2, 3]
show arrays.total(a)
```

Documented expected output:

```text
6
```

### arrays.count

`func arrays.count(a:ints, value:int) -> int`

EN: Count occurrences of a value.

ES: Cuenta las apariciones de un valor.

Native example flag: false.

```tripol
import tripol.arrays

a = [2, 5, 2]
show arrays.count(a, 2)
```

Documented expected output:

```text
2
```

### arrays.contains

`func arrays.contains(a:ints, value:int) -> bool`

EN: Check whether a value appears in an array.

ES: Comprueba si un valor aparece en un arreglo.

Native example flag: false.

```tripol
import tripol.arrays

a = [1, 2, 3]
show arrays.contains(a, 2)
```

Documented expected output:

```text
true
```

### arrays.reversed

`func arrays.reversed(a:ints) -> ints`

EN: Return a new array in reverse order.

ES: Devuelve un arreglo nuevo en orden inverso.

Native example flag: false.

```tripol
import tripol.arrays

a = [1, 2, 3]
show arrays.reversed(a)
```

Documented expected output:

```text
[3, 2, 1]
```

### arrays.copy

`func arrays.copy(a:ints) -> ints`

EN: Create an independent copy of an array.

ES: Crea una copia independiente de un arreglo.

Native example flag: false.

```tripol
import tripol.arrays

a = [1, 2]
b = arrays.copy(a)
b[0] = 9
show a, b
```

Documented expected output:

```text
[1, 2] [9, 2]
```

### arrays.length

`func arrays.length(a:ints) -> int`

EN: Count elements in an array.

ES: Cuenta los elementos de un arreglo.

Native example flag: false.

```tripol
import tripol.arrays

a = [10, 20, 30]
show arrays.length(a)
```

Documented expected output:

```text
3
```

### arrays.copy_into

`func arrays.copy_into(destination:ints, source:ints) -> bool`

EN: Copy values into an existing dynamic array of equal length.

ES: Copia valores a un arreglo dinámico existente de igual longitud.

Native example flag: false.

Note: Returns false if lengths differ; the destination is unchanged.

Nota: Devuelve false si las longitudes difieren; el destino no cambia.

```tripol
import tripol.arrays

source = [4, 7]
destination = ints(2)
show arrays.copy_into(destination, source)
show destination
```

Documented expected output:

```text
true
[4, 7]
```

---

## tripol.cli — 5 functions

Import: `import tripol.cli`

### cli.option_name

`func cli.option_name(name:string) -> bool`

EN: Validate an option name without the leading two dashes.

ES: Valida un nombre de opción sin los dos guiones iniciales.

Native example flag: false.

```tripol
import tripol.cli

show cli.option_name("retries")
```

Documented expected output:

```text
true
```

### cli.find_option

`func cli.find_option(line:string,name:string) -> int`

EN: Find the byte offset of the last named option in the input text, or return -1.

ES: Busca la posición en bytes de la última opción indicada o devuelve -1.

Native example flag: false.

```tripol
import tripol.cli

show cli.find_option("build --fast","fast")
```

Documented expected output:

```text
6
```

### cli.has

`func cli.has(line:string,name:string) -> bool`

EN: Check for a named option in whitespace-delimited input text.

ES: Comprueba si aparece la opción indicada en texto separado por espacios.

Native example flag: false.

```tripol
import tripol.cli

show cli.has("--fast --count=3","fast")
```

Documented expected output:

```text
true
```

### cli.value

`func cli.value(line:string,name:string,fallback="") -> string`

EN: Read the last value assigned to the named option, using a fallback when absent.

ES: Lee el último valor de la opción indicada o usa el valor por defecto si falta.

Native example flag: false.

```tripol
import tripol.cli

show cli.value("--mode=fast","mode","slow")
```

Documented expected output:

```text
fast
```

### cli.integer

`func cli.integer(line:string,name:string,fallback=0) -> int`

EN: Read an integer option or use the given fallback value.

ES: Lee una opción entera o usa el valor por defecto indicado.

Native example flag: false.

```tripol
import tripol.cli

show cli.integer("--count=3","count",1)
```

Documented expected output:

```text
3
```

---

## tripol.collections — 6 functions

Import: `import tripol.collections`

### collections.append

`func collections.append(a:ints,value:int) -> ints`

EN: Return a new array with an extra element at the end.

ES: Devuelve un arreglo nuevo con un elemento adicional al final.

Native example flag: false.

Note: The original array is not modified.

Nota: El arreglo original no cambia.

```tripol
import tripol.collections

a = [1, 2]
show collections.append(a, 3)
```

Documented expected output:

```text
[1, 2, 3]
```

### collections.concat

`func collections.concat(a:ints,b:ints) -> ints`

EN: Join two arrays into a new one.

ES: Une dos arreglos en uno nuevo.

Native example flag: false.

```tripol
import tripol.collections

a = [1, 2]
b = [3, 4]
show collections.concat(a, b)
```

Documented expected output:

```text
[1, 2, 3, 4]
```

### collections.part

`func collections.part(a:ints,start:int,end:int) -> ints`

EN: Copy a range, including start but excluding end.

ES: Copia un intervalo, incluyendo start pero excluyendo end.

Native example flag: false.

```tripol
import tripol.collections

a = [10, 20, 30, 40]
show collections.part(a, 1, 3)
```

Documented expected output:

```text
[20, 30]
```

### collections.remove_at

`func collections.remove_at(a:ints,index:int) -> ints`

EN: Copy an array without an indexed element.

ES: Copia un arreglo sin el elemento indicado.

Native example flag: false.

```tripol
import tripol.collections

a = [10, 20, 30]
show collections.remove_at(a, 1)
```

Documented expected output:

```text
[10, 30]
```

### collections.insert

`func collections.insert(a:ints,index:int,value:int) -> ints`

EN: Make a new array with a value inserted at an index.

ES: Crea un arreglo nuevo insertando un valor en un índice.

Native example flag: false.

```tripol
import tripol.collections

a = [10, 30]
show collections.insert(a, 1, 20)
```

Documented expected output:

```text
[10, 20, 30]
```

### collections.equal

`func collections.equal(a:ints,b:ints) -> bool`

EN: Check whether two arrays contain the same ordered values.

ES: Comprueba si dos arreglos contienen los mismos valores en el mismo orden.

Native example flag: false.

```tripol
import tripol.collections

a = [1, 2]
b = [1, 2]
show collections.equal(a, b)
```

Documented expected output:

```text
true
```

---

## tripol.console — 3 functions

Import: `import tripol.console`

### console.ask

`func console.ask(prompt:string) -> string`

EN: Display a prompt and read a line of text.

ES: Muestra una pregunta y lee una línea de texto.

Native example flag: false.

Note: Enter Ada when prompted; output includes the prompt and Ada.

Nota: Introduce Ada cuando se solicite; la salida incluye la pregunta y Ada.

```tripol
import tripol.console

name = console.ask("Name:")
show name
```

Documented input:

```text
Ada
```

### console.ask_int

`func console.ask_int(prompt:string) -> int`

EN: Display a prompt and read an integer.

ES: Muestra una pregunta y lee un entero.

Native example flag: false.

Note: Enter 16 to obtain 17 after the prompt.

Nota: Introduce 16 para obtener 17 después de la pregunta.

```tripol
import tripol.console

age = console.ask_int("Age:")
show age + 1
```

Documented input:

```text
16
```

### console.line

`func console.line(message:string) -> int`

EN: Print a line and return zero.

ES: Imprime una línea y devuelve cero.

Native example flag: false.

```tripol
import tripol.console

console.line("Hola TRIPOL")
```

Documented expected output:

```text
Hola TRIPOL
```

---

## tripol.convert — 8 functions

Import: `import tripol.convert`

### convert.celsius_to_fahrenheit

`func convert.celsius_to_fahrenheit(c:float) -> float`

EN: Convert Celsius to Fahrenheit.

ES: Convierte grados Celsius a Fahrenheit.

Native example flag: false.

```tripol
import tripol.convert

show convert.celsius_to_fahrenheit(20.0)
```

Documented expected output:

```text
68.000000
```

### convert.fahrenheit_to_celsius

`func convert.fahrenheit_to_celsius(f:float) -> float`

EN: Convert Fahrenheit to Celsius.

ES: Convierte grados Fahrenheit a Celsius.

Native example flag: false.

```tripol
import tripol.convert

show convert.fahrenheit_to_celsius(68.0)
```

Documented expected output:

```text
20.000000
```

### convert.km_to_miles

`func convert.km_to_miles(k:float) -> float`

EN: Convert kilometers to miles.

ES: Convierte kilómetros a millas.

Native example flag: false.

```tripol
import tripol.convert

show convert.km_to_miles(1.609344)
```

Documented expected output:

```text
1.000000
```

### convert.miles_to_km

`func convert.miles_to_km(m:float) -> float`

EN: Convert miles to kilometers.

ES: Convierte millas a kilómetros.

Native example flag: false.

```tripol
import tripol.convert

show convert.miles_to_km(1.0)
```

Documented expected output:

```text
1.609344
```

### convert.kg_to_lb

`func convert.kg_to_lb(k:float) -> float`

EN: Convert kilograms to pounds.

ES: Convierte kilogramos a libras.

Native example flag: false.

```tripol
import tripol.convert

show convert.kg_to_lb(0.45359237)
```

Documented expected output:

```text
1.000000
```

### convert.lb_to_kg

`func convert.lb_to_kg(l:float) -> float`

EN: Convert pounds to kilograms.

ES: Convierte libras a kilogramos.

Native example flag: false.

```tripol
import tripol.convert

show convert.lb_to_kg(1.0)
```

Documented expected output:

```text
0.453592
```

### convert.celsius_to_kelvin

`func convert.celsius_to_kelvin(c:float) -> float`

EN: Convert Celsius to kelvin.

ES: Convierte grados Celsius a kelvin.

Native example flag: false.

```tripol
import tripol.convert

show convert.celsius_to_kelvin(0.0)
```

Documented expected output:

```text
273.150000
```

### convert.kelvin_to_celsius

`func convert.kelvin_to_celsius(k:float) -> float`

EN: Convert kelvin to Celsius.

ES: Convierte kelvin a grados Celsius.

Native example flag: false.

```tripol
import tripol.convert

show convert.kelvin_to_celsius(273.15)
```

Documented expected output:

```text
0.000000
```

---

## tripol.csv — 3 functions

Import: `import tripol.csv`

### csv.escape

`func csv.escape(field:string) -> string`

EN: Quote one CSV field and double any embedded double quote.

ES: Entrecomilla un campo CSV y duplica las comillas dobles internas.

Native example flag: false.

```tripol
import tripol.csv

show csv.escape("a,b")
```

Documented expected output:

```text
"a,b"
```

### csv.fields

`func csv.fields(record:string) -> int`

EN: Count fields in one comma-separated record, returning -1 for invalid quoting.

ES: Cuenta los campos de un registro separado por comas y devuelve -1 si las comillas son inválidas.

Native example flag: false.

```tripol
import tripol.csv

show csv.fields("\"a,b\",c")
```

Documented expected output:

```text
2
```

### csv.valid

`func csv.valid(record:string) -> bool`

EN: Validate a single CSV record, including quoted fields.

ES: Valida un registro CSV, incluidos los campos entrecomillados.

Native example flag: false.

```tripol
import tripol.csv

show csv.valid("a,b,c")
```

Documented expected output:

```text
true
```

---

## tripol.datetime — 7 functions

Import: `import tripol.datetime`

### datetime.leap

`func datetime.leap(y:int) -> bool`

EN: Check whether a year has 366 days.

ES: Comprueba si un año es bisiesto y tiene 366 días.

Native example flag: false.

```tripol
import tripol.datetime

show datetime.leap(2024)
```

Documented expected output:

```text
true
```

### datetime.days_in_month

`func datetime.days_in_month(y:int,m:int) -> int`

EN: Count the days in a month and year.

ES: Devuelve cuántos días tiene un mes de un año concreto.

Native example flag: false.

```tripol
import tripol.datetime

show datetime.days_in_month(2024, 2)
```

Documented expected output:

```text
29
```

### datetime.valid

`func datetime.valid(y:int,m:int,d:int) -> bool`

EN: Check whether a date exists.

ES: Comprueba si una fecha es válida.

Native example flag: false.

```tripol
import tripol.datetime

show datetime.valid(2024, 2, 29)
```

Documented expected output:

```text
true
```

### datetime.day_of_year

`func datetime.day_of_year(y:int,m:int,d:int) -> int`

EN: Find the position of a day within its year.

ES: Devuelve la posición de un día dentro de su año.

Native example flag: false.

```tripol
import tripol.datetime

show datetime.day_of_year(2024, 2, 1)
```

Documented expected output:

```text
32
```

### datetime.ordinal

`func datetime.ordinal(y:int,m:int,d:int) -> int`

EN: Return the day number counted from year 1.

ES: Devuelve el número de día contado desde el año 1.

Native example flag: false.

Note: January 1 of year 1 is day 1; Gregorian years 1–9999.

Nota: El 1 de enero del año 1 es el día 1; años gregorianos de 1 a 9999.

```tripol
import tripol.datetime

show datetime.ordinal(1, 1, 1)
show datetime.ordinal(1, 1, 2)
```

Documented expected output:

```text
1
2
```

### datetime.weekday

`func datetime.weekday(y:int,m:int,d:int) -> int`

EN: Return a weekday number, with Monday = 1.

ES: Devuelve el número del día de la semana, con lunes = 1.

Native example flag: false.

```tripol
import tripol.datetime

show datetime.weekday(2024, 1, 1)
```

Documented expected output:

```text
1
```

### datetime.days_between

`func datetime.days_between(y1:int,m1:int,d1:int,y2:int,m2:int,d2:int) -> int`

EN: Calculate the signed difference in days between dates.

ES: Calcula la diferencia con signo en días entre fechas.

Native example flag: false.

```tripol
import tripol.datetime

show datetime.days_between(2024, 1, 1, 2024, 1, 3)
```

Documented expected output:

```text
2
```

---

## tripol.errors — 7 functions

Import: `import tripol.errors`

### errors.is_int

`func errors.is_int(s:string) -> bool`

EN: Check whether text represents a valid signed 64-bit integer.

ES: Comprueba si el texto representa un entero válido de 64 bits con signo.

Native example flag: false.

```tripol
import tripol.errors

show errors.is_int("-42")
show errors.is_int("abc")
```

Documented expected output:

```text
true
false
```

### errors.int_or

`func errors.int_or(s:string,fallback:int) -> int`

EN: Convert text to an integer or return a fallback.

ES: Convierte texto a entero o devuelve un valor alternativo.

Native example flag: false.

```tripol
import tripol.errors

show errors.int_or("abc", 7)
```

Documented expected output:

```text
7
```

### errors.can_divide

`func errors.can_divide(a:int,b:int) -> bool`

EN: Check whether integer division avoids zero and overflow.

ES: Comprueba si una división entera evita el cero y el desbordamiento.

Native example flag: false.

```tripol
import tripol.errors

show errors.can_divide(10, 0)
```

Documented expected output:

```text
false
```

### errors.div_or

`func errors.div_or(a:int,b:int,fallback:int) -> int`

EN: Divide safely or return a fallback.

ES: Divide de forma segura o devuelve un valor alternativo.

Native example flag: false.

```tripol
import tripol.errors

show errors.div_or(10, 0, -1)
```

Documented expected output:

```text
-1
```

### errors.get_or

`func errors.get_or(a:ints,index:int,fallback:int) -> int`

EN: Read an array element or return a fallback for an invalid index.

ES: Lee un elemento o devuelve un valor alternativo si el índice no es válido.

Native example flag: false.

```tripol
import tripol.errors

a = [10, 20]
show errors.get_or(a, 5, -1)
```

Documented expected output:

```text
-1
```

### errors.message

`func errors.message(code:int) -> string`

EN: Get the English description of an exception code.

ES: Obtiene la descripción en inglés de un código de excepción.

Native example flag: false.

```tripol
import tripol.errors

show errors.message(11)
```

Documented expected output:

```text
Index out of bounds
```

### errors.message_es

`func errors.message_es(code:int) -> string`

EN: Get the Spanish description of an exception code.

ES: Obtiene la descripción en español de un código de excepción.

Native example flag: false.

```tripol
import tripol.errors

show errors.message_es(11)
```

Documented expected output:

```text
Indice fuera de rango
```

---

## tripol.files — 5 functions

Import: `import tripol.files`

### files.read

`func files.read(path:string) -> string`

EN: Read text from a file.

ES: Lee texto de un archivo.

Native example flag: false.

```tripol
import tripol.files

files.write("demo.txt", "Hola")
show files.read("demo.txt")
```

Documented expected output:

```text
Hola
```

### files.write

`func files.write(path:string, content:string) -> int`

EN: Write text, replacing previous content.

ES: Escribe texto, reemplazando el contenido anterior.

Native example flag: false.

```tripol
import tripol.files

files.write("demo.txt", "Hola")
show files.read("demo.txt")
```

Documented expected output:

```text
Hola
```

### files.append

`func files.append(path:string, content:string) -> int`

EN: Add text at the end of a file.

ES: Añade texto al final de un archivo.

Native example flag: false.

```tripol
import tripol.files

files.write("demo.txt", "Hola")
files.append("demo.txt", " TRIPOL")
show files.read("demo.txt")
```

Documented expected output:

```text
Hola TRIPOL
```

### files.write_line

`func files.write_line(path:string, content:string) -> int`

EN: Replace file contents with a line ending in a newline.

ES: Reemplaza el contenido con una línea terminada en salto de línea.

Native example flag: false.

```tripol
import tripol.files

files.write_line("demo.txt", "Hola")
show len(files.read("demo.txt"))
```

Documented expected output:

```text
5
```

### files.append_line

`func files.append_line(path:string, content:string) -> int`

EN: Append a line ending in a newline.

ES: Añade una línea terminada en salto de línea.

Native example flag: false.

```tripol
import tripol.files

files.write_line("demo.txt", "Uno")
files.append_line("demo.txt", "Dos")
show len(files.read("demo.txt"))
```

Documented expected output:

```text
8
```

---

## tripol.math — 12 functions

Import: `import tripol.math`

### math.square

`func math.square(n:int) -> int`

EN: Multiply a number by itself.

ES: Multiplica un número por sí mismo.

Native example flag: false.

```tripol
import tripol.math

show math.square(5)
```

Documented expected output:

```text
25
```

### math.cube

`func math.cube(n:int) -> int`

EN: Multiply a number by itself three times.

ES: Multiplica un número por sí mismo tres veces.

Native example flag: false.

```tripol
import tripol.math

show math.cube(3)
```

Documented expected output:

```text
27
```

### math.clamp

`func math.clamp(n:int, low:int, high:int) -> int`

EN: Keep a number within two limits.

ES: Mantiene un número entre dos límites.

Native example flag: false.

Note: Use a lower limit no greater than the upper limit.

Nota: El límite inferior no debe superar el superior.

```tripol
import tripol.math

show math.clamp(120, 0, 100)
```

Documented expected output:

```text
100
```

### math.power

`func math.power(base:int, exponent:int) -> int`

EN: Raise an integer to a nonnegative exponent.

ES: Eleva un entero a un exponente no negativo.

Native example flag: false.

```tripol
import tripol.math

show math.power(2, 4)
```

Documented expected output:

```text
16
```

### math.gcd

`func math.gcd(a:int, b:int) -> int`

EN: Find the greatest common divisor of two integers.

ES: Calcula el máximo común divisor de dos enteros.

Native example flag: false.

```tripol
import tripol.math

show math.gcd(12, 18)
```

Documented expected output:

```text
6
```

### math.root

`func math.root(x:float) -> float`

EN: Calculate the square root of a nonnegative decimal.

ES: Calcula la raíz cuadrada de un decimal no negativo.

Native example flag: false.

```tripol
import tripol.math

show math.root(81.0)
```

Documented expected output:

```text
9.000000
```

### math.pi

`func math.pi() -> float`

EN: Return the mathematical constant pi.

ES: Devuelve la constante matemática pi.

Native example flag: false.

```tripol
import tripol.math

show math.pi()
```

Documented expected output:

```text
3.141593
```

### math.percent

`func math.percent(value:float, rate:float) -> float`

EN: Calculate a percentage of a number.

ES: Calcula un porcentaje de un número.

Native example flag: false.

```tripol
import tripol.math

show math.percent(200.0, 10.0)
```

Documented expected output:

```text
20.000000
```

### math.sqrt_value

`func math.sqrt_value(x:float) -> float`

EN: Calculate a square root, just like math.root.

ES: Calcula una raíz cuadrada, igual que math.root.

Native example flag: false.

```tripol
import tripol.math

show math.sqrt_value(25.0)
```

Documented expected output:

```text
5.000000
```

### math.absolute

`func math.absolute(n:int) -> int`

EN: Return the distance of an integer from zero.

ES: Devuelve la distancia de un entero respecto a cero.

Native example flag: false.

```tripol
import tripol.math

show math.absolute(-8)
```

Documented expected output:

```text
8
```

### math.minimum

`func math.minimum(a:int,b:int) -> int`

EN: Return the smaller integer.

ES: Devuelve el menor de dos enteros.

Native example flag: false.

```tripol
import tripol.math

show math.minimum(3, 8)
```

Documented expected output:

```text
3
```

### math.maximum

`func math.maximum(a:int,b:int) -> int`

EN: Return the larger integer.

ES: Devuelve el mayor de dos enteros.

Native example flag: false.

```tripol
import tripol.math

show math.maximum(3, 8)
```

Documented expected output:

```text
8
```

---

## tripol.path — 6 functions

Import: `import tripol.path`

### path.separator

`func path.separator(b:int) -> bool`

EN: Check whether a byte is a slash or backslash.

ES: Comprueba si un byte representa una barra normal o invertida.

Native example flag: false.

```tripol
import tripol.path

show path.separator(47)
show path.separator(92)
```

Documented expected output:

```text
true
true
```

### path.last_separator

`func path.last_separator(s:string) -> int`

EN: Find the final separator's byte index, or -1.

ES: Busca el índice por bytes del último separador, o -1.

Native example flag: false.

```tripol
import tripol.path

show path.last_separator("a/b.txt")
```

Documented expected output:

```text
1
```

### path.basename

`func path.basename(s:string) -> string`

EN: Return the last path component.

ES: Devuelve el último componente de una ruta.

Native example flag: false.

```tripol
import tripol.path

show path.basename("docs/main.tri")
```

Documented expected output:

```text
main.tri
```

### path.dirname

`func path.dirname(s:string) -> string`

EN: Return the path before the last separator.

ES: Devuelve la ruta anterior al último separador.

Native example flag: false.

```tripol
import tripol.path

show path.dirname("docs/main.tri")
```

Documented expected output:

```text
docs
```

### path.extension

`func path.extension(s:string) -> string`

EN: Return the last filename extension, including its dot.

ES: Devuelve la última extensión del nombre, incluyendo el punto.

Native example flag: false.

```tripol
import tripol.path

show path.extension("main.tri")
```

Documented expected output:

```text
.tri
```

### path.join

`func path.join(a:string,b:string) -> string`

EN: Join path text with a slash when needed.

ES: Une textos de ruta con una barra cuando hace falta.

Native example flag: false.

Note: Lexical only: does not resolve traversal or check the filesystem.

Nota: Solo manipula texto: no resuelve rutas relativas ni comprueba el sistema de archivos.

```tripol
import tripol.path

show path.join("docs", "main.tri")
```

Documented expected output:

```text
docs/main.tri
```

---

## tripol.random — 3 functions

Import: `import tripol.random`

### random.next

`func random.next(seed:int) -> int`

EN: Advance the deterministic pseudorandom seed.

ES: Avanza la semilla pseudoaleatoria determinista.

Native example flag: false.

Note: Not cryptographically secure; seed must be 1–2147483646.

Nota: No es criptográficamente seguro; semilla de 1 a 2147483646.

```tripol
import tripol.random

show random.next(1)
```

Documented expected output:

```text
48271
```

### random.unit

`func random.unit(seed:int) -> float`

EN: Generate a deterministic decimal between zero and one.

ES: Genera un decimal determinista entre cero y uno.

Native example flag: false.

```tripol
import tripol.random

show random.unit(1) > 0.0
```

Documented expected output:

```text
true
```

### random.bounded

`func random.bounded(seed:int,limit:int) -> int`

EN: Generate a deterministic integer from 0 to limit-1.

ES: Genera un entero determinista desde 0 hasta limit-1.

Native example flag: false.

```tripol
import tripol.random

show random.bounded(1, 10)
```

Documented expected output:

```text
0
```

---

## tripol.search — 4 functions

Import: `import tripol.search`

### search.linear

`func search.linear(a:ints, value:int) -> int`

EN: Find the first occurrence of a value, or -1.

ES: Busca la primera aparición de un valor, o -1.

Native example flag: false.

```tripol
import tripol.search

a = [8, 3, 8]
show search.linear(a, 8)
```

Documented expected output:

```text
0
```

### search.lower_bound

`func search.lower_bound(a:ints, value:int) -> int`

EN: Find the insertion index in an ascending sorted array.

ES: Busca el índice de inserción en un arreglo ordenado ascendentemente.

Native example flag: false.

```tripol
import tripol.search

a = [1, 3, 5]
show search.lower_bound(a, 4)
```

Documented expected output:

```text
2
```

### search.binary

`func search.binary(a:ints, value:int) -> int`

EN: Find a value in an ascending sorted array, or -1.

ES: Busca un valor en un arreglo ordenado ascendentemente, o -1.

Native example flag: false.

```tripol
import tripol.search

a = [1, 3, 5]
show search.binary(a, 5)
```

Documented expected output:

```text
2
```

### search.last

`func search.last(a:ints,value:int) -> int`

EN: Find the last occurrence of a value, or -1.

ES: Busca la última aparición de un valor, o -1.

Native example flag: false.

```tripol
import tripol.search

a = [8, 3, 8]
show search.last(a, 8)
```

Documented expected output:

```text
2
```

---

## tripol.sort — 4 functions

Import: `import tripol.sort`

### sort.sift

`func sort.sift(a:ints, root:int, size:int) -> int`

EN: Restore the max-heap property below a selected root in a dynamic array.

ES: Restaura la propiedad de montículo máximo bajo una raíz en un arreglo dinámico.

Native example flag: false.

Note: Advanced helper; use ascending() for ordinary sorting.

Nota: Función avanzada; usa ascending() para ordenar normalmente.

```tripol
import tripol.sort

a = ints(3)
a[0] = 1
a[1] = 3
a[2] = 2
sort.sift(a, 0, 3)
show a
```

Documented expected output:

```text
[3, 1, 2]
```

### sort.ascending

`func sort.ascending(a:ints) -> ints`

EN: Return a sorted copy from smallest to largest.

ES: Devuelve una copia ordenada de menor a mayor.

Native example flag: false.

```tripol
import tripol.sort

a = [3, 1, 2]
show sort.ascending(a)
```

Documented expected output:

```text
[1, 2, 3]
```

### sort.descending

`func sort.descending(a:ints) -> ints`

EN: Return a sorted copy from largest to smallest.

ES: Devuelve una copia ordenada de mayor a menor.

Native example flag: false.

```tripol
import tripol.sort

a = [3, 1, 2]
show sort.descending(a)
```

Documented expected output:

```text
[3, 2, 1]
```

### sort.sorted

`func sort.sorted(a:ints) -> bool`

EN: Check whether an array is sorted in ascending order.

ES: Comprueba si un arreglo está ordenado ascendentemente.

Native example flag: false.

```tripol
import tripol.sort

a = [1, 2, 3]
show sort.sorted(a)
```

Documented expected output:

```text
true
```

---

## tripol.text — 9 functions

Import: `import tripol.text`

### text.size

`func text.size(s:string) -> int`

EN: Count UTF-8 bytes, not necessarily visible characters.

ES: Cuenta bytes UTF-8, no necesariamente caracteres visibles.

Native example flag: false.

```tripol
import tripol.text

show text.size("TRIPOL")
```

Documented expected output:

```text
6
```

### text.empty

`func text.empty(s:string) -> bool`

EN: Check whether a string is empty.

ES: Comprueba si una cadena está vacía.

Native example flag: false.

```tripol
import tripol.text

show text.empty("")
```

Documented expected output:

```text
true
```

### text.from_int

`func text.from_int(n:int) -> string`

EN: Convert an integer to text.

ES: Convierte un entero en texto.

Native example flag: false.

```tripol
import tripol.text

show text.from_int(42)
```

Documented expected output:

```text
42
```

### text.to_int

`func text.to_int(s:string) -> int`

EN: Convert valid integer text into a number.

ES: Convierte texto entero válido en un número.

Native example flag: false.

```tripol
import tripol.text

show text.to_int("42") + 1
```

Documented expected output:

```text
43
```

### text.part

`func text.part(s:string, start:int, end:int) -> string`

EN: Copy bytes from start (included) to end (excluded).

ES: Copia bytes desde start (incluido) hasta end (excluido).

Native example flag: false.

Note: UTF-8 characters may use several bytes.

Nota: Los caracteres UTF-8 pueden ocupar varios bytes.

```tripol
import tripol.text

show text.part("TRIPOL", 0, 3)
```

Documented expected output:

```text
TRI
```

### text.starts_with

`func text.starts_with(s:string, prefix:string) -> bool`

EN: Check whether text begins with a prefix.

ES: Comprueba si un texto comienza por un prefijo.

Native example flag: false.

```tripol
import tripol.text

show text.starts_with("TRIPOL", "TRI")
```

Documented expected output:

```text
true
```

### text.ends_with

`func text.ends_with(s:string, suffix:string) -> bool`

EN: Check whether text ends with a suffix.

ES: Comprueba si un texto termina con un sufijo.

Native example flag: false.

```tripol
import tripol.text

show text.ends_with("main.tri", ".tri")
```

Documented expected output:

```text
true
```

### text.index_of

`func text.index_of(s:string,needle:string) -> int`

EN: Find the first byte index of a substring, or -1.

ES: Encuentra el primer índice por bytes de un fragmento, o -1.

Native example flag: false.

```tripol
import tripol.text

show text.index_of("TRIPOL", "POL")
```

Documented expected output:

```text
3
```

### text.contains

`func text.contains(s:string,needle:string) -> bool`

EN: Check whether text contains a substring.

ES: Comprueba si un texto contiene un fragmento.

Native example flag: false.

```tripol
import tripol.text

show text.contains("Hello TRIPOL", "TRIPOL")
```

Documented expected output:

```text
true
```

---

## tripol.validate — 6 functions

Import: `import tripol.validate`

### validate.digit

`func validate.digit(n:int) -> bool`

EN: Check whether an ASCII byte represents a digit.

ES: Comprueba si un byte ASCII representa un dígito.

Native example flag: false.

```tripol
import tripol.validate

show validate.digit(57)
```

Documented expected output:

```text
true
```

### validate.letter

`func validate.letter(n:int) -> bool`

EN: Check whether an ASCII byte represents an English letter.

ES: Comprueba si un byte ASCII representa una letra inglesa.

Native example flag: false.

```tripol
import tripol.validate

show validate.letter(65)
```

Documented expected output:

```text
true
```

### validate.whitespace

`func validate.whitespace(n:int) -> bool`

EN: Check for space, tab, line feed or carriage return.

ES: Comprueba si es espacio, tabulador, salto de línea o retorno de carro.

Native example flag: false.

```tripol
import tripol.validate

show validate.whitespace(32)
```

Documented expected output:

```text
true
```

### validate.integer

`func validate.integer(s:string) -> bool`

EN: Check the spelling of a signed integer, without checking its numeric range.

ES: Comprueba la escritura de un entero con signo, sin comprobar su rango numérico.

Native example flag: false.

Note: Use errors.is_int() to check the full int64 range.

Nota: Usa errors.is_int() para comprobar todo el rango int64.

```tripol
import tripol.validate

show validate.integer("-42")
show validate.integer("3.5")
```

Documented expected output:

```text
true
false
```

### validate.ascii

`func validate.ascii(s:string) -> bool`

EN: Check whether every byte is in the ASCII range.

ES: Comprueba si todos los bytes pertenecen al rango ASCII.

Native example flag: false.

```tripol
import tripol.validate

show validate.ascii("TRIPOL")
show validate.ascii("á")
```

Documented expected output:

```text
true
false
```

### validate.identifier

`func validate.identifier(s:string) -> bool`

EN: Check an ASCII identifier beginning with a letter or underscore.

ES: Comprueba un identificador ASCII que comienza por letra o guion bajo.

Native example flag: false.

Note: Does not check whether the name is a reserved keyword.

Nota: No comprueba si el nombre es una palabra reservada.

```tripol
import tripol.validate

show validate.identifier("player_1")
show validate.identifier("1player")
```

Documented expected output:

```text
true
false
```

---

## tripol.animation — 6 functions

Import: `import tripol.animation`

### animation.clamp01

`func animation.clamp01(t:float) -> float`

EN: Limit a decimal to the range 0–1.

ES: Limita un decimal al intervalo de 0 a 1.

Native example flag: false.

```tripol
import tripol.animation

show animation.clamp01(1.5)
```

Documented expected output:

```text
1.000000
```

### animation.linear

`func animation.linear(a:float,b:float,t:float) -> float`

EN: Interpolate between two values with progress from 0 to 1.

ES: Interpola entre dos valores con un progreso de 0 a 1.

Native example flag: false.

```tripol
import tripol.animation

show animation.linear(0.0, 10.0, 0.5)
```

Documented expected output:

```text
5.000000
```

### animation.ease_in

`func animation.ease_in(t:float) -> float`

EN: Start an animation slowly and accelerate.

ES: Inicia una animación lentamente y acelera.

Native example flag: false.

```tripol
import tripol.animation

show animation.ease_in(0.5)
```

Documented expected output:

```text
0.250000
```

### animation.ease_out

`func animation.ease_out(t:float) -> float`

EN: Start quickly and slow down near the end.

ES: Comienza rápido y desacelera al final.

Native example flag: false.

```tripol
import tripol.animation

show animation.ease_out(0.5)
```

Documented expected output:

```text
0.750000
```

### animation.smoothstep

`func animation.smoothstep(t:float) -> float`

EN: Smooth the beginning and end of a transition.

ES: Suaviza el inicio y el final de una transición.

Native example flag: false.

```tripol
import tripol.animation

show animation.smoothstep(0.5)
```

Documented expected output:

```text
0.500000
```

### animation.frame

`func animation.frame(elapsed:int,frame_ms:int,count:int) -> int`

EN: Find the current looping frame from elapsed milliseconds.

ES: Calcula el fotograma actual de una animación repetitiva a partir de milisegundos.

Native example flag: false.

Note: Frame duration and count must be positive.

Nota: La duración y cantidad de fotogramas deben ser positivas.

```tripol
import tripol.animation

show animation.frame(250, 100, 4)
```

Documented expected output:

```text
2
```

---

## tripol.buffer — 8 functions

Import: `import tripol.buffer`

### buffer.create

`func buffer.create(initial=4) -> ints`

EN: Create an empty growable integer buffer with an initial capacity.

ES: Crea un búfer dinámico de enteros vacío con capacidad inicial.

Native example flag: false.

```tripol
import tripol.buffer

b = buffer.create(2)
show buffer.size(b)
```

Documented expected output:

```text
0
```

### buffer.capacity

`func buffer.capacity(b:ints) -> int`

EN: Return the number of storage positions currently allocated.

ES: Devuelve cuántas posiciones de almacenamiento están reservadas.

Native example flag: false.

```tripol
import tripol.buffer

b = buffer.create(4)
show buffer.capacity(b)
```

Documented expected output:

```text
4
```

### buffer.size

`func buffer.size(b:ints) -> int`

EN: Return the number of stored values.

ES: Devuelve cuántos valores contiene el búfer.

Native example flag: false.

```tripol
import tripol.buffer

b = buffer.create()
b = buffer.push(b, 10)
show buffer.size(b)
```

Documented expected output:

```text
1
```

### buffer.push

`func buffer.push(b:ints,value:int) -> ints`

EN: Append an integer and grow storage when necessary.

ES: Añade un entero y amplía el almacenamiento si hace falta.

Native example flag: false.

Note: Always reassign the returned buffer.

Nota: Reasigna siempre el búfer devuelto.

```tripol
import tripol.buffer

b = buffer.create(1)
b = buffer.push(b, 10)
b = buffer.push(b, 20)
show buffer.values(b)
```

Documented expected output:

```text
[10, 20]
```

### buffer.get

`func buffer.get(b:ints,index:int) -> int`

EN: Read an element at a zero-based index.

ES: Lee un elemento por su índice, empezando en cero.

Native example flag: false.

```tripol
import tripol.buffer

b = buffer.create()
b = buffer.push(b, 42)
show buffer.get(b, 0)
```

Documented expected output:

```text
42
```

### buffer.set

`func buffer.set(b:ints,index:int,value:int) -> int`

EN: Replace an existing element and return the new value.

ES: Reemplaza un elemento existente y devuelve el nuevo valor.

Native example flag: false.

```tripol
import tripol.buffer

b = buffer.create()
b = buffer.push(b, 1)
show buffer.set(b, 0, 9)
```

Documented expected output:

```text
9
```

### buffer.pop

`func buffer.pop(b:ints) -> int`

EN: Remove and return the last stored integer.

ES: Elimina y devuelve el último entero guardado.

Native example flag: false.

Note: The buffer must not be empty.

Nota: El búfer no debe estar vacío.

```tripol
import tripol.buffer

b = buffer.create()
b = buffer.push(b, 7)
show buffer.pop(b)
```

Documented expected output:

```text
7
```

### buffer.values

`func buffer.values(b:ints) -> ints`

EN: Copy the logical contents into a new array.

ES: Copia el contenido lógico a un arreglo nuevo.

Native example flag: false.

```tripol
import tripol.buffer

b = buffer.create()
b = buffer.push(b, 5)
show buffer.values(b)
```

Documented expected output:

```text
[5]
```

---

## tripol.camera — 4 functions

Import: `import tripol.camera`

### camera.screen

`func camera.screen(world:int,origin:int,zoom=1) -> int`

EN: Convert a world coordinate to a screen coordinate.

ES: Convierte una coordenada del mundo en una de pantalla.

Native example flag: false.

```tripol
import tripol.camera

show camera.screen(12, 2, 3)
```

Documented expected output:

```text
30
```

### camera.world

`func camera.world(screen_pos:int,origin:int,zoom=1) -> int`

EN: Convert a screen coordinate back into world space using integer division.

ES: Convierte una coordenada de pantalla al mundo con división entera.

Native example flag: false.

```tripol
import tripol.camera

show camera.world(30, 2, 3)
```

Documented expected output:

```text
12
```

### camera.follow

`func camera.follow(pos:int,view:int,limit:int) -> int`

EN: Center a camera on a position while staying inside world bounds.

ES: Centra la cámara en una posición sin salir de los límites del mundo.

Native example flag: false.

```tripol
import tripol.camera

show camera.follow(50, 20, 100)
```

Documented expected output:

```text
40
```

### camera.visible

`func camera.visible(x:int,y:int,w:int,h:int,vw:int,vh:int) -> bool`

EN: Check whether a rectangle overlaps the visible viewport.

ES: Comprueba si un rectángulo se cruza con el área visible.

Native example flag: false.

```tripol
import tripol.camera

show camera.visible(5, 5, 10, 10, 100, 100)
```

Documented expected output:

```text
true
```

---

## tripol.canvas — 10 functions

Import: `import tripol.canvas`

### canvas.size

`func canvas.size(w:int,h:int) -> int`

EN: Calculate the number of pixels in a valid canvas.

ES: Calcula la cantidad de píxeles de un lienzo válido.

Native example flag: false.

```tripol
import tripol.canvas

show canvas.size(4, 3)
```

Documented expected output:

```text
12
```

### canvas.create

`func canvas.create(w:int,h:int) -> ints`

EN: Allocate a zero-filled RGB framebuffer.

ES: Reserva un arreglo de píxeles RGB inicializado en cero.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(2, 2)
show p
```

Documented expected output:

```text
[0, 0, 0, 0]
```

### canvas.clear

`func canvas.clear(pixels:ints,color:int) -> int`

EN: Fill every pixel with a packed RGB color.

ES: Rellena todos los píxeles con un color RGB empaquetado.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(2, 2)
canvas.clear(p, 255)
show p
```

Documented expected output:

```text
[255, 255, 255, 255]
```

### canvas.pixel

`func canvas.pixel(pixels:ints,w:int,h:int,x:int,y:int,color:int) -> bool`

EN: Write one pixel; return false if outside the canvas.

ES: Escribe un píxel; devuelve false si está fuera del lienzo.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(2, 2)
show canvas.pixel(p, 2, 2, 1, 0, 255)
show p
```

Documented expected output:

```text
true
[0, 255, 0, 0]
```

### canvas.get

`func canvas.get(pixels:ints,w:int,h:int,x:int,y:int) -> int`

EN: Read the packed color at a pixel coordinate.

ES: Lee el color empaquetado de un píxel.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(2, 2)
canvas.pixel(p, 2, 2, 1, 1, 255)
show canvas.get(p, 2, 2, 1, 1)
```

Documented expected output:

```text
255
```

### canvas.rectangle

`func canvas.rectangle(pixels:ints,w:int,h:int,x:int,y:int,rw:int,rh:int,color:int) -> int`

EN: Draw a filled rectangle clipped to the canvas.

ES: Dibuja un rectángulo relleno recortado a los límites del lienzo.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(3, 2)
show canvas.rectangle(p, 3, 2, 1, 0, 2, 1, 255)
show p
```

Documented expected output:

```text
2
[0, 255, 255, 0, 0, 0]
```

### canvas.line

`func canvas.line(p:ints,w:int,h:int,x:int,y:int,ex:int,ey:int,color:int) -> int`

EN: Draw a line and return the number of pixels written.

ES: Dibuja una línea y devuelve cuántos píxeles escribió.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(3, 3)
show canvas.line(p, 3, 3, 0, 0, 2, 2, 255)
```

Documented expected output:

```text
3
```

### canvas.circle

`func canvas.circle(p:ints,w:int,h:int,x:int,y:int,r:int,color:int) -> int`

EN: Draw a filled disk and return the number of visible pixels written.

ES: Dibuja un círculo relleno y devuelve cuántos píxeles visibles escribió.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(3, 3)
show canvas.circle(p, 3, 3, 1, 1, 1, 255)
```

Documented expected output:

```text
5
```

### canvas.blend

`func canvas.blend(p:ints,w:int,h:int,x:int,y:int,color:int,opacity:int) -> bool`

EN: Mix a color into one pixel using opacity from 0 to 100.

ES: Mezcla un color con un píxel usando opacidad de 0 a 100.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(1, 1)
show canvas.blend(p, 1, 1, 0, 0, 255, 50)
show p
```

Documented expected output:

```text
true
[127]
```

### canvas.tile_map

`func canvas.tile_map(dst:ints,w:int,h:int,tiles:ints,cols:int,rows:int,palette:ints,side:int,ox:int,oy:int) -> int`

EN: Draw visible palette-indexed tiles; tile zero is transparent.

ES: Dibuja mosaicos visibles mediante índices de paleta; el mosaico cero es transparente.

Native example flag: false.

```tripol
import tripol.canvas

p = canvas.create(2, 1)
tiles = [0, 1]
palette = [0, 255]
show canvas.tile_map(p, 2, 1, tiles, 2, 1, palette, 1, 0, 0)
show p
```

Documented expected output:

```text
1
[0, 255]
```

---

## tripol.colors — 8 functions

Import: `import tripol.colors`

### colors.rgb

`func colors.rgb(r:int,g:int,b:int) -> int`

EN: Pack red, green and blue channels into one RGB24 integer.

ES: Empaqueta los canales rojo, verde y azul en un entero RGB24.

Native example flag: false.

```tripol
import tripol.colors

show colors.rgb(255, 0, 0)
```

Documented expected output:

```text
16711680
```

### colors.valid

`func colors.valid(c:int) -> bool`

EN: Check whether a number is a valid RGB24 color.

ES: Comprueba si un número representa un color RGB24 válido.

Native example flag: false.

```tripol
import tripol.colors

show colors.valid(16777215)
```

Documented expected output:

```text
true
```

### colors.red

`func colors.red(c:int) -> int`

EN: Extract the red channel.

ES: Extrae el canal rojo.

Native example flag: false.

```tripol
import tripol.colors

show colors.red(colors.rgb(10, 20, 30))
```

Documented expected output:

```text
10
```

### colors.green

`func colors.green(c:int) -> int`

EN: Extract the green channel.

ES: Extrae el canal verde.

Native example flag: false.

```tripol
import tripol.colors

show colors.green(colors.rgb(10, 20, 30))
```

Documented expected output:

```text
20
```

### colors.blue

`func colors.blue(c:int) -> int`

EN: Extract the blue channel.

ES: Extrae el canal azul.

Native example flag: false.

```tripol
import tripol.colors

show colors.blue(colors.rgb(10, 20, 30))
```

Documented expected output:

```text
30
```

### colors.invert

`func colors.invert(c:int) -> int`

EN: Invert all three RGB channels.

ES: Invierte los tres canales RGB.

Native example flag: false.

```tripol
import tripol.colors

show colors.invert(0)
```

Documented expected output:

```text
16777215
```

### colors.gray

`func colors.gray(v:int) -> int`

EN: Create a gray color with equal RGB channels.

ES: Crea un color gris con los tres canales iguales.

Native example flag: false.

```tripol
import tripol.colors

show colors.gray(1)
```

Documented expected output:

```text
65793
```

### colors.mix

`func colors.mix(a:int,b:int,percent:int) -> int`

EN: Blend two colors using an integer percentage from 0 to 100.

ES: Mezcla dos colores usando un porcentaje entero de 0 a 100.

Native example flag: false.

```tripol
import tripol.colors

show colors.mix(0, 255, 50)
```

Documented expected output:

```text
127
```

---

## tripol.dictionary — 8 functions

Import: `import tripol.dictionary`

### dictionary.create

`func dictionary.create(capacity=32) -> ints`

EN: Create a fixed-capacity map from text keys to integers.

ES: Crea un mapa de capacidad fija con claves de texto y valores enteros.

Native example flag: false.

```tripol
import tripol.dictionary

d = dictionary.create(2)
show dictionary.size(d)
```

Documented expected output:

```text
0
```

### dictionary.capacity

`func dictionary.capacity(d:ints) -> int`

EN: Return the maximum number of entries.

ES: Devuelve la cantidad máxima de entradas.

Native example flag: false.

```tripol
import tripol.dictionary

d = dictionary.create(4)
show dictionary.capacity(d)
```

Documented expected output:

```text
4
```

### dictionary.locate

`func dictionary.locate(d:ints,key:string) -> int`

EN: Find the internal slot offset of a key, or -1.

ES: Busca el desplazamiento interno de una clave, o -1.

Native example flag: false.

Note: Prefer get() or contains() for everyday use.

Nota: Prefiere get() o contains() en el uso cotidiano.

```tripol
import tripol.dictionary

d = dictionary.create(2)
show dictionary.locate(d, "score")
```

Documented expected output:

```text
-1
```

### dictionary.contains

`func dictionary.contains(d:ints,key:string) -> bool`

EN: Check if a key exists.

ES: Comprueba si existe una clave.

Native example flag: false.

```tripol
import tripol.dictionary

d = dictionary.create(2)
dictionary.set(d, "score", 42)
show dictionary.contains(d, "score")
```

Documented expected output:

```text
true
```

### dictionary.get

`func dictionary.get(d:ints,key:string,fallback=0) -> int`

EN: Read an integer by key, returning a fallback if missing.

ES: Lee un entero por su clave, devolviendo un valor alternativo si falta.

Native example flag: false.

```tripol
import tripol.dictionary

d = dictionary.create(2)
show dictionary.get(d, "score", -1)
```

Documented expected output:

```text
-1
```

### dictionary.size

`func dictionary.size(d:ints) -> int`

EN: Count occupied entries.

ES: Cuenta las entradas ocupadas.

Native example flag: false.

```tripol
import tripol.dictionary

d = dictionary.create(2)
dictionary.set(d, "score", 42)
show dictionary.size(d)
```

Documented expected output:

```text
1
```

### dictionary.set

`func dictionary.set(d:ints,key:string,value:int) -> bool`

EN: Insert or update an integer; return false if full and the key is new.

ES: Inserta o actualiza un entero; devuelve false si está lleno y la clave es nueva.

Native example flag: false.

Note: Keys must contain 1–32 UTF-8 bytes.

Nota: Las claves deben tener entre 1 y 32 bytes UTF-8.

```tripol
import tripol.dictionary

d = dictionary.create(2)
show dictionary.set(d, "score", 42)
show dictionary.get(d, "score")
```

Documented expected output:

```text
true
42
```

### dictionary.remove

`func dictionary.remove(d:ints,key:string) -> bool`

EN: Remove an existing key; return false if absent.

ES: Elimina una clave existente; devuelve false si no existe.

Native example flag: false.

```tripol
import tripol.dictionary

d = dictionary.create(2)
dictionary.set(d, "score", 42)
show dictionary.remove(d, "score")
show dictionary.size(d)
```

Documented expected output:

```text
true
0
```

---

## tripol.geometry — 6 functions

Import: `import tripol.geometry`

### geometry.rectangle

`func geometry.rectangle(w:float,h:float) -> float`

EN: Calculate a rectangle's area.

ES: Calcula el área de un rectángulo.

Native example flag: false.

```tripol
import tripol.geometry

show geometry.rectangle(4.0, 3.0)
```

Documented expected output:

```text
12.000000
```

### geometry.circle

`func geometry.circle(r:float) -> float`

EN: Calculate the area of a circle.

ES: Calcula el área de un círculo.

Native example flag: false.

```tripol
import tripol.geometry

show geometry.circle(1.0)
```

Documented expected output:

```text
3.141593
```

### geometry.circumference

`func geometry.circumference(r:float) -> float`

EN: Calculate the perimeter of a circle.

ES: Calcula el perímetro de un círculo.

Native example flag: false.

```tripol
import tripol.geometry

show geometry.circumference(1.0)
```

Documented expected output:

```text
6.283185
```

### geometry.triangle

`func geometry.triangle(b:float,h:float) -> float`

EN: Calculate triangle area from base and height.

ES: Calcula el área de un triángulo a partir de base y altura.

Native example flag: false.

```tripol
import tripol.geometry

show geometry.triangle(6.0, 4.0)
```

Documented expected output:

```text
12.000000
```

### geometry.distance

`func geometry.distance(x:float,y:float,u:float,v:float) -> float`

EN: Calculate Euclidean distance between two 2D points.

ES: Calcula la distancia euclidiana entre dos puntos 2D.

Native example flag: false.

```tripol
import tripol.geometry

show geometry.distance(0.0, 0.0, 3.0, 4.0)
```

Documented expected output:

```text
5.000000
```

### geometry.lerp

`func geometry.lerp(a:float,b:float,t:float) -> float`

EN: Interpolate or extrapolate between two values.

ES: Interpola o extrapola entre dos valores.

Native example flag: false.

Note: Unlike animation.linear, t is not clamped.

Nota: A diferencia de animation.linear, t no se limita.

```tripol
import tripol.geometry

show geometry.lerp(10.0, 20.0, 0.5)
```

Documented expected output:

```text
15.000000
```

---

## tripol.heap — 6 functions

Import: `import tripol.heap`

### heap.create

`func heap.create(limit=32) -> ints`

EN: Allocate a min-heap with space for the specified number of integers.

ES: Crea un montículo mínimo con espacio para la cantidad indicada de enteros.

Native example flag: false.

```tripol
import tripol.heap

h=heap.create(3)
show heap.size(h)
```

Documented expected output:

```text
0
```

### heap.capacity

`func heap.capacity(h:ints) -> int`

EN: Read the maximum number of values accepted by the heap.

ES: Lee la cantidad máxima de valores admitidos por el montículo.

Native example flag: false.

```tripol
import tripol.heap

h=heap.create(3)
show heap.capacity(h)
```

Documented expected output:

```text
3
```

### heap.size

`func heap.size(h:ints) -> int`

EN: Count the currently stored values in the heap.

ES: Cuenta los valores almacenados actualmente en el montículo.

Native example flag: false.

```tripol
import tripol.heap

h=heap.create(3)
heap.push(h,7)
show heap.size(h)
```

Documented expected output:

```text
1
```

### heap.push

`func heap.push(h:ints,value:int) -> bool`

EN: Insert an integer and return false if the heap is full.

ES: Inserta un entero y devuelve false si el montículo está lleno.

Native example flag: false.

```tripol
import tripol.heap

h=heap.create(1)
show heap.push(h,7),heap.push(h,4)
```

Documented expected output:

```text
true false
```

### heap.peek

`func heap.peek(h:ints) -> int`

EN: Read the smallest value of a nonempty heap.

ES: Lee el menor valor de un montículo que contiene elementos.

Native example flag: false.

```tripol
import tripol.heap

h=heap.create(3)
heap.push(h,7)
heap.push(h,4)
show heap.peek(h)
```

Documented expected output:

```text
4
```

### heap.pop

`func heap.pop(h:ints) -> int`

EN: Remove and return the smallest value of a nonempty heap.

ES: Retira y devuelve el menor valor de un montículo que contiene elementos.

Native example flag: false.

```tripol
import tripol.heap

h=heap.create(2)
heap.push(h,7)
heap.push(h,4)
show heap.pop(h),heap.pop(h)
```

Documented expected output:

```text
4 7
```

---

## tripol.image — 6 functions

Import: `import tripol.image`

### image.flip_x

`func image.flip_x(pixels:ints,w:int,h:int) -> ints`

EN: Return a horizontally mirrored image.

ES: Devuelve una imagen reflejada horizontalmente.

Native example flag: false.

```tripol
import tripol.image

p = [1, 2, 3, 4]
show image.flip_x(p, 2, 2)
```

Documented expected output:

```text
[2, 1, 4, 3]
```

### image.flip_y

`func image.flip_y(pixels:ints,w:int,h:int) -> ints`

EN: Return a vertically mirrored image.

ES: Devuelve una imagen reflejada verticalmente.

Native example flag: false.

```tripol
import tripol.image

p = [1, 2, 3, 4]
show image.flip_y(p, 2, 2)
```

Documented expected output:

```text
[3, 4, 1, 2]
```

### image.crop

`func image.crop(pixels:ints,w:int,h:int,x:int,y:int,cw:int,ch:int) -> ints`

EN: Copy a rectangular area into an independent image.

ES: Copia un área rectangular a una imagen independiente.

Native example flag: false.

```tripol
import tripol.image

p = [1, 2, 3, 4]
show image.crop(p, 2, 2, 1, 0, 1, 2)
```

Documented expected output:

```text
[2, 4]
```

### image.grayscale

`func image.grayscale(pixels:ints) -> ints`

EN: Convert RGB pixels to grayscale using weighted channels.

ES: Convierte píxeles RGB a escala de grises ponderando los canales.

Native example flag: false.

```tripol
import tripol.image

p = [0, 16777215]
show image.grayscale(p)
```

Documented expected output:

```text
[0, 16777215]
```

### image.ppm

`func image.ppm(pixels:ints,w:int,h:int) -> string`

EN: Encode a small RGB image in ASCII PPM (P3) format.

ES: Codifica una imagen RGB pequeña en formato PPM ASCII (P3).

Native example flag: false.

Note: At most 256 pixels; not PNG or JPEG.

Nota: Máximo 256 píxeles; no es PNG ni JPEG.

```tripol
import tripol.image

p = [16711680]
show image.ppm(p, 1, 1)
```

### image.save_ppm

`func image.save_ppm(path:string,pixels:ints,w:int,h:int) -> int`

EN: Save a small image as an ASCII PPM file.

ES: Guarda una imagen pequeña como archivo PPM ASCII.

Native example flag: true.

Note: Requires file output support. At most 256 pixels.

Nota: Requiere soporte de archivos. Máximo 256 píxeles.

```tripol
import tripol.image

p = [16711680]
show image.save_ppm("demo.ppm", p, 1, 1) > 0
```

Documented expected output:

```text
true
```

---

## tripol.io — 2 functions

Import: `import tripol.io`

### io.read

`func io.read(path:string,limit=65536) -> ints`

EN: Read raw bytes from a path up to the specified limit.

ES: Lee los bytes de una ruta hasta el límite indicado.

Native example flag: false.

```tripol
import tripol.io

data=ints(2)
data[0]=255
io.write("data.bin",data)
show io.read("data.bin",2)
```

Documented expected output:

```text
[255, 0]
```

### io.write

`func io.write(path:string,data:ints) -> int`

EN: Write raw byte values to the specified path and return the byte count.

ES: Escribe los bytes en la ruta indicada y devuelve la cantidad escrita.

Native example flag: false.

```tripol
import tripol.io

data=ints(2)
data[0]=255
show io.write("data.bin",data)
```

Documented expected output:

```text
2
```

---

## tripol.json — 18 functions + 3 methods

Import: `import tripol.json`

### json.quote

`func json.quote(s:string) -> string`

EN: Escape text and surround it with JSON quotes.

ES: Escapa un texto y lo rodea de comillas JSON.

Native example flag: false.

Note: The native encoder allocates one output string. The encoded result respects the TRIPOL string capacity.

Nota: El codificador nativo reserva una cadena de salida. El resultado respeta la capacidad de cadena de TRIPOL.

```tripol
import tripol.json

show json.quote("Hola \"TRIPOL\"")
```

Documented expected output:

```text
"Hola \"TRIPOL\""
```

### json.integer

`func json.integer(n:int) -> string`

EN: Convert an integer to JSON numeric text.

ES: Convierte un entero en texto numérico JSON.

Native example flag: false.

```tripol
import tripol.json

show json.integer(-42)
```

Documented expected output:

```text
-42
```

### json.boolean

`func json.boolean(b:bool) -> string`

EN: Convert a boolean into JSON text.

ES: Convierte un booleano en texto JSON.

Native example flag: false.

```tripol
import tripol.json

show json.boolean(true)
show json.boolean(false)
```

Documented expected output:

```text
true
false
```

### json.pair

`func json.pair(key:string,value:string) -> string`

EN: Create a JSON object containing one string property.

ES: Crea un objeto JSON con una propiedad de texto.

Native example flag: false.

Note: Both key and value are strings, not raw JSON values.

Nota: La clave y el valor son textos, no valores JSON sin procesar.

```tripol
import tripol.json

show json.pair("name", "TRIPOL")
```

Documented expected output:

```text
{"name":"TRIPOL"}
```

### json.integers

`func json.integers(a:ints) -> string`

EN: Convert an integer array into JSON array text.

ES: Convierte un arreglo de enteros en texto de arreglo JSON.

Native example flag: false.

Note: Maximum 128 integers.

Nota: Máximo 128 enteros.

```tripol
import tripol.json

numbers = [1, 2, 3]
show json.integers(numbers)
```

Documented expected output:

```text
[1,2,3]
```

### json.ws

`func json.ws(s:string,i:int,n:int) -> int`

EN: Skip JSON whitespace from a byte offset.

ES: Salta los espacios en blanco JSON desde un índice de bytes.

Native example flag: false.

Note: Advanced scanner helper; offsets start at zero and n is the exclusive end.

Nota: Función avanzada del analizador; los índices empiezan en cero y n es el límite exclusivo.

```tripol
import tripol.json

s = "  42"
show json.ws(s, 0, len(s))
```

Documented expected output:

```text
2
```

### json.hex

`func json.hex(b:int) -> bool`

EN: Check whether a byte is an ASCII hexadecimal digit.

ES: Comprueba si un byte es un dígito hexadecimal ASCII.

Native example flag: false.

Note: 65 is A, 71 is G.

Nota: 65 representa A y 71 representa G.

```tripol
import tripol.json

show json.hex(65)
show json.hex(71)
```

Documented expected output:

```text
true
false
```

### json.digit

`func json.digit(b:int) -> bool`

EN: Check whether a byte is an ASCII decimal digit.

ES: Comprueba si un byte es un dígito decimal ASCII.

Native example flag: false.

Note: 57 is 9, 65 is A.

Nota: 57 representa 9 y 65 representa A.

```tripol
import tripol.json

show json.digit(57)
show json.digit(65)
```

Documented expected output:

```text
true
false
```

### json.quoted

`func json.quoted(s:string,i:int,n:int) -> int`

EN: Scan a JSON quoted string and return the byte offset immediately after it.

ES: Analiza una cadena JSON entre comillas y devuelve el índice de bytes justo después de ella.

Native example flag: false.

Note: Returns -1 if the string is invalid. Use quote() for creating strings.

Nota: Devuelve -1 si la cadena no es válida. Usa quote() para crear cadenas.

```tripol
import tripol.json

s = "\"hola\""
show json.quoted(s, 0, len(s))
```

Documented expected output:

```text
6
```

### json.literal

`func json.literal(s:string,i:int,n:int,w:string) -> int`

EN: Match a literal token at the current byte offset.

ES: Comprueba si un literal coincide desde el índice de bytes indicado.

Native example flag: false.

Note: Returns the offset after the match or -1. It does not validate a complete JSON document.

Nota: Devuelve el índice posterior a la coincidencia o -1. No valida un documento JSON completo.

```tripol
import tripol.json

s = "true"
show json.literal(s, 0, len(s), "true")
show json.literal(s, 0, len(s), "null")
```

Documented expected output:

```text
4
-1
```

### json.number

`func json.number(s:string,i:int,n:int) -> int`

EN: Scan a JSON number and return where it ends, without converting it to a number.

ES: Analiza un número JSON y devuelve dónde termina, sin convertirlo a un número.

Native example flag: false.

Note: Returns -1 for invalid number syntax; use valid() to validate the complete input.

Nota: Devuelve -1 para sintaxis numérica inválida; usa valid() para validar toda la entrada.

```tripol
import tripol.json

s = "-12.5"
show json.number(s, 0, len(s))
```

Documented expected output:

```text
5
```

### json.value

`func json.value(s:string,i:int,n:int,depth:int) -> int`

EN: Scan one JSON value (object, array, number, string, boolean or null).

ES: Analiza un valor JSON (objeto, arreglo, número, texto, booleano o null).

Native example flag: false.

Note: Advanced scanner; returns the next byte offset or -1. Nesting is limited to depth 32.

Nota: Analizador avanzado; devuelve el siguiente índice de bytes o -1. Profundidad máxima 32.

```tripol
import tripol.json

s = "[1,2]"
show json.value(s, 0, len(s), 0)
```

Documented expected output:

```text
5
```

### json.valid

`func json.valid(s:string) -> bool`

EN: Check whether the entire string contains valid JSON.

ES: Comprueba si toda la cadena contiene JSON válido.

Native example flag: false.

Note: Maximum input 4096 bytes; nesting depth limited to 32.

Nota: Máximo 4096 bytes de entrada y profundidad 32.

```tripol
import tripol.json

show json.valid("{\"score\":100}")
show json.valid("{score:100}")
```

Documented expected output:

```text
true
false
```

### json.get_raw

`func json.get_raw(s:string,key:string) -> string`

EN: Read the original JSON text of a top-level property.

ES: Lee el texto JSON original de una propiedad del objeto principal.

Native example flag: false.

Note: Maximum 4096 bytes. String values retain their quotes. Missing or invalid properties return empty text; this is not a general object decoder.

Nota: Máximo 4096 bytes. Los valores de texto conservan sus comillas. Las propiedades ausentes o inválidas devuelven texto vacío; no es un decodificador general.

```tripol
import tripol.json

profile = "{\"name\":\"TRIPOL\",\"score\":100}"
show json.get_raw(profile, "name")
show json.get_raw(profile, "score")
```

Documented expected output:

```text
"TRIPOL"
100
```

### json.get_raw_large

`func json.get_raw_large(s:string,key:string) -> string`

EN: Read a top-level property from a larger JSON document.

ES: Lee una propiedad del objeto principal de un documento JSON más grande.

Native example flag: false.

Note: Maximum 65535 input bytes and 512 key bytes. Returns raw JSON, not a decoded object.

Nota: Máximo 65535 bytes de entrada y 512 bytes de clave. Devuelve JSON sin procesar, no un objeto decodificado.

```tripol
import tripol.json

event = "{\"data\":{\"score\":100}}"
show json.get_raw_large(event, "data")
```

Documented expected output:

```text
{"score":100}
```

### json.string_equals

`func json.string_equals(raw:string, value:string) -> bool`

EN: Compare a quoted JSON string to UTF-8 text after interpreting JSON escapes.

ES: Compara una cadena JSON entre comillas con texto UTF-8 después de interpretar sus escapes.

Native example flag: false.

Note: The comparison accepts up to 512 decoded UTF-8 bytes and returns false for invalid escape sequences or unpaired surrogates.

Nota: La comparación acepta hasta 512 bytes UTF-8 decodificados y devuelve false para escapes inválidos o pares sustitutos incompletos.

```tripol
import tripol.json

raw = "\"caf\\u00e9\\/\\ud83d\\ude00\""
show json.string_equals(raw, "café/😀")
show json.string_equals(raw, "cafe/😀")
```

Documented expected output:

```text
true
false
```

### json.unquote

`func json.unquote(raw:string) -> string`

EN: Decode JSON escapes and Unicode surrogate pairs.

ES: Decodifica escapes JSON y pares de sustitutos Unicode.

Native example flag: false.

```tripol
import tripol.json

show json.unquote("\"caf\\u00e9 \\ud83d\\ude00\"")
```

Documented expected output:

```text
café 😀
```

---

### json.object

Create an empty typed JSON object. / Crea un objeto JSON tipado vacío.

```tripol
func json.object() -> Object
```

```tripol
import tripol.json

obj=json.object()
show obj.encode()
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
{}
```

### json.Object.set

Store a string, integer, boolean or nested object by key. / Guarda texto, entero, booleano u objeto anidado por clave.

```tripol
func json.Object.set(key:string,value:string) -> bool
func json.Object.set(key:string,value:int) -> bool
func json.Object.set(key:string,value:bool) -> bool
func json.Object.set(key:string,value:Object) -> bool
```

```tripol
import tripol.json

obj=json.object()
show obj.set("n",42)
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### json.Object.with

Set a field and return the object for chaining. / Asigna un campo y devuelve el objeto para encadenar.

```tripol
func json.Object.with(key:string,value:string) -> Object
func json.Object.with(key:string,value:int) -> Object
func json.Object.with(key:string,value:bool) -> Object
func json.Object.with(key:string,value:Object) -> Object
```

```tripol
import tripol.json

obj=json.object().with("status",true)
show obj.encode()
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
{"status":true}
```

### json.Object.encode

Serialize the current object as JSON text. / Serializa el objeto como texto JSON.

```tripol
func json.Object.encode() -> string
```

```tripol
import tripol.json

obj=json.object()
obj.set("n",42)
show obj.encode()
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
{"n":42}
```

## tripol.matrix — 6 functions

Import: `import tripol.matrix`

### matrix.size

`func matrix.size(rows:int,cols:int) -> int`

EN: Calculate the number of elements in a valid matrix shape.

ES: Calcula cuántos elementos tiene una matriz de dimensiones válidas.

Native example flag: false.

```tripol
import tripol.matrix

show matrix.size(2, 3)
```

Documented expected output:

```text
6
```

### matrix.zeros

`func matrix.zeros(rows:int,cols:int) -> ints`

EN: Create a flat zero-filled integer matrix.

ES: Crea una matriz plana de enteros llena de ceros.

Native example flag: false.

```tripol
import tripol.matrix

show matrix.zeros(2, 2)
```

Documented expected output:

```text
[0, 0, 0, 0]
```

### matrix.identity

`func matrix.identity(n:int) -> ints`

EN: Create a square identity matrix in row-major order.

ES: Crea una matriz identidad cuadrada ordenada por filas.

Native example flag: false.

```tripol
import tripol.matrix

show matrix.identity(2)
```

Documented expected output:

```text
[1, 0, 0, 1]
```

### matrix.transpose

`func matrix.transpose(a:ints,rows:int,cols:int) -> ints`

EN: Swap rows and columns of a flat matrix.

ES: Intercambia filas y columnas de una matriz plana.

Native example flag: false.

```tripol
import tripol.matrix

a = [1, 2, 3, 4, 5, 6]
show matrix.transpose(a, 2, 3)
```

Documented expected output:

```text
[1, 4, 2, 5, 3, 6]
```

### matrix.add

`func matrix.add(a:ints,b:ints) -> ints`

EN: Add matching elements of two equally sized matrices.

ES: Suma elementos correspondientes de dos matrices de igual tamaño.

Native example flag: false.

```tripol
import tripol.matrix

a = [1, 2]
b = [3, 4]
show matrix.add(a, b)
```

Documented expected output:

```text
[4, 6]
```

### matrix.dot

`func matrix.dot(a:ints,b:ints) -> int`

EN: Calculate the dot product of two integer vectors.

ES: Calcula el producto escalar de dos vectores enteros.

Native example flag: false.

```tripol
import tripol.matrix

a = [1, 2]
b = [3, 4]
show matrix.dot(a, b)
```

Documented expected output:

```text
11
```

---

## tripol.particles — 6 functions

Import: `import tripol.particles`

### particles.create

`func particles.create(n:int) -> ints`

EN: Allocate a pool of particles with zero life.

ES: Reserva un grupo de partículas con vida inicial cero.

Native example flag: false.

```tripol
import tripol.particles

p = particles.create(3)
show particles.count(p)
```

Documented expected output:

```text
3
```

### particles.count

`func particles.count(p:ints) -> int`

EN: Return the total number of particle slots.

ES: Devuelve el total de espacios para partículas.

Native example flag: false.

```tripol
import tripol.particles

p = particles.create(3)
show particles.count(p)
```

Documented expected output:

```text
3
```

### particles.step

`func particles.step(p:ints) -> int`

EN: Advance positions and reduce life by one tick; count survivors.

ES: Avanza las posiciones y reduce la vida un tick; cuenta las supervivientes.

Native example flag: false.

```tripol
import tripol.particles

p = particles.create(2)
particles.burst(p, 5, 5, 1)
show particles.step(p)
```

Documented expected output:

```text
0
```

### particles.alive

`func particles.alive(p:ints) -> int`

EN: Count particles whose life is positive.

ES: Cuenta las partículas cuya vida es positiva.

Native example flag: false.

```tripol
import tripol.particles

p = particles.create(3)
particles.burst(p, 5, 5)
show particles.alive(p)
```

Documented expected output:

```text
3
```

### particles.burst

`func particles.burst(p:ints,x:int,y:int,life=16) -> int`

EN: Reset the pool into a deterministic burst at a position.

ES: Reinicia el grupo como una explosión determinista en una posición.

Native example flag: false.

```tripol
import tripol.particles

p = particles.create(3)
show particles.burst(p, 10, 10, 4)
```

Documented expected output:

```text
3
```

### particles.draw

`func particles.draw(p:ints,dst:ints,w:int,h:int,ox:int,oy:int,color:int) -> int`

EN: Draw living particles into a clipped RGB framebuffer.

ES: Dibuja partículas vivas dentro de un arreglo RGB recortado.

Native example flag: false.

```tripol
import tripol.particles

import tripol.canvas

p = particles.create(1)
particles.burst(p, 1, 1)
frame = canvas.create(3, 3)
show particles.draw(p, frame, 3, 3, 0, 0, 255)
```

Documented expected output:

```text
1
```

---

## tripol.physics — 5 functions

Import: `import tripol.physics`

### physics.overlap

`func physics.overlap(x:int,y:int,w:int,h:int,u:int,v:int,a:int,b:int) -> bool`

EN: Check if two axis-aligned rectangles overlap with positive area.

ES: Comprueba si dos rectángulos alineados a los ejes se superponen con área positiva.

Native example flag: false.

```tripol
import tripol.physics

show physics.overlap(0, 0, 10, 10, 5, 5, 10, 10)
```

Documented expected output:

```text
true
```

### physics.contains

`func physics.contains(x:int,y:int,w:int,h:int,px:int,py:int) -> bool`

EN: Check whether a point is inside a rectangle.

ES: Comprueba si un punto está dentro de un rectángulo.

Native example flag: false.

```tripol
import tripol.physics

show physics.contains(0, 0, 10, 10, 5, 5)
```

Documented expected output:

```text
true
```

### physics.clamp_position

`func physics.clamp_position(pos:int,extent:int,limit:int) -> int`

EN: Keep an object of a given extent inside a movement limit.

ES: Mantiene un objeto de tamaño dado dentro de un límite de movimiento.

Native example flag: false.

```tripol
import tripol.physics

show physics.clamp_position(95, 10, 100)
```

Documented expected output:

```text
90
```

### physics.integrate

`func physics.integrate(pos:float,velocity:float,dt:float) -> float`

EN: Advance a position at constant velocity for a time step.

ES: Avanza una posición con velocidad constante durante un intervalo.

Native example flag: false.

```tripol
import tripol.physics

show physics.integrate(0.0, 10.0, 2.0)
```

Documented expected output:

```text
20.000000
```

### physics.gravity

`func physics.gravity(velocity:float,acceleration:float,dt:float) -> float`

EN: Update velocity using acceleration and a time step.

ES: Actualiza la velocidad usando aceleración y un intervalo de tiempo.

Native example flag: false.

```tripol
import tripol.physics

show physics.gravity(0.0, 9.8, 1.0)
```

Documented expected output:

```text
9.800000
```

---

## tripol.queue — 7 functions

Import: `import tripol.queue`

### queue.create

`func queue.create(limit=32) -> ints`

EN: Create an empty bounded FIFO queue.

ES: Crea una cola FIFO vacía de capacidad fija.

Native example flag: false.

```tripol
import tripol.queue

q = queue.create(3)
show queue.size(q)
```

Documented expected output:

```text
0
```

### queue.capacity

`func queue.capacity(q:ints) -> int`

EN: Return the maximum number of stored integers.

ES: Devuelve la cantidad máxima de enteros almacenados.

Native example flag: false.

```tripol
import tripol.queue

q = queue.create(3)
show queue.capacity(q)
```

Documented expected output:

```text
3
```

### queue.size

`func queue.size(q:ints) -> int`

EN: Return the number of stored integers.

ES: Devuelve la cantidad de enteros almacenados.

Native example flag: false.

```tripol
import tripol.queue

q = queue.create(3)
queue.push(q, 10)
show queue.size(q)
```

Documented expected output:

```text
1
```

### queue.push

`func queue.push(q:ints,value:int) -> bool`

EN: Add an integer; return false when the container is full.

ES: Añade un entero; devuelve false cuando el contenedor está lleno.

Native example flag: false.

```tripol
import tripol.queue

q = queue.create(3)
show queue.push(q, 10)
```

Documented expected output:

```text
true
```

### queue.peek

`func queue.peek(q:ints) -> int`

EN: Read the oldest value without removing it.

ES: Lee el valor más antiguo sin eliminarlo.

Native example flag: false.

Note: Must not be empty.

Nota: No debe estar vacío.

```tripol
import tripol.queue

q = queue.create(3)
queue.push(q, 10)
queue.push(q, 20)
show queue.peek(q)
```

Documented expected output:

```text
10
```

### queue.pop

`func queue.pop(q:ints) -> int`

EN: Remove and return the oldest value.

ES: Elimina y devuelve el valor más antiguo.

Native example flag: false.

Note: Must not be empty.

Nota: No debe estar vacío.

```tripol
import tripol.queue

q = queue.create(3)
queue.push(q, 10)
queue.push(q, 20)
show queue.pop(q)
```

Documented expected output:

```text
10
```

### queue.clear

`func queue.clear(q:ints) -> int`

EN: Remove all logical entries and return the previous count.

ES: Elimina todas las entradas lógicas y devuelve la cantidad anterior.

Native example flag: false.

```tripol
import tripol.queue

q = queue.create(3)
queue.push(q, 10)
queue.push(q, 20)
show queue.clear(q)
show queue.size(q)
```

Documented expected output:

```text
2
0
```

---

## tripol.regex — 2 functions

Import: `import tripol.regex`

### regex.valid

`func regex.valid(pattern:string) -> bool`

EN: Check if a pattern uses the supported literal, dot and star syntax.

ES: Comprueba si el patrón usa la sintaxis admitida de literales, punto y asterisco.

Native example flag: false.

```tripol
import tripol.regex

show regex.valid("ab.*")
```

Documented expected output:

```text
true
```

### regex.matches

`func regex.matches(pattern:string,value:string) -> bool`

EN: Match the complete string against a bounded byte-oriented pattern.

ES: Compara toda la cadena con un patrón acotado que opera sobre bytes.

Native example flag: false.

```tripol
import tripol.regex

show regex.matches("ab.*","abacus")
```

Documented expected output:

```text
true
```

---

## tripol.sequences — 4 functions

Import: `import tripol.sequences`

### sequences.arithmetic

`func sequences.arithmetic(n:int,start:int,step:int) -> ints`

EN: Create a sequence by adding a constant step.

ES: Crea una secuencia sumando un paso constante.

Native example flag: false.

```tripol
import tripol.sequences

show sequences.arithmetic(4, 2, 3)
```

Documented expected output:

```text
[2, 5, 8, 11]
```

### sequences.fibonacci

`func sequences.fibonacci(n:int) -> ints`

EN: Create the first n Fibonacci numbers, beginning at zero.

ES: Crea los primeros n números de Fibonacci empezando por cero.

Native example flag: false.

```tripol
import tripol.sequences

show sequences.fibonacci(7)
```

Documented expected output:

```text
[0, 1, 1, 2, 3, 5, 8]
```

### sequences.factorial

`func sequences.factorial(n:int) -> int`

EN: Calculate n factorial for integers from 0 to 20.

ES: Calcula el factorial de un entero de 0 a 20.

Native example flag: false.

```tripol
import tripol.sequences

show sequences.factorial(5)
```

Documented expected output:

```text
120
```

### sequences.triangular

`func sequences.triangular(n:int) -> int`

EN: Calculate the sum of integers from 1 to n.

ES: Calcula la suma de enteros desde 1 hasta n.

Native example flag: false.

```tripol
import tripol.sequences

show sequences.triangular(5)
```

Documented expected output:

```text
15
```

---

## tripol.sprite — 6 functions

Import: `import tripol.sprite`

### sprite.draw_key

`func sprite.draw_key(dst:ints,w:int,h:int,src:ints,sw:int,sh:int,x:int,y:int,key:int) -> int`

EN: Draw a sprite while skipping pixels matching a transparent color key.

ES: Dibuja un sprite omitiendo los píxeles que coinciden con el color transparente.

Native example flag: false.

```tripol
import tripol.sprite

dst = ints(2)
src = [0, 255]
show sprite.draw_key(dst, 2, 1, src, 2, 1, 0, 0, 0)
show dst
```

Documented expected output:

```text
1
[0, 255]
```

### sprite.draw

`func sprite.draw(dst:ints,w:int,h:int,src:ints,sw:int,sh:int,x:int,y:int) -> int`

EN: Copy visible sprite pixels into a destination framebuffer.

ES: Copia píxeles visibles de un sprite al arreglo de destino.

Native example flag: false.

```tripol
import tripol.sprite

dst = ints(2)
src = [255]
show sprite.draw(dst, 2, 1, src, 1, 1, 1, 0)
show dst
```

Documented expected output:

```text
1
[0, 255]
```

### sprite.solid

`func sprite.solid(w:int,h:int,color:int) -> ints`

EN: Create a rectangular sprite of one color.

ES: Crea un sprite rectangular de un único color.

Native example flag: false.

```tripol
import tripol.sprite

show sprite.solid(2, 1, 255)
```

Documented expected output:

```text
[255, 255]
```

### sprite.rotate90

`func sprite.rotate90(src:ints,w:int,h:int) -> ints`

EN: Rotate a sprite clockwise by 90 degrees.

ES: Gira un sprite 90 grados en sentido horario.

Native example flag: false.

```tripol
import tripol.sprite

src = [1, 2, 3, 4]
show sprite.rotate90(src, 2, 2)
```

Documented expected output:

```text
[3, 1, 4, 2]
```

### sprite.scale_nearest

`func sprite.scale_nearest(src:ints,w:int,h:int,nw:int,nh:int) -> ints`

EN: Resize using nearest-neighbor sampling.

ES: Cambia el tamaño usando vecino más cercano.

Native example flag: false.

```tripol
import tripol.sprite

src = [1, 2]
show sprite.scale_nearest(src, 2, 1, 4, 1)
```

Documented expected output:

```text
[1, 1, 2, 2]
```

### sprite.draw_scaled

`func sprite.draw_scaled(dst:ints,w:int,h:int,src:ints,sw:int,sh:int,x:int,y:int,dw:int,dh:int,key=-1,opacity=100) -> int`

EN: Draw a scaled sprite with clipping, optional color key and opacity.

ES: Dibuja un sprite escalado con recorte, color transparente opcional y opacidad.

Native example flag: false.

```tripol
import tripol.sprite

dst = ints(2)
src = [255]
show sprite.draw_scaled(dst, 2, 1, src, 1, 1, 0, 0, 2, 1)
show dst
```

Documented expected output:

```text
2
[255, 255]
```

---

## tripol.stack — 7 functions

Import: `import tripol.stack`

### stack.create

`func stack.create(limit=32) -> ints`

EN: Create an empty bounded LIFO stack.

ES: Crea una pila LIFO vacía de capacidad fija.

Native example flag: false.

```tripol
import tripol.stack

s = stack.create(3)
show stack.size(s)
```

Documented expected output:

```text
0
```

### stack.capacity

`func stack.capacity(s:ints) -> int`

EN: Return the maximum number of stored integers.

ES: Devuelve la cantidad máxima de enteros almacenados.

Native example flag: false.

```tripol
import tripol.stack

s = stack.create(3)
show stack.capacity(s)
```

Documented expected output:

```text
3
```

### stack.size

`func stack.size(s:ints) -> int`

EN: Return the number of stored integers.

ES: Devuelve la cantidad de enteros almacenados.

Native example flag: false.

```tripol
import tripol.stack

s = stack.create(3)
stack.push(s, 10)
show stack.size(s)
```

Documented expected output:

```text
1
```

### stack.push

`func stack.push(s:ints,value:int) -> bool`

EN: Add an integer; return false when the container is full.

ES: Añade un entero; devuelve false cuando el contenedor está lleno.

Native example flag: false.

```tripol
import tripol.stack

s = stack.create(3)
show stack.push(s, 10)
```

Documented expected output:

```text
true
```

### stack.peek

`func stack.peek(s:ints) -> int`

EN: Read the last inserted value without removing it.

ES: Lee el último valor insertado sin eliminarlo.

Native example flag: false.

Note: Must not be empty.

Nota: No debe estar vacío.

```tripol
import tripol.stack

s = stack.create(3)
stack.push(s, 10)
stack.push(s, 20)
show stack.peek(s)
```

Documented expected output:

```text
20
```

### stack.pop

`func stack.pop(s:ints) -> int`

EN: Remove and return the last inserted value.

ES: Elimina y devuelve el último valor insertado.

Native example flag: false.

Note: Must not be empty.

Nota: No debe estar vacío.

```tripol
import tripol.stack

s = stack.create(3)
stack.push(s, 10)
stack.push(s, 20)
show stack.pop(s)
```

Documented expected output:

```text
20
```

### stack.clear

`func stack.clear(s:ints) -> int`

EN: Remove all logical entries and return the previous count.

ES: Elimina todas las entradas lógicas y devuelve la cantidad anterior.

Native example flag: false.

```tripol
import tripol.stack

s = stack.create(3)
stack.push(s, 10)
stack.push(s, 20)
show stack.clear(s)
show stack.size(s)
```

Documented expected output:

```text
2
0
```

---

## tripol.stats — 6 functions

Import: `import tripol.stats`

### stats.mean

`func stats.mean(a:ints) -> float`

EN: Calculate the arithmetic average of a nonempty array.

ES: Calcula el promedio aritmético de un arreglo no vacío.

Native example flag: false.

```tripol
import tripol.stats

a = [2, 4, 6]
show stats.mean(a)
```

Documented expected output:

```text
4.000000
```

### stats.minimum

`func stats.minimum(a:ints) -> int`

EN: Find the smallest element of a nonempty array.

ES: Busca el menor elemento de un arreglo no vacío.

Native example flag: false.

```tripol
import tripol.stats

a = [2, 4, 6]
show stats.minimum(a)
```

Documented expected output:

```text
2
```

### stats.maximum

`func stats.maximum(a:ints) -> int`

EN: Find the largest element of a nonempty array.

ES: Busca el mayor elemento de un arreglo no vacío.

Native example flag: false.

```tripol
import tripol.stats

a = [2, 4, 6]
show stats.maximum(a)
```

Documented expected output:

```text
6
```

### stats.variance

`func stats.variance(a:ints) -> float`

EN: Calculate population variance (divide by n).

ES: Calcula la varianza poblacional (divide entre n).

Native example flag: false.

```tripol
import tripol.stats

a = [2, 4]
show stats.variance(a)
```

Documented expected output:

```text
1.000000
```

### stats.deviation

`func stats.deviation(a:ints) -> float`

EN: Calculate population standard deviation.

ES: Calcula la desviación estándar poblacional.

Native example flag: false.

```tripol
import tripol.stats

a = [2, 4]
show stats.deviation(a)
```

Documented expected output:

```text
1.000000
```

### stats.total

`func stats.total(a:ints) -> int`

EN: Sum all array elements.

ES: Suma todos los elementos del arreglo.

Native example flag: false.

```tripol
import tripol.stats

a = [2, 4, 6]
show stats.total(a)
```

Documented expected output:

```text
12
```

---

## tripol.stream — 6 functions

Import: `import tripol.stream`

### stream.create

`func stream.create() -> ints`

EN: Allocate a byte cursor positioned at the start of an array.

ES: Crea un cursor de bytes situado al comienzo de un arreglo.

Native example flag: false.

```tripol
import tripol.stream

cursor=stream.create()
show stream.position(cursor)
```

Documented expected output:

```text
0
```

### stream.position

`func stream.position(cursor:ints) -> int`

EN: Read the cursor's zero-based byte offset.

ES: Lee la posición en bytes del cursor, a partir de cero.

Native example flag: false.

```tripol
import tripol.stream

cursor=stream.create()
show stream.position(cursor)
```

Documented expected output:

```text
0
```

### stream.remaining

`func stream.remaining(cursor:ints,data:ints) -> int`

EN: Count bytes available from the current offset to the end of the data.

ES: Cuenta los bytes disponibles desde la posición actual hasta el final de los datos.

Native example flag: false.

```tripol
import tripol.stream

cursor=stream.create()
data=ints(3)
show stream.remaining(cursor,data)
```

Documented expected output:

```text
3
```

### stream.seek

`func stream.seek(cursor:ints,data:ints,offset:int) -> int`

EN: Set the cursor to an offset within the data, including the end.

ES: Coloca el cursor en una posición de los datos, incluido el final.

Native example flag: false.

```tripol
import tripol.stream

cursor=stream.create()
data=ints(3)
show stream.seek(cursor,data,2)
```

Documented expected output:

```text
2
```

### stream.read_byte

`func stream.read_byte(cursor:ints,data:ints) -> int`

EN: Read a byte and advance the cursor by one.

ES: Lee un byte y avanza el cursor una posición.

Native example flag: false.

```tripol
import tripol.stream

cursor=stream.create()
data=ints(1)
data[0]=255
show stream.read_byte(cursor,data)
```

Documented expected output:

```text
255
```

### stream.read_u16

`func stream.read_u16(cursor:ints,data:ints) -> int`

EN: Read two bytes as a little-endian unsigned integer and advance the cursor.

ES: Lee dos bytes como un entero sin signo en orden little-endian y avanza el cursor.

Native example flag: false.

```tripol
import tripol.stream

cursor=stream.create()
data=ints(2)
data[0]=52
data[1]=18
show stream.read_u16(cursor,data)
```

Documented expected output:

```text
4660
```

---

## tripol.testing — 4 functions

Import: `import tripol.testing`

### testing.equal_int

`func testing.equal_int(actual:int,expected:int) -> bool`

EN: Assert that two integers are equal; return true when they match.

ES: Comprueba que dos enteros sean iguales; devuelve true si coinciden.

Native example flag: false.

Note: A mismatch triggers an assertion rather than returning false.

Nota: Una diferencia provoca una aserción en vez de devolver false.

```tripol
import tripol.testing

show testing.equal_int(2 + 3, 5)
```

Documented expected output:

```text
true
```

### testing.equal_text

`func testing.equal_text(actual:string,expected:string) -> bool`

EN: Assert that two strings are equal.

ES: Comprueba que dos textos sean iguales.

Native example flag: false.

```tripol
import tripol.testing

show testing.equal_text("TRI" + "POL", "TRIPOL")
```

Documented expected output:

```text
true
```

### testing.equal_bool

`func testing.equal_bool(actual:bool,expected:bool) -> bool`

EN: Assert that two boolean values are equal.

ES: Comprueba que dos valores booleanos sean iguales.

Native example flag: false.

```tripol
import tripol.testing

show testing.equal_bool(5 > 2, true)
```

Documented expected output:

```text
true
```

### testing.near

`func testing.near(actual:float,expected:float,tolerance:float) -> bool`

EN: Assert that a decimal is within a nonnegative tolerance of the expected value.

ES: Comprueba que un decimal esté dentro de una tolerancia no negativa del valor esperado.

Native example flag: false.

```tripol
import tripol.testing

show testing.near(3.14, 3.141, 0.01)
```

Documented expected output:

```text
true
```

---

## tripol.time — 6 functions

Import: `import tripol.time`

### time.now_ms

`func time.now_ms() -> int`

EN: Read a monotonic timestamp in milliseconds.

ES: Lee una marca de tiempo monotónica en milisegundos.

Native example flag: false.

Note: Not a calendar date or Unix timestamp.

Nota: No es una fecha ni una marca de tiempo Unix.

```tripol
import tripol.time

show time.now_ms()
```

### time.elapsed

`func time.elapsed(start:int) -> int`

EN: Calculate milliseconds elapsed since a previous timestamp.

ES: Calcula los milisegundos transcurridos desde una marca anterior.

Native example flag: false.

```tripol
import tripol.time

start = time.now_ms()
show time.elapsed(start) >= 0
```

Documented expected output:

```text
true
```

### time.sleep

`func time.sleep(ms:int) -> int`

EN: Pause execution for a number of milliseconds.

ES: Pausa la ejecución durante varios milisegundos.

Native example flag: true.

Note: Browser sleep requires an isolated worker with the relevant capability.

Nota: La pausa web requiere un Worker aislado con la capacidad correspondiente.

```tripol
import tripol.time

time.sleep(100)
show "Ready"
```

Documented expected output:

```text
Ready
```

### time.seconds

`func time.seconds(ms:int) -> float`

EN: Convert milliseconds into decimal seconds.

ES: Convierte milisegundos en segundos decimales.

Native example flag: false.

```tripol
import tripol.time

show time.seconds(2500)
```

Documented expected output:

```text
2.500000
```

### time.now_ns

`func time.now_ns() -> int`

EN: Read a high-resolution monotonic timestamp in nanoseconds.

ES: Lee una marca monotónica de alta resolución en nanosegundos.

Native example flag: false.

Note: Origin and precision depend on the target.

Nota: El origen y la precisión dependen del destino.

```tripol
import tripol.time

show time.now_ns()
```

### time.elapsed_ns

`func time.elapsed_ns(start:int) -> int`

EN: Calculate nanoseconds elapsed since a previous timestamp.

ES: Calcula los nanosegundos transcurridos desde una marca anterior.

Native example flag: false.

```tripol
import tripol.time

start = time.now_ns()
show time.elapsed_ns(start) >= 0
```

Documented expected output:

```text
true
```

---

## tripol.utf8 — 5 functions

Import: `import tripol.utf8`

### utf8.continuation

`func utf8.continuation(b:int) -> bool`

EN: Check whether a byte is a UTF-8 continuation byte.

ES: Comprueba si el valor corresponde a un byte de continuación UTF-8.

Native example flag: false.

```tripol
import tripol.utf8

show utf8.continuation(169)
```

Documented expected output:

```text
true
```

### utf8.width_at

`func utf8.width_at(s:string,i:int) -> int`

EN: Return the byte width of the UTF-8 scalar at a byte offset, or zero for an invalid sequence.

ES: Devuelve el ancho en bytes del carácter Unicode en la posición indicada o cero si la secuencia es inválida.

Native example flag: false.

```tripol
import tripol.utf8

show utf8.width_at("é",0)
```

Documented expected output:

```text
2
```

### utf8.valid

`func utf8.valid(s:string) -> bool`

EN: Validate all UTF-8 sequences in the text.

ES: Valida todas las secuencias UTF-8 del texto.

Native example flag: false.

```tripol
import tripol.utf8

show utf8.valid("café")
```

Documented expected output:

```text
true
```

### utf8.count

`func utf8.count(s:string) -> int`

EN: Count the Unicode scalar values in valid UTF-8 text.

ES: Cuenta los valores escalares Unicode de un texto UTF-8 válido.

Native example flag: false.

```tripol
import tripol.utf8

show utf8.count("café")
```

Documented expected output:

```text
4
```

### utf8.scalar_at

`func utf8.scalar_at(s:string,i:int) -> int`

EN: Decode the Unicode scalar beginning at the given byte offset.

ES: Decodifica el carácter Unicode que comienza en la posición en bytes indicada.

Native example flag: false.

```tripol
import tripol.utf8

show utf8.scalar_at("é",0)
```

Documented expected output:

```text
233
```

---

## tripol.vector — 7 functions

Import: `import tripol.vector`

### vector.dot2

`func vector.dot2(x:float,y:float,u:float,v:float) -> float`

EN: Calculate the 2D dot product.

ES: Calcula el producto escalar de dos vectores 2D.

Native example flag: false.

```tripol
import tripol.vector

show vector.dot2(1.0, 2.0, 3.0, 4.0)
```

Documented expected output:

```text
11.000000
```

### vector.cross2

`func vector.cross2(x:float,y:float,u:float,v:float) -> float`

EN: Calculate the signed scalar 2D cross product.

ES: Calcula el producto cruzado escalar con signo de dos vectores 2D.

Native example flag: false.

```tripol
import tripol.vector

show vector.cross2(1.0, 0.0, 0.0, 1.0)
```

Documented expected output:

```text
1.000000
```

### vector.length2

`func vector.length2(x:float,y:float) -> float`

EN: Calculate the length of a 2D vector.

ES: Calcula la longitud de un vector 2D.

Native example flag: false.

```tripol
import tripol.vector

show vector.length2(3.0, 4.0)
```

Documented expected output:

```text
5.000000
```

### vector.distance2

`func vector.distance2(x:float,y:float,u:float,v:float) -> float`

EN: Calculate the distance between two 2D points.

ES: Calcula la distancia entre dos puntos 2D.

Native example flag: false.

```tripol
import tripol.vector

show vector.distance2(0.0, 0.0, 3.0, 4.0)
```

Documented expected output:

```text
5.000000
```

### vector.normalized_x

`func vector.normalized_x(x:float,y:float) -> float`

EN: Return the x component after normalizing a nonzero vector.

ES: Devuelve la componente x después de normalizar un vector no nulo.

Native example flag: false.

```tripol
import tripol.vector

show vector.normalized_x(3.0, 4.0)
```

Documented expected output:

```text
0.600000
```

### vector.normalized_y

`func vector.normalized_y(x:float,y:float) -> float`

EN: Return the y component after normalizing a nonzero vector.

ES: Devuelve la componente y después de normalizar un vector no nulo.

Native example flag: false.

```tripol
import tripol.vector

show vector.normalized_y(3.0, 4.0)
```

Documented expected output:

```text
0.800000
```

### vector.scale

`func vector.scale(a:ints,k:int) -> ints`

EN: Return a new integer array multiplied by a scalar.

ES: Devuelve un arreglo nuevo de enteros multiplicados por un escalar.

Native example flag: false.

```tripol
import tripol.vector

a = [1, 2, 3]
show vector.scale(a, 2)
```

Documented expected output:

```text
[2, 4, 6]
```

---

## tripol.atomic — 4 functions

Import: `import tripol.atomic`

### atomic.load

`func atomic.load(data:ints,index:int) -> int`

EN: Read the signed integer at an aligned array index.

ES: Lee el entero con signo en un índice alineado del arreglo.

Native example flag: false.

```tripol
import tripol.atomic

a=ints(1)
a[0]=7
show atomic.load(a,0)
```

Documented expected output:

```text
7
```

### atomic.exchange

`func atomic.exchange(data:ints,index:int,value:int) -> int`

EN: Replace the indexed value and return its previous value.

ES: Reemplaza el valor del índice y devuelve el anterior.

Native example flag: false.

```tripol
import tripol.atomic

a=ints(1)
a[0]=7
show atomic.exchange(a,0,9),a[0]
```

Documented expected output:

```text
7 9
```

### atomic.compare_exchange

`func atomic.compare_exchange(data:ints,index:int,expected:int,value:int) -> bool`

EN: Replace the value only if it equals the expected value.

ES: Reemplaza el valor solo si coincide con el esperado.

Native example flag: false.

```tripol
import tripol.atomic

a=ints(1)
show atomic.compare_exchange(a,0,0,9),a[0]
```

Documented expected output:

```text
true 9
```

### atomic.add

`func atomic.add(data:ints,index:int,delta:int) -> int`

EN: Add a signed delta with overflow checking and return the new value.

ES: Suma un incremento con comprobación de desbordamiento y devuelve el valor nuevo.

Native example flag: false.

```tripol
import tripol.atomic

a=ints(1)
show atomic.add(a,0,4)
```

Documented expected output:

```text
4
```

---

## tripol.audio — 6 functions

Import: `import tripol.audio`

### audio.wav_mono

`func audio.wav_mono(samples:ints,rate=22050) -> ints`

EN: Encode signed 16-bit mono samples as a PCM WAV byte array at the specified sample rate.

ES: Codifica muestras mono de 16 bits con signo en un arreglo PCM WAV a la frecuencia indicada.

Native example flag: false.

```tripol
import tripol.audio

a=ints(1)
show len(audio.wav_mono(a))
```

Documented expected output:

```text
46
```

### audio.square_wave

`func audio.square_wave(hz:int,milliseconds:int,rate=22050,amplitude=5000) -> ints`

EN: Generate samples for a square wave using frequency, duration, sample rate and amplitude.

ES: Genera muestras de onda cuadrada con frecuencia, duración, frecuencia de muestreo y amplitud.

Native example flag: false.

```tripol
import tripol.audio

show len(audio.square_wave(440,10,8000))
```

Documented expected output:

```text
80
```

### audio.save_wav

`func audio.save_wav(path:string,samples:ints,rate=22050) -> int`

EN: Write PCM samples to the given WAV path and return the byte count.

ES: Escribe muestras PCM en la ruta WAV y devuelve el número de bytes.

Native example flag: false.

```tripol
import tripol.audio

a=ints(1)
show audio.save_wav("tone.wav",a)
```

Documented expected output:

```text
46
```

### audio.playback_supported

`func audio.playback_supported() -> bool`

EN: Check whether the current target provides native WAV playback.

ES: Comprueba si el destino actual ofrece reproducción WAV nativa.

Native example flag: false.

```tripol
import tripol.audio

show audio.playback_supported()
```

Documented expected output:

```text
false
```

### audio.play_file

`func audio.play_file(path:string) -> bool`

EN: Play a WAV file from its path on a supported Windows host.

ES: Reproduce un archivo WAV por su ruta en Windows compatible.

Native example flag: true.

Target: "windows".

Note: Requires native Windows.

Nota: Requiere Windows nativo.

```tripol
import tripol.audio

a=ints(1)
audio.save_wav("tone.wav",a)
show audio.play_file("tone.wav")
```

### audio.beep

`func audio.beep(hz:int,milliseconds:int) -> bool`

EN: Play a brief system tone with the specified frequency and duration.

ES: Reproduce un tono breve con la frecuencia y duración indicadas.

Native example flag: true.

Target: "windows".

Note: Requires native Windows.

Nota: Requiere Windows nativo.

```tripol
import tripol.audio

show audio.beep(440,100)
```

---

## tripol.binary — 8 functions

Import: `import tripol.binary`

### binary.is_byte

`func binary.is_byte(value:int) -> bool`

EN: Check if an integer fits in one byte.

ES: Comprueba si un entero cabe en un byte.

Native example flag: false.

```tripol
import tripol.binary

show binary.is_byte(255)
show binary.is_byte(256)
```

Documented expected output:

```text
true
false
```

### binary.valid

`func binary.valid(data:ints) -> bool`

EN: Check whether all elements are bytes from 0 to 255.

ES: Comprueba si todos los elementos son bytes de 0 a 255.

Native example flag: false.

```tripol
import tripol.binary

a = [0, 128, 255]
show binary.valid(a)
```

Documented expected output:

```text
true
```

### binary.read_u16_le

`func binary.read_u16_le(data:ints,offset:int) -> int`

EN: Read a 16-bit unsigned number in little-endian order.

ES: Lee un entero sin signo de 16 bits en orden little-endian.

Native example flag: false.

```tripol
import tripol.binary

a = [1, 2]
show binary.read_u16_le(a, 0)
```

Documented expected output:

```text
513
```

### binary.read_u32_le

`func binary.read_u32_le(data:ints,offset:int) -> int`

EN: Read four bytes as a little-endian unsigned integer.

ES: Lee cuatro bytes como entero sin signo en orden little-endian.

Native example flag: false.

```tripol
import tripol.binary

a = [1, 2, 0, 0]
show binary.read_u32_le(a, 0)
```

Documented expected output:

```text
513
```

### binary.u16_le

`func binary.u16_le(value:int) -> ints`

EN: Encode an unsigned 16-bit value with the least significant byte first.

ES: Codifica un entero sin signo de 16 bits con el byte menos significativo primero.

Native example flag: false.

```tripol
import tripol.binary

show binary.u16_le(513)
```

Documented expected output:

```text
[1, 2]
```

### binary.u32_le

`func binary.u32_le(value:int) -> ints`

EN: Encode an unsigned 32-bit value into four little-endian bytes.

ES: Codifica un entero sin signo de 32 bits en cuatro bytes little-endian.

Native example flag: false.

```tripol
import tripol.binary

show binary.u32_le(65536)
```

Documented expected output:

```text
[0, 0, 1, 0]
```

### binary.read_u16_be

`func binary.read_u16_be(data:ints,offset:int) -> int`

EN: Read two bytes with the most significant byte first.

ES: Lee dos bytes con el byte más significativo primero.

Native example flag: false.

```tripol
import tripol.binary

a = [2, 1]
show binary.read_u16_be(a, 0)
```

Documented expected output:

```text
513
```

### binary.u16_be

`func binary.u16_be(value:int) -> ints`

EN: Encode an unsigned 16-bit value in big-endian order.

ES: Codifica un entero sin signo de 16 bits en orden big-endian.

Native example flag: false.

```tripol
import tripol.binary

show binary.u16_be(513)
```

Documented expected output:

```text
[2, 1]
```

---

## tripol.bits — 6 functions

Import: `import tripol.bits`

### bits.mask

`func bits.mask(index:int) -> int`

EN: Create a number with only one bit set.

ES: Crea un número con un único bit activado.

Native example flag: false.

```tripol
import tripol.bits

show bits.mask(3)
```

Documented expected output:

```text
8
```

### bits.test

`func bits.test(n:int,index:int) -> bool`

EN: Check whether a bit is set.

ES: Comprueba si un bit está activado.

Native example flag: false.

```tripol
import tripol.bits

show bits.test(8, 3)
```

Documented expected output:

```text
true
```

### bits.set

`func bits.set(n:int,index:int) -> int`

EN: Turn on a bit without changing other bits.

ES: Activa un bit sin modificar los demás.

Native example flag: false.

```tripol
import tripol.bits

show bits.set(1, 2)
```

Documented expected output:

```text
5
```

### bits.clear

`func bits.clear(n:int,index:int) -> int`

EN: Turn off a bit.

ES: Desactiva un bit.

Native example flag: false.

```tripol
import tripol.bits

show bits.clear(7, 1)
```

Documented expected output:

```text
5
```

### bits.toggle

`func bits.toggle(n:int,index:int) -> int`

EN: Reverse the state of a bit.

ES: Invierte el estado de un bit.

Native example flag: false.

```tripol
import tripol.bits

show bits.toggle(5, 2)
```

Documented expected output:

```text
1
```

### bits.count

`func bits.count(n:int) -> int`

EN: Count the number of set bits.

ES: Cuenta los bits activados.

Native example flag: false.

```tripol
import tripol.bits

show bits.count(7)
```

Documented expected output:

```text
3
```

---

## tripol.bitset — 6 functions

Import: `import tripol.bitset`

### bitset.create

`func bitset.create(count:int) -> ints`

EN: Create an empty compact set with a fixed number of bit positions.

ES: Crea un conjunto compacto vacío con una cantidad fija de posiciones.

Native example flag: false.

```tripol
import tripol.bitset

b = bitset.create(100)
show bitset.count(b)
```

Documented expected output:

```text
0
```

### bitset.capacity

`func bitset.capacity(a:ints) -> int`

EN: Return the number of addressable bit positions.

ES: Devuelve la cantidad de posiciones de bits disponibles.

Native example flag: false.

```tripol
import tripol.bitset

b = bitset.create(100)
show bitset.capacity(b)
```

Documented expected output:

```text
100
```

### bitset.contains

`func bitset.contains(a:ints,index:int) -> bool`

EN: Check whether an indexed bit is set.

ES: Comprueba si está activado un bit por su índice.

Native example flag: false.

```tripol
import tripol.bitset

b = bitset.create(100)
bitset.add(b, 65)
show bitset.contains(b, 65)
```

Documented expected output:

```text
true
```

### bitset.add

`func bitset.add(a:ints,index:int) -> bool`

EN: Set a bit and return true only if it was previously unset.

ES: Activa un bit y devuelve true solo si antes estaba desactivado.

Native example flag: false.

```tripol
import tripol.bitset

b = bitset.create(100)
show bitset.add(b, 65)
show bitset.add(b, 65)
```

Documented expected output:

```text
true
false
```

### bitset.remove

`func bitset.remove(a:ints,index:int) -> bool`

EN: Clear a bit and return whether it was previously set.

ES: Desactiva un bit y devuelve si antes estaba activado.

Native example flag: false.

```tripol
import tripol.bitset

b = bitset.create(100)
bitset.add(b, 65)
show bitset.remove(b, 65)
```

Documented expected output:

```text
true
```

### bitset.count

`func bitset.count(a:ints) -> int`

EN: Count the active bits in the set.

ES: Cuenta los bits activados en el conjunto.

Native example flag: false.

```tripol
import tripol.bitset

b = bitset.create(100)
bitset.add(b, 1)
bitset.add(b, 65)
show bitset.count(b)
```

Documented expected output:

```text
2
```

---

## tripol.checksum — 3 functions

Import: `import tripol.checksum`

### checksum.adler32

`func checksum.adler32(s:string) -> int`

EN: Calculate the Adler-32 checksum of the UTF-8 bytes.

ES: Calcula la suma de comprobación Adler-32 de los bytes UTF-8.

Native example flag: false.

```tripol
import tripol.checksum

show checksum.adler32("abc")
```

Documented expected output:

```text
38600999
```

### checksum.sum8

`func checksum.sum8(s:string) -> int`

EN: Sum bytes modulo 256.

ES: Suma los bytes módulo 256.

Native example flag: false.

```tripol
import tripol.checksum

show checksum.sum8("abc")
```

Documented expected output:

```text
38
```

### checksum.sum16

`func checksum.sum16(s:string) -> int`

EN: Sum bytes modulo 65536.

ES: Suma los bytes módulo 65536.

Native example flag: false.

```tripol
import tripol.checksum

show checksum.sum16("abc")
```

Documented expected output:

```text
294
```

---

## tripol.compression — 3 functions

Import: `import tripol.compression`

### compression.rle_encode

`func compression.rle_encode(a:ints) -> ints`

EN: Replace runs of equal integers with value/count pairs.

ES: Reemplaza secuencias de enteros iguales por parejas valor/cantidad.

Native example flag: false.

```tripol
import tripol.compression

a = [7, 7, 7, 2]
show compression.rle_encode(a)
```

Documented expected output:

```text
[7, 3, 2, 1]
```

### compression.rle_valid

`func compression.rle_valid(a:ints) -> bool`

EN: Check whether RLE pairs can be safely decoded.

ES: Comprueba si las parejas RLE se pueden decodificar correctamente.

Native example flag: false.

```tripol
import tripol.compression

a = [7, 3, 2, 1]
show compression.rle_valid(a)
```

Documented expected output:

```text
true
```

### compression.rle_decode

`func compression.rle_decode(a:ints) -> ints`

EN: Expand RLE value/count pairs into the original array.

ES: Expande las parejas RLE valor/cantidad al arreglo original.

Native example flag: false.

Note: Decoded length is limited to 65535 integers.

Nota: La longitud decodificada se limita a 65535 enteros.

```tripol
import tripol.compression

a = [7, 3, 2, 1]
show compression.rle_decode(a)
```

Documented expected output:

```text
[7, 7, 7, 2]
```

---

## tripol.crypto — 3 functions

Import: `import tripol.crypto`

### crypto.sha256_bytes

`func crypto.sha256_bytes(data:ints) -> ints`

EN: Hash byte values with SHA-256 and return the 32 digest bytes.

ES: Calcula SHA-256 sobre los bytes y devuelve los 32 bytes del resumen.

Native example flag: false.

```tripol
import tripol.crypto

data=ints(1)
data[0]=97
show len(crypto.sha256_bytes(data))
```

Documented expected output:

```text
32
```

### crypto.sha256

`func crypto.sha256(data:ints) -> string`

EN: Hash a byte array with SHA-256 and return its hexadecimal digest.

ES: Calcula SHA-256 sobre el arreglo de bytes y devuelve el resumen hexadecimal.

Native example flag: false.

```tripol
import tripol.crypto

data=ints(1)
data[0]=97
show crypto.sha256(data)
```

Documented expected output:

```text
ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb
```

### crypto.sha256_text

`func crypto.sha256_text(value:string) -> string`

EN: Hash the UTF-8 bytes of a string and return the hexadecimal SHA-256 digest.

ES: Calcula SHA-256 sobre los bytes UTF-8 del texto y devuelve el resumen hexadecimal.

Native example flag: false.

```tripol
import tripol.crypto

show crypto.sha256_text("abc")
```

Documented expected output:

```text
ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
```

---

## tripol.discord — 99 functions + 206 methods

Import: `import tripol.discord`

[Categorized web reference](https://tripollas.github.io/Tripol/docs/libraries/discord.html#functions). The current catalog and native guide take precedence over historical Discord R5 descriptions. Live networking remains experimental.

### Getting started & connection / Primeros pasos y conexión (19)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-start)

#### discord.bot

`func discord.bot() -> Bot`

EN: Create a bot with its own command table and presence settings.

ES: Crea un bot con su propia tabla de comandos y configuración de presencia.

Native example flag: false.

Note: Runs offline on native TRIPOL and Web Engine R3.

Nota: Funciona sin conexión en TRIPOL nativo y Web Engine R3.

```tripol
import tripol.discord

bot = discord.bot()
bot.command("hola", "¡Hola!")
show bot.preview("hola")
```

Documented expected output:

```text
¡Hola!
```

#### discord.gateway

`func discord.gateway() -> string`

EN: Discover the public Discord Gateway address without a bot token.

ES: Consulta la dirección pública del Gateway de Discord sin utilizar un token.

Native example flag: true.

Note: Requires an Internet connection and native HTTPS support. The response is JSON.

Nota: Requiere conexión a Internet y HTTPS nativo. La respuesta es JSON.

Return metadata: ["JSON returned by Discord.","JSON devuelto por Discord."]

```tripol
import tripol.discord

show discord.gateway()
```

#### discord.api_host

`func discord.api_host() -> string`

EN: Get the hostname used by the Discord API.

ES: Obtiene el nombre del servidor de la API de Discord.

Native example flag: true.

```tripol
import tripol.discord

show discord.api_host()
```

Documented expected output:

```text
discord.com
```

#### discord.token

`func discord.token() -> string`

EN: Read DISCORD_TOKEN from the local environment.

ES: Lee DISCORD_TOKEN desde el entorno local.

Native example flag: true.

Note: Set DISCORD_TOKEN locally. This example checks that a token exists without revealing it. Never print the token itself.

Nota: Configura DISCORD_TOKEN localmente. El ejemplo comprueba que existe sin revelarlo. Nunca imprimas el token.

Return metadata: ["Secret token; do not display it.","Token secreto; no lo muestres."]

```tripol
import tripol.discord

secret = discord.token()
show len(secret) > 0
```

Documented expected output:

```text
true
```

#### discord.login

`func discord.login(secret:string) -> bool`

EN: Configure the bot using a token supplied to the function.

ES: Configura el bot utilizando un token recibido por la función.

Native example flag: true.

Note: Set DISCORD_TOKEN locally. Prefer login_env() for everyday use; never hardcode secrets. Configuration does not validate the token with Discord.

Nota: Configura DISCORD_TOKEN localmente. Prefiere login_env() para uso habitual; nunca escribas secretos en el código. Configurar no valida el token con Discord.

```tripol
import tripol.discord

secret = discord.token()
show discord.login(secret)
```

Documented expected output:

```text
true
```

#### discord.login_env

`func discord.login_env() -> bool`

EN: Load the bot token from DISCORD_TOKEN into private local configuration.

ES: Carga el token de DISCORD_TOKEN en la configuración privada local.

Native example flag: true.

Note: Set a valid-format token locally before running. This does not open a network connection.

Nota: Configura localmente un token con formato válido. Esto no abre una conexión de red.

```tripol
import tripol.discord

show discord.login_env()
```

Documented expected output:

```text
true
```

#### discord.reset

`func discord.reset() -> bool`

EN: Clear the local bot configuration, including token, commands and channel.

ES: Borra la configuración local del bot, incluidos token, comandos y canal.

Native example flag: true.

Note: Does not delete anything on Discord.

Nota: No elimina nada en Discord.

```tripol
import tripol.discord

discord.reply("!hola", "Hola!")
show discord.reset()
show discord.preview("!hola") == ""
```

Documented expected output:

```text
true
true
```

#### discord.configured

`func discord.configured() -> bool`

EN: Check whether a token has been stored in the local bot configuration.

ES: Comprueba si hay un token guardado en la configuración local.

Native example flag: true.

Note: A true result does not guarantee the token is accepted by Discord.

Nota: Un resultado true no garantiza que Discord acepte el token.

```tripol
import tripol.discord

discord.reset()
show discord.configured()
```

Documented expected output:

```text
false
```

#### discord.setup

`func discord.setup(secret:string,id:string) -> bool`

EN: Configure the token and channel in one call.

ES: Configura el token y el canal en una sola llamada.

Native example flag: true.

Note: Configure DISCORD_TOKEN locally; replace 123 with the real ID for live use. No network connection is opened.

Nota: Configura DISCORD_TOKEN localmente; sustituye 123 por el ID real para conectarte. No abre una conexión.

```tripol
import tripol.discord

secret = discord.token()
show discord.setup(secret, "123")
```

Documented expected output:

```text
true
```

#### discord.status

`func discord.status(value:string) -> bool`

EN: Choose the presence status sent when the bot identifies itself to Discord.

ES: Elige el estado de presencia que se envía al identificarse el bot.

Native example flag: true.

Note: Allowed: online, idle, dnd, invisible. Configure before run().

Nota: Permitidos: online, idle, dnd, invisible. Configura antes de run().

```tripol
import tripol.discord

show discord.status("idle")
```

Documented expected output:

```text
true
```

#### discord.activity

`func discord.activity(value:string) -> bool`

EN: Choose the activity text displayed by the bot when connecting.

ES: Elige el texto de actividad que muestra el bot al conectarse.

Native example flag: true.

Note: Maximum 128 UTF-8 bytes; empty text clears the activity.

Nota: Máximo 128 bytes UTF-8; el texto vacío elimina la actividad.

```tripol
import tripol.discord

show discord.activity("Programando con TRIPOL")
```

Documented expected output:

```text
true
```

#### discord.connected

`func discord.connected() -> bool`

EN: Check the local Gateway ready state.

ES: Comprueba si el Gateway está listo según su estado local.

Native example flag: true.

Note: This is a local state check, not a fresh network test.

Nota: Comprueba el estado local, no realiza una prueba de red nueva.

```tripol
import tripol.discord

discord.reset()
show discord.connected()
```

Documented expected output:

```text
false
```

#### discord.ping

`func discord.ping() -> int`

EN: Get the most recent Gateway heartbeat round-trip time in milliseconds.

ES: Obtiene la última latencia de heartbeat del Gateway en milisegundos.

Native example flag: true.

Note: Returns -1 until a heartbeat acknowledgement is recorded.

Nota: Devuelve -1 hasta que se registra una confirmación de heartbeat.

```tripol
import tripol.discord

discord.reset()
show discord.ping()
```

Documented expected output:

```text
-1
```

#### discord.run

`func discord.run(intents=33281) -> int`

EN: Start the bot and keep it connected to Discord to answer registered commands.

ES: Inicia el bot y lo mantiene conectado a Discord para responder comandos registrados.

Native example flag: true.

Note: Windows and Linux x86-64. Set DISCORD_TOKEN and enable Message Content Intent for text commands. Linux supports bounded reconnect and Resume.

Nota: Windows y Linux x86-64. Configura DISCORD_TOKEN y activa Message Content Intent para comandos de texto. Linux admite reconexión acotada y Resume.

Return metadata: ["Returns 1 after the Gateway loop closes normally; failures may stop execution.","Devuelve 1 cuando el bucle Gateway termina normalmente; los fallos pueden detener la ejecución."]

```tripol
import tripol.discord

discord.login_env()
discord.status("online")
discord.activity("TRIPOL")
discord.reply("!hola", "Hola!")
discord.run()
```

#### discord.run_session

`func discord.run_session(handler:func(Message)->string,intents:int,keys:string_list,handlers:list[func(Interaction)->string]) -> int`

EN: Run a configured Gateway with a message and interaction registry.

ES: Ejecuta un Gateway con sus manejadores de mensajes e interacciones.

Native example flag: true.

Note: Requires an authorized Discord bot and a native run on Linux or Windows.

Nota: Requiere un bot autorizado de Discord y ejecución nativa en Linux o Windows.

```tripol
import tripol.discord

discord.login_env()
discord.run_session(&discord.ignore_message,1,string_list(),discord.empty_interactions())
```

#### discord.Bot.prefix

`func discord.Bot.prefix(value:string) -> bool`

EN: Set the command prefix. The initial prefix is empty.

ES: Define el prefijo del comando. El prefijo inicial está vacío.

Native example flag: false.

Note: Up to 32 UTF-8 bytes.

Nota: Hasta 32 bytes UTF-8.

```tripol
import tripol.discord

bot = discord.bot()
bot.command("hola", "¡Hola!")
bot.prefix("!")
show bot.preview("!hola")
```

Documented expected output:

```text
¡Hola!
```

#### discord.Bot.status

`func discord.Bot.status(value:string) -> bool`

EN: Set the initial presence to online, idle, dnd or invisible.

ES: Define la presencia inicial como online, idle, dnd o invisible.

Native example flag: false.

```tripol
import tripol.discord

bot = discord.bot()
show bot.status("idle")
```

Documented expected output:

```text
true
```

#### discord.Bot.activity

`func discord.Bot.activity(value:string) -> bool`

EN: Set the activity text shown in the initial presence.

ES: Define el texto de actividad de la presencia inicial.

Native example flag: false.

Note: Up to 128 UTF-8 bytes.

Nota: Hasta 128 bytes UTF-8.

```tripol
import tripol.discord

bot = discord.bot()
show bot.activity("Programando en TRIPOL")
```

Documented expected output:

```text
true
```

#### discord.Bot.start

`func discord.Bot.start() -> int`

EN: Start a blocking Gateway session. Set Message Content Intent when using text commands.

ES: Inicia una sesión Gateway bloqueante. Activa Message Content Intent al usar comandos de texto.

Native example flag: true.

Note: Run on Windows or Linux x86-64 with DISCORD_TOKEN configured.

Nota: Ejecuta en Windows o Linux x86-64 con DISCORD_TOKEN configurado.

```tripol
import tripol.discord

bot = discord.bot()
bot.command("hola", "¡Hola!")
bot.start()
```

### Commands & context menus / Comandos y menús contextuales (55)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-commands)

#### discord.empty_interactions

`func discord.empty_interactions() -> list[func(Interaction)->string]`

EN: Create an empty typed list for interaction handlers.

ES: Crea una lista tipada vacía para callbacks de interacción.

Native example flag: false.

```tripol
import tripol.discord

handlers=discord.empty_interactions()
show handlers.size()
```

Documented expected output:

```text
0
```

#### discord.interaction_selector

`func discord.interaction_selector(event:string) -> string`

EN: Read the dispatch key from an interaction event.

ES: Lee la clave de despacho de un evento de interacción.

Native example flag: false.

```tripol
import tripol.discord

show discord.interaction_selector("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"id\":\"123\",\"token\":\"fixture-token\",\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"options\":[{\"type\":3,\"name\":\"text\",\"value\":\"Ana\"}]}}}")
```

Documented expected output:

```text
s"hola"/
```

#### discord.interaction_matches

`func discord.interaction_matches(key:string,selector:string) -> bool`

EN: Compare a registered key with a Discord interaction selector.

ES: Compara una clave registrada con el selector de un evento.

Native example flag: false.

```tripol
import tripol.discord

show discord.interaction_matches("shola", "s\"hola\"")
```

Documented expected output:

```text
true
```

#### discord.dispatch_interaction

`func discord.dispatch_interaction(event:string,key:string,handler:func(Interaction)->string) -> string`

EN: Invoke a typed callback and build its initial response.

ES: Ejecuta el callback tipado y construye su respuesta inicial.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
show discord.dispatch_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"id\":\"123\",\"token\":\"fixture-token\",\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"options\":[{\"type\":3,\"name\":\"text\",\"value\":\"Ana\"}]}}}","shola",&answer)
```

Documented expected output:

```text
{"type":4,"data":{"content":"Hola, Ana","flags":64,"allowed_mentions":{"parse":[]}}}
```

#### discord.deliver_interaction

`func discord.deliver_interaction(event:string,key:string,handler:func(Interaction)->string) -> bool`

EN: Send the returned response to the interaction callback endpoint.

ES: Envía la respuesta devuelta al endpoint de la interacción.

Native example flag: true.

Note: Requires an authorized Discord bot and a native run on Linux or Windows.

Nota: Requiere un bot autorizado de Discord y ejecución nativa en Linux o Windows.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
discord.deliver_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"id\":\"123\",\"token\":\"fixture-token\",\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"options\":[{\"type\":3,\"name\":\"text\",\"value\":\"Ana\"}]}}}","shola",&answer)
```

#### discord.slash_name

`func discord.slash_name(name:string) -> bool`

EN: Validate the supported slash name syntax.

ES: Valida la sintaxis admitida para un nombre slash.

Native example flag: false.

```tripol
import tripol.discord

show discord.slash_name("hola-2")
```

Documented expected output:

```text
true
```

#### discord.slash_is_current

`func discord.slash_is_current(remote:string,definition:string) -> bool`

EN: Compare a local slash definition with the remote command array.

ES: Compara una definición slash local con el arreglo remoto.

Native example flag: false.

```tripol
import tripol.discord

show discord.slash_is_current("[{\"name\":\"hola\",\"description\":\"Saluda\",\"type\":1}]","{\"name\":\"hola\",\"description\":\"Saluda\",\"type\":1}")
```

Documented expected output:

```text
true
```

#### discord.sync_slash_commands

`func discord.sync_slash_commands(specs:string_list,guilds:string_list) -> bool`

EN: Synchronize application commands in the specified scopes.

ES: Sincroniza comandos de aplicación en los ámbitos indicados.

Native example flag: true.

Note: Requires an authorized Discord bot and a native run on Linux or Windows.

Nota: Requiere un bot autorizado de Discord y ejecución nativa en Linux o Windows.

```tripol
import tripol.discord

discord.login_env()
specs=string_list(["{\"name\":\"hola\",\"description\":\"Saluda\",\"type\":1}"])
show discord.sync_slash_commands(specs,string_list([""]))
```

#### discord.slash_options

`func discord.slash_options() -> SlashOptions`

EN: Create a builder for flat application command options.

ES: Crea un constructor de opciones planas para comandos.

Native example flag: false.

```tripol
import tripol.discord

opts=discord.slash_options()
show opts.payload()
```

Documented expected output:

```text
[]
```

#### discord.slash_option_signature

`func discord.slash_option_signature(options:string,depth=0) -> string`

EN: Normalize supported slash option fields for comparison.

ES: Normaliza los campos de opciones slash para compararlos.

Native example flag: false.

```tripol
import tripol.discord

opts=discord.slash_options()
opts.add("text","Nombre","string",true)
show len(discord.slash_option_signature(opts.payload()))>0
```

Documented expected output:

```text
true
```

#### discord.interaction_message

`func discord.interaction_message(content:string,ephemeral=false) -> string`

EN: Build a text message for edits or follow-ups, with mentions disabled.

ES: Prepara un mensaje de texto para edición o seguimiento con las menciones desactivadas.

Native example flag: false.

```tripol
import tripol.discord

show discord.interaction_message("Hola", true)
```

Documented expected output:

```text
{"content":"Hola","flags":64,"allowed_mentions":{"parse":[]}}
```

#### discord.interaction_context

`func discord.interaction_context(payload:string,key:string) -> Interaction`

EN: Build an offline interaction context from the event data and a dispatch key.

ES: Crea un contexto de interacción sin conexión a partir de los datos del evento y una clave de despacho.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"id\":\"123\"}","shola")
show ctx.id,ctx.name
```

Documented expected output:

```text
123 hola
```

#### discord.interaction_option

`func discord.interaction_option(data:string,name:string,focused=false) -> string`

EN: Read the JSON object for a named or focused option, including within a subcommand.

ES: Lee el objeto JSON de una opción o del campo activo, incluso dentro de un subcomando.

Native example flag: false.

```tripol
import tripol.discord

show discord.interaction_option("{\"type\":1,\"name\":\"order\",\"options\":[{\"type\":2,\"name\":\"drinks\",\"options\":[{\"type\":1,\"name\":\"find\",\"options\":[{\"name\":\"query\",\"type\":3,\"value\":\"ca\",\"focused\":true}]}]}]}","",true)
```

Documented expected output:

```text
{"name":"query","type":3,"value":"ca","focused":true}
```

#### discord.first_command_option

`func discord.first_command_option(data:string) -> string`

EN: Read the first option object in interaction or subcommand data.

ES: Lee el primer objeto de opción en los datos de una interacción o subcomando.

Native example flag: false.

```tripol
import tripol.discord

show discord.first_command_option("{\"type\":1,\"name\":\"order\",\"options\":[{\"type\":2,\"name\":\"drinks\",\"options\":[{\"type\":1,\"name\":\"find\",\"options\":[{\"name\":\"query\",\"type\":3,\"value\":\"ca\",\"focused\":true}]}]}]}")
```

Documented expected output:

```text
{"type":2,"name":"drinks","options":[{"type":1,"name":"find","options":[{"name":"query","type":3,"value":"ca","focused":true}]}]}
```

#### discord.command_choices

`func discord.command_choices() -> CommandChoices`

EN: Create a typed choice list for a command or autocomplete callback.

ES: Crea una lista tipada de elecciones para un comando o autocompletado.

Native example flag: false.

```tripol
import tripol.discord

choices=discord.command_choices()
show choices.payload()
```

Documented expected output:

```text
[]
```

#### discord.slash_choice_signature

`func discord.slash_choice_signature(choices:string) -> string`

EN: Normalize a Discord choice list for comparison when syncing commands.

ES: Normaliza una lista de elecciones para compararla al sincronizar comandos.

Native example flag: false.

```tripol
import tripol.discord

show discord.slash_choice_signature("[{\"value\":10,\"name\":\"D\\u00edez\"}]")
```

Documented expected output:

```text
"Díez":1e1;
```

#### discord.slash_number_signature

`func discord.slash_number_signature(value:string) -> string`

EN: Normalize equivalent Discord number spellings for command sync.

ES: Normaliza representaciones equivalentes de números para sincronizar comandos.

Native example flag: false.

```tripol
import tripol.discord

show discord.slash_number_signature("1.250")==discord.slash_number_signature("125e-2")
```

Documented expected output:

```text
true
```

#### discord.Bot.command

`func discord.Bot.command(name:string, response:string) -> bool`

EN: Register an exact text command or replace its reply.

ES: Registra un comando de texto exacto o sustituye su respuesta.

Native example flag: false.

Note: Up to 16 names. Names and replies accept 1–96 and 1–512 UTF-8 bytes, respectively.

Nota: Hasta 16 nombres. Los nombres y las respuestas aceptan 1–96 y 1–512 bytes UTF-8, respectivamente.

```tripol
import tripol.discord

bot = discord.bot()
bot.command("ping", "Pong!")
show bot.preview("ping")
```

Documented expected output:

```text
Pong!
```

#### discord.Bot.preview

`func discord.Bot.preview(content:string) -> string`

EN: Look up an exact reply locally, with the configured prefix.

ES: Busca la respuesta exacta localmente con el prefijo configurado.

Native example flag: false.

Note: Runs offline on native TRIPOL and Web Engine R3.

Nota: Funciona sin conexión en TRIPOL nativo y Web Engine R3.

```tripol
import tripol.discord

bot = discord.bot()
bot.command("ping", "Pong!")
show bot.preview("ping"), bot.preview("other") == ""
```

Documented expected output:

```text
Pong! true
```

#### discord.Bot.bind_interaction

`func discord.Bot.bind_interaction(key:string,handler:func(Interaction)->string,definition="",guild_id="") -> bool`

EN: Bind a callback to a protocol-prefixed route key.

ES: Asocia un callback a una clave con prefijo de protocolo.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
show b.bind_interaction("shola",&answer)
```

Documented expected output:

```text
true
```

#### discord.Bot.slash

`func discord.Bot.slash(name:string,description:string,handler:func(Interaction)->string,guild_id="",default_permissions="") -> bool`

EN: Register a typed slash command in a global or guild scope.

ES: Registra un comando slash tipado global o por servidor.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
b.slash("hola","Saluda",&answer)
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"id\":\"123\",\"token\":\"fixture-token\",\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"options\":[{\"type\":3,\"name\":\"text\",\"value\":\"Ana\"}]}}}")
```

Documented expected output:

```text
{"type":4,"data":{"content":"Hola, Ana","flags":64,"allowed_mentions":{"parse":[]}}}
```

#### discord.Bot.slash_with_options

`func discord.Bot.slash_with_options(name:string,description:string,handler:func(Interaction)->string,options:SlashOptions,guild_id="",default_permissions="") -> bool`

EN: Register a slash command with a typed option schema.

ES: Registra un comando slash con opciones tipadas.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
opts=discord.slash_options()
opts.add("text","Nombre","string",true)
b.slash_with_options("hola","Saluda",&answer,opts,"123")
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"id\":\"123\",\"token\":\"fixture-token\",\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"guild_id\":\"123\",\"options\":[{\"type\":3,\"name\":\"text\",\"value\":\"Ana\"}]}}}")
```

Documented expected output:

```text
{"type":4,"data":{"content":"Hola, Ana","flags":64,"allowed_mentions":{"parse":[]}}}
```

#### discord.Bot.preview_interaction

`func discord.Bot.preview_interaction(event:string) -> string`

EN: Execute a typed interaction callback without connecting.

ES: Ejecuta un callback tipado de interacción sin conexión.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
b.slash("hola","Saluda",&answer)
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"id\":\"123\",\"token\":\"fixture-token\",\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"options\":[{\"type\":3,\"name\":\"text\",\"value\":\"Ana\"}]}}}")
```

Documented expected output:

```text
{"type":4,"data":{"content":"Hola, Ana","flags":64,"allowed_mentions":{"parse":[]}}}
```

#### discord.Interaction.reply

`func discord.Interaction.reply(content:string,ephemeral=false) -> string
func discord.Interaction.reply(embed:Embed,ephemeral=false) -> string
func discord.Interaction.reply(batch:Embeds,ephemeral=false) -> string
func discord.Interaction.reply(layout:Components,ephemeral=false) -> string`

EN: Build an initial text, embed, embed batch or Components V2 reply.

ES: Construye una respuesta inicial de texto, embed, lote de embeds o Components V2.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
b.slash("hola","Saluda",&answer)
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"id\":\"123\",\"token\":\"fixture-token\",\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"options\":[{\"type\":3,\"name\":\"text\",\"value\":\"Ana\"}]}}}")
```

Documented expected output:

```text
{"type":4,"data":{"content":"Hola, Ana","flags":64,"allowed_mentions":{"parse":[]}}}
```

#### discord.Interaction.defer

`func discord.Interaction.defer(ephemeral=false) -> string`

EN: Acknowledge an interaction for a later response.

ES: Confirma una interacción para responder más tarde.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
show ctx.defer(true)
```

Documented expected output:

```text
{"type":5,"data":{"flags":64}}
```

#### discord.Interaction.ack

`func discord.Interaction.ack() -> string`

EN: Acknowledge a component interaction on its existing message.

ES: Confirma la interacción de un componente sobre su mensaje.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
show ctx.ack()
```

Documented expected output:

```text
{"type":6}
```

#### discord.Interaction.option

`func discord.Interaction.option(name:string) -> string`

EN: Read a raw JSON command option, including inside a selected subcommand.

ES: Lee una opción como JSON, incluso dentro del subcomando seleccionado.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(2,"hola","123","456","789","123","567","{\"options\":[{\"name\":\"text\",\"value\":\"Ana\"}]}",false)
show ctx.option("text")
```

Documented expected output:

```text
"Ana"
```

#### discord.Interaction.option_text

`func discord.Interaction.option_text(name:string) -> string`

EN: Decode a string option from the interaction.

ES: Decodifica una opción de texto de la interacción.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(2,"hola","123","456","789","123","567","{\"options\":[{\"name\":\"text\",\"value\":\"Ana\"}]}",false)
show ctx.option_text("text")
```

Documented expected output:

```text
Ana
```

#### discord.SlashOptions.add

`func discord.SlashOptions.add(name:string,description:string,kind="string",required=false,autocomplete=false) -> bool`

EN: Add a typed flat slash option with optional requiredness.

ES: Añade una opción slash tipada y permite marcarla requerida.

Native example flag: false.

```tripol
import tripol.discord

opts=discord.slash_options()
show opts.add("text","Nombre","string",true)
```

Documented expected output:

```text
true
```

#### discord.SlashOptions.payload

`func discord.SlashOptions.payload() -> string`

EN: Serialize the declared command options.

ES: Serializa las opciones declaradas del comando.

Native example flag: false.

```tripol
import tripol.discord

opts=discord.slash_options()
opts.add("text","Nombre","string",true)
show len(opts.payload())>2
```

Documented expected output:

```text
true
```

#### discord.Interaction.respond

`func discord.Interaction.respond(body:string) -> bool`

EN: Send the initial response once before continuing in a connected handler.

ES: Envía la respuesta inicial una sola vez antes de continuar en un callback conectado.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

func report(ctx:discord.Interaction)->string {
    if ctx.respond(ctx.defer(true)) {
        ctx.edit_reply("Informe listo")
    }
    return ""
}
bot=discord.bot()
bot.slash("informe","Prepara un informe",&report)
bot.start()
```

#### discord.Interaction.edit_response

`func discord.Interaction.edit_response(body:string) -> bool`

EN: Edit the original interaction reply with a Discord message payload.

ES: Edita la respuesta original de la interacción con un payload de mensaje de Discord.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

func report(ctx:discord.Interaction)->string {
    if ctx.respond(ctx.defer(true)) {
        ctx.edit_response(discord.interaction_message("Informe listo"))
    }
    return ""
}
bot=discord.bot()
bot.slash("informe","Prepara un informe",&report)
bot.start()
```

#### discord.Interaction.edit_reply

`func discord.Interaction.edit_reply(content:string) -> bool`

EN: Fill a deferred response or update the text of the original reply.

ES: Completa una respuesta diferida o actualiza el texto de la respuesta original.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

func report(ctx:discord.Interaction)->string {
    if ctx.respond(ctx.defer(true)) {
        ctx.edit_reply("Informe listo")
    }
    return ""
}
bot=discord.bot()
bot.slash("informe","Prepara un informe",&report)
bot.start()
```

#### discord.Interaction.followup

`func discord.Interaction.followup(content:string,ephemeral=false) -> string`

EN: Send a text follow-up and return its message ID, or empty text on failure.

ES: Envía un seguimiento de texto y devuelve el ID del mensaje, o una cadena vacía si falla.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

func report(ctx:discord.Interaction)->string {
    if ctx.respond(ctx.reply("Procesado",true)) {
        show ctx.followup("Informe listo", true)
    }
    return ""
}
bot=discord.bot()
bot.slash("informe","Prepara un informe",&report)
bot.start()
```

#### discord.Interaction.command_group

`func discord.Interaction.command_group() -> string`

EN: Read the selected slash subcommand group name.

ES: Lee el nombre del grupo de subcomandos seleccionado.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":4,\"data\":{\"type\":1,\"name\":\"order\",\"options\":[{\"type\":2,\"name\":\"drinks\",\"options\":[{\"type\":1,\"name\":\"find\",\"options\":[{\"name\":\"query\",\"type\":3,\"value\":\"ca\",\"focused\":true}]}]}]}}","aorder")
show ctx.command_group()
```

Documented expected output:

```text
drinks
```

#### discord.Interaction.subcommand

`func discord.Interaction.subcommand() -> string`

EN: Read a selected slash subcommand, including a grouped subcommand.

ES: Lee el subcomando slash seleccionado, incluso dentro de un grupo.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":4,\"data\":{\"type\":1,\"name\":\"order\",\"options\":[{\"type\":2,\"name\":\"drinks\",\"options\":[{\"type\":1,\"name\":\"find\",\"options\":[{\"name\":\"query\",\"type\":3,\"value\":\"ca\",\"focused\":true}]}]}]}}","aorder")
show ctx.subcommand()
```

Documented expected output:

```text
find
```

#### discord.Interaction.focused

`func discord.Interaction.focused() -> string`

EN: Read the option name currently being filled during autocomplete.

ES: Lee el nombre de la opción que se está rellenando.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":4,\"data\":{\"type\":1,\"name\":\"order\",\"options\":[{\"type\":2,\"name\":\"drinks\",\"options\":[{\"type\":1,\"name\":\"find\",\"options\":[{\"name\":\"query\",\"type\":3,\"value\":\"ca\",\"focused\":true}]}]}]}}","aorder")
show ctx.focused()
```

Documented expected output:

```text
query
```

#### discord.Bot.on_autocomplete

`func discord.Bot.on_autocomplete(name:string,handler:func(Interaction)->string,guild_id="") -> bool`

EN: Register an autocomplete handler for a global or server command.

ES: Registra un callback de autocompletado global o por servidor.

Native example flag: false.

```tripol
import tripol.discord

func suggest(ctx:discord.Interaction)->string {
    choices=discord.command_choices()
    choices.add("Café","cafe")
    return ctx.autocomplete(choices)
}
b=discord.bot()
b.on_autocomplete("order",&suggest)
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":4,\"data\":{\"type\":1,\"name\":\"order\",\"options\":[{\"type\":2,\"name\":\"drinks\",\"options\":[{\"type\":1,\"name\":\"find\",\"options\":[{\"name\":\"query\",\"type\":3,\"value\":\"ca\",\"focused\":true}]}]}]}}}")
```

Documented expected output:

```text
{"type":8,"data":{"choices":[{"name":"Café","value":"cafe"}]}}
```

#### discord.SlashOptions.subcommand

`func discord.SlashOptions.subcommand(name:string,description:string,options:SlashOptions) -> bool`

EN: Add a subcommand with its argument schema.

ES: Añade un subcomando con su esquema de argumentos.

Native example flag: false.

```tripol
import tripol.discord

arguments=discord.slash_options()
arguments.add("query","Qué buscar","string",true,true)
commands=discord.slash_options()
show commands.subcommand("find","Busca",arguments)
```

Documented expected output:

```text
true
```

#### discord.SlashOptions.group

`func discord.SlashOptions.group(name:string,description:string,subcommands:SlashOptions) -> bool`

EN: Add a group containing one or more subcommands.

ES: Añade un grupo con uno o más subcomandos.

Native example flag: false.

```tripol
import tripol.discord

arguments=discord.slash_options()
arguments.add("query","Qué buscar","string",true,true)
commands=discord.slash_options()
commands.subcommand("find","Busca",arguments)
root=discord.slash_options()
show root.group("drinks","Bebidas",commands)
```

Documented expected output:

```text
true
```

#### discord.SlashOptions.choose

`func discord.SlashOptions.choose(name:string,description:string,choices:CommandChoices,required=false) -> bool`

EN: Add a fixed text or integer choice option.

ES: Añade una opción con elecciones fijas de texto o de entero.

Native example flag: false.

```tripol
import tripol.discord

choices=discord.command_choices()
choices.integer("Diez",10)
arguments=discord.slash_options()
show arguments.choose("limit","Límite",choices)
```

Documented expected output:

```text
true
```

#### discord.CommandChoices.add

`func discord.CommandChoices.add(label:string,value:string) -> bool`

EN: Add a text value and its displayed label.

ES: Añade un valor de texto con su etiqueta visible.

Native example flag: false.

```tripol
import tripol.discord

choices=discord.command_choices()
show choices.add("Café","cafe")
```

Documented expected output:

```text
true
```

#### discord.CommandChoices.integer

`func discord.CommandChoices.integer(label:string,value:int) -> bool`

EN: Add an integer value and its displayed label.

ES: Añade un valor entero con su etiqueta visible.

Native example flag: false.

```tripol
import tripol.discord

choices=discord.command_choices()
show choices.integer("Diez",10)
```

Documented expected output:

```text
true
```

#### discord.CommandChoices.payload

`func discord.CommandChoices.payload() -> string`

EN: Serialize a list of choices as Discord JSON.

ES: Serializa una lista de elecciones en el JSON de Discord.

Native example flag: false.

```tripol
import tripol.discord

choices=discord.command_choices()
choices.integer("Diez",10)
show choices.payload()
```

Documented expected output:

```text
[{"name":"Diez","value":10}]
```

#### discord.Interaction.autocomplete

`func discord.Interaction.autocomplete(choices:CommandChoices) -> string`

EN: Build the type-8 callback response from matching choices.

ES: Genera una respuesta de autocompletado tipo 8 con elecciones compatibles.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":4,\"data\":{\"type\":1,\"name\":\"order\",\"options\":[{\"type\":2,\"name\":\"drinks\",\"options\":[{\"type\":1,\"name\":\"find\",\"options\":[{\"name\":\"query\",\"type\":3,\"value\":\"ca\",\"focused\":true}]}]}]}}","aorder")
choices=discord.command_choices()
choices.add("Café","cafe")
show ctx.autocomplete(choices)
```

Documented expected output:

```text
{"type":8,"data":{"choices":[{"name":"Café","value":"cafe"}]}}
```

#### discord.Interaction.edit_followup

`func discord.Interaction.edit_followup(message_id:string,content:string) -> bool`

EN: Edit the text of an interaction follow-up.

ES: Edita el texto de un seguimiento de interacción.

Native example flag: true.

Note: Run in a native callback with an authorized Discord app.

Nota: Ejecuta en un callback nativo con una aplicación Discord autorizada.

```tripol
import tripol.discord

func reply(ctx) {
    if ctx.respond(ctx.reply("Hola", true)) {
        id=ctx.followup("Mensaje")
        if id!="" {
            ctx.edit_followup(id,"Actualizado")
        }
    }
    return ""
}
bot=discord.bot()
bot.slash("hola","Saluda",reply)
bot.start()
```

#### discord.Interaction.edit_followup_payload

`func discord.Interaction.edit_followup_payload(message_id:string,body:string) -> bool`

EN: Edit a follow-up with a message payload.

ES: Edita un seguimiento con un payload de mensaje.

Native example flag: true.

Note: Run in a native callback with an authorized Discord app.

Nota: Ejecuta en un callback nativo con una aplicación Discord autorizada.

```tripol
import tripol.discord

func reply(ctx) {
    if ctx.respond(ctx.reply("Hola", true)) {
        id=ctx.followup("Mensaje")
        if id!="" {
            ctx.edit_followup_payload(id,discord.interaction_message("Actualizado"))
        }
    }
    return ""
}
bot=discord.bot()
bot.slash("hola","Saluda",reply)
bot.start()
```

#### discord.Interaction.delete_followup

`func discord.Interaction.delete_followup(message_id:string) -> bool`

EN: Delete a follow-up by its message ID.

ES: Borra un seguimiento por ID de mensaje.

Native example flag: true.

Note: Run in a native callback with an authorized Discord app.

Nota: Ejecuta en un callback nativo con una aplicación Discord autorizada.

```tripol
import tripol.discord

func reply(ctx) {
    if ctx.respond(ctx.reply("Hola", true)) {
        id=ctx.followup("Mensaje")
        if id!="" {
            ctx.delete_followup(id)
        }
    }
    return ""
}
bot=discord.bot()
bot.slash("hola","Saluda",reply)
bot.start()
```

#### discord.Bot.user_menu

`func discord.Bot.user_menu(name:string,handler:func(Interaction)->string,guild_id="",default_permissions="") -> bool`

EN: Register a typed user context menu.

ES: Registra un menú contextual de usuario tipado.

Native example flag: false.

```tripol
import tripol.discord

func selected(ctx) { ctx.reply("ID: "+ctx.target_id(),true) }
bot=discord.bot()
show bot.user_menu("Ver perfil",selected,"123","8192")
```

Documented expected output:

```text
true
```

#### discord.Bot.message_menu

`func discord.Bot.message_menu(name:string,handler:func(Interaction)->string,guild_id="",default_permissions="") -> bool`

EN: Register a typed message context menu.

ES: Registra un menú contextual de mensaje tipado.

Native example flag: false.

```tripol
import tripol.discord

func selected(ctx) { ctx.reply("ID: "+ctx.target_id(),true) }
bot=discord.bot()
show bot.message_menu("Ver mensaje",selected)
```

Documented expected output:

```text
true
```

#### discord.Interaction.target_id

`func discord.Interaction.target_id() -> string`

EN: Return the selected Discord user or message ID.

ES: Devuelve el ID del usuario o mensaje seleccionado.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"data\":{\"type\":2,\"target_id\":\"42\",\"resolved\":{\"users\":{\"42\":{\"id\":\"42\",\"username\":\"Ada\"}}}}}","uPerfil/")
show ctx.target_id()
```

Documented expected output:

```text
42
```

#### discord.Interaction.target

`func discord.Interaction.target() -> string`

EN: Read the selected user or message from the resolved interaction data.

ES: Lee el usuario o mensaje seleccionado de los datos resueltos.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"data\":{\"type\":2,\"target_id\":\"42\",\"resolved\":{\"users\":{\"42\":{\"id\":\"42\",\"username\":\"Ada\"}}}}}","uPerfil/")
show ctx.target()
```

Documented expected output:

```text
{"id":"42","username":"Ada"}
```

#### discord.Interaction.target_username

`func discord.Interaction.target_username() -> string`

EN: Read the username selected by a user context menu.

ES: Lee el nombre de usuario seleccionado en un menú de usuario.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"guild_id\":\"123\",\"member\":{\"permissions\":\"8192\"},\"data\":{\"type\":2,\"target_id\":\"42\",\"resolved\":{\"users\":{\"42\":{\"id\":\"42\",\"username\":\"Ada\"}}}}}","uProfile")
show ctx.target_username()
```

Documented expected output:

```text
Ada
```

#### discord.Interaction.target_content

`func discord.Interaction.target_content() -> string`

EN: Read the text selected by a message context menu.

ES: Lee el texto seleccionado en un menú de mensaje.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"data\":{\"type\":3,\"target_id\":\"42\",\"resolved\":{\"messages\":{\"42\":{\"id\":\"42\",\"content\":\"Hola\"}}}}}","pQuote")
show ctx.target_content()
```

Documented expected output:

```text
Hola
```

#### discord.CommandChoices.number

`func discord.CommandChoices.number(label:string,value:string) -> bool`

EN: Add a fractional choice from decimal text.

ES: Añade una elección fraccionaria a partir de texto decimal.

Native example flag: false.

```tripol
import tripol.discord

choices=discord.command_choices()
choices.number("Half","0.5")
show choices.payload()
```

Documented expected output:

```text
[{"name":"Half","value":0.5}]
```

### Events & messages / Eventos y mensajes (22)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-events)

#### discord.event_channel

`func discord.event_channel(event:string) -> string`

EN: Extract the channel ID of a MESSAGE_CREATE Gateway event and apply the channel filter.

ES: Extrae el ID de canal de un evento Gateway MESSAGE_CREATE y aplica el filtro de canal.

Native example flag: true.

Return metadata: ["Channel ID text or empty text if the event is rejected.","ID de canal como texto o texto vacío si el evento se rechaza."]

```tripol
import tripol.discord

event = "{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"!hola\",\"author\":{\"bot\":false}}}"
show discord.event_channel(event)
```

Documented expected output:

```text
123
```

#### discord.gateway_reply

`func discord.gateway_reply(event:string) -> string`

EN: Resolve a registered reply from a Gateway event without opening a connection.

ES: Busca la respuesta registrada para un evento Gateway sin abrir una conexión.

Native example flag: true.

Note: Ignores bots, webhooks and events outside the configured channel filter.

Nota: Ignora bots, webhooks y eventos fuera del filtro de canal.

```tripol
import tripol.discord

discord.reply("!hola", "Hola!")
event = "{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"!hola\",\"author\":{\"bot\":false}}}"
show discord.gateway_reply(event)
```

Documented expected output:

```text
Hola!
```

#### discord.ignore_message

`func discord.ignore_message(message:Message) -> string`

EN: Return an empty response for the default message handler.

ES: Devuelve texto vacío para el manejador predeterminado.

Native example flag: false.

```tripol
import tripol.discord

m = discord.Message("","","123","456")
show discord.ignore_message(m) == ""
```

Documented expected output:

```text
true
```

#### discord.event_id

`func discord.event_id(raw:string) -> string`

EN: Extract a numeric Discord ID from a quoted JSON value.

ES: Extrae un ID numérico de Discord a partir de un valor JSON entre comillas.

Native example flag: false.

```tripol
import tripol.discord

show discord.event_id("\"123\"")
```

Documented expected output:

```text
123
```

#### discord.dispatch_message

`func discord.dispatch_message(event:string,handler:func(Message)->string) -> string`

EN: Validate a MESSAGE_CREATE event and invoke a typed message callback.

ES: Valida un evento MESSAGE_CREATE y ejecuta un callback tipado.

Native example flag: false.

```tripol
import tripol.discord

func answer(m:discord.Message)->string {
    if m.matches("hola") {return m.reply("¡Hola, " + m.author_id + "!")}
    return ""
}
show discord.dispatch_message("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"hola\",\"author\":{\"id\":\"456\"}}}", &answer)
```

Documented expected output:

```text
¡Hola, 456!
```

#### discord.run_handler

`func discord.run_handler(handler:func(Message)->string,intents=33281) -> int`

EN: Run the configured Gateway with a typed message handler.

ES: Ejecuta el Gateway configurado con un manejador de mensajes tipado.

Native example flag: true.

Note: Run on Windows or Linux x86-64 with DISCORD_TOKEN configured.

Nota: Ejecuta en Windows o Linux x86-64 con DISCORD_TOKEN configurado.

```tripol
import tripol.discord

func answer(m:discord.Message)->string {
    if m.matches("hola") {return m.reply("¡Hola, " + m.author_id + "!")}
    return ""
}
discord.login_env()
discord.run_handler(&answer)
```

#### discord.ignore_gateway_event

`func discord.ignore_gateway_event(event:GatewayEvent) -> bool`

EN: Return false for an event without an active handler.

ES: Devuelve false si un evento carece de manejador activo.

Native example flag: false.

```tripol
import tripol.discord

event=discord.GatewayEvent("MESSAGE_DELETE",8,"321","123","","{\"id\":\"999\"}")
show discord.ignore_gateway_event(event)
```

Documented expected output:

```text
false
```

#### discord.events

`func discord.events() -> GatewayEvents`

EN: Create an independent registry for Gateway events and intents.

ES: Crea un registro independiente para eventos e intents del Gateway.

Native example flag: false.

```tripol
import tripol.discord

registry=discord.events()
show registry.intent_mask()
```

Documented expected output:

```text
1
```

#### discord.gateway_intent

`func discord.gateway_intent(name:string) -> int`

EN: Look up the bit index of an official Gateway intent.

ES: Consulta la posición del bit de un intent oficial del Gateway.

Native example flag: false.

```tripol
import tripol.discord

show discord.gateway_intent("GUILD_MEMBERS")
```

Documented expected output:

```text
1
```

#### discord.gateway_event_name

`func discord.gateway_event_name(name:string) -> string`

EN: Translate a member alias or validate an uppercase dispatch name.

ES: Traduce un alias de miembros o valida un nombre de evento en mayúsculas.

Native example flag: false.

```tripol
import tripol.discord

show discord.gateway_event_name("member_join")
```

Documented expected output:

```text
GUILD_MEMBER_ADD
```

#### discord.deliver_gateway_event

`func discord.deliver_gateway_event(event:string,name:string,handler:func(GatewayEvent)->bool) -> bool`

EN: Invoke a typed handler for a matching Gateway dispatch.

ES: Ejecuta un callback tipado para un evento del Gateway que coincide.

Native example flag: false.

```tripol
import tripol.discord

func ready(event:discord.GatewayEvent)->bool {return event.name=="READY"}
show discord.deliver_gateway_event("{\"op\":0,\"t\":\"READY\",\"s\":1,\"d\":{}}","READY",&ready)
```

Documented expected output:

```text
true
```

#### discord.post_message

`func discord.post_message(channel_id:string,content:string) -> string`

EN: Post a text message with mentions disabled and return its ID.

ES: Envía un mensaje con las menciones desactivadas y devuelve su ID.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

show discord.post_message("123","Hola")
```

#### discord.Message.matches

`func discord.Message.matches(content:string) -> bool`

EN: Compare the decoded message content with exact UTF-8 text.

ES: Compara el contenido decodificado del mensaje con texto UTF-8 exacto.

Native example flag: false.

```tripol
import tripol.discord

func answer(m:discord.Message)->string {
    if m.matches("hola") {return m.reply("¡Hola, " + m.author_id + "!")}
    return ""
}
bot = discord.bot()
bot.on_message(&answer)
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"hola\",\"author\":{\"id\":\"456\"}}}")
```

Documented expected output:

```text
¡Hola, 456!
```

#### discord.Message.reply

`func discord.Message.reply(content:string) -> string`

EN: Validate and return the message a handler wants to send.

ES: Valida y devuelve el mensaje que quiere enviar el manejador.

Native example flag: false.

```tripol
import tripol.discord

func answer(m:discord.Message)->string {
    if m.matches("hola") {return m.reply("¡Hola, " + m.author_id + "!")}
    return ""
}
bot = discord.bot()
bot.on_message(&answer)
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"hola\",\"author\":{\"id\":\"456\"}}}")
```

Documented expected output:

```text
¡Hola, 456!
```

#### discord.Bot.on_message

`func discord.Bot.on_message(handler:func(Message)->string) -> bool`

EN: Register a typed callback for messages that do not match fixed commands.

ES: Registra un callback tipado para mensajes que no coinciden con comandos fijos.

Native example flag: false.

```tripol
import tripol.discord

func answer(m:discord.Message)->string {
    if m.matches("hola") {return m.reply("¡Hola, " + m.author_id + "!")}
    return ""
}
bot = discord.bot()
bot.on_message(&answer)
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"hola\",\"author\":{\"id\":\"456\"}}}")
```

Documented expected output:

```text
¡Hola, 456!
```

#### discord.Bot.preview_event

`func discord.Bot.preview_event(event:string) -> string`

EN: Dispatch an example Gateway event offline.

ES: Procesa sin conexión un evento de ejemplo del Gateway.

Native example flag: false.

```tripol
import tripol.discord

func answer(m:discord.Message)->string {
    if m.matches("hola") {return m.reply("¡Hola, " + m.author_id + "!")}
    return ""
}
bot = discord.bot()
bot.on_message(&answer)
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"hola\",\"author\":{\"id\":\"456\"}}}")
```

Documented expected output:

```text
¡Hola, 456!
```

#### discord.GatewayEvent.field

`func discord.GatewayEvent.field(name:string) -> string`

EN: Read a property of the event body as raw JSON.

ES: Lee una propiedad del cuerpo del evento como JSON.

Native example flag: false.

```tripol
import tripol.discord

event=discord.GatewayEvent("MESSAGE_DELETE",8,"321","123","","{\"id\":\"999\"}")
show event.field("id")
```

Documented expected output:

```text
"999"
```

#### discord.GatewayEvents.intent

`func discord.GatewayEvents.intent(name:string) -> bool`

EN: Enable a named Gateway intent in the registry.

ES: Activa un intent del Gateway por nombre.

Native example flag: false.

```tripol
import tripol.discord

registry=discord.events()
show registry.intent("GUILD_MESSAGE_REACTIONS")
```

Documented expected output:

```text
true
```

#### discord.GatewayEvents.intent_mask

`func discord.GatewayEvents.intent_mask() -> int`

EN: Read the bit mask to identify with the registry intents.

ES: Lee la máscara de bits configurada en el registro.

Native example flag: false.

```tripol
import tripol.discord

registry=discord.events()
registry.intent("GUILD_MEMBERS")
show registry.intent_mask()
```

Documented expected output:

```text
3
```

#### discord.GatewayEvents.on

`func discord.GatewayEvents.on(name:string,handler:func(GatewayEvent)->bool) -> bool`

EN: Register or replace the typed callback for a Gateway event.

ES: Registra o reemplaza el callback tipado de un evento del Gateway.

Native example flag: false.

```tripol
import tripol.discord

func ready(event:discord.GatewayEvent)->bool {return event.name=="READY"}
registry=discord.events()
show registry.on("READY",&ready)
```

Documented expected output:

```text
true
```

#### discord.GatewayEvents.preview

`func discord.GatewayEvents.preview(event:string) -> bool`

EN: Dispatch a fixture through the registry without live credentials.

ES: Despacha un evento de prueba en el registro sin credenciales.

Native example flag: false.

```tripol
import tripol.discord

func ready(event:discord.GatewayEvent)->bool {return event.name=="READY"}
registry=discord.events()
registry.on("READY",&ready)
show registry.preview("{\"op\":0,\"t\":\"READY\",\"s\":1,\"d\":{}}")
```

Documented expected output:

```text
true
```

#### discord.Bot.start_with_events

`func discord.Bot.start_with_events(registry:GatewayEvents) -> int`

EN: Start a native bot and its Gateway event registry.

ES: Inicia un bot nativo y su registro de eventos del Gateway.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

func ready(event:discord.GatewayEvent)->bool {return event.name=="READY"}
registry=discord.events()
registry.on("READY",&ready)
bot=discord.bot()
bot.start_with_events(registry)
```

### Embeds / Embeds (23)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-embeds)

#### discord.embed

`func discord.embed() -> Embed`

EN: Create an empty Discord embed builder.

ES: Crea un constructor de embeds de Discord vacío.

Native example flag: false.

Note: Runs offline on native TRIPOL and Web Engine R3.

Nota: Funciona sin conexión en TRIPOL nativo y Web Engine R3.

```tripol
import tripol.discord

embed = discord.embed()
embed.title("Bienvenido")
show embed.payload()
```

Documented expected output:

```text
{"embeds":[{"title":"Bienvenido"}],"allowed_mentions":{"parse":[]}}
```

#### discord.media_url

`func discord.media_url(value:string) -> bool`

EN: Validate the HTTP and HTTPS URL format supported by Discord builders.

ES: Valida el formato HTTP y HTTPS admitido por los constructores de Discord.

Native example flag: false.

```tripol
import tripol.discord

show discord.media_url("https://example.org/image.png")
```

Documented expected output:

```text
true
```

#### discord.timestamp_valid

`func discord.timestamp_valid(value:string) -> bool`

EN: Validate a calendar date and time with an RFC 3339 offset.

ES: Valida una fecha con zona horaria RFC 3339.

Native example flag: false.

```tripol
import tripol.discord

show discord.timestamp_valid("2028-02-29T10:15:30Z")
```

Documented expected output:

```text
true
```

#### discord.embeds

`func discord.embeds() -> Embeds`

EN: Create a batch of up to ten traditional embeds.

ES: Crea un lote de hasta diez embeds tradicionales.

Native example flag: false.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
batch=discord.embeds()
batch.add(e)
show len(batch.payload())>0
```

Documented expected output:

```text
true
```

#### discord.rgb_color

`func discord.rgb_color(value:string) -> int`

EN: Convert a #RRGGBB color into its RGB integer.

ES: Convierte un color #RRGGBB a su entero RGB.

Native example flag: false.

```tripol
import tripol.discord

show discord.rgb_color("#FFFFFF")
```

Documented expected output:

```text
16777215
```

#### discord.Embed.title

`func discord.Embed.title(value:string) -> bool`

EN: Set an embed title of up to 256 UTF-8 bytes.

ES: Define un título de hasta 256 bytes UTF-8.

Native example flag: false.

```tripol
import tripol.discord

embed = discord.embed()
show embed.title("Ayuda")
```

Documented expected output:

```text
true
```

#### discord.Embed.description

`func discord.Embed.description(value:string) -> bool`

EN: Set an embed description of up to 512 UTF-8 bytes.

ES: Define una descripción de hasta 512 bytes UTF-8.

Native example flag: false.

```tripol
import tripol.discord

embed = discord.embed()
show embed.description("Escribe !hola")
```

Documented expected output:

```text
true
```

#### discord.Embed.footer

`func discord.Embed.footer(value:string,icon_url="") -> bool`

EN: Set a footer of up to 512 UTF-8 bytes.

ES: Define un pie de hasta 512 bytes UTF-8.

Native example flag: false.

```tripol
import tripol.discord

embed = discord.embed()
show embed.footer("DISCORD REVOLUTION")
```

Documented expected output:

```text
true
```

#### discord.Embed.color

`func discord.Embed.color(value:string) -> bool`

EN: Set the RGB accent as six hex digits preceded by #.

ES: Define el color RGB con seis dígitos hexadecimales precedidos por #.

Native example flag: false.

```tripol
import tripol.discord

embed = discord.embed()
show embed.color("#00A0FF")
```

Documented expected output:

```text
true
```

#### discord.Embed.field

`func discord.Embed.field(label:string, value:string, inline=false) -> bool`

EN: Add one of up to 25 fields. Set inline to true for a compact field.

ES: Añade uno de hasta 25 campos. Usa inline en true para un campo compacto.

Native example flag: false.

Note: Labels accept 1–256 UTF-8 bytes. Values accept 1–512.

Nota: Las etiquetas aceptan 1–256 bytes UTF-8. Los valores aceptan 1–512.

```tripol
import tripol.discord

embed = discord.embed()
show embed.field("Versión", "0.28.0", true)
```

Documented expected output:

```text
true
```

#### discord.Embed.payload

`func discord.Embed.payload() -> string`

EN: Generate Discord create-message JSON with mentions disabled.

ES: Genera el JSON de creación de mensajes de Discord con menciones desactivadas.

Native example flag: false.

Note: Requires text or media. A message can include up to 6,000 text characters across embeds.

Nota: Requiere texto o medios. Un mensaje admite hasta 6.000 caracteres de texto entre todos sus embeds.

```tripol
import tripol.discord

embed = discord.embed()
embed.title("Hola")
show embed.payload()
```

Documented expected output:

```text
{"embeds":[{"title":"Hola"}],"allowed_mentions":{"parse":[]}}
```

#### discord.Embed.send

`func discord.Embed.send(channel_id:string) -> bool`

EN: Send the embed to a Discord channel over REST.

ES: Envía el embed a un canal de Discord por REST.

Native example flag: true.

Note: Replace 123 with a channel ID. Windows or Linux x86-64, DISCORD_TOKEN and Embed Links permission required.

Nota: Sustituye 123 por el ID. Requiere Windows o Linux x86-64, DISCORD_TOKEN y permiso para insertar enlaces.

```tripol
import tripol.discord

embed = discord.embed()
embed.title("Hola")
embed.send("123")
```

#### discord.Embed.url

`func discord.Embed.url(value:string) -> bool`

EN: Set the link associated with the title.

ES: Define el enlace del título.

Native example flag: false.

```tripol
import tripol.discord

embed = discord.embed()
show embed.url("https://example.org/")
```

Documented expected output:

```text
true
```

#### discord.Embed.image

`func discord.Embed.image(value:string) -> bool`

EN: Set the main image.

ES: Define la imagen principal.

Native example flag: false.

```tripol
import tripol.discord

embed = discord.embed()
show embed.image("https://example.org/main.png")
```

Documented expected output:

```text
true
```

#### discord.Embed.thumbnail

`func discord.Embed.thumbnail(value:string) -> bool`

EN: Set the thumbnail image.

ES: Define la miniatura.

Native example flag: false.

```tripol
import tripol.discord

embed = discord.embed()
show embed.thumbnail("https://example.org/thumb.png")
```

Documented expected output:

```text
true
```

#### discord.Interaction.reply_embeds

`func discord.Interaction.reply_embeds(batch:Embeds,ephemeral=false) -> string`

EN: Build an initial reply containing a traditional embed batch.

ES: Construye una respuesta inicial con un lote de embeds tradicionales.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
e=discord.embed()
e.title("Hola")
batch=discord.embeds()
batch.add(e)
show len(ctx.reply_embeds(batch,true))>0
```

Documented expected output:

```text
true
```

#### discord.Embed.text_size

`func discord.Embed.text_size() -> int`

EN: Count characters subject to the combined embed limit.

ES: Cuenta los caracteres sujetos al límite combinado.

Native example flag: false.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
show e.text_size()
```

Documented expected output:

```text
4
```

#### discord.Embed.entry

`func discord.Embed.entry() -> string`

EN: Build a single embed object for a message.

ES: Genera la estructura de un embed para un mensaje.

Native example flag: false.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
show len(e.entry())>0
```

Documented expected output:

```text
true
```

#### discord.Embed.author

`func discord.Embed.author(name:string,url="",icon_url="") -> bool`

EN: Set an author name, link and icon.

ES: Define el nombre, el enlace y el icono del autor.

Native example flag: false.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
show e.author("Ada","https://example.org","https://example.org/icon.png")
```

Documented expected output:

```text
true
```

#### discord.Embed.timestamp

`func discord.Embed.timestamp(value:string) -> bool`

EN: Set an RFC 3339 embed timestamp.

ES: Define la fecha del embed mediante RFC 3339.

Native example flag: false.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
show e.timestamp("2028-02-29T10:15:30Z")
```

Documented expected output:

```text
true
```

#### discord.Embeds.add

`func discord.Embeds.add(embed:Embed) -> bool`

EN: Append an embed and enforce the combined text budget.

ES: Añade un embed y valida el límite de texto combinado.

Native example flag: false.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
batch=discord.embeds()
show batch.add(e)
```

Documented expected output:

```text
true
```

#### discord.Embeds.payload

`func discord.Embeds.payload() -> string`

EN: Build the message payload for multiple traditional embeds.

ES: Genera el JSON de un mensaje con varios embeds.

Native example flag: false.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
batch=discord.embeds()
batch.add(e)
show len(batch.payload())>0
```

Documented expected output:

```text
true
```

#### discord.Embeds.send

`func discord.Embeds.send(channel_id:string) -> bool`

EN: Send a batch of embeds to a channel.

ES: Envía un lote de embeds a un canal.

Native example flag: true.

Note: Requires an authorized Discord bot and a native run on Linux or Windows.

Nota: Requiere un bot autorizado de Discord y ejecución nativa en Linux o Windows.

```tripol
import tripol.discord

e=discord.embed()
e.title("Hola")
batch=discord.embeds()
batch.add(e)
batch.send("123")
```

### Components V2 / Components V2 (16)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-components)

#### discord.components

`func discord.components() -> Components`

EN: Create a Components V2 message builder.

ES: Crea un constructor de mensajes Components V2.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
layout.text("Hola")
show layout.payload()
```

Documented expected output:

```text
{"flags":32768,"components":[{"type":10,"content":"Hola"}],"allowed_mentions":{"parse":[]}}
```

#### discord.container

`func discord.container(color="#FFFFFF") -> Components`

EN: Create a Components V2 panel with an accent color.

ES: Crea un panel Components V2 con color de acento.

Native example flag: false.

```tripol
import tripol.discord

panel=discord.container("#FFFFFF")
show panel.text("Hola")
```

Documented expected output:

```text
true
```

#### discord.Components.append_entry

`func discord.Components.append_entry(entry:string,weight:int) -> bool`

EN: Append a pre-serialized entry with a component count.

ES: Añade un componente serializado y su recuento.

Native example flag: false.

Note: Low-level method. Prefer the typed builders for shape validation.

Nota: Método de bajo nivel. Usa los constructores tipados para validar la estructura.

```tripol
import tripol.discord

layout = discord.components()
show layout.append_entry("{\"type\":10,\"content\":\"Hola\"}",1)
```

Documented expected output:

```text
true
```

#### discord.Components.entries_json

`func discord.Components.entries_json() -> string`

EN: Encode the component array for a message or container.

ES: Serializa el arreglo de componentes para un mensaje o contenedor.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
layout.text("Hola")
show layout.entries_json()
```

Documented expected output:

```text
[{"type":10,"content":"Hola"}]
```

#### discord.Components.text

`func discord.Components.text(content:string) -> bool`

EN: Add a Text Display within the aggregate 4,000-character budget.

ES: Añade texto dentro del límite total de 4.000 caracteres.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
show layout.text("Bienvenido")
```

Documented expected output:

```text
true
```

#### discord.Components.separator

`func discord.Components.separator(spacing=1,divider=true) -> bool`

EN: Add a separator with small or large spacing.

ES: Añade un separador con espacio pequeño o grande.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
show layout.separator(2,false)
```

Documented expected output:

```text
true
```

#### discord.Components.section

`func discord.Components.section(content:string,thumbnail_url:string,description="") -> bool`

EN: Add a section with a text display and thumbnail accessory.

ES: Añade una sección con texto y una miniatura.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
show layout.section("Bienvenido","https://example.org/image.png","Foto")
```

Documented expected output:

```text
true
```

#### discord.Components.gallery

`func discord.Components.gallery(urls:string_list) -> bool`

EN: Add a Media Gallery with 1 to 10 public URLs.

ES: Añade una galería multimedia con 1 a 10 URLs públicas.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
show layout.gallery(string_list(["https://example.org/image.png"]))
```

Documented expected output:

```text
true
```

#### discord.Components.container

`func discord.Components.container(children:Components,accent=-1,spoiler=false) -> bool`

EN: Wrap a component array in a container with an optional RGB accent.

ES: Agrupa componentes en un contenedor con color RGB opcional.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
layout.text("Hola")
message=discord.components()
show message.container(layout,41215)
```

Documented expected output:

```text
true
```

#### discord.Components.payload

`func discord.Components.payload() -> string`

EN: Build a message body with the Components V2 flag and mentions disabled.

ES: Genera el mensaje con la bandera Components V2 y menciones desactivadas.

Native example flag: false.

Note: Encoded text must fit in the TRIPOL string capacity of 65,535 bytes.

Nota: El JSON debe caber en la capacidad de cadena de TRIPOL de 65.535 bytes.

```tripol
import tripol.discord

layout = discord.components()
layout.text("Hola")
show layout.payload()
```

Documented expected output:

```text
{"flags":32768,"components":[{"type":10,"content":"Hola"}],"allowed_mentions":{"parse":[]}}
```

#### discord.Components.send

`func discord.Components.send(channel_id:string) -> bool`

EN: Send a Components V2 message to a channel.

ES: Envía un mensaje Components V2 a un canal.

Native example flag: true.

Note: Run on Windows or Linux x86-64 with DISCORD_TOKEN configured.

Nota: Ejecuta en Windows o Linux x86-64 con DISCORD_TOKEN configurado.

```tripol
import tripol.discord

layout = discord.components()
layout.text("Hola")
layout.send("123")
```

#### discord.Components.append_text

`func discord.Components.append_text(entry:string,weight:int,characters:int) -> bool`

EN: Append a serialized entry and account for its Text Display characters.

ES: Inserta un componente serializado y contabiliza su texto.

Native example flag: false.

```tripol
import tripol.discord

layout=discord.components()
show layout.append_text("{\"type\":10,\"content\":\"Hola\"}",1,4)
```

Documented expected output:

```text
true
```

#### discord.Components.row

`func discord.Components.row(buttons:Buttons) -> bool`

EN: Add a group of buttons and validate its custom IDs.

ES: Añade un grupo de botones y valida sus IDs.

Native example flag: false.

```tripol
import tripol.discord

buttons=discord.button_row()
buttons.button("Aceptar","confirm",1)
layout=discord.components()
show layout.row(buttons)
```

Documented expected output:

```text
true
```

#### discord.Interaction.reply_components

`func discord.Interaction.reply_components(layout:Components,ephemeral=false) -> string`

EN: Build an initial reply containing a V2 layout.

ES: Construye una respuesta inicial con componentes V2.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
layout=discord.components()
layout.text("Hola")
show len(ctx.reply_components(layout,true))>0
```

Documented expected output:

```text
true
```

#### discord.Interaction.update_components

`func discord.Interaction.update_components(layout:Components) -> string`

EN: Update the source message with a V2 layout.

ES: Actualiza el mensaje de origen con componentes V2.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
layout=discord.components()
layout.text("Hola")
show len(ctx.update_components(layout))>0
```

Documented expected output:

```text
true
```

#### discord.Channel.edit_components

`func discord.Channel.edit_components(message_id:string,layout:Components) -> bool`

EN: Replace the layout of a Components V2 message owned by the bot.

ES: Sustituye el panel de un mensaje Components V2 del propio bot.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

room=discord.channel_ref("123")
panel=discord.container()
panel.button("Done","done","success","✅",true)
show room.edit_components("456",panel)
```

### Buttons & emojis / Botones y emojis (11)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-buttons)

#### discord.button_row

`func discord.button_row() -> Buttons`

EN: Create a row of up to five buttons.

ES: Crea una fila de hasta cinco botones.

Native example flag: false.

```tripol
import tripol.discord

buttons=discord.button_row()
buttons.button("Aceptar","confirm",1)
layout=discord.components()
show layout.row(buttons)
```

Documented expected output:

```text
true
```

#### discord.button_style

`func discord.button_style(name:string) -> int`

EN: Look up a custom button style by name.

ES: Consulta un estilo de botón por nombre.

Native example flag: false.

```tripol
import tripol.discord

show discord.button_style("primary")
```

Documented expected output:

```text
1
```

#### discord.emoji_payload

`func discord.emoji_payload(emoji:string) -> string`

EN: Serialize a Unicode, custom or animated Discord emoji.

ES: Serializa un emoji Unicode, personalizado o animado.

Native example flag: false.

```tripol
import tripol.discord

show discord.emoji_payload("✅")
```

Documented expected output:

```text
{"name":"✅"}
```

#### discord.Components.link

`func discord.Components.link(label:string,url:string,emoji="",disabled=false) -> bool`

EN: Add a link button inside an action row.

ES: Añade un botón de enlace dentro de una fila.

Native example flag: false.

```tripol
import tripol.discord

layout = discord.components()
show layout.link("Docs", "https://example.org/docs")
```

Documented expected output:

```text
true
```

#### discord.Components.button

`func discord.Components.button(label:string,custom_id:string,style=1,emoji="",disabled=false) -> bool
func discord.Components.button(label:string,custom_id:string,style:string,emoji="",disabled=false) -> bool`

EN: Add a button with a unique custom ID inside an action row.

ES: Añade un botón con un ID único dentro de una fila.

Native example flag: false.

Note: Register a callback for the same custom ID with bot.on_button.

Nota: Registra un callback para el mismo ID con bot.on_button.

```tripol
import tripol.discord

layout = discord.components()
show layout.button("OK","accept",3)
```

Documented expected output:

```text
true
```

#### discord.Bot.on_button

`func discord.Bot.on_button(id:string,handler:func(Interaction)->string) -> bool`

EN: Register a handler for a custom button ID.

ES: Registra un manejador para el ID de un botón.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
show b.on_button("confirm",&answer)
```

Documented expected output:

```text
true
```

#### discord.Buttons.button

`func discord.Buttons.button(label:string,id:string,style=1,disabled=false,emoji="") -> bool
func discord.Buttons.button(label:string,id:string,style:string,disabled=false,emoji="") -> bool`

EN: Add a custom button to a multi-button row.

ES: Añade un botón personalizado a una fila.

Native example flag: false.

```tripol
import tripol.discord

buttons=discord.button_row()
show buttons.button("Aceptar","confirm",3)
```

Documented expected output:

```text
true
```

#### discord.Buttons.link

`func discord.Buttons.link(label:string,url:string,disabled=false,emoji="") -> bool`

EN: Add a link button to a multi-button row.

ES: Añade un botón de enlace a una fila.

Native example flag: false.

```tripol
import tripol.discord

buttons=discord.button_row()
show buttons.link("Docs","https://example.org/docs")
```

Documented expected output:

```text
true
```

#### discord.Components.premium

`func discord.Components.premium(sku_id:string,disabled=false) -> bool`

EN: Add a purchase button for an application SKU.

ES: Añade un botón de compra para un SKU de la aplicación.

Native example flag: false.

```tripol
import tripol.discord

panel=discord.container()
show panel.premium("123456")
```

Documented expected output:

```text
true
```

#### discord.Buttons.premium

`func discord.Buttons.premium(sku_id:string,disabled=false) -> bool`

EN: Add an SKU button to a row without a custom ID.

ES: Añade a una fila un botón de SKU sin ID personalizado.

Native example flag: false.

```tripol
import tripol.discord

row=discord.button_row()
show row.premium("123456")
```

Documented expected output:

```text
true
```

#### discord.Buttons.disable

`func discord.Buttons.disable(id:string,disabled=true) -> bool`

EN: Disable or enable a button in a row by its custom ID.

ES: Deshabilita o habilita un botón de la fila por su ID.

Native example flag: false.

```tripol
import tripol.discord

row=discord.button_row()
row.button("Done","done","success",false,"✅")
show row.disable("done")
```

Documented expected output:

```text
true
```

### Select menus / Selectores (11)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-selects)

#### discord.select

`func discord.select(id:string,kind="string",minimum=1,maximum=1) -> Select`

EN: Create a String, User, Role, Mentionable or Channel Select.

ES: Crea un selector de texto, usuario, rol, mencionable o canal.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("choice","string",0,1)
menu.option("Azul","b")
show len(menu.payload())>0
```

Documented expected output:

```text
true
```

#### discord.interaction_values

`func discord.interaction_values(values:string) -> string_list`

EN: Decode a list of selected string values.

ES: Decodifica una lista de valores de texto seleccionados.

Native example flag: false.

```tripol
import tripol.discord

values=discord.interaction_values("[\"Ada\"]")
show values[0]
```

Documented expected output:

```text
Ada
```

#### discord.Components.select

`func discord.Components.select(menu:Select) -> bool`

EN: Add a select menu as one action row.

ES: Añade un selector en una fila.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("choice","string",0,1)
menu.option("Azul","b")
layout=discord.components()
show layout.select(menu)
```

Documented expected output:

```text
true
```

#### discord.Bot.on_select

`func discord.Bot.on_select(id:string,handler:func(Interaction)->string) -> bool`

EN: Register a handler for a custom select ID.

ES: Registra un manejador para el ID de un selector.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
show b.on_select("choice",&answer)
```

Documented expected output:

```text
true
```

#### discord.Interaction.values

`func discord.Interaction.values() -> string_list`

EN: Return submitted select values or entity IDs.

ES: Devuelve los valores o IDs seleccionados.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
values=ctx.values()
show values[0],values[1]
```

Documented expected output:

```text
42 43
```

#### discord.Select.placeholder

`func discord.Select.placeholder(value:string) -> bool`

EN: Set the empty-state text of a select menu.

ES: Define el texto inicial de un selector.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("choice","string",0,1)
show menu.placeholder("Elige un color")
```

Documented expected output:

```text
true
```

#### discord.Select.option

`func discord.Select.option(label:string,value:string,description="",selected=false,emoji="") -> bool`

EN: Add a label and value to a String Select.

ES: Añade una etiqueta y un valor a un selector de texto.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("choice","string",0,1)
show menu.option("Azul","b","Un color")
```

Documented expected output:

```text
true
```

#### discord.Select.default_id

`func discord.Select.default_id(id:string,kind:string) -> bool`

EN: Add a default entity to an auto-populated select.

ES: Añade una entidad predeterminada a un selector automático.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("people","user")
show menu.default_id("123","user")
```

Documented expected output:

```text
true
```

#### discord.Select.channels

`func discord.Select.channels(kinds:ints) -> bool`

EN: Limit the channel types shown in a Channel Select.

ES: Limita los tipos de canal de un selector de canales.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("room","channel")
allowed=ints(1)
allowed[0]=2
show menu.channels(allowed)
```

Documented expected output:

```text
true
```

#### discord.Select.payload

`func discord.Select.payload() -> string`

EN: Build the complete JSON for a select component.

ES: Genera el JSON completo de un selector.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("choice","string",0,1)
menu.option("Azul","b")
show len(menu.payload())>0
```

Documented expected output:

```text
true
```

#### discord.Select.disable

`func discord.Select.disable(disabled=true) -> bool`

EN: Disable or enable a select menu in a message.

ES: Deshabilita o habilita un selector en un mensaje.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("status")
menu.option("Ready","ready","",false,"✅")
show menu.disable()
```

Documented expected output:

```text
true
```

### Forms & modals / Formularios (21)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-forms)

#### discord.modal_field_value

`func discord.modal_field_value(node:string,id:string,property:string,depth:int) -> string`

EN: Find a property in a submitted modal component by custom ID.

ES: Busca una propiedad de un campo entregado por su ID.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(5,"hola","123","456","789","123","567","{\"components\":[{\"type\":18,\"component\":{\"custom_id\":\"topic\",\"value\":\"Hola\"}},{\"type\":18,\"component\":{\"custom_id\":\"file\",\"values\":[\"42\"]}}]}",false)
show discord.modal_field_value(ctx.data_json,"topic","value",0)
```

Documented expected output:

```text
"Hola"
```

#### discord.modal

`func discord.modal(id:string,title:string) -> Modal`

EN: Create a modal with a custom ID and a title.

ES: Crea un formulario con ID y título.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("sample","Cuéntanos")
form.text("topic","Tema")
show len(form.payload())>0
```

Documented expected output:

```text
true
```

#### discord.file_types_payload

`func discord.file_types_payload(types:string_list) -> string`

EN: Serialize attachment extension filters.

ES: Serializa filtros por extensión de adjuntos.

Native example flag: false.

```tripol
import tripol.discord

show discord.file_types_payload(list(["image",".pdf"]))
```

Documented expected output:

```text
["image",".pdf"]
```

#### discord.Bot.on_modal

`func discord.Bot.on_modal(id:string,handler:func(Interaction)->string) -> bool`

EN: Register a handler for a submitted modal ID.

ES: Registra un manejador para el ID de un formulario enviado.

Native example flag: false.

```tripol
import tripol.discord

func answer(ctx:discord.Interaction)->string {return ctx.reply("Hola, "+ctx.option_text("text"),true)}
b=discord.bot()
show b.on_modal("sample",&answer)
```

Documented expected output:

```text
true
```

#### discord.Interaction.field

`func discord.Interaction.field(id:string) -> string`

EN: Read the decoded text of a submitted input.

ES: Lee el texto decodificado de un campo enviado.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(5,"hola","123","456","789","123","567","{\"components\":[{\"type\":18,\"component\":{\"custom_id\":\"topic\",\"value\":\"Hola\"}},{\"type\":18,\"component\":{\"custom_id\":\"file\",\"values\":[\"42\"]}}]}",false)
show ctx.field("topic")
```

Documented expected output:

```text
Hola
```

#### discord.Interaction.field_data

`func discord.Interaction.field_data(id:string) -> string`

EN: Read the submitted value array of a select or file upload.

ES: Lee el arreglo enviado por un selector o un campo de archivos.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(5,"hola","123","456","789","123","567","{\"components\":[{\"type\":18,\"component\":{\"custom_id\":\"topic\",\"value\":\"Hola\"}},{\"type\":18,\"component\":{\"custom_id\":\"file\",\"values\":[\"42\"]}}]}",false)
show ctx.field_data("file")
```

Documented expected output:

```text
["42"]
```

#### discord.Interaction.open_modal

`func discord.Interaction.open_modal(form:Modal) -> string`

EN: Build the response that displays a modal.

ES: Construye la respuesta que abre un formulario.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
form=discord.modal("sample","Cuéntanos")
form.text("topic","Tema")
show len(ctx.open_modal(form))>0
```

Documented expected output:

```text
true
```

#### discord.Modal.add_label

`func discord.Modal.add_label(id:string,label:string,component:string,description="") -> bool`

EN: Wrap a serialized form field in a Label component.

ES: Envuelve un campo serializado en un componente Label.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("sample","Cuéntanos")
show form.add_label("topic","Tema","{\"type\":4,\"custom_id\":\"topic\",\"style\":1}")
```

Documented expected output:

```text
true
```

#### discord.Modal.text

`func discord.Modal.text(id:string,label:string,style=1,required=true,min_length=0,max_length=4000,placeholder="",value="",description="") -> bool`

EN: Add a single-line or paragraph text input.

ES: Añade un campo de texto de una línea o de párrafo.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("sample","Cuéntanos")
show form.text("topic","Tema",2,true,2,4000)
```

Documented expected output:

```text
true
```

#### discord.Modal.payload

`func discord.Modal.payload() -> string`

EN: Build the modal JSON payload.

ES: Genera el JSON del formulario.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("sample","Cuéntanos")
form.text("topic","Tema")
show len(form.payload())>0
```

Documented expected output:

```text
true
```

#### discord.Modal.select

`func discord.Modal.select(label:string,menu:Select,required=true,description="") -> bool`

EN: Add a select menu to a form label.

ES: Añade un selector a una etiqueta de formulario.

Native example flag: false.

```tripol
import tripol.discord

menu=discord.select("choice","string",0,1)
menu.option("Azul","b")
form=discord.modal("sample","Cuéntanos")
show form.select("Color",menu,false)
```

Documented expected output:

```text
true
```

#### discord.Modal.help_text

`func discord.Modal.help_text(content:string) -> bool`

EN: Add explanatory Text Display content to a form.

ES: Añade texto de ayuda a un formulario.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("sample","Cuéntanos")
show form.help_text("Cuéntanos más")
```

Documented expected output:

```text
true
```

#### discord.Modal.file

`func discord.Modal.file(id:string,label:string,minimum=1,maximum=1,required=true,description="") -> bool
func discord.Modal.file(id:string,label:string,types:string_list,minimum=1,maximum=1,required=true,description="") -> bool`

EN: Add an input for files submitted by the user.

ES: Añade un campo para archivos enviados por el usuario.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("sample","Cuéntanos")
show form.file("file","Adjunto",0,2,false)
```

Documented expected output:

```text
true
```

#### discord.Modal.choice_group

`func discord.Modal.choice_group(label:string,menu:Select,radio:bool,required:bool,description="") -> bool`

EN: Build a radio or checkbox group from manually specified options.

ES: Construye un grupo de opción única o de casillas con opciones definidas.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("colors","Colores")
menu=discord.select("color")
menu.option("Azul","blue")
menu.option("Rojo","red")
show form.choice_group("Color",menu,true,true)
```

Documented expected output:

```text
true
```

#### discord.Modal.radio

`func discord.Modal.radio(label:string,menu:Select,required=true,description="") -> bool`

EN: Add a Radio Group with one selected option.

ES: Añade un grupo de opción única.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("colors","Colores")
menu=discord.select("color")
menu.option("Azul","blue")
menu.option("Rojo","red")
show form.radio("Color",menu)
```

Documented expected output:

```text
true
```

#### discord.Modal.checkbox_group

`func discord.Modal.checkbox_group(label:string,menu:Select,required=true,description="") -> bool`

EN: Add a Checkbox Group with multiple options.

ES: Añade un grupo de casillas con varias opciones.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("colors","Colores")
menu=discord.select("color")
menu.option("Azul","blue")
menu.option("Rojo","red")
show form.checkbox_group("Colores",menu)
```

Documented expected output:

```text
true
```

#### discord.Modal.checkbox

`func discord.Modal.checkbox(id:string,label:string,checked=false,description="") -> bool`

EN: Add a single checkbox with an optional default.

ES: Añade una casilla individual con un valor predeterminado opcional.

Native example flag: false.

```tripol
import tripol.discord

form=discord.modal("colors","Colores")
show form.checkbox("agree","Acepto",true)
```

Documented expected output:

```text
true
```

#### discord.Interaction.checked

`func discord.Interaction.checked(id:string) -> bool`

EN: Read the submitted value of a checkbox.

ES: Lee el valor enviado de una casilla.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(5,"form","","","","","","{\"components\":[{\"type\":18,\"component\":{\"type\":23,\"custom_id\":\"agree\",\"value\":true}}]}",false)
show ctx.checked("agree")
```

Documented expected output:

```text
true
```

#### discord.Interaction.has_field

`func discord.Interaction.has_field(id:string) -> bool`

EN: Check whether a field was submitted, including when empty.

ES: Comprueba si se envió un campo, incluso vacío.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":5,\"data\":{\"components\":[{\"type\":18,\"component\":{\"type\":4,\"custom_id\":\"name\",\"value\":\"\"}},{\"type\":18,\"component\":{\"type\":6,\"custom_id\":\"roles\",\"values\":[\"100\",\"200\"]}},{\"type\":18,\"component\":{\"type\":19,\"custom_id\":\"file\",\"values\":[\"42\"]}}],\"resolved\":{\"attachments\":{\"42\":{\"id\":\"42\",\"filename\":\"report.pdf\",\"url\":\"https://example.com/report.pdf\",\"size\":42,\"content_type\":\"application/pdf\"}}}}}","mFeedback")
show ctx.has_field("name"),ctx.has_field("missing")
```

Documented expected output:

```text
true false
```

#### discord.Interaction.field_values

`func discord.Interaction.field_values(id:string) -> string_list`

EN: Read the selected values as a string list.

ES: Lee los valores seleccionados como una lista de cadenas.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":5,\"data\":{\"components\":[{\"type\":18,\"component\":{\"type\":4,\"custom_id\":\"name\",\"value\":\"\"}},{\"type\":18,\"component\":{\"type\":6,\"custom_id\":\"roles\",\"values\":[\"100\",\"200\"]}},{\"type\":18,\"component\":{\"type\":19,\"custom_id\":\"file\",\"values\":[\"42\"]}}],\"resolved\":{\"attachments\":{\"42\":{\"id\":\"42\",\"filename\":\"report.pdf\",\"url\":\"https://example.com/report.pdf\",\"size\":42,\"content_type\":\"application/pdf\"}}}}}","mFeedback")
roles=ctx.field_values("roles")
show roles[0],roles[1]
```

Documented expected output:

```text
100 200
```

#### discord.Interaction.field_attachment

`func discord.Interaction.field_attachment(id:string,index=0) -> UploadedFile`

EN: Read the metadata of an uploaded file by field and index.

ES: Lee los metadatos de un archivo subido por campo e índice.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":5,\"data\":{\"components\":[{\"type\":18,\"component\":{\"type\":4,\"custom_id\":\"name\",\"value\":\"\"}},{\"type\":18,\"component\":{\"type\":6,\"custom_id\":\"roles\",\"values\":[\"100\",\"200\"]}},{\"type\":18,\"component\":{\"type\":19,\"custom_id\":\"file\",\"values\":[\"42\"]}}],\"resolved\":{\"attachments\":{\"42\":{\"id\":\"42\",\"filename\":\"report.pdf\",\"url\":\"https://example.com/report.pdf\",\"size\":42,\"content_type\":\"application/pdf\"}}}}}","mFeedback")
file=ctx.field_attachment("file")
show file.filename,file.size
```

Documented expected output:

```text
report.pdf 42
```

### Files & attachments / Archivos y adjuntos (10)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-files)

#### discord.attachment_name_valid

`func discord.attachment_name_valid(filename:string) -> bool`

EN: Check the filename used by a Discord attachment reference.

ES: Valida el nombre utilizado en una referencia a un adjunto de Discord.

Native example flag: false.

```tripol
import tripol.discord

show discord.attachment_name_valid("report.txt")
```

Documented expected output:

```text
true
```

#### discord.attachment_url

`func discord.attachment_url(filename:string) -> string`

EN: Build a reference to a file uploaded in the same message.

ES: Crea una referencia a un archivo adjunto al mismo mensaje.

Native example flag: false.

```tripol
import tripol.discord

show discord.attachment_url("report.txt")
```

Documented expected output:

```text
attachment://report.txt
```

#### discord.attachment_media_url

`func discord.attachment_media_url(value:string) -> bool`

EN: Accept an HTTP URL or a valid message attachment reference.

ES: Valida una URL HTTP o una referencia a un archivo adjunto.

Native example flag: false.

```tripol
import tripol.discord

show discord.attachment_media_url("attachment://report.txt")
```

Documented expected output:

```text
true
```

#### discord.send_attachments

`func discord.send_attachments(channel_id:string,payload:string,paths:string_list) -> string`

EN: Stream files and a Discord message payload to a channel.

ES: Envía archivos y el payload de un mensaje a un canal.

Native example flag: true.

Target: "native".

Note: Multipart uploads use the native Windows or Linux launcher.

Nota: Los envíos multipart usan el lanzador nativo de Windows o Linux.

```tripol
import tripol.discord

show discord.send_attachments("123",discord.interaction_message("Informe"),list(["report.txt"]))
```

#### discord.upload_envelope

`func discord.upload_envelope(payload:string,paths:string_list) -> string`

EN: Serialize a payload and a list of paths for the launcher.

ES: Serializa un payload y una lista de rutas para el lanzador.

Native example flag: false.

```tripol
import tripol.discord

show discord.upload_envelope("{}",list(["report.txt"]))
```

Documented expected output:

```text
{"payload":{},"files":["report.txt"]}
```

#### discord.Components.send_files

`func discord.Components.send_files(channel_id:string,paths:string_list) -> string`

EN: Publish a V2 layout with files and return the new message ID.

ES: Publica un panel V2 con archivos y devuelve el ID del mensaje.

Native example flag: true.

Target: "native".

Note: Multipart uploads use the native Windows or Linux launcher.

Nota: Los envíos multipart usan el lanzador nativo de Windows o Linux.

```tripol
import tripol.discord

panel=discord.container()
panel.file("report.txt")
show panel.send_files("123",list(["report.txt"]))
```

#### discord.Embed.send_files

`func discord.Embed.send_files(channel_id:string,paths:string_list) -> string`

EN: Publish an embed with its referenced files.

ES: Publica un embed con los archivos a los que hace referencia.

Native example flag: true.

Target: "native".

Note: Multipart uploads use the native Windows or Linux launcher.

Nota: Los envíos multipart usan el lanzador nativo de Windows o Linux.

```tripol
import tripol.discord

embed=discord.embed()
embed.image(discord.attachment_url("photo.png"))
show embed.send_files("123",list(["photo.png"]))
```

#### discord.Embeds.send_files

`func discord.Embeds.send_files(channel_id:string,paths:string_list) -> string`

EN: Publish an embed batch and attached files.

ES: Publica un lote de embeds y sus archivos adjuntos.

Native example flag: true.

Target: "native".

Note: Multipart uploads use the native Windows or Linux launcher.

Nota: Los envíos multipart usan el lanzador nativo de Windows o Linux.

```tripol
import tripol.discord

embed=discord.embed()
embed.image(discord.attachment_url("photo.png"))
batch=discord.embeds()
batch.add(embed)
show batch.send_files("123",list(["photo.png"]))
```

#### discord.Components.file

`func discord.Components.file(filename:string,spoiler=false) -> bool`

EN: Display a file attached to the same V2 message.

ES: Muestra un archivo adjunto al mismo mensaje V2.

Native example flag: false.

```tripol
import tripol.discord

panel=discord.container()
show panel.file("report.txt",true)
```

Documented expected output:

```text
true
```

#### discord.Interaction.reply_files

`func discord.Interaction.reply_files(content:string,paths:string_list,ephemeral=false) -> bool
func discord.Interaction.reply_files(embed:Embed,paths:string_list,ephemeral=false) -> bool
func discord.Interaction.reply_files(batch:Embeds,paths:string_list,ephemeral=false) -> bool
func discord.Interaction.reply_files(layout:Components,paths:string_list,ephemeral=false) -> bool`

EN: Send the initial interaction response immediately with local files.

ES: Envía de inmediato la respuesta inicial a la interacción con archivos locales.

Native example flag: true.

Target: "native".

Note: Run through the Windows or Linux launcher with an authorized Discord application.

Nota: Ejecuta mediante el lanzador de Windows o Linux con una aplicación de Discord autorizada.

```tripol
import tripol.discord

func report(ctx) {
 panel=discord.components()
 panel.file("report.txt")
 show ctx.reply_files(panel,list(["report.txt"]),true)
 ""
}
bot=discord.bot()
bot.slash("report","Send a test file",report)
bot.start()
```

### Moderation & permissions / Moderación y permisos (25)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-moderation)

#### discord.permissions_bitset_valid

`func discord.permissions_bitset_valid(value:string) -> bool`

EN: Check the canonical decimal format of a permission bitset.

ES: Comprueba el formato decimal canónico de una máscara de permisos.

Native example flag: false.

```tripol
import tripol.discord

show discord.permissions_bitset_valid("8192")
```

Documented expected output:

```text
true
```

#### discord.command_permission_payload

`func discord.command_permission_payload(default_permissions:string) -> string`

EN: Serialize default application command permissions.

ES: Serializa los permisos predeterminados de un comando de aplicación.

Native example flag: false.

```tripol
import tripol.discord

show discord.command_permission_payload("8192")
```

Documented expected output:

```text
,"default_member_permissions":"8192"
```

#### discord.permission_index

`func discord.permission_index(name:string) -> int`

EN: Find the bit index of a named Discord permission.

ES: Obtén el índice del bit de un permiso de Discord.

Native example flag: false.

```tripol
import tripol.discord

show discord.permission_index("manage_messages")
```

Documented expected output:

```text
13
```

#### discord.Member.kick

`func discord.Member.kick() -> bool`

EN: Remove the member from the guild.

ES: Expulsa al miembro del servidor.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.kick()
```

#### discord.Member.ban

`func discord.Member.ban(delete_message_seconds=0) -> bool`

EN: Ban the member with optional deletion of recent messages.

ES: Banea al miembro con borrado opcional de mensajes recientes.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.ban(3600)
```

#### discord.Member.unban

`func discord.Member.unban() -> bool`

EN: Revoke a ban for the stored user ID.

ES: Revoca un baneo para el ID de usuario asociado.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.unban()
```

#### discord.Member.nickname

`func discord.Member.nickname(value:string) -> bool`

EN: Set or clear the member’s nickname.

ES: Cambia o limpia el apodo del miembro.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.nickname("Ada")
```

#### discord.Member.timeout

`func discord.Member.timeout(until="") -> bool`

EN: Set an ISO 8601 timeout deadline or clear the existing timeout.

ES: Establece una fecha ISO 8601 para el aislamiento o retíralo.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.timeout()
```

#### discord.Member.role

`func discord.Member.role(role_id:string,enabled=true) -> bool`

EN: Add or remove an existing server role for the member.

ES: Añade o retira al miembro un rol existente en el servidor.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.role("789",false)
```

#### discord.Channel.slowmode

`func discord.Channel.slowmode(seconds:int) -> bool`

EN: Set the channel delay in seconds.

ES: Fija la espera entre mensajes en segundos.

Native example flag: true.

Note: Requires a native run and an authorized bot in the target channel.

Nota: Requiere ejecución nativa y un bot autorizado en el canal.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.slowmode(15)
```

#### discord.Channel.pin

`func discord.Channel.pin(message_id:string,pinned=true) -> bool`

EN: Pin or unpin a message in the channel.

ES: Fija o desfija un mensaje del canal.

Native example flag: true.

Note: Requires a native run and an authorized bot in the target channel.

Nota: Requiere ejecución nativa y un bot autorizado en el canal.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.pin("456",true)
```

#### discord.Channel.delete_message

`func discord.Channel.delete_message(message_id:string) -> bool`

EN: Delete a single message by ID.

ES: Borra un mensaje por ID.

Native example flag: true.

Note: Requires a native run and an authorized bot in the target channel.

Nota: Requiere ejecución nativa y un bot autorizado en el canal.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.delete_message("456")
```

#### discord.Channel.purge

`func discord.Channel.purge(message_ids:string_list) -> bool`

EN: Bulk-delete 2–100 distinct message IDs.

ES: Borra en bloque entre 2 y 100 IDs distintos.

Native example flag: true.

Note: Requires a native run and an authorized bot in the target channel.

Nota: Requiere ejecución nativa y un bot autorizado en el canal.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.purge(list(["456","789"]))
```

#### discord.Channel.react

`func discord.Channel.react(message_id:string,emoji:string,enabled=true) -> bool`

EN: Add or remove the bot’s own reaction on a message.

ES: Añade o retira la reacción del bot a un mensaje.

Native example flag: true.

Note: Requires a native run and an authorized bot in the target channel.

Nota: Requiere ejecución nativa y un bot autorizado en el canal.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.react("456","🔥")
```

#### discord.Role.permissions

`func discord.Role.permissions() -> string`

EN: Fetch the role and @everyone combined guild permissions.

ES: Consulta los permisos combinados del rol y @everyone en el servidor.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
role=server.role_ref("42")
show role.permissions()
```

#### discord.Interaction.can

`func discord.Interaction.can(permission:string) -> bool`

EN: Check the effective member permissions supplied by Discord.

ES: Comprueba los permisos efectivos del miembro informados por Discord.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"guild_id\":\"123\",\"member\":{\"permissions\":\"8192\"},\"data\":{\"type\":2,\"target_id\":\"42\",\"resolved\":{\"users\":{\"42\":{\"id\":\"42\",\"username\":\"Ada\"}}}}}","uProfile")
show ctx.can("manage_messages"),ctx.can("kick_members")
```

Documented expected output:

```text
true false
```

#### discord.Bot.permission

`func discord.Bot.permission(name:string,permission:string,target_option="",guild_id="") -> bool`

EN: Require effective permissions and check a moderation target when applicable.

ES: Exige permisos efectivos y comprueba el objetivo cuando corresponde.

Native example flag: false.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {ctx.reply("Done",true)})
show bot.permission("work","MANAGE_MESSAGES")
```

Documented expected output:

```text
true
```

#### discord.Interaction.bot_can

`func discord.Interaction.bot_can(permission:string) -> bool`

EN: Check the bot’s effective permission in the interaction’s channel.

ES: Comprueba el permiso efectivo del bot en el canal de la interacción.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"guild_id\":\"321\",\"app_permissions\":\"8192\"}","swork/")
show ctx.bot_can("MANAGE_MESSAGES")
```

Documented expected output:

```text
true
```

#### discord.Interaction.can_moderate

`func discord.Interaction.can_moderate(target:string,permission:string) -> bool`

EN: Check effective permissions, the owner and member role hierarchy against a target.

ES: Comprueba permisos efectivos, propietario y jerarquía frente a un objetivo.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("check","Verify role order",func(ctx) {ctx.reply(text_status(ctx.can_moderate("456","BAN_MEMBERS")),true)})
func text_status(allowed:bool)->string {if allowed {return "Allowed"};return "Denied"}
bot.start()
```

#### discord.Channel.lock

`func discord.Channel.lock() -> bool`

EN: Save the @everyone overwrite and deny Send Messages.

ES: Guarda el permiso de @everyone y deniega Enviar mensajes.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.lock()
show room.unlock()
```

#### discord.Channel.unlock

`func discord.Channel.unlock() -> bool`

EN: Restore the saved overwrite in the same launcher.

ES: Restaura el permiso guardado en el mismo lanzador.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.lock()
show room.unlock()
```

#### discord.Channel.overwrite

`func discord.Channel.overwrite(id:string,allow="0",deny="0",kind="role") -> bool`

EN: Set allow and deny permission masks for a user or role.

ES: Establece las máscaras de permisos permitidos y denegados para un usuario o rol.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.overwrite("456","1024","0","user")
```

#### discord.Channel.clear_overwrite

`func discord.Channel.clear_overwrite(id:string) -> bool`

EN: Remove a user or role permission overwrite.

ES: Elimina el permiso explícito de un usuario o rol.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.clear_overwrite("456")
```

#### discord.Bot.restrict

`func discord.Bot.restrict(name:string,scope:string,id:string,guild_id="") -> bool`

EN: Allow specific users, roles, channels or guilds to invoke a registered slash command.

ES: Permite que ciertos usuarios, roles, canales o servidores invoquen un comando slash registrado.

Native example flag: false.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {ctx.reply("Done")})
show bot.restrict("work","channel","456")
```

Documented expected output:

```text
true
```

#### discord.Interaction.has_role

`func discord.Interaction.has_role(id:string) -> bool`

EN: Check whether the invoking member has a specified role.

ES: Comprueba si el miembro que invoca tiene un rol determinado.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.interaction_context("{\"type\":2,\"guild_id\":\"321\",\"member\":{\"roles\":[\"456\"]}}","swork/")
show ctx.has_role("456")
```

Documented expected output:

```text
true
```

### Cooldowns & economy / Cooldowns y economía (7)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-economy)

#### discord.cooldown_bucket

`func discord.cooldown_bucket(command:string,scope:string,subject="",guild_id="") -> CooldownBucket`

EN: Select a cooldown record by command, scope and subject.

ES: Selecciona un cooldown por comando, ámbito y sujeto.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

bucket=discord.cooldown_bucket("work","user","456")
show bucket.remaining()
```

#### discord.Bot.cooldown

`func discord.Bot.cooldown(name:string,seconds:int,scope="user",message="",guild_id="") -> bool`

EN: Require a minimum interval between executions of a registered slash command.

ES: Exige un intervalo mínimo entre ejecuciones de un comando slash registrado.

Native example flag: false.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {ctx.reply("Done",true)})
show bot.cooldown("work",60,"user","Try again in")
```

Documented expected output:

```text
true
```

#### discord.Bot.persist_cooldowns

`func discord.Bot.persist_cooldowns(path:string) -> bool`

EN: Save active cooldown reservations across restarts.

ES: Guarda las reservas de cooldown activas al reiniciar.

Native example flag: false.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {ctx.reply("Done",true)})
show bot.persist_cooldowns("cooldowns.json")
```

Documented expected output:

```text
true
```

#### discord.Interaction.cooldown_remaining

`func discord.Interaction.cooldown_remaining() -> int`

EN: Return the remaining time for the reserved bucket in seconds.

ES: Devuelve los segundos restantes de la reserva actual.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {ctx.reply(text(ctx.cooldown_remaining()))})
bot.cooldown("work",60)
bot.start()
```

#### discord.Interaction.reset_cooldown

`func discord.Interaction.reset_cooldown() -> bool`

EN: Reset the current interaction’s cooldown bucket.

ES: Restablece el cooldown de la interacción actual.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {show ctx.reset_cooldown();ctx.reply("Done")})
bot.cooldown("work",60)
bot.start()
```

#### discord.CooldownBucket.remaining

`func discord.CooldownBucket.remaining() -> int`

EN: Query the seconds remaining for a specific cooldown bucket.

ES: Consulta los segundos restantes para un cooldown específico.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

bucket=discord.cooldown_bucket("work","user","456")
show bucket.remaining()
```

#### discord.CooldownBucket.reset

`func discord.CooldownBucket.reset() -> bool`

EN: Reset a cooldown for a specific scope and subject.

ES: Restablece un cooldown para un ámbito y sujeto específicos.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

bucket=discord.cooldown_bucket("work","user","456")
show bucket.reset()
```

### Servers, channels & roles / Servidores, canales y roles (27)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-servers)

#### discord.guild

`func discord.guild(id:string) -> Guild`

EN: Build a typed server reference from its ID.

ES: Crea una referencia tipada a un servidor por su ID.

Native example flag: false.

```tripol
import tripol.discord

server=discord.guild("123")
show server.id
```

Documented expected output:

```text
123
```

#### discord.channel_ref

`func discord.channel_ref(id:string) -> Channel`

EN: Create a typed reference to a Discord channel by ID.

ES: Crea una referencia tipada a un canal de Discord por ID.

Native example flag: false.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.id()
```

Documented expected output:

```text
123
```

#### discord.role_payload

`func discord.role_payload(name:string,permissions:string,color:string) -> string`

EN: Build a role update with only the supplied fields.

ES: Genera una actualización de rol con los campos indicados.

Native example flag: false.

```tripol
import tripol.discord

show discord.role_payload("Reviewer","0","#FFFFFF")
```

Documented expected output:

```text
{"name":"Reviewer","permissions":"0","colors":{"primary_color":16777215}}
```

#### discord.Guild.member

`func discord.Guild.member(id:string) -> Member`

EN: Build a typed member reference from a user ID within the server.

ES: Crea una referencia tipada a un miembro a partir de su ID de usuario.

Native example flag: false.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.id()
```

Documented expected output:

```text
456
```

#### discord.Member.id

`func discord.Member.id() -> string`

EN: Return a member’s user ID as text.

ES: Devuelve como texto el ID de usuario del miembro.

Native example flag: false.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.id()
```

Documented expected output:

```text
456
```

#### discord.Member.read

`func discord.Member.read() -> string`

EN: Fetch the member’s Discord JSON data.

ES: Consulta los datos JSON del miembro en Discord.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.read()
```

#### discord.Member.dm

`func discord.Member.dm(content:string) -> string`

EN: Send a direct message and return its message ID.

ES: Envía un mensaje directo y devuelve su ID.

Native example flag: true.

Note: Run with an authorized Discord application on Linux or Windows.

Nota: Ejecuta con una aplicación de Discord autorizada en Linux o Windows.

```tripol
import tripol.discord

server=discord.guild("123")
member=server.member("456")
show member.dm("Hola")
```

#### discord.Channel.id

`func discord.Channel.id() -> string`

EN: Return the channel ID as text.

ES: Devuelve el ID del canal como texto.

Native example flag: false.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.id()
```

Documented expected output:

```text
123
```

#### discord.Channel.read

`func discord.Channel.read() -> string`

EN: Fetch the channel JSON from Discord.

ES: Consulta el JSON del canal en Discord.

Native example flag: true.

Note: Requires a native run and an authorized bot in the target channel.

Nota: Requiere ejecución nativa y un bot autorizado en el canal.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.read()
```

#### discord.Guild.channels

`func discord.Guild.channels() -> string`

EN: Fetch the guild’s channel list as JSON.

ES: Consulta la lista de canales del servidor como JSON.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
show server.channels()
```

#### discord.Guild.create_channel

`func discord.Guild.create_channel(name:string,kind="text",parent_id="",topic="") -> string`

EN: Create a text, voice, category or announcement channel.

ES: Crea un canal de texto, voz, categoría o anuncios.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
show server.create_channel("help","text","456","Support")
```

#### discord.Guild.audit

`func discord.Guild.audit(limit=25,before="",action=0) -> string`

EN: Fetch a page of audit logs with optional filters.

ES: Consulta una página del registro de auditoría con filtros opcionales.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
show server.audit(25,"",10)
```

#### discord.Guild.roles

`func discord.Guild.roles() -> string`

EN: Fetch the server’s roles as JSON.

ES: Consulta los roles del servidor como JSON.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
show server.roles()
```

#### discord.Guild.create_role

`func discord.Guild.create_role(name:string,permissions="0",color="") -> string`

EN: Create a role with decimal permissions and an optional hex color.

ES: Crea un rol con permisos decimales y un color HEX opcional.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
show server.create_role("Reviewer","0","#FFFFFF")
```

#### discord.Guild.role_ref

`func discord.Guild.role_ref(id:string) -> Role`

EN: Create a typed role reference within a server.

ES: Crea una referencia tipada a un rol de un servidor.

Native example flag: false.

```tripol
import tripol.discord

server=discord.guild("321")
role=server.role_ref("42")
show role.id()
```

Documented expected output:

```text
42
```

#### discord.Role.id

`func discord.Role.id() -> string`

EN: Return the role ID as text.

ES: Devuelve el ID del rol como texto.

Native example flag: false.

```tripol
import tripol.discord

server=discord.guild("321")
role=server.role_ref("42")
show role.id()
```

Documented expected output:

```text
42
```

#### discord.Role.read

`func discord.Role.read() -> string`

EN: Fetch the Discord JSON for the role.

ES: Consulta los datos JSON del rol en Discord.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
role=server.role_ref("42")
show role.read()
```

#### discord.Role.edit

`func discord.Role.edit(name="",permissions="",color="") -> bool`

EN: Update the supplied role name, permissions or primary color.

ES: Actualiza el nombre, los permisos o el color principal indicados.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
role=server.role_ref("42")
show role.edit("Review","8192")
```

#### discord.Role.delete

`func discord.Role.delete() -> bool`

EN: Delete the selected role.

ES: Elimina el rol seleccionado.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

server=discord.guild("321")
role=server.role_ref("42")
show role.delete()
```

#### discord.Channel.rename

`func discord.Channel.rename(name:string) -> bool`

EN: Rename a Discord channel.

ES: Cambia el nombre de un canal de Discord.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.rename("support")
```

#### discord.Channel.topic

`func discord.Channel.topic(topic:string) -> bool`

EN: Change the topic of a channel that supports it.

ES: Cambia el tema de un canal compatible.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.topic("Updates")
```

#### discord.Channel.move

`func discord.Channel.move(parent_id="") -> bool`

EN: Move a channel to a category or unset its parent.

ES: Mueve un canal a una categoría o retira la categoría actual.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.move("456")
```

#### discord.Channel.delete

`func discord.Channel.delete() -> bool`

EN: Delete the selected Discord channel.

ES: Elimina el canal de Discord seleccionado.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.delete()
```

#### discord.Channel.edit_message

`func discord.Channel.edit_message(message_id:string,content:string) -> bool`

EN: Edit the text of a message owned by the bot.

ES: Edita el texto de un mensaje del propio bot.

Native example flag: true.

Note: Requires an authorized Discord application and native runtime.

Nota: Requiere una aplicación de Discord autorizada y ejecución nativa.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.edit_message("456","Listo")
```

#### discord.Guild.move_role

`func discord.Guild.move_role(role_id:string,position:int) -> bool`

EN: Move a role to a position above @everyone.

ES: Mueve un rol a una posición por encima de @everyone.

Native example flag: true.

Target: "native".

Note: Run through the Windows or Linux launcher with an authorized Discord application.

Nota: Ejecuta mediante el lanzador de Windows o Linux con una aplicación de Discord autorizada.

```tripol
import tripol.discord

server=discord.guild("123")
show server.move_role("456",2)
```

#### discord.Guild.move_channel

`func discord.Guild.move_channel(channel_id:string,position:int) -> bool`

EN: Move a channel to a numeric position in the server.

ES: Mueve un canal a una posición numérica del servidor.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

server=discord.guild("321")
show server.move_channel("123",2)
```

#### discord.Guild.private_channel

`func discord.Guild.private_channel(name:string,user_id:string,parent_id="") -> string`

EN: Create a text channel for a member with @everyone denied view access.

ES: Crea un canal de texto para un miembro y deniega la vista a @everyone.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

server=discord.guild("321")
show server.private_channel("support","456")
```

### Webhooks & invitations / Webhooks e invitaciones (14)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-webhooks)

#### discord.invite_path

`func discord.invite_path(code:string) -> string`

EN: Build a Discord API path from a validated invitation code.

ES: Construye una ruta de Discord a partir de un código de invitación válido.

Native example flag: false.

```tripol
import tripol.discord

show discord.invite_path("sample")
```

Documented expected output:

```text
/api/v10/invites/sample
```

#### discord.invite_info

`func discord.invite_info(code:string,counts=false) -> string`

EN: Read information and expiration for an invitation.

ES: Consulta los datos y la caducidad de una invitación.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

show discord.invite_info("sample",true)
```

#### discord.revoke_invite

`func discord.revoke_invite(code:string) -> bool`

EN: Delete a Discord invitation code.

ES: Elimina un código de invitación de Discord.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

show discord.revoke_invite("sample")
```

#### discord.webhook_ref

`func discord.webhook_ref(id:string,credential="") -> Webhook`

EN: Create a typed webhook reference without contacting Discord.

ES: Crea una referencia tipada de webhook sin contactar con Discord.

Native example flag: false.

```tripol
import tripol.discord

hook=discord.webhook_ref("901")
show hook.id()
```

Documented expected output:

```text
901
```

#### discord.Guild.invites

`func discord.Guild.invites() -> string`

EN: List the invitations in a server.

ES: Enumera las invitaciones de un servidor.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

server=discord.guild("321")
show server.invites()
```

#### discord.Channel.create_webhook

`func discord.Channel.create_webhook(name:string) -> Webhook`

EN: Create a typed incoming webhook in a channel.

ES: Crea un acceso webhook entrante tipado en un canal.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
hook=room.create_webhook("Updates")
show hook.id()
```

#### discord.Channel.invite

`func discord.Channel.invite(max_age=86400,max_uses=0,temporary=false,unique=true) -> string`

EN: Create an invitation with an expiration and optional use limit.

ES: Crea una invitación con caducidad y límite de usos opcional.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.invite(600,1)
```

#### discord.Channel.invites

`func discord.Channel.invites() -> string`

EN: List the invitations in a channel.

ES: Enumera las invitaciones de un canal.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.invites()
```

#### discord.Webhook.id

`func discord.Webhook.id() -> string`

EN: Return the webhook identifier without exposing its credentials.

ES: Devuelve el ID del webhook sin exponer sus credenciales.

Native example flag: false.

```tripol
import tripol.discord

hook=discord.webhook_ref("901")
show hook.id()
```

Documented expected output:

```text
901
```

#### discord.Webhook.read

`func discord.Webhook.read() -> WebhookInfo`

EN: Fetch webhook metadata while omitting the credential.

ES: Consulta los metadatos del webhook y omite sus credenciales.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

hook=discord.webhook_ref("901")
info=hook.read()
show info.name
```

#### discord.Webhook.edit

`func discord.Webhook.edit(name:string,channel_id="") -> bool`

EN: Update the name of a webhook and optionally its channel.

ES: Actualiza el nombre de un webhook y, opcionalmente, su canal.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

hook=discord.webhook_ref("901")
show hook.edit("News")
```

#### discord.Webhook.delete

`func discord.Webhook.delete() -> bool`

EN: Delete an incoming webhook with bot authorization.

ES: Elimina un webhook entrante con autorización del bot.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

hook=discord.webhook_ref("901")
show hook.delete()
```

#### discord.Webhook.send

`func discord.Webhook.send(content:string) -> string
func discord.Webhook.send(embed:Embed) -> string
func discord.Webhook.send(batch:Embeds) -> string
func discord.Webhook.send(layout:Components) -> string`

EN: Send text, an embed, an embed batch or a supported Components V2 layout.

ES: Envía texto, un embed, un lote de embeds o un diseño compatible de Components V2.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

hook=discord.channel_ref("123").create_webhook("Updates")
show hook.send("Release notes")
```

#### discord.Webhook.send_files

`func discord.Webhook.send_files(layout:Components,paths:string_list) -> string
func discord.Webhook.send_files(embed:Embed,paths:string_list) -> string
func discord.Webhook.send_files(content:string,paths:string_list) -> string`

EN: Send a message with local attachments through the native launcher.

ES: Envía un mensaje con archivos locales mediante el lanzador nativo.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

hook=discord.channel_ref("123").create_webhook("Updates")
show hook.send_files("Report",list(["report.txt"]))
```

### Threads & forums / Hilos y foros (12)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-threads)

#### discord.archive_minutes

`func discord.archive_minutes(duration:string) -> int`

EN: Convert a named Discord thread archive period to minutes.

ES: Convierte en minutos un periodo de archivado con nombre.

Native example flag: false.

```tripol
import tripol.discord

show discord.archive_minutes("day")
```

Documented expected output:

```text
1440
```

#### discord.forum_tags

`func discord.forum_tags() -> ForumTags`

EN: Construct the replacement set for available forum labels.

ES: Construye el conjunto de sustitución de las etiquetas disponibles.

Native example flag: false.

```tripol
import tripol.discord

tags=discord.forum_tags()
tags.tag("Ideas","",false,"✅")
show tags.payload()
```

Documented expected output:

```text
[{"name":"Ideas","moderated":false,"emoji_name":"✅"}]
```

#### discord.Channel.thread

`func discord.Channel.thread(name:string,archive="day",private=false) -> string`

EN: Start a public or private thread with a named archive period.

ES: Inicia un hilo público o privado con un periodo de archivado por nombre.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.thread("Planning","hour",true)
```

#### discord.Channel.thread_from

`func discord.Channel.thread_from(message_id:string,name:string,archive="day") -> string`

EN: Start a thread under an existing message.

ES: Inicia un hilo a partir de un mensaje existente.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.thread_from("456","Planning")
```

#### discord.Channel.archive

`func discord.Channel.archive(enabled=true) -> bool`

EN: Archive or unarchive a thread or forum publication.

ES: Archiva o desarchiva un hilo o publicación de foro.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.archive()
```

#### discord.Channel.lock_thread

`func discord.Channel.lock_thread(enabled=true) -> bool`

EN: Change the locked state of a thread or forum publication.

ES: Cambia el estado de bloqueo de un hilo o publicación de foro.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.lock_thread()
```

#### discord.Channel.thread_member

`func discord.Channel.thread_member(user_id="@me",enabled=true) -> bool`

EN: Add or remove a member, or join or leave as the bot.

ES: Añade o retira un miembro o entra y sal como bot.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.thread_member()
```

#### discord.Channel.post

`func discord.Channel.post(title:string,content:string,archive="day") -> string
func discord.Channel.post(title:string,content:string,tags:string_list,archive="day") -> string`

EN: Create a forum post with optional tag IDs and archive period.

ES: Crea una publicación de foro con IDs de etiquetas y archivado opcionales.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.post("Planning","First post",list(["456"]))
```

#### discord.Channel.apply_tags

`func discord.Channel.apply_tags(tags:string_list) -> bool`

EN: Set up to five tag IDs on a forum publication.

ES: Aplica hasta cinco IDs de etiquetas a una publicación de foro.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
show room.apply_tags(list(["456"]))
```

#### discord.ForumTags.tag

`func discord.ForumTags.tag(name:string,id="",moderated=false,emoji="") -> bool`

EN: Add a new or existing label and optional emoji to the set.

ES: Añade una etiqueta nueva o existente y un emoji opcional al conjunto.

Native example flag: false.

```tripol
import tripol.discord

tags=discord.forum_tags()
show tags.tag("Ideas","",false,"✅")
```

Documented expected output:

```text
true
```

#### discord.ForumTags.payload

`func discord.ForumTags.payload() -> string`

EN: Serialize the configured set of forum tags.

ES: Serializa el conjunto configurado de etiquetas del foro.

Native example flag: false.

```tripol
import tripol.discord

tags=discord.forum_tags()
tags.tag("Ideas","",false,"✅")
show tags.payload()
```

Documented expected output:

```text
[{"name":"Ideas","moderated":false,"emoji_name":"✅"}]
```

#### discord.Channel.set_forum_tags

`func discord.Channel.set_forum_tags(tags:ForumTags) -> bool`

EN: Replace the full set of available forum labels.

ES: Sustituye el conjunto completo de etiquetas disponibles en el foro.

Native example flag: true.

Target: "native".

Note: Use the Windows or Linux launcher and an authorized Discord app for network operations.

Nota: Usa el lanzador de Windows o Linux y una aplicación de Discord autorizada para las operaciones de red.

```tripol
import tripol.discord

room=discord.channel_ref("123")
tags=discord.forum_tags()
tags.tag("Ideas","",false,"✅")
show room.set_forum_tags(tags)
```

### Low-level helpers & legacy API / Utilidades de bajo nivel y API clásica (21)

[Web category](https://tripollas.github.io/Tripol/docs/libraries/discord.html#api-helpers)

#### discord.channel_path

`func discord.channel_path(channel:string) -> string`

EN: Build the API path used to access messages in a channel.

ES: Construye la ruta de la API que permite acceder a los mensajes de un canal.

Native example flag: true.

Note: The channel ID must contain 1–20 digits. This example does not contact Discord.

Nota: El ID debe contener entre 1 y 20 dígitos. Este ejemplo no contacta con Discord.

```tripol
import tripol.discord

show discord.channel_path("123")
```

Documented expected output:

```text
/api/v10/channels/123/messages
```

#### discord.send

`func discord.send(channel:string,message:string) -> bool`

EN: Send a message using the token stored in the local environment.

ES: Envía un mensaje usando el token configurado en el entorno local.

Native example flag: true.

Note: Replace the channel ID and set DISCORD_TOKEN. Requires permission to send messages. Message: 1–512 UTF-8 bytes.

Nota: Sustituye el ID y configura DISCORD_TOKEN. Requiere permiso para enviar mensajes. Mensaje: 1–512 bytes UTF-8.

Return metadata: ["True if the request returns a nonempty response; network errors may stop execution.","true si la petición devuelve una respuesta no vacía; los fallos de red pueden detener el programa."]

```tripol
import tripol.discord

show discord.send("YOUR_CHANNEL_ID", "Hola!")
```

#### discord.latest

`func discord.latest(channel:string) -> string`

EN: Fetch the newest message in a channel as JSON.

ES: Obtiene el mensaje más reciente de un canal como JSON.

Native example flag: true.

Note: Set DISCORD_TOKEN and replace the ID. The bot needs access to the channel and message history.

Nota: Configura DISCORD_TOKEN y sustituye el ID. El bot necesita acceso al canal y al historial.

Return metadata: ["Message JSON or empty text if no usable message is found.","JSON del mensaje o texto vacío si no se encuentra un mensaje utilizable."]

```tripol
import tripol.discord

message = discord.latest("YOUR_CHANNEL_ID")
show len(message) > 0
```

#### discord.message_id

`func discord.message_id(message:string) -> int`

EN: Extract a numeric ID from a message JSON object.

ES: Extrae el ID numérico de un mensaje JSON.

Native example flag: true.

Note: Returns zero if the ID is absent, invalid or outside the integer range.

Nota: Devuelve cero si el ID falta, no es válido o está fuera del rango entero.

```tripol
import tripol.discord

message = "{\"id\":\"123\"}"
show discord.message_id(message)
```

Documented expected output:

```text
123
```

#### discord.matches

`func discord.matches(message:string,command:string) -> bool`

EN: Check whether a message is an exact command from a non-bot, non-webhook author.

ES: Comprueba si un mensaje coincide exactamente con un comando y no procede de un bot ni webhook.

Native example flag: false.

Note: Compares the decoded message text.

Nota: Compara el texto del mensaje después de interpretar los escapes JSON.

```tripol
import tripol.discord

message = "{\"content\":\"!caf\\u00e9\",\"author\":{\"bot\":false}}"
show discord.matches(message, "!café")
show discord.matches(message, "!adios")
```

Documented expected output:

```text
true
false
```

#### discord.listen

`func discord.listen(channel:string,command:string,response:string) -> int`

EN: Start a legacy bot loop for one command in one channel.

ES: Inicia un bot antiguo que escucha un único comando en un canal.

Native example flag: true.

Note: Native Windows REST polling; needs DISCORD_TOKEN and channel permissions. Polls the latest message every 2500 ms; use run() for Gateway bots. Stop with Ctrl+C.

Nota: Sondeo REST en Windows nativo; requiere DISCORD_TOKEN y permisos. Consulta el último mensaje cada 2500 ms; usa run() para bots Gateway. Detén con Ctrl+C.

```tripol
import tripol.discord

discord.listen("YOUR_CHANNEL_ID", "!hola", "Hola!")
```

#### discord.channel

`func discord.channel(id:string) -> bool`

EN: Select the channel filter for incoming Gateway commands and the destination for say().

ES: Selecciona el filtro de canal para comandos Gateway y el destino de say().

Native example flag: true.

Note: The example only configures a numeric ID; use a real channel ID when connecting.

Nota: El ejemplo solo configura un ID numérico; usa un ID real al conectarte.

```tripol
import tripol.discord

show discord.channel("123")
show discord.has_channel()
```

Documented expected output:

```text
true
true
```

#### discord.reply

`func discord.reply(command:string,response:string) -> bool`

EN: Register or replace the response to an exact command.

ES: Registra o reemplaza la respuesta de un comando exacto.

Native example flag: true.

Note: Maximum 16 commands; command 1–128 bytes, response 1–512 bytes.

Nota: Máximo 16 comandos; comando de 1–128 bytes, respuesta de 1–512 bytes.

```tripol
import tripol.discord

discord.reply("!hola", "Hola!")
show discord.preview("!hola")
```

Documented expected output:

```text
Hola!
```

#### discord.preview

`func discord.preview(command:string) -> string`

EN: Look up a registered command response without sending a message.

ES: Consulta una respuesta registrada sin enviar un mensaje.

Native example flag: true.

Return metadata: ["Response text or empty text if no command matches.","Texto de respuesta o texto vacío si ningún comando coincide."]

```tripol
import tripol.discord

discord.reply("!hola", "Hola!")
show discord.preview("!hola")
```

Documented expected output:

```text
Hola!
```

#### discord.dispatch

`func discord.dispatch(message:string) -> string`

EN: Find the registered response to a raw message JSON without making a network request.

ES: Busca la respuesta registrada para un mensaje JSON sin realizar peticiones de red.

Native example flag: true.

Note: Bot and webhook messages are ignored.

Nota: Se ignoran mensajes de bots y webhooks.

```tripol
import tripol.discord

discord.reply("!hola", "Hola!")
message = "{\"id\":\"123\",\"content\":\"!hola\",\"author\":{\"bot\":false}}"
show discord.dispatch(message)
```

Documented expected output:

```text
Hola!
```

#### discord.send_to

`func discord.send_to(id:string,message:string) -> bool`

EN: Send a message to a specific channel using the configured token.

ES: Envía un mensaje a un canal concreto utilizando el token configurado.

Native example flag: true.

Note: Replace the ID, configure DISCORD_TOKEN and grant the bot channel permissions.

Nota: Sustituye el ID, configura DISCORD_TOKEN y concede permisos al bot.

```tripol
import tripol.discord

discord.login_env()
show discord.send_to("YOUR_CHANNEL_ID", "Hola!")
```

#### discord.has_channel

`func discord.has_channel() -> bool`

EN: Check whether a channel has been selected.

ES: Comprueba si hay un canal seleccionado.

Native example flag: true.

```tripol
import tripol.discord

discord.channel("123")
show discord.has_channel()
```

Documented expected output:

```text
true
```

#### discord.all_channels

`func discord.all_channels() -> bool`

EN: Remove the channel filter so Gateway replies can work in permitted channels.

ES: Elimina el filtro para responder en canales permitidos mediante Gateway.

Native example flag: true.

Note: say() and welcome messages still require a selected channel.

Nota: say() y los mensajes de bienvenida siguen necesitando un canal seleccionado.

```tripol
import tripol.discord

discord.channel("123")
show discord.all_channels()
show discord.has_channel()
```

Documented expected output:

```text
true
false
```

#### discord.say

`func discord.say(message:string) -> bool`

EN: Send a message to the currently selected channel.

ES: Envía un mensaje al canal seleccionado.

Native example flag: true.

Note: Requires login, a real channel ID and permission to send messages.

Nota: Requiere iniciar sesión, un ID real y permiso para enviar mensajes.

```tripol
import tripol.discord

discord.login_env()
discord.channel("YOUR_CHANNEL_ID")
show discord.say("Hola!")
```

#### discord.parse_latest

`func discord.parse_latest(result:string) -> string`

EN: Extract the first message object from a JSON array returned by REST.

ES: Extrae el primer objeto de mensaje de un arreglo JSON devuelto por REST.

Native example flag: true.

Return metadata: ["First object as JSON or empty text on invalid input.","Primer objeto como JSON o texto vacío si la entrada no es válida."]

```tripol
import tripol.discord

result = "[{\"id\":\"123\"},{\"id\":\"456\"}]"
show discord.parse_latest(result)
```

Documented expected output:

```text
{"id":"123"}
```

#### discord.run_rest

`func discord.run_rest() -> int`

EN: Start the legacy REST polling loop using the configured channel and commands.

ES: Inicia el bucle antiguo de sondeo REST con el canal y comandos configurados.

Native example flag: true.

Note: Native Windows. Requires a selected channel, token and permissions. Polls the newest message; stop with Ctrl+C.

Nota: Windows nativo. Requiere canal, token y permisos. Consulta el último mensaje; detén con Ctrl+C.

```tripol
import tripol.discord

discord.login_env()
discord.channel("YOUR_CHANNEL_ID")
discord.reply("!hola", "Hola!")
discord.run_rest()
```

#### discord.welcome

`func discord.welcome(value:string) -> bool`

EN: Set a message to send once the Gateway is ready.

ES: Configura un mensaje que se envía cuando el Gateway está listo.

Native example flag: true.

Note: Maximum 512 bytes. Select a real channel before run(); empty text disables the greeting.

Nota: Máximo 512 bytes. Selecciona un canal real antes de run(); el texto vacío desactiva el saludo.

```tripol
import tripol.discord

discord.channel("123")
show discord.welcome("Bot listo!")
```

Documented expected output:

```text
true
```

#### discord.identify_payload

`func discord.identify_payload(intents=33281) -> string`

EN: Build the Gateway identification JSON, including the private token.

ES: Construye el JSON de identificación del Gateway, que incluye el token privado.

Native example flag: true.

Note: Advanced helper. The payload contains the token: never print, log or share it. run() calls this internally. Supported intent sets are 33281 for guild text commands and 1 for GUILDS only.

Nota: Función avanzada. El resultado contiene el token: nunca lo imprimas, registres ni compartas. run() lo utiliza internamente. Admite 33281 para comandos de texto y 1 para solo GUILDS.

```tripol
import tripol.discord

discord.login_env()
payload = discord.identify_payload()
show len(payload) > 0
```

Documented expected output:

```text
true
```

#### discord.heartbeat_payload

`func discord.heartbeat_payload(sequence:int) -> string`

EN: Build the JSON heartbeat for a Gateway sequence number.

ES: Construye el JSON de heartbeat para un número de secuencia.

Native example flag: true.

Note: Advanced protocol helper; run() sends heartbeats automatically.

Nota: Función avanzada del protocolo; run() envía heartbeats automáticamente.

```tripol
import tripol.discord

show discord.heartbeat_payload(7)
show discord.heartbeat_payload(-1)
```

Documented expected output:

```text
{"op":1,"d":7}
{"op":1,"d":null}
```

#### discord.ignore_interaction

`func discord.ignore_interaction(context:Interaction) -> string`

EN: Return an empty response for an unused callback.

ES: Devuelve texto vacío para un callback sin respuesta.

Native example flag: false.

```tripol
import tripol.discord

ctx=discord.Interaction(3,"hola","123","456","789","123","567","{\"values\":[\"42\",\"43\"]}",true)
show discord.ignore_interaction(ctx)==""
```

Documented expected output:

```text
true
```

#### discord.interaction_key_end

`func discord.interaction_key_end(key:string) -> int`

EN: Locate the start of a command route’s guild suffix.

ES: Localiza dónde comienza el sufijo de servidor en una ruta de comando.

Native example flag: false.

```tripol
import tripol.discord

show discord.interaction_key_end("ushow/321")
```

Documented expected output:

```text
5
```

---

### discord.GatewayEvents.on_http

Register a callback for a background HTTP response. / Registra un callback para una respuesta HTTP en segundo plano.

```tripol
func discord.GatewayEvents.on_http(name:string,handler:func(GatewayEvent)->bool) -> bool
```

```tripol
import tripol.discord

import tripol.http

func fetch(ctx) {
    request=http.request("GET","http://127.0.0.1:8765/health")
    ctx.defer_http(request,"health")
    ""
}
func completed(event) {
    response=event.http_response()
    event.complete_http(response.text)
}
events=discord.events()
events.on_http("health",completed)
bot=discord.bot()
bot.slash("health","Read the local status",fetch)
bot.start_with_events(events)
```

Requires an authorized private Discord server and the native launcher. / Requiere un servidor privado autorizado y el lanzador nativo.

### discord.Interaction.defer_http

Defer an interaction and start a background HTTP request. / Difiere una interacción e inicia una petición HTTP en segundo plano.

```tripol
func discord.Interaction.defer_http(request:http.Request,event:string,ephemeral=true) -> bool
```

```tripol
import tripol.discord

import tripol.http

func fetch(ctx) {
    request=http.request("GET","http://127.0.0.1:8765/health")
    ctx.defer_http(request,"health")
    ""
}
func completed(event) {
    response=event.http_response()
    event.complete_http(response.text)
}
events=discord.events()
events.on_http("health",completed)
bot=discord.bot()
bot.slash("health","Read the local status",fetch)
bot.start_with_events(events)
```

Requires an authorized private Discord server and the native launcher. / Requiere un servidor privado autorizado y el lanzador nativo.

### discord.GatewayEvent.http_response

Read a completed HTTP response through a typed result. / Lee una respuesta HTTP terminada mediante un resultado tipado.

```tripol
func discord.GatewayEvent.http_response() -> http.Response
```

```tripol
import tripol.discord

import tripol.http

func fetch(ctx) {
    request=http.request("GET","http://127.0.0.1:8765/health")
    ctx.defer_http(request,"health")
    ""
}
func completed(event) {
    response=event.http_response()
    event.complete_http(response.text)
}
events=discord.events()
events.on_http("health",completed)
bot=discord.bot()
bot.slash("health","Read the local status",fetch)
bot.start_with_events(events)
```

Requires an authorized private Discord server and the native launcher. / Requiere un servidor privado autorizado y el lanzador nativo.

### discord.GatewayEvent.complete_http

Edit the deferred reply once with the result. / Edita una vez la respuesta diferida con el resultado.

```tripol
func discord.GatewayEvent.complete_http(message:string) -> bool
func discord.GatewayEvent.complete_http(embed:Embed) -> bool
func discord.GatewayEvent.complete_http(batch:Embeds) -> bool
func discord.GatewayEvent.complete_http(layout:Components) -> bool
```

```tripol
import tripol.discord

import tripol.http

func fetch(ctx) {
    request=http.request("GET","http://127.0.0.1:8765/health")
    ctx.defer_http(request,"health")
    ""
}
func completed(event) {
    response=event.http_response()
    event.complete_http(response.text)
}
events=discord.events()
events.on_http("health",completed)
bot=discord.bot()
bot.slash("health","Read the local status",fetch)
bot.start_with_events(events)
```

Requires an authorized private Discord server and the native launcher. / Requiere un servidor privado autorizado y el lanzador nativo.

### discord.cancel_http

Cancel an upstream request owned by a deferred Discord reply. / Cancela la petición de una respuesta diferida de Discord.

```text
func discord.cancel_http(receipt:int) -> bool
```

```tripol
import tripol.discord

show discord.cancel_http(0)
```

Run offline.

Funciona sin conexión.

### discord.GatewayEvent.http_receipt

Read the receipt for an HTTP result or delivery status. / Lee el recibo del resultado HTTP o de su entrega.

```text
func discord.GatewayEvent.http_receipt() -> int
```

```tripol
import tripol.discord

ev=discord.GatewayEvent("TRIPOL_HTTP:status",0,"","","","{\"receipt\":7}")
show ev.http_receipt()
```

Run offline.

Funciona sin conexión.

### discord.GatewayEvent.http_delivered

Check whether a queued Discord edit succeeded. / Comprueba si se completó una edición de Discord.

```text
func discord.GatewayEvent.http_delivered() -> bool
```

```tripol
import tripol.discord

ev=discord.GatewayEvent("TRIPOL_HTTP_DELIVERY:status",0,"","","","{\"receipt\":7,\"ok\":true}")
show ev.http_delivered()
```

Run offline.

Funciona sin conexión.

### discord.GatewayEvent.queue_http

Queue a single edit with text, an embed or Components V2. / Encola una sola edición con texto, un embed o Components V2.

```text
func discord.GatewayEvent.queue_http(message:string) -> bool
func discord.GatewayEvent.queue_http(embed:Embed) -> bool
func discord.GatewayEvent.queue_http(layout:Components) -> bool
```

```tripol
import tripol.discord

func result(event:discord.GatewayEvent)->bool {
 panel=discord.components()
 panel.text("# Status")
 return event.queue_http(panel)
}
```

Requires the Windows or Linux launcher. Follow the authorized Discord validation guide or use the local HTTP fixture.

Requiere el lanzador de Windows o Linux. Usa la guía de validación autorizada de Discord o el servidor HTTP local.

### discord.GatewayEvents.on_http_delivery

Register a status callback for a queued edit. / Registra un callback para el estado de una edición encolada.

```text
func discord.GatewayEvents.on_http_delivery(name:string,handler:func(GatewayEvent)->bool) -> bool
```

```tripol
import tripol.discord

func delivered(event) {event.http_delivered()}
events=discord.events()
show events.on_http_delivery("status",delivered)
```

Run offline.

Funciona sin conexión.

### discord.Interaction.defer_http_task

Start a cancelable request after deferring a reply. / Inicia una petición cancelable tras diferir la respuesta.

```text
func discord.Interaction.defer_http_task(request:http.Request,event:string,ephemeral=true) -> int
```

```tripol
import tripol.discord

import tripol.http
func load(ctx:discord.Interaction)->int {
 r=http.request("GET","https://example.com/")
 return ctx.defer_http_task(r,"status")
}
```

Requires the Windows or Linux launcher. Follow the authorized Discord validation guide or use the local HTTP fixture.

Requiere el lanzador de Windows o Linux. Usa la guía de validación autorizada de Discord o el servidor HTTP local.

### discord.Interaction.defer_http_update

Start a cancelable request to update the original component message. / Inicia una petición cancelable para editar el mensaje original del componente.

```text
func discord.Interaction.defer_http_update(request:http.Request,event:string) -> int
```

```tripol
import tripol.discord

import tripol.http
func load(ctx:discord.Interaction)->int {
 r=http.request("GET","https://example.com/")
 return ctx.defer_http_update(r,"status")
}
```

Requires the Windows or Linux launcher. Follow the authorized Discord validation guide or use the local HTTP fixture.

Requiere el lanzador de Windows o Linux. Usa la guía de validación autorizada de Discord o el servidor HTTP local.

## tripol.encoding — 4 functions

Import: `import tripol.encoding`

### encoding.hex_digit

`func encoding.hex_digit(n:int) -> string`

EN: Convert 0–15 to one hexadecimal digit.

ES: Convierte un entero entre 0 y 15 en un dígito hexadecimal.

Native example flag: false.

```tripol
import tripol.encoding

show encoding.hex_digit(15)
```

Documented expected output:

```text
f
```

### encoding.hex_byte

`func encoding.hex_byte(n:int) -> string`

EN: Convert a byte to two hexadecimal digits.

ES: Convierte un byte en dos dígitos hexadecimales.

Native example flag: false.

```tripol
import tripol.encoding

show encoding.hex_byte(255)
```

Documented expected output:

```text
ff
```

### encoding.hex_int

`func encoding.hex_int(n:int) -> string`

EN: Convert a nonnegative integer into hexadecimal text.

ES: Convierte un entero no negativo en texto hexadecimal.

Native example flag: false.

```tripol
import tripol.encoding

show encoding.hex_int(256)
```

Documented expected output:

```text
100
```

### encoding.hex_text

`func encoding.hex_text(s:string) -> string`

EN: Encode UTF-8 bytes as hexadecimal text.

ES: Codifica los bytes UTF-8 en texto hexadecimal.

Native example flag: false.

Note: Maximum 1024 input bytes.

Nota: Máximo 1024 bytes de entrada.

```tripol
import tripol.encoding

show encoding.hex_text("Hi")
```

Documented expected output:

```text
4869
```

---

## tripol.env — 1 functions

Import: `import tripol.env`

### env.get

`func env.get(name:string) -> string`

EN: Read a local process environment variable.

ES: Lee una variable del entorno del proceso local.

Native example flag: true.

Note: Native Windows or Linux x86-64. Set TRIPOL_DEMO locally. Missing or oversized values return empty text.

Nota: Windows o Linux x86-64 nativos. Configura TRIPOL_DEMO localmente. Si falta o supera el límite, devuelve texto vacío.

```tripol
import tripol.env

show len(env.get("TRIPOL_DEMO")) > 0
```

---

## tripol.events — 14 functions

Import: `import tripol.events`

### events.pressed

`func events.pressed(key:int) -> bool`

EN: Check whether a key is held in the focused game window.

ES: Comprueba si se mantiene pulsada una tecla en la ventana del juego.

Native example flag: true.

Note: Native Windows; the result depends on keyboard state.

Nota: Windows nativo; el resultado depende del teclado.

```tripol
import tripol.events

import tripol.game

game.window("Keys", 320, 240)
show events.pressed(events.space())
game.close()
```

### events.left

`func events.left() -> int`

EN: Return the key code for left arrow.

ES: Devuelve el código de la tecla flecha izquierda.

Native example flag: false.

```tripol
import tripol.events

show events.left()
```

Documented expected output:

```text
37
```

### events.right

`func events.right() -> int`

EN: Return the key code for right arrow.

ES: Devuelve el código de la tecla flecha derecha.

Native example flag: false.

```tripol
import tripol.events

show events.right()
```

Documented expected output:

```text
39
```

### events.up

`func events.up() -> int`

EN: Return the key code for up arrow.

ES: Devuelve el código de la tecla flecha arriba.

Native example flag: false.

```tripol
import tripol.events

show events.up()
```

Documented expected output:

```text
38
```

### events.down

`func events.down() -> int`

EN: Return the key code for down arrow.

ES: Devuelve el código de la tecla flecha abajo.

Native example flag: false.

```tripol
import tripol.events

show events.down()
```

Documented expected output:

```text
40
```

### events.escape

`func events.escape() -> int`

EN: Return the key code for Escape.

ES: Devuelve el código de la tecla Escape.

Native example flag: false.

```tripol
import tripol.events

show events.escape()
```

Documented expected output:

```text
27
```

### events.space

`func events.space() -> int`

EN: Return the key code for space bar.

ES: Devuelve el código de la tecla barra espaciadora.

Native example flag: false.

```tripol
import tripol.events

show events.space()
```

Documented expected output:

```text
32
```

### events.pad_connected

`func events.pad_connected(controller=0) -> bool`

EN: Check whether a controller is connected.

ES: Comprueba si un mando está conectado.

Native example flag: true.

Note: Native Windows XInput; controller index 0–3.

Nota: XInput de Windows nativo; índice de mando de 0 a 3.

```tripol
import tripol.events

show events.pad_connected(0)
```

### events.pad_buttons

`func events.pad_buttons(controller=0) -> int`

EN: Read the controller button bitmask.

ES: Lee la máscara de botones del mando.

Native example flag: true.

Note: Native Windows XInput; value depends on the controller.

Nota: XInput de Windows nativo; depende del mando.

```tripol
import tripol.events

show events.pad_buttons(0)
```

### events.pad_axis

`func events.pad_axis(axis:int,controller=0) -> int`

EN: Read one of four stick axes.

ES: Lee uno de los cuatro ejes de los sticks.

Native example flag: true.

Note: Native Windows XInput; axis and controller indexes 0–3.

Nota: XInput de Windows nativo; índices de eje y mando de 0 a 3.

```tripol
import tripol.events

show events.pad_axis(0, 0)
```

### events.pad_trigger

`func events.pad_trigger(trigger:int,controller=0) -> int`

EN: Read one of two controller triggers.

ES: Lee uno de los dos gatillos del mando.

Native example flag: true.

Note: Native Windows XInput; trigger 0–1, controller 0–3.

Nota: XInput de Windows nativo; gatillo 0–1, mando 0–3.

```tripol
import tripol.events

show events.pad_trigger(0, 0)
```

### events.mouse_x

`func events.mouse_x(width:int) -> int`

EN: Read the horizontal mouse coordinate scaled to the logical width.

ES: Lee la coordenada horizontal del ratón adaptada al ancho lógico.

Native example flag: true.

Note: Native Windows focused window; position varies.

Nota: Ventana enfocada de Windows nativo; posición variable.

```tripol
import tripol.events

import tripol.game

game.window("Mouse", 320, 240)
show events.mouse_x(320)
game.close()
```

### events.mouse_y

`func events.mouse_y(height:int) -> int`

EN: Read the vertical mouse coordinate scaled to the logical height.

ES: Lee la coordenada vertical del ratón adaptada al alto lógico.

Native example flag: true.

Note: Native Windows focused window; position varies.

Nota: Ventana enfocada de Windows nativo; posición variable.

```tripol
import tripol.events

import tripol.game

game.window("Mouse", 320, 240)
show events.mouse_y(240)
game.close()
```

### events.mouse_left

`func events.mouse_left() -> bool`

EN: Check whether the left mouse button is held.

ES: Comprueba si el botón izquierdo del ratón está pulsado.

Native example flag: true.

Note: Native Windows focused window.

Nota: Ventana enfocada de Windows nativo.

```tripol
import tripol.events

import tripol.game

game.window("Mouse", 320, 240)
show events.mouse_left()
game.close()
```

---

## tripol.game — 7 functions

Import: `import tripol.game`

### game.window

`func game.window(title:string,w:int,h:int) -> bool`

EN: Open a native game window with a title and size.

ES: Abre una ventana nativa de juego con título y tamaño.

Native example flag: false.

Note: Bounded browser preview; live input is not implemented.

Nota: Vista previa web limitada; entrada en tiempo real aún no disponible.

```tripol
import tripol.game
import tripol.canvas
import tripol.colors

w = 160
h = 90
game.window("TRIPOL Canvas", w, h)
pixels = canvas.create(w, h)
canvas.clear(pixels, colors.rgb(22, 28, 45))
canvas.rectangle(pixels, w, h, 20, 20, 55, 35, colors.rgb(75, 180, 230))
canvas.circle(pixels, w, h, 110, 45, 20, colors.rgb(255, 170, 70))
game.present(pixels, w, h)
game.close()
```

### game.running

`func game.running() -> bool`

EN: Check whether the game window remains open.

ES: Comprueba si la ventana del juego sigue abierta.

Native example flag: false.

Note: Bounded browser preview; live input is not implemented.

Nota: Vista previa web limitada; entrada en tiempo real aún no disponible.

```tripol
import tripol.game

game.window("Preview", 64, 64)
show game.running()
game.close()
```

### game.present

`func game.present(pixels:ints,w:int,h:int) -> bool`

EN: Display a packed RGB framebuffer in the window.

ES: Muestra un arreglo de píxeles RGB en la ventana.

Native example flag: false.

Note: Bounded browser preview; live input is not implemented.

Nota: Vista previa web limitada; entrada en tiempo real aún no disponible.

```tripol
import tripol.game
import tripol.canvas
import tripol.colors

w = 160
h = 90
game.window("TRIPOL Canvas", w, h)
pixels = canvas.create(w, h)
canvas.clear(pixels, colors.rgb(22, 28, 45))
canvas.rectangle(pixels, w, h, 20, 20, 55, 35, colors.rgb(75, 180, 230))
canvas.circle(pixels, w, h, 110, 45, 20, colors.rgb(255, 170, 70))
game.present(pixels, w, h)
game.close()
```

### game.key

`func game.key(code:int) -> bool`

EN: Read a virtual-key code while the game has focus.

ES: Lee una tecla virtual mientras el juego tiene el foco.

Native example flag: true.

Note: 32 is the space bar; the result depends on input.

Nota: 32 representa la barra espaciadora; el resultado depende de la entrada.

```tripol
import tripol.game

game.window("Demo", 320, 240)
show game.key(32)
game.close()
```

### game.close

`func game.close() -> bool`

EN: Close the current game window.

ES: Cierra la ventana actual del juego.

Native example flag: false.

Note: Bounded browser preview; live input is not implemented.

Nota: Vista previa web limitada; entrada en tiempo real aún no disponible.

```tripol
import tripol.game

game.window("Preview", 64, 64)
show game.close()
```

### game.wait

`func game.wait(ms:int) -> int`

EN: Pause the game loop for a number of milliseconds.

ES: Pausa el bucle del juego durante varios milisegundos.

Native example flag: false.

Note: Bounded browser preview; live input is not implemented.

Nota: Vista previa web limitada; entrada en tiempo real aún no disponible.

```tripol
import tripol.game
import tripol.canvas

w = 96
h = 64
game.window("Moving square", w, h)
pixels = canvas.create(w, h)
for x in range(0, 40) {
    canvas.clear(pixels, 1052688)
    canvas.rectangle(pixels, w, h, x * 2, 24, 12, 12, 65535)
    game.present(pixels, w, h)
    game.wait(25)
}
game.close()
```

### game.gl_frame

`func game.gl_frame(vertices:ints,background=0) -> bool`

EN: Present experimental OpenGL triangles as x, y, RGB triples.

ES: Presenta triángulos OpenGL experimentales mediante grupos x, y, RGB.

Native example flag: false.

Note: Bounded browser preview; live input is not implemented.

Nota: Vista previa web limitada; entrada en tiempo real aún no disponible.

```tripol
import tripol.game

game.window("RGB Triangle", 160, 120)
v = ints(9)
v[0] = -800
v[1] = -800
v[2] = 16711680
v[3] = 800
v[4] = -800
v[5] = 65280
v[6] = 0
v[7] = 800
v[8] = 255
game.gl_frame(v)
game.close()
```

---

## tripol.http — 16 functions + 42 methods

### http.get_raw

Send a legacy HTTP/1.0 request and return the raw message. / Hace una petición HTTP/1.0 antigua y devuelve el mensaje completo.

```tripol
func http.get_raw(address:int,port:int,host:string,path:string) -> string
```

```tripol
import tripol.http

show http.get_raw(2130706433,8765,"localhost","/health")
```

Plain HTTP legacy transport. Uses the original platform TCP support. / Transporte antiguo HTTP sin TLS. Usa el TCP original de la plataforma.

### http.status

Extract the status from a raw HTTP/1.x message. / Extrae el estado de un mensaje HTTP/1.x sin procesar.

```tripol
func http.status(response:string) -> int
```

```tripol
import tripol.http

show http.status("HTTP/1.0 200 OK\r\n\r\nHello")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
200
```

### http.body

Extract the body of a raw HTTP/1.x message. / Extrae el cuerpo de un mensaje HTTP/1.x.

```tripol
func http.body(response:string) -> string
```

```tripol
import tripol.http

show http.body("HTTP/1.0 200 OK\r\n\r\nHello")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
Hello
```

### http.request

Build a request from a method and URL. / Construye una petición a partir del método y la URL.

```tripol
func http.request(method:string,url:string) -> Request
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.method
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
GET
```

### http.secure_request

Build a request that requires HTTPS. / Construye una petición que exige HTTPS.

```tripol
func http.secure_request(method:string,url:string) -> Request
```

```tripol
import tripol.http

r=http.secure_request("GET","https://example.com/")
show r.url
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
https://example.com/
```

### http.parse_response

Decode the launcher response into a Response object. / Decodifica el resultado del lanzador en un objeto Response.

```tripol
func http.parse_response(raw:string) -> Response
```

```tripol
import tripol.http

r=http.parse_response("{\"status\":200,\"text\":\"{\\\"ok\\\":true}\",\"headers\":{\"x-test\":\"ready\"}}")
show r.status
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
200
```

### http.secure_client

Create a client restricted to HTTPS. / Crea un cliente restringido a HTTPS.

```tripol
func http.secure_client() -> Client
```

```tripol
import tripol.http

c=http.secure_client()
show c.base_url("https://example.com")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.client

Create a client with connection pooling and optional cookies. / Crea un cliente con conexiones reutilizadas y cookies opcionales.

```tripol
func http.client() -> Client
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.timeout(5000)
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.get

Read a resource. / Lee un recurso.

```tripol
func http.get(url:string) -> Response
```

```tripol
import tripol.http

r=http.get("http://127.0.0.1:8765/health")
show r.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.head

Fetch resource headers without its body. / Consulta los encabezados de un recurso.

```tripol
func http.head(url:string) -> Response
```

```tripol
import tripol.http

r=http.head("http://127.0.0.1:8765/health")
show r.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.options

Send an OPTIONS request. / Envía una petición OPTIONS.

```tripol
func http.options(url:string) -> Response
```

```tripol
import tripol.http

r=http.options("http://127.0.0.1:8765/health")
show r.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.delete

Send a DELETE request. / Envía una petición DELETE.

```tripol
func http.delete(url:string) -> Response
```

```tripol
import tripol.http

r=http.delete("http://127.0.0.1:8765/health")
show r.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.post

Send JSON in a POST. / Envía JSON en una petición POST.

```tripol
func http.post(url:string,body:string) -> Response
func http.post(url:string,body:json.Object) -> Response
```

```tripol
import tripol.http

r=http.post("http://127.0.0.1:8765/echo", {"name":"Alex","active":true})
show r.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.put

Send JSON in a PUT. / Envía JSON en una petición PUT.

```tripol
func http.put(url:string,body:string) -> Response
```

```tripol
import tripol.http

r=http.put("http://127.0.0.1:8765/echo", "{}")
show r.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.patch

Send JSON in a PATCH. / Envía JSON en una petición PATCH.

```tripol
func http.patch(url:string,body:string) -> Response
```

```tripol
import tripol.http

r=http.patch("http://127.0.0.1:8765/echo", "{}")
show r.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.download

Stream a file to a new path with optional progress updates. / Transmite un archivo a una ruta nueva con progreso opcional.

```tripol
func http.download(url:string,path:string) -> Response
func http.download(url:string,path:string,callback:func(Progress)->void) -> Response
```

```tripol
import tripol.http

r=http.download("http://127.0.0.1:8765/health","health.json",func(progress) {show progress.percent})
show r.ok()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.get

Read a resource. Reuse the client session. / Lee un recurso. Reutiliza la sesión.

```tripol
func http.Client.get(url:string) -> Response
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.get("http://127.0.0.1:8765/health")
show r.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.head

Fetch resource headers without its body. Reuse the client session. / Consulta los encabezados de un recurso. Reutiliza la sesión.

```tripol
func http.Client.head(url:string) -> Response
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.head("http://127.0.0.1:8765/health")
show r.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.options

Send an OPTIONS request. Reuse the client session. / Envía una petición OPTIONS. Reutiliza la sesión.

```tripol
func http.Client.options(url:string) -> Response
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.options("http://127.0.0.1:8765/health")
show r.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.delete

Send a DELETE request. Reuse the client session. / Envía una petición DELETE. Reutiliza la sesión.

```tripol
func http.Client.delete(url:string) -> Response
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.delete("http://127.0.0.1:8765/health")
show r.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.post

Send JSON in a POST. Reuse the client session. / Envía JSON en una petición POST. Reutiliza la sesión.

```tripol
func http.Client.post(url:string,body:string) -> Response
func http.Client.post(url:string,body:json.Object) -> Response
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.post("http://127.0.0.1:8765/echo", {"name":"Alex","active":true})
show r.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.put

Send JSON in a PUT. Reuse the client session. / Envía JSON en una petición PUT. Reutiliza la sesión.

```tripol
func http.Client.put(url:string,body:string) -> Response
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.put("http://127.0.0.1:8765/echo", "{}")
show r.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.patch

Send JSON in a PATCH. Reuse the client session. / Envía JSON en una petición PATCH. Reutiliza la sesión.

```tripol
func http.Client.patch(url:string,body:string) -> Response
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.patch("http://127.0.0.1:8765/echo", "{}")
show r.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Request.header

Add a header to a request. / Añade un encabezado a la petición.

```tripol
func http.Request.header(name:string,value:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.header("Accept","application/json")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.query

Append a URL-encoded query value. / Añade un valor codificado a la consulta.

```tripol
func http.Request.query(name:string,value:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.query("q","a & b")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.form

Append a URL-encoded or multipart field. / Añade un campo codificado en URL o multipart.

```tripol
func http.Request.form(name:string,value:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.form("name","Alex")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.file

Add a local file part to a multipart upload. / Añade un archivo local a una subida multipart.

```tripol
func http.Request.file(name:string,path:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.file("report","report.txt")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.bearer

Set Bearer authorization. / Configura la autorización Bearer.

```tripol
func http.Request.bearer(token:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.bearer("fixture-token")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.basic

Set HTTP Basic authorization. / Configura la autorización HTTP Basic.

```tripol
func http.Request.basic(user:string,password:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.basic("demo","demo")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.api_key

Set an API key header. / Configura un encabezado de clave API.

```tripol
func http.Request.api_key(name:string,value:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.api_key("X-API-Key","fixture")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.timeout

Set the timeout in milliseconds. / Configura el tiempo de espera en milisegundos.

```tripol
func http.Request.timeout(milliseconds:int) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.timeout(5000)
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.redirects

Choose no redirects or same-origin GET and HEAD redirects. / Elige sin redirecciones o redirecciones GET y HEAD del mismo origen.

```tripol
func http.Request.redirects(policy="same-origin") -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.redirects("same-origin")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.retries

Set the number of safe retries for bodyless reads. / Configura los reintentos seguros de lecturas sin cuerpo.

```tripol
func http.Request.retries(count=1) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.retries(2)
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.text

Set a plain text request body. / Configura un cuerpo de texto.

```tripol
func http.Request.text(value:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.text("Hola")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.json

Set a JSON string or typed object as the body. / Configura una cadena JSON u objeto tipado como cuerpo.

```tripol
func http.Request.json(value:string) -> bool
func http.Request.json(value:json.Object) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.json({"active":true})
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.download_to

Set a destination for a streaming GET. / Configura un destino para un GET por flujo.

```tripol
func http.Request.download_to(path:string) -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
show r.download_to("health.json")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Request.send

Send the prepared request and wait for its response. / Envía la petición preparada y espera su respuesta.

```tripol
func http.Request.send() -> Response
func http.Request.send(callback:func(Progress)->void) -> Response
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
response=r.send()
show response.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Request.start

Execute the request in the launcher background. / Ejecuta la petición en segundo plano en el lanzador.

```tripol
func http.Request.start() -> Task
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
task=r.start()
response=task.await()
show response.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Task.id

Return the current job identifier. / Devuelve el identificador de la tarea actual.

```tripol
func http.Task.id() -> int
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
task=r.start()
show task.id()
task.await()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Task.poll

Collect a completed response without blocking. / Recoge una respuesta terminada sin bloquear.

```tripol
func http.Task.poll() -> Response
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
task=r.start()
response=task.poll()
show response.error
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Task.await

Wait for a response and collect it once. / Espera la respuesta y recógela una vez.

```tripol
func http.Task.await() -> Response
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
task=r.start()
response=task.await()
show response.status
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Task.cancel

Request cancellation of a job. / Solicita la cancelación de la tarea.

```tripol
func http.Task.cancel() -> bool
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
task=r.start()
show task.cancel()
task.await()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Task.progress

Read the bytes transferred and the expected total. / Lee los bytes transferidos y el total esperado.

```tripol
func http.Task.progress() -> Progress
```

```tripol
import tripol.http

r=http.request("GET","http://127.0.0.1:8765/health")
task=r.start()
progress=task.progress()
show progress.bytes
task.await()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Response.ok

Check for a successful 2xx result. / Comprueba si el resultado es 2xx.

```tripol
func http.Response.ok() -> bool
```

```tripol
import tripol.http

r=http.parse_response("{\"status\":200,\"text\":\"{\\\"ok\\\":true}\",\"headers\":{\"x-test\":\"ready\"}}")
show r.ok()
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Response.json

Read the validated response body as JSON text. / Lee el cuerpo validado como texto JSON.

```tripol
func http.Response.json() -> string
```

```tripol
import tripol.http

r=http.parse_response("{\"status\":200,\"text\":\"{\\\"ok\\\":true}\",\"headers\":{\"x-test\":\"ready\"}}")
show r.json()
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
{"ok":true}
```

### http.Response.header

Read the first value of a response header. / Lee el primer valor de un encabezado.

```tripol
func http.Response.header(name:string) -> string
```

```tripol
import tripol.http

r=http.parse_response("{\"status\":200,\"text\":\"{\\\"ok\\\":true}\",\"headers\":{\"x-test\":\"ready\"}}")
show r.header("X-Test")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
ready
```

### http.Client.base_url

Set the base URL of the client. / Configura la URL base del cliente.

```tripol
func http.Client.base_url(url:string) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.base_url("https://example.com/")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.cookies

Enable or disable cookies before opening a session. / Habilita o deshabilita cookies antes de abrir la sesión.

```tripol
func http.Client.cookies(enabled=true) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.cookies(false)
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.timeout

Set the default timeout in milliseconds. / Configura el tiempo de espera predeterminado.

```tripol
func http.Client.timeout(milliseconds:int) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.timeout(5000)
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.bearer

Set the client Bearer credential. / Configura la credencial Bearer del cliente.

```tripol
func http.Client.bearer(token:string) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.bearer("fixture-token")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.basic

Set the client Basic credential. / Configura la credencial Basic del cliente.

```tripol
func http.Client.basic(user:string,password:string) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.basic("demo","demo")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.api_key

Set the client API key header. / Configura el encabezado de clave API del cliente.

```tripol
func http.Client.api_key(name:string,value:string) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.api_key("X-API-Key","fixture")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.redirects

Set the default redirect policy. / Configura la política de redirección.

```tripol
func http.Client.redirects(policy="same-origin") -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.redirects("same-origin")
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.retries

Set the client retry count for safe methods. / Configura los reintentos para métodos seguros.

```tripol
func http.Client.retries(count=1) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
show c.retries(2)
```

Runs offline. The example does not open a socket. / Funciona sin red. El ejemplo no abre una conexión.

Expected / Resultado

```text
true
```

### http.Client.request

Prepare a request using the client defaults. / Prepara una petición con los valores predeterminados del cliente.

```tripol
func http.Client.request(method:string,url:string) -> Request
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.request("GET","/health")
response=r.send()
show response.status
c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.close

Close the connection pool and cancel its jobs. / Cierra las conexiones y cancela las tareas del cliente.

```tripol
func http.Client.close() -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("http://127.0.0.1:8765")
r=c.get("/health")
show c.close()
```

Run with the Windows or Linux launcher. Examples using 127.0.0.1 use the local fixture described in the guide. / Ejecuta con el lanzador en Windows o Linux. Los ejemplos con 127.0.0.1 usan el servidor local descrito en la guía.

### http.Client.cookie_file

Save expiring cookies in an isolated file for the base origin. / Guarda cookies con caducidad en un archivo aislado para el origen base.

```text
func http.Client.cookie_file(path:string) -> bool
```

```tripol
import tripol.http

c=http.client()
c.base_url("https://example.com")
show c.cookie_file("private-cookies.json")
```

Configure before opening the client. Use a private directory and select a separate file for each client.

Configura antes de abrir el cliente. Usa un directorio privado y un archivo por cliente.

### http.Client.metrics

Read request, error, connection, transfer and elapsed time counters for this client. / Consulta contadores de peticiones, errores, conexiones, transferencia y latencia para este cliente.

```text
func http.Client.metrics() -> ClientMetrics
```

```tripol
import tripol.http

c=http.client()
r=c.get("http://127.0.0.1:8765/health")
stats=c.metrics()
show stats.completed,stats.connections
c.close()
```

Requires the Windows or Linux launcher. Follow the authorized Discord validation guide or use the local HTTP fixture.

Requiere el lanzador de Windows o Linux. Usa la guía de validación autorizada de Discord o el servidor HTTP local.

## tripol.https — 3 functions

### https.get

Get a full HTTPS URL or use the original host/path call. / Consulta una URL HTTPS completa o usa la llamada antigua de host y ruta.

```tripol
func https.get(host:string,path:string) -> string
func https.get(url:string) -> http.Response
```

```tripol
import tripol.https

r=https.get("https://example.com/")
show r.status
```

Uses TLS with the system trust store. The older get(host,path) overload uses Windows WinHTTP. / Usa TLS con certificados del sistema. La sobrecarga antigua get(host,path) usa WinHTTP en Windows.

### https.request

Build a request restricted to TLS. / Construye una petición que requiere TLS.

```tripol
func https.request(method:string,url:string) -> http.Request
```

```tripol
import tripol.https

r=https.request("HEAD","https://example.com/")
result=r.send()
show result.status
```

Uses TLS with the system trust store. The older get(host,path) overload uses Windows WinHTTP. / Usa TLS con certificados del sistema. La sobrecarga antigua get(host,path) usa WinHTTP en Windows.

### https.client

Create an HTTPS-only client. / Crea un cliente que solo admite HTTPS.

```tripol
func https.client() -> http.Client
```

```tripol
import tripol.https

c=https.client()
c.base_url("https://example.com")
response=c.head("/")
show response.status
c.close()
```

Uses TLS with the system trust store. The older get(host,path) overload uses Windows WinHTTP. / Usa TLS con certificados del sistema. La sobrecarga antigua get(host,path) usa WinHTTP en Windows.

## tripol.memory — 4 functions

Import: `import tripol.memory`

### memory.capacity

`func memory.capacity() -> int`

EN: Report the target's arena capacity in bytes.

ES: Informa la capacidad de la arena de memoria del destino en bytes.

Native example flag: false.

Note: Native and browser limits differ; this is not total RAM.

Nota: Los límites nativos y web difieren; no es la RAM total.

```tripol
import tripol.memory

show memory.capacity()
```

### memory.used

`func memory.used() -> int`

EN: Report the currently charged arena bytes.

ES: Informa los bytes de arena actualmente utilizados.

Native example flag: false.

```tripol
import tripol.memory

a = ints(10)
show memory.used()
```

### memory.committed

`func memory.committed() -> int`

EN: Report committed arena bytes, according to the target.

ES: Informa los bytes comprometidos de la arena según el destino.

Native example flag: false.

Note: Web reports logical charged bytes, not process memory.

Nota: La web informa bytes lógicos, no memoria del proceso.

```tripol
import tripol.memory

show memory.committed()
```

### memory.available

`func memory.available() -> int`

EN: Report arena capacity minus usage.

ES: Informa la capacidad de arena menos su uso.

Native example flag: false.

```tripol
import tripol.memory

show memory.available()
```

---

## tripol.network — 3 functions

Import: `import tripol.network`

### network.ipv4

`func network.ipv4(a:int,b:int,c:int,d:int) -> int`

EN: Pack four IPv4 octets into an unsigned 32-bit integer.

ES: Empaqueta cuatro octetos IPv4 en un entero de 32 bits sin signo.

Native example flag: false.

```tripol
import tripol.network

show network.ipv4(127, 0, 0, 1)
```

Documented expected output:

```text
2130706433
```

### network.exchange

`func network.exchange(address:int,port:int,request:string) -> string`

EN: Send a TCP request and read text until the peer closes.

ES: Envía una petición TCP y lee texto hasta que el otro extremo cierra.

Native example flag: true.

Note: Requires a local TCP test server on port 8080 and native sockets; no raw TCP in browsers.

Nota: Requiere un servidor TCP local de prueba en el puerto 8080 y sockets nativos; no hay TCP directo en navegadores.

```tripol
import tripol.network

response = network.exchange(network.loopback(), 8080, "hello")
show response
```

### network.loopback

`func network.loopback() -> int`

EN: Return the packed address 127.0.0.1.

ES: Devuelve la dirección empaquetada 127.0.0.1.

Native example flag: false.

```tripol
import tripol.network

show network.loopback()
```

Documented expected output:

```text
2130706433
```

---

## tripol.process — 3 functions

Import: `import tripol.process`

### process.id

`func process.id() -> int`

EN: Read the operating system identifier of the current native process.

ES: Lee el identificador que el sistema operativo asignó al proceso nativo actual.

Native example flag: true.

Target: "native".

Note: Requires a native executable.

Nota: Requiere un ejecutable nativo.

```tripol
import tripol.process

show process.id()>0
```

### process.platform

`func process.platform() -> string`

EN: Return the current target name for the running program.

ES: Devuelve el nombre del destino en el que se ejecuta el programa.

Native example flag: false.

```tripol
import tripol.process

show process.platform()!="unknown"
```

Documented expected output:

```text
true
```

### process.exit

`func process.exit(code=0) -> int`

EN: End the current native process with the given status code.

ES: Termina el proceso nativo actual con el código indicado.

Native example flag: true.

Target: "native".

Note: Requires a native executable.

Nota: Requiere un ejecutable nativo.

```tripol
import tripol.process

show "Done"
process.exit(0)
```

---

## tripol.profiler — 2 functions

Import: `import tripol.profiler`

### profiler.start

`func profiler.start() -> int`

EN: Start measuring monotonic elapsed time in nanoseconds.

ES: Inicia una medición de tiempo monotónico en nanosegundos.

Native example flag: false.

Note: Measures elapsed time, not CPU or heap samples.

Nota: Mide tiempo transcurrido, no muestras de CPU ni de memoria.

```tripol
import tripol.profiler

start = profiler.start()
show profiler.elapsed(start) >= 0
```

Documented expected output:

```text
true
```

### profiler.elapsed

`func profiler.elapsed(started:int) -> int`

EN: Return elapsed nanoseconds since a profiler timestamp.

ES: Devuelve los nanosegundos transcurridos desde una marca del perfilador.

Native example flag: false.

Note: The sum is 4950; the time varies by execution.

Nota: La suma es 4950; el tiempo varía según la ejecución.

```tripol
import tripol.profiler

start = profiler.start()
total = 0
for i in range(100) {total += i}
show total
show profiler.elapsed(start)
```

---

## tripol.simd — 1 functions

Import: `import tripol.simd`

### simd.equal

`func simd.equal(a:ints,b:ints) -> bool`

EN: Compare lengths and contents of two integer arrays.

ES: Compara la longitud y el contenido de dos arreglos de enteros.

Native example flag: false.

```tripol
import tripol.simd

a = [1, 2]
b = [1, 2]
show simd.equal(a, b)
```

Documented expected output:

```text
true
```

---

## tripol.system — 5 functions

Import: `import tripol.system`

### system.platform

`func system.platform() -> string`

EN: Return the current target platform name.

ES: Devuelve el nombre de la plataforma de ejecución.

Native example flag: false.

Note: Examples: windows, linux or web on the corresponding target.

Nota: Ejemplos: windows, linux o web según el destino.

```tripol
import tripol.system

show system.platform()
```

### system.architecture

`func system.architecture() -> string`

EN: Return the target architecture name.

ES: Devuelve el nombre de la arquitectura de ejecución.

Native example flag: false.

Note: Native x86_64; web target wasm32 where the adapter is enabled.

Nota: x86_64 nativo; wasm32 en la web cuando el adaptador está habilitado.

```tripol
import tripol.system

show system.architecture()
```

### system.compiler

`func system.compiler() -> string`

EN: Return the compiler version string.

ES: Devuelve el texto con la versión del compilador.

Native example flag: false.

```tripol
import tripol.system

show system.compiler()
```

### system.is_windows

`func system.is_windows() -> bool`

EN: Check whether the target is Windows.

ES: Comprueba si el destino es Windows.

Native example flag: false.

```tripol
import tripol.system

show system.is_windows()
```

### system.is_linux

`func system.is_linux() -> bool`

EN: Check whether the target is Linux.

ES: Comprueba si el destino es Linux.

Native example flag: false.

```tripol
import tripol.system

show system.is_linux()
```

---

## tripol.thread — 3 functions

Import: `import tripol.thread`

### thread.supported

`func thread.supported() -> bool`

EN: Check whether the Windows native worker is available on the target.

ES: Comprueba si el destino ofrece el trabajador nativo de Windows.

Native example flag: false.

```tripol
import tripol.thread

show thread.supported()
```

Documented expected output:

```text
false
```

### thread.sum

`func thread.sum(data:ints) -> int`

EN: Sum integers with checked arithmetic in a joined Windows worker.

ES: Suma enteros con aritmética comprobada en un trabajador de Windows al que se espera antes de continuar.

Native example flag: true.

Target: "windows".

Note: Requires native Windows.

Nota: Requiere Windows nativo.

```tripol
import tripol.thread

data=ints(2)
data[0]=3
data[1]=4
show thread.sum(data)
```

### thread.sleep

`func thread.sleep(milliseconds:int) -> int`

EN: Wait for the given number of milliseconds under the target's limits.

ES: Espera la cantidad de milisegundos indicada dentro de los límites del destino.

Native example flag: false.

```tripol
import tripol.thread

show thread.sleep(0)
```

Documented expected output:

```text
0
```

---

## tripol.ui — 7 functions

Import: `import tripol.ui`

### ui.message

`func ui.message(title:string, message:string) -> bool`

EN: Show an informational dialog and report whether it was acknowledged.

ES: Muestra un diálogo informativo e indica si fue aceptado.

Native example flag: true.

Note: Native Windows, ASCII text only; result depends on the user.

Nota: Windows nativo, solo ASCII; el resultado depende del usuario.

```tripol
import tripol.ui

show ui.message("TRIPOL", "Hello!")
```

### ui.confirm

`func ui.confirm(title:string, message:string) -> bool`

EN: Ask a yes/no question and return true for Yes.

ES: Hace una pregunta de sí o no y devuelve true al elegir Sí.

Native example flag: true.

Note: Native Windows, ASCII text only.

Nota: Windows nativo, solo texto ASCII.

```tripol
import tripol.ui

show ui.confirm("Continue?", "Start the demo?")
```

### ui.form

`func ui.form(title:string,label:string,button:string) -> string`

EN: Open a modal form and read text entered by the user.

ES: Abre un formulario modal y lee el texto introducido.

Native example flag: true.

Note: Native Windows, ASCII labels; input depends on the user.

Nota: Windows nativo, etiquetas ASCII; la entrada depende del usuario.

```tripol
import tripol.ui

name = ui.form("Profile", "Name", "Save")
show name
```

### ui.warning

`func ui.warning(title:string,message:string) -> bool`

EN: Show a warning dialog.

ES: Muestra un diálogo de advertencia.

Native example flag: true.

Note: Native Windows, ASCII text only.

Nota: Windows nativo, solo ASCII.

```tripol
import tripol.ui

show ui.warning("Warning", "Check your input")
```

### ui.error

`func ui.error(title:string,message:string) -> bool`

EN: Show an error dialog.

ES: Muestra un diálogo de error.

Native example flag: true.

Note: Native Windows, ASCII text only.

Nota: Windows nativo, solo ASCII.

```tripol
import tripol.ui

show ui.error("Error", "Invalid input")
```

### ui.retry

`func ui.retry(title:string,message:string) -> bool`

EN: Ask whether to retry; return true when Retry is selected.

ES: Pregunta si se desea reintentar; devuelve true al elegir Reintentar.

Native example flag: true.

Note: Native Windows, ASCII text only.

Nota: Windows nativo, solo ASCII.

```tripol
import tripol.ui

show ui.retry("Retry", "Try again?")
```

### ui.choose

`func ui.choose(title:string,message:string) -> int`

EN: Display Yes, No and Cancel choices and return the selected code.

ES: Muestra las opciones Sí, No y Cancelar y devuelve el código elegido.

Native example flag: true.

Note: Native Windows: Yes = 6, No = 7, Cancel = 2.

Nota: Windows nativo: Sí = 6, No = 7, Cancelar = 2.

```tripol
import tripol.ui

choice = ui.choose("Save", "Save changes?")
show choice
```

---

# Additional official guides

---

## Source: docs/genesis.md

# GENESIS R3 — QUANTUM X

The production Go compiler builds a native stage written in TRIPOL. GENESIS emits Linux x86-64 ELF and experimental Windows PE32+ console executables. It now supports integer functions (up to 8 int parameters), local integer variables, recursion, forward and mutual calls, if/else-if/else, while, for/range, checked + - * / %, and repeated show. A function must end in a return; an earlier return in a branch is permitted. Global variable capture, a complete type system are not supported. Boolean conditions have short-circuit operators. Up to 32 function declarations, 32 slots per frame, 4096 source bytes, 16 statement nesting levels. Declarations must precede control flow in each frame. Loop counters are immutable and scoped.

## Direct binary output

```sh
tripol build compiler/direct.tri --target linux -o genesis-stage
printf 'examples/projects/genesis_functions.tri\ngenesis-demo\n' | ./genesis-stage
chmod +x genesis-demo
./genesis-demo
```

On Windows build the stage with `--target windows -o genesis-stage.exe` then run it and enter the source and output paths on separate lines. The stage prints bytes written and produces an ELF for Linux x86-64. Choose a new output path. Both the direct driver and the hex transport (`compiler/file.tri` with `tripol genesis`) reject existing destinations, including paths that resolve to the source file.

Expected output of genesis_functions.tri: 0, 5, 10, ..., 45, each on its own line. Full self-compilation and additional language features remain future work.

## Windows PE R1

```cmd
tripol build compiler/windows.tri --target windows -o genesis-win.exe
genesis-win.exe
```
Enter `examples/projects/genesis_functions.tri` and a new output path such as `genesis-demo.exe` on separate lines. Run `genesis-demo.exe`. Expected ten lines from 0 to 45 in steps of 5. The compiler emits PE32+ directly, with RX code and a RW import table for GetStdHandle, WriteFile and ExitProcess. Linux can host this stage as well. The generated programs use kernel32 and a fixed image base. General Windows APIs and full self-hosting remain future work.

### GENESIS R3 additional control flow

The declaration pass now builds the function signature table before emitting code. Calls target stable native jump stubs patched to actual function bodies, enabling forward calls and mutual recursion. Conditions support integer comparisons, boolean literals, `not`, `and`, `or` (also !, &&, ||), short-circuit evaluation and boolean grouping; `else if` chains are supported. Booleans are now first-class scalar values in local variables, parameters and function returns; `show` remains integer-only in this stage. Integer and boolean expressions now pass through a separate bounded AST; statements still lower directly from tokens, functions cannot access main's local slots, and self-compilation is not yet possible. See examples/projects/genesis_control.tri.

### Loop control

`break` exits the nearest `while` or `for` loop. `continue` rechecks a `while` condition or advances a `for` counter. Both are available inside nested conditionals. GENESIS rejects loop control outside a loop.

## Integer expression AST

`compiler/ast.tri` represents literals, variable references, unary negation, checked arithmetic, calls and linked argument nodes in a bounded integer array (1024 nodes per expression). Parsing is separate from native code emission. The backend preserves left-to-right argument evaluation, precedence and checked overflow/division. Existing mutable integer variables support `+=`, `-=`, `*=`, `/=`, `%=`; range counters remain immutable. `examples/projects/genesis_ast.tri` exercises the AST on both native targets. This is not yet a full program AST or a self-hosting compiler.

### Typed scalar expressions

GENESIS now infers int/bool local types and checks assignment, function parameters/returns, arithmetic, comparisons and conditions during compilation. Boolean equality and short-circuit operators use the same AST as arithmetic, including grouped expressions. Boolean values can be returned or passed to functions, but cannot be printed directly yet; use an integer-returning function or branch. Integer range counters remain immutable. `examples/projects/genesis_types.tri` demonstrates the two scalar types. Strings, arrays, records and full program AST remain pending.

---

## Source: docs/update.md

# TRIPOL UPDATE R2 — first stage

Interactive commands can print a notice on stderr when a newer numbered release is available. Checks are cached for 24 hours (including failed offline attempts), with a 700 ms HTTP deadline on a cache miss. Redirected output, CI, benchmarks, machine-oriented commands and --version are unaffected. Set TRIPOL_NO_UPDATE=1 to disable automatic checks. Explicit commands still work.

- `tripol update --check`: refresh the official GitHub release list and report status.
- `tripol update`: download the platform-specific official distribution into a private cache directory, validate its exact published size and GitHub SHA-256 digest, then print installation instructions.
- No installation, execution or source-code changes occur automatically. Windows setup remains unsigned. Numbered public pre-releases are eligible because TRIPOL is a public preview; drafts and unrecognized version tags are excluded.

No tokens are required. Release checks contact api.github.com and explicit downloads contact GitHub release assets. No source files are sent. The cache holds public release metadata. SHA-256 detects changes to the downloaded asset relative to GitHub's release metadata; it does not replace code signing.

New in QUANTUM X: interactive transfer progress, `tripol update --download-only` (verified download without installation instructions), and `tripol update --history` (local verified download receipts, not installed-version history). Default `tripol update` still prints manual installation instructions. Automated installation, backup and rollback are not yet implemented. Receipts contain public asset name/digest and local path; they remain on the user's machine.

## UPDATE R2 transaction (0.27.1)

`tripol update --install` downloads the official compiler (raw Windows asset or Linux archive), verifies SHA-256, asks for confirmation and replaces only the managed per-user compiler. Interactive default update offers installation when a supported asset is available; `--download-only` never installs. Windows copies the current compiler to a temporary worker and waits for the parent to exit before replacement. The worker reports the result in its printed local log. Linux extracts only the known regular `tripol` member; other paths, special files and oversized archives are rejected. The candidate and installed binary are probed using --version. An exclusive lock prevents simultaneous replacement. The previous executable and its digest are retained in a .tripol-rollback-* directory next to the installation; `tripol update --rollback` confirms before restoring the newest intact backup. If a verification step fails the original is restored when possible; if restoration is blocked, the error identifies the preserved original. Close other TRIPOL processes. Portable/system/symlink installs remain manual. No admin rights or PATH changes are requested. Unsigned Windows binaries remain experimental; digests are from official release metadata, not independent code signing. The 0.27.0 updater requires manual installation of 0.27.1, and the worker is used for subsequent releases.

---

## Source: docs/ecosystem-release.md

# Public ecosystem release / Publicación

Repository: https://github.com/Tripollas/Tripol

## Release

After the ecosystem PR is reviewed and merged, open Actions → Public ecosystem → Run workflow on `main`. The workflow scans reachable Git history, tests builds on Linux and Windows, builds the unsigned Windows setup, Linux compiler archive and VSIX, and creates a **draft pre-release** with SHA-256 checksums. Review the actual logs, artifacts and unsigned installer before publishing the draft. Never replace existing release assets silently. If the workflow is not available, follow [build.md](build.md).

## VS Code Marketplace

The repository includes an extension package. Install its VSIX locally for testing. To publish publicly, create/verify a Visual Studio Marketplace publisher and set `vscode/package.json` publisher to that owned ID; use `vsce publish` from a trusted machine with a scoped token stored outside source. Do not expose tokens in CI or logs.

## Community

An administrator must enable Discussions in repository Settings → General → Features → Discussions. Suggested categories: Announcements, Q&A, Ideas, Show and tell. Use the existing code of conduct and security policy. The API connector cannot change this setting.

## Security

The workflow scans Git history of all fetched refs with Gitleaks in redacted mode; it does not scan every external artifact, release, Actions log, unreachable commit or binary. Review those separately. Enable GitHub secret scanning and push protection if available in Settings. Rotate any exposed credentials immediately, even after deleting them from a file. Never publish a real secret in a ticket or discussion. Keep private vulnerability reporting enabled where available.

---

## Source: docs/public-release.md

# Public preview release gate

- [ ] Synchronize the complete ORIGIN snapshot and verify the remote tree matches it.
- [ ] Review all branches, tags, history, releases, artifacts and logs for secrets, personal information and material not authorized for publication; revoke any leaked credentials. A current-tree scan is not a history audit.
- [ ] Confirm all public files are authorized for Apache-2.0 distribution; keep private research outside this repository.
- [ ] Verify Windows and Linux CI, then run live Discord and graphics smoke tests on an isolated test server; no credentials in CI logs.
- [ ] Publish verified installer/portable/source archives with SHA-256 checksums as a pre-release, not v1.0. Installer is unsigned.
- [ ] Enable branch protection, required CI, private vulnerability reporting and community discussions where available.
- [ ] Change visibility in Settings only after approval of these checks. Public Git history may be copied permanently.

The current repository is a private staging area until the maintainer completes these steps. No claim of a public release or full audit is made by this checklist.

---

## Source: docs/open-core.md

# Open ecosystem and licensing

The buildable compiler, runtime, official libraries, examples, tests and documentation in this repository are distributed under Apache-2.0 unless a file explicitly states another compatible license. Anyone may use, modify and redistribute them under the license, including commercial use. Preserve required notices. The project remains free to download and build.

Future private research may live in a separate TRIPOL Labs repository under independent terms. No such repository or proprietary dependency is required by the public compiler. Code already distributed under Apache-2.0 remains available under that license; private research is not represented as open source. A source-available or closed component must be labeled separately. Changes to this policy must not remove rights previously granted.

---

## Source: docs/website.md

# Website / GitHub Pages

The public website is maintained exclusively on the [gh-pages branch](https://github.com/Tripollas/Tripol/tree/gh-pages), not in main. Its public URL is https://tripollas.github.io/Tripol/.

## Local preview

Clone the website branch into a separate directory:

```sh
git clone --branch gh-pages --single-branch https://github.com/Tripollas/Tripol.git tripol-website
cd tripol-website
python -m http.server 8000
```

Open http://localhost:8000. The branch includes static files, documentation and versioned downloads. The installation scripts verify pinned SHA-256 hashes; Windows setup is unsigned.

## Deployment

In Settings > Pages, use Deploy from a branch > gh-pages > /(root). Keep changes to the website on gh-pages and compiler changes on main. Before publishing a new version, regenerate verified downloads and checksums, update installer hashes and links, and test both platforms. Never commit secrets. GitHub Releases and website downloads must be verified independently; no automatic synchronization is implied.

---

# Guided web lessons / Lecciones guiadas

---

## Lesson 1 — docs/learn-hello.html

LESSON 01 OF 14 (ES: LECCIÓN 01 DE 14)

### Hello world (ES: Hola mundo)
Write and run your first TRIPOL program. (ES: Escribe y ejecuta tu primer programa en TRIPOL.)

#### Your first program (ES: Tu primer programa)
The show instruction writes values to the console. Text is written between double quotes. No imports or main function are needed for a simple program. (ES: La instrucción show escribe valores en la consola. El texto se escribe entre comillas dobles. No necesitas importar bibliotecas ni declarar una función main para un programa sencillo.)

```tripol
show "Hello, TRIPOL!"
```

Expected output (ES: Resultado esperado)

```output
Hello, TRIPOL!
```

#### Several values and comments (ES: Varios valores y comentarios)
Use commas to show several values separated by spaces. A new show instruction starts a new output line. Comments begin with // or # and are ignored. (ES: Utiliza comas para mostrar varios valores separados por espacios. Una nueva instrucción show comienza otra línea. Los comentarios empiezan con // o # y se ignoran.)

```tripol
// My first program
show "TRIPOL", 42
# Another line
show "Ready"
```

Expected output (ES: Resultado esperado)

```output
TRIPOL 42
Ready
```

#### Run a file locally (ES: Ejecutar un archivo localmente)
Save your program as main.tri. With TRIPOL installed, open a terminal in the folder and run these commands. The browser examples can be opened with Try it instead. (ES: Guarda el programa como main.tri. Con TRIPOL instalado, abre una terminal en esa carpeta y ejecuta estos comandos. Para los ejemplos del navegador, utiliza Probar.)
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```cmd
tripol run main.tri
tripol build main.tri
```

Important (ES: Importante)The build command produces a native executable; it does not require Python or Go for users of the distribution. (ES: El comando build genera un ejecutable nativo; los usuarios de la distribución no necesitan Python ni Go.)

Installation and quick start (ES: Instalación e inicio rápido)

---

## Lesson 2 — docs/learn-variables.html

LESSON 02 OF 14 (ES: LECCIÓN 02 DE 14)

### Variables (ES: Variables)
Store values, update them and define constants. (ES: Guarda valores, modifícalos y define constantes.)

#### Create a variable (ES: Crear una variable)
A bare assignment creates a mutable local if it does not exist, or updates it if it does. TRIPOL infers the type from the assigned value. let explicitly declares a mutable binding. (ES: Una asignación directa crea una variable local mutable si no existe, o la actualiza si ya existe. TRIPOL deduce el tipo del valor asignado. let declara explícitamente una variable mutable.)

```tripol
score = 10
let bonus = 5
score = score + bonus
show score
```

Expected output (ES: Resultado esperado)

```output
15
```

#### Update a value (ES: Modificar un valor)
Use = to replace a value and compound assignments such as += or -= to update it. A variable cannot change its type after declaration. (ES: Usa = para reemplazar un valor y asignaciones compuestas como += o -= para actualizarlo. Una variable no puede cambiar de tipo después de declararse.)

```tripol
health = 100
health -= 20
health += 5
show health
```

Expected output (ES: Resultado esperado)

```output
85
```

#### Constants and scope (ES: Constantes y ámbito)
const prevents reassignment of its binding. Names are ASCII identifiers; avoid shadowing. A constant reference is not a guarantee that all shared data are immutable. (ES: const impide reasignar el nombre declarado. Los nombres son identificadores ASCII; no se admite ocultar una variable con otra del mismo nombre. Una referencia constante no garantiza que todos los datos compartidos sean inmutables.)

```tripol
const max_health = 100
health = 75
show health, max_health
```

Expected output (ES: Resultado esperado)

```output
75 100
```

Important (ES: Importante)Do not assign to max_health after its declaration. (ES: No asignes un nuevo valor a max_health después de declararla.)

---

## Lesson 3 — docs/learn-types.html

LESSON 03 OF 14 (ES: LECCIÓN 03 DE 14)

### Data types (ES: Tipos de datos)
Learn int, float, bool and string. (ES: Aprende int, float, bool y string.)

#### Four scalar types (ES: Cuatro tipos escalares)
int is a checked signed 64-bit integer; float is an IEEE binary64 decimal approximation; bool is true or false; string is immutable UTF-8 byte-oriented text. (ES: int es un entero de 64 bits con signo y comprobaciones; float es un número de coma flotante IEEE binary64; bool es true o false; string es texto inmutable orientado a bytes UTF-8.)

```tripol
age = 16
price = 2.5
active = true
name = "TRIPOL"
show age, price, active, name
```

Expected output (ES: Resultado esperado)

```output
16 2.500000 true TRIPOL
```

#### Explicit conversion (ES: Conversión explícita)
Integer and float arithmetic require matching types. Use float(n) to convert an integer, text(n) to format an integer, and parse_int(s) for valid integer text. (ES: La aritmética de int y float exige tipos iguales. Usa float(n) para convertir un entero, text(n) para convertirlo en texto y parse_int(s) para interpretar texto entero válido.)

```tripol
show float(9) / 2.0
show "Score: " + text(42)
show parse_int("42") + 1
```

Expected output (ES: Resultado esperado)

```output
4.500000
Score: 42
43
```

Important (ES: Importante)For untrusted integer text use errors.int_or or errors.is_int; invalid parse_int input may fail. (ES: Para texto entero no comprobado utiliza errors.int_or o errors.is_int; parse_int puede fallar con datos inválidos.)

Safe integer conversion (ES: Conversión segura de enteros)

#### Other data types (ES: Otros tipos de datos)
TRIPOL also supports fixed integer arrays, runtime-sized ints, typed growable lists and named structures. Learn these in the following lessons. (ES: TRIPOL también admite arreglos enteros fijos, ints de tamaño determinado en ejecución, listas dinámicas tipadas y estructuras con nombre. Aprende estos tipos en las siguientes lecciones.)

```tripol
numbers = [10, 20, 30]
show numbers[0]
```

Expected output (ES: Resultado esperado)

```output
10
```

---

## Lesson 4 — docs/learn-operators.html

LESSON 04 OF 14 (ES: LECCIÓN 04 DE 14)

### Operators (ES: Operadores)
Arithmetic, comparisons, logic and precedence. (ES: Aritmética, comparaciones, lógica y prioridad.)

#### Arithmetic (ES: Aritmética)
Use + for addition, - for subtraction or negation, * for multiplication, / for division, and % for integer remainder. These operators need no library import. (ES: Utiliza + para sumar, - para restar o negar, * para multiplicar, / para dividir y % para el residuo entero. No necesitan importar una biblioteca.)

```tripol
show 7 + 3, 7 - 3, 7 * 3
show 7 / 3, 7 % 3
```

Expected output (ES: Resultado esperado)

```output
10 4 21
2 1
```

#### Signed division and remainder (ES: División y residuo con negativos)
Integer division truncates toward zero. The remainder follows the sign of the dividend (or is zero). Use floats for decimal division. (ES: La división entera trunca hacia cero. El residuo tiene el signo del dividendo (o es cero). Usa float para una división decimal.)

```tripol
show -7 / 3, -7 % 3
show 7 / -3, 7 % -3
show -7 / -3, -7 % -3
```

Expected output (ES: Resultado esperado)

```output
-2 -1
-2 1
2 -1
```

```tripol
show 7.0 / 3.0
show float(7) / 3.0
```

Expected output (ES: Resultado esperado)

```output
2.333333
2.333333
```

Important (ES: Importante)Division or remainder by zero and signed integer overflow are checked errors. (ES: La división o residuo entre cero y el desbordamiento de enteros con signo son errores comprobados.)

#### Comparison and logic (ES: Comparaciones y lógica)
Compare using ==, !=, <, >, <= and >=. Combine boolean expressions with and, or, and not; and/or short-circuit. (ES: Compara con ==, !=, <, >, <= y >=. Combina expresiones booleanas con and, or y not; and/or usan cortocircuito.)

```tripol
show 7 == 7, 7 != 3
show 7 > 3, 7 < 3
show 7 >= 7, 3 <= 7
show true and false, true or false, not true
```

Expected output (ES: Resultado esperado)

```output
true true
true false
true true
false true false
```

#### Precedence and updates (ES: Prioridad y asignaciones)
Multiplication, division and remainder take precedence over addition and subtraction. Parentheses override the order. +=, -=, *=, /= and %= modify an existing mutable integer variable. (ES: Multiplicación, división y residuo tienen prioridad sobre suma y resta. Los paréntesis cambian el orden. +=, -=, *=, /= y %= modifican una variable entera mutable.)

```tripol
show 2 + 3 * 4
show (2 + 3) * 4
```

Expected output (ES: Resultado esperado)

```output
14
20
```

```tripol
n = 10
n += 5
n -= 2
n *= 3
n /= 2
n %= 4
show n
```

Expected output (ES: Resultado esperado)

```output
3
```

Important (ES: Importante)For powers, roots and percentages see the math library. (ES: Para potencias, raíces y porcentajes consulta la biblioteca math.)

Complete math reference (ES: Referencia completa de matemáticas)

---

## Lesson 5 — docs/learn-conditions.html

LESSON 05 OF 14 (ES: LECCIÓN 05 DE 14)

### Conditionals (ES: Condicionales)
Choose what your program does with if, else if and else. (ES: Decide qué hace tu programa con if, else if y else.)

#### if (ES: if)
An if block executes only when its boolean condition is true. Use braces for the block. (ES: Un bloque if se ejecuta únicamente cuando su condición booleana es true. Usa llaves para delimitarlo.)

```tripol
score = 80
if score >= 60 {
    show "Passed"
}
```

Expected output (ES: Resultado esperado)

```output
Passed
```

#### else and else if (ES: else y else if)
else runs when the earlier condition is false. else if checks another condition; the first matching branch wins. (ES: else se ejecuta si la condición anterior es falsa. else if comprueba otra condición; se ejecuta la primera rama que coincide.)

```tripol
score = 75
if score >= 90 {
    show "Excellent"
} else if score >= 60 {
    show "Passed"
} else {
    show "Try again"
}
```

Expected output (ES: Resultado esperado)

```output
Passed
```

#### Combine conditions (ES: Combinar condiciones)
Use and when both conditions must be true, or when at least one must be true, and not to reverse a boolean. (ES: Usa and cuando ambas condiciones deben ser verdaderas, or cuando al menos una debe ser verdadera y not para invertir un booleano.)

```tripol
age = 16
allowed = true
if age >= 13 and allowed {
    show "Welcome"
} else {
    show "Unavailable"
}
```

Expected output (ES: Resultado esperado)

```output
Welcome
```

---

## Lesson 6 — docs/learn-loops.html

LESSON 06 OF 14 (ES: LECCIÓN 06 DE 14)

### Loops (ES: Bucles)
Repeat work with for, while and repeat. (ES: Repite tareas con for, while y repeat.)

#### for and range (ES: for y range)
range(end) starts at zero; range(start,end) starts at start. The end is excluded. for also iterates collection values. (ES: range(end) comienza en cero; range(start,end) comienza en start. El límite final se excluye. for también recorre valores de colecciones.)

```tripol
for i in range(1, 4) {
    show i
}
```

Expected output (ES: Resultado esperado)

```output
1
2
3
```

```tripol
numbers = [10, 20, 30]
for i, value in numbers {
    show i, value
}
```

Expected output (ES: Resultado esperado)

```output
0 10
1 20
2 30
```

#### while (ES: while)
while repeats while its boolean condition is true. Update the condition's variables so a finite loop can finish. (ES: while repite mientras la condición booleana sea verdadera. Actualiza las variables de la condición para que un bucle finito pueda terminar.)

```tripol
n = 0
while n < 3 {
    n += 1
}
show n
```

Expected output (ES: Resultado esperado)

```output
3
```

#### repeat (ES: repeat)
repeat executes its body a specified number of times. (ES: repeat ejecuta su bloque una cantidad determinada de veces.)

```tripol
total = 0
repeat 3 {
    total += 2
}
show total
```

Expected output (ES: Resultado esperado)

```output
6
```

#### break and continue (ES: break y continue)
break exits the current loop; continue skips the rest of the current iteration. (ES: break termina el bucle actual; continue salta el resto de la iteración actual.)

```tripol
for i in range(5) {
    if i == 2 {continue}
    if i == 4 {break}
    show i
}
```

Expected output (ES: Resultado esperado)

```output
0
1
3
```

---

## Lesson 7 — docs/learn-functions.html

LESSON 07 OF 14 (ES: LECCIÓN 07 DE 14)

### Functions (ES: Funciones)
Create reusable code with parameters, returns and defaults. (ES: Crea código reutilizable con parámetros, retornos y valores predeterminados.)

#### Define and call a function (ES: Definir y llamar una función)
func declares a named function. Parameters receive arguments, -> specifies a return type, and return sends a value back. Declare functions before the program statements in these examples. (ES: func declara una función con nombre. Los parámetros reciben argumentos, -> indica el tipo de retorno y return devuelve un valor. En estos ejemplos declara las funciones antes de las instrucciones del programa.)

```tripol
func add(a:int, b:int)->int {
    return a + b
}
show add(4, 5)
```

Expected output (ES: Resultado esperado)

```output
9
```

#### Optional parameters (ES: Parámetros opcionales)
Trailing scalar parameters may have constant default values. Required parameters come first; omit optional arguments to use their defaults. Calls are positional. (ES: Los parámetros escalares finales pueden tener valores predeterminados constantes. Los obligatorios van primero; omite los opcionales para usar sus valores predeterminados. Las llamadas son posicionales.)

```tripol
func greet(name="Guest")->string {
    return "Hello " + name
}
show greet()
show greet("Ada")
```

Expected output (ES: Resultado esperado)

```output
Hello Guest
Hello Ada
```

#### Short functions and inference (ES: Funciones cortas e inferencia)
Use => for a single returning expression. Current TRIPOL can infer simple scalar signatures, but explicit annotations are useful for a stable contract. Functions have no globals or closures. (ES: Usa => para una única expresión de retorno. TRIPOL puede inferir firmas escalares sencillas, pero los tipos explícitos son útiles para definir un contrato estable. Las funciones no tienen globales ni cierres.)

```tripol
func twice(n:int)->int => n * 2
show twice(21)
```

Expected output (ES: Resultado esperado)

```output
42
```

---

## Lesson 8 — docs/learn-arrays.html

LESSON 08 OF 14 (ES: LECCIÓN 08 DE 14)

### Arrays (ES: Arreglos)
Create, index, change and iterate integer arrays. (ES: Crea, consulta, modifica y recorre arreglos de enteros.)

#### Fixed arrays (ES: Arreglos fijos)
A literal like [10,20,30] creates a fixed integer array. Indexes begin at zero; len counts elements. Its length cannot grow. (ES: Un literal como [10,20,30] crea un arreglo entero fijo. Los índices comienzan en cero; len cuenta los elementos. Su longitud no puede crecer.)

```tripol
a = [10, 20, 30]
show a[0], len(a)
a[1] += 5
show a
```

Expected output (ES: Resultado esperado)

```output
10 3
[10, 25, 30]
```

Important (ES: Importante)Accessing an index outside the array is an error. (ES: Acceder a un índice fuera del arreglo es un error.)

#### Runtime-sized arrays (ES: Arreglos de tamaño en ejecución)
ints(n) creates an integer array of n zero-initialized elements. Its length is chosen at runtime but remains fixed. Assignment of dynamic arrays shares their backing reference. (ES: ints(n) crea un arreglo entero de n elementos inicializados en cero. Su longitud se decide en ejecución, pero permanece fija. La asignación de arreglos dinámicos comparte su referencia.)

```tripol
a = ints(3)
for i in range(3) {
    a[i] = i * 2
}
show a, sum(a)
```

Expected output (ES: Resultado esperado)

```output
[0, 2, 4] 6
```

#### Array utilities (ES: Utilidades de arreglos)
Import arrays or collections for operations such as sum, copying and appending into a new array. Mutating library functions generally require a dynamic array. (ES: Importa arrays o collections para operaciones como sumar, copiar o añadir a un arreglo nuevo. Las funciones de biblioteca que modifican el arreglo generalmente requieren uno dinámico.)

```tripol
import tripol.collections

a = [1, 2]
b = collections.append(a, 3)
show a
show b
```

Expected output (ES: Resultado esperado)

```output
[1, 2]
[1, 2, 3]
```

Important (ES: Importante)For a collection that grows in place through a handle, use list(). (ES: Para una colección que crece mediante una referencia estable, utiliza list().)

Array library (ES: Biblioteca arrays) · Collection library (ES: Biblioteca collections)

---

## Lesson 9 — docs/learn-lists.html

LESSON 09 OF 14 (ES: LECCIÓN 09 DE 14)

### Lists (ES: Listas)
Use growable, homogeneous collections. (ES: Utiliza colecciones dinámicas de un mismo tipo.)

#### Create and grow a list (ES: Crear y ampliar una lista)
list() creates an empty integer list. list([values]) initializes a list from homogeneous scalar values. add appends an element and returns the new size. (ES: list() crea una lista de enteros vacía. list([valores]) inicializa una lista con valores escalares del mismo tipo. add añade un elemento y devuelve el nuevo tamaño.)

```tripol
numbers = list([10, 20, 30])
numbers.add(40)
show numbers
show numbers.get(3), len(numbers)
```

Expected output (ES: Resultado esperado)

```output
[10, 20, 30, 40]
40 4
```

#### Read, update and remove (ES: Leer, modificar y eliminar)
Use zero-based indexing, get, set and remove_at. remove_at returns the removed element. clear resets the logical length; size and capacity report current size and reserved capacity. (ES: Usa índices desde cero, get, set y remove_at. remove_at devuelve el elemento eliminado. clear reinicia la longitud lógica; size y capacity indican el tamaño actual y la capacidad reservada.)

```tripol
numbers = list([10, 20, 30])
numbers[0] += 5
show numbers.remove_at(1)
show numbers
```

Expected output (ES: Resultado esperado)

```output
20
[15, 30]
```

#### Typed lists (ES: Listas tipadas)
Nonempty homogeneous initializers infer string, boolean or float lists as well. For empty typed lists, use string_list(), bool_list() or float_list(). Mixed element types are rejected. (ES: Los inicializadores no vacíos y homogéneos también permiten inferir listas de texto, booleanos o decimales. Para listas tipadas vacías utiliza string_list(), bool_list() o float_list(). No se pueden mezclar tipos.)

```tripol
names = list(["Ada", "Alex"])
names.add("Sam")
show names[2]
```

Expected output (ES: Resultado esperado)

```output
Sam
```

Important (ES: Importante)Lists are distinct from ints. Array-only library functions do not accept lists. Aliases share one list handle; integer, boolean and string lists have contains, remove and sort; float lists do not have these methods in R1. (ES: Las listas son distintas de ints. Las funciones exclusivas de arreglos no aceptan listas. Los alias comparten una referencia; las listas de enteros, booleanos y textos tienen contains, remove y sort; las de float no incluyen esos métodos en R1.)

---

## Lesson 10 — docs/learn-strings.html

LESSON 10 OF 14 (ES: LECCIÓN 10 DE 14)

### Strings (ES: Strings)
Build, inspect and interpolate immutable text. (ES: Crea, consulta e interpola textos inmutables.)

#### Text and concatenation (ES: Texto y concatenación)
Write text inside double quotes and join strings with +. Use escape sequences such as \n and \" for a newline or quotation mark. (ES: Escribe texto entre comillas dobles y une cadenas con +. Utiliza secuencias de escape como \n y \" para saltos de línea o comillas.)

```tripol
name = "TRIPOL"
show "Hello " + name
```

Expected output (ES: Resultado esperado)

```output
Hello TRIPOL
```

```tripol
show "Line 1\nLine 2"
```

Expected output (ES: Resultado esperado)

```output
Line 1
Line 2
```

#### Interpolation (ES: Interpolación)
An f-prefixed string inserts int, bool or string expressions between braces. Double braces escape literal braces. Float interpolation is not currently supported. (ES: Una cadena con prefijo f inserta expresiones int, bool o string entre llaves. Las llaves dobles representan llaves literales. Actualmente no se admite interpolar float.)

```tripol
name = "Alex"
score = 42
show f"{name}: {score}"
show f"Total: {2 + 3}"
```

Expected output (ES: Resultado esperado)

```output
Alex: 42
Total: 5
```

#### Bytes and text helpers (ES: Bytes y utilidades de texto)
len measures UTF-8 bytes rather than visible characters. The text library provides byte slicing, prefix/suffix tests and substring search. Slicing inside a multibyte character can produce invalid text. (ES: len mide bytes UTF-8 en vez de caracteres visibles. La biblioteca text ofrece recortes por bytes, comprobaciones de prefijos y sufijos y búsqueda de fragmentos. Recortar dentro de un carácter multibyte puede producir texto inválido.)

```tripol
import tripol.text

show len("á")
show text.part("TRIPOL", 0, 3)
show text.contains("TRIPOL", "POL")
```

Expected output (ES: Resultado esperado)

```output
2
TRI
true
```

Important (ES: Importante)Strings are immutable; create a new value instead of changing an individual byte. (ES: Las cadenas son inmutables; crea un nuevo valor en lugar de modificar un byte individual.)

Text library (ES: Biblioteca text)

---

## Lesson 11 — docs/learn-structs.html

LESSON 11 OF 14 (ES: LECCIÓN 11 DE 14)

### Structures (ES: Estructuras)
Group named fields with struct. (ES: Agrupa campos con nombre mediante struct.)

#### Declare a structure (ES: Declarar una estructura)
A struct describes named, typed fields. Construct an instance by calling its name with positional field values. Use dot notation to read and write fields. (ES: Un struct describe campos con nombre y tipo. Crea una instancia llamando a su nombre con los valores en orden. Utiliza un punto para leer y escribir campos.)

```tripol
struct Player {name:string;score:int}
p = Player("Alex", 10)
p.score = 20
show p.name, p.score
```

Expected output (ES: Resultado esperado)

```output
Alex 20
```

#### Default fields (ES: Campos predeterminados)
Trailing scalar fields may have constant defaults. You can omit them when calling the constructor. (ES: Los campos escalares finales pueden tener valores predeterminados constantes. Puedes omitirlos al llamar al constructor.)

```tripol
struct Player {name:string="Alex";score:int=0}
p = Player()
p.score = 100
show f"{p.name}: {p.score}"
```

Expected output (ES: Resultado esperado)

```output
Alex: 100
```

#### Reference semantics (ES: Semántica de referencia)
Assigning a structure to another variable shares the same instance. Fields can be scalar values or ints; nested named records and methods are not supported in this version. (ES: Asignar una estructura a otra variable comparte la misma instancia. Los campos pueden contener escalares o ints; esta versión no admite estructuras con nombre anidadas ni métodos.)

```tripol
struct Counter {value:int}
a = Counter(1)
b = a
b.value = 7
show a.value
```

Expected output (ES: Resultado esperado)

```output
7
```

---

## Lesson 12 — docs/learn-exceptions.html

LESSON 12 OF 14 (ES: LECCIÓN 12 DE 14)

### Error handling (ES: Manejo de errores)
Recover from selected failures with try, catch, finally and throw. (ES: Recupérate de fallos seleccionados con try, catch, finally y throw.)

#### Catch an arithmetic error (ES: Capturar un error aritmético)
On the native compiler, try runs protected code and catch receives the numeric error code if a recoverable failure occurs. Arithmetic errors use code 10; bounds errors use 11; failed assertions use 12. (ES: En el compilador nativo, try ejecuta código protegido y catch recibe el código numérico si ocurre un fallo recuperable. Los errores aritméticos usan 10; los de índices 11; las aserciones fallidas 12.)
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```tripol
a = 10
b = 0
try {
    show a / b
} catch code {
    show code
}
```

Expected output (ES: Resultado esperado)

```output
10
```

#### finally and throw (ES: finally y throw)
finally is optional and runs after success or recoverable unwinding. throw raises an application exception; use codes 100–255. Values outside that range normalize to 100. (ES: finally es opcional y se ejecuta tras el éxito o durante una recuperación. throw genera una excepción de aplicación; utiliza códigos de 100 a 255. Los valores fuera del rango se normalizan a 100.)
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```tripol
try {
    throw 101
} catch code {
    show code
} finally {
    show "Finished"
}
```

Expected output (ES: Resultado esperado)

```output
101
Finished
```

#### Prevent errors (ES: Prevenir errores)
Check values before dangerous operations. The errors library offers fallback helpers that avoid selected arithmetic and parsing failures without using exceptions. (ES: Comprueba los valores antes de operaciones peligrosas. La biblioteca errors ofrece funciones con valores alternativos que evitan algunos fallos aritméticos y de conversión sin utilizar excepciones.)

```tripol
import tripol.errors

show errors.div_or(10, 0, -1)
show errors.int_or("abc", 7)
```

Expected output (ES: Resultado esperado)

```output
-1
7
```

Important (ES: Importante)try/catch/finally/throw are native features not yet supported by Web Engine R3. Native I/O, allocation and OS faults are not generally recoverable. return inside protected blocks is rejected. (ES: try/catch/finally/throw son funciones nativas aún no disponibles en Web Engine R3. Los fallos nativos de E/S, asignación de memoria y sistema operativo no suelen ser recuperables. return dentro de bloques protegidos se rechaza.)

Technical error guide (ES: Guía técnica de errores) · Error library (ES: Biblioteca errors)

---

## Lesson 13 — docs/learn-memory.html

LESSON 13 OF 14 (ES: LECCIÓN 13 DE 14)

### Memory (ES: Memoria)
Understand temporary allocation and scratch regions. (ES: Comprende las asignaciones temporales y las regiones scratch.)

#### The arena (ES: La arena de memoria)
Native TRIPOL stores ordinary arena allocations until program exit unless it can reclaim a proven temporary scope. The native arena reserves 256 MiB of virtual space; this is not total RAM or a promise of free physical memory. (ES: TRIPOL nativo conserva las asignaciones normales de la arena hasta terminar el programa, salvo cuando puede recuperar un ámbito temporal comprobado. La arena nativa reserva 256 MiB virtuales; no es la RAM total ni una garantía de memoria física disponible.)

```tripol
import tripol.memory

show memory.capacity()
show memory.used()
show memory.available()
```

Important (ES: Importante)The browser adapter reports a smaller logical budget, not actual process memory. Values depend on the target and program. (ES: El adaptador web informa un presupuesto lógico menor, no la memoria real del proceso. Los valores dependen del destino y del programa.)

#### Temporary memory with scratch (ES: Memoria temporal con scratch)
scratch saves and restores the allocation frontier. Allocate temporary arrays or strings inside it; only permitted outer scalar values can be updated and survive the block. (ES: scratch guarda y restaura el límite de asignación. Reserva arreglos o textos temporales dentro del bloque; solo los escalares externos permitidos pueden modificarse y sobrevivir al bloque.)
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```tripol
total = 0
repeat 3 {
    scratch {
        values = ints(100)
        values[0] = 7
        total += values[0]
    }
}
show total
```

Expected output (ES: Resultado esperado)

```output
21
```

#### Lifetime rules (ES: Reglas de duración)
Scratch locals cannot escape. Outer strings are read-only and outer aggregates cannot be accessed inside the region. Control flow cannot leave the protected region. Scratch does not undo file writes or other external side effects, and retained memory is not securely erased. (ES: Las variables locales de scratch no pueden escapar. Los textos externos son de solo lectura y no se puede acceder a agregados externos dentro de la región. El control de flujo no puede salir de la región protegida. scratch no deshace escrituras de archivos ni otros efectos externos, y la memoria retenida no se borra de forma segura.)
Important (ES: Importante)scratch lifetime semantics require native TRIPOL; Web Engine R3 does not implement full scratch parity. General garbage collection is not available. (ES: La semántica de duración de scratch requiere TRIPOL nativo; Web Engine R3 no ofrece paridad completa. No hay recolector de basura general.)

Memory library (ES: Biblioteca memory)

---

## Lesson 14 — docs/learn-modules.html

LESSON 14 OF 14 (ES: LECCIÓN 14 DE 14)

### Modules (ES: Módulos)
Import standard libraries and create your own reusable modules. (ES: Importa bibliotecas estándar y crea tus propios módulos reutilizables.)

#### Import a standard library (ES: Importar una biblioteca estándar)
Place import statements before declarations and program instructions. Standard library imports use tripol.name; call functions with their module prefix. (ES: Coloca import antes de las declaraciones e instrucciones del programa. Las bibliotecas estándar se importan con tripol.nombre; llama las funciones usando el prefijo del módulo.)

```tripol
import tripol.math

show math.square(7)
```

Expected output (ES: Resultado esperado)

```output
49
```

Standard library catalog (ES: Catálogo de bibliotecas)

#### Create a local library (ES: Crear una biblioteca local)
Create two files in the same project directory. Local .tri modules can contain function definitions and imports, not top-level executable statements. Imported local functions share the top-level namespace and are called without a prefix. (ES: Crea dos archivos en el mismo directorio del proyecto. Los módulos .tri locales pueden contener funciones e importaciones, no instrucciones ejecutables de nivel superior. Las funciones locales importadas comparten el espacio de nombres y se llaman sin prefijo.)
helpers.tri
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```tripol
func double(n:int)->int {
    return n * 2
}
```

main.tri
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```tripol
import "helpers.tri"

show double(21)
```

Expected output (ES: Resultado esperado)

```output
42
```

Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```cmd
tripol run main.tri
```

Important (ES: Importante)Paths are relative to the importing file and must stay inside the project. Local module examples need real project files; the browser compiler does not provide them yet. (ES: Las rutas son relativas al archivo que importa y deben permanecer dentro del proyecto. Los ejemplos de módulos locales necesitan archivos reales; el compilador web todavía no los proporciona.)

#### Package your code with Vault (ES: Empaquetar código con Vault)
Vault installs a local .tri library as a named package with a version and source hash in the lockfile. Use the helpers.tri file from the previous example and run these commands in the native project. (ES: Vault instala una biblioteca .tri local como paquete con nombre, versión y hash de código en el archivo de bloqueo. Utiliza helpers.tri del ejemplo anterior y ejecuta estos comandos en el proyecto nativo.)
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```cmd
tripol add helpers helpers.tri --version 1.0.0
tripol packages
tripol verify
```

main.tri
Native TRIPOL · local execution (ES: TRIPOL nativo · ejecución local)

```tripol
import vault.helpers

show helpers.double(21)
```

Expected output (ES: Resultado esperado)

```output
42
```

Important (ES: Importante)Vault is local and has no public remote package registry. Only use trusted code; keep the lockfile and local package cache for reproducible offline builds. (ES: Vault es local y no dispone de un registro público remoto. Utiliza únicamente código de confianza; conserva el archivo de bloqueo y la caché local para compilaciones reproducibles sin conexión.)

Build and package reference (ES: Referencia de compilación y paquetes)

---

End of TRIPOL v0.28.0 AI reference.
