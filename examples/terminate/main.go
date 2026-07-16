package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Terminate kills the process immediately.

	// Example: terminate
	proc := execx.Command("sleep", "2").Start()
	_ = proc.Terminate()
	res, _ := proc.Wait()
	fmt.Printf("%+v", res)
	// {Stdout: Stderr: ExitCode:-1 Err:<nil> Duration:70.654ms signal:killed}
}
