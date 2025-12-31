//go:build ignore
// +build ignore

package main

import "github.com/goforj/execx"

func main() {
	// ShadowOff disables shadow printing for this command chain, preserving configuration.

	// Example: shadow off
	_, _ = execx.Command("printf", "hi").ShadowPrint().ShadowOff().Run()
}
