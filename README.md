<div align="center">

## v0.29.0 DIPO HTTP 2 — development

Typed HTTP requests, connection pooling, validated HTTPS, streamed files, background requests and deferred Discord replies. The development compiler builds for Windows and Linux x86-64. [HTTP guide](docs/http.md) · [Validation and status](docs/http-status.md)

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
