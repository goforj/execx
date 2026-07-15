package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// OnStdout registers a line callback for stdout. The final unterminated line is
	// delivered after the process exits. Output callbacks and writers are serialized
	// across one pipeline.

	// Example: stdout lines
	_, _ = execx.Command("printf", "hi\n").
		OnStdout(func(line string) { fmt.Println(line) }).
		Run()
	// hi
}
