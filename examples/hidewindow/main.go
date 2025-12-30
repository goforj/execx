//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// HideWindow hides console windows and sets CREATE_NO_WINDOW for console apps.

	// Example: hide window
	out, _ := execx.Command("printf", "ok").HideWindow(true).Output()
	fmt.Print(out)
	// ok
}
