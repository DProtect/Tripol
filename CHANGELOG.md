# Changelog

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
