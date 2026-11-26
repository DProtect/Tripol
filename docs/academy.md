# TRIPOL Academy — SIMPLE R1

Offline beginner course. No account, payment or network required. The compiler includes nine lessons; their exact source examples run in the native regression suite. Start with `tripol learn es list` or `tripol learn en list`, then `tripol learn es 1`. Each lesson displays source, expected output and an exercise. Copy it into `main.tri` and run `tripol run`. The older `tripol learn es` and `tripol learn en` overview commands still work. The lessons do not save progress or automatically grade answers yet.

## 1. Primer programa / First program

```tripol
show "Hola, TRIPOL!"
```

Expected:
```text
Hola, TRIPOL!
```

show muestra texto en la terminal. Cambia el saludo. / show prints text to the terminal. Change the greeting.

## 2. Variables / Variables

```tripol
name = "Alex"
score = 10
show f"{name}: {score * 2}"
```

Expected:
```text
Alex: 20
```

Las variables infieren su tipo. Cambia score y vuelve a ejecutar. / Variables infer their type. Change score and run again.

## 3. Condiciones / Conditions

```tripol
score = 80
if score >= 60 {
    show "Passed"
} else {
    show "Try again"
}
```

Expected:
```text
Passed
```

Cambia score a 40 para probar la otra rama. / Change score to 40 to test the other branch.

## 4. Bucles / Loops

```tripol
total = 0
for n in range(1, 6) {
    total += n
}
show total
```

Expected:
```text
15
```

range excluye el limite final. Prueba sumar del 1 al 10. / range excludes its end. Try summing 1 through 10.

## 5. Funciones / Functions

```tripol
func twice(n) => n * 2
show twice(7)
```

Expected:
```text
14
```

Una funcion puede devolver un resultado. Prueba otro numero. / A function can return a result. Try another number.

## 6. Listas / Lists

```tripol
numbers = list([10, 20, 30])
numbers.add(40)
show numbers
```

Expected:
```text
[10, 20, 30, 40]
```

Una lista puede crecer. Los indices comienzan en cero. / A list can grow. Indices start at zero.

## 7. Estructuras / Records

```tripol
struct Player {
    score: int = 0
}
p = Player()
p.score = 5
show p.score
```

Expected:
```text
5
```

Una estructura agrupa datos relacionados. / A record groups related data.

## 8. Bibliotecas / Libraries

```tripol
import tripol.math
show math.square(7)
```

Expected:
```text
49
```

tripol libs enumera bibliotecas; tripol docs math muestra sus funciones. / tripol libs lists libraries; tripol docs math shows their functions.

## 9. Errores / Errors

```tripol
import tripol.errors
try {
    show 10 / 0
} catch code {
    show errors.message(code)
}
```

Expected:
```text
Integer overflow or division by zero
```

try/catch recupera errores comprobados. No todos los fallos son recuperables. / try/catch handles checked errors. Not all failures are recoverable.

## Test your projects

Write `assert` checks in files ending in `.test.tri` and run `tripol test examples/academy` from the repository root. `tripol test` selects only immediate `*.test.tri` files in the current directory; `tripol test FILE.tri` explicitly selects one file. Every selected program is compiled with the balanced profile, run in a separate temporary working directory and checked for exit status. Ten seconds per program, two minutes total execution budget, at most 64 tests. Output from failed tests is limited to 64 KiB. A caught assertion is not a failure. This is not a sandbox; only run trusted project tests, which have ordinary native process permissions. No package installation or network access is performed by the runner itself.

## Interactive shell

Run `tripol shell` to try scalar expressions, assignments and interpolated text without making a project. Each accepted command compiles to a temporary native executable. Pure scalar assignments persist by deterministic replay; failed commands do not change session source. `:help`, `:code`, `:reset`, `:exit` are available. The shell rejects arbitrary calls, imports, loops, aggregate mutation and IO so replay cannot repeat observable effects. It is a teaching shell, not a full incremental REPL.
