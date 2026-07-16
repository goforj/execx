package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// PipeStrict selects the first stage with an execution error or non-zero exit as
	// the primary result. Every stage is still started concurrently.

	// Example: strict
	res, _ := execx.Command("false").
		Pipe("printf", "ok").
		PipeStrict().
		Run()
	fmt.Println(res.ExitCode != 0)
	// #bool true
}
