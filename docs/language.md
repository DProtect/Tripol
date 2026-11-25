# Language reference — 0.21

Source files are `.tri`. Imports must precede declarations and statements; use `import tripol.math`, `import vault.name`, or a quoted project-local `.tri` path. `//` and `#` start comments; newline or semicolon separates statements. Blocks use braces. Strings use JSON-style quoting and escape sequences (no embedded NUL). Identifiers are ASCII; text is byte-oriented. `let` declares mutable bindings, `const` immutable bindings; bare assignment creates or updates a mutable local. No shadowing. Types: signed checked int64 (`int`), `bool`, IEEE binary64 (`float`), immutable `string`, runtime-length `ints`, named `struct`. Integer and float operations require matching types. `float(n)` converts integers; `show` prints floats with six decimals. Struct fields support scalar, string and `ints` types, not nested records. Struct and dynamic array assignments alias references, and const is not deep immutability.

Functions use `func name(a:int,b:int)->int {return a+b}`; omitted types default to int. Functions have no globals or closures. `if/else if/else`, `while`, `for i in range(start,end)`, `for n in array`, `repeat`, `break`, `continue` are supported. `and/or/not` short circuit; `+=` etc. modify bindings. Fixed arrays of up to 400 integers use stack space; runtime-sized `ints(n)` has 0..4194304 elements and cannot grow directly; `tripol.buffer` implements an explicit growable integer backing array. The 401 stack-slot limit includes temporary locals and exception records. Recursion can exhaust native stack. `assert condition,"message"` checks runtime conditions.

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

The process-lifetime arena reserves 256 MiB of contiguous virtual address space lazily, enabling 4 KiB pages as allocations advance. On Windows this uses MEM_RESERVE/MEM_COMMIT; on Linux mmap(PROT_NONE)/mprotect. The `memory` module exposes used, accessible/committed and capacity bytes (not RSS). Pages already enabled stay enabled after scratch restores the frontier. `scratch` records and restores the allocation frontier (and input scratch pointer) on normal exit and recoverable unwinding. Locals in the region cannot escape. Outer scalar bindings may be read and modified; outer strings are read-only and outer aggregate bindings cannot be accessed. `return` and loop control leaving a protected region are rejected. No general automatic memory reclamation. Scratch does not sandbox I/O or reverse external side effects. Buffers are reused, not securely wiped. Ordinary allocations outside scratch live until exit. Text concatenation and file/network response caps are 65535 bytes; input is 4096 bytes. See library guide.

## Infinity syntax

`f"Hello {name}, score {score}"` interpolates existing scalar bindings of type string, int or bool. Braces may be escaped as `{{` and `}}`; ordinary quoted strings remain unchanged, including JSON strings. Only simple identifiers are allowed inside braces; float formatting, fields and arbitrary expressions are not supported. Interpolation creates immutable strings using the same bounded arena, with the usual 65535-byte result cap. Use `scratch` for repeated temporary strings when safe.

Functions support trailing optional scalar parameters: `func greet(name="Guest", score:int=10)->string {return f"{name}: {score}"}`. Defaults are side-effect-free scalar constants (including checked foldable integer expressions), evaluated as literals at call sites; type inference for omitted parameter annotations is based on the default, otherwise int. Required parameters must precede optional parameters. Struct constructors still require all fields. No keyword arguments, dynamic defaults, overloads or callbacks.

`tripol ir FILE.tri` validates the project and prints the subset of small pure checked int64 expressions eligible for TIR lowering. This is not a full program IR. Nano may inline small zero/one-parameter functions; turbo can inline up to four parameters; trace/debug retain call boundaries. The inliner evaluates arguments exactly once and keeps overflow checks. The general backend remains unchanged.

## Quantum syntax and growable buffers

`func twice(n) => n * 2` is shorthand for a single returning expression. The default return type is int; annotate other return types, e.g. `func title()->string => "TRIPOL"`. No anonymous functions or struct methods.

```tripol
import tripol.buffer
b = buffer.create()
b = buffer.push(b, 10)
b = buffer.push(b, 20)
show buffer.values(b)
```

`buffer.push` returns a backing array and must be assigned back. It doubles capacity as necessary up to 4,194,303 values, retaining old allocations until exit or safe scratch reclamation. Existing aliases share in-place updates until growth replaces the backing array; this is not an opaque list object. Mutating library helpers require dynamic arrays rather than borrowed stack literals. See library guide.

## Nova collections

`dictionary.create(capacity)` provides a fixed-capacity byte-string-to-int64 map using an integer array, 1..4096 entries, keys 1..32 bytes. `dictionary.set` returns false when full; `get` has an optional fallback. `remove` clears a slot for reuse. This is not general dictionary literal syntax or a hash table. Avoid mutating its raw storage.

CORE preserves the 0.15 syntax and signed integer division semantics (truncation toward zero, remainder sign follows dividend). Positive literal powers of two can compile to shifts and bias correction. Arithmetic overflow checks remain active for potentially overflowing operations. No new public syntax or memory ownership model in this release.

HYPER adds no new syntax. `balanced` and `turbo` preserve checked integer semantics, argument evaluation order and all source validation. Register lowering is limited to pure scalar leaves and call-free straight-line integer range loops; guards and tracing disable loop register promotion.

TITAN adds no new syntax. Eligible call-free integer range loops may have up to four promoted outer accumulators and conditional branches with break/continue. All checked arithmetic and validation remain.

APEX adds no syntax. The optimizer may reuse a total invariant expression and specialize restricted binary recurrences; integer arithmetic remains checked.

OMEGA adds no grammar changes. ArrayCore specializes eligible dynamic integer array reduction loops without changing checked arithmetic or bounds semantics.

ECLIPSE adds the stateful Discord REST API (see [Discord guide](discord.md)) and native array copy/search; no new grammar.

Direct self tail calls with full argument lists reuse the frame in all profiles unless tracing. Other recursion remains subject to the OS stack limit.

## Indexed compound assignment (v0.27)

Mutable fixed and dynamic integer arrays support `a[i] += x`, `-=`, `*=`, `/=`, `%=`. The index is checked and evaluated exactly once before loading the old value, then the RHS is evaluated once. Checked arithmetic and recoverable traps are unchanged. Constant arrays cannot be modified. Dynamic collection method syntax remains pending.
