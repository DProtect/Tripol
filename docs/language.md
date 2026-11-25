# Language reference — 0.12

Source files are `.tri`. Imports must precede declarations and statements; use `import tripol.math`, `import vault.name`, or a quoted project-local `.tri` path. `//` and `#` start comments; newline or semicolon separates statements. Blocks use braces. Strings use JSON-style quoting and escape sequences (no embedded NUL). Identifiers are ASCII; text is byte-oriented. `let` declares mutable bindings, `const` immutable bindings; bare assignment creates or updates a mutable local. No shadowing. Types: signed checked int64 (`int`), `bool`, IEEE binary64 (`float`), immutable `string`, runtime-length `ints`, named `struct`. Integer and float operations require matching types. `float(n)` converts integers; `show` prints floats with six decimals. Struct fields support scalar, string and `ints` types, not nested records. Struct and dynamic array assignments alias references, and const is not deep immutability.

Functions use `func name(a:int,b:int)->int {return a+b}`; omitted types default to int. Functions have no globals or closures. `if/else if/else`, `while`, `for i in range(start,end)`, `for n in array`, `repeat`, `break`, `continue` are supported. `and/or/not` short circuit; `+=` etc. modify bindings. Fixed arrays of up to 400 integers use stack space; runtime-sized `ints(n)` has 0..65535 elements and cannot grow. The 401 stack-slot limit includes temporary locals and exception records. Recursion can exhaust native stack. `assert condition,"message"` checks runtime conditions.

## Exceptions

`try { ... } catch code { ... } finally { ... }` (finally optional). Catch receives integer 10 for arithmetic, 11 for bounds, 12 for assertions, or 100..255 from `throw n` (other values normalize to 100). Catches work across function calls and nested blocks. Finally runs after success, recovery, or a recoverable failure in catch; an exception in finally replaces a pending one. Native input, file, socket, float-format, UI and allocator failures remain fatal. Return inside protected blocks is rejected; break/continue may only target loops inside the same protected block. See [errors](errors.md). No exceptions for OS faults or native stack overflow.

## Temporary memory

```tripol
total=0
repeat 1000 {
    scratch {
        values=ints(65535)
        values[0]=7
        total+=values[0]
    }
}
show total
```

The process-lifetime arena remains limited to 8 MiB and is initialized lazily. `scratch` records and restores the allocation frontier (and input scratch pointer) on normal exit and recoverable unwinding. Locals in the region cannot escape. Outer scalar bindings may be read and modified; outer strings are read-only and outer aggregate bindings cannot be accessed. `return` and loop control leaving a protected region are rejected. No general automatic memory reclamation. Scratch does not sandbox I/O or reverse external side effects. Buffers are reused, not securely wiped. Ordinary allocations outside scratch live until exit. Text concatenation and file/network response caps are 65535 bytes; input is 4096 bytes. See library guide.
