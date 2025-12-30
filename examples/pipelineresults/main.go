//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// PipelineResults executes the command and returns per-stage results and any error.

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
