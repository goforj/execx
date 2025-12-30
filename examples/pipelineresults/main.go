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
	results, err := execx.Command("printf", "go").
		Pipe("tr", "a-z", "A-Z").
		PipelineResults()
	fmt.Println(err == nil && len(results) == 2)
	// #bool true
}
