package main

import "github.com/goforj/execx"

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// WithPrefix sets the shadow print prefix.

	// Example: shadow prefix
	_, _ = execx.Command("printf", "hi").ShadowPrint(execx.WithPrefix("run")).Run()
	// run > printf hi
	// run > printf hi (1ms)
}
