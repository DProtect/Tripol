# TRIPOL WEB ENGINE R3 — compatibility contract

The web backend reuses the native parser and x86 validation pass, then emits WebAssembly with an explicit browser host ABI. It is not yet a universal shared IR or a full replacement for the native backend. The existing educational interpreter remains available.

## Implemented in source (subject to CI validation)
- int64, bool, string and float scalar lowering; string interpolation, concatenation and equality.
- Scalar functions, defaults, recursion, control flow.
- Fixed and dynamic integer arrays, indexed reads and writes, iteration, sum and length.
- Typed growable lists, collection methods and indexed access.
- Record constructors, defaults, typed field reads/writes; shaped dictionaries with literal-key reads/writes and len.
- Console input; selected pure math functions.
- Bounded browser host handles; checked integer arithmetic, array bounds, step and output limits.

## Remaining parity work
- Exceptions: native try/catch/finally/throw and recoverable faults across calls.
- Scratch memory, allocation frontier, native memory metrics, and identical lifetime behavior.
- Complete aggregate display, full library intrinsic coverage, precise float/UTF-8 behavior, and alias/lifetime edge cases.
- Virtual project files, package imports, virtual filesystem, permitted browser networking, graphics and audio adapters.
- Shared whole-program typed IR; native-equivalent optimization, SIMD and incremental cache.
- Windows differential testing and full platform-specific capability matrix.

## Verification
The workflow builds the Go compiler for wasm, runs actual generated WASM modules, and compares the 100 Playground sources to Linux native executables. The parity-report.json records matches, differences, and rejections. Required smoke tests and all 100 native differential cases gate subsequent publication. Earlier R3 builds may have a partial parity report. No claim of full language compatibility is made. The browser discovers the highest verified ABI and preserves explicit engine selection; it does not silently replace a WASM failure with the interpreter.

Generated modules require the matching host ABI. The browser is sandboxed; native OS access, arbitrary filesystem paths and raw sockets cannot be provided directly. Unsupported features must produce explicit diagnostics rather than fabricated output.

## WEB LIBRARIES — portable expansion
The R3 ABI now has source implementations for byte-oriented strings, byte slicing and search, integer conversion, input, array equality/fill/range fill, checked u32 operations and monotonic clocks. Library imports use the original embedded TRIPOL sources rather than JavaScript reimplementations of library algorithms.

The browser import allowlist includes 50 candidate modules: math, arrays, collections, datetime, bits, bitset, binary, buffer, checksum, colors, compression, console, convert, dictionary, encoding, geometry, json, matrix, path, physics, random, search, sequences, simd, sort, stats, text, validate, vector, stack, queue, animation, camera, particles, canvas, sprite, testing, errors, network, http, image, time, profiler, heap, utf8, csv, stream, regex, cli, crypto.

Allowlisted is not equivalent to fully compatible. Required smoke tests for datetime, text and JSON run before publication. A separate library-parity-report.json compares one real program per candidate module against native Linux; passing a probe does not certify every function. Browser limits remain intentionally smaller than native limits. Native-only operations such as TCP exchange, image file export, sleep and native device access may still be rejected even if their containing module is allowed. Unsupported modules require dedicated capability adapters. No simulated successful output.
