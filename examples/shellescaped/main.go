//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// ShellEscaped returns a shell-escaped string for logging only.

	// Example: shell escaped
	cmd := execx.Command("echo", "hello world", "it's")
	fmt.Println(cmd.ShellEscaped())
	// #string echo 'hello world' "it's"
}
