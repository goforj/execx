//go:build ignore
// +build ignore

package main

import "github.com/goforj/execx"

func main() {
	// ShadowOn enables shadow printing using the previously configured options.

	// Example: shadow on
	cmd := execx.Command("printf", "hi").
		ShadowPrint(execx.WithPrefix("run"))
	cmd.ShadowOff()
	_, _ = cmd.ShadowOn().Run()
	// run > printf hi
	// run > printf hi (1ms)
}
