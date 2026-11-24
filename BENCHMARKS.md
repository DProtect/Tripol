# ORIGIN — preliminary measurements

Environment: Linux x86-64, Python 3.13.5, Clang 17, GCC 14.2.

| Example | Compiler | Profile | Binary size |
|---|---|---|---:|
| hello | Clang | nano | 14,488 bytes |
| fibonacci | Clang | nano | 14,504 bytes |
| demo | GCC | nano | 14,472 bytes |

Fibonacci result: fib(30) = 832040. GNU time reported peak RSS 512 KiB
on five runs and elapsed time rounded to 0.00 seconds. These are very short
processes; RSS can be underreported and the elapsed measurement is too coarse
for performance comparison. The binaries link dynamically to the system C
runtime. These figures are not measurements of the compiler itself or Windows.

Reproduce on Linux:

```sh
python tripol.py build examples/fibonacci.tri --cc clang -o fibonacci
for i in 1 2 3 4 5; do /usr/bin/time -f '%M KiB %e s' ./fibonacci; done
```

No comparative Rust/Python performance or energy benchmark has been run.
