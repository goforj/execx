//go:build ignore
// +build ignore

package main

import "github.com/goforj/execx"

func main() {
	// ShadowPrintOff disables shadow printing for this command chain.

	// Example: shadow print off
	_, _ = execx.Command("printf", "hi").ShadowPrint().ShadowPrintOff().Run()
}
