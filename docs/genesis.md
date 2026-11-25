# GENESIS R3 — bootstrap compiler

The production compiler remains Go. GENESIS is written in TRIPOL and emits Linux x86-64 ELF machine code. The file driver reads a source path on stdin and outputs hex transport; the Go CLI writes those bytes to a new file. No Python is needed for the new command, but standalone binary I/O in TRIPOL remains pending.

```sh
tripol build compiler/file.tri --target linux -o genesis-stage
tripol genesis ./genesis-stage examples/projects/genesis_r3.tri -o genesis-demo
./genesis-demo
```

Expected: 30. On Windows use `--target windows -o genesis-stage.exe` and `tripol genesis genesis-stage.exe examples/projects/genesis_r3.tri -o genesis-demo`. The emitted program still requires Linux x86-64. Use a new output file; overwrite is rejected.

R3 supports up to 32 predeclared integer variables, reassignment, multiline/semicolon separators, integer comparisons, nested if/else and while blocks, checked runtime + - * / % and one final show expression. Variable declarations are required before control flow, avoiding uninitialized slots; nesting limited to 16 statement levels and 64 expression levels, source max 4096 bytes subject to host file limits. Runtime divide-by-zero and signed overflow exit 10. The old stdin driver and R1 constant show mode remain available for compatibility.

No functions, for loops, boolean short-circuit expressions, general type system, AST/SSA backend, PE emission, direct binary writes or self-compilation yet. The host driver invokes an explicitly selected local stage executable; use a trusted stage. Full bootstrap parity and multi-generation verification remain future work.
