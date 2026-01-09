//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/goforj/execx"
)

func main() {
	// WithPTY attaches stdout/stderr to a pseudo-terminal.
	// Output is combined; OnStdout and OnStderr receive the same lines.

	// Example: with pty
	_, _ = execx.Command("printf", "hi\n").
		WithPTY().
		OnStdout(func(line string) { fmt.Println(line) }).
		Run()
	// hi
}
