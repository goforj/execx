//go:build ignore
// +build ignore

package main

import "github.com/goforj/execx"

func main() {
	// ShadowPrint writes the shell-escaped command to stderr before and after execution.

	// Example: shadow print
	_, _ = execx.Command("printf", "hi").ShadowPrint().Run()
	// execx > printf hi
	// execx > printf hi (1ms)
}
