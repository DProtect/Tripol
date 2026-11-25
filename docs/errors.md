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
