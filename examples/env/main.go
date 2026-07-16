package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Env adds environment variables to the command.

	// Example: set env
	cmd := execx.Command("go", "env", "GOOS").Env("MODE=prod")
	fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "MODE=prod"))
	// #bool true
}
