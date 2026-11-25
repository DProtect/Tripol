# Production compiler

The Go implementation is packaged here. The root `main.go` calls `core.Run()`, preserving `go build .` and `go run .` from the repository root. White-box regression tests stay beside the compiler so they can access private implementation details. Run `go test ./core ./tests` from the root.

`../compiler/` is GENESIS, the experimental compiler written in TRIPOL. `../stdlib/` supplies embedded libraries; `../native/` contains the assembly source and regeneration tool.
