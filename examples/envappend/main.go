package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// EnvAppend merges variables into the inherited environment.

	// Example: append env
	cmd := execx.Command("go", "env", "GOOS").EnvAppend(map[string]string{"A": "1"})
	fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "A=1"))
	// #bool true
}
