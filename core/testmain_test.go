package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Preserve the repository-relative fixture paths used by the original suite.
func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	root := cwd
	if _, err := os.Stat(filepath.Join(root, "compiler", "main.tri")); err != nil {
		root = filepath.Dir(cwd)
	}
	if _, err := os.Stat(filepath.Join(root, "compiler", "main.tri")); err != nil {
		fmt.Fprintln(os.Stderr, "cannot locate TRIPOL repository fixtures:", err)
		os.Exit(1)
	}
	if err := os.Chdir(root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
