//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Run executes the command and returns the result and any error.

	// Example: run
	res, err := execx.Command("go", "env", "GOOS").Run()
	fmt.Println(err == nil && res.ExitCode == 0)
	// #bool true
}
