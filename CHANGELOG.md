# Changelog

## 0.12.0 — SENTINEL (experimental)

- Native recoverable try/catch/finally and explicit throw, with restricted control flow.
- Optional `--trace` language-frame trace; spelling suggestions in text and JSON diagnostics.
- Explicit `scratch` arena regions with escape restrictions and recovery cleanup.
- SSE2 length scanning; aligned loads do not cross page boundaries.
- JSON grammar validation and top-level raw value extraction.
- Windows modal form, WinHTTP HTTPS GET, and public Discord gateway discovery helper (PE inspected; Windows runtime testing pending).
- 21 embedded standard libraries; syntax highlighting updated; documentation reorganized.

Earlier releases and their full history remain in Git. v0.11 — VISION & GUARD introduced initial diagnostics and dialogs; v0.10 — CONNECT native TCP/HTTP; v0.9 — ECOSYSTEM Vault; v0.8 — ATLAS standard libraries. The original bootstrap was replaced by a native code emitter in v0.2.
