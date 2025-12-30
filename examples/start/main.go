//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Start executes the command asynchronously.

	// Example: start
	proc := execx.Command("go", "env", "GOOS").Start()
	res, _ := proc.Wait()
	fmt.Println(res.ExitCode == 0)
	// #bool true
}
