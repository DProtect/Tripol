# Language reference — QUANTUM X 0.27.1

## DISCORD REVOLUTION 0.28 development

A named function can be stored as a typed callback. `&format_name` takes a reference to the function. Declare the corresponding parameter or field as `func(string)->string` and call a local callback as a regular function.

```tripol
func format_name(name:string)->string {
    return "Hola, " + name
}

func welcome(render:func(string)->string,name:string)->string {
    return render(name)
}

show welcome(&format_name, "Ada")
```

Callback signatures must match the argument and result types. Pass all arguments when invoking a callback, including arguments whose defaults apply to direct calls. References work in native builds and WebAssembly. The current implementation takes references to named functions with ordinary parameter and result types.

Use `list([&callback_a, &callback_b])` to create a list with a single callback signature. Annotate fields and parameters with `list[func(string)->string]` when they store or receive such a list. A scratch block can call `get`, `size` and `capacity` on a surrounding list. These operations expose immutable values and measurements. Mutations of surrounding lists remain outside scratch blocks.

## QUANTUM X 0.27.1 reference

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

The process-lifetime arena reserves 256 MiB of contiguous virtual address space lazily, enabling 4 KiB pages as allocations advance. On Windows this uses MEM_RESERVE/MEM_COMMIT; on Linux mmap(PROT_NONE)/mprotect. The `memory` module exposes used, accessible/committed and capacity bytes (not RSS). Pages already enabled stay enabled after scratch restores the frontier. `scratch` records and restores the allocation frontier (and input scratch pointer) on normal exit and recoverable unwinding. Locals in the region cannot escape. Outer scalar bindings may be read and modified; outer strings are read-only and outer aggregate bindings cannot be accessed. `return` and loop control leaving a protected region are rejected. Eligible loop-local arrays, strings and lists are automatically reclaimed using proven scratch regions, including supported conditional scopes. General automatic memory reclamation remains pending. Scratch does not sandbox I/O or reverse external side effects. Buffers are reused, not securely wiped. Ordinary allocations outside scratch live until exit. Text concatenation and file/network response caps are 65535 bytes; input is 4096 bytes. See library guide.

## Infinity syntax

`f"Hello {name}, score {score}"` interpolates existing scalar bindings of type string, int or bool. Braces may be escaped as `{{` and `}}`; ordinary quoted strings remain unchanged, including JSON strings. Expressions of type int, bool or string are allowed inside braces, including arithmetic, calls, field access and indexing. Float formatting is not supported. Each field is evaluated once, left to right. Interpolation creates immutable strings using the same bounded arena, with the usual 65535-byte result cap. Use `scratch` for repeated temporary strings when safe.

Functions support trailing optional scalar parameters: `func greet(name="Guest", score:int=10)->string {return f"{name}: {score}"}`. Defaults are side-effect-free scalar constants (including checked foldable integer expressions), evaluated as literals at call sites; type inference for omitted parameter annotations is based on the default, otherwise int. Required parameters must precede optional parameters. Struct constructors may omit trailing fields with scalar defaults. No keyword arguments, dynamic defaults, overloads or callbacks.

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

Mutable fixed and dynamic integer arrays support `a[i] += x`, `-=`, `*=`, `/=`, `%=`. The index is checked and evaluated exactly once before loading the old value, then the RHS is evaluated once. Checked arithmetic and recoverable traps are unchanged. Constant arrays cannot be modified. Dynamic integer lists also support these indexed assignments.

## Growable integer lists (0.27.1)

`nums=list()` or `list(initial_capacity)` creates a stable reference to a growable integer collection. `add(value)` returns new size; `get(index)` and `set(index,value)` check bounds; `remove_at(index)` returns removed value; `clear()` resets size; `size()`, `capacity()`, `len(nums)` and `show nums` are supported. Capacity is at most 1048576, subject to process arena limits. Aliases share one handle across growth; old backing allocations remain in the arena. Direct mutation through a constant binding is rejected; another mutable alias can still mutate the shared collection. This is not general memory reclamation. Lists are distinct from `ints`; methods and checked indexing are supported, but array-only functions do not accept lists.

### Indexed collection iteration and default record fields

```tripol
struct Player {name:string="Guest";score:int=0}
p=Player()
numbers=list()
numbers.add(10)
numbers.add(20)
numbers[0]+=5
for i,value in numbers {show i,value}
```
Default fields must be trailing scalar constants of the declared type. Constructors may omit those fields. Iteration supports fixed arrays, dynamic integer arrays and lists; the bound is captured once on entry and each access remains checked. List indexing and checked compound indexed assignment are now supported. A list index is evaluated once; assignment rechecks the current backing array after the RHS (which may mutate an alias). A shrinking list may cause a bounds trap during iteration rather than an invalid memory access.

## SIMPLE R1

`numbers=list([10,20,30])` initializes a growable integer list. Unlike `numbers=[10,20,30]` (fixed stack array), it can grow with `numbers.add(40)`. The initializer is evaluated left to right once, with up to 400 integer expressions. `list([])` creates an empty list. Existing list(capacity) and fixed arrays retain their behavior. Lists now also infer homogeneous string, bool and float elements from a nonempty initializer. Use `string_list()`, `bool_list()` or `float_list()` for an empty typed list. Mixed lists are rejected. Methods, iteration, display and checked indexing work for these scalar lists. Numeric `sum` and compound indexed assignment remain integer-only.

```tripol
import tripol.errors
numbers=list([10,20,30])
numbers.add(40)
show f"Total items: {len(numbers)}"
try {show numbers[99]} catch code {show errors.message(code)}
```

### Natural dictionaries

```tripol
player = {
    "name": "Alex",
    "health": 100,
    "active": true
}
player["health"] -= 20
show player
```
A dictionary literal is a statically shaped reference with up to 64 unique nonempty literal keys and independently inferred scalar types. It supports checked literal-key lookup, assignment, integer compound assignment, aliases, display and len. Keys and their types cannot change; no dynamic lookup or insertion in this literal type. Use `tripol.dictionary` for runtime string-to-int keys. These are native arena references, not a dynamically typed object model. General reclamation remains limited.

### Simple function inference

A function with a single return expression and no explicit return annotation can infer a scalar return type and omitted parameter types from local expression constraints and calls in the same source unit. For example `func greet(name) => "Hello "+name` and `func double(n) => n*2`. The function is monomorphic; inconsistent calls are rejected by normal type checking. Defaults and explicit signatures are preserved. Ambiguous parameters retain the integer default. Multiple statements, imported call sites and recursive return inference may need explicit annotations.

Integer, boolean and string lists support `contains(value)`, `remove(value)` (first match, returns bool) and `sort()` (in-place, stable, returns size; maximum 4096 items). Strings sort lexicographically by UTF-8 bytes, not locale collation. Float lists retain indexed access, add/remove_at and iteration; value-search and sorting are not included for floats in R1.

Compatibility note: a single-return function with no explicit `->int` can now infer a non-integer return type. Write `->int` if an integer return contract is required, including for exported Vault functions. Explicit mismatches remain compile errors even in unused functions.


## SIMPLE R2

TRIPOL infers a single signature for each function across a project. Calls and typed callback parameters supply context. For example, `apply(twice)` gives `twice` the signature required by `apply`. A function name is a function reference when a local variable does not shadow it. `&twice` remains valid.

```tripol
func apply(action:func(int)->int) { action(21) }
func twice(x) { x*2 }

show apply(twice)
show apply(func(x) { x+1 })
```

An inferred function returns the value of its last expression when it reaches the end. Final `if/else` branches can each yield a value. Explicit `return` remains available. An anonymous function can capture a scalar literal declared with `const`. The compiler copies that value into the function. Mutable captures, runtime values and aggregate captures require explicit arguments. Method overloads use distinct parameter types. The compiler selects one signature from the provided arguments and available defaults. Conflicting evidence requires a type annotation or separate functions.

TRIPOL infiere una firma por función en todo el proyecto. Las llamadas y los parámetros de callback tipados aportan contexto. En el ejemplo, `apply(twice)` aporta a `twice` la firma requerida por `apply`. El nombre de una función actúa como referencia si ninguna variable local lo oculta. `&twice` sigue siendo válido.

Las funciones con retorno inferido devuelven el valor de la última expresión al llegar al final. Las ramas terminales de un `if/else` pueden aportar un valor cada una. También puedes escribir `return`. Una función anónima puede capturar un literal escalar declarado con `const`. El compilador copia ese valor a la función. Pasa como argumentos explícitos las variables mutables, los valores de ejecución y las estructuras. Las sobrecargas de un método se distinguen por los tipos de sus parámetros. El compilador elige una firma a partir de los argumentos y los valores predeterminados. Si hay conflicto, anota el tipo o escribe funciones separadas.

## SIMPLE R3

Functions and methods can overload explicit parameter types. Define callbacks with `->void` when their result is unused. Calls expecting a void callback can infer it for a returnless function literal. Nested structures use typed references. A dictionary literal passed directly where `json.Object` is expected becomes a typed JSON object with nested objects, strings, booleans and integers.

Las funciones y métodos admiten sobrecargas con parámetros explícitos. Usa `->void` para callbacks sin resultado. El compilador puede inferirlo para una función anónima sin retorno cuando la firma esperada es void. Las estructuras anidadas conservan referencias tipadas. Un literal de diccionario usado directamente donde se espera `json.Object` se convierte en un objeto JSON tipado y admite objetos anidados, texto, booleanos y enteros.

## SIMPLE R3 collection and JSON types / Colecciones y tipos JSON

Use `list([Point(1), Point(2)])` to infer a homogeneous list of records. Declare an empty typed list with `let points:list[Point]=list()`. A nested type such as `list[list[Point]]` keeps the element type at each level. Use `rows.get(0).get(1)` to read a nested element and bind the inner list to a local variable before mutating it. Lists share a checked handle when assigned to another variable. They are allocated in the program arena. A `scratch` scope cannot extract an aggregate element from an outer list since it could let temporary storage escape. Search and sorting on object lists require an explicit field selection.

Usa `list([Point(1), Point(2)])` para inferir una lista homogénea de estructuras. Declara una lista vacía con `let points:list[Point]=list()`. Un tipo como `list[list[Point]]` conserva el tipo de cada nivel. Lee un elemento con `rows.get(0).get(1)` y asigna la lista interna a una variable local si vas a modificarla. Los alias comparten un identificador verificado. Las listas se asignan en la arena del programa. Un bloque `scratch` restringe el acceso a elementos compuestos de listas exteriores para que la memoria temporal no salga del bloque.

`json.decode(Report, raw)` uses the struct declaration at compile time. It handles user structs, nested structs, typed lists, text, booleans and signed 64-bit integers. It raises an assertion on malformed JSON, missing required fields, duplicate keys or incompatible values. Unknown fields are accepted. See Final Polish below for decimals, defaults and optional fields. Function references are outside the JSON schema. Use the same source on native targets and in the Playground. [Runnable example](../examples/simple_r3_objects.tri)

`json.decode(Report, raw)` usa la estructura en tiempo de compilación. Admite estructuras anidadas, listas tipadas, cadenas, booleanos y enteros de 64 bits. Genera un error de aserción ante JSON inválido, campos requeridos ausentes, claves duplicadas o valores incompatibles. Admite campos adicionales. Final Polish añade decimales, valores predeterminados y campos opcionales. Las referencias a funciones quedan fuera del esquema JSON. [Ejemplo ejecutable](../examples/simple_r3_objects.tri)

## SIMPLE R3 Final Polish / Acabado final de SIMPLE R3

Use named arguments after positional arguments. A call such as `Player(speed: 2.75, name: "Lina")` fills omitted parameters or fields that have defaults. Expressions run once in the order written, even when the names appear out of declaration order. In an overloaded call, the names and types both participate in selecting the signature.

Usa argumentos nombrados después de los posicionales. Una llamada como `Player(speed: 2.75, name: "Lina")` conserva los valores predeterminados de los campos omitidos. Cada expresión se evalúa una vez, en el orden escrito. Los nombres y tipos intervienen en la selección de una sobrecarga.

Declare an optional value with `T?`. Use `null` for absence and `some(value)` for presence. Read it with `.has()`, `.unwrap()` or the lazy fallback `.or_else(value)`. A nested list can be read as `rows[0][1]`. Computed indices are checked on native and WebAssembly targets.

Declara un valor opcional con `T?`. Usa `null` para la ausencia y `some(value)` para un valor presente. Puedes comprobarlo con `.has()`, leerlo con `.unwrap()` o proporcionar un valor alternativo con `.or_else(value)`. La alternativa solo se evalúa si falta el valor. Usa `rows[0][1]` para leer una lista anidada con comprobación de límites.

An anonymous callback takes a snapshot of captured integer, Boolean, decimal and text values when it is created. Later changes to the original variables do not change the snapshot. Callbacks may outlive the function where they were created. Temporary `scratch` memory cannot be assigned to an outer callback. Mutable aggregate capture and general parallel execution of callback bodies remain pending.

Un callback anónimo captura el valor de las variables escalares cuando se crea. Los cambios posteriores de las variables originales no cambian la captura. El callback puede sobrevivir a la función en la que se creó. Un bloque `scratch` impide asignarlo a un callback exterior cuando su memoria es temporal. Siguen pendientes las capturas de agregados mutables y la ejecución paralela general de callbacks.

`json.decode(Player, text)` uses defaults for absent fields and validates the types of present fields. Optional fields accept `null`. `json.encode(Player, player)` serializes the fields in declaration order, with exact finite binary64 decimal output. Invalid recognizable fields produce a diagnostic naming the struct and field. Malformed input, overflow and unrepresentable string contents can be reported by the underlying value reader. Cycles fail the serialization depth check.

`json.decode(Player, texto)` usa los valores predeterminados cuando faltan campos y valida los que están presentes. Los opcionales aceptan `null`. `json.encode(Player, player)` serializa los campos en orden de declaración y conserva el valor exacto de los decimales binarios finitos. Una incompatibilidad de tipo identifica la estructura y el campo cuando el lector puede reconocerlo. La serialización limita la profundidad y rechaza los ciclos. [Ejemplo / Example](../examples/simple_r3_polish.tri)
