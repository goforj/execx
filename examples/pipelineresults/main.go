package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// PipelineResults executes the command and returns per-stage results and the first
	// execution error. Non-zero exit codes remain data in their corresponding Result.

	// Example: pipeline results
	results, _ := execx.Command("printf", "go").
		Pipe("tr", "a-z", "A-Z").
		PipelineResults()
	fmt.Printf("%+v", results)
	// [
	//	{Stdout:go Stderr: ExitCode:0 Err:<nil> Duration:6.367208ms signal:<nil>}
	//	{Stdout:GO Stderr: ExitCode:0 Err:<nil> Duration:4.976291ms signal:<nil>}
	// ]
}
