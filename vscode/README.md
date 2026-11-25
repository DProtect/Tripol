# TRIPOL Language 0.27.1

Syntax highlighting, snippets, keyword/module completion, catalog-based member completion and hover, local function definition lookup, and explicit Check/Build/Run commands. The VSIX release workflow generates the full catalog from actual stdlib signatures. Set `tripol.compilerPath` if the compiler is not on PATH. Compiler commands require a trusted workspace and saved files; optional check-on-save is off by default. Run invokes the compiler explicitly, not automatically. This is not a full LSP, debugger or formatter. Run `node generate-catalog.js` in the extension folder for a current local catalog.

## QUANTUM X formatting

Format Document uses the installed `tripol fmt --stdin` command in a trusted workspace. The formatter preserves tokens, comments, line endings and interpolated strings; it adjusts indentation only. The native completion/hover/definition providers remain active. A standalone stdio language server is also available through `tripol lsp` for editors that can launch an external LSP server. It supports full-buffer changes, syntax/type diagnostics for self-contained and standard-library code, library completions, hover, same-document definitions, document symbols and formatting. Local/Vault import resolution, incremental changes and cross-file navigation are not yet implemented in this server.
