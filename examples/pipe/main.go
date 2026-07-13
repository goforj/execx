package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Pipe appends a new command to the pipeline. Pipelines run on all platforms.

	// Example: pipe
	out, _ := execx.Command("printf", "go").
		Pipe("tr", "a-z", "A-Z").
		OutputTrimmed()
	fmt.Println(out)
	// #string GO
}
