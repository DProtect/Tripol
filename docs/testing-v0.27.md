# QUANTUM X — Windows validation

Run only after the v0.27.1 release is published. Save v0.26 results first, on the same PC and with the same profile, workload and power settings:

```cmd
tripol --version
tripol arena --runs 20 --profile turbo --tripol-only --json arena-v026.json
tripol update --check
tripol update
```

UPDATE R1 downloads and verifies the official installer. Close TRIPOL processes, run the verified unsigned installer from the printed path and reopen the terminal. Do not install if verification failed.

## Installed compiler

```cmd
where tripol
tripol --version
tripol doctor
tripol new Test027
cd Test027
tripol check
tripol run
tripol build
main.exe
```

Expected: TRIPOL 0.27.1 and Hello, TRIPOL! from the generated program. Check all four profiles: nano, balanced, turbo, debug. Each must give equivalent results.

## Validation Center

From a clone of the repository, in PowerShell:

```powershell
powershell -NoProfile -File .\tests\manual\test-v027.ps1
```

It verifies the installed version, native compiler execution, all four profiles, compound array operations, recoverable arithmetic/bounds failures, new libraries and invalid source rejection. It writes `tripol-validation-v027.json` in the invocation directory. Add `-Benchmark` for the full Arena run, or `-CheckUpdates` for a live network check. It never installs updates. If script execution is restricted, review the source and use a process-scoped execution policy if permitted, or run the commands manually.

## GENESIS R3

```cmd
tripol build compiler/file.tri --target windows -o genesis-stage.exe
tripol genesis genesis-stage.exe examples/projects/genesis_r3.tri -o genesis-demo
```

Expected: an ELF64 Linux x86-64 executable. On a Linux x86-64 system run `chmod +x genesis-demo && ./genesis-demo`; expected output 30. The stage is written in TRIPOL; Go CLI transports hexadecimal bytes to disk (no Python). Linux can execute the emitted ELF; Windows only hosts the compiler. No Windows PE generation from GENESIS or full self-hosting yet.

## Libraries and memory

```cmd
tripol run examples/projects/quantum_x_libraries.tri
tripol run examples/projects/memory_stress.tri
```

Library output: 10 20 30, then 8, then 3. The stress example performs one million push/pop pairs on a bounded heap and verifies unchanged used arena bytes. Arena bytes are not total process RSS; use Task Manager separately. This does not demonstrate general automatic reclamation.

## Graphics

```cmd
tripol build examples/projects/quantum_x/main.tri -o quantum-x.exe
quantum-x.exe
```

Use WASD/arrows, collect the gold squares, ESC to quit. The console prints score, frames and elapsed nanoseconds. Check controls and stability. CPU software renderer; no GPU acceleration claimed.

## Editor

Install the official v0.27.1 VSIX through VS Code's Extensions: Install from VSIX command. Open a .tri file, try import completion, function completion/hover and local function definition lookup. Save, then run TRIPOL: Check source. Try invalid code and inspect Problems. Explicit Build/Run commands require trusted workspace; check-on-save is optional and off by default. No formatter/full LSP yet.

## Updates and offline

Run `tripol update --check`, `tripol update --history`, and (when a newer release exists) `tripol update --download-only`. Disconnect the network and build a local program: compilation must remain independent of the network. Explicit online checks may fail cleanly. Installation remains manual in this implementation; no automatic rollback. Use local Go HTTP tests to exercise invalid metadata, tampered assets and interrupted downloads; do not publish fake releases.

## Repeat measurements and source tests

```cmd
tripol arena --runs 20 --profile turbo --tripol-only --json arena-v027.json
go test -count=1 -v ./core ./tests
go vet . ./core ./stdlib ./tests
go build -o tripol.exe .
tripol.exe --version
```

Compare medians and p95 for the same case names in both JSON files. Record hardware, versions, profile, output equivalence and sample counts. No universal speed claims. Keep any failure logs; do not disable tests.

For the 0.27.1 additions: `tripol run examples/projects/dynamic_list.tri`, `tripol run examples/projects/binary_io.tri` from a disposable working directory (writes quantum-bytes.bin), and `tripol build compiler/direct.tri -o genesis-stage.exe`. See [GENESIS](genesis.md) for the direct driver. Retain prior version reports for comparison.
