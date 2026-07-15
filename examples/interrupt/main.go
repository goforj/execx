package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Interrupt sends an interrupt signal to the process.

	// Example: interrupt
	proc := execx.Command("sleep", "2").Start()
	_ = proc.Interrupt()
	res, _ := proc.Wait()
	fmt.Printf("%+v", res)
	// {Stdout: Stderr: ExitCode:-1 Err:<nil> Duration:75.987ms signal:interrupt}
}
