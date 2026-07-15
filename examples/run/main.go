package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Run executes the command and returns the result and any execution error. A
	// non-zero exit code is represented in Result and does not itself produce an error.

	// Example: run
	res, _ := execx.Command("go", "env", "GOOS").Run()
	fmt.Println(res.ExitCode == 0)
	// #bool true
}
