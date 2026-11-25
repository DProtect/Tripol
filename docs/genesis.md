# GENESIS R1 — bootstrap compiler

The production compiler remains the Go implementation. GENESIS is a separate, native TRIPOL executable built from `compiler/main.tri`. It uses the existing experimental lexer and arithmetic parser, with an ELF writer implemented in TRIPOL. Its entire accepted language is one line `show INTEGER_EXPRESSION` (decimal literals, parentheses, unary +/- and checked + - * / %). Arithmetic is evaluated at compile time; the generated binary prints the result and exits. It does not compile itself yet, nor support variables, functions or Windows output.

On Linux x86-64 with TRIPOL and Python 3 installed:

```sh
tripol build compiler/main.tri --target linux -o genesis
printf 'show (6+1)*6\n' | ./genesis | python3 compiler/unpack.py demo
./demo
```

The output should be 42. Python only decodes the hex transport and writes the bytes; all parsing and ELF/x86-64 generation happen in TRIPOL. The output path must not exist. The stage compiler can also be built for Windows, where it emits Linux ELF images, not PE files.

The two frontend files in compiler/ are initially copied from bootstrap/ because project imports are confined to their source root. The historical bootstrap demos remain intact. Later changes should maintain equivalence or explicitly document divergence.

Next milestones: runtime expression emission, binary file I/O and command-line arguments, statements and variables, typed functions, PE generation, then repeated self-compilation and comparison to the Go implementation. Do not claim self-hosting until the new compiler can build itself and pass reproducible multi-generation tests.
