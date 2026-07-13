package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// Args returns the argv slice used for execution.

	// Example: args
	cmd := execx.Command("go", "env", "GOOS")
	fmt.Println(strings.Join(cmd.Args(), " "))
	// #string go env GOOS
}
