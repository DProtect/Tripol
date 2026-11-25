# TRIPOL v0.7.0 — HORIZON

Free, experimental Apache-2.0 programming language. Standalone Go compiler emits x86-64 machine code and PE32+ (Windows) or ELF64 (Linux) directly. No Python, GCC, Clang or Go installation is required to use the distributed compiler. Go 1.23+ is needed only to rebuild it.

## Windows x64

Run `TRIPOL-v0.7-Setup-x64.exe` to update the per-user installation at `%LOCALAPPDATA%\Programs\TRIPOL`. Existing example files are preserved, and the installer adds the directory to the user's PATH. Close and reopen CMD; run `tripol --version`. The installer uses Windows PowerShell for PATH modification and is unsigned. v0.6 was confirmed by the project owner on Windows (calculator template and offline tutorial); v0.7 is cross-compiled and PE-inspected, with Windows execution pending. Portable ZIP is also available.

```cmd
tripol --version
tripol new MyProject
cd MyProject
tripol run main.tri
tripol build main.tri
main.exe
tripol bench main.tri --runs 5
```

Commands: `tripol check FILE.tri`, `tripol run FILE.tri [--profile nano|debug]`, `tripol build FILE.tri [-o OUTPUT] [--target windows|linux] [--profile nano|debug]`, `tripol bench FILE.tri [--runs 1..100] [--profile nano|debug]`, `tripol new NAME`. Default profile nano. Debug disables the AST/native optimization passes, not a source debugger. Benchmark includes OS process startup, compiles once, excludes one warm-up, closes stdin and discards stdout. Use noninteractive programs only; no built-in timeout. Targets x86-64 only.

## Horizon: typed structures, assertions and Spanish tutorial

```tripol
struct Jugador {
    nombre: string
    puntos: int
    activo: bool
}
func bonus(j: Jugador, cantidad: int) -> Jugador {
    j.puntos += cantidad
    return j
}
j = Jugador(input("Tu nombre?"), 10, true)
j = bonus(j, 5)
assert j.puntos == 15, "Puntos incorrectos"
show "Jugador:", j.nombre
show "Puntos:", j.puntos
show "Activo:", j.activo
```

- `struct` declares a named type at module scope, with 1–64 named fields. Types are int, bool, float, string and ints. A constructor requires every field in declaration order; fields can be read with dot syntax and assigned on a mutable named variable, including `+=` and other supported compound operations. Structures can be passed to or returned from functions and imported with local modules. No nested record fields, default arguments, methods, inheritance, null, or structural equality yet.
- Records have reference semantics: assignments and parameters alias the same allocation. `const` forbids direct mutation through that binding but not through another alias. Each record reserves eight bytes per field in the existing shared 8 MiB arena and lives until exit. No per-object reclamation yet. Array fields are dynamic integer-array references; read into a variable before indexing.
- `assert condition` or `assert condition, "message"` checks a boolean and exits with status 1 and a stderr diagnostic on failure. The message must be a string literal. It remains active in nano and debug profiles. `struct` and `assert` are newly reserved names.
- `tripol learn es` opens the Spanish offline tutorial; `tripol learn` or `tripol learn en` opens the English tutorial. `tripol new Name --template records` generates a playable console example. VS Code highlighting includes the new keywords.
- Native string-table emission now maintains first-use order directly, eliminating the former nested search over labels and distinct literals. Native code for existing programs is unchanged; the benchmark measures compiler time, not application runtime.

## Nexus: easier syntax

```tripol
name = input("Your name?")
age = input_int("Your age?")
show "Hello", name, age
numbers = [10, 20, 30]
for n in numbers { show n }
show "Total:", sum(numbers)
repeat 3 { show "TRIPOL" }
for i in range(5) { show i }
if age < 13 { show "Child" } else if age < 18 { show "Teen" } else { show "Adult" }
```

- `show a, b, c` writes scalar values separated by one space and ends with a newline. It streams values left to right without allocating concatenated strings. A later failure can leave partial output. Single-value `show` is unchanged.
- `input(prompt)` and `input_int(prompt)` optionally print a text prompt on its own line. `input_int` uses the existing strict decimal integer parser; invalid input traps. `input()` remains valid.
- `repeat n` and `range(n)` evaluate their integer count once; negative or zero counts perform zero iterations. `break` and `continue` work normally. `for value in array` iterates a named integer array with an immutable loop value; the initial length is captured, but each indexed access uses the current binding and retains bounds checking. Avoid rebinding the array during iteration.
- `+=`, `-=`, `*=`, `/=`, `%=` work on existing mutable scalar variables, with checked arithmetic and existing type rules. `and`, `or`, `not` are short-circuit logical operators; `not` binds less tightly than comparison but more tightly than `and`/`or`. `repeat`, `and`, `or`, `not` are new reserved words. Older programs that used them as variable names need renaming.
- `sum(array)` is a native checked integer accumulation over a named fixed or dynamic integer array; `sum` of an empty dynamic array is zero. `min(a,b)`, `max(a,b)`, `abs(a)` operate on int64; overflow traps. User-defined functions continue to take precedence over public built-ins.
- `tripol run`, `tripol build`, `tripol check`, and `tripol bench` default to `main.tri` if no source is supplied. `tripol learn` displays an offline tutorial; `tripol doctor` reports the compiler location and supported host/targets (not a full hardware diagnostic).
- `tripol new MyApp --template calculator` or `--template files` creates a ready-to-run example. `basic` is the default. Existing directories are not overwritten.

Nexus adds language conveniences and native math. Horizon adds typed structures; existing memory and platform limitations still apply.

## Interactive programs

```tripol
show "Como te llamas?"
nombre = input()
show "Hola, " + nombre
show "Escribe un entero:"
n = parse_int(input())
show n * 2
show sqrt(81.0)
```

`input()` reads one line (LF or CRLF), returning empty string on EOF, with maximum 4096 bytes before line ending. NUL, long lines and I/O errors trap; Windows console encoding follows the console's code page (ASCII examples are recommended until Windows testing). `parse_int` accepts optional +/- followed by decimal digits only; overflow and invalid text trap. `string + string` concatenates at runtime; `len(string)` counts bytes, not Unicode characters. Strings are immutable, can be passed and returned safely within a process, and use NUL termination. Concatenation has a maximum length of 65535 bytes. The 8 MiB arena is allocated lazily with mmap/VirtualAlloc, has a reusable input scratch buffer and retains each resulting string until process exit. No garbage collector, per-object reclamation or general heap management; arena exhaustion traps. Unused runtime helpers and arena are omitted from nano builds.

## Float64

```tripol
func area(ancho: float, alto: float) -> float {
    return ancho * alto
}
show area(2.5, 4.0)
show float(7) / 2.0
```

Float literals use digits on both sides of the dot (no exponent syntax). Arithmetic + - * /, negation, comparisons, sqrt(float) and float(int) use IEEE-754 binary64/SSE2 with no implicit int/float coercion. NaN comparisons are unordered. `show` rounds to fixed six fractional digits, approximately up to 9.22e12 in magnitude; nonfinite or out-of-range values trap when printed. Negative values rounded to zero print positive zero. This is not shortest-roundtrip or arbitrary precision formatting.

## Local modules

```tripol
import "lib/math.tri"
show area(2.5, 4.0)
```

See `examples/modules.tri` and `examples/lib/math.tri`. Imports must appear first in each file and use quoted local relative `.tri` paths, resolved relative to the importing file and confined to the entry file's project directory. Modules may define functions only; they share a flat function namespace. Duplicate imports are deduplicated, cycles and duplicate function names rejected. No external package manager or namespaced import syntax yet.

## FileSystem R1

```tripol
write_text("note.txt", "Hello")
append_text("note.txt", " TRIPOL!")
show read_text("note.txt")
```

`write_text(path, text)` creates or truncates a file; `append_text` creates or appends.
Both return the number of bytes written. `read_text` returns up to 65535 bytes,
rejects embedded NUL and traps on errors or larger files. The reader reserves
65536 bytes temporarily and reclaims the unused suffix of that reservation.
Strings are byte sequences without embedded NUL. Windows paths use CreateFileA
and therefore the system ANSI code page; ASCII paths are recommended. Relative
paths refer to the program's current working directory. Operations are not atomic,
not sandboxed, and do not create parent directories. Do not use untrusted paths.

## Runtime-sized arrays and text

```tripol
n = parse_int(input())
a = ints(n)
for i in range(0, len(a)) { a[i] = i * 2 }
show a[0]
show text(42)
show slice("abcdef", 1, 4)
```

`ints(n)` allocates a zero-initialized array of 0..65535 signed 64-bit integers.
It is a reference type (`ints`), may be passed to and returned from functions,
and assignment aliases the same array. A const binding rejects direct element
assignment but is not deep immutability. Length cannot be changed; this is not
an expandable vector. Indexing checks bounds, and array references cannot be
printed or compared. Allocation uses the existing shared 8 MiB arena; memory
is not individually reclaimed and exhaustion traps. Fixed local literal arrays
continue to use stack storage. `text(int)` returns decimal text, including the
full signed 64-bit range. `slice(s,start,end)` returns a new string over checked,
zero-based byte offsets with exclusive end; it may split UTF-8 sequences.

## Previous features and limitations

v0.3 features preserved: int64, bool, literal strings, functions, recursion, checked integer arithmetic, while, if/else, short-circuiting, fixed-size local integer arrays with bounds checks, range loops, break/continue, lexical stack reuse, dead function elimination and direct self-tail-call optimization. Fixed local literal arrays cannot be copied, resized, nested, returned or passed. Runtime-sized arrays have the reference semantics described above. UTF-8 BOM accepted. Max 400 local slots. Not a sandbox; ordinary recursion can exhaust the stack. No general allocator, TIR, networking, GUI, methods or production-grade ABI yet. The dynamic arrays and text share the bounded arena. Native Linux binaries use syscalls; Windows binaries use kernel32.dll. Go is only the compiler implementation. Runtime templates are hand-written x86-64 assembly; the checked-in generated Go file allows rebuilding without assembly tools or Python. `native/generate.py` uses Python and GNU binutils only if a maintainer edits assembly templates.

## Build from source

```sh
go test -count=1 -v .
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o tripol .
mkdir -p installer/payload
cp examples/hello.tri installer/payload/hello.tri
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o installer/payload/tripol.exe .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o TRIPOL-v0.7-Setup-x64.exe ./installer
```

See BENCHMARKS.md and test-results.txt. 35 test groups and 198 subtests passed on Linux; go vet passed. No claim of outperforming Rust or C. License Apache-2.0.
