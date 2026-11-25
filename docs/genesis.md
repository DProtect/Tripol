# GENESIS R2 — bootstrap compiler

The production compiler remains Go. GENESIS is a native executable built from TRIPOL source in compiler/. It reads one line of source and emits Linux x86-64 ELF bytes as hex. Python only decodes the hex transport and writes a new file.

On Linux x86-64:

```sh
tripol build compiler/main.tri --target linux -o genesis
printf 'a=20;b=30;show a*b\n' | ./genesis | python3 compiler/unpack.py demo
./demo
```

Expected output: 600. The assignment mode emits native machine instructions for integer variables and arithmetic evaluated when the generated binary runs. The older single constant show expression mode remains supported, including division and remainder.

## R2 subset

- One line with semicolon-separated assignments followed by one show expression, at most 32 integer variables; ASCII identifiers up to 32 bytes.
- Runtime assignment expressions support integer literals, known variables, parentheses, unary signs, +, - and *. Overflow exits with status 10.
- Source lexer limit 4096 bytes, subject to the input intrinsic's own limit. No functions, control flow, arbitrary input, runtime division, Windows PE output or self-compilation in R2.
- Windows can host GENESIS but emitted programs target Linux. The unpack script does not compile anything.
- No direct binary file I/O yet; the compiler itself still needs the Go implementation for bootstrapping.

Next: runtime division, binary file I/O, typed statements, control flow, functions, PE output and multi-generation self-compilation.
