# Roadmap

ORIGIN is a public-preview candidate, not stable v1.0. Goals, not promised versions or dates.

- Public release gate: sync complete source, audit history and credentials, review Apache-2.0 boundaries, Windows/Linux CI, verify downloadable assets.
- Gateway reliability: reconnect/resume, rate limits, slash interactions, Unicode and bounded queue; preserve simple API.
- Safe memory reclamation, stronger ownership/reference analysis and dynamic collections.
- Full typed control-flow TIR, register allocation, checked optimizations, broader reproducible workloads.
- Graphics, audio, crypto via reviewed APIs, database and packaging.
- Stable specification, compatibility policy, release candidates and reproducible builds before v1.0.

## GENESIS self-hosting

R1: native ELF emission from one constant integer statement. R2: runtime expression emission and binary output; R3: statements and typed functions; R4: bootstrap compiler capable of compiling itself, with multi-generation equivalence tests. Preserve the Go compiler until parity is verified.
