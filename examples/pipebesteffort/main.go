package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// PipeBestEffort selects the last stage as the primary result while surfacing the
	// first execution error. Every stage is started concurrently.

	// Example: best effort
	res, _ := execx.Command("false").
		Pipe("printf", "ok").
		PipeBestEffort().
		Run()
	fmt.Print(res.Stdout)
	// ok
}
