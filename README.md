<div align="center">

<img src="assets/tripol-emblem.png" alt="TRIPOL official logo" width="180" />

# TRIPOL — Official Website

**Simple by design. Native by nature.**

[Visit the website](https://tripollas.github.io/Tripol/) · [Documentation](https://tripollas.github.io/Tripol/docs/) · [Playground](https://tripollas.github.io/Tripol/playground/) · [Compiler source](https://github.com/Tripollas/Tripol/tree/main)

</div>

## About this branch

The `gh-pages` branch hosts the official static website for TRIPOL, a free experimental programming language that compiles directly to native x86-64 executables.

The compiler and standard library source code live in [`main`](https://github.com/Tripollas/Tripol/tree/main).

## QUANTUM X — v0.27.1 experimental release

- Official landing page with a monochrome glassmorphism design and the TRIPOL emblem.
- English by default, with Spanish navigation and library overviews. Code examples remain in English.
- Reference pages for all 64 standard libraries and 383 callable functions.
- Searchable documentation, benchmark results, install scripts and official GitHub Releases downloads, GENESIS R3, UPDATE R2 and new Bitset/Profiler guides.

## Downloads

| Platform | Distribution |
| --- | --- |
| Windows x64 | [Installer ZIP](https://github.com/Tripollas/Tripol/releases/download/v0.27.1/TRIPOL-v0.27-Installer.zip) |
| Windows x64 | [Setup EXE](https://github.com/Tripollas/Tripol/releases/download/v0.27.1/TRIPOL-v0.27-Setup-x64.exe) |
| Linux x64 | [Tarball](https://github.com/Tripollas/Tripol/releases/download/v0.27.1/TRIPOL-v0.27-Linux-x64.tar.gz) |
| Source | [Source ZIP](https://github.com/Tripollas/Tripol/archive/refs/tags/v0.27.1.zip) |
| VS Code | [VSIX](https://github.com/Tripollas/Tripol/releases/download/v0.27.1/tripol-language-0.27.1.vsix) |
| Integrity | [SHA-256 checksums](https://github.com/Tripollas/Tripol/releases/download/v0.27.1/SHA256SUMS.txt) |

QUANTUM X is a pre-release. The previous v0.24 downloads remain archived in downloads/. The Windows installer is unsigned. Review installation scripts and verify the checksums before running downloaded software.

## Playground R3 / WEB ENGINE R1 (experimental)

[Open the playground](https://tripollas.github.io/Tripol/playground/?v=3). 100 editable bilingual examples in eight categories, with search, syntax highlighting, completion, file open/download, share links, local draft, diagnostics and console. All 100 examples have verified output in the educational interpreter.

The optional Web Engine R1 uses the production Go parser and validation pass to emit actual WASM for a scalar subset (int64, bool, checked arithmetic, branches, loops, scalar functions and selected pure math functions). The browser compiler binary is built and tested by the main branch's Web Engine workflow. Until that build succeeds and its assets deploy, the interpreter remains available. No silent fallback from WASM mode. Generated modules require the TRIPOL host ABI, and are not standalone executables. Strings, arrays, input and general libraries are not supported by the WASM backend yet.

Run `node playground/test-r3.cjs` for the interpreter and editor regression suite. Use `go test ./core -run TestWebEngineR1` on main and the Web Engine workflow for actual WASM validation.

## Local development

Clone this branch and serve it locally:

```bash
git clone --branch gh-pages --single-branch https://github.com/Tripollas/Tripol.git
cd Tripol
python -m http.server 8000
```

Open http://localhost:8000. No Node.js or build step is required.

## Deployment

GitHub Pages serves the root of `gh-pages`. Keep `index.html`, `.nojekyll`, `assets/`, `docs/` and `downloads/` in place. Website changes belong in this branch; compiler changes belong in `main`.

## Contributing

Report issues through the [official repository](https://github.com/Tripollas/Tripol/issues). See the [contribution guide](https://github.com/Tripollas/Tripol/blob/main/CONTRIBUTING.md), [security policy](https://github.com/Tripollas/Tripol/blob/main/SECURITY.md) and [license](https://github.com/Tripollas/Tripol/blob/main/LICENSE).

## Source architecture

The production Go compiler is now in [`core/`](https://github.com/Tripollas/Tripol/tree/main/core); GENESIS is in `compiler/`. CLI integration tests and historical logs are in `tests/`. The root Go entry point remains compatible.


## AI-ready reference / Referencia para IA

[Download TRIPOL-AI.md](TRIPOL-AI.md) — complete v0.27.1 reference with 64 modules, 383 functions, 18 guides and 14 lessons. Attach it to an AI assistant to ask questions or write TRIPOL code.
