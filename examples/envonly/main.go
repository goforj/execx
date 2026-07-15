package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// EnvOnly ignores the parent environment.

	// Example: replace env
	cmd := execx.Command("go", "env", "GOOS").EnvOnly(map[string]string{"A": "1"})
	fmt.Println(strings.Join(cmd.EnvList(), ","))
	// #string A=1
}
