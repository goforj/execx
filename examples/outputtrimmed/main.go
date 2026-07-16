package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// OutputTrimmed executes the command and returns trimmed stdout and any error.

	// Example: output trimmed
	out, _ := execx.Command("printf", "hello\n").OutputTrimmed()
	fmt.Println(out)
	// #string hello
}
