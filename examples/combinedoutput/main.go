//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// CombinedOutput executes the command and returns stdout+stderr and any error.

	// Example: combined output
	out, _ := execx.Command("go", "env", "GOOS").CombinedOutput()
	fmt.Println(out != "")
	// #bool true
}
