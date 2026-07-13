package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// String returns a human-readable representation of the command.

	// Example: string
	cmd := execx.Command("echo", "hello world", "it's")
	fmt.Println(cmd.String())
	// #string echo "hello world" it's
}
