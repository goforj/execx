//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// WithMask applies a masker to the shadow-printed command string.

	// Example: shadow mask
	mask := func(cmd string) string {
		return strings.ReplaceAll(cmd, "secret", "***")
	}
	_, _ = execx.Command("printf", "secret").ShadowPrint(execx.WithMask(mask)).Run()
	// execx > printf ***
	// execx > printf *** (1ms)
}
