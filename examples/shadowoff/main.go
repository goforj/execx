package main

import "github.com/goforj/execx"

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// ShadowOff disables shadow printing for this command chain, preserving configuration.

	// Example: shadow off
	_, _ = execx.Command("printf", "hi").ShadowPrint().ShadowOff().Run()
}
