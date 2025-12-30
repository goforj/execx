//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// HideWindow is a no-op on non-Windows platforms; on Windows it hides console windows.

	// Example: hide window
	out, _ := execx.Command("printf", "ok").HideWindow(true).Output()
	fmt.Print(out)
	// ok
}
