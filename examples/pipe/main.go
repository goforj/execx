package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Pipe appends a new command to the pipeline. Pipelines run on all platforms.
	// Configuration called on the returned Cmd applies to the new stage; configure
	// the current stage before Pipe when it needs different environment, context, or I/O.

	// Example: pipe
	out, _ := execx.Command("printf", "go").
		Pipe("tr", "a-z", "A-Z").
		OutputTrimmed()
	fmt.Println(out)
	// #string GO
}
