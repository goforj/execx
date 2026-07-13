package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// WithPTY attaches stdout/stderr to a pseudo-terminal.
	//
	// When enabled, stdout and stderr are merged into a single stream. OnStdout and
	// OnStderr both receive the same lines, and Result.Stderr remains empty.
	// Platforms without PTY support return an error when the command runs.

	// Example: with pty
	_, _ = execx.Command("printf", "hi").
		WithPTY().
		OnStdout(func(line string) { fmt.Println(line) }).
		Run()
	// hi
}
