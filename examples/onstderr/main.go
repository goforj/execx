package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// OnStderr registers a line callback for stderr. The final unterminated line is
	// delivered after the process exits. Output callbacks and writers are serialized
	// across one pipeline.

	// Example: stderr lines
	_, err := execx.Command("go", "env", "-badflag").
		OnStderr(func(line string) {
			fmt.Println(line)
		}).
		Run()
	fmt.Println(err == nil)
	// flag provided but not defined: -badflag
	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
	// Run 'go help env' for details.
	// false
}
