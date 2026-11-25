<div align="center">

<img src="assets/tripol-emblem.png" alt="TRIPOL official logo" width="180" />

# TRIPOL — Official Website

**Simple by design. Native by nature.**

[Visit the website](https://tripollas.github.io/Tripol/) · [Documentation](https://tripollas.github.io/Tripol/docs/) · [Compiler source](https://github.com/Tripollas/Tripol/tree/main)

</div>

## About this branch

The `gh-pages` branch hosts the official static website for TRIPOL, a free experimental programming language that compiles directly to native x86-64 executables.

The compiler and standard library source code live in [`main`](https://github.com/Tripollas/Tripol/tree/main).

## OVERDRIVE — v0.26.0

- Official landing page with a monochrome glassmorphism design and the TRIPOL emblem.
- English by default, with Spanish navigation and library overviews. Code examples remain in English.
- Reference pages for all 52 standard libraries and 319 callable functions.
- Searchable documentation, benchmark results, install scripts and official GitHub Releases downloads, GENESIS R2, UPDATE R1 and new Bitset/Profiler guides.

## Downloads

| Platform | Distribution |
| --- | --- |
| Windows x64 | [Installer ZIP](https://github.com/Tripollas/Tripol/releases/download/v0.26.0/TRIPOL-v0.26-Installer.zip) |
| Windows x64 | [Setup EXE](https://github.com/Tripollas/Tripol/releases/download/v0.26.0/TRIPOL-v0.26-Setup-x64.exe) |
| Linux x64 | [Tarball](https://github.com/Tripollas/Tripol/releases/download/v0.26.0/TRIPOL-v0.26-Linux-x64.tar.gz) |
| Source | [Source ZIP](https://github.com/Tripollas/Tripol/archive/refs/tags/v0.26.0.zip) |
| VS Code | [VSIX](https://github.com/Tripollas/Tripol/releases/download/v0.26.0/tripol-language-0.26.0.vsix) |
| Integrity | [SHA-256 checksums](https://github.com/Tripollas/Tripol/releases/download/v0.26.0/SHA256SUMS.txt) |

OVERDRIVE is a pre-release. The previous v0.24 downloads remain archived in downloads/. The Windows installer is unsigned. Review installation scripts and verify the checksums before running downloaded software.

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
