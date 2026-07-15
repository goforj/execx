package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Start executes the command asynchronously. Startup still completes synchronously,
	// so a returned Process represents either a fully started pipeline or a completed error.

	// Example: start
	proc := execx.Command("go", "env", "GOOS").Start()
	res, _ := proc.Wait()
	fmt.Println(res.ExitCode == 0)
	// #bool true
}
