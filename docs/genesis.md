# GENESIS R3 — QUANTUM X

The production Go compiler builds a native stage written in TRIPOL. GENESIS emits Linux x86-64 ELF and experimental Windows PE32+ console executables. It now supports integer functions (up to 8 int parameters), local integer variables, recursion, forward and mutual calls, if/else-if/else, while, for/range, checked + - * / %, and repeated show. A function must end in a return; an earlier return in a branch is permitted. Global variable capture, a complete type system are not supported. Boolean conditions have short-circuit operators. Up to 32 function declarations, 32 slots per frame, 4096 source bytes, 16 statement nesting levels. Declarations must precede control flow in each frame. Loop counters are immutable and scoped.

## Direct binary output

```sh
tripol build compiler/direct.tri --target linux -o genesis-stage
printf 'examples/projects/genesis_functions.tri\ngenesis-demo\n' | ./genesis-stage
chmod +x genesis-demo
./genesis-demo
```

On Windows build the stage with `--target windows -o genesis-stage.exe` then run it and enter the source and output paths on separate lines. The stage prints bytes written and produces an ELF for Linux x86-64. Choose a new output path. Both the direct driver and the hex transport (`compiler/file.tri` with `tripol genesis`) reject existing destinations, including paths that resolve to the source file.

Expected output of genesis_functions.tri: 0, 5, 10, ..., 45, each on its own line. Full self-compilation and additional language features remain future work.

## Windows PE R1

```cmd
tripol build compiler/windows.tri --target windows -o genesis-win.exe
genesis-win.exe
```
Enter `examples/projects/genesis_functions.tri` and a new output path such as `genesis-demo.exe` on separate lines. Run `genesis-demo.exe`. Expected ten lines from 0 to 45 in steps of 5. The compiler emits PE32+ directly, with RX code and a RW import table for GetStdHandle, WriteFile and ExitProcess. Linux can host this stage as well. The generated programs use kernel32 and a fixed image base. General Windows APIs and full self-hosting remain future work.

### GENESIS R3 additional control flow

The declaration pass now builds the function signature table before emitting code. Calls target stable native jump stubs patched to actual function bodies, enabling forward calls and mutual recursion. Conditions support integer comparisons, boolean literals, `not`, `and`, `or` (also !, &&, ||), short-circuit evaluation and boolean grouping; `else if` chains are supported. Booleans are now first-class scalar values in local variables, parameters and function returns; `show` remains integer-only in this stage. Integer and boolean expressions now pass through a separate bounded AST; statements still lower directly from tokens, functions cannot access main's local slots, and self-compilation is not yet possible. See examples/projects/genesis_control.tri.

### Loop control

`break` exits the nearest `while` or `for` loop. `continue` rechecks a `while` condition or advances a `for` counter. Both are available inside nested conditionals. GENESIS rejects loop control outside a loop.

## Integer expression AST

`compiler/ast.tri` represents literals, variable references, unary negation, checked arithmetic, calls and linked argument nodes in a bounded integer array (1024 nodes per expression). Parsing is separate from native code emission. The backend preserves left-to-right argument evaluation, precedence and checked overflow/division. Existing mutable integer variables support `+=`, `-=`, `*=`, `/=`, `%=`; range counters remain immutable. `examples/projects/genesis_ast.tri` exercises the AST on both native targets. This is not yet a full program AST or a self-hosting compiler.

### Typed scalar expressions

GENESIS now infers int/bool local types and checks assignment, function parameters/returns, arithmetic, comparisons and conditions during compilation. Boolean equality and short-circuit operators use the same AST as arithmetic, including grouped expressions. Boolean values can be returned or passed to functions, but cannot be printed directly yet; use an integer-returning function or branch. Integer range counters remain immutable. `examples/projects/genesis_types.tri` demonstrates the two scalar types. Strings, arrays, records and full program AST remain pending.
