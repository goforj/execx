package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Output executes the command and returns stdout and any error.

	// Example: output
	out, _ := execx.Command("printf", "hello").Output()
	fmt.Print(out)
	// hello
}
