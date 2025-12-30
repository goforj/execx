//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Wait waits for the command to complete and returns the result and any error.

	// Example: wait
	proc := execx.Command("go", "env", "GOOS").Start()
	res, err := proc.Wait()
	fmt.Println(err == nil && res.ExitCode == 0)
	// #bool true
}
