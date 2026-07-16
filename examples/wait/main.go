package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Wait waits for the command to complete and returns the result and any error.

	// Example: wait
	proc := execx.Command("go", "env", "GOOS").Start()
	res, _ := proc.Wait()
	fmt.Printf("%+v", res)
	// {Stdout:darwin
	// Stderr: ExitCode:0 Err:<nil> Duration:1.234ms signal:<nil>}
}
