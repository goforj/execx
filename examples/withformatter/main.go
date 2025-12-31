//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// WithFormatter sets a formatter for ShadowPrint output.

	// Example: shadow formatter
	formatter := func(ev execx.ShadowEvent) string {
		return fmt.Sprintf("shadow: %s %s", ev.Phase, ev.Command)
	}
	_, _ = execx.Command("printf", "hi").ShadowPrint(execx.WithFormatter(formatter)).Run()
	// shadow: before printf hi
	// shadow: after printf hi
}
