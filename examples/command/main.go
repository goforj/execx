//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Command constructs a new command without executing it.

	// Example: command
	cmd := execx.Command("printf", "hello")
	out, _ := cmd.Output()
	fmt.Print(out)
	// hello
}
