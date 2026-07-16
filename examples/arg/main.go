package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Arg appends arguments to the command.

	// Example: add args
	cmd := execx.Command("printf").Arg("hello")
	out, _ := cmd.Output()
	fmt.Print(out)
	// hello
}
