<div align="center">

<img src="assets/tripol-emblem.png" alt="TRIPOL official logo" width="180" />

# TRIPOL — Official Website

**Simple by design. Native by nature.**

[Visit the website](https://dprotect.github.io/Tripol/) · [Documentation](https://dprotect.github.io/Tripol/docs/) · [Compiler source](https://github.com/DProtect/Tripol/tree/main)

</div>

## About this branch

The `gh-pages` branch hosts the official static website for TRIPOL, a free experimental programming language that compiles directly to native x86-64 executables.

The compiler and standard library source code live in [`main`](https://github.com/DProtect/Tripol/tree/main).

## ASCENSION — v0.24.0

- Official landing page with a monochrome glassmorphism design and the TRIPOL emblem.
- English by default, with Spanish navigation and library overviews. Code examples remain in English.
- Reference pages for all 48 standard libraries and 297 callable functions.
- Searchable documentation, benchmark results, install scripts and downloadable distributions.

## Downloads

| Platform | Distribution |
| --- | --- |
| Windows x64 | [Installer ZIP](downloads/TRIPOL-v0.24-Installer.zip) |
| Windows x64 | [Portable ZIP](downloads/TRIPOL-v0.24-Windows-x64.zip) |
| Linux x64 | [Tarball](downloads/TRIPOL-v0.24-Linux-x64.tar.gz) |
| Source | [Source ZIP](downloads/TRIPOL-v0.24-source.zip) |
| Integrity | [SHA-256 checksums](downloads/TRIPOL-v0.24-SHA256.txt) |

The Windows installer is unsigned. Review installation scripts and verify the checksums before running downloaded software.

## Local development

Clone this branch and serve it locally:

```bash
git clone --branch gh-pages --single-branch https://github.com/DProtect/Tripol.git
cd Tripol
python -m http.server 8000
```

Open http://localhost:8000. No Node.js or build step is required.

## Deployment

GitHub Pages serves the root of `gh-pages`. Keep `index.html`, `.nojekyll`, `assets/`, `docs/` and `downloads/` in place. Website changes belong in this branch; compiler changes belong in `main`.

## Contributing

Report issues through the [official repository](https://github.com/DProtect/Tripol/issues). See the [contribution guide](https://github.com/DProtect/Tripol/blob/main/CONTRIBUTING.md), [security policy](https://github.com/DProtect/Tripol/blob/main/SECURITY.md) and [license](https://github.com/DProtect/Tripol/blob/main/LICENSE).
