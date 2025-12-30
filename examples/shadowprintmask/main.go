//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// ShadowPrintMask sets a command masker for ShadowPrint output.

	// Example: shadow print mask
	mask := func(cmd string) string {
		return strings.ReplaceAll(cmd, "secret", "***")
	}
	_, _ = execx.Command("printf", "secret").ShadowPrintMask(mask).Run()
	// execx > printf ***
	// execx > printf *** (1ms)
}
