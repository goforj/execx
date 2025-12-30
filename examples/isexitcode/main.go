//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// IsExitCode reports whether the exit code matches.

	// Example: exit code
	res, _ := execx.Command("go", "env", "GOOS").Run()
	fmt.Println(res.IsExitCode(0))
	// #bool true
}
