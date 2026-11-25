# GENESIS R3 — QUANTUM X

The production Go compiler builds a native stage written in TRIPOL. GENESIS emits Linux x86-64 ELF, not Windows PE. It now supports integer functions (up to 8 int parameters), local integer variables, recursion, calls to previously declared functions, if/else, while, for/range, checked + - * / %, and repeated show. A function must end in a return; an earlier return in a branch is permitted. Forward calls, mutual recursion, global variable capture, bool expressions and a complete type system are not supported. Up to 32 function declarations, 32 slots per frame, 4096 source bytes, 16 statement nesting levels. Declarations must precede control flow in each frame. Loop counters are immutable and scoped.

## Direct binary output

```sh
tripol build compiler/direct.tri --target linux -o genesis-stage
printf 'examples/projects/genesis_functions.tri\ngenesis-demo\n' | ./genesis-stage
chmod +x genesis-demo
./genesis-demo
```

On Windows build the stage with `--target windows -o genesis-stage.exe` then run it and enter the source and output paths on separate lines. The stage prints bytes written and produces a Linux ELF; execute on Linux x86-64. The direct driver replaces the specified output file, so choose a new path. The existing hex transport (`compiler/file.tri` with `tripol genesis`) remains available and refuses overwrite.

Expected output of genesis_functions.tri: 0, 5, 10, ..., 45, each on its own line. Full self-compilation, PE generation and additional language features remain future work.
