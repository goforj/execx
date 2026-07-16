package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// ShellEscaped returns a shell-escaped string for logging only.

	// Example: shell escaped
	cmd := execx.Command("echo", "hello world", "it's")
	fmt.Println(cmd.ShellEscaped())
	// #string echo 'hello world' "it's"
}
