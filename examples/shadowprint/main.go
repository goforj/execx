//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// ShadowPrint configures shadow printing for this command chain.

	// Example: shadow print
	_, _ = execx.Command("bash", "-c", `echo "hello world"`).
		ShadowPrint().
		OnStdout(func(line string) { fmt.Println(line) }).
		Run()
	// execx > bash -c 'echo "hello world"'
	// hello world
	// execx > bash -c 'echo "hello world"' (1ms)
}
