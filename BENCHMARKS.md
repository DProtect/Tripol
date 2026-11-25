# NANO v0.3 — validation and measurements

Linux x86-64 environment, Go 1.23.2. 9 separate process executions each of the same recursive fib(35), output verified as 9227465. Wall-clock includes process start; no algorithm substitution or memoization. The user's v0.2 Windows measurement (23,350.0547 ms) is not directly comparable with these results.

| Test | v0.2 | v0.3 |
|---|---:|---:|
| Linux recursive fib(35) executable | 827 B | 616 B |
| Linux fib(35) median elapsed | 83.40 ms | 75.52 ms |
| Linux hello executable | 662 B | 380 B |
| Windows hello PE | 2048 B | 1536 B |
| Windows recursive Fibonacci PE | 2048 B | 1536 B |

Windows v0.3 files were generated and PE format inspected; actual execution remains pending. The prior v0.2 Windows test was supplied by the project owner. No cross-language speed, energy or RSS claims. Timing varies by machine and load.

Regression suite includes 44 baseline cases plus format/determinism checks, 15 additional valid cases, 13 additional invalid cases, 3 bounds traps, optimization size and profile equivalence checks. The tail-recursive test counts to 1,000,000 without growing a frame per call. See test-results.txt.

The exact Windows example provided by the owner, with an extra `show "TRIPOL FIBONACCI TEST"` before the numeric output, remains 2048 B in both versions due to PE section alignment. The 1536 B Fibonacci result above is for the version printing only the number.
