package main

import (
	"errors"
	"os"
	"os/exec"
)

func runNative(p string) {
	c := exec.Command(p)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if e := c.Run(); e != nil {
		var x *exec.ExitError
		if errors.As(e, &x) {
			die("program exited with status %d", x.ExitCode())
		}
		die("%v", e)
	}
}
