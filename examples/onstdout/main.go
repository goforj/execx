//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// OnStdout registers a line callback for stdout.

	// Example: stdout lines
	_, _ = execx.Command("printf", "hi\n").
		OnStdout(func(line string) { fmt.Println(line) }).
		Run()
	// hi
}
