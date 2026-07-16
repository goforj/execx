package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// EnvAppend merges variables into the inherited environment.

	// Example: append env
	cmd := execx.Command("go", "env", "GOOS").EnvAppend(map[string]string{"A": "1"})
	fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "A=1"))
	// #bool true
}
