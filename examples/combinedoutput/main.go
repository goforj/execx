package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// CombinedOutput executes the command and returns stdout and stderr in observed
	// chunk order plus any execution error. Exact byte interleaving between the two
	// operating-system streams is inherently scheduler-dependent.

	// Example: combined output
	out, err := execx.Command("go", "env", "-badflag").CombinedOutput()
	fmt.Print(out)
	fmt.Println(err == nil)
	// flag provided but not defined: -badflag
	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
	// Run 'go help env' for details.
	// false
}
