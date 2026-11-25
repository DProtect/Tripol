# Integration tests

Run from the repository root: `go test ./core ./tests`.

The Go compiler's white-box tests remain inside `core/` to retain access to unexported internals. This directory contains the public CLI integration suite. `history/test-results.txt` is an archived historical run, not the latest CI result. Current CI evidence is in GitHub Actions.
