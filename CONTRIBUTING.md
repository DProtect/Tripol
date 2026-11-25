# Contributing

Build from [source](docs/build.md), write a minimal failing `.tri` example for bugs, and run `go test -count=1 -v ./core ./tests` and `go vet . ./core ./stdlib ./tests`. Add tests for both valid and invalid programs. Keep native code ABI assumptions and platform limitations explicit. Do not commit generated installers, credentials, or local caches. Regenerate native templates only when changing assembly. Benchmark with identical algorithms and outputs, warm-up and repeated runs; report hardware/OS and avoid extrapolating. Describe Windows tests honestly (cross-compilation is not execution). All contributions are subject to the Apache-2.0 license and repository access policy.

The compiler implementation and white-box tests are in `core/`. Run `go test ./core ./tests` and `go vet . ./core ./stdlib ./tests` from the repository root. The existing `go build -o tripol .` entry point remains supported.
