# Roadmap

TRIPOL is experimental. These are goals, not promised versions or dates.

1. Validate Windows UI/HTTPS and native ABI; add automated Windows tests and signed release distribution when feasible.
2. Expand error recovery with safe cleanup for native resources, richer diagnostics and source locations.
3. Extend memory management beyond explicit scratch regions, with escape analysis and ownership or other proven approaches.
4. Improve JSON typed access, HTTPS headers/POST, WebSocket and authenticated bot APIs; never embed secrets in source.
5. Add custom windows, layouts, events, drawing and audio.
6. Design a typed intermediate representation, optimizer and reproducible cross-language benchmarks; add ARM64/WebAssembly where practical.
7. Grow Vault into a trustworthy package ecosystem and eventually self-host the compiler in TRIPOL.

Correctness, simplicity, reproducibility and free access come before performance marketing.
