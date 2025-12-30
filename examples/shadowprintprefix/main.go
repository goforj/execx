//go:build ignore
// +build ignore

package main

import "github.com/goforj/execx"

func main() {
	// ShadowPrintPrefix sets the prefix used by ShadowPrint.

	// Example: shadow print prefix
	_, _ = execx.Command("printf", "hi").ShadowPrintPrefix("run").Run()
	// run > printf hi
	// run > printf hi (1ms)
}
