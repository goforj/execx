//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// PipelineResults executes the command and returns per-stage results.

	// Example: pipeline results
	results := execx.Command("printf", "go").
		Pipe("tr", "a-z", "A-Z").
		PipelineResults()
	fmt.Println(len(results))
	// #int 2
}
