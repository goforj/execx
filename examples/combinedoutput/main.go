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
	out, err := execx.Command("go", "env", "-badflag").CombinedOutput()
	fmt.Print(out)
	fmt.Println(err == nil)
	// flag provided but not defined: -badflag
	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
	// Run 'go help env' for details.
	// false
}
