package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// EnvInherit restores default environment inheritance.

	// Example: inherit env
	cmd := execx.Command("go", "env", "GOOS").EnvInherit()
	fmt.Println(len(cmd.EnvList()) > 0)
	// #bool true
}
