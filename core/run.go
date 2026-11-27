package core

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

type discordStudioOptions struct{ metricsFile string }

func runNative(p string, studio *discordStudioOptions) {
	c := exec.Command(p)
	c.Stdout = os.Stdout
	var stderr bytes.Buffer
	if diagnosticJSON {
		c.Stderr = &stderr
	} else {
		c.Stderr = os.Stderr
	}
	c.Stdin = os.Stdin
	if e := runNativeCommand(c, studio); e != nil {
		var x *exec.ExitError
		if errors.As(e, &x) {
			d := runtimeDiagnostic(x.ExitCode())
			d.Detail = strings.TrimSpace(stderr.String())
			panic(d)
		}
		die("%v", e)
	}
}
