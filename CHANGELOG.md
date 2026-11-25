# v0.27.1 QUANTUM X (development)

- GENESIS PE R1 console output and native Windows executable generation.
- GENESIS integer functions, direct calls, range loops, repeated show, native ELF output without host byte transport via compiler/direct.tri.
- Bounded native binary file operations, byte cursor, limited regex matching.
- Growable integer lists with stable aliases and checked methods.
- Versioned follow-up to published v0.27.0; old release assets are preserved.

# v0.27 QUANTUM X (development)

- Experimental GENESIS R3: checked division/remainder, native conditional branches and loops, multiline file driver; Go CLI transports ELF bytes without Python.
- Checked compound indexed assignments, including single evaluation of the index and recoverable arithmetic/bounds failures.
- 55 standard modules: heap, utf8 and csv added; new reuse stress example.
- UPDATE R2 first stage: download progress and verified receipt history. Manual installation remains required.
- TIR R6 first-stage control-flow graph inspection; existing native backend retained.
- Editor completion, hover, local definitions and saved-file compiler commands; software-rendered game example.
- Windows Validation Center, new native regression tests and release workflow.
- No claim of full SSA, general reclamation, self-hosting, PE generation by GENESIS, GPU acceleration, automatic installation/rollback or benchmark gains.

# v0.26 OVERDRIVE (development)

- UPDATE R1: opt-out interactive release notices, daily cache, explicit checks and verified downloads. No background installation.
- GENESIS R2: semicolon-separated integer assignments, symbol table, checked native +, - and * on Linux x86-64; retains the R1 constant-expression mode. No full self-hosting, PE emitter or direct binary I/O yet.
- 52 modules: new bitset and profiler. Mutable collection operations now require dynamic backing arrays instead of borrowed fixed arrays.
- New tests for release selection, caching, download tampering, bitset and runtime ELF generation. No benchmark improvements claimed before measurement.

# v0.25 EVOLUTION (development)

- GENESIS R1: TRIPOL-written frontend and ELF x86-64 writer for a single constant integer show statement; deterministic hex transport, Linux execution and Windows-host structural test coverage proposed.
- Bounded reusable integer stack and FIFO ring queue: no new storage after creation during valid operations.
- TIR R5: eliminate pure constant left identities without removing checked right-hand arithmetic.
- Version metadata, VS Code grammar and release workflow updated. Full self-hosting, general memory reclamation, SSA and expanded syntax remain future work.
- No new performance numbers claimed before measurement.

# v0.24 ASCENSION

- Native allocation-free array copy into an existing buffer, including self-copy and size mismatch checks.
- Arena: 120-second default timeout, configurable 1–900 seconds, single retry of Rust version probe on timeout, copy-reuse workload and p95 report.
- Monochrome responsive static website, bilingual interface, searchable library catalog and checksum-verified user installers.
- Discord unchanged. General optimizer and memory reclamation work remains planned.

# ORIGIN v0.23.0

- Discord Gateway replies route automatically to the received channel; optional channel filter remains. `send_to`, `event_channel`, `all_channels`, `has_channel` added. No reconnect, slash commands or live Gateway validation yet.
- `tripol.binary`: eight validated byte-array and endian conversion functions.
- AVX2 search checks two vectors per loop while preserving first match and scalar fallback.
- Public release preparation: CI configuration, security reporting, conduct, contribution templates, open-core boundaries.
- Linux tests: 139 groups, 371 subtests. Windows PE cross-built; Windows CI not executed locally.

# Changelog

## 0.22.0 — NEXUS

- Experimental native Windows Gateway (TLS WebSocket, bounded receive worker, Hello/Identify, sequence/heartbeat/ACK, READY and MESSAGE_CREATE). `run` now uses Gateway; `run_rest` is explicit fallback. Stop on disconnect, reconnect request, invalid session or missed ACK; no Resume or automatic reconnect. Offline protocol/PE tests only, live Windows verification pending.
- Easier Discord API: `setup`, `welcome`, `activity`, `status`, `ping`, `connected`; 30 callable Discord functions including protocol helpers. `json.get_raw_large` supports Gateway payloads up to 65535 bytes.
- AVX2 first-match int64 search with CPUID and XGETBV checks, scalar fallback and tail. 21-sample Linux search median 0.574797 → 0.156361 ms.
- 47 libraries, 284 callable functions; 137 top-level test groups and 367 subtests passed on Linux. No Rust toolchain or native Windows execution in release environment.

## 0.21.0 — ECLIPSE

- Discord R3: `login`, `login_env`, `channel`, `say`, `reply`, `run`, `reset`, `configured`, `preview`, `dispatch` and parsing helper. Up to 16 exact replies, immutable configuration snapshots in bounded process-local storage, scratch reuse for polling and registration. Existing Discord R2 API retained. No Gateway, live connection or presence validated.
- ArrayCore R2: native independent `arrays.copy` and bounded `search.linear` / `arrays.contains`, using `rep movsq` / `repne scasq`; no unchecked arithmetic introduced.
- Arena: copy and search cases, original four unchanged. Linux 21-sample comparison: copy 4.299 → 2.865 ms, search 1.021 → 0.581 ms.
- 47 libraries, 273 callable functions, including 20 Discord functions (helpers included).
- General memory reclamation, AVX2, WebSocket, native Windows execution CI and full SSA remain pending.

## 0.20.0 — OMEGA

- ArrayCore R1: checked pointer traversal and four-way sequential unrolling of verified zero-based dynamic array reductions in Balanced/Turbo. Overflow order and bounds errors preserved, including overflow before an eventual bounds failure. Guards, tracing and unsupported loops use generic code.
- VectorCore R1: `tripol.simd.equal`, SSE2 exact equality of two int64 arrays. No automatic AVX2 yet.
- Discord R2: Windows-only native HTTPS authenticated GET and POST to fixed discord.com host; one-channel REST polling, environment token, safe mention default, fail-stop on HTTP errors including rate limits. Template `tripol new NAME --template discord`. Not a Gateway client; no live Windows or Discord connection verified. See docs/discord.md.
- `tripol.env.get` (Windows); 47 modules and 262 callable functions (including helpers). 127 groups and 347 subtests pass on Linux; Windows cross-build and import inspection only.
- Reproducible Arena comparisons: array reduction 5.32x faster, other workloads roughly unchanged. No general compiler speed gain or fresh Rust result claimed.

## 0.19.0 — APEX

- RegisterFlow R3 (restricted): private register argument convention for validated two-call integer recurrences with a literal base threshold. Preserves recursive call order and checked combining operation; no benchmark-name special cases. Generic ABI for guards, tracing and unsupported forms.
- DivisionCore: hoist one total, invariant power-of-two operation in call-free Turbo loops; specialize nonnegative loop indexes; lower modulo-zero tests directly; combine adjacent quotient/remainder updates while preserving each overflow trap.
- Arena: checked array summation, `--json FILE` exports samples, medians and output verification. Existing three workloads unchanged.
- 122 top-level test groups / 334 first-level subtests in Linux. No Rust installed on release host; Windows compiler/installer cross-built, native execution pending.
- No new syntax or libraries. Full CFG/SSA, SIMD, general memory reclamation and additional platforms remain future work.

## 0.18.0 — TITAN

- RegisterFlow R2: up to four outer integer accumulators, call-free conditional branches and loop control, with conservative fallback and explicit register save/restore.
- Shorter checked signed power-of-two division and remainder lowering; 150 additional deterministic random signed cases.
- Built-in `tripol arena` runs three paired TRIPOL/Rust workloads with checked i64 arithmetic, runtime inputs, output validation, alternating samples and timeouts. No Python/Cargo requirement. Rust execution not tested here because rustc is unavailable. TRIPOL-only mode verified.
- 45 modules, 253 callable library functions; 119 top-level test groups and 310 subtests pass on Linux. Windows cross-build only. Full SSA, AVX2, general memory reclamation and incremental compilation remain pending.

## 0.17.0 — HYPER

- Register-backed straight-line integer range loops: index, bound, one outer accumulator and one local scalar; direct checked updates. Conservative fallback around calls, nested flow, guards and tracing.
- Private R8–R11 argument ABI for eligible pure int64 leaf functions; register-based expression lowering and selective spills; balanced and turbo profiles.
- High-resolution monotonic nanosecond timing and explicit internal benchmark markers. Calibrated build batches and profile selection.
- 45 modules, 253 callable functions. 115 top-level test groups and 300 subtests; existing runtime test helper also exercised under balanced/turbo. 100 deterministic four-argument differential cases. Linux runtime tests and Windows cross-build/PE inspection only.
- Full control-flow TIR, general register allocation, SIMD AVX2 and memory reclamation remain pending.

## 0.16.0 — CORE

- Single-pass byte lexer, differential-tested against 0.15 over 20,008 inputs (including invalid syntax, UTF-8 and byte positions).
- O(1) average relocation-reference lookup and per-build scalar TIR cache. Full unoptimized source-body validation remains; discarded validation executable no longer linked.
- Checked arithmetic directly on stack memory operands, direct integer variable comparison branches, range counter increment proof and signed constant power-of-two division/remainder.
- 45 modules and 251 functions retained; 109 top-level tests and 288 subtests passed on Linux. Windows cross-build only.
- Reproducible compilation and runtime comparisons in benchmarks/core.py. No general GC, register calling convention or full TIR control-flow rewrite yet.

## 0.15.0 — NOVA

- 45 embedded modules, 251 callable functions.
- Dynamic int arrays up to 4,194,304 elements; growable buffers up to 4,194,303 values. Arena remains 256 MiB, no general reclamation.
- Software framebuffer up to 2048×2048, 16 MiB Win32 staging allocation, checked native horizontal span fills, clipped RGB sprites, line/circle drawing and opacity blending.
- Fixed-capacity string-to-int maps, camera transforms, particle helpers and testing assertions.
- TIR multi-argument scalar arithmetic operates directly on argument stack slots while preserving overflow checks.
- Self-hosted experimental integer expression parser/evaluator (not the production parser).
- 104 test groups, 288 subtests; Linux execution, Windows cross-build and PE inspection. No GPU, PNG, audio, multithreading, general GUI toolkit or full self-hosting.

## 0.14.0 — QUANTUM

- 40 standard modules, 225 callable functions (including library helpers).
- Native software framebuffer, Windows GDI game window and foreground key polling; Linux canvas rendering, Windows PE verification only.
- Demand-enabled 256 MiB arena reservation; scratch reuse and memory metrics.
- Growable integer buffer with explicit reassignment, physics, animation, vector, image and lexical path helpers.
- Short returning functions, experimental turbo profile and broader scalar TIR inlining.
- Non-mutating Guard Fix preview, Benchmark Arena and experimental self-hosted lexer.
- No general GC, full GUI toolkit, arbitrary sprites, audio, Discord bot or complete self-hosting.

## 0.13.0 — INFINITY (experimental)

- Explicit `f"...{name}..."` interpolation and scalar constant default parameters.
- Nine new standard modules: bits, checksum, collections, colors, compression, datetime, geometry, matrix and sequences; 30 modules and 164 callable functions total.
- Four additional Windows native message-box variants. No general desktop window/event toolkit yet.
- TIR R1 typed checked int64 expression subset and experimental small-function inlining; `tripol ir` inspection.
- Fold literal string lengths and remove safe arithmetic identities in nano builds.
- 89 test groups and 285 subtests passed in Linux; Windows binaries cross-built and PE checked.
- Existing arena, networking, Discord and exception restrictions remain.

## 0.12.0 — SENTINEL (experimental)

- Native recoverable try/catch/finally and explicit throw, with restricted control flow.
- Optional `--trace` language-frame trace; spelling suggestions in text and JSON diagnostics.
- Explicit `scratch` arena regions with escape restrictions and recovery cleanup.
- SSE2 length scanning; aligned loads do not cross page boundaries.
- JSON grammar validation and top-level raw value extraction.
- Windows modal form, WinHTTP HTTPS GET, and public Discord gateway discovery helper (PE inspected; Windows runtime testing pending).
- 21 embedded standard libraries; syntax highlighting updated; documentation reorganized.

Earlier releases and their full history remain in Git. v0.11 — VISION & GUARD introduced initial diagnostics and dialogs; v0.10 — CONNECT native TCP/HTTP; v0.9 — ECOSYSTEM Vault; v0.8 — ATLAS standard libraries. The original bootstrap was replaced by a native code emitter in v0.2.

Direct self tail calls with full argument lists reuse the frame in all profiles unless tracing. Other recursion remains subject to the OS stack limit.
