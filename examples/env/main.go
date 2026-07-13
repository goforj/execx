package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// Env adds environment variables to the command.

	// Example: set env
	cmd := execx.Command("go", "env", "GOOS").Env("MODE=prod")
	fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "MODE=prod"))
	// #bool true
}
