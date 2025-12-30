//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// CreationFlags is a no-op on non-Windows platforms; on Windows it sets process creation flags.

	// Example: creation flags
	out, _ := execx.Command("printf", "ok").CreationFlags(execx.CreateNewProcessGroup).Output()
	fmt.Print(out)
	// ok
}
