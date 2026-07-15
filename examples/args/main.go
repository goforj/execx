package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Args returns the argv slice used for execution.

	// Example: args
	cmd := execx.Command("go", "env", "GOOS")
	fmt.Println(strings.Join(cmd.Args(), " "))
	// #string go env GOOS
}
