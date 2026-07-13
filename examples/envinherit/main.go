package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// EnvInherit restores default environment inheritance.

	// Example: inherit env
	cmd := execx.Command("go", "env", "GOOS").EnvInherit()
	fmt.Println(len(cmd.EnvList()) > 0)
	// #bool true
}
