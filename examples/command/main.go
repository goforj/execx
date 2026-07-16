package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Command constructs a new command without executing it.

	// Example: command
	cmd := execx.Command("printf", "hello")
	out, _ := cmd.Output()
	fmt.Print(out)
	// hello
}
